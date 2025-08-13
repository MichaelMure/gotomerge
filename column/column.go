package column

import (
	"fmt"
	"io"
	"iter"
	"math"

	"github.com/jcalabro/leb128"

	"gotomerge/column/rle"
	"gotomerge/lbuf"
)

func ReadUlebColumn(r *lbuf.Reader) iter.Seq2[rle.NullableValue[uint64], error] {
	return rle.ReadUint64RLE(r)
}

// ReadDeltaColumn reads the values of a delta column of length bytes
func ReadDeltaColumn(r *lbuf.Reader) iter.Seq2[rle.NullableValue[uint64], error] {
	return func(yield func(rle.NullableValue[uint64], error) bool) {
		var prev, res uint64
		for signed, err := range rle.ReadInt64RLE(r) {
			if err != nil {
				yield(nil, err)
				return
			}
			val, valid := signed.Value()

			switch {
			case !valid:
				// automerge in rust consider null values as "stay the same" and yield null
				if !yield(rle.NewNullUint64(), nil) {
					return
				}
				continue
			case val > 0:
				if prev > math.MaxUint64-uint64(val) {
					yield(nil, fmt.Errorf("overflow in delta column"))
					return
				}
				res = prev + uint64(val)
			case val < 0:
				if prev < uint64(-val) {
					yield(nil, fmt.Errorf("underflow in delta column"))
					return
				}
				res = prev - uint64(-val)
			}
			prev = res
			if !yield(rle.NewNullableUint64(res), nil) {
				return
			}
		}
	}
}

// ReadStringColumn reads the values of a string column of length bytes
func ReadStringColumn(r *lbuf.Reader) iter.Seq2[rle.NullableValue[string], error] {
	return rle.ReadStringRLE(r)
}

func ReadBooleanColumn(r *lbuf.Reader) iter.Seq2[bool, error] {
	return func(yield func(bool, error) bool) {
		var val bool
		for {
			count, err := leb128.DecodeU64(r)
			if err == io.EOF {
				return
			}
			if err != nil {
				yield(false, err)
				return
			}
			for i := uint64(0); i < count; i++ {
				if !yield(val, nil) {
					return
				}
			}
			val = !val
		}
	}
}
