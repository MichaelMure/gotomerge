package opset

import (
	"fmt"

	"gotomerge/format"
	"gotomerge/types"
)

// changeActorMap translates chunk-local actor indices to OpSet global indices.
// Index 0 is the change's own actor; indices 1..N are OtherActors[i-1].
type changeActorMap []uint32

func (m changeActorMap) opId(id types.OpId) types.OpId {
	return types.OpId{ActorIdx: m[id.ActorIdx], Counter: id.Counter}
}

func (m changeActorMap) objectId(obj types.ObjectId) types.ObjectId {
	if obj.IsRoot() {
		return obj
	}
	return types.ObjectId(m.opId(types.OpId(obj)))
}

func (m changeActorMap) key(key types.Key) types.Key {
	if k, ok := key.(types.KeyOpId); ok {
		return types.KeyOpId{ActorIdx: m[k.ActorIdx], Counter: k.Counter}
	}
	return key
}

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
	cam := make(changeActorMap, 1+len(cc.OtherActors))
	cam[0] = s.internActor(cc.Actor)
	for i, a := range cc.OtherActors {
		cam[i+1] = s.internActor(a)
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
			resolvedPred := cam.opId(pred)
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
			Id:     cam.opId(changeOp.Id),
			Object: cam.objectId(changeOp.Object),
			Key:    cam.key(changeOp.Key),
			Insert: changeOp.Insert,
			Action: changeOp.Action,
		}

		idx := uint32(len(s.delta.ops))
		s.delta.ops = append(s.delta.ops, op)
		s.delta.byId[op.Id] = idx
		s.delta.byObj[op.Object] = append(s.delta.byObj[op.Object], idx)
		if op.Id.Counter > s.maxOpCounter[op.Id.ActorIdx] {
			s.maxOpCounter[op.Id.ActorIdx] = op.Id.Counter
		}
	}

	// Update heads: remove deps (they now have a descendant), add this change.
	for _, dep := range cc.Dependencies {
		delete(s.heads, dep)
	}
	s.heads[cc.Hash] = struct{}{}
	s.appliedHashes[cc.Hash] = struct{}{}

	return nil
}
