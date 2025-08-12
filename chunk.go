package gotomerge

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"iter"
	"strings"

	"github.com/jcalabro/leb128"

	"gotomerge/column"
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
	r = io.LimitReader(r, int64(length))

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

	skip := func(t any, l uint64) {
		fmt.Printf("SKIP: %s %v\n", t, l)
		_, err := io.CopyN(io.Discard, r, int64(l))
		if err != nil {
			panic(err)
		}
	}

	res.Changes = make([][]any, len(res.ChangeMetadata))
	for i, metadatum := range res.ChangeMetadata {
		switch metadatum.Spec.Type() {
		case column.TypeGroup:
			skip(metadatum.Spec, metadatum.Length)
		case column.TypeActor:
			skip(metadatum.Spec, metadatum.Length)
		case column.TypeULEB128:
			res.Changes[i] = acc(column.ReadUlebColumn(r, metadatum.Length))
		case column.TypeDelta:
			res.Changes[i] = acc(column.ReadDeltaColumn(r, metadatum.Length))
		case column.TypeBool:
			res.Changes[i] = acc(column.ReadBooleanColumn(r, metadatum.Length))
		case column.TypeString:
			res.Changes[i] = acc(column.ReadStringColumn(r, metadatum.Length))
		case column.TypeValueMetadata:
			res.Changes[i] = acc(column.ReadValueMetadataColumn(r, metadatum.Length))
		case column.TypeValue:
			skip(metadatum.Spec, metadatum.Length)
		}
	}

	res.Operations = make([][]any, len(res.OperationMetadata))
	for i, metadatum := range res.OperationMetadata {
		switch metadatum.Spec.Type() {
		case column.TypeGroup:
			skip(metadatum.Spec, metadatum.Length)
		case column.TypeActor:
			skip(metadatum.Spec, metadatum.Length)
		case column.TypeULEB128:
			res.Operations[i] = acc(column.ReadUlebColumn(r, metadatum.Length))
		case column.TypeDelta:
			res.Operations[i] = acc(column.ReadDeltaColumn(r, metadatum.Length))
		case column.TypeBool:
			res.Operations[i] = acc(column.ReadBooleanColumn(r, metadatum.Length))
		case column.TypeString:
			res.Operations[i] = acc(column.ReadStringColumn(r, metadatum.Length))
		case column.TypeValueMetadata:
			res.Operations[i] = acc(column.ReadValueMetadataColumn(r, metadatum.Length))
		case column.TypeValue:
			skip(metadatum.Spec, metadatum.Length)
		}
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
