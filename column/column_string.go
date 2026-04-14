package column

import (
	"io"

	"github.com/MichaelMure/gotomerge/column/rle"
	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

type StringReader = rle.StringReader

// PeekStringReader creates a reader over a snapshot of r. See PeekActorReader.
func PeekStringReader(r ioutil.SubReader) *StringReader {
	return rle.NewStringReader(r)
}

type StringWriter = rle.Writer[string]

func NewStringWriter(w io.Writer) *StringWriter {
	return rle.NewStringWriter(w)
}
