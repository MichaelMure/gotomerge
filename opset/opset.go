package opset

import (
	"bytes"

	"github.com/MichaelMure/gotomerge/types"
)

// OpSet holds every operation that has been applied to a document.
//
// All OpIds stored here use indices into the OpSet's own actors slice.
// Indices from incoming chunks (DocumentChunk.Actors, ChangeChunk.Actor /
// OtherActors) are translated at apply time.
type OpSet struct {
	// actors is the document-wide actor table. OpId.ActorIdx values throughout
	// this OpSet are indices into this slice.
	actors   []types.ActorId
	actorIdx map[string]uint32 // string(actorId) → index in actors

	// snapshot holds operations from the most recently applied DocumentChunk,
	// stored as column references plus sparse metadata. Nil if no DocumentChunk
	// has been applied yet.
	snapshot *snapshotStore

	// delta holds operations from applied ChangeChunks, stored as mutable
	// structs. Nil if no ChangeChunks have been applied yet.
	delta *deltaStore

	// heads is the set of change hashes that no other applied change depends on.
	// It represents the current version of the document.
	heads map[types.ChangeHash]struct{}

	// appliedHashes tracks hashes of all applied ChangeChunks for dependency
	// validation. Changes in a DocumentChunk are not enumerated here; the
	// presence of snapshot implies all pre-snapshot deps are satisfied.
	appliedHashes map[types.ChangeHash]struct{}

	// maxOpCounter tracks the highest counter seen per actor (global index).
	// Used by Begin to compute startOp for new transactions.
	maxOpCounter map[uint32]uint32

	// counterDeltas accumulates increment deltas for counter ops.
	// Key: OpId of the ActionSet(Counter) op; Value: sum of all applied Inc deltas.
	// Maintained by ApplyChange (ActionInc ops) and ApplyDocument (post-scan pass).
	counterDeltas map[types.OpId]int64

	// deltaSuccessors maps predecessor OpId → slice of successor OpIds for all
	// ops applied via ApplyChange. Built incrementally so ExportDocument does not
	// need to re-parse every raw change chunk. Reset to nil when ApplyDocument is
	// called, because the new snapshot already encodes within-snapshot successors.
	deltaSuccessors map[types.OpId][]types.OpId
}

func New() *OpSet {
	return &OpSet{
		actorIdx:      make(map[string]uint32),
		heads:         make(map[types.ChangeHash]struct{}),
		appliedHashes: make(map[types.ChangeHash]struct{}),
		maxOpCounter:  make(map[uint32]uint32),
		counterDeltas: make(map[types.OpId]int64),
	}
}

// Actor returns the ActorId for the given index.
func (s *OpSet) Actor(idx uint32) types.ActorId {
	return s.actors[idx]
}

// Heads returns the hashes of the changes that no other applied change depends
// on — the current version of the document.
func (s *OpSet) Heads() []types.ChangeHash {
	out := make([]types.ChangeHash, 0, len(s.heads))
	for h := range s.heads {
		out = append(out, h)
	}
	return out
}

// AppliedHashes returns the set of all ChangeChunk hashes that have been applied.
// The returned map must not be modified by the caller.
func (s *OpSet) AppliedHashes() map[types.ChangeHash]struct{} {
	return s.appliedHashes
}

// internActor registers an actor and returns its index in the OpSet actor table.
// If the actor is already registered, the existing index is returned.
func (s *OpSet) internActor(id types.ActorId) uint32 {
	key := string(id)
	if idx, ok := s.actorIdx[key]; ok {
		return idx
	}
	idx := uint32(len(s.actors))
	s.actors = append(s.actors, id)
	s.actorIdx[key] = idx
	return idx
}

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
