/*
Package opset maintains the merged set of all operations that make up an
Automerge document and provides queries over it.

# Architecture overview

	DocumentChunk                      ChangeChunk(s)
	     │ ApplyDocument                    │ ApplyChange
	     ▼                                  ▼
	┌──────────────────────────────────────────────────────────────┐
	│  OpSet                                                       │
	│                                                              │
	│  actors []ActorId  ◄── global actor table (all sources)      │
	│  actorIdx map                                                │
	│  heads / appliedHashes  ◄── change-graph version tracking    │
	│                                                              │
	│  ┌───────────────────────────────┐  ┌──────────────────────┐ │
	│  │  snapshotStore    (snapshot/) │  │  deltaStore (delta/) │ │
	│  │                               │  │                      │ │
	│  │  succCount []uint32           │  │  ops   []Op          │ │
	│  │  insert    bitset.Bitset      │  │  byId  map[OpId]     │ │
	│  │                               │  │  byObj map[ObjectId] │ │
	│  │  objRanges  map[ObjectId]     │  └──────────────────────┘ │
	│  │  objCreators map[ObjectId]    │                           │
	│  │  byId       map[OpId]uint32   │                           │
	│  │                               │                           │
	│  │  SubReaders (key/opId/action) │  ← zero-copy col refs     │
	│  │       │                       │                           │
	│  │       ▼                       │                           │
	│  │  seekIndex  (seek.go)         │  ← checkpoints            │
	│  │  []seekPoint{                 │     every seekStride ops  │
	│  │    opIdx, *KeyReader,         │                           │
	│  │    *OpIdReader, *ActionReader │                           │
	│  │  }                            │                           │
	│  └───────────────────────────────┘                           │
	└─────────────────────────┬────────────────────────────────────┘
	                          │
	       ┌──────────────────┼──────────────────┐
	       ▼                  ▼                  ▼
	  MapGet/All          MapKeys            ObjType
	  (ListElements)      (Text)

# Data sources

An OpSet is populated from two kinds of chunks:

DocumentChunk — a merged snapshot of the entire document history. Its op
columns are stored as SubReader references directly into the paged reader's
pages (zero-copy). The snapshotStore wraps these column refs and builds a
small amount of derived metadata in a single pass:

  - succCount []uint32  - how many successors each op has (live iff 0)
  - insert    bitset    - set for list-insertion ops (skipped by map queries)
  - objRanges           - contiguous op range per object, for O(1) lookup
  - objCreators         - action kind of Make* ops, for ObjType() without scan
  - byId                - opId → array index, for predecessor updates

ChangeChunk — an incremental change from one actor. Each change is applied
on top of whatever is already in the OpSet (snapshot or prior changes). Ops
are stored as mutable structs in deltaStore rather than column refs for two
reasons: change chunks carry no successor column (a change cannot know what
will supersede it), so SuccCount must be writable in place; and byId provides
O(1) predecessor lookup when applying further changes. Applying a change:

 1. Validates that all declared dependencies are already present.
 2. Translates local actor indices to the OpSet's global actor table.
 3. For each predecessor listed in an op, increments SuccCount on that
    predecessor (in snapshotStore.succCount or in deltaStore.ops).
 4. Appends the new op to the changes slice.
 5. Updates heads: removes satisfied deps, adds this change's hash.

# Actor tables

Each chunk has its own local actor numbering. The OpSet maintains a
document-global actor table (actors []ActorId) and an actorIdx map for
deduplication. At apply time, local indices are translated to global ones so
all stored OpIds are self-consistent within the OpSet.

# Query model

A query such as MapGet(obj, "key") proceeds in two phases:

Phase 1 — locate the ops for the target object.

	For snapshot ops: objRanges[obj] gives an opRange{start, end} — the
	contiguous slice of column positions that belong to obj. This is O(1).
	The contiguous range is a guarantee of the DocumentChunk binary format:
	document chunks store operations grouped by object (all ops for object A,
	then all for object B, …), unlike change chunks which store ops in creation
	order across all objects. The single-pass in ApplyDocument records each
	object's span as ops arrive in that grouped order.

	For change ops: changesByObj[obj] gives a []uint32 of indices into the
	changes slice. Also O(1).

Phase 2 — scan and filter.

	For snapshot ops, scanSnapshotRange is called with the opRange. It cannot
	simply jump to position r.start in the column readers, because the readers
	are stateful streams — seeking means replaying from the beginning. The seek
	index (docSeekIndex) solves this:

	  1. seekIdx.seek(r.start) returns the checkpoint just before r.start and
	     the number of ops to skip from there (at most seekStride, currently 64).
	  2. Each checkpoint (docSeekPoint) holds pre-forked KeyReader, OpIdReader,
	     and ActionReader positioned at a known op boundary. Forking a reader
	     calls SubReaderOffset(0) on its underlying SubReader, which returns a
	     new independent read cursor into the same zero-copy page memory —
	     no data is copied.
	  3. The forked readers are advanced skip times to reach r.start, then
	     iterated from r.start to r.end. At each position, succCount[idx] and
	     insert.Get(idx) are read from the flat arrays (random access, O(1)),
	     while the key, opId, and action values are decoded from the column
	     readers in step.

	Each SubReader covers the full byte span of one column as decoded from
	the document chunk (e.g. all key-actor bytes for every op in the snapshot).
	Forking is cheap: it copies the reader state (a few integers) and creates
	a new independent cursor into the same bytes — no data is copied.

	For delta ops, each Op struct is read directly by index from deltaStore.ops
	— no column decoding needed.

	Both scans apply the same three filters:

	  - snapshotStore.succCount[i] > 0 / deltaStore.ops[i].SuccCount > 0  — superseded
	  - insert                             — list-insertion marker (map queries)
	  - ActionDelete                       — live tombstone, key was deleted

	When both a snapshotStore and a deltaStore are present, both are scanned and
	results are merged. For map queries the highest-OpId winner is kept (counter
	first, actor bytes as tiebreaker).

# Heads

heads tracks which applied changes no other applied change depends on. It
represents the current document version. After ApplyDocument it is
initialised from DocumentChunk.Heads. Each ApplyChange removes the new
change's declared deps from heads and adds its own hash.
*/
package opset
