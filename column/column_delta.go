package column

import (
	"bytes"
	"fmt"
	"iter"
	"math"

	"gotomerge/column/rle"
)

type DeltaColumn []byte

type DeltaColumnIter = iter.Seq2[rle.NullableValue[uint64], error]

func DeltaColumnFromBytes(b []byte) DeltaColumn {
	return DeltaColumn(b)
}

func (d DeltaColumn) Iter() DeltaColumnIter {
	return func(yield func(rle.NullableValue[uint64], error) bool) {
		var prev, res uint64
		for signed, err := range rle.ReadInt64RLE(bytes.NewReader(d)) {
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
