package opset

import (
	"gotomerge/types"
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
}

func New() *OpSet {
	return &OpSet{
		actorIdx:      make(map[string]uint32),
		heads:         make(map[types.ChangeHash]struct{}),
		appliedHashes: make(map[types.ChangeHash]struct{}),
		maxOpCounter:  make(map[uint32]uint32),
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
