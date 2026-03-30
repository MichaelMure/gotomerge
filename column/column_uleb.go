package column

import (
	"gotomerge/column/rle"
	ioutil "gotomerge/utils/io"
)

type UlebReader = rle.Uint64Reader

func NewUlebReader(r ioutil.SubReader) *UlebReader {
	return rle.NewUint64Reader(r)
}
