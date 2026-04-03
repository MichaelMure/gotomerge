package format

import (
	"bytes"
	"io"

	"github.com/jcalabro/leb128"
)

// WriteDocument serialises dc as a complete Automerge document chunk and
// writes it to w. Column data is written uncompressed regardless of whether
// the original was deflated; the metadata headers reflect this (deflate bit
// cleared).
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
	if err := writeRawMetadata(w, dc.rawChangeColumns); err != nil {
		return err
	}
	if err := writeRawMetadata(w, dc.rawOpColumns); err != nil {
		return err
	}
	for _, col := range dc.rawChangeColumns {
		if _, err := w.Write(col.data); err != nil {
			return err
		}
	}
	for _, col := range dc.rawOpColumns {
		if _, err := w.Write(col.data); err != nil {
			return err
		}
	}
	return writeHeadIndexes(w, dc.HeadIndexes)
}

// writeRawMetadata writes a column metadata header: leb128(count) followed
// by (leb128(spec), leb128(len)) pairs. This mirrors the on-wire format used
// by column.WriteMetadata / column.ReadMetadata.
func writeRawMetadata(w io.Writer, cols []rawColumn) error {
	if _, err := w.Write(leb128.EncodeU64(uint64(len(cols)))); err != nil {
		return err
	}
	for _, col := range cols {
		if _, err := w.Write(leb128.EncodeU32(col.specBits)); err != nil {
			return err
		}
		if _, err := w.Write(leb128.EncodeU64(uint64(len(col.data)))); err != nil {
			return err
		}
	}
	return nil
}
