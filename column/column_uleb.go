package column

import (
	"io"

	"github.com/MichaelMure/gotomerge/column/rle"
	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

type UlebReader = rle.Uint64Reader

// PeekUlebReader creates a reader over a snapshot of r. See PeekActorReader.
func PeekUlebReader(r ioutil.SubReader) *UlebReader {
	return rle.NewUint64Reader(r)
}

type UlebWriter = rle.Writer[uint64]

func NewUlebWriter(w io.Writer) *UlebWriter {
	return rle.NewUint64Writer(w)
}
