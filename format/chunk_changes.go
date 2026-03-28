package format

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"math"
	"strings"

	"github.com/jcalabro/leb128"

	"gotomerge/column"
	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

// ChangeChunk is the parsed form of an Automerge change chunk.
//
// A change records a set of operations authored by a single actor at a point in
// time. Changes are the unit of replication: peers exchange changes to converge
// on the same document state.
//
// The operations are stored in a columnar layout: instead of one record per
// operation, each field (object, key, action, …) is stored as its own compressed
// column. OperationColumns holds the decoded iterators for those columns.
//
// OtherActors lists every actor referenced *by the operations in this change*
// other than the change's own Actor. Operation actor-index 0 always refers to
// Actor; indices 1..N refer to OtherActors[i-1]. This local numbering keeps
// repeated actor IDs compact on the wire.
//
// ExtraBytes is a reserved extension field in the binary format. It is preserved
// intact so files using future extensions can be round-tripped without loss.
type ChangeChunk struct {
	// Hash is set by ReadChunk after the checksum is verified. It is the full
	// 32-byte SHA-256 of the serialized change and serves as the change's
	// globally unique identifier used in dependency references between changes.
	Hash         types.ChangeHash
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

	unknownColumns []rawColumn
}

func (ChangeChunk) chunk() {}

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
		var rawCol ioutil.SubReader
		rawCol, err = r.SubReader(offset, metadatum.Length)
		if err != nil {
			return nil, fmt.Errorf("error reading operation column: %w", err)
		}
		offset += metadatum.Length
		if metadatum.Spec.Deflate() {
			return nil, fmt.Errorf("deflate not supported in change chunk column")
		}

		switch metadatum.Spec {
		case 1: // ID: 0, type: actor
			res.OpColumns.ObjectActorId = column.ReadActorColumn(rawCol)
		case 2: // ID: 0, type: uleb128
			res.OpColumns.ObjectCounter = column.ReadUlebColumn(rawCol)
		case 17: // ID: 1, type: actor
			res.OpColumns.KeyActorId = column.ReadActorColumn(rawCol)
		case 19: // ID: 1, type: delta
			res.OpColumns.KeyCounter = column.ReadDeltaColumn(rawCol)
		case 21: // ID: 1, type: string
			res.OpColumns.KeyString = column.ReadStringColumn(rawCol)
		case 33: // ID: 2, type: actor
			res.OpColumns.ActorId = column.ReadActorColumn(rawCol)
		case 35: // ID: 2, type: delta
			res.OpColumns.Counter = column.ReadDeltaColumn(rawCol)
		case 52: // ID: 3, type: bool
			res.OpColumns.Insert = column.ReadBooleanColumn(rawCol)
		case 66: // ID: 4, type: uleb128
			res.OpColumns.Action = column.ReadUlebColumn(rawCol)
		case 86: // ID: 5, type: value_metadata
			res.OpColumns.ValueMetadata = column.ReadValueMetadataColumn(rawCol)
		case 87: // ID: 5, type: value
			res.OpColumns.Value = column.NewValueColumn(rawCol)
		case 112: // ID: 7, type: group
			res.OpColumns.PredecessorGroup = column.ReadGroupColumn(rawCol)
		case 113: // ID: 7, type: actor
			res.OpColumns.PredecessorActorId = column.ReadActorColumn(rawCol)
		case 115: // ID: 7, type: delta
			res.OpColumns.PredecessorCounter = column.ReadDeltaColumn(rawCol)
		case 128: // ID: 8, type: group
			res.OpColumns.SuccessorGroup = column.ReadGroupColumn(rawCol)
		case 129: // ID: 8, type: actor
			res.OpColumns.SuccessorActorId = column.ReadActorColumn(rawCol)
		case 131: // ID: 8, type: delta
			res.OpColumns.SuccessorCounter = column.ReadDeltaColumn(rawCol)
		case 148: // ID: 9, type: bool
			res.OpColumns.ExpandControl = column.ReadBooleanColumn(rawCol)
		case 165: // ID: 10, type: string
			res.OpColumns.Mark = column.ReadStringColumn(rawCol)
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

	return func(yield func(types.ChangeOperation, error) bool) {
		defer objIter.Stop()
		defer keyIter.Stop()
		defer insertIter.Stop()
		defer actionIter.Stop()
		defer predIter.Stop()

		var opIdx uint64
		for {
			action, errAction := actionIter.Next()
			if errAction != nil && !errors.Is(errAction, column.ErrDone) {
				yield(types.ChangeOperation{}, errAction)
				return
			}

			// The action column drives iteration length: it is always present and
			// its exhaustion signals that all operations have been yielded. Other
			// columns may be absent (nil iterator) or shorter if their values were
			// all null/default.
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

			counter := cc.StartOp + opIdx
			if counter > math.MaxUint32 {
				yield(types.ChangeOperation{}, fmt.Errorf("operation counter overflow at op %d", opIdx))
				return
			}

			if !yield(types.ChangeOperation{
				Id:           types.OpId{ActorIdx: 0, Counter: uint32(counter)},
				Object:       obj,
				Key:          key,
				Insert:       insert,
				Action:       action,
				Predecessors: pred,
			}, nil) {
				return
			}
			opIdx++
		}
	}
}

