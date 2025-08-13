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

func ReadUlebColumn(r *lbuf.Reader) iter.Seq2[rle.NullableUint64, error] {
	return rle.ReadUint64RLE(r)
}

// ReadDeltaColumn reads the values of a delta column of length bytes
func ReadDeltaColumn(r *lbuf.Reader) iter.Seq2[uint64, error] {
	return func(yield func(uint64, error) bool) {
		var prev, res uint64
		for signed, err := range rle.ReadInt64RLE(r) {
			if err != nil {
				yield(0, err)
				return
			}
			val, valid := signed.Value()

			// automerge in rust consider null values as "stay the same"
			if !valid {
				if !yield(res, nil) {
					return
				}
			}

			switch {
			case val > 0:
				if prev > math.MaxUint64-uint64(val) {
					yield(0, fmt.Errorf("overflow in delta column"))
					return
				}
				res = prev + uint64(val)
			case val < 0:
				if prev < uint64(-val) {
					yield(0, fmt.Errorf("underflow in delta column"))
					return
				}
				res = prev - uint64(-val)
			}
			prev = res
			if !yield(res, nil) {
				return
			}
		}
	}
}

// ReadStringColumn reads the values of a string column of length bytes
func ReadStringColumn(r *lbuf.Reader) iter.Seq2[rle.NullableString, error] {
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
