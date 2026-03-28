package format

import (
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/jcalabro/leb128"

	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

type ChunkType byte

const (
	ChunkTypeDocument         ChunkType = 0x00
	ChunkTypeChange           ChunkType = 0x01
	ChunkTypeCompressedChange ChunkType = 0x02
)

var magicBytes = []byte{0x85, 0x6f, 0x4a, 0x83}

// Chunk is the top-level unit of an Automerge binary file.
//
// An Automerge file is a sequence of one or more chunks. Each chunk is either a
// document snapshot (DocumentChunk) or a single peer edit (ChangeChunk). A file
// with just one DocumentChunk is a compact, fully-merged snapshot. A file may
// also hold a series of ChangeChunks representing the edit history of a document,
// in which case the reader applies them in dependency order to reconstruct the
// current state.
//
// The only concrete implementations are *DocumentChunk and *ChangeChunk. The
// unexported chunk() method prevents other packages from satisfying this interface,
// making those two the exhaustive set of variants. Callers must type-assert to the
// concrete type to access chunk-specific fields.
type Chunk interface {
	chunk()
}

// rawColumn holds the raw bytes of an operation or change column whose spec number
// is not recognised by this implementation. Storing it intact allows the file to be
// round-tripped without data loss if a newer format version adds column types we
// don't yet know about.
type rawColumn struct {
	specBits uint32 // encoded column spec; high bits are the column ID, low bits are the type
	data     []byte
}

// ReadChunk reads one chunk from r. The second return value is the total number
// of payload bytes belonging to this chunk. It is returned even on error so the
// caller can decide whether to skip and continue.
//
// The file format for each chunk is:
//
//	[4 magic][4 checksum][1 type][varint length][...payload...]
//
// The checksum is the first 4 bytes of SHA-256(type || length || payload).
// ReadChunk verifies the checksum before returning.
func ReadChunk(r ioutil.SubReader) (Chunk, int, error) {
	// Open a sub-reader so we can track exactly how many bytes we consume
	// and return an accurate skip count independent of parsing success.
	r, err := r.SubReaderOffset(0)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading chunk: %w", err)
	}

	// reading buffer of 16 bytes, but sized to read magic+checksum
	buf := make([]byte, 4+4, 16)
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return nil, r.Consumed(), fmt.Errorf("error reading chunk header: %w", err)
	}
	if n != (4 + 4) {
		return nil, r.Consumed(), fmt.Errorf("unexpected end of chunk")
	}
	if !bytes.Equal(buf[0:4], magicBytes) {
		return nil, r.Consumed(), fmt.Errorf("invalid magic bytes")
	}
	checksum := buf[4:8]

	// Run the remaining reads through SHA256 to compute the checksum
	h := sha256.New()
	tee := io.TeeReader(r, h)

	// reuse the allocated buffer to read a single byte
	buf = buf[0:1]
	n, err = io.ReadFull(tee, buf)
	if err != nil {
		return nil, r.Consumed(), fmt.Errorf("error reading chunk type: %w", err)
	}
	if n != 1 {
		return nil, r.Consumed(), fmt.Errorf("unexpected end of chunk")
	}
	chunkType := ChunkType(buf[0])

	length, err := leb128.DecodeU64(tee)
	if err != nil {
		return nil, r.Consumed(), fmt.Errorf("error reading chunk length: %w", err)
	}
	toSkip := r.Consumed() + int(length)
	sub, err := r.SubReader(0, length)
	if err != nil {
		return nil, toSkip, fmt.Errorf("error reading chunk: %w", err)
	}

	var res Chunk

	switch chunkType {
	case ChunkTypeDocument:
		err = r.Peek(h, int(length))
		if err != nil {
			return nil, toSkip, fmt.Errorf("error hashing chunk: %w", err)
		}
		res, err = readDocumentChunk(sub)
	case ChunkTypeChange:
		err = r.Peek(h, int(length))
		if err != nil {
			return nil, toSkip, fmt.Errorf("error hashing chunk: %w", err)
		}
		res, err = readChangeChunk(sub)
	case ChunkTypeCompressedChange:
		// A compressed change has the same content as a plain change, but the payload
		// is DEFLATE-compressed. Crucially, the *canonical hash* of a compressed change
		// is defined as the hash of its *uncompressed* equivalent, not the compressed bytes.
		// This means two peers — one that stored the change compressed and one that stored
		// it uncompressed — will agree on the hash and therefore on the change's identity.
		//
		// The downside: because the uncompressed length is only known after full decompression,
		// we cannot stream-hash this chunk the way we do for plain changes. We must decompress
		// entirely into memory and then hash type || decompressed-length || decompressed-bytes.
		//
		// Benchmarks show this costs roughly +11% time, +7% allocation count, +2% allocation size
		// compared to plain changes.
		decompressed, err := io.ReadAll(flate.NewReader(sub))
		if err != nil {
			return nil, toSkip, fmt.Errorf("error reading compressed change: %w", err)
		}

		h = sha256.New()
		_, err = h.Write([]byte{byte(ChunkTypeChange)})
		if err != nil {
			return nil, toSkip, err
		}
		_, err = h.Write(leb128.EncodeU64(uint64(len(decompressed))))
		if err != nil {
			return nil, toSkip, err
		}
		_, err = h.Write(decompressed)
		if err != nil {
			return nil, toSkip, err
		}

		res, err = readChangeChunk(ioutil.NewBytesReader(decompressed))
	default:
		return nil, toSkip, fmt.Errorf("invalid chunk type: %d", chunkType)
	}
	if err != nil {
		return nil, toSkip, fmt.Errorf("error reading chunk: %w", err)
	}

	digest := h.Sum(nil)
	if !bytes.Equal(checksum, digest[0:4]) {
		return nil, toSkip, fmt.Errorf("invalid checksum")
	}

	// A change's identity in the dependency graph is its full 32-byte hash.
	// We have already computed it for checksum verification, so store it in
	// the chunk rather than requiring callers to recompute it later.
	// DocumentChunk has no hash identity in the protocol, so we only do this for changes.
	if cc, ok := res.(*ChangeChunk); ok {
		copy(cc.Hash[:], digest)
	}

	return res, toSkip, nil
}

func readChangeHashes(r io.Reader) ([]types.ChangeHash, error) {
	n, err := leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading change hash length: %w", err)
	}
	if n == 0 {
		return nil, nil
	}
	// limit pre-allocation to avoid DOS
	allocate := n
	// TODO: adjust with reasonable value
	if n > 128 {
		allocate = 128
	}
	res := make([]types.ChangeHash, 0, allocate)

	for i := uint64(0); i < n; i++ {
		var h types.ChangeHash
		_, err = io.ReadFull(r, h[:])
		if err != nil {
			return nil, fmt.Errorf("error reading change hash: %w", err)
		}
		res = append(res, h)
	}
	return res, nil
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

func readActorIds(r io.Reader) ([]types.ActorId, error) {
	n, err := leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading actor ids length: %w", err)
	}
	if n == 0 {
		return nil, nil
	}
	// limit pre-allocation to avoid DOS
	allocate := n
	// TODO: tune with reasonable value
	if n > 128 {
		allocate = 128
	}
	res := make([]types.ActorId, 0, allocate)

	var prevId types.ActorId
	for i := uint64(0); i < n; i++ {
		id, err := types.ReadLengthEncodedActorId(r)
		if err != nil {
			return nil, fmt.Errorf("error reading actor id: %w", err)
		}
		if bytes.Compare(id, prevId) <= 0 {
			return nil, fmt.Errorf("actor IDs must be sorted")
		}
		prevId = id
		res = append(res, id)
	}
	return res, nil
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

func readHeadIndexes(r io.Reader, count int) ([]uint64, error) {
	res := make([]uint64, count)
	for i := 0; i < count; i++ {
		index, err := leb128.DecodeU64(r)
		if err != nil {
			if i == 0 && err == io.EOF {
				res = nil
				// HeadIndexes were added in a later revision of the format.
				// Pre-existing documents may not have them; treat EOF as "no indexes".
				break
			}
			return nil, err
		}
		res[i] = index
	}
	return res, nil
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
