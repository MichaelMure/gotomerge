package column

import (
	"io"

	"github.com/MichaelMure/gotomerge/column/rle"
	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

type ActorReader = rle.Uint64Reader

// PeekActorReader creates a reader over a snapshot of r. The original cursor
// is not advanced; r is passed by value so the reader owns an independent copy.
func PeekActorReader(r ioutil.SubReader) *ActorReader {
	return rle.NewUint64Reader(r)
}

type ActorWriter = rle.Writer[uint64]

func NewActorWriter(w io.Writer) *ActorWriter {
	return rle.NewUint64Writer(w)
}
