package column

import (
	"io"
	"iter"

	"gotomerge/column/rle"
)

// type UlebColumnIter = iter.Seq2[rle.NullableValue[uint64], error]

type UlebColumnIter = iter.Seq2[rle.NullableValue[uint64], error]

func ReadUlebColumn(r io.Reader) UlebColumnIter {
	return rle.ReadUint64RLE(r)
}

// type UlebColumnIter struct {
// 	next func() (rle.NullableValue[uint64], error)
// 	stop func()
// }
//
// func ReadUlebColumn(r io.Reader) *UlebColumnIter {
// 	next, stop := iter.Pull2(rle.ReadUint64RLE(r))
// 	return &UlebColumnIter{
// 		next: func() (rle.NullableValue[uint64], error) {
// 			val, err, ok := next()
// 			if !ok {
// 				return rle.NewNullUint64(), nil
// 			}
// 			if err != nil {
// 				return rle.NewNullUint64(), err
// 			}
// 			return val, err
// 		},
// 		stop: stop,
// 	}
// }
//
// func (a *UlebColumnIter) Next() (rle.NullableValue[uint64], error) {
// 	if a == nil {
// 		return rle.NewNullUint64(), nil
// 	}
// 	return a.next()
// }
//
// func (a *UlebColumnIter) Stop() {
// 	if a == nil {
// 		return
// 	}
// 	a.stop()
// }
