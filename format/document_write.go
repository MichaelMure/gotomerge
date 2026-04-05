package format

import (
	"bytes"
	"compress/flate"
	"io"

	"github.com/MichaelMure/leb128"

	"gotomerge/column"
	ioutil "gotomerge/utils/io"
)

// deflateMinSize is the threshold above which individual columns are
// DEFLATE-compressed, matching the Rust reference (DEFLATE_MIN_SIZE = 256).
const deflateMinSize = 256

// WriteDocument serialises dc as a complete Automerge document chunk and
// writes it to w. Columns larger than deflateMinSize are compressed and
// written with the deflate bit set in their spec.
func WriteDocument(w io.Writer, dc *DocumentChunk) error {
	var payload bytes.Buffer
	if err := writeDocumentPayload(&payload, dc); err != nil {
		return err
	}
	return writeChunk(w, ChunkTypeDocument, payload.Bytes(), nil)
}

func writeDocumentPayload(w io.Writer, dc *DocumentChunk) error {
	if err := writeActorIds(w, dc.Actors); err != nil {
		return err
	}
	if err := writeChangeHashes(w, dc.Heads); err != nil {
		return err
	}

	// The wire layout is: change metadata, op metadata, change data, op data.
	// We build metadata and data in parallel buffers, then write them in order.
	var changeMeta, changeData, opMeta, opData bytes.Buffer
	if err := encodeColumns(&changeMeta, &changeData, dc.ChangeMetadata, changeColumnReader(&dc.ChangesColumns)); err != nil {
		return err
	}
	if err := encodeColumns(&opMeta, &opData, dc.OpMetadata, opColumnReader(&dc.OpColumns)); err != nil {
		return err
	}

	for _, b := range []*bytes.Buffer{&changeMeta, &opMeta, &changeData, &opData} {
		if _, err := w.Write(b.Bytes()); err != nil {
			return err
		}
	}
	return writeHeadIndexes(w, dc.HeadIndexes)
}

// encodeColumns iterates metadata, reads each SubReader, optionally compresses,
// and writes (spec, len) pairs to meta and raw bytes to data.
func encodeColumns(meta, data *bytes.Buffer, cols column.Metadata, readerFor func(uint32) *ioutil.SubReader) error {
	meta.Write(leb128.EncodeU64(uint64(len(cols))))
	for _, m := range cols {
		spec := uint32(m.Spec) &^ 8 // canonical spec without deflate bit
		sr := readerFor(spec)

		if sr.HasAtLeast(deflateMinSize + 1) {
			// Large column: stream directly into the deflate writer.
			var compressed bytes.Buffer
			fw, err := flate.NewWriter(&compressed, flate.DefaultCompression)
			if err != nil {
				return err
			}
			_, err = io.Copy(fw, sr)
			if cerr := fw.Close(); err == nil {
				err = cerr
			}
			if err != nil {
				return err
			}
			meta.Write(leb128.EncodeU32(spec | 8))
			meta.Write(leb128.EncodeU64(uint64(compressed.Len())))
			data.Write(compressed.Bytes())
		} else {
			// Small column: stream directly; io.Copy returns the byte count.
			n, err := io.Copy(data, sr)
			if err != nil {
				return err
			}
			meta.Write(leb128.EncodeU32(spec))
			meta.Write(leb128.EncodeU64(uint64(n)))
		}
	}
	return nil
}

// changeColumnReader returns a function that maps a column spec to the
// corresponding SubReader in ChangesColumns.
func changeColumnReader(cols *ChangeColumns) func(uint32) *ioutil.SubReader {
	return func(spec uint32) *ioutil.SubReader {
		switch spec {
		case colDocChgActor:
			return cols.ActorId
		case colDocChgSeqNum:
			return cols.SeqNum
		case colDocChgMaxOp:
			return cols.MaxOp
		case colDocChgTime:
			return cols.Time
		case colDocChgMessage:
			return cols.Message
		case colDocChgDepsGrp:
			return cols.DependenciesGroup
		case colDocChgDepsIdx:
			return cols.DependenciesIndex
		case colValMeta:
			return cols.ExtraMetadata
		case colVal:
			return cols.ExtraData
		}
		return nil
	}
}

// opColumnReader returns a function that maps a column spec to the
// corresponding SubReader in OperationColumns.
func opColumnReader(cols *OperationColumns) func(uint32) *ioutil.SubReader {
	return func(spec uint32) *ioutil.SubReader {
		switch spec {
		case colObjActor:
			return cols.ObjectActorId
		case colObjCtr:
			return cols.ObjectCounter
		case colKeyActor:
			return cols.KeyActorId
		case colKeyCtr:
			return cols.KeyCounter
		case colKeyStr:
			return cols.KeyString
		case colDocOpActor:
			return cols.ActorId
		case colDocOpCtr:
			return cols.Counter
		case colInsert:
			return cols.Insert
		case colAction:
			return cols.Action
		case colValMeta:
			return cols.ValueMetadata
		case colVal:
			return cols.Value
		case colPredGrp:
			return cols.PredecessorGroup
		case colPredActor:
			return cols.PredecessorActorId
		case colPredCtr:
			return cols.PredecessorCounter
		case colDocSuccGrp:
			return cols.SuccessorGroup
		case colDocSuccActor:
			return cols.SuccessorActorId
		case colDocSuccCtr:
			return cols.SuccessorCounter
		case colExpandControl:
			return cols.ExpandControl
		case colMark:
			return cols.Mark
		}
		return nil
	}
}
