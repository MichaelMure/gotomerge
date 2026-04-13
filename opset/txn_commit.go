package opset

import (
	"fmt"
	"io"

	"github.com/MichaelMure/gotomerge/format"
	"github.com/MichaelMure/gotomerge/types"
)

// HasOps reports whether the transaction has any buffered operations.
func (t *Transaction) HasOps() bool { return len(t.ops) > 0 }

// Commit encodes the buffered operations as a ChangeChunk, applies it locally,
// and writes the serialised chunk to w. If there are no operations Commit is a
// no-op: it writes nothing and returns nil.
func (t *Transaction) Commit(w io.Writer) error {
	if len(t.ops) == 0 {
		return nil
	}

	others := t.otherActors()

	// Map global actorIdx → local index: 0 = own actor, 1..N = others in sort order.
	mapper := types.NewActorMapper(len(t.s.actors))
	mapper.Add(t.actorIdx)
	for _, a := range others {
		mapper.Add(t.s.actorIdx[string(a)])
	}

	ops := format.NewChangeOpsWriter()
	for _, op := range t.ops {
		ops.Append(op.obj, op.key, op.insert, op.action, op.preds, mapper)
	}

	cc := &format.ChangeChunk{
		Dependencies: t.deps,
		Actor:        t.actor,
		SeqNum:       t.seqNum,
		StartOp:      uint64(t.startOp),
		OtherActors:  others,
	}
	if err := format.WriteChange(w, cc, ops); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return t.applyCommitted(cc.Hash)
}

// applyCommitted applies the transaction's ops directly to s.delta.
// All OpIds in t.ops already use global actor indices,
// so no translation is needed — unlike ApplyChange which must translate from
// change-local indices.
func (t *Transaction) applyCommitted(hash types.ChangeHash) error {
	s := t.s

	// Validate that all declared dependencies are already present.
	for _, dep := range t.deps {
		if _, ok := s.appliedHashes[dep]; ok {
			continue
		}
		if s.snapshot != nil {
			continue
		}
		return fmt.Errorf("unsatisfied dependency: %s", dep)
	}

	if s.delta == nil {
		s.delta = newDeltaStore()
	}
	if s.deltaSuccessors == nil {
		s.deltaSuccessors = make(map[types.OpId][]types.OpId)
	}

	for i, txOp := range t.ops {
		opId := types.OpId{ActorIdx: t.actorIdx, Counter: t.startOp + uint32(i)}

		// For ActionInc: accumulate the counter delta without killing the
		// counter op (same logic as ApplyChange / incValue).
		// For all other ops: increment SuccCount on each predecessor.
		if txOp.action.Kind == types.ActionInc {
			delta := incValue(txOp.action.Value)
			for _, pred := range txOp.preds {
				s.counterDeltas[pred] += delta
				s.deltaSuccessors[pred] = append(s.deltaSuccessors[pred], opId)
			}
		} else {
			for _, pred := range txOp.preds {
				if s.snapshot != nil {
					if predIdx, ok := s.snapshot.byId[pred]; ok {
						wasZero := s.snapshot.succCount[predIdx] == 0
						s.snapshot.succCount[predIdx]++
						if wasZero {
							if obj, ok := s.snapshot.objectForOp(predIdx); ok {
								if s.snapshot.insert.Get(predIdx) {
									s.onInsertKilledInListTreap(obj, pred)
								} else {
									s.invalidateListTreap(obj)
								}
							}
						}
					}
				}
				if predIdx, ok := s.delta.byId[pred]; ok {
					predOp := &s.delta.ops[predIdx]
					wasZero := predOp.SuccCount == 0
					predOp.SuccCount++
					if wasZero {
						if predOp.Insert {
							s.onInsertKilledInListTreap(predOp.Object, pred)
						} else {
							s.invalidateListTreap(predOp.Object)
						}
					}
				}
				s.deltaSuccessors[pred] = append(s.deltaSuccessors[pred], opId)
			}
		}

		op := Op{
			Id:     opId,
			Object: txOp.obj,
			Key:    txOp.key,
			Insert: txOp.insert,
			Action: txOp.action,
		}

		idx := uint32(len(s.delta.ops))
		s.delta.ops = append(s.delta.ops, op)
		s.delta.byId[op.Id] = idx
		s.delta.byObj[op.Object] = append(s.delta.byObj[op.Object], idx)
		if k, ok := op.Key.(types.KeyString); ok && !op.Insert &&
			op.Action.Kind != types.ActionDelete && op.Action.Kind != types.ActionInc {
			s.delta.addToMapKeys(op.Object, string(k), idx)
		}
		if op.Insert {
			s.insertOpInListTreap(op)
		} else if _, isPos := op.Key.(types.KeyOpId); isPos && op.Action.Kind == types.ActionSet {
			// Non-insert ActionSet at a list position: the winning value at
			// that position may change, so the cached treap must be rebuilt.
			s.invalidateListTreap(op.Object)
		}
		// A MakeText op may have had its treap built without a rope during
		// the transaction (before the op was in the delta). Drop it so the
		// next access rebuilds correctly.
		if op.Action.Kind == types.ActionMakeText {
			s.invalidateListTreap(types.ObjectId(op.Id))
		}
		if op.Id.Counter > s.maxOpCounter[op.Id.ActorIdx] {
			s.maxOpCounter[op.Id.ActorIdx] = op.Id.Counter
		}
	}

	// Update heads: remove deps (they now have a descendant), add this change.
	for _, dep := range t.deps {
		delete(s.heads, dep)
	}
	s.heads[hash] = struct{}{}
	s.appliedHashes[hash] = struct{}{}

	return nil
}
