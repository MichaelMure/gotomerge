package format

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/jcalabro/leb128"

	"gotomerge/column"
	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)


// WriteChange serialises cc as a complete Automerge change chunk and writes it
// to w. It buffers the payload internally (required by the wire format, which
// places the checksum before the payload bytes), computes the SHA-256 hash,
// stores it in cc.Hash, and then writes magic + checksum + type + len + payload.
//
// After WriteChange returns, cc.Hash is set and cc is ready for ApplyChange.
func WriteChange(w io.Writer, cc *ChangeChunk, enc *ChangeOpsWriter) error {
	cc.OpColumns = enc.opColumns()

	var payload bytes.Buffer
	if err := writeChangePayload(&payload, cc, enc.writePayloadColumns); err != nil {
		return err
	}

	return writeChunk(w, ChunkTypeChange, payload.Bytes(), &cc.Hash)
}

// writeChangePayload writes the serialised payload of a ChangeChunk to w.
func writeChangePayload(w io.Writer, cc *ChangeChunk, writeColumns func(io.Writer) error) error {
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

	return writeColumns(w)
}

// writeChunk writes a framed Automerge chunk to w:
// magic(4) + checksum(4) + type(1) + leb128(len) + payload.
// The checksum is the first 4 bytes of sha256(type || leb128(len) || payload).
// If hash is non-nil, the computed digest is stored there.
func writeChunk(w io.Writer, chunkType ChunkType, payload []byte, hash *types.ChangeHash) error {
	h := sha256.New()
	h.Write([]byte{byte(chunkType)})
	h.Write(leb128.EncodeU64(uint64(len(payload))))
	h.Write(payload)
	digest := h.Sum(nil)
	if hash != nil {
		copy(hash[:], digest)
	}

	if _, err := w.Write(magicBytes); err != nil {
		return err
	}
	if _, err := w.Write(digest[:4]); err != nil {
		return err
	}
	if _, err := w.Write([]byte{byte(chunkType)}); err != nil {
		return err
	}
	if _, err := w.Write(leb128.EncodeU64(uint64(len(payload)))); err != nil {
		return err
	}
	_, err := w.Write(payload)
	return err
}

// ChangeOpsWriter streams operations into per-column writers for a change chunk.
// Call Append for each operation, then Finalise once, then pass to WriteChange.
type ChangeOpsWriter struct {
	objActorBuf, objCtrBuf bytes.Buffer
	keyActorBuf, keyCtrBuf bytes.Buffer
	keyStrBuf              bytes.Buffer
	insertBuf              bytes.Buffer
	actionBuf              bytes.Buffer
	valueMetaBuf, valueBuf bytes.Buffer
	predGrpBuf             bytes.Buffer
	predActorBuf           bytes.Buffer
	predCtrBuf             bytes.Buffer

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

// Finalise flushes all column writers. Must be called once before WriteChange.
func (w *ChangeOpsWriter) Finalise() error {
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
	if w.insert.HasInserts() {
		add(colInsert, w.insertBuf.Bytes())
	}
	add(colAction, w.actionBuf.Bytes())
	if w.action.HasValues() {
		add(colValMeta, w.valueMetaBuf.Bytes())
		add(colVal, w.valueBuf.Bytes())
	}
	if w.preds.HasPreds() {
		add(colPredGrp, w.predGrpBuf.Bytes())
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

// opColumns builds OperationColumns SubReaders from the finalised buffers.
// Absent columns are nil, treated as all-null by readers.
func (w *ChangeOpsWriter) opColumns() OperationColumns {
	var oc OperationColumns
	if w.obj.HasNonRoot() {
		oc.ObjectActorId = maybeReader(w.objActorBuf.Bytes())
		oc.ObjectCounter = maybeReader(w.objCtrBuf.Bytes())
	}
	if w.key.HasOpId() {
		if w.key.HasNonNullActor() {
			oc.KeyActorId = maybeReader(w.keyActorBuf.Bytes())
		}
		oc.KeyCounter = maybeReader(w.keyCtrBuf.Bytes())
	}
	if w.key.HasString() {
		oc.KeyString = maybeReader(w.keyStrBuf.Bytes())
	}
	if w.insert.HasInserts() {
		oc.Insert = maybeReader(w.insertBuf.Bytes())
	}
	oc.Action = maybeReader(w.actionBuf.Bytes())
	if w.action.HasValues() {
		oc.ValueMetadata = maybeReader(w.valueMetaBuf.Bytes())
		oc.Value = maybeReader(w.valueBuf.Bytes())
	}
	if w.preds.HasPreds() {
		oc.PredecessorGroup = maybeReader(w.predGrpBuf.Bytes())
		oc.PredecessorActorId = maybeReader(w.predActorBuf.Bytes())
		oc.PredecessorCounter = maybeReader(w.predCtrBuf.Bytes())
	}
	return oc
}

func maybeReader(b []byte) ioutil.SubReader {
	if len(b) == 0 {
		return nil
	}
	return ioutil.NewBytesReader(b)
}
