package gotomerge

import (
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"fmt"
	"io"
	"iter"

	"github.com/jcalabro/leb128"

	"gotomerge/types"
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
		res, err = readDocumentChunk(r)
	case ChunkTypeChange:
		res, err = readChangeChunk(r)
	case ChunkTypeCompressedChange:
		r = flate.NewReader(r)
		res, err = readDocumentChunk(r)
	default:
		return nil, fmt.Errorf("invalid chunk type: %d", chunkType)
	}
	if err != nil {
		return nil, fmt.Errorf("error reading chunk: %w", err)
	}

	// TODO: remove
	// rest, err := io.ReadAll(r)
	// if err != nil {
	// 	panic("")
	// }
	// fmt.Printf("REST READ: %v\n", len(rest))

	if !bytes.Equal(checksum, h.Sum(nil)[0:4]) {
		return nil, fmt.Errorf("invalid checksum")
	}

	return res, nil
}

// TODO: remove
func acc[T any](it iter.Seq2[T, error]) []any {
	var res []any
	for t, err := range it {
		if err != nil {
			panic(err)
		}
		res = append(res, t)
	}
	return res
}

func readChangeHashes(r io.Reader) ([]types.ChangeHash, error) {
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
	// limit pre-allocation to avoid DOS
	allocate := n
	// TODO: tune with reasonable value
	if n > 128 {
		allocate = 128
	}
	res := make([]types.ActorId, 0, allocate)

	for i := uint64(0); i < n; i++ {
		id, err := types.ReadLengthEncodedActorId(r)
		if err != nil {
			return nil, fmt.Errorf("error reading actor id: %w", err)
		}
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
				// old document may not have this
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
