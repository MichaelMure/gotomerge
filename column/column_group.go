package column

import (
	"gotomerge/column/rle"
	ioutil "gotomerge/utils/io"
)

type GroupReader = rle.Uint64Reader

func NewGroupReader(r ioutil.SubReader) *GroupReader {
	return rle.NewUint64Reader(r)
}
