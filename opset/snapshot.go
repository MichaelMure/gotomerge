package opset

import (
	"fmt"
	"io"

	"gotomerge/format"
	"gotomerge/types"
	"gotomerge/utils/bitset"
	ioutil "gotomerge/utils/io"
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
//   - SubReaders for key/action/opId columns — zero-copy references into the
//     paged reader's pages. The paged reader must remain alive while these are
//     used; see DocumentChunk for the lifetime contract.
type snapshotStore struct {
	opCount   uint32
	succCount []uint32
	insert    bitset.Bitset

	objRanges   map[types.ObjectId]opRange
	objCreators map[types.ObjectId]types.ActionKind
	byId        map[types.OpId]uint32

	// Column SubReaders for query-time scanning. Re-readable via SubReaderOffset(0).
	keyActorId    ioutil.SubReader
	keyCounter    ioutil.SubReader
	keyString     ioutil.SubReader
	opActorId     ioutil.SubReader
	opCounter     ioutil.SubReader
	action        ioutil.SubReader
	valueMetadata ioutil.SubReader
	value         ioutil.SubReader

	// seek is a sparse checkpoint list built from all query columns, spaced
	// seekStride ops apart. scanSnapshotRange uses it to jump near the target
	// range without replaying from column position 0.
	seek seekIndex
}

// ApplyDocument loads all operations from a document snapshot.
// The OpSet must be empty; applying a DocumentChunk on top of existing ops
// is not supported.
//
// The SubReaders stored in the resulting snapshotStore point into the paged
// reader's pages. The caller must not call Skip on the paged reader (past the
// document chunk's bytes) while this OpSet is in use.
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

	// Single pass to build all metadata structures.
	var opIdx uint32
	for docOp, err := range doc.Operations() {
		if err != nil {
			return fmt.Errorf("reading operation: %w", err)
		}

		ss.succCount = append(ss.succCount, uint32(len(docOp.Successors)))
		if docOp.Insert {
			ss.insert.Set(opIdx)
		}
		ss.byId[docOp.Id] = opIdx

		// Extend the contiguous range for this object.
		r := ss.objRanges[docOp.Object]
		if r.end == opIdx {
			r.end = opIdx + 1
		} else if r.start == r.end {
			r = opRange{start: opIdx, end: opIdx + 1}
		} else {
			r.end = opIdx + 1
		}
		ss.objRanges[docOp.Object] = r

		// Cache action kind of Make* ops for ObjType() lookups.
		switch docOp.Action.Kind {
		case types.ActionMakeMap, types.ActionMakeList, types.ActionMakeText:
			ss.objCreators[types.ObjectId(docOp.Id)] = docOp.Action.Kind
		}

		opIdx++
	}
	ss.opCount = opIdx

	// Store column SubReaders for query-time iteration (zero-copy page refs).
	ss.keyActorId = doc.OpColumns.KeyActorId
	ss.keyCounter = doc.OpColumns.KeyCounter
	ss.keyString = doc.OpColumns.KeyString
	ss.opActorId = doc.OpColumns.ActorId
	ss.opCounter = doc.OpColumns.Counter
	ss.action = doc.OpColumns.Action
	ss.valueMetadata = doc.OpColumns.ValueMetadata
	ss.value = doc.OpColumns.Value

	ss.seek = buildSeekIndex(ss)

	s.snapshot = ss

	for _, h := range doc.Heads {
		s.heads[h] = struct{}{}
	}

	return nil
}

// scanSnapshotRange iterates over operations [r.start, r.end) in the snapshot,
// calling fn for each with its position index and decoded Op. fn returns true
// to continue, false to stop early. Decode errors abort the scan silently.
//
// The seek index lets us jump to within seekStride ops of r.start in O(1) and
// then advance at most seekStride ops to reach the target.
func scanSnapshotRange(ss *snapshotStore, r opRange, fn func(idx uint32, op Op) bool) {
	if r.start >= r.end {
		return
	}

	pt, skip := ss.seek.seek(r.start)

	keyIter, err := pt.key.Fork()
	if err != nil {
		return
	}
	opIdIter, err := pt.opId.Fork()
	if err != nil {
		return
	}
	actionIter, err := pt.action.Fork()
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
