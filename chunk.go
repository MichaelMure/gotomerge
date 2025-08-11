package gotomerge

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"github.com/jcalabro/leb128"
)

type ChunkType byte

const (
	ChunkTypeDocument         ChunkType = 0x00
	ChunkTypeChange           ChunkType = 0x01
	ChunkTypeCompressedChange ChunkType = 0x02
)

var magicBytes = []byte{0x85, 0x6f, 0x4a, 0x83}

type chunk interface {
}

type DocumentChunk struct {
	Actors []ActorId
	Heads  []changeHash
}

func (d DocumentChunk) String() string {
	var res strings.Builder
	res.WriteString("DocumentChunk {\n")
	res.WriteString(fmt.Sprintf("  Actors: %v\n", d.Actors))
	res.WriteString(fmt.Sprintf("  Heads: %v\n", d.Heads))
	res.WriteString("}\n")
	return res.String()
}

func readChunk(r io.Reader) (chunk, error) {
	// reading buffer of 16 bytes, but sized to read magic+checksum
	buf := make([]byte, 4+4, 16)
	n, err := r.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk header: %w", err)
	}
	if n != (4 + 4) {
		return nil, fmt.Errorf("unexpected end of chunk")
	}
	if !bytes.Equal(buf[0:4], magicBytes) {
		return nil, fmt.Errorf("invalid magic bytes")
	}
	checksum := buf[4:8]

	// run the remaining reads through SHA256 to compute the checksum
	h := sha256.New()
	r = io.TeeReader(r, h)

	// reuse the allocated buffer to read a single byte
	buf = buf[0:1]
	n, err = r.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk type: %w", err)
	}
	if n != 1 {
		return nil, fmt.Errorf("unexpected end of chunk")
	}
	chunkType := ChunkType(buf[0])

	length, err := leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading chunk length: %w", err)
	}
	r = io.LimitReader(r, int64(length))

	var res chunk

	switch chunkType {
	case ChunkTypeDocument:
		res, err = readDocumentChunk(r, length)
	case ChunkTypeChange:
		panic("not implemented")
	case ChunkTypeCompressedChange:
		panic("not implemented")
	default:
		return nil, fmt.Errorf("invalid chunk type: %d", chunkType)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading chunk: %w", err)
	}

	// TODO: remove
	rest, err := io.ReadAll(r)
	if err != nil {
		panic("")
	}
	fmt.Printf("REST READ: %v\n", len(rest))

	if !bytes.Equal(checksum, h.Sum(nil)[0:4]) {
		return nil, fmt.Errorf("invalid checksum")
	}

	return res, nil
}

func readDocumentChunk(r io.Reader, length uint64) (*DocumentChunk, error) {
	var res DocumentChunk
	var err error

	res.Actors, err = readActorIds(r)
	if err != nil {
		return nil, fmt.Errorf("error reading actors: %w", err)
	}

	res.Heads, err = readChangeHashes(r)
	if err != nil {
		return nil, fmt.Errorf("error reading heads: %w", err)
	}

	// TODO: columns

	return &res, nil
}

type changeHash [32]byte

func (ch changeHash) String() string {
	return hex.EncodeToString(ch[:])
}

func readChangeHashes(r io.Reader) ([]changeHash, error) {
	n, err := leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading change hash length: %w", err)
	}
	// limit pre-allocation to avoid DOS
	allocate := n
	// TODO: adjust with reasonable value
	if n > 128 {
		allocate = 128
	}
	res := make([]changeHash, 0, allocate)

	for i := uint64(0); i < n; i++ {
		var h changeHash
		_, err = io.ReadFull(r, h[:])
		if err != nil {
			return nil, fmt.Errorf("error reading change hash: %w", err)
		}
		res = append(res, h)
	}
	return res, nil
}

func writeChangeHashes(w io.Writer, hashes []changeHash) error {
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
