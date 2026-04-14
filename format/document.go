package format

import (
	"fmt"
	"iter"
	"strings"

	"github.com/MichaelMure/gotomerge/column"
	"github.com/MichaelMure/gotomerge/types"
	ioutil "github.com/MichaelMure/gotomerge/utils/io"
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

		switch metadatum.Spec.WithoutDeflate() {
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

		switch metadatum.Spec.WithoutDeflate() {
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
// Changes iterates over every change's raw metadata as stored in the column data.
// ActorIdx is an index into d.Actors; Time and Message are nil when absent.
// Dependency references are indices into the document's change array (not hashes).
func (d DocumentChunk) Changes() iter.Seq2[column.RawChangeMeta, error] {
	if d.ChangesColumns.ActorId == nil {
		return func(yield func(column.RawChangeMeta, error) bool) {}
	}
	cr := column.NewChangesReader(
		peek(d.ChangesColumns.ActorId, column.PeekActorReader),
		peek(d.ChangesColumns.SeqNum, column.PeekDeltaReader),
		peek(d.ChangesColumns.MaxOp, column.PeekDeltaReader),
		peek(d.ChangesColumns.Time, column.PeekDeltaReader),
		peek(d.ChangesColumns.Message, column.PeekStringReader),
		peek(d.ChangesColumns.DependenciesGroup, column.PeekGroupReader),
		peek(d.ChangesColumns.DependenciesIndex, column.PeekDeltaReader),
	)

	return func(yield func(column.RawChangeMeta, error) bool) {
		for {
			raw, err := cr.Next()
			if isDone(err) {
				return
			}
			if !yield(raw, err) || err != nil {
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
	// Zero-ops document: action column absent → yield nothing.
	if d.OpColumns.Action == nil {
		return func(yield func(types.DocOperation, error) bool) {}
	}
	objActor := peek(d.OpColumns.ObjectActorId, column.PeekActorReader)
	objCounter := peek(d.OpColumns.ObjectCounter, column.PeekUlebReader)
	opActor := peek(d.OpColumns.ActorId, column.PeekActorReader)
	opCounter := peek(d.OpColumns.Counter, column.PeekDeltaReader)
	actionKind := peek(d.OpColumns.Action, column.PeekUlebReader)
	keyActor := peek(d.OpColumns.KeyActorId, column.PeekActorReader)
	keyCounter := peek(d.OpColumns.KeyCounter, column.PeekDeltaReader)
	keyString := peek(d.OpColumns.KeyString, column.PeekStringReader)
	insert := peek(d.OpColumns.Insert, column.PeekBoolReader)
	valueMeta := peek(d.OpColumns.ValueMetadata, column.PeekValueMetadataReader)
	value := peek(d.OpColumns.Value, column.PeekValueReader)
	succGroup := peek(d.OpColumns.SuccessorGroup, column.PeekGroupReader)
	succActor := peek(d.OpColumns.SuccessorActorId, column.PeekActorReader)
	succCounter := peek(d.OpColumns.SuccessorCounter, column.PeekDeltaReader)

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
