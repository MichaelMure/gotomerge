package gotomerge

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"strings"

	"github.com/jcalabro/leb128"

	"gotomerge/column"
	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

type ChangeChunk struct {
	Dependencies []types.ChangeHash
	Actor        types.ActorId
	SeqNum       uint64
	StartOp      uint64
	Time         types.Timestamp
	Message      string
	OtherActors  []types.ActorId

	OpMetadata column.Metadata
	OpColumns  OperationColumns

	ExtraBytes []byte
}

func (cc ChangeChunk) String() string {
	var res strings.Builder
	res.WriteString("ChangeChunk {\n")
	res.WriteString(fmt.Sprintf("  Dependencies: %v\n", cc.Dependencies))
	res.WriteString(fmt.Sprintf("  Actor: %v\n", cc.Actor))
	res.WriteString(fmt.Sprintf("  SeqNum: %v\n", cc.SeqNum))
	res.WriteString(fmt.Sprintf("  StartOp: %v\n", cc.StartOp))
	res.WriteString(fmt.Sprintf("  Time: %v\n", cc.Time))
	res.WriteString(fmt.Sprintf("  Message: %v\n", cc.Message))
	res.WriteString(fmt.Sprintf("  OtherActors: %v\n", cc.OtherActors))
	for i, metadatum := range cc.OpMetadata {
		res.WriteString(fmt.Sprintf("  OperationMetadata[%d]: %v\n", i, metadatum))
	}
	i := 0
	for operation, err := range cc.Operations() {
		if err != nil {
			res.WriteString(fmt.Sprintf("  Operation[%d]: %v\n", i, err))
		} else {
			res.WriteString(fmt.Sprintf("  Operation[%d]: %v\n", i, operation))
		}
		i++
	}
	res.WriteString(fmt.Sprintf("  ExtraBytes: %v\n", cc.ExtraBytes))
	res.WriteString("}\n")
	return res.String()
}

func readChangeChunk(r ioutil.SubReader) (*ChangeChunk, error) {
	var res ChangeChunk
	var err error

	res.Dependencies, err = readChangeHashes(r)
	if err != nil {
		return nil, fmt.Errorf("error reading dependencies: %w", err)
	}

	res.Actor, err = types.ReadLengthEncodedActorId(r)
	if err != nil {
		return nil, fmt.Errorf("error reading actor id: %w", err)
	}

	res.SeqNum, err = leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading seq num: %w", err)
	}

	res.StartOp, err = leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading start op: %w", err)
	}

	time, err := leb128.DecodeS64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading time: %w", err)
	}
	res.Time = types.Timestamp(time)

	msgLen, err := leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading message length: %w", err)
	}
	res.Message, err = ioutil.ReadStringLimitedPrealloc(r, msgLen)
	if err != nil {
		return nil, fmt.Errorf("error reading message: %w", err)
	}

	res.OtherActors, err = readActorIds(r)
	if err != nil {
		return nil, fmt.Errorf("error reading other actors: %w", err)
	}

	res.OpMetadata, err = column.ReadMetadata(r)
	if err != nil {
		return nil, fmt.Errorf("error reading operation metadata: %w", err)
	}

	var offset uint64
	for _, metadatum := range res.OpMetadata {
		rCol, err := r.SubReader(offset, metadatum.Length)
		if err != nil {
			return nil, fmt.Errorf("error reading operation column: %w", err)
		}
		offset += metadatum.Length
		if metadatum.Spec.Deflate() {
			return nil, fmt.Errorf("deflate not supported in change chunk column")
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

	extra, err := r.SubReaderOffset(offset)
	if err != nil {
		return nil, fmt.Errorf("error reading extra: %w", err)
	}

	if extra.Empty() {
		// Don't try to read the extra bytes if we know there is none.
		// This avoids an allocation in io.ReadAll(), as we know that virtually
		// all the changes we read don't have those extra bytes.
		return &res, nil
	}

	res.ExtraBytes, err = io.ReadAll(extra)
	if err != nil {
		return nil, fmt.Errorf("error reading extra bytes: %w", err)
	}

	return &res, nil
}

func (cc ChangeChunk) Operations() iter.Seq2[types.ChangeOperation, error] {
	objIter := column.ObjectColumn(cc.OpColumns.ObjectActorId, cc.OpColumns.ObjectCounter)
	keyIter := column.KeyColumn(cc.OpColumns.KeyActorId, cc.OpColumns.KeyCounter, cc.OpColumns.KeyString)
	insertIter := column.InsertColumn(cc.OpColumns.Insert)
	actionIter := column.ActionColumn(cc.OpColumns.Action, cc.OpColumns.ValueMetadata, cc.OpColumns.Value)
	predIter := column.GroupedOperationIdColumn("predecessors", cc.OpColumns.PredecessorGroup, cc.OpColumns.PredecessorActorId, cc.OpColumns.PredecessorCounter)

	// TODO: text formatting

	return func(yield func(types.ChangeOperation, error) bool) {
		defer objIter.Stop()
		defer keyIter.Stop()
		defer insertIter.Stop()
		defer actionIter.Stop()

		for {
			action, errAction := actionIter.Next()
			if errAction != nil && !errors.Is(errAction, column.ErrDone) {
				yield(types.ChangeOperation{}, errAction)
				return
			}

			// Action act as the marker for how long we should iterate
			// (from the rust codebase)
			if errors.Is(errAction, column.ErrDone) {
				return
			}

			obj, err := objIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.ChangeOperation{}, err)
				return
			}

			key, err := keyIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.ChangeOperation{}, err)
				return
			}

			insert, err := insertIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.ChangeOperation{}, err)
				return
			}

			pred, err := predIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.ChangeOperation{}, err)
				return
			}

			if !yield(types.ChangeOperation{
				Object:       obj,
				Key:          key,
				Insert:       insert,
				Action:       action,
				Predecessors: pred,
			}, nil) {
				return
			}
		}
	}
}

// type Change struct {
// 	ActorId   types.ActorId
// 	SeqNum    uint64
// 	Ops       []Operation
// 	Deps      []types.ChangeHash
// 	Time      types.Timestamp
// 	Message   string
// 	ExtraData any
// }

func (cc ChangeChunk) ToChange() {

}
