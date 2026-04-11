package opset

import (
	"sort"

	"github.com/MichaelMure/gotomerge/column"
)

const seekStride = 64

// seekPoint is a combined checkpoint for all snapshot query columns,
// positioned immediately before op opIdx is read.
type seekPoint struct {
	opIdx  uint32
	key    *column.KeyReader
	opId   *column.OpIdReader
	action *column.ActionReader
}

// seekIndex is a sorted list of combined column checkpoints, spaced roughly
// seekStride operations apart. The first entry always has opIdx == 0.
type seekIndex []seekPoint

// seek returns the checkpoint at or before targetOp and the number of ops to
// skip from that checkpoint to reach targetOp. The remaining count is at most
// seekStride.
func (idx seekIndex) seek(targetOp uint32) (seekPoint, uint32) {
	if len(idx) == 0 {
		return seekPoint{}, targetOp
	}
	i := sort.Search(len(idx), func(i int) bool { return idx[i].opIdx > targetOp }) - 1
	if i < 0 {
		i = 0
	}
	pt := idx[i]
	return pt, targetOp - pt.opIdx
}
