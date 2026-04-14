package column

import (
	"io"

	"github.com/MichaelMure/gotomerge/column/rle"
	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

type GroupReader = rle.Uint64Reader

// PeekGroupReader creates a reader over a snapshot of r. See PeekActorReader.
func PeekGroupReader(r ioutil.SubReader) *GroupReader {
	return rle.NewUint64Reader(r)
}

type GroupWriter = rle.Writer[uint64]

func NewGroupWriter(w io.Writer) *GroupWriter {
	return rle.NewUint64Writer(w)
}
