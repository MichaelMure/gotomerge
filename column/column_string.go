package column

import (
	"iter"

	"gotomerge/column/rle"
	"gotomerge/lbuf"
)

type StringColumn []byte

type StringColumnIter = iter.Seq2[rle.NullableValue[string], error]

func StringColumnFromBytes(b []byte) StringColumn {
	return StringColumn(b)
}

func (sc StringColumn) Iter() StringColumnIter {
	return rle.ReadStringRLE(lbuf.FromBytes(sc))
}
