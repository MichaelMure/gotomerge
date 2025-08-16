package column

import (
	"bytes"
	"iter"

	"gotomerge/column/rle"
)

type GroupColumn []byte

type GroupColumnIter = iter.Seq2[rle.NullableValue[uint64], error]

func GroupColumnFromBytes(b []byte) GroupColumn {
	return GroupColumn(b)
}

func (gc GroupColumn) Iter() GroupColumnIter {
	return rle.ReadUint64RLE(bytes.NewReader(gc))
}
