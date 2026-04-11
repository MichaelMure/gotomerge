package column

import (
	"io"

	"github.com/MichaelMure/gotomerge/column/rle"
	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

type ActorReader = rle.Uint64Reader

func NewActorReader(r *ioutil.SubReader) *ActorReader {
	return rle.NewUint64Reader(r)
}

type ActorWriter = rle.Writer[uint64]

func NewActorWriter(w io.Writer) *ActorWriter {
	return rle.NewUint64Writer(w)
}
