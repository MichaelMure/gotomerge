package opset

import (
	"bytes"

	"github.com/MichaelMure/gotomerge/types"
)

// ObjType returns the ActionKind of the operation that created obj
// (ActionMakeMap, ActionMakeList, or ActionMakeText). The root object
// always returns ActionMakeMap. Returns false if the object is unknown.
func (s *OpSet) ObjType(obj types.ObjectId) (types.ActionKind, bool) {
	if obj.IsRoot() {
		return types.ActionMakeMap, true
	}
	if s.snapshot != nil {
		if kind, ok := s.snapshot.objCreators[obj]; ok {
			return kind, true
		}
	}
	if s.delta != nil {
		if idx, ok := s.delta.byId[types.OpId(obj)]; ok {
			return s.delta.ops[idx].Action.Kind, true
		}
	}
	return 0, false
}

// applyCounterDelta returns op with its Counter value adjusted by any
// accumulated increments. For non-counter ops it is a no-op.
func (s *OpSet) applyCounterDelta(op Op) Op {
	if op.Action.Kind != types.ActionSet {
		return op
	}
	base, isCounter := op.Action.Value.(types.Counter)
	if !isCounter {
		return op
	}
	if delta, ok := s.counterDeltas[op.Id]; ok {
		op.Action.Value = types.Counter(int64(base) + delta)
	}
	return op
}

// opIdGreater reports whether a is greater than b.
// Higher Counter wins; equal counters are broken by actor bytes (higher wins).
func (s *OpSet) opIdGreater(a, b types.OpId) bool {
	if a.Counter != b.Counter {
		return a.Counter > b.Counter
	}
	return bytes.Compare(s.actors[a.ActorIdx], s.actors[b.ActorIdx]) > 0
}

// Materialize converts an Op to a plain Go value. For scalar ops it returns
// op.Action.Value directly. For Make* ops it returns the ObjectId of the
// created object so the caller can wrap it in a reader — the docproxy layer
// does this wrapping, keeping opset free of view types.
func (s *OpSet) Materialize(op Op) any {
	switch op.Action.Kind {
	case types.ActionMakeMap, types.ActionMakeList, types.ActionMakeText:
		return types.ObjectId(op.Id)
	default:
		return op.Action.Value
	}
}

// sortOpIdsDesc sorts OpIds in descending order (highest first).
func sortOpIdsDesc(ids []types.OpId, s *OpSet) {
	for i := 1; i < len(ids); i++ {
		for j := i; j > 0 && s.opIdGreater(ids[j], ids[j-1]); j-- {
			ids[j], ids[j-1] = ids[j-1], ids[j]
		}
	}
}

// sortOpsDesc sorts ops in descending OpId order (highest first).
func sortOpsDesc(ops []Op, s *OpSet) {
	for i := 1; i < len(ops); i++ {
		for j := i; j > 0 && s.opIdGreater(ops[j].Id, ops[j-1].Id); j-- {
			ops[j], ops[j-1] = ops[j-1], ops[j]
		}
	}
}
