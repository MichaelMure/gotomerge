package format

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/MichaelMure/leb128"

	"github.com/MichaelMure/gotomerge/column"
	"github.com/MichaelMure/gotomerge/types"
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

	return ops.writeColumns(w)
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
// m maps global actor indices to local indices in the change.
func (w *ChangeOpsWriter) Append(obj types.ObjectId, key types.Key, insert bool, action types.Action, preds []types.OpId, mapper types.ActorMapper) {
	w.obj.Append(obj, mapper)
	w.key.Append(key, mapper)
	w.insert.Append(insert)
	w.action.Append(action)
	w.preds.Append(preds, mapper)
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

// writeColumns writes the column section of a change payload to out:
// column count (LEB128), then (spec, length) metadata pairs, then column bytes
// — all in ascending spec order. Only non-empty columns are included.
func (w *ChangeOpsWriter) writeColumns(out io.Writer) error {
	if err := w.flush(); err != nil {
		return fmt.Errorf("ops writer finalize: %w", err)
	}

	type col struct {
		spec uint32
		data []byte
	}
	var cols []col

	// Each column is included iff its buffer is non-empty after flushing.
	// An all-null (or all-false for bool) column produces zero bytes and is
	// omitted, matching Rust's filter(!c.data.is_empty()). INSERT is a bool
	// column but false is a non-null value, so it always produces bytes when
	// ops exist — no special-casing needed unlike Rust's write_unless_empty.
	for _, pair := range []struct {
		spec   int
		buffer bytes.Buffer
	}{
		{colObjActor, w.objActorBuf},
		{colObjCtr, w.objCtrBuf},
		{colKeyActor, w.keyActorBuf},
		{colKeyCtr, w.keyCtrBuf},
		{colKeyStr, w.keyStrBuf},
		{colInsert, w.insertBuf},
		{colAction, w.actionBuf},
		{colValMeta, w.valueMetaBuf},
		{colVal, w.valueBuf},
		{colPredGrp, w.predGrpBuf},
		{colPredActor, w.predActorBuf},
		{colPredCtr, w.predCtrBuf},
	} {
		if pair.buffer.Len() > 0 {
			cols = append(cols, col{uint32(pair.spec), pair.buffer.Bytes()})
		}
	}

	sort.Slice(cols, func(i, j int) bool { return cols[i].spec < cols[j].spec })

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
