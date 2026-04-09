package opset

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"gotomerge/format"
	"gotomerge/types"
)

// ExportDocument writes a complete document chunk to w representing the full
// current state of the OpSet. It merges snapshot (if any) and delta operations
// into a single object-sorted sequence, reconstructing successor lists from the
// provided applied change chunks.
//
// changes must be the slice of all ChangeChunks that have been applied (in
// application order). Their predecessor lists are used to reconstruct the
// successor map. For snapshot ops, within-snapshot successors are recovered from
// the original DocumentChunk stored during ApplyDocument.
//
// Change metadata (per-change actor/seqNum/maxOp/deps) is not included.
// ApplyDocument does not require it; it can be added later for sync support.
func (s *OpSet) ExportDocument(changes []*format.ChangeChunk, w io.Writer) error {
	// The document format requires actors sorted lexicographically.
	// Build a sorted copy of the actor table and a remapping from OpSet global
	// indices to sorted-document indices, used as localOf when encoding ops.
	sortedActors, localOf := sortedActorTable(s.actors)

	// Build delta successor map: predecessor OpId → list of successor OpIds.
	deltaSuccessors := buildDeltaSuccessors(s, changes)

	ops := format.NewDocOpsWriter()

	if s.snapshot != nil {
		if err := exportSnapshotOps(s, ops, localOf, deltaSuccessors); err != nil {
			return fmt.Errorf("export snapshot ops: %w", err)
		}
	}

	if s.delta != nil {
		exportDeltaOps(s, ops, localOf, deltaSuccessors)
	}

	return format.WriteDocument(w, sortedActors, s.Heads(), nil, nil, ops)
}

// buildDeltaSuccessors builds a map from each predecessor OpId to the list of
// delta operation OpIds that declare it as a predecessor.
func buildDeltaSuccessors(s *OpSet, changes []*format.ChangeChunk) map[types.OpId][]types.OpId {
	result := make(map[types.OpId][]types.OpId)
	for _, cc := range changes {
		// Rebuild the local → global actor index mapping used in ApplyChange.
		cam := make(changeActorMap, 1+len(cc.OtherActors))
		cam[0] = s.actorIdx[string(cc.Actor)]
		for i, a := range cc.OtherActors {
			cam[i+1] = s.actorIdx[string(a)]
		}

		for changeOp, err := range cc.Operations() {
			if err != nil {
				continue // best-effort; ApplyChange already validated these
			}
			myId := cam.opId(changeOp.Id)
			for _, pred := range changeOp.Predecessors {
				predGlobal := cam.opId(pred)
				result[predGlobal] = append(result[predGlobal], myId)
			}
		}
	}
	return result
}

// exportSnapshotOps re-emits all snapshot ops in their original document order,
// augmenting each op's successor list with any new delta successors.
func exportSnapshotOps(s *OpSet, ops *format.DocOpsWriter, mapper types.ActorMapper, deltaSuccessors map[types.OpId][]types.OpId) error {
	ss := s.snapshot
	dc := ss.docChunk
	if dc == nil {
		return fmt.Errorf("snapshot has no stored DocumentChunk")
	}

	// docChunk.Operations() yields DocOperation{Object, Key, Id, Insert, Action, Successors}.
	// Actor indices inside equal the OpSet's global indices: ApplyDocument registered actors
	// in the same order, and the OpSet was empty at that time.
	for docOp, err := range dc.Operations() {
		if err != nil {
			return fmt.Errorf("iterating snapshot ops: %w", err)
		}
		// Merge within-snapshot successors with new delta successors.
		allSuccessors := docOp.Successors
		if ds := deltaSuccessors[docOp.Id]; len(ds) > 0 {
			allSuccessors = append(allSuccessors, ds...)
		}
		ops.Append(docOp.Object, docOp.Key, docOp.Id, docOp.Insert, docOp.Action, allSuccessors, mapper)
	}
	return nil
}

// exportDeltaOps emits all delta ops sorted by object (ObjectId creation order),
// then by their position within each object (delta application order).
func exportDeltaOps(s *OpSet, enc *format.DocOpsWriter, mapper types.ActorMapper, deltaSuccessors map[types.OpId][]types.OpId) {
	// Collect all unique objects that have delta ops.
	objects := make([]types.ObjectId, 0, len(s.delta.byObj))
	for obj := range s.delta.byObj {
		objects = append(objects, obj)
	}
	sortObjects(objects)

	for _, obj := range objects {
		for _, idx := range s.delta.byObj[obj] {
			op := s.delta.ops[idx]
			succs := deltaSuccessors[op.Id]
			enc.Append(op.Object, op.Key, op.Id, op.Insert, op.Action, succs, mapper)
		}
	}
}

// sortObjects sorts a slice of ObjectIds in deterministic document order:
// root first, then by (counter ascending, actorIdx ascending).
func sortObjects(objects []types.ObjectId) {
	sort.Slice(objects, func(i, j int) bool {
		a, b := objects[i], objects[j]
		if a.IsRoot() {
			return true
		}
		if b.IsRoot() {
			return false
		}
		aId, bId := types.OpId(a), types.OpId(b)
		if aId.Counter != bId.Counter {
			return aId.Counter < bId.Counter
		}
		return aId.ActorIdx < bId.ActorIdx
	})
}

// sortedActorTable returns a lexicographically sorted copy of actors and an
// ActorMapper that translates original OpSet actor indices to their new sorted
// positions. The document format requires actors in sorted order.
func sortedActorTable(actors []types.ActorId) ([]types.ActorId, types.ActorMapper) {
	n := len(actors)
	// Build index slice, sort it by actor bytes.
	indices := make([]uint32, n)
	for i := range indices {
		indices[i] = uint32(i)
	}
	sort.Slice(indices, func(i, j int) bool {
		return bytes.Compare(actors[indices[i]], actors[indices[j]]) < 0
	})

	sorted := make([]types.ActorId, n)
	mapper := types.NewActorMapper(n)
	for newIdx, oldIdx := range indices {
		sorted[newIdx] = actors[oldIdx]
		mapper.Add(oldIdx)
	}
	return sorted, mapper
}
