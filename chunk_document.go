package gotomerge

import (
	"compress/flate"
	"fmt"
	"io"
	"iter"
	"strings"

	"gotomerge/column"
	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

type DocumentChunk struct {
	Actors      []types.ActorId
	Heads       []types.ChangeHash
	HeadIndexes []uint64

	// TODO: should that stays?
	ChangeMetadata column.Metadata
	Changes        [][]any

	OperationMetadata column.Metadata
	OperationColumns  OperationColumns
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
	}
	for operation, err := range d.Operations() {
		if err != nil {
			res.WriteString(fmt.Sprintf("  Operation[i]: %v\n", err))
		} else {
			res.WriteString(fmt.Sprintf("  Operation[i]: %v\n", operation))
		}
	}
	res.WriteString("}\n")
	return res.String()
}

func readDocumentChunk(r ioutil.SubReader) (*DocumentChunk, error) {
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

	var offset uint64
	res.Changes = make([][]any, len(res.ChangeMetadata))
	for _, metadatum := range res.ChangeMetadata {
		var rCol io.Reader
		rCol, err = r.SubReader(offset, metadatum.Length)
		if err != nil {
			return nil, fmt.Errorf("error reading change column: %w", err)
		}
		offset += metadatum.Length
		if metadatum.Spec.Deflate() {
			rCol = flate.NewReader(rCol)
		}

		// TODO

		// switch metadatum.Spec.Type() {
		// case column.TypeGroup:
		// 	res.Changes[i] = acc(rle.ReadUint64RLE(rCol))
		// case column.TypeActor:
		// 	res.Changes[i] = acc(rle.ReadUint64RLE(rCol))
		// case column.TypeULEB128:
		// 	res.Changes[i] = acc(column.ReadUlebColumn(rCol))
		// case column.TypeDelta:
		// 	res.Changes[i] = acc(column.ReadDeltaColumn(rCol))
		// case column.TypeBool:
		// 	res.Changes[i] = acc(column.ReadBooleanColumn(rCol))
		// case column.TypeString:
		// 	res.Changes[i] = acc(column.ReadStringColumn(rCol))
		// case column.TypeValueMetadata:
		// 	res.Changes[i] = acc(column.ReadValueMetadataColumn(rCol))
		// case column.TypeValue:
		// 	skip(rCol, metadatum.Spec, metadatum.Length)
		// }
	}

	for _, metadatum := range res.OperationMetadata {
		var rCol io.Reader
		rCol, err = r.SubReader(offset, metadatum.Length)
		if err != nil {
			return nil, fmt.Errorf("error reading change column: %w", err)
		}
		offset += metadatum.Length
		if metadatum.Spec.Deflate() {
			rCol = flate.NewReader(rCol)
		}

		switch metadatum.Spec {
		case 1: // ID: 0, type: actor
			res.OperationColumns.ObjectActorId = column.ReadActorColumn(rCol)
		case 2: // ID: 0, type: uleb128
			res.OperationColumns.ObjectCounter = column.ReadUlebColumn(rCol)
		case 17: // ID: 1, type: actor
			res.OperationColumns.KeyActorId = column.ReadActorColumn(rCol)
		case 19: // ID: 1, type: delta
			res.OperationColumns.KeyCounter = column.ReadDeltaColumn(rCol)
		case 21: // ID: 1, type: string
			res.OperationColumns.KeyString = column.ReadStringColumn(rCol)
		case 33: // ID: 2, type: actor
			res.OperationColumns.ActorId = column.ReadActorColumn(rCol)
		case 35: // ID: 2, type: delta
			res.OperationColumns.Counter = column.ReadDeltaColumn(rCol)
		case 52: // ID: 3, type: bool
			res.OperationColumns.Insert = column.ReadBooleanColumn(rCol)
		case 66: // ID: 4, type: uleb128
			res.OperationColumns.Action = column.ReadUlebColumn(rCol)
		case 86: // ID: 5, type: value_metadata
			res.OperationColumns.ValueMetadata = column.ReadValueMetadataColumn(rCol)
		case 87: // ID: 5, type: value
			res.OperationColumns.Value = column.NewValueColumn(rCol)
		case 112: // ID: 7, type: group
			res.OperationColumns.PredecessorGroup = column.ReadGroupColumn(rCol)
		case 113: // ID: 7, type: actor
			res.OperationColumns.PredecessorActorId = column.ReadActorColumn(rCol)
		case 115: // ID: 7, type: delta
			res.OperationColumns.PredecessorCounter = column.ReadDeltaColumn(rCol)
		case 128: // ID: 8, type: group
			res.OperationColumns.SuccessorGroup = column.ReadGroupColumn(rCol)
		case 129: // ID: 8, type: actor
			res.OperationColumns.SuccessorActorId = column.ReadActorColumn(rCol)
		case 131: // ID: 8, type: delta
			res.OperationColumns.SuccessorCounter = column.ReadDeltaColumn(rCol)
		case 148: // ID: 9, type: bool
			res.OperationColumns.ExpandControl = column.ReadBooleanColumn(rCol)
		case 165: // ID: 10, type: string
			res.OperationColumns.Mark = column.ReadStringColumn(rCol)
		default:
			// TODO: unknown column should be maintained
			panic(fmt.Sprintf("unknown column type: %v", metadatum.Spec))
		}
	}

	res.HeadIndexes, err = readHeadIndexes(r, len(res.Heads))
	if err != nil {
		return nil, fmt.Errorf("error reading head indexes: %w", err)
	}

	return &res, nil
}

func (d DocumentChunk) Operations() iter.Seq2[Operation, error] {
	return d.OperationColumns.Operations()
}
