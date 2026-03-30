package column

import (
	"gotomerge/column/rle"
	ioutil "gotomerge/utils/io"
)

type ActorReader = rle.Uint64Reader

func NewActorReader(r ioutil.SubReader) *ActorReader {
	return rle.NewUint64Reader(r)
}
