package gotomerge

import (
	"bytes"
	"compress/flate"
	"crypto/sha256"
	"fmt"
	"io"
	"iter"
	"strings"

	"github.com/jcalabro/leb128"

	"gotomerge/column"
	"gotomerge/column/rle"
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

type DocumentChunk struct {
	Actors      []types.ActorId
	Heads       []types.ChangeHash
	HeadIndexes []uint64

	// TODO: should that stays?
	ChangeMetadata column.Metadata
	Changes        [][]any

	OperationMetadata column.Metadata
	Operations        [][]any
}

func (d DocumentChunk) String() string {
	var res strings.Builder
	res.WriteString("DocumentChunk {\n")
	res.WriteString(fmt.Sprintf("  Actors: %v\n", d.Actors))
	res.WriteString(fmt.Sprintf("  Heads: %v\n", d.Heads))
	res.WriteString(fmt.Sprintf("  HeadIndexes: %v\n", d.HeadIndexes))
	for i, metadatum := range d.ChangeMetadata {
		res.WriteString(fmt.Sprintf("  ChangeMetadata[%d]: %v\n", i, metadatum))
		res.WriteString(fmt.Sprintf("    Values: %v\n", d.Changes[i]))
	}
	for i, metadatum := range d.OperationMetadata {
		res.WriteString(fmt.Sprintf("  OperationMetadata[%d]: %v\n", i, metadatum))
		res.WriteString(fmt.Sprintf("    Values: %v\n", d.Operations[i]))
	}
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
		res, err = readDocumentChunk(r)
	case ChunkTypeChange:
		panic("not implemented")
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

func readDocumentChunk(r io.Reader) (*DocumentChunk, error) {
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

	res.ChangeMetadata, err = column.ReadMetadata(r)
	if err != nil {
		return nil, fmt.Errorf("error reading change metadata: %w", err)
	}

	res.OperationMetadata, err = column.ReadMetadata(r)
	if err != nil {
		return nil, fmt.Errorf("error reading operation metadata: %w", err)
	}

	// TODO: remove
	skip := func(r io.Reader, t any, l uint64) {
		fmt.Printf("SKIP: %s %v\n", t, l)
		_, err := io.CopyN(io.Discard, r, int64(l))
		if err != nil {
			panic(err)
		}
	}

	res.Changes = make([][]any, len(res.ChangeMetadata))
	for i, metadatum := range res.ChangeMetadata {
		rCol := io.LimitReader(r, int64(metadatum.Length))
		if metadatum.Spec.Deflate() {
			rCol = flate.NewReader(rCol)
		}

		switch metadatum.Spec.Type() {
		case column.TypeGroup:
			res.Changes[i] = acc(rle.ReadUint64RLE(rCol))
		case column.TypeActor:
			res.Changes[i] = acc(rle.ReadUint64RLE(rCol))
		case column.TypeULEB128:
			res.Changes[i] = acc(column.ReadUlebColumn(rCol))
		case column.TypeDelta:
			res.Changes[i] = acc(column.ReadDeltaColumn(rCol))
		case column.TypeBool:
			res.Changes[i] = acc(column.ReadBooleanColumn(rCol))
		case column.TypeString:
			res.Changes[i] = acc(column.ReadStringColumn(rCol))
		case column.TypeValueMetadata:
			res.Changes[i] = acc(column.ReadValueMetadataColumn(rCol))
		case column.TypeValue:
			skip(rCol, metadatum.Spec, metadatum.Length)
		}
	}

	var prevValueMetadata []column.ValueMetadata

	res.Operations = make([][]any, len(res.OperationMetadata))
	for i, metadatum := range res.OperationMetadata {
		rCol := io.LimitReader(r, int64(metadatum.Length))
		if metadatum.Spec.Deflate() {
			rCol = flate.NewReader(rCol)
		}

		switch metadatum.Spec.Type() {
		case column.TypeGroup:
			res.Operations[i] = acc(rle.ReadUint64RLE(rCol))
		case column.TypeActor:
			res.Operations[i] = acc(rle.ReadUint64RLE(rCol))
		case column.TypeULEB128:
			res.Operations[i] = acc(column.ReadUlebColumn(rCol))
		case column.TypeDelta:
			res.Operations[i] = acc(column.ReadDeltaColumn(rCol))
		case column.TypeBool:
			res.Operations[i] = acc(column.ReadBooleanColumn(rCol))
		case column.TypeString:
			res.Operations[i] = acc(column.ReadStringColumn(rCol))
		case column.TypeValueMetadata:
			// TODO: HACK just for early visualisation
			buf := &bytes.Buffer{}
			tee := io.TeeReader(rCol, buf)
			res.Operations[i] = acc(column.ReadValueMetadataColumn(tee))
			prevValueMetadata = nil
			for vm, err := range column.ReadValueMetadataColumn(buf) {
				if err != nil {
					panic(err)
				}
				prevValueMetadata = append(prevValueMetadata, vm)
			}
		case column.TypeValue:
			// skip(rCol, metadatum.Spec, metadatum.Length)

			// TODO: HACK just for early visualisation
			res.Operations[i] = acc(column.ReadValueColumn(r, prevValueMetadata))
		}
	}

	res.HeadIndexes, err = readHeadIndexes(r, len(res.Heads))
	if err != nil {
		return nil, fmt.Errorf("error reading head indexes: %w", err)
	}

	return &res, nil
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
	if n > 128 {
		allocate = 128
	}
	res := make([]types.ActorId, 0, allocate)

	for i := uint64(0); i < n; i++ {
		l, err := leb128.DecodeU64(r)
		if err != nil {
			return nil, fmt.Errorf("error reading actor id length: %w", err)
		}
		if l > 32 {
			return nil, fmt.Errorf("unexpectedly large actor id length")
		}
		id := make([]byte, l)
		_, err = io.ReadFull(r, id[:])
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
