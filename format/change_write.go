package format

import (
	"bytes"
	"fmt"
	"io"

	"github.com/MichaelMure/leb128"

	"gotomerge/column"
	"gotomerge/types"
)

// WriteChange serialises cc as a complete Automerge change chunk and writes it
// to w. Payloads larger than deflateMinSize bytes are written as
// ChunkTypeCompressedChange; smaller ones as plain ChunkTypeChange.
//
// After WriteChange returns, cc.Hash is set and cc is ready for ApplyChange.
func WriteChange(w io.Writer, cc *ChangeChunk, ops *ChangeOpsWriter) error {
	var payload bytes.Buffer
	if err := writeChangePayload(&payload, cc, ops); err != nil {
		return err
	}

	if payload.Len() > deflateMinSize {
		return writeChangeChunkCompressed(w, payload.Bytes(), &cc.Hash)
	}
	return writeChunk(w, ChunkTypeChange, payload.Bytes(), &cc.Hash)
}

// writeChangePayload writes the serialised payload of a ChangeChunk to w.
func writeChangePayload(w io.Writer, cc *ChangeChunk, ops *ChangeOpsWriter) error {
	// Dependencies
	if _, err := w.Write(leb128.EncodeU64(uint64(len(cc.Dependencies)))); err != nil {
		return err
	}
	for _, dep := range cc.Dependencies {
		if _, err := w.Write(dep[:]); err != nil {
			return err
		}
	}

	// Actor (length-prefixed bytes)
	if _, err := w.Write(leb128.EncodeU64(uint64(len(cc.Actor)))); err != nil {
		return err
	}
	if _, err := w.Write(cc.Actor); err != nil {
		return err
	}

	// SeqNum, StartOp, Time, Message
	if _, err := w.Write(leb128.EncodeU64(cc.SeqNum)); err != nil {
		return err
	}
	if _, err := w.Write(leb128.EncodeU64(cc.StartOp)); err != nil {
		return err
	}
	if _, err := w.Write(leb128.EncodeS64(int64(cc.Time))); err != nil {
		return err
	}
	if _, err := w.Write(leb128.EncodeU64(uint64(len(cc.Message)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(cc.Message)); err != nil {
		return err
	}

	// OtherActors: count then each as (length, bytes)
	if _, err := w.Write(leb128.EncodeU64(uint64(len(cc.OtherActors)))); err != nil {
		return err
	}
	for _, a := range cc.OtherActors {
		if _, err := w.Write(leb128.EncodeU64(uint64(len(a)))); err != nil {
			return err
		}
		if _, err := w.Write(a); err != nil {
			return err
		}
	}

	return ops.writePayloadColumns(w)
}

// ChangeOpsWriter streams operations into per-column writers for a change chunk.
// Call Append for each operation, then flush once, then pass to WriteChange.
type ChangeOpsWriter struct {
	objActorBuf, objCtrBuf               bytes.Buffer // obj
	keyActorBuf, keyCtrBuf, keyStrBuf    bytes.Buffer // key
	insertBuf                            bytes.Buffer // insert
	actionBuf, valueMetaBuf, valueBuf    bytes.Buffer // action
	predGrpBuf, predActorBuf, predCtrBuf bytes.Buffer // preds

	obj    *column.ObjectWriter
	key    *column.KeyWriter
	insert *column.InsertWriter
	action *column.ActionWriter
	preds  *column.GroupedOpIdWriter
}

func NewChangeOpsWriter() *ChangeOpsWriter {
	w := &ChangeOpsWriter{}
	w.obj = column.NewObjectWriter(&w.objActorBuf, &w.objCtrBuf)
	w.key = column.NewKeyWriter(&w.keyActorBuf, &w.keyCtrBuf, &w.keyStrBuf)
	w.insert = column.NewInsertWriter(&w.insertBuf)
	w.action = column.NewActionWriter(&w.actionBuf, &w.valueMetaBuf, &w.valueBuf)
	w.preds = column.NewGroupedOpIdWriter(&w.predGrpBuf, &w.predActorBuf, &w.predCtrBuf)
	return w
}

// Append encodes one operation into the per-column writers.
// localOf maps global actor indices to local indices in the change.
func (w *ChangeOpsWriter) Append(obj types.ObjectId, key types.Key, insert bool, action types.Action, preds []types.OpId, localOf map[uint32]uint32) {
	w.obj.Append(obj, localOf)
	w.key.Append(key, localOf)
	w.insert.Append(insert)
	w.action.Append(action)
	w.preds.Append(preds, localOf)
}

// flush flushes all column writers.
func (w *ChangeOpsWriter) flush() error {
	for _, f := range []interface{ Flush() error }{
		w.obj, w.key, w.insert, w.action, w.preds,
	} {
		if err := f.Flush(); err != nil {
			return fmt.Errorf("ops writer finalise: %w", err)
		}
	}
	return nil
}

// writePayloadColumns writes the column section of a change payload to out:
// column count (LEB128), then (spec, length) metadata pairs, then column bytes
// — all in ascending spec order. Only non-empty columns are included.
func (w *ChangeOpsWriter) writePayloadColumns(out io.Writer) error {
	if err := w.flush(); err != nil {
		return fmt.Errorf("ops writer finalize: %w", err)
	}

	type col struct {
		spec uint32
		data []byte
	}
	var cols []col
	add := func(spec uint32, data []byte) {
		cols = append(cols, col{spec, data})
	}

	if w.obj.HasNonRoot() {
		add(colObjActor, w.objActorBuf.Bytes())
		add(colObjCtr, w.objCtrBuf.Bytes())
	}
	if w.key.HasOpId() {
		if w.key.HasNonNullActor() {
			add(colKeyActor, w.keyActorBuf.Bytes())
		}
		add(colKeyCtr, w.keyCtrBuf.Bytes())
	}
	if w.key.HasString() {
		add(colKeyStr, w.keyStrBuf.Bytes())
	}
	// insert, value metadata, and predecessor group are written unconditionally
	// even when trivially empty (all-false / all-null / all-zero). The Rust
	// reference implementation always emits them, so omitting them would produce
	// a different hash for the same logical content, breaking cross-implementation
	// compatibility.
	add(colInsert, w.insertBuf.Bytes())
	add(colAction, w.actionBuf.Bytes())
	add(colValMeta, w.valueMetaBuf.Bytes())
	if w.action.HasValues() {
		add(colVal, w.valueBuf.Bytes())
	}
	add(colPredGrp, w.predGrpBuf.Bytes())
	if w.preds.HasPreds() {
		add(colPredActor, w.predActorBuf.Bytes())
		add(colPredCtr, w.predCtrBuf.Bytes())
	}

	if _, err := out.Write(leb128.EncodeU64(uint64(len(cols)))); err != nil {
		return err
	}
	for _, c := range cols {
		if _, err := out.Write(leb128.EncodeU32(c.spec)); err != nil {
			return err
		}
		if _, err := out.Write(leb128.EncodeU64(uint64(len(c.data)))); err != nil {
			return err
		}
	}
	for _, c := range cols {
		if _, err := out.Write(c.data); err != nil {
			return err
		}
	}
	return nil
}
