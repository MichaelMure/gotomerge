package column

import (
	"fmt"
	"io"
	"iter"

	"gotomerge/column/rle"
)

// ReadDeltaColumn reads the values of a delta column of length bytes
// TODO: is that returning int64 or uint64 ?
func ReadDeltaColumn(r io.Reader, length uint64) iter.Seq2[int64, error] {
	r = io.LimitReader(r, int64(length))
	return func(yield func(int64, error) bool) {
		var prev int64
		for signed, err := range rle.ReadInt64RLE(r) {
			if err != nil {
				yield(0, err)
				return
			}
			val, valid := signed.Value()
			if !valid {
				yield(0, fmt.Errorf("null value in delta column"))
				return
			}

			res := prev + val
			if (res > prev) == (val < 0) {
				yield(0, fmt.Errorf("overflow or underflow in delta column"))
				return
			}
			prev = res
			if !yield(res, nil) {
				return
			}
		}
	}
}

// ReadStringColumn reads the values of a string column of length bytes
func ReadStringColumn(r io.Reader, length uint64) iter.Seq2[rle.NullableString, error] {
	return rle.ReadStringRLE(io.LimitReader(r, int64(length)))
}
