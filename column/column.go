package column

import (
	"fmt"
	"io"
	"iter"

	"github.com/jcalabro/leb128"

	"gotomerge/column/rle"
)

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
				res = prev + val
				if (res > prev) == (val < 0) {
					yield(0, fmt.Errorf("overflow or underflow in delta column"))
					return
				}
				prev = res
			}

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
