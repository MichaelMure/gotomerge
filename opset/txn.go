package opset

import (
	"bytes"
	"sort"

	"github.com/MichaelMure/gotomerge/types"
)

// txnOp is one pending operation buffered in a Transaction.
// All OpIds and ObjectIds use global OpSet actor indices.
type txnOp struct {
	obj    types.ObjectId
	key    types.Key
	insert bool
	action types.Action
	preds  []types.OpId
}

// Transaction accumulates operations to be committed as a single ChangeChunk.
// Reads go through the live OpSet; writes are buffered until Commit.
type Transaction struct {
	s        *OpSet
	actor    types.ActorId
	actorIdx uint32 // index in s.actors
	seqNum   uint64
	startOp  uint32
	deps     []types.ChangeHash // heads snapshot at Begin time
	ops      []txnOp
	lists    map[types.ObjectId]*workingList // eager working state per list/text object
}

// Begin starts a new transaction. seqNum must increase monotonically for
// this actor across all changes.
func (s *OpSet) Begin(actor types.ActorId, seqNum uint64) *Transaction {
	actorIdx := s.internActor(actor)
	return &Transaction{
		s:        s,
		actor:    actor,
		actorIdx: actorIdx,
		seqNum:   seqNum,
		startOp:  s.maxOpCounter[actorIdx] + 1,
		deps:     s.Heads(),
	}
}

// nextOpId returns the global OpId the next operation will receive.
func (t *Transaction) nextOpId() types.OpId {
	return types.OpId{ActorIdx: t.actorIdx, Counter: t.startOp + uint32(len(t.ops))}
}

// opsToIds extracts the OpId from each Op.
func opsToIds(ops []Op) []types.OpId {
	ids := make([]types.OpId, len(ops))
	for i, op := range ops {
		ids[i] = op.Id
	}
	return ids
}

// otherActors returns all actors referenced in t.ops other than t.actor,
// sorted ascending by actor bytes (required by the binary format).
func (t *Transaction) otherActors() []types.ActorId {
	seen := make(map[uint32]struct{})
	add := func(globalIdx uint32) {
		if globalIdx != t.actorIdx {
			seen[globalIdx] = struct{}{}
		}
	}
	for _, op := range t.ops {
		if !op.obj.IsRoot() {
			add(op.obj.ActorIdx)
		}
		if k, ok := op.key.(types.KeyOpId); ok && k.Counter != 0 {
			add(k.ActorIdx)
		}
		for _, pred := range op.preds {
			add(pred.ActorIdx)
		}
	}
	actors := make([]types.ActorId, 0, len(seen))
	for idx := range seen {
		actors = append(actors, t.s.actors[idx])
	}
	sort.Slice(actors, func(i, j int) bool {
		return bytes.Compare(actors[i], actors[j]) < 0
	})
	return actors
}
