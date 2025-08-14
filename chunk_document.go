package gotomerge

import (
	"fmt"
	"io"
	"strings"

	"gotomerge/column"
	"gotomerge/column/rle"
	"gotomerge/lbuf"
	"gotomerge/types"
)

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

func readDocumentChunk(r *lbuf.Reader) (*DocumentChunk, error) {
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
		rCol := r
		if metadatum.Spec.Deflate() {
			rCol = rCol.Deflate()
		}
		rCol = r.Limit(int64(metadatum.Length))

		// rCol := r.Limit(int64(metadatum.Length))
		// if metadatum.Spec.Deflate() {
		// 	rCol = rCol.AddProcessor(func(reader io.Reader) io.Reader {
		// 		return flate.NewReader(reader)
		// 	})
		// }

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
		rCol := r
		if metadatum.Spec.Deflate() {
			rCol = rCol.Deflate()
		}
		rCol = rCol.Limit(int64(metadatum.Length))

		// rCol := r.Limit(int64(metadatum.Length))
		// if metadatum.Spec.Deflate() {
		// 	rCol = rCol.AddProcessor(func(reader io.Reader) io.Reader {
		// 		return flate.NewReader(reader)
		// 	})
		// }

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
			it := column.ReadValueMetadataColumn(rCol)
			for vm, err := range it {
				if err != nil {
					panic(err)
				}
				res.Operations[i] = append(res.Operations[i], vm)
				prevValueMetadata = append(prevValueMetadata, vm)
			}
		case column.TypeValue:
			// skip(rCol, metadatum.Spec, metadatum.Length)

			// TODO: HACK just for early visualisation
			res.Operations[i] = acc(column.ReadValueColumn(rCol, prevValueMetadata))
		}
	}

	res.HeadIndexes, err = readHeadIndexes(r, len(res.Heads))
	if err != nil {
		return nil, fmt.Errorf("error reading head indexes: %w", err)
	}

	return &res, nil
}
