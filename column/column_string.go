package column

import (
	"io"
	"iter"

	"gotomerge/column/rle"
)

type StringColumnIter = iter.Seq2[rle.NullableValue[string], error]

func ReadStringColumn(r io.Reader) StringColumnIter {
	return rle.ReadStringRLE(r)
}
