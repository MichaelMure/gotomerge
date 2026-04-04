package format

import (
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/MichaelMure/leb128"

	"gotomerge/types"
)

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

// writeChangeChunkCompressed DEFLATE-compresses payload and writes it as a
// ChunkTypeCompressedChange. The canonical hash — stored in hash and used
// as the chunk's identity in the dependency graph — is computed over the
// uncompressed form (as ChunkTypeChange), so compressed and uncompressed
// peers agree on the same hash.
func writeChangeChunkCompressed(w io.Writer, payload []byte, hash *types.ChangeHash) error {
	h := sha256.New()
	h.Write([]byte{byte(ChunkTypeChange)})
	h.Write(leb128.EncodeU64(uint64(len(payload))))
	h.Write(payload)
	digest := h.Sum(nil)
	if hash != nil {
		copy(hash[:], digest)
	}

	var compressed bytes.Buffer
	fw, err := flate.NewWriter(&compressed, flate.DefaultCompression)
	if err != nil {
		return fmt.Errorf("compressed change: %w", err)
	}
	_, err = fw.Write(payload)
	if cerr := fw.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return fmt.Errorf("compressed change: %w", err)
	}

	if _, err = w.Write(magicBytes); err != nil {
		return err
	}
	if _, err = w.Write(digest[:4]); err != nil {
		return err
	}
	if _, err = w.Write([]byte{byte(ChunkTypeCompressedChange)}); err != nil {
		return err
	}
	if _, err = w.Write(leb128.EncodeU64(uint64(compressed.Len()))); err != nil {
		return err
	}
	_, err = w.Write(compressed.Bytes())
	return err
}

func writeChangeHashes(w io.Writer, hashes []types.ChangeHash) error {
	_, err := w.Write(leb128.EncodeU64(uint64(len(hashes))))
	if err != nil {
		return err
	}
	for _, h := range hashes {
		_, err = w.Write(h[:])
		if err != nil {
			return err
		}
	}
	return nil
}

func writeActorIds(w io.Writer, ids []types.ActorId) error {
	_, err := w.Write(leb128.EncodeU64(uint64(len(ids))))
	if err != nil {
		return err
	}
	for _, id := range ids {
		_, err = w.Write(leb128.EncodeU64(uint64(len(id))))
		if err != nil {
			return err
		}
		_, err = w.Write(id)
		if err != nil {
			return err
		}
	}
	return nil
}

func writeHeadIndexes(w io.Writer, indexes []uint64) error {
	for _, index := range indexes {
		_, err := w.Write(leb128.EncodeU64(index))
		if err != nil {
			return err
		}
	}
	return nil
}
