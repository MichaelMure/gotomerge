package column

import (
	"io"
	"iter"

	"gotomerge/column/rle"
)

type GroupColumnIter = iter.Seq2[rle.NullableValue[uint64], error]

func ReadGroupColumn(r io.Reader) GroupColumnIter {
	return rle.ReadUint64RLE(r)
}
