package opset

import (
	"github.com/MichaelMure/gotomerge/types"
	"github.com/MichaelMure/gotomerge/utils/treap"
)

// listTreap is a lazily built, cached ordered view of the live elements of a
// list or text object. It stores only live (non-tombstone) ops in CRDT
// traversal order, so callers can iterate without liveness checks.
//
// The treap is invalidated (removed from OpSet.listTreaps) whenever any op
// for the object arrives, ensuring it is never stale. The byId map lets the
// transaction seed its working copy with an O(n) treap scan instead of an
// O(n) map-and-DFS reconstruction.
type listTreap struct {
	r    *treap.Treap[Op]
	byId map[types.OpId]*treap.Node[Op]
}

// getOrBuildListTreap returns the cached listTreap for obj, building and
// caching it on first access.
func (s *OpSet) getOrBuildListTreap(obj types.ObjectId) *listTreap {
	if lt, ok := s.listTreaps[obj]; ok {
		return lt
	}
	lt := s.buildListTreap(obj)
	s.listTreaps[obj] = lt
	return lt
}

// invalidateListTreap drops the cached treap for obj. Called by ApplyChange
// and applyCommitted whenever ops arrive for that object.
func (s *OpSet) invalidateListTreap(obj types.ObjectId) {
	delete(s.listTreaps, obj)
}

// buildListTreap runs the same scan-and-DFS as ListElements but populates a
// treap instead of a slice. The scan is shared via iterListElements.
func (s *OpSet) buildListTreap(obj types.ObjectId) *listTreap {
	lt := &listTreap{
		r:    treap.New[Op](),
		byId: make(map[types.OpId]*treap.Node[Op]),
	}
	s.iterListElements(obj, func(op Op) {
		lt.byId[op.Id] = lt.r.PushBack(op)
	})
	return lt
}

// iterListElements runs the CRDT list reconstruction for obj and calls fn for
// each live element in order. It is the shared core used by both ListElements
// and buildListTreap so the algorithm is not duplicated.
func (s *OpSet) iterListElements(obj types.ObjectId, fn func(Op)) {
	// children maps each predecessor key to the insert ops immediately after it.
	children := make(map[types.Key][]types.OpId)
	// liveValues maps a list position (OpId of its Insert=true op) to the
	// winning live value op at that position.
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
	var traverse func(key types.Key)
	traverse = func(key types.Key) {
		for _, id := range children[key] {
			if op, ok := liveValues[id]; ok {
				fn(op)
			}
			traverse(types.KeyOpId(id))
		}
	}
	traverse(types.KeyOpId{ActorIdx: 0, Counter: 0})
}
