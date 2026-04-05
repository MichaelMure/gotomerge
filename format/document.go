package format

import (
	"errors"
	"fmt"
	"iter"
	"strings"

	"gotomerge/column"
	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

// DocumentChunk is the parsed form of an Automerge document snapshot chunk.
//
// A document chunk is a merged representation of an entire document's history.
// Rather than storing individual changes, it stores:
//
//   - Every operation from every change, in columnar form (OpColumns). Each
//     operation carries a successor list: the operations that later overwrote it.
//     An operation with no successors is the current live value; one with successors
//     has been overwritten or deleted.
//   - Per-change metadata in ChangesColumns: actor, sequence number, maxOp,
//     timestamp, message, and dependency indices. This is enough to reconstruct
//     the change graph and map operations back to their originating change.
//     The operations themselves are not duplicated here; the range of operations
//     belonging to change i is derived from MaxOp[i-1]+1 .. MaxOp[i].
//     Dependency references are integer indices (not hashes) into this same array,
//     which is stored in topological order: if change A depends on change B,
//     B always appears at a lower index than A.
//
// Heads contains the hashes of the "tip" changes — the changes that no other
// change in this document depends on. They identify the document's current version.
// HeadIndexes is a parallel array: HeadIndexes[i] is the index of Heads[i] in the
// change summary table, for fast lookup without scanning all changes.
//
// Actors is the document-wide actor table. All actor references in both the change
// and operation columns are stored as indices into this table.
//
// The SubReaders in OpColumns and ChangesColumns point directly into the paged
// reader's pages — no copy is made. Callers must not call Skip on the paged reader
// while this DocumentChunk (or any OpSet derived from it) is still in use.
type DocumentChunk struct {
	Actors      []types.ActorId
	Heads       []types.ChangeHash
	HeadIndexes []uint64

	ChangeMetadata column.Metadata
	ChangesColumns ChangeColumns

	OpMetadata column.Metadata
	OpColumns  OperationColumns
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

func readDocumentChunk(r *ioutil.SubReader) (*DocumentChunk, error) {
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
		colReader, err := r.SubReader(offset, metadatum.Length)
		if err != nil {
			return nil, fmt.Errorf("error reading change column: %w", err)
		}
		offset += metadatum.Length

		if metadatum.Spec.Deflate() {
			colReader, _, err = colReader.Deflate()
			if err != nil {
				return nil, fmt.Errorf("error deflating change column: %w", err)
			}
		}

		switch metadatum.Spec {
		case colDocChgActor:
			res.ChangesColumns.ActorId = colReader
		case colDocChgSeqNum:
			res.ChangesColumns.SeqNum = colReader
		case colDocChgMaxOp:
			res.ChangesColumns.MaxOp = colReader
		case colDocChgTime:
			res.ChangesColumns.Time = colReader
		case colDocChgMessage:
			res.ChangesColumns.Message = colReader
		case colDocChgDepsGrp:
			res.ChangesColumns.DependenciesGroup = colReader
		case colDocChgDepsIdx:
			res.ChangesColumns.DependenciesIndex = colReader
		case colValMeta:
			res.ChangesColumns.ExtraMetadata = colReader
		case colVal:
			res.ChangesColumns.ExtraData = colReader
		}
	}

	for _, metadatum := range res.OpMetadata {
		colReader, err := r.SubReader(offset, metadatum.Length)
		if err != nil {
			return nil, fmt.Errorf("error reading op column: %w", err)
		}
		offset += metadatum.Length

		if metadatum.Spec.Deflate() {
			colReader, _, err = colReader.Deflate()
			if err != nil {
				return nil, fmt.Errorf("error deflating op column: %w", err)
			}
		}

		switch metadatum.Spec {
		case colObjActor:
			res.OpColumns.ObjectActorId = colReader
		case colObjCtr:
			res.OpColumns.ObjectCounter = colReader
		case colKeyActor:
			res.OpColumns.KeyActorId = colReader
		case colKeyCtr:
			res.OpColumns.KeyCounter = colReader
		case colKeyStr:
			res.OpColumns.KeyString = colReader
		case colDocOpActor:
			res.OpColumns.ActorId = colReader
		case colDocOpCtr:
			res.OpColumns.Counter = colReader
		case colInsert:
			res.OpColumns.Insert = colReader
		case colAction:
			res.OpColumns.Action = colReader
		case colValMeta:
			res.OpColumns.ValueMetadata = colReader
		case colVal:
			res.OpColumns.Value = colReader
		case colPredGrp:
			res.OpColumns.PredecessorGroup = colReader
		case colPredActor:
			res.OpColumns.PredecessorActorId = colReader
		case colPredCtr:
			res.OpColumns.PredecessorCounter = colReader
		case colDocSuccGrp:
			res.OpColumns.SuccessorGroup = colReader
		case colDocSuccActor:
			res.OpColumns.SuccessorActorId = colReader
		case colDocSuccCtr:
			res.OpColumns.SuccessorCounter = colReader
		case colExpandControl:
			res.OpColumns.ExpandControl = colReader
		case colMark:
			res.OpColumns.Mark = colReader
		}
	}

	// SubReader() calls do not advance r, so r is still at the start of the
	// column data. Skip past all column bytes before reading HeadIndexes.
	if err := r.Skip(int(offset)); err != nil {
		return nil, fmt.Errorf("error skipping column data: %w", err)
	}

	res.HeadIndexes, err = readHeadIndexes(r, len(res.Heads))
	if err != nil {
		return nil, fmt.Errorf("error reading head indexes: %w", err)
	}

	return &res, nil
}

// Changes iterates over the change summaries embedded in this document snapshot.
//
// Each yielded DocChange has actor IDs resolved from the document's Actors table,
// but dependency references remain as integer indices (DocChange.Deps) rather than
// hashes. This is because the document chunk does not store per-change hashes;
// computing them requires re-serializing each change's content.
func (d DocumentChunk) Changes() iter.Seq2[types.DocChange, error] {
	actor, e1 := column.Req(d.ChangesColumns.ActorId, column.NewActorReader, "changes.actorId")
	seqNum, e2 := column.Req(d.ChangesColumns.SeqNum, column.NewDeltaReader, "changes.seqNum")
	maxOp, e3 := column.Req(d.ChangesColumns.MaxOp, column.NewDeltaReader, "changes.maxOp")
	time, e4 := column.Opt(d.ChangesColumns.Time, column.NewDeltaReader)
	message, e5 := column.Opt(d.ChangesColumns.Message, column.NewStringReader)
	depsGroup, e6 := column.Opt(d.ChangesColumns.DependenciesGroup, column.NewGroupReader)
	depsIndex, e7 := column.Opt(d.ChangesColumns.DependenciesIndex, column.NewDeltaReader)
	if err := errors.Join(e1, e2, e3, e4, e5, e6, e7); err != nil {
		return errSeq[types.DocChange](err)
	}
	cr := column.NewChangesReader(actor, seqNum, maxOp, time, message, depsGroup, depsIndex)

	return func(yield func(types.DocChange, error) bool) {
		for {
			raw, err := cr.Next()
			if isDone(err) {
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

// Operations iterates over every operation stored in this document snapshot.
//
// Unlike a change chunk — where operations are listed in creation order —
// a document chunk stores operations grouped by the object they belong to.
// This object-local ordering is part of the document chunk format definition.
func (d DocumentChunk) Operations() iter.Seq2[types.DocOperation, error] {
	objActor, e1 := column.Opt(d.OpColumns.ObjectActorId, column.NewActorReader)
	objCounter, e2 := column.Opt(d.OpColumns.ObjectCounter, column.NewUlebReader)
	opActor, e3 := column.Req(d.OpColumns.ActorId, column.NewActorReader, "op.actorId")
	opCounter, e4 := column.Req(d.OpColumns.Counter, column.NewDeltaReader, "op.counter")
	actionKind, e5 := column.Req(d.OpColumns.Action, column.NewUlebReader, "action")
	keyActor, e6 := column.Opt(d.OpColumns.KeyActorId, column.NewActorReader)
	keyCounter, e7 := column.Opt(d.OpColumns.KeyCounter, column.NewDeltaReader)
	keyString, e8 := column.Opt(d.OpColumns.KeyString, column.NewStringReader)
	insert, e9 := column.Opt(d.OpColumns.Insert, column.NewBoolReader)
	valueMeta, e10 := column.Opt(d.OpColumns.ValueMetadata, column.NewValueMetadataReader)
	value, e11 := column.Opt(d.OpColumns.Value, column.NewValueReader)
	succGroup, e12 := column.Opt(d.OpColumns.SuccessorGroup, column.NewGroupReader)
	succActor, e13 := column.Opt(d.OpColumns.SuccessorActorId, column.NewActorReader)
	succCounter, e14 := column.Opt(d.OpColumns.SuccessorCounter, column.NewDeltaReader)
	if err := errors.Join(e1, e2, e3, e4, e5, e6, e7, e8, e9, e10, e11, e12, e13, e14); err != nil {
		return errSeq[types.DocOperation](err)
	}

	objReader := column.NewObjectReader(objActor, objCounter)
	keyReader := column.NewKeyReader(keyActor, keyCounter, keyString)
	opIdReader := column.NewOpIdReader(opActor, opCounter)
	insertReader := column.NewInsertReader(insert)
	actionReader := column.NewActionReader(actionKind, valueMeta, value)
	succReader := column.NewGroupedOpIdReader("successor", succGroup, succActor, succCounter)

	return func(yield func(types.DocOperation, error) bool) {
		for {
			action, errAction := actionReader.Next()
			if errAction != nil && !isDone(errAction) {
				yield(types.DocOperation{}, errAction)
				return
			}

			// The action column drives iteration length: it is always present and
			// its exhaustion signals that all operations have been yielded. Other
			// columns may be absent (nil reader) or shorter if their values were
			// all null/default.
			if isDone(errAction) {
				return
			}

			obj, err := objReader.Next()
			if err != nil && !isDone(err) {
				yield(types.DocOperation{}, err)
				return
			}

			key, err := keyReader.Next()
			if err != nil && !isDone(err) {
				yield(types.DocOperation{}, err)
				return
			}

			id, err := opIdReader.Next()
			if err != nil && !isDone(err) {
				yield(types.DocOperation{}, err)
				return
			}

			insert, err := insertReader.Next()
			if err != nil && !isDone(err) {
				yield(types.DocOperation{}, err)
				return
			}

			succ, err := succReader.Next()
			if err != nil && !isDone(err) {
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
