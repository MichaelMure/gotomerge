package format

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"

	"github.com/MichaelMure/leb128"

	"gotomerge/column"
	"gotomerge/types"
)

// deflateMinSize is the threshold above which individual columns are
// DEFLATE-compressed, matching the Rust reference (DEFLATE_MIN_SIZE = 256).
const deflateMinSize = 256

// WriteDocument serialises a document chunk. Pass nil for changes to omit the
// change metadata section (valid, but omits sync metadata).
func WriteDocument(w io.Writer, actors []types.ActorId, heads []types.ChangeHash, headIndexes []uint64, changes *ChangeMetaWriter, ops *DocOpsWriter) error {
	var payload bytes.Buffer
	if err := writeDocumentPayload(&payload, actors, heads, headIndexes, changes, ops); err != nil {
		return err
	}
	return writeChunk(w, ChunkTypeDocument, payload.Bytes(), nil)
}

func writeDocumentPayload(w io.Writer, actors []types.ActorId, heads []types.ChangeHash, headIndexes []uint64, changes *ChangeMetaWriter, ops *DocOpsWriter) error {
	if err := writeActorIds(w, actors); err != nil {
		return err
	}
	if err := writeChangeHashes(w, heads); err != nil {
		return err
	}

	// Wire layout: change_meta_header | op_meta_header | change_data | op_data
	var chgMeta, chgData bytes.Buffer
	if changes != nil {
		if err := changes.encodeSection(&chgMeta, &chgData); err != nil {
			return err
		}
	} else {
		chgMeta.Write(leb128.EncodeU64(0))
	}
	var opMeta, opData bytes.Buffer
	if err := ops.encodeSection(&opMeta, &opData); err != nil {
		return err
	}
	for _, b := range [][]byte{chgMeta.Bytes(), opMeta.Bytes(), chgData.Bytes(), opData.Bytes()} {
		if _, err := w.Write(b); err != nil {
			return err
		}
	}
	return writeHeadIndexes(w, headIndexes)
}

// ChangeMetaWriter streams per-change metadata into column writers for a document chunk.
// Call Append for each change in order, then pass to WriteDocument.
type ChangeMetaWriter struct {
	actorBuf, seqBuf, maxOpBuf, timeBuf bytes.Buffer // changes
	msgBuf, depsGrpBuf, depsIdxBuf      bytes.Buffer // changes
	extraMetaBuf                        bytes.Buffer // extraMeta

	changes   *column.ChangesWriter
	extraMeta *column.ValueMetadataWriter
}

// NewChangeMetaWriter creates a ChangeMetaWriter ready to accept Append calls.
func NewChangeMetaWriter() *ChangeMetaWriter {
	w := &ChangeMetaWriter{}
	w.changes = column.NewChangesWriter(
		&w.actorBuf, &w.seqBuf, &w.maxOpBuf, &w.timeBuf,
		&w.msgBuf, &w.depsGrpBuf, &w.depsIdxBuf,
	)
	w.extraMeta = column.NewValueMetadataWriter(&w.extraMetaBuf)
	return w
}

// Append encodes one change's metadata. m.ActorIdx must be an index into the
// document's actor table (same mapping used for op columns).
func (w *ChangeMetaWriter) Append(m column.RawChangeMeta) {
	w.changes.Append(m)
	// Always emit a "0 bytes, bytes type" extra-metadata entry per change,
	// matching the Rust reference encoder which always writes this column.
	w.extraMeta.Append(column.NewValueMetadata(column.ValueTypeBytes, 0))
}

// encodeSection writes the change metadata column section to metaBuf and dataBuf.
// Change columns are never DEFLATE-compressed (matching the Rust reference).
func (w *ChangeMetaWriter) encodeSection(metaBuf, dataBuf *bytes.Buffer) error {
	if err := w.changes.Flush(); err != nil {
		return fmt.Errorf("change meta writer flush: %w", err)
	}
	if err := w.extraMeta.Flush(); err != nil {
		return fmt.Errorf("change meta writer flush extra: %w", err)
	}

	type col struct {
		spec uint32
		data []byte
	}
	var cols []col
	add := func(spec uint32, data []byte) {
		cols = append(cols, col{spec, data})
	}

	add(colDocChgActor, w.actorBuf.Bytes())
	add(colDocChgSeqNum, w.seqBuf.Bytes())
	add(colDocChgMaxOp, w.maxOpBuf.Bytes())
	if w.changes.HasTime() {
		add(colDocChgTime, w.timeBuf.Bytes())
	}
	if w.changes.HasMessage() {
		add(colDocChgMessage, w.msgBuf.Bytes())
	}
	add(colDocChgDepsGrp, w.depsGrpBuf.Bytes())
	if w.changes.HasDeps() {
		add(colDocChgDepsIdx, w.depsIdxBuf.Bytes())
	}
	add(colValMeta, w.extraMetaBuf.Bytes())

	metaBuf.Write(leb128.EncodeU64(uint64(len(cols))))
	for _, c := range cols {
		metaBuf.Write(leb128.EncodeU32(c.spec))
		metaBuf.Write(leb128.EncodeU64(uint64(len(c.data))))
	}
	for _, c := range cols {
		dataBuf.Write(c.data)
	}
	return nil
}

// DocOpsWriter streams operations into per-column writers for a document chunk.
// Unlike ChangeOpsWriter, it encodes an explicit OpId (actorId + counter) for
// each operation and writes a successor list instead of a predecessor list.
// Call Append for each operation in object order, then pass to WriteDocument.
type DocOpsWriter struct {
	objActorBuf, objCtrBuf               bytes.Buffer // obj
	keyActorBuf, keyCtrBuf, keyStrBuf    bytes.Buffer // key
	opActorBuf, opCtrBuf                 bytes.Buffer // opId
	insertBuf                            bytes.Buffer // insert
	actionBuf, valueMetaBuf, valueBuf    bytes.Buffer // action
	succGrpBuf, succActorBuf, succCtrBuf bytes.Buffer // succs

	obj    *column.ObjectWriter
	key    *column.KeyWriter
	opId   *column.OpIdWriter
	insert *column.InsertWriter
	action *column.ActionWriter
	succs  *column.GroupedOpIdWriter
}

// NewDocOpsWriter creates a DocOpsWriter ready to accept Append calls.
func NewDocOpsWriter() *DocOpsWriter {
	w := &DocOpsWriter{}
	w.obj = column.NewObjectWriter(&w.objActorBuf, &w.objCtrBuf)
	w.key = column.NewKeyWriter(&w.keyActorBuf, &w.keyCtrBuf, &w.keyStrBuf)
	w.opId = column.NewOpIdWriter(&w.opActorBuf, &w.opCtrBuf)
	w.insert = column.NewInsertWriter(&w.insertBuf)
	w.action = column.NewActionWriter(&w.actionBuf, &w.valueMetaBuf, &w.valueBuf)
	w.succs = column.NewGroupedOpIdWriter(&w.succGrpBuf, &w.succActorBuf, &w.succCtrBuf)
	return w
}

// Append encodes one operation into the per-column writers.
// m maps global actor indices to local indices in the document chunk.
func (w *DocOpsWriter) Append(obj types.ObjectId, key types.Key, id types.OpId, insert bool, action types.Action, succs []types.OpId, mapper types.ActorMapper) {
	w.obj.Append(obj, mapper)
	w.key.Append(key, mapper)
	w.opId.Append(id, mapper)
	w.insert.Append(insert)
	w.action.Append(action)
	w.succs.Append(succs, mapper)
}

// flush finalises all column writers.
func (w *DocOpsWriter) flush() error {
	for _, f := range []interface{ Flush() error }{
		w.obj, w.key, w.opId, w.insert, w.action, w.succs,
	} {
		if err := f.Flush(); err != nil {
			return fmt.Errorf("doc ops writer flush: %w", err)
		}
	}
	return nil
}

// encodeSection encodes all non-empty op columns into metaBuf and dataBuf.
// metaBuf receives: count (LEB128) then N×(spec LEB128, length LEB128).
// dataBuf receives the raw (or deflate-compressed) column bytes.
// Columns larger than deflateMinSize are compressed; their spec has the deflate bit set.
// Specs are written in ascending order as required by the format.
// Must be called after flush.
func (w *DocOpsWriter) encodeSection(metaBuf, dataBuf *bytes.Buffer) error {
	if err := w.flush(); err != nil {
		return fmt.Errorf("doc ops writer finalize: %w", err)
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
	add(colDocOpActor, w.opActorBuf.Bytes())
	add(colDocOpCtr, w.opCtrBuf.Bytes())
	add(colInsert, w.insertBuf.Bytes())
	add(colAction, w.actionBuf.Bytes())
	add(colValMeta, w.valueMetaBuf.Bytes())
	if w.action.HasValues() {
		add(colVal, w.valueBuf.Bytes())
	}
	add(colDocSuccGrp, w.succGrpBuf.Bytes())
	if w.succs.HasPreds() {
		add(colDocSuccActor, w.succActorBuf.Bytes())
		add(colDocSuccCtr, w.succCtrBuf.Bytes())
	}

	metaBuf.Write(leb128.EncodeU64(uint64(len(cols))))
	for _, c := range cols {
		if len(c.data) > deflateMinSize {
			var compressed bytes.Buffer
			fw, err := flate.NewWriter(&compressed, flate.DefaultCompression)
			if err != nil {
				return err
			}
			if _, err = fw.Write(c.data); err != nil {
				_ = fw.Close()
				return err
			}
			if err = fw.Close(); err != nil {
				return err
			}
			spec := column.Specification(c.spec).WithDeflate()
			metaBuf.Write(leb128.EncodeU32(uint32(spec)))
			metaBuf.Write(leb128.EncodeU64(uint64(compressed.Len())))
			dataBuf.Write(compressed.Bytes())
		} else {
			metaBuf.Write(leb128.EncodeU32(c.spec))
			metaBuf.Write(leb128.EncodeU64(uint64(len(c.data))))
			dataBuf.Write(c.data)
		}
	}
	return nil
}
