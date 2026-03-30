// Package opset maintains the set of all operations that make up a document
// and provides queries over them.
//
// An OpSet is built by applying one or more chunks read from the format layer.
// Two sources are supported:
//
//   - A DocumentChunk, which already contains the fully merged operation set
//     with successor lists. Applying it reads column references directly into a
//     docStore — no data is copied for non-compressed columns. The paged reader's
//     pages act as the backing store and must remain alive (no Skip) while the
//     docStore is in use.
//
//   - ChangeChunks, applied one at a time in dependency order. Each change
//     carries predecessor references (the ops it supersedes). Applying a change
//     increments SuccCnt on predecessors and appends the new ops as structs.
//     Changes are rare additions on top of a document snapshot, so the
//     materialized-struct cost is bounded.
//
// After loading, the OpSet answers document-level queries: current values,
// conflict detection, key enumeration, and object type lookup.
package opset

import (
	"encoding/hex"
	"fmt"

	"gotomerge/format"
	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

// opRange marks a contiguous span of operations in the docStore that all
// belong to the same object. In a DocumentChunk, operations are grouped by
// object, so this is always a single range.
type opRange struct{ start, end uint32 }

// docStore holds the document snapshot in a memory-efficient form.
//
// Rather than keeping every operation as a struct, it stores:
//   - succCnt and insert as flat arrays, indexed by operation position. These
//     are the hot-path filters during queries: succCnt[i] == 0 means op i is
//     the current live value; insert[i] marks list insertions (skipped by map
//     queries).
//   - objRanges for O(1) lookup of which operations belong to each object.
//     Because document ops are grouped by object, each object occupies exactly
//     one contiguous range.
//   - objCreators so that ObjType() can answer without scanning the action column.
//   - byId for predecessor-update lookups when a subsequent ChangeChunk
//     supersedes a document op.
//   - SubReaders for key/action/opId columns — re-readable references into the
//     paged reader's pages (no copy). The paged reader must remain alive while
//     these are used; see DocumentChunk for the lifetime contract.
type docStore struct {
	opCount     uint32
	succCnt     []uint32
	insert      []bool
	objRanges   map[types.ObjectId]opRange
	objCreators map[types.ObjectId]types.ActionKind
	byId        map[types.OpId]uint32

	// Column SubReaders for query-time scanning. These point into the paged
	// reader's pages and are re-readable without copying (SubReaderOffset(0) forks).
	keyActorId    ioutil.SubReader
	keyCounter    ioutil.SubReader
	keyString     ioutil.SubReader
	opActorId     ioutil.SubReader
	opCounter     ioutil.SubReader
	action        ioutil.SubReader
	valueMetadata ioutil.SubReader
	value         ioutil.SubReader

	// seekIdx is a sparse combined checkpoint list built from all query columns.
	// scanDocRange uses it to seek past the ops before the target range in O(1)
	// instead of O(r.start).
	seekIdx docSeekIndex
}

// OpSet holds every operation that has been applied to a document.
//
// All OpIds stored here use indices into the OpSet's own actors slice.
// Indices from incoming chunks (DocumentChunk.Actors, ChangeChunk.Actor /
// OtherActors) are translated at apply time.
type OpSet struct {
	// actors is the document-wide actor table. OpId.ActorIdx values throughout
	// this OpSet are indices into this slice.
	actors   []types.ActorId
	actorIdx map[string]uint32 // hex(actorId) → index in actors

	// doc holds operations from the most recently applied DocumentChunk, stored
	// as column references plus sparse metadata. Nil if no DocumentChunk has been
	// applied yet.
	doc *docStore

	// changes holds operations from ChangeChunks applied after the document
	// snapshot. These are stored as structs because they are incremental
	// additions and their count is bounded relative to the snapshot.
	changes      []Op
	changesById  map[types.OpId]uint32       // opId → index in changes
	changesByObj map[types.ObjectId][]uint32 // objectId → indices in changes
}

func New() *OpSet {
	return &OpSet{
		actorIdx:     make(map[string]uint32),
		changesById:  make(map[types.OpId]uint32),
		changesByObj: make(map[types.ObjectId][]uint32),
	}
}

// Actor returns the ActorId for the given index.
func (s *OpSet) Actor(idx uint32) types.ActorId {
	return s.actors[idx]
}

// internActor registers an actor and returns its index in the OpSet actor table.
// If the actor is already registered, the existing index is returned.
func (s *OpSet) internActor(id types.ActorId) uint32 {
	key := hex.EncodeToString(id)
	if idx, ok := s.actorIdx[key]; ok {
		return idx
	}
	idx := uint32(len(s.actors))
	s.actors = append(s.actors, id)
	s.actorIdx[key] = idx
	return idx
}

// ApplyDocument loads all operations from a document snapshot.
// The OpSet must be empty; applying a DocumentChunk on top of existing ops
// is not supported.
//
// The SubReaders stored in the resulting docStore point into the paged reader's
// pages. The caller must not call Skip on the paged reader (past the document
// chunk's bytes) while this OpSet is in use.
func (s *OpSet) ApplyDocument(doc *format.DocumentChunk) error {
	if s.doc != nil || len(s.changes) > 0 {
		return fmt.Errorf("ApplyDocument called on non-empty OpSet")
	}

	// Register actors in the order they appear in the document. Because the
	// OpSet is empty, each doc actor index equals its OpSet index — no
	// translation is needed for the column data itself.
	for _, a := range doc.Actors {
		s.internActor(a)
	}

	ds := &docStore{
		objRanges:   make(map[types.ObjectId]opRange),
		objCreators: make(map[types.ObjectId]types.ActionKind),
		byId:        make(map[types.OpId]uint32),
	}

	// Single pass over all operations to build the metadata structures.
	// succCnt, insert, objRanges, objCreators, and byId are all populated here.
	var opIdx uint32
	for docOp, err := range doc.Operations() {
		if err != nil {
			return fmt.Errorf("reading operation: %w", err)
		}

		ds.succCnt = append(ds.succCnt, uint32(len(docOp.Successors)))
		ds.insert = append(ds.insert, docOp.Insert)
		ds.byId[docOp.Id] = opIdx

		// Track the contiguous range of ops that belong to each object.
		r := ds.objRanges[docOp.Object]
		if r.end == opIdx {
			r.end = opIdx + 1
		} else if r.start == r.end {
			r = opRange{start: opIdx, end: opIdx + 1}
		} else {
			r.end = opIdx + 1
		}
		ds.objRanges[docOp.Object] = r

		// Record the action kind of Make* ops so ObjType() works without
		// scanning the action column.
		switch docOp.Action.Kind {
		case types.ActionMakeMap, types.ActionMakeList, types.ActionMakeText:
			objId := types.ObjectId(docOp.Id)
			ds.objCreators[objId] = docOp.Action.Kind
		}

		opIdx++
	}
	ds.opCount = opIdx

	// Store SubReaders for query-time iteration. These are re-readable
	// via SubReaderOffset(0) without copying any bytes.
	ds.keyActorId = doc.OpColumns.KeyActorId
	ds.keyCounter = doc.OpColumns.KeyCounter
	ds.keyString = doc.OpColumns.KeyString
	ds.opActorId = doc.OpColumns.ActorId
	ds.opCounter = doc.OpColumns.Counter
	ds.action = doc.OpColumns.Action
	ds.valueMetadata = doc.OpColumns.ValueMetadata
	ds.value = doc.OpColumns.Value

	ds.seekIdx = buildDocSeekIndex(ds)

	s.doc = ds
	return nil
}

// ApplyChange applies a single change to the OpSet. Changes must be applied in
// dependency order: all changes listed in cc.Dependencies must already be present.
func (s *OpSet) ApplyChange(cc *format.ChangeChunk) error {
	// Build a mapping from the change's local actor indices to ours.
	// Index 0 is the change's own actor; indices 1..N are OtherActors[i-1].
	changeActors := make([]uint32, 1+len(cc.OtherActors))
	changeActors[0] = s.internActor(cc.Actor)
	for i, a := range cc.OtherActors {
		changeActors[i+1] = s.internActor(a)
	}

	translateId := func(id types.OpId) types.OpId {
		return types.OpId{ActorIdx: changeActors[id.ActorIdx], Counter: id.Counter}
	}
	translateObj := func(obj types.ObjectId) types.ObjectId {
		if obj.IsRoot() {
			return obj
		}
		return types.ObjectId(translateId(types.OpId(obj)))
	}

	for changeOp, err := range cc.Operations() {
		if err != nil {
			return fmt.Errorf("reading operation: %w", err)
		}

		// Increment SuccCnt on each predecessor this op supersedes.
		for _, pred := range changeOp.Predecessors {
			resolvedPred := translateId(pred)
			if s.doc != nil {
				if predIdx, ok := s.doc.byId[resolvedPred]; ok {
					s.doc.succCnt[predIdx]++
				}
			}
			if predIdx, ok := s.changesById[resolvedPred]; ok {
				s.changes[predIdx].SuccCnt++
			}
		}

		op := Op{
			Id:     translateId(changeOp.Id),
			Object: translateObj(changeOp.Object),
			Key:    translateKey(changeOp.Key, changeActors),
			Insert: changeOp.Insert,
			Action: changeOp.Action,
		}

		idx := uint32(len(s.changes))
		s.changes = append(s.changes, op)
		s.changesById[op.Id] = idx
		s.changesByObj[op.Object] = append(s.changesByObj[op.Object], idx)
	}

	return nil
}

// translateKey remaps any actor indices embedded in a Key using the provided
// actor index mapping.
func translateKey(key types.Key, actorMap []uint32) types.Key {
	if k, ok := key.(types.KeyOpId); ok {
		return types.KeyOpId{ActorIdx: actorMap[k.ActorIdx], Counter: k.Counter}
	}
	return key
}
