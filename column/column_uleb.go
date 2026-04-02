package column

import (
	"io"

	"gotomerge/column/rle"
	ioutil "gotomerge/utils/io"
)

type UlebReader = rle.Uint64Reader

func NewUlebReader(r ioutil.SubReader) *UlebReader {
	return rle.NewUint64Reader(r)
}

type UlebWriter = rle.Writer[uint64]

func NewUlebWriter(w io.Writer) *UlebWriter {
	return rle.NewUint64Writer(w)
}
