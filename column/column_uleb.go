package column

import (
	"io"
	"iter"

	"gotomerge/column/rle"
)

type UlebColumnIter = iter.Seq2[rle.NullableValue[uint64], error]

func ReadUlebColumn(r io.Reader) UlebColumnIter {
	return rle.ReadUint64RLE(r)
}
