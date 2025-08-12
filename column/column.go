package column

import (
	"fmt"
	"io"
	"iter"
	"math"

	"github.com/jcalabro/leb128"

	"gotomerge/column/rle"
)

func ReadUlebColumn(r io.Reader, length uint64) iter.Seq2[rle.NullableUint64, error] {
	// TODO: does that really return null values?
	return rle.ReadUint64RLE(io.LimitReader(r, int64(length)))
}

// ReadDeltaColumn reads the values of a delta column of length bytes
// TODO: is that returning int64 or uint64 ?
func ReadDeltaColumn(r io.Reader, length uint64) iter.Seq2[int64, error] {
	return func(yield func(int64, error) bool) {
		r = io.LimitReader(r, int64(length))
		var prev, res int64
		for signed, err := range rle.ReadInt64RLE(r) {
			if err != nil {
				yield(0, err)
				return
			}
			val, valid := signed.Value()
			// automerge in rust consider null values as "stay the same"
			if valid {
				if val > 0 && prev > math.MaxInt64-val {
					yield(0, fmt.Errorf("overflow in delta column"))
					return
				}
				if val < 0 && prev < math.MinInt64-val {
					yield(0, fmt.Errorf("underflow in delta column"))
					return
				}
				res = prev + val
			}

			if !yield(res, nil) {
				return
			}
		}
	}
}

// ReadStringColumn reads the values of a string column of length bytes
func ReadStringColumn(r io.Reader, length uint64) iter.Seq2[rle.NullableString, error] {
	// TODO: does that really return null values?
	return rle.ReadStringRLE(io.LimitReader(r, int64(length)))
}

func ReadBooleanColumn(r io.Reader, length uint64) iter.Seq2[bool, error] {
	return func(yield func(bool, error) bool) {
		r = io.LimitReader(r, int64(length))
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
