package column

import (
	"gotomerge/column/rle"
	ioutil "gotomerge/utils/io"
)

type StringReader = rle.StringReader

func NewStringReader(r ioutil.SubReader) *StringReader {
	return rle.NewStringReader(r)
}
