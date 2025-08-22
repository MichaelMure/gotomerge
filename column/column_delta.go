package column

import (
	"fmt"
	"io"
	"iter"
	"math"

	"gotomerge/column/rle"
)

type DeltaColumnIter = iter.Seq2[rle.NullableValue[int64], error]

func ReadDeltaColumn(r io.Reader) DeltaColumnIter {
	return func(yield func(rle.NullableValue[int64], error) bool) {
		var prev, res int64
		for signed, err := range rle.ReadInt64RLE(r) {
			if err != nil {
				yield(nil, err)
				return
			}
			val, valid := signed.Value()

			switch {
			case !valid:
				// automerge in rust consider null values as "stay the same" and yield null
				if !yield(rle.NewNullInt64(), nil) {
					return
				}
				continue
			case val > 0:
				if prev > math.MaxInt64-int64(val) {
					yield(nil, fmt.Errorf("overflow in delta column"))
					return
				}
				res = prev + val
			case val < 0:
				if prev < math.MinInt64-int64(val) {
					yield(nil, fmt.Errorf("underflow in delta column"))
					return
				}
				res = prev - -val
			}
			prev = res
			if !yield(rle.NewNullableInt64(res), nil) {
				return
			}
		}
	}
}

// type DeltaColumnIter struct {
// 	next func() (rle.NullableValue[uint64], error)
// 	stop func()
// }
//
// func ReadDeltaColumn(r io.Reader) *DeltaColumnIter {
// 	next, stop := iter.Pull2(func(yield func(rle.NullableValue[uint64], error) bool) {
// 		var prev, res uint64
// 		for signed, err := range rle.ReadInt64RLE(r) {
// 			if err != nil {
// 				yield(nil, err)
// 				return
// 			}
// 			val, valid := signed.Value()
//
// 			switch {
// 			case !valid:
// 				// automerge in rust consider null values as "stay the same" and yield null
// 				if !yield(rle.NewNullUint64(), nil) {
// 					return
// 				}
// 				continue
// 			case val > 0:
// 				if prev > math.MaxUint64-uint64(val) {
// 					yield(nil, fmt.Errorf("overflow in delta column"))
// 					return
// 				}
// 				res = prev + uint64(val)
// 			case val < 0:
// 				if prev < uint64(-val) {
// 					yield(nil, fmt.Errorf("underflow in delta column"))
// 					return
// 				}
// 				res = prev - uint64(-val)
// 			}
// 			prev = res
// 			if !yield(rle.NewNullableUint64(res), nil) {
// 				return
// 			}
// 		}
// 	})
// 	return &DeltaColumnIter{
// 		next: func() (rle.NullableValue[uint64], error) {
// 			val, err, ok := next()
// 			if !ok {
// 				return rle.NewNullUint64(), ErrDone
// 			}
// 			if err != nil {
// 				return rle.NewNullUint64(), err
// 			}
// 			return val, err
// 		},
// 	}
// }
