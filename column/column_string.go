package column

import (
	"io"

	"gotomerge/column/rle"
	ioutil "gotomerge/utils/io"
)

type StringReader = rle.StringReader

func NewStringReader(r ioutil.SubReader) *StringReader {
	return rle.NewStringReader(r)
}

type StringWriter = rle.Writer[string]

func NewStringWriter(w io.Writer) *StringWriter {
	return rle.NewStringWriter(w)
}
