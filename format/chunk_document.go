package format

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

	ChangeMetadata column.Metadata
	ChangesColumns ChangeColumns

	OpMetadata column.Metadata
	OpColumns  OperationColumns

	unknownColumns []rawColumn
}

func (DocumentChunk) chunk() {}

func (d DocumentChunk) String() string {
	var res strings.Builder
	res.WriteString("DocumentChunk {\n")
	res.WriteString(fmt.Sprintf("  Actors: %v\n", d.Actors))
	res.WriteString(fmt.Sprintf("  Heads: %v\n", d.Heads))
	res.WriteString(fmt.Sprintf("  HeadIndexes: %v\n", d.HeadIndexes))
	for i, metadatum := range d.ChangeMetadata {
		res.WriteString(fmt.Sprintf("  ChangeMetadata[%d]: %v\n", i, metadatum))
	}
	for i, metadatum := range d.OpMetadata {
		res.WriteString(fmt.Sprintf("  OperationMetadata[%d]: %v\n", i, metadatum))
	}
	i := 0
	for change, err := range d.Changes() {
		if err != nil {
			res.WriteString(fmt.Sprintf("  Change[%d]: %v\n", i, err))
		} else {
			res.WriteString(fmt.Sprintf("  Change[%d]: %v\n", i, change))
		}
		i++
	}
	i = 0
	for operation, err := range d.Operations() {
		if err != nil {
			res.WriteString(fmt.Sprintf("  Operation[%d]: %v\n", i, err))
		} else {
			res.WriteString(fmt.Sprintf("  Operation[%d]: %v\n", i, operation))
		}
		i++
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
	for _, metadatum := range res.ChangeMetadata {
		var rawCol ioutil.SubReader
		rawCol, err = r.SubReader(offset, metadatum.Length)
		if err != nil {
			return nil, fmt.Errorf("error reading change column: %w", err)
		}
		offset += metadatum.Length
		var rCol io.Reader = rawCol
		if metadatum.Spec.Deflate() {
			rCol = flate.NewReader(rawCol)
		}

		switch metadatum.Spec {
		case 1:
			res.ChangesColumns.ActorId = column.ReadActorColumn(rCol)
		case 3:
			res.ChangesColumns.SeqNum = column.ReadDeltaColumn(rCol)
		case 19:
			res.ChangesColumns.MaxOp = column.ReadDeltaColumn(rCol)
		case 35:
			res.ChangesColumns.Time = column.ReadDeltaColumn(rCol)
		case 53:
			res.ChangesColumns.Message = column.ReadStringColumn(rCol)
		case 64:
			res.ChangesColumns.DependenciesGroup = column.ReadGroupColumn(rCol)
		case 67:
			res.ChangesColumns.DependenciesIndex = column.ReadDeltaColumn(rCol)
		case 86:
			res.ChangesColumns.ExtraMetadata = column.ReadValueMetadataColumn(rCol)
		case 87:
			res.ChangesColumns.ExtraData = column.NewValueColumn(rCol)
		default:
			data, err := io.ReadAll(rawCol)
			if err != nil {
				return nil, fmt.Errorf("error reading unknown change column: %w", err)
			}
			res.unknownColumns = append(res.unknownColumns, rawColumn{
				specBits: uint32(metadatum.Spec),
				data:     data,
			})
		}
	}

	for _, metadatum := range res.OpMetadata {
		var rawCol ioutil.SubReader
		rawCol, err = r.SubReader(offset, metadatum.Length)
		if err != nil {
			return nil, fmt.Errorf("error reading op column: %w", err)
		}
		offset += metadatum.Length
		var rCol io.Reader = rawCol
		if metadatum.Spec.Deflate() {
			rCol = flate.NewReader(rawCol)
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
			data, err := io.ReadAll(rawCol)
			if err != nil {
				return nil, fmt.Errorf("error reading unknown op column: %w", err)
			}
			res.unknownColumns = append(res.unknownColumns, rawColumn{
				specBits: uint32(metadatum.Spec),
				data:     data,
			})
		}
	}

	res.HeadIndexes, err = readHeadIndexes(r, len(res.Heads))
	if err != nil {
		return nil, fmt.Errorf("error reading head indexes: %w", err)
	}

	return &res, nil
}

func (d DocumentChunk) Changes() iter.Seq2[types.DocChange, error] {
	changesIter := column.NewChangesIter(
		d.ChangesColumns.ActorId,
		d.ChangesColumns.SeqNum,
		d.ChangesColumns.MaxOp,
		d.ChangesColumns.Time,
		d.ChangesColumns.Message,
		d.ChangesColumns.DependenciesGroup,
		d.ChangesColumns.DependenciesIndex,
	)

	return func(yield func(types.DocChange, error) bool) {
		defer changesIter.Stop()

		for {
			raw, err := changesIter.Next()
			if errors.Is(err, column.ErrDone) {
				return
			}
			if err != nil {
				yield(types.DocChange{}, err)
				return
			}

			if raw.ActorIdx >= uint64(len(d.Actors)) {
				yield(types.DocChange{}, fmt.Errorf("actor index out of range: %d (have %d actors)", raw.ActorIdx, len(d.Actors)))
				return
			}
			actor := d.Actors[raw.ActorIdx]

			var t types.Timestamp
			if raw.Time != nil {
				t = types.Timestamp(*raw.Time)
			}

			if !yield(types.DocChange{
				ActorId: actor,
				SeqNum:  raw.SeqNum,
				MaxOp:   raw.MaxOp,
				Deps:    raw.Deps,
				Time:    t,
				Message: raw.Message,
			}, nil) {
				return
			}
		}
	}
}

func (d DocumentChunk) Operations() iter.Seq2[types.DocOperation, error] {
	objIter := column.ObjectColumn(d.OpColumns.ObjectActorId, d.OpColumns.ObjectCounter)
	keyIter := column.KeyColumn(d.OpColumns.KeyActorId, d.OpColumns.KeyCounter, d.OpColumns.KeyString)
	OpIdIter := column.OperationIdColumn(d.OpColumns.ActorId, d.OpColumns.Counter)
	insertIter := column.InsertColumn(d.OpColumns.Insert)
	actionIter := column.ActionColumn(d.OpColumns.Action, d.OpColumns.ValueMetadata, d.OpColumns.Value)
	succIter := column.GroupedOperationIdColumn("successor", d.OpColumns.SuccessorGroup, d.OpColumns.SuccessorActorId, d.OpColumns.SuccessorCounter)

	return func(yield func(types.DocOperation, error) bool) {
		defer objIter.Stop()
		defer keyIter.Stop()
		defer OpIdIter.Stop()
		defer insertIter.Stop()
		defer actionIter.Stop()
		defer succIter.Stop()

		for {
			action, errAction := actionIter.Next()
			if errAction != nil && !errors.Is(errAction, column.ErrDone) {
				yield(types.DocOperation{}, errAction)
				return
			}

			// Actions act as the marker for how long we should iterate
			// (from the rust codebase)
			if errors.Is(errAction, column.ErrDone) {
				return
			}

			obj, err := objIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.DocOperation{}, err)
				return
			}

			key, err := keyIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.DocOperation{}, err)
				return
			}

			id, err := OpIdIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.DocOperation{}, err)
				return
			}

			insert, err := insertIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.DocOperation{}, err)
				return
			}

			succ, err := succIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.DocOperation{}, err)
				return
			}

			if !yield(types.DocOperation{
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

type ChangeColumns struct {
	ActorId column.ActorColumnIter
	SeqNum  column.DeltaColumnIter

	MaxOp column.DeltaColumnIter

	Time column.DeltaColumnIter

	Message column.StringColumnIter

	DependenciesGroup column.GroupColumnIter
	DependenciesIndex column.DeltaColumnIter

	ExtraMetadata column.ValueMetadataColumnIter
	ExtraData     column.ValueColumn
}
