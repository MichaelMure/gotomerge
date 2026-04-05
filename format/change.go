package format

import (
	"errors"
	"fmt"
	"io"
	"iter"
	"math"
	"strings"

	"github.com/MichaelMure/leb128"

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
// column. OperationColumns holds re-readable SubReaders for those columns; fresh
// readers are created each time Operations() is called.
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

func readChangeChunk(r *ioutil.SubReader) (*ChangeChunk, error) {
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
		rawCol, err := r.SubReader(offset, metadatum.Length)
		if err != nil {
			return nil, fmt.Errorf("error reading operation column: %w", err)
		}
		offset += metadatum.Length
		if metadatum.Spec.Deflate() {
			return nil, fmt.Errorf("deflate not supported in change chunk column")
		}

		switch metadatum.Spec {
		case colObjActor:
			res.OpColumns.ObjectActorId = rawCol
		case colObjCtr:
			res.OpColumns.ObjectCounter = rawCol
		case colKeyActor:
			res.OpColumns.KeyActorId = rawCol
		case colKeyCtr:
			res.OpColumns.KeyCounter = rawCol
		case colKeyStr:
			res.OpColumns.KeyString = rawCol
		case colDocOpActor:
			res.OpColumns.ActorId = rawCol
		case colDocOpCtr:
			res.OpColumns.Counter = rawCol
		case colInsert:
			res.OpColumns.Insert = rawCol
		case colAction:
			res.OpColumns.Action = rawCol
		case colValMeta:
			res.OpColumns.ValueMetadata = rawCol
		case colVal:
			res.OpColumns.Value = rawCol
		case colPredGrp:
			res.OpColumns.PredecessorGroup = rawCol
		case colPredActor:
			res.OpColumns.PredecessorActorId = rawCol
		case colPredCtr:
			res.OpColumns.PredecessorCounter = rawCol
		case colDocSuccGrp:
			res.OpColumns.SuccessorGroup = rawCol
		case colDocSuccActor:
			res.OpColumns.SuccessorActorId = rawCol
		case colDocSuccCtr:
			res.OpColumns.SuccessorCounter = rawCol
		case colExpandControl:
			res.OpColumns.ExpandControl = rawCol
		case colMark:
			res.OpColumns.Mark = rawCol
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
	objActor, e1 := column.Opt(cc.OpColumns.ObjectActorId, column.NewActorReader)
	objCounter, e2 := column.Opt(cc.OpColumns.ObjectCounter, column.NewUlebReader)
	keyActor, e3 := column.Opt(cc.OpColumns.KeyActorId, column.NewActorReader)
	keyCounter, e4 := column.Opt(cc.OpColumns.KeyCounter, column.NewDeltaReader)
	keyString, e5 := column.Opt(cc.OpColumns.KeyString, column.NewStringReader)
	insert, e6 := column.Opt(cc.OpColumns.Insert, column.NewBoolReader)
	actionKind, e7 := column.Opt(cc.OpColumns.Action, column.NewUlebReader)
	valueMeta, e8 := column.Opt(cc.OpColumns.ValueMetadata, column.NewValueMetadataReader)
	value, e9 := column.Opt(cc.OpColumns.Value, column.NewValueReader)
	predGroup, e10 := column.Opt(cc.OpColumns.PredecessorGroup, column.NewGroupReader)
	predActor, e11 := column.Opt(cc.OpColumns.PredecessorActorId, column.NewActorReader)
	predCounter, e12 := column.Opt(cc.OpColumns.PredecessorCounter, column.NewDeltaReader)
	if err := errors.Join(e1, e2, e3, e4, e5, e6, e7, e8, e9, e10, e11, e12); err != nil {
		return errSeq[types.ChangeOperation](err)
	}

	objReader := column.NewObjectReader(objActor, objCounter)
	keyReader := column.NewKeyReader(keyActor, keyCounter, keyString)
	insertReader := column.NewInsertReader(insert)
	actionReader := column.NewActionReader(actionKind, valueMeta, value)
	predReader := column.NewGroupedOpIdReader("predecessors", predGroup, predActor, predCounter)

	return func(yield func(types.ChangeOperation, error) bool) {
		var opIdx uint64
		for {
			action, errAction := actionReader.Next()
			if errors.Is(errAction, column.ErrDone) {
				return
			}
			if errAction != nil {
				yield(types.ChangeOperation{}, errAction)
				return
			}

			// The action column drives iteration length: it is always present and
			// its exhaustion signals that all operations have been yielded. Other
			// columns may be absent (nil reader) or shorter if their values were
			// all null/default.

			obj, err := objReader.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.ChangeOperation{}, err)
				return
			}

			key, err := keyReader.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.ChangeOperation{}, err)
				return
			}

			insert, err := insertReader.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(types.ChangeOperation{}, err)
				return
			}

			pred, err := predReader.Next()
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
