package column

import (
	"io"
	"iter"

	"gotomerge/column/rle"
)

type ActorColumnIter = iter.Seq2[rle.NullableValue[uint64], error]

func ReadActorColumn(r io.Reader) ActorColumnIter {
	return rle.ReadUint64RLE(r)
}

//
// type ActorColumnIter struct {
// 	next func() (rle.NullableValue[uint64], error)
// 	stop func()
// }
//
// func ReadActorColumn(r io.Reader) *ActorColumnIter {
// 	next, stop := iter.Pull2(rle.ReadUint64RLE(r))
// 	return &ActorColumnIter{
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
// 		stop: stop,
// 	}
// }
//
// func (a *ActorColumnIter) Next() (rle.NullableValue[uint64], error) {
// 	if a == nil {
// 		return rle.NewNullUint64(), nil
// 	}
// 	return a.next()
// }
//
// func (a *ActorColumnIter) Stop() {
// 	if a == nil {
// 		return
// 	}
// 	a.stop()
// }
