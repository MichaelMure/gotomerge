package opset

import (
	"fmt"

	"github.com/MichaelMure/gotomerge/types"
	"github.com/MichaelMure/gotomerge/utils/rope"
	"github.com/MichaelMure/gotomerge/utils/treap"
)

// listTreap is a lazily built, incrementally maintained ordered view of
// a list or text object. It stores ALL insert ops — both live and
// tombstoned — in CRDT traversal order using soft-deletion (MarkDead).
//
// Node values are the insert ops themselves (Insert=true). This keeps
// node.Value().Key equal to the CRDT predecessor key, which is required
// by the insertion algorithm.
//
// For list objects where an update op (non-insert, non-delete) supersedes
// the original insert at a position, the winning value is stored in
// liveValues and returned by ListElements. Text objects never have update
// ops so liveValues stays nil.
//
// For text objects, rope holds the rendered text as chunked character slices
// and is kept in sync with the treap. Text() reads from the rope in O(n/K)
// instead of the O(n) treap traversal.
//
// byId maps insert-op Id → treap node, enabling O(1) predecessor lookup
// followed by O(log n) InsertBefore / MarkDead.
type listTreap struct {
	r          *treap.Treap[Op]
	byId       map[types.OpId]*treap.Node[Op]
	liveValues map[types.OpId]Op // non-nil only when update ops exist
	rope       *rope.Rope        // non-nil for text objects only
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

// invalidateListTreap drops the cached treap for obj. Used as a fallback
// for operations that cannot be applied incrementally (e.g. update ops
// that change a position's winning value).
func (s *OpSet) invalidateListTreap(obj types.ObjectId) {
	delete(s.listTreaps, obj)
}

// buildListTreap runs the full CRDT reconstruction for obj and populates
// a treap with every insert op (alive and dead) in DFS order. For text
// objects it also builds the rope for fast subsequent reads.
func (s *OpSet) buildListTreap(obj types.ObjectId) *listTreap {
	lt := &listTreap{
		r:    treap.New[Op](),
		byId: make(map[types.OpId]*treap.Node[Op]),
	}
	if kind, ok := s.ObjType(obj); ok && kind == types.ActionMakeText {
		lt.rope = &rope.Rope{}
	}

	// children maps predecessor key → insert op ids at that position (DFS edges).
	children := make(map[types.Key][]types.OpId)
	// insertOps stores the original insert op for every position.
	insertOps := make(map[types.OpId]Op)
	// liveValues stores the winning (highest-Id live) op at each position.
	liveValues := make(map[types.OpId]Op)

	addInsert := func(op Op, succCount uint32) {
		children[op.Key] = append(children[op.Key], op.Id)
		insertOps[op.Id] = op
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
				op.Object = obj // scanRange does not set Object
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

	for key := range children {
		sortOpIdsDesc(children[key], s)
	}

	// DFS: add every insert op to the treap in CRDT order.
	var traverse func(key types.Key)
	traverse = func(key types.Key) {
		for _, id := range children[key] {
			node := lt.r.PushBack(insertOps[id])
			lt.byId[id] = node

			winner, alive := liveValues[id]
			if !alive {
				lt.r.MarkDead(node)
			} else {
				if winner.Id != id {
					// An update op is the winner: track it in liveValues.
					if lt.liveValues == nil {
						lt.liveValues = make(map[types.OpId]Op)
					}
					lt.liveValues[id] = winner
				}
				if lt.rope != nil {
					if ch, ok := insertOps[id].Action.Value.(string); ok {
						lt.rope.PushBack(ch)
					}
				}
			}

			traverse(types.KeyOpId(id))
		}
	}
	traverse(types.KeyOpId{ActorIdx: 0, Counter: 0})

	return lt
}

// insertOp adds a new insert op to the treap at the correct CRDT position.
// The op's predecessor must already be in the treap (or be the sentinel).
// The node starts alive; call markDead immediately after if needed.
// For text objects the rope is updated in sync.
func (lt *listTreap) insertOp(op Op, s *OpSet) {
	predKey, ok := op.Key.(types.KeyOpId)
	if !ok {
		panic(fmt.Sprintf("insert op has non-KeyOpId predecessor key: %T", op.Key))
	}
	predId := types.OpId(predKey)

	var predNode *treap.Node[Op]
	if predId != (types.OpId{}) {
		predNode = lt.byId[predId]
	}

	insertBefore := lt.findInsertionPoint(predId, op.Id, predNode, s)

	var node *treap.Node[Op]
	if insertBefore == nil {
		node = lt.r.PushBack(op)
	} else {
		node = lt.r.InsertBefore(op, insertBefore)
	}
	lt.byId[op.Id] = node

	if lt.rope != nil {
		if ch, ok := op.Action.Value.(string); ok {
			liveRank := lt.r.LiveRank(node)
			lt.rope.InsertAt(liveRank, ch)
		}
	}
}

// markDead marks the insert op with the given id as a tombstone. The node
// remains in the treap as a position marker but is excluded from Len and All.
// For text objects the rope is updated in sync before the treap node is
// marked dead (LiveRank must be queried while the node is still alive).
func (lt *listTreap) markDead(id types.OpId) {
	node, ok := lt.byId[id]
	if !ok {
		return
	}

	if lt.rope != nil {
		liveRank := lt.r.LiveRank(node)
		lt.rope.DeleteAt(liveRank)
	}

	lt.r.MarkDead(node)
	if lt.liveValues != nil {
		delete(lt.liveValues, id)
	}
}

// insertOpInListTreap applies an insert op incrementally to the cached treap
// for op.Object, if one exists. No-op if the treap has not been built yet.
func (s *OpSet) insertOpInListTreap(op Op) {
	if !op.Insert {
		panic("insertOpInListTreap called with non-insert op")
	}
	lt, ok := s.listTreaps[op.Object]
	if !ok {
		return
	}
	lt.insertOp(op, s)
}

// onInsertKilledInListTreap is called when a previously live insert op has
// its SuccCount incremented from 0 to 1 (it just became dead). obj is the
// object the insert op belongs to.
//
// If a live update op exists at that position the position remains alive and
// we fall back to a full rebuild, because we would otherwise need a Revive
// operation on the treap node.
func (s *OpSet) onInsertKilledInListTreap(obj types.ObjectId, insertId types.OpId) {
	lt, ok := s.listTreaps[obj]
	if !ok {
		return
	}
	if s.hasLiveUpdateAt(insertId) {
		s.invalidateListTreap(obj)
		return
	}
	lt.markDead(insertId)
}

// hasLiveUpdateAt reports whether any delta op targets posId as a non-delete
// update (Insert=false, Action≠Delete) with SuccCount==0.
func (s *OpSet) hasLiveUpdateAt(posId types.OpId) bool {
	if s.deltaSuccessors == nil {
		return false
	}
	for _, succId := range s.deltaSuccessors[posId] {
		idx, ok := s.delta.byId[succId]
		if !ok {
			continue
		}
		op := s.delta.ops[idx]
		if !op.Insert && op.Action.Kind != types.ActionDelete && op.SuccCount == 0 {
			return true
		}
	}
	return false
}

// findInsertionPoint finds the treap node before which newOp should be
// inserted, given its CRDT predecessor pred and the predecessor's treap node
// predNode (nil if pred is the sentinel).
//
// It scans forward using Next(), skipping higher-priority siblings (same pred,
// higher OpId) and their entire CRDT subtrees.
func (lt *listTreap) findInsertionPoint(
	pred types.OpId,
	newOpId types.OpId,
	predNode *treap.Node[Op],
	s *OpSet,
) *treap.Node[Op] {
	var cur *treap.Node[Op]
	if predNode == nil {
		cur = lt.r.Front()
	} else {
		cur = treap.Next(predNode)
	}

	var skipRoot types.OpId
	skipActive := false

	for cur != nil {
		curKey, ok := cur.Value().Key.(types.KeyOpId)
		if !ok {
			panic(fmt.Sprintf("list treap node has non-KeyOpId key: %T", cur.Value().Key))
		}
		curParent := types.OpId(curKey)

		if skipActive {
			// Are we still inside the skip-root's CRDT subtree?
			if curParent == skipRoot || lt.isCrdtDescendant(curParent, skipRoot) {
				cur = treap.Next(cur)
				continue
			}
			skipActive = false
		}

		if curParent == pred {
			if s.opIdGreater(cur.Value().Id, newOpId) {
				// Higher-priority sibling: skip it and its entire subtree.
				skipRoot = cur.Value().Id
				skipActive = true
				cur = treap.Next(cur)
				continue
			}
			// Lower-priority sibling: newOp goes before cur.
			return cur
		}

		// cur is not a sibling of newOp; we've gone past pred's children.
		return cur
	}

	return nil // insert at end
}

// isCrdtDescendant reports whether ancestor is in the CRDT ancestry chain
// of child (walking up via each node's predecessor key).
func (lt *listTreap) isCrdtDescendant(child, ancestor types.OpId) bool {
	cur := child
	for {
		if cur == (types.OpId{}) {
			return false
		}
		node, ok := lt.byId[cur]
		if !ok {
			return false
		}
		k, ok := node.Value().Key.(types.KeyOpId)
		if !ok {
			return false
		}
		cur = types.OpId(k)
		if cur == ancestor {
			return true
		}
	}
}

// iterListElements runs the CRDT list reconstruction for obj and calls fn for
// each live element in order. This is kept for future use or as a fallback
// but is no longer called internally (buildListTreap has its own scan).
func (s *OpSet) iterListElements(obj types.ObjectId, fn func(Op)) {
	children := make(map[types.Key][]types.OpId)
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

	for key := range children {
		sortOpIdsDesc(children[key], s)
	}

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
