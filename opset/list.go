package opset

import (
	"strings"

	"github.com/MichaelMure/gotomerge/types"
)

// ListElements returns the live elements of a list or text object in order.
// Each returned Op carries the winning live value at that position.
//
// A position is live if at least one op at that position (the insert op or a
// later update) has SuccCount==0 and action != Delete. Dead positions
// (tombstones) are skipped in the result but still traversed, so elements
// inserted after a deleted element remain in the list.
func (s *OpSet) ListElements(obj types.ObjectId) []Op {
	// children maps each predecessor key to the insert ops immediately after it.
	children := make(map[types.Key][]types.OpId)
	// liveValues maps a list position (OpId of its Insert=true op) to the
	// winning live value op at that position. A position absent from this map
	// is deleted (tombstone).
	liveValues := make(map[types.OpId]Op)

	addInsert := func(op Op, succCount uint32) {
		children[op.Key] = append(children[op.Key], op.Id)
		if succCount == 0 && op.Action.Kind != types.ActionDelete {
			liveValues[op.Id] = op
		}
	}

	addUpdate := func(op Op, succCount uint32) {
		key, ok := op.Key.(types.KeyOpId)
		if !ok {
			return
		}
		target := types.OpId(key)
		if succCount == 0 && op.Action.Kind != types.ActionDelete {
			if cur, exists := liveValues[target]; !exists || s.opIdGreater(op.Id, cur.Id) {
				liveValues[target] = op
			}
		}
	}

	if s.snapshot != nil {
		if r, ok := s.snapshot.objRanges[obj]; ok {
			s.snapshot.scanRange(r, func(idx uint32, op Op) bool {
				if op.Insert {
					addInsert(op, s.snapshot.succCount[idx])
				} else {
					addUpdate(op, s.snapshot.succCount[idx])
				}
				return true
			})
		}
	}

	if s.delta != nil {
		for _, idx := range s.delta.byObj[obj] {
			op := s.delta.ops[idx]
			if op.Insert {
				addInsert(op, op.SuccCount)
			} else {
				addUpdate(op, op.SuccCount)
			}
		}
	}

	// Sort concurrent insertions after the same predecessor: higher OpId first.
	for key := range children {
		sortOpIdsDesc(children[key], s)
	}

	// DFS traversal of the insertion tree in list order.
	var result []Op
	var traverse func(key types.Key)
	traverse = func(key types.Key) {
		for _, id := range children[key] {
			if op, ok := liveValues[id]; ok {
				result = append(result, op)
			}
			traverse(types.KeyOpId(id))
		}
	}
	// The head-of-list sentinel is KeyOpId{0,0}: null actor index + counter 0.
	traverse(types.KeyOpId{ActorIdx: 0, Counter: 0})

	return result
}

// ListGet returns the element at the given zero-based index.
// Returns false if index is out of range.
func (s *OpSet) ListGet(obj types.ObjectId, index int) (Op, bool) {
	elems := s.ListElements(obj)
	if index < 0 || index >= len(elems) {
		return Op{}, false
	}
	return elems[index], true
}

// Text concatenates the string values of all live elements of a text object
// in order.
func (s *OpSet) Text(obj types.ObjectId) string {
	elems := s.ListElements(obj)
	var b strings.Builder
	for _, op := range elems {
		if v, ok := op.Action.Value.(string); ok {
			b.WriteString(v)
		}
	}
	return b.String()
}
