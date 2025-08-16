package column

import (
	"iter"

	"gotomerge/column/rle"
	"gotomerge/lbuf"
)

type UlebColumn []byte

type UlebColumnIter = iter.Seq2[rle.NullableValue[uint64], error]

func UlebColumnFromBytes(b []byte) UlebColumn {
	return UlebColumn(b)
}

func (uc UlebColumn) Iter() UlebColumnIter {
	return rle.ReadUint64RLE(lbuf.FromBytes(uc))
}
