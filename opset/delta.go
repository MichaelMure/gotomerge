package opset

import (
	"fmt"

	"gotomerge/format"
	"gotomerge/types"
)

// deltaStore holds operations from applied ChangeChunks as mutable structs.
//
// Change chunks carry no successor column — a change cannot know what will
// supersede it later — so SuccCount must be writable in place. Two indexes
// support efficient access: byId for O(1) predecessor lookup when applying
// subsequent changes, and byObj for O(1) per-object scans during queries.
//
// The delta can grow unboundedly as changes accumulate. It can be compacted
// back into a snapshot by serializing the full document (Save) and reloading
// it as a new DocumentChunk.
type deltaStore struct {
	ops   []Op
	byId  map[types.OpId]uint32       // opId → index in ops
	byObj map[types.ObjectId][]uint32 // objectId → indices in ops
}

func newDeltaStore() *deltaStore {
	return &deltaStore{
		byId:  make(map[types.OpId]uint32),
		byObj: make(map[types.ObjectId][]uint32),
	}
}

// ApplyChange applies a single change to the OpSet. Changes must be applied in
// dependency order: all changes listed in cc.Dependencies must already be present.
func (s *OpSet) ApplyChange(cc *format.ChangeChunk) error {
	// Validate that all declared dependencies are already present.
	for _, dep := range cc.Dependencies {
		if _, ok := s.appliedHashes[dep]; ok {
			continue
		}
		// A snapshot implicitly satisfies all pre-snapshot dependencies.
		if s.snapshot != nil {
			continue
		}
		return fmt.Errorf("unsatisfied dependency: %s", dep)
	}

	// Build a mapping from the change's local actor indices to ours.
	// Index 0 is the change's own actor; indices 1..N are OtherActors[i-1].
	changeActors := make([]uint32, 1+len(cc.OtherActors))
	changeActors[0] = s.internActor(cc.Actor)
	for i, a := range cc.OtherActors {
		changeActors[i+1] = s.internActor(a)
	}

	translateId := func(id types.OpId) types.OpId {
		return types.OpId{ActorIdx: changeActors[id.ActorIdx], Counter: id.Counter}
	}
	translateObj := func(obj types.ObjectId) types.ObjectId {
		if obj.IsRoot() {
			return obj
		}
		return types.ObjectId(translateId(types.OpId(obj)))
	}
	translateKey := func(key types.Key, actorMap []uint32) types.Key {
		if k, ok := key.(types.KeyOpId); ok {
			return types.KeyOpId{ActorIdx: actorMap[k.ActorIdx], Counter: k.Counter}
		}
		return key
	}

	if s.delta == nil {
		s.delta = newDeltaStore()
	}

	for changeOp, err := range cc.Operations() {
		if err != nil {
			return fmt.Errorf("reading operation: %w", err)
		}

		// Increment SuccCount on each predecessor this op supersedes.
		for _, pred := range changeOp.Predecessors {
			resolvedPred := translateId(pred)
			if s.snapshot != nil {
				if predIdx, ok := s.snapshot.byId[resolvedPred]; ok {
					s.snapshot.succCount[predIdx]++
				}
			}
			if predIdx, ok := s.delta.byId[resolvedPred]; ok {
				s.delta.ops[predIdx].SuccCount++
			}
		}

		op := Op{
			Id:     translateId(changeOp.Id),
			Object: translateObj(changeOp.Object),
			Key:    translateKey(changeOp.Key, changeActors),
			Insert: changeOp.Insert,
			Action: changeOp.Action,
		}

		idx := uint32(len(s.delta.ops))
		s.delta.ops = append(s.delta.ops, op)
		s.delta.byId[op.Id] = idx
		s.delta.byObj[op.Object] = append(s.delta.byObj[op.Object], idx)
	}

	// Update heads: remove deps (they now have a descendant), add this change.
	for _, dep := range cc.Dependencies {
		delete(s.heads, dep)
	}
	s.heads[cc.Hash] = struct{}{}
	s.appliedHashes[cc.Hash] = struct{}{}

	return nil
}
