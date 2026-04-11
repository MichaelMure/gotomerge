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

// MapSet sets obj[key] = value, superseding any existing live values.
func (t *Transaction) MapSet(obj types.ObjectId, key string, value any) types.OpId {
	id := t.nextOpId()
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionSet, Value: value},
		preds:  opsToIds(t.s.MapGetAll(obj, key)),
	})
	return id
}

// MapIncrement adds delta to the counter at obj[key]. The key must currently
// hold a Counter value; if it does not, this is a no-op. Panics if read-only.
func (t *Transaction) MapIncrement(obj types.ObjectId, key string, delta int64) {
	preds := opsToIds(t.s.MapGetAll(obj, key))
	if len(preds) == 0 {
		return
	}
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionInc, Value: delta},
		preds:  preds,
	})
}

// MapDelete deletes obj[key] by superseding all live values there.
// No-ops if the key has no live value.
func (t *Transaction) MapDelete(obj types.ObjectId, key string) {
	preds := opsToIds(t.s.MapGetAll(obj, key))
	if len(preds) == 0 {
		return
	}
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionDelete},
		preds:  preds,
	})
}

// MakeMap creates a nested map at obj[key] and returns its ObjectId and OpId.
func (t *Transaction) MakeMap(obj types.ObjectId, key string) (types.ObjectId, types.OpId) {
	id := t.nextOpId()
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionMakeMap},
		preds:  opsToIds(t.s.MapGetAll(obj, key)),
	})
	return types.ObjectId(id), id
}

// MakeList creates a list at obj[key] and returns its ObjectId and OpId.
func (t *Transaction) MakeList(obj types.ObjectId, key string) (types.ObjectId, types.OpId) {
	id := t.nextOpId()
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionMakeList},
		preds:  opsToIds(t.s.MapGetAll(obj, key)),
	})
	return types.ObjectId(id), id
}

// MakeText creates a text object at obj[key] and returns its ObjectId and OpId.
func (t *Transaction) MakeText(obj types.ObjectId, key string) (types.ObjectId, types.OpId) {
	id := t.nextOpId()
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionMakeText},
		preds:  opsToIds(t.s.MapGetAll(obj, key)),
	})
	return types.ObjectId(id), id
}

// ListInsert inserts a new element after pred in list obj.
// Use types.KeyOpId{} (zero Counter = head sentinel) to insert at the front.
func (t *Transaction) ListInsert(obj types.ObjectId, pred types.Key, value any) types.OpId {
	id := t.nextOpId()
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    pred,
		insert: true,
		action: types.Action{Kind: types.ActionSet, Value: value},
	})
	return id
}

// ListDelete marks element posId as deleted. liveOpId is the OpId of the
// current live value at that position (use posId when the position has never
// been updated by a separate op).
func (t *Transaction) ListDelete(obj types.ObjectId, posId types.OpId, liveOpId types.OpId) {
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyOpId(posId),
		action: types.Action{Kind: types.ActionDelete},
		preds:  []types.OpId{liveOpId},
	})
}

// ListElements returns the live elements of a list or text object in order,
// including any operations buffered in this transaction that have not yet been
// committed. Use this inside a transaction instead of OpSet.ListElements when
// you need write-your-own-reads (e.g. multiple splices in one Change).
func (t *Transaction) ListElements(obj types.ObjectId) []Op {
	base := t.s.ListElements(obj)

	// Fast path: no pending insert/delete ops for this object.
	hasPending := false
	for _, op := range t.ops {
		if op.obj == obj && (op.insert || op.action.Kind == types.ActionDelete) {
			hasPending = true
			break
		}
	}
	if !hasPending {
		return base
	}

	// Apply pending ops sequentially on top of the committed list.
	working := make([]Op, len(base))
	copy(working, base)

	findById := func(id types.OpId) int {
		for i, e := range working {
			if e.Id == id {
				return i
			}
		}
		return -1
	}

	for i, op := range t.ops {
		if op.obj != obj {
			continue
		}
		opId := types.OpId{ActorIdx: t.actorIdx, Counter: t.startOp + uint32(i)}

		if op.insert {
			newOp := Op{Id: opId, Object: obj, Key: op.key, Insert: true, Action: op.action}
			pred, isPred := op.key.(types.KeyOpId)
			insertAt := len(working) // default: append
			if isPred && (pred.ActorIdx != 0 || pred.Counter != 0) {
				if idx := findById(types.OpId(pred)); idx >= 0 {
					insertAt = idx + 1
				}
			} else {
				insertAt = 0 // head sentinel: prepend
			}
			working = append(working, Op{})
			copy(working[insertAt+1:], working[insertAt:])
			working[insertAt] = newOp
		} else if op.action.Kind == types.ActionDelete {
			for _, pred := range op.preds {
				if idx := findById(pred); idx >= 0 {
					working = append(working[:idx], working[idx+1:]...)
					break
				}
			}
		}
	}

	return working
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
