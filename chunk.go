package gotomerge

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

type chunk interface {
	// TODO?
}

func readChunk(r ioutil.SubReader) (chunk, int, error) {
	// take a subreader to read without consuming, we'll
	// return how many bytes to skip when done with the chunk
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

	var res chunk

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
		// The compressed change format required that we hash the equivalent uncompressed chunk, that is:
		// - the normal change chunk type (0x01)
		// - the length of the *decompressed* bytes
		// - the decompressed bytes of the remaining chunk (excluding type+size)
		// As we can't know the decompressed size before fully decompressing it, it's not possible to stream decode the
		// chunk like with the other types. A small format change would allow that.
		// Instead, we have to fully decompress in memory and re-do the hashing.
		//
		// Some benchmarking shows that this cost +11% speed, +7% allocation count, +2% allocation size.
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

	// TODO: remove
	// rest, err := io.ReadAll(r)
	// if err != nil {
	// 	panic("")
	// }
	// fmt.Printf("REST READ: %v\n", len(rest))

	if !bytes.Equal(checksum, h.Sum(nil)[0:4]) {
		return nil, toSkip, fmt.Errorf("invalid checksum")
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
