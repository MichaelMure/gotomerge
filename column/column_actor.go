package column

import (
	"bytes"
	"iter"

	"gotomerge/column/rle"
)

type ActorColumn []byte

type ActorColumnIter = iter.Seq2[rle.NullableValue[uint64], error]

func ActorColumnFromBytes(b []byte) ActorColumn {
	return ActorColumn(b)
}

func (ac ActorColumn) Iter() ActorColumnIter {
	return rle.ReadUint64RLE(bytes.NewReader(ac))
}
