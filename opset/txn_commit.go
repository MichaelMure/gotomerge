package opset

import (
	"fmt"
	"io"

	"gotomerge/format"
	"gotomerge/types"
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

// applyCommitted applies the transaction's ops directly to s.delta, bypassing
// the binary format entirely. All OpIds in t.ops already use global actor indices,
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

	for i, txOp := range t.ops {
		opId := types.OpId{ActorIdx: t.actorIdx, Counter: t.startOp + uint32(i)}

		// Increment SuccCount on each predecessor this op supersedes.
		for _, pred := range txOp.preds {
			if s.snapshot != nil {
				if predIdx, ok := s.snapshot.byId[pred]; ok {
					s.snapshot.succCount[predIdx]++
				}
			}
			if predIdx, ok := s.delta.byId[pred]; ok {
				s.delta.ops[predIdx].SuccCount++
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
