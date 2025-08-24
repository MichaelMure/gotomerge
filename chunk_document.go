package gotomerge

import (
	"compress/flate"
	"errors"
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

	OpMetadata column.Metadata
	OpColumns  OperationColumns
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
	for i, metadatum := range d.OpMetadata {
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

	res.OpMetadata, err = column.ReadMetadata(r)
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

	for _, metadatum := range res.OpMetadata {
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
			res.OpColumns.ObjectActorId = column.ReadActorColumn(rCol)
		case 2: // ID: 0, type: uleb128
			res.OpColumns.ObjectCounter = column.ReadUlebColumn(rCol)
		case 17: // ID: 1, type: actor
			res.OpColumns.KeyActorId = column.ReadActorColumn(rCol)
		case 19: // ID: 1, type: delta
			res.OpColumns.KeyCounter = column.ReadDeltaColumn(rCol)
		case 21: // ID: 1, type: string
			res.OpColumns.KeyString = column.ReadStringColumn(rCol)
		case 33: // ID: 2, type: actor
			res.OpColumns.ActorId = column.ReadActorColumn(rCol)
		case 35: // ID: 2, type: delta
			res.OpColumns.Counter = column.ReadDeltaColumn(rCol)
		case 52: // ID: 3, type: bool
			res.OpColumns.Insert = column.ReadBooleanColumn(rCol)
		case 66: // ID: 4, type: uleb128
			res.OpColumns.Action = column.ReadUlebColumn(rCol)
		case 86: // ID: 5, type: value_metadata
			res.OpColumns.ValueMetadata = column.ReadValueMetadataColumn(rCol)
		case 87: // ID: 5, type: value
			res.OpColumns.Value = column.NewValueColumn(rCol)
		case 112: // ID: 7, type: group
			res.OpColumns.PredecessorGroup = column.ReadGroupColumn(rCol)
		case 113: // ID: 7, type: actor
			res.OpColumns.PredecessorActorId = column.ReadActorColumn(rCol)
		case 115: // ID: 7, type: delta
			res.OpColumns.PredecessorCounter = column.ReadDeltaColumn(rCol)
		case 128: // ID: 8, type: group
			res.OpColumns.SuccessorGroup = column.ReadGroupColumn(rCol)
		case 129: // ID: 8, type: actor
			res.OpColumns.SuccessorActorId = column.ReadActorColumn(rCol)
		case 131: // ID: 8, type: delta
			res.OpColumns.SuccessorCounter = column.ReadDeltaColumn(rCol)
		case 148: // ID: 9, type: bool
			res.OpColumns.ExpandControl = column.ReadBooleanColumn(rCol)
		case 165: // ID: 10, type: string
			res.OpColumns.Mark = column.ReadStringColumn(rCol)
		default:
			// TODO: unknown column should be maintained
			return nil, fmt.Errorf("unknown column type (TODO implementation): %v", metadatum.Spec)
		}
	}

	res.HeadIndexes, err = readHeadIndexes(r, len(res.Heads))
	if err != nil {
		return nil, fmt.Errorf("error reading head indexes: %w", err)
	}

	return &res, nil
}

func (d DocumentChunk) Operations() iter.Seq2[DocOperation, error] {
	objIter := column.ObjectColumn(d.OpColumns.ObjectActorId, d.OpColumns.ObjectCounter)
	keyIter := column.KeyColumn(d.OpColumns.KeyActorId, d.OpColumns.KeyCounter, d.OpColumns.KeyString)
	OpIdIter := column.OperationIdColumn(d.OpColumns.ActorId, d.OpColumns.Counter)
	insertIter := column.InsertColumn(d.OpColumns.Insert)
	actionIter := column.ActionColumn(d.OpColumns.Action, d.OpColumns.ValueMetadata, d.OpColumns.Value)
	succIter := column.GroupedOperationIdColumn("successor", d.OpColumns.SuccessorGroup, d.OpColumns.SuccessorActorId, d.OpColumns.SuccessorCounter)

	// TODO: text formatting

	return func(yield func(DocOperation, error) bool) {
		defer objIter.Stop()
		defer keyIter.Stop()
		defer OpIdIter.Stop()
		defer insertIter.Stop()
		defer actionIter.Stop()

		for {
			action, errAction := actionIter.Next()
			if errAction != nil && !errors.Is(errAction, column.ErrDone) {
				yield(DocOperation{}, errAction)
				return
			}

			// Action act as the marker for how long we should iterate
			// (from the rust codebase)
			if errors.Is(errAction, column.ErrDone) {
				return
			}

			obj, err := objIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(DocOperation{}, err)
				return
			}

			key, err := keyIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(DocOperation{}, err)
				return
			}

			id, err := OpIdIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(DocOperation{}, err)
				return
			}

			insert, err := insertIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(DocOperation{}, err)
				return
			}

			succ, err := succIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(DocOperation{}, err)
				return
			}

			if !yield(DocOperation{
				Object:     obj,
				Key:        key,
				Id:         id,
				Insert:     insert,
				Action:     action,
				Successors: succ,
			}, nil) {
				return
			}
		}
	}
}

type DocOperation struct {
	Object     types.ObjectId
	Key        types.Key
	Id         types.OpId
	Insert     bool
	Action     types.Action
	Successors []types.OpId
}

func (o DocOperation) String() string {
	var res strings.Builder
	res.WriteString("Operation {\n")
	res.WriteString(fmt.Sprintf("  \tObject: %v\n", o.Object))
	res.WriteString(fmt.Sprintf("  \tKey: %v\n", o.Key))
	res.WriteString(fmt.Sprintf("  \tId: %v\n", o.Id))
	res.WriteString(fmt.Sprintf("  \tInsert: %v\n", o.Insert))
	res.WriteString(fmt.Sprintf("  \tAction: %v\n", o.Action))
	res.WriteString(fmt.Sprintf("  \tSuccessors: %v\n", o.Successors))
	for i, succ := range o.Successors {
		res.WriteString(fmt.Sprintf("  \tSuccessors[%d]: %v\n", i, succ))
	}
	res.WriteString("  }")
	return res.String()
}
