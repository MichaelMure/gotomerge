package opset

import (
	"errors"
	"fmt"
	"io"

	"github.com/MichaelMure/gotomerge/column"
	"github.com/MichaelMure/gotomerge/format"
	"github.com/MichaelMure/gotomerge/types"
	"github.com/MichaelMure/gotomerge/utils/bitset"
)

// opRange marks a contiguous span of column positions belonging to one object.
// Contiguity is guaranteed by the DocumentChunk binary format, which groups all
// ops for the same object together (unlike change chunks, which use creation
// order). Each object therefore occupies exactly one range.
type opRange struct{ start, end uint32 }

// snapshotStore holds the decoded state of a single DocumentChunk in a
// memory-efficient form. "Snapshot" refers to the DocumentChunk itself — a
// point-in-time materialization of the full op history — not the full document
// state, which is represented by the OpSet as a whole.
//
// Rather than keeping every operation as a struct, it stores:
//   - succCount as a flat array and insert as a bitset, indexed by op position.
//     These are the hot-path filters during queries: succCount[i] == 0 means
//     op i is the current live value; insert[i] marks list insertions (skipped
//     by map queries). The bitset packs 64 flags per word, 8x less than []bool.
//   - objRanges for O(1) lookup of which ops belong to each object.
//   - objCreators so that ObjType() can answer without scanning the action column.
//   - byId for predecessor-update lookups when a subsequent ChangeChunk
//     supersedes a snapshot op.
//   - seek: a sparse checkpoint list of forked column readers, built in the
//     same pass as the metadata above. Each checkpoint holds pre-forked
//     KeyReader, OpIdReader, and ActionReader positioned at a known op boundary.
//     The underlying column bytes are shared and never copied; the checkpoint
//     is just a cursor state (a few integers).
type snapshotStore struct {
	opCount   uint32
	succCount []uint32
	insert    bitset.Bitset

	objRanges   map[types.ObjectId]opRange
	objCreators map[types.ObjectId]types.ActionKind
	byId        map[types.OpId]uint32

	seek seekIndex

	// docChunk is the original DocumentChunk loaded via ApplyDocument. It is
	// kept so ExportDocument can re-iterate ops with their full column data
	// (including the original within-snapshot successor lists) when building a
	// new document chunk. The SubReaders inside docChunk share the same
	// underlying bytes; do not call Skip() on the original reader while the
	// OpSet is in use.
	docChunk *format.DocumentChunk
}

// isDone reports whether err signals that a column reader is exhausted.
// RLE readers return io.EOF; GroupedOpIdReader returns column.ErrDone when
// the group column is absent.
func isDone(err error) bool {
	return err == io.EOF || errors.Is(err, column.ErrDone)
}

// ApplyDocument loads all operations from a document snapshot.
// The OpSet must be empty; applying a DocumentChunk on top of existing ops
// is not supported.
//
// Column readers are set up once and iterated in a single pass to build both
// the metadata structures and the seek index. Seek checkpoints are forks of
// the live readers, so no second column scan is needed.
//
// The seek checkpoints hold column readers that reference the same underlying
// bytes as the DocumentChunk's SubReaders. The caller must not call Skip on the
// paged reader (past the document chunk's bytes) while this OpSet is in use.
func (s *OpSet) ApplyDocument(doc *format.DocumentChunk) error {
	if s.snapshot != nil || s.delta != nil {
		return fmt.Errorf("ApplyDocument called on non-empty OpSet")
	}

	// Register actors in document order. Because the OpSet is empty, each doc
	// actor index equals its OpSet index — no translation needed for column data.
	for _, a := range doc.Actors {
		s.internActor(a)
	}

	ss := &snapshotStore{
		objRanges:   make(map[types.ObjectId]opRange),
		objCreators: make(map[types.ObjectId]types.ActionKind),
		byId:        make(map[types.OpId]uint32),
	}

	// Set up all column readers in one shot. These are the only readers for
	// this document chunk — no parallel set is needed.
	objActor, e1 := column.Opt(doc.OpColumns.ObjectActorId, column.NewActorReader)
	objCounter, e2 := column.Opt(doc.OpColumns.ObjectCounter, column.NewUlebReader)
	keyActor, e3 := column.Opt(doc.OpColumns.KeyActorId, column.NewActorReader)
	keyCounter, e4 := column.Opt(doc.OpColumns.KeyCounter, column.NewDeltaReader)
	keyStr, e5 := column.Opt(doc.OpColumns.KeyString, column.NewStringReader)
	opActor, e6 := column.Req(doc.OpColumns.ActorId, column.NewActorReader, "op.actorId")
	opCounter, e7 := column.Req(doc.OpColumns.Counter, column.NewDeltaReader, "op.counter")
	insertCol, e8 := column.Opt(doc.OpColumns.Insert, column.NewBoolReader)
	actionKind, e9 := column.Req(doc.OpColumns.Action, column.NewUlebReader, "action")
	valueMeta, e10 := column.Opt(doc.OpColumns.ValueMetadata, column.NewValueMetadataReader)
	valueCol, e11 := column.Opt(doc.OpColumns.Value, column.NewValueReader)
	succGroup, e12 := column.Opt(doc.OpColumns.SuccessorGroup, column.NewGroupReader)
	succActor, e13 := column.Opt(doc.OpColumns.SuccessorActorId, column.NewActorReader)
	succCounter, e14 := column.Opt(doc.OpColumns.SuccessorCounter, column.NewDeltaReader)
	if err := errors.Join(e1, e2, e3, e4, e5, e6, e7, e8, e9, e10, e11, e12, e13, e14); err != nil {
		return fmt.Errorf("column setup: %w", err)
	}

	objReader := column.NewObjectReader(objActor, objCounter)
	keyReader := column.NewKeyReader(keyActor, keyCounter, keyStr)
	opIdReader := column.NewOpIdReader(opActor, opCounter)
	insertReader := column.NewInsertReader(insertCol)
	actionReader := column.NewActionReader(actionKind, valueMeta, valueCol)
	succReader := column.NewGroupedOpIdReader("successor", succGroup, succActor, succCounter)

	// Checkpoint at op 0 before reading anything.
	k0, ek := keyReader.Fork()
	oi0, eoi := opIdReader.Fork()
	a0, ea := actionReader.Fork()
	if err := errors.Join(ek, eoi, ea); err != nil {
		return fmt.Errorf("seek index initial checkpoint: %w", err)
	}
	var seek seekIndex

	// incDeltas maps Inc op IDs to their delta values. Used in the post-scan
	// pass to adjust succCount for counter ops whose successors are all Inc ops.
	incDeltas := make(map[types.OpId]int64)
	// counterSuccessors maps counter Set op IDs to their snapshot successor IDs.
	// Populated only for counter ops that have at least one successor.
	type counterEntry struct {
		opId types.OpId
		succ []types.OpId
	}
	var counterSuccessors []counterEntry

	var opIdx uint32
	for {
		// Checkpoint at stride boundaries before reading the next op.
		// The readers are positioned exactly before op opIdx here.
		if opIdx > 0 && opIdx%seekStride == 0 {
			kf, err1 := keyReader.Fork()
			oif, err2 := opIdReader.Fork()
			af, err3 := actionReader.Fork()
			if err := errors.Join(err1, err2, err3); err == nil {
				seek = append(seek, seekPoint{opIdx: opIdx, key: kf, opId: oif, action: af})
			}
		}

		// The action column drives iteration — its exhaustion signals end of ops.
		action, err := actionReader.Next()
		if isDone(err) {
			break
		}
		if err != nil {
			return fmt.Errorf("reading action at op %d: %w", opIdx, err)
		}

		obj, err := objReader.Next()
		if err != nil && !isDone(err) {
			return fmt.Errorf("reading object at op %d: %w", opIdx, err)
		}

		// Key is advanced to keep the reader in lockstep for seek checkpoints;
		// the decoded value is not needed for metadata.
		if _, err = keyReader.Next(); err != nil && !isDone(err) {
			return fmt.Errorf("reading key at op %d: %w", opIdx, err)
		}

		id, err := opIdReader.Next()
		if err != nil && !isDone(err) {
			return fmt.Errorf("reading opId at op %d: %w", opIdx, err)
		}

		ins, err := insertReader.Next()
		if err != nil && !isDone(err) {
			return fmt.Errorf("reading insert at op %d: %w", opIdx, err)
		}

		succ, err := succReader.Next()
		if err != nil && !isDone(err) {
			return fmt.Errorf("reading successors at op %d: %w", opIdx, err)
		}

		ss.succCount = append(ss.succCount, uint32(len(succ)))
		if ins {
			ss.insert.Set(opIdx)
		}
		ss.byId[id] = opIdx
		if id.Counter > s.maxOpCounter[id.ActorIdx] {
			s.maxOpCounter[id.ActorIdx] = id.Counter
		}

		r := ss.objRanges[obj]
		if r.end == opIdx {
			r.end = opIdx + 1
		} else if r.start == r.end {
			r = opRange{start: opIdx, end: opIdx + 1}
		} else {
			r.end = opIdx + 1
		}
		ss.objRanges[obj] = r

		switch action.Kind {
		case types.ActionMakeMap, types.ActionMakeList, types.ActionMakeText:
			ss.objCreators[types.ObjectId(id)] = action.Kind
		case types.ActionInc:
			incDeltas[id] = incDelta(action.Value)
		case types.ActionSet:
			if _, isCounter := action.Value.(types.Counter); isCounter && len(succ) > 0 {
				succCopy := make([]types.OpId, len(succ))
				copy(succCopy, succ)
				counterSuccessors = append(counterSuccessors, counterEntry{opId: id, succ: succCopy})
			}
		}

		opIdx++
	}
	ss.opCount = opIdx

	if opIdx > 0 {
		ss.seek = append(seekIndex{seekPoint{opIdx: 0, key: k0, opId: oi0, action: a0}}, seek...)
	}

	// Post-scan: adjust succCount for counter ops whose successors are Inc ops.
	// Inc ops do not kill the counter — they accumulate a delta instead.
	for _, entry := range counterSuccessors {
		opIdx := ss.byId[entry.opId]
		for _, succId := range entry.succ {
			if delta, ok := incDeltas[succId]; ok {
				ss.succCount[opIdx]--
				s.counterDeltas[entry.opId] += delta
			}
		}
	}

	s.snapshot = ss
	ss.docChunk = doc

	for _, h := range doc.Heads {
		s.heads[h] = struct{}{}
	}

	return nil
}

// scanRange iterates over operations [r.start, r.end) in the snapshot,
// calling fn for each with its position index and decoded Op. fn returns true
// to continue, false to stop early. Decode errors abort the scan silently.
//
// The seek index lets us jump to within seekStride ops of r.start in O(1) and
// then advance at most seekStride ops to reach the target.
func (ss *snapshotStore) scanRange(r opRange, fn func(idx uint32, op Op) bool) {
	if r.start >= r.end {
		return
	}

	point, skip := ss.seek.seek(r.start)

	keyIter, err := point.key.Fork()
	if err != nil {
		return
	}
	opIdIter, err := point.opId.Fork()
	if err != nil {
		return
	}
	actionIter, err := point.action.Fork()
	if err != nil {
		return
	}

	// Advance from the checkpoint to r.start (at most seekStride ops).
	for i := uint32(0); i < skip; i++ {
		if _, err := keyIter.Next(); err != nil && err != io.EOF {
			return
		}
		if _, err := opIdIter.Next(); err != nil && err != io.EOF {
			return
		}
		if _, err := actionIter.Next(); err != nil && err != io.EOF {
			return
		}
	}

	for idx := r.start; idx < r.end; idx++ {
		key, err := keyIter.Next()
		if err != nil {
			return
		}
		id, err := opIdIter.Next()
		if err != nil {
			return
		}
		action, err := actionIter.Next()
		if err != nil {
			return
		}

		op := Op{
			Id:     id,
			Key:    key,
			Insert: ss.insert.Get(idx),
			Action: action,
		}
		if !fn(idx, op) {
			return
		}
	}
}
