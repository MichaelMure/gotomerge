package opset

import (
	"fmt"
	"io"

	"gotomerge/format"
)

// Commit encodes the buffered operations as a ChangeChunk, applies it locally,
// and writes the serialised chunk to w.
func (t *Transaction) Commit(w io.Writer) error {
	if len(t.ops) == 0 {
		return fmt.Errorf("commit: no operations")
	}

	others := t.otherActors()

	// TODO: this logic could be extracted, kinda like a MMU handling redirection
	// Map global actorIdx → local index: 0 = own actor, 1..N = others in sort order.
	localOf := make(map[uint32]uint32, 1+len(others))
	localOf[t.actorIdx] = 0
	for i, a := range others {
		localOf[t.s.actorIdx[string(a)]] = uint32(i + 1)
	}

	enc := format.NewChangeOpsWriter()
	for _, op := range t.ops {
		enc.Append(op.obj, op.key, op.insert, op.action, op.preds, localOf)
	}
	if err := enc.Finalise(); err != nil {
		return fmt.Errorf("commit: encode: %w", err)
	}

	cc := &format.ChangeChunk{
		Dependencies: t.deps,
		Actor:        t.actor,
		SeqNum:       t.seqNum,
		StartOp:      uint64(t.startOp),
		OtherActors:  others,
	}
	if err := format.WriteChange(w, cc, enc); err != nil {
		return fmt.Errorf("commit: %w", err)
	}
	return t.s.ApplyChange(cc)
}
