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
