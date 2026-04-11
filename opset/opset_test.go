package opset

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/gotomerge/column/rle"
	"github.com/MichaelMure/gotomerge/format"
	"github.com/MichaelMure/gotomerge/types"
	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

// valueMetaString encodes a value metadata entry for a string of the given length.
// ValueMetadata = (length << 4) | type; string type = 6.
func valueMetaString(length int) uint64 {
	return uint64(length)<<4 | 6
}

// changeChunkSetKey builds a minimal ChangeChunk with one op: set root[key] = value (string).
func changeChunkSetKey(actor types.ActorId, seqNum, startOp uint64, deps []types.ChangeHash, hash types.ChangeHash, key, value string) *format.ChangeChunk {
	return &format.ChangeChunk{
		Hash:         hash,
		Dependencies: deps,
		Actor:        actor,
		SeqNum:       seqNum,
		StartOp:      startOp,
		OpColumns: format.OperationColumns{
			KeyString:     ioutil.NewSubReader(rle.EncodeString(key)),
			Action:        ioutil.NewSubReader(rle.EncodeUint64(uint64(types.ActionSet))),
			ValueMetadata: ioutil.NewSubReader(rle.EncodeUint64(valueMetaString(len(value)))),
			Value:         ioutil.NewSubReader([]byte(value)),
		},
	}
}

// changeChunkDelete builds a minimal ChangeChunk with one op: delete the op
// identified by predActorIdx/predCounter (in the change's local actor space).
func changeChunkDelete(actor types.ActorId, seqNum, startOp uint64, deps []types.ChangeHash, hash types.ChangeHash, key string, predActorIdx uint32, predCounter int64) *format.ChangeChunk {
	return &format.ChangeChunk{
		Hash:         hash,
		Dependencies: deps,
		Actor:        actor,
		SeqNum:       seqNum,
		StartOp:      startOp,
		OpColumns: format.OperationColumns{
			KeyString:          ioutil.NewSubReader(rle.EncodeString(key)),
			Action:             ioutil.NewSubReader(rle.EncodeUint64(uint64(types.ActionDelete))),
			PredecessorGroup:   ioutil.NewSubReader(rle.EncodeUint64(1)),
			PredecessorActorId: ioutil.NewSubReader(rle.EncodeUint64(uint64(predActorIdx))),
			PredecessorCounter: ioutil.NewSubReader(rle.EncodeInt64(predCounter)),
		},
	}
}

func applyFile(t testing.TB, path string) *OpSet {
	t.Helper()

	// Load the entire file into memory so that SubReaders can outlive Skip calls.
	// bytesReader.Skip only advances a position offset; the underlying []byte is
	// never freed. SubReaders created from this reader remain valid for the
	// lifetime of the OpSet.
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	s := New()
	r := ioutil.NewSubReader(data)
	for !r.Empty() {
		c, toSkip, err := format.ReadChunk(r)
		require.NoError(t, err)
		switch cc := c.(type) {
		case *format.DocumentChunk:
			require.NoError(t, s.ApplyDocument(cc))
		case *format.ChangeChunk:
			require.NoError(t, s.ApplyChange(cc))
		}
		require.NoError(t, r.Skip(toSkip))
	}
	return s
}

// TestApplyDocumentChunk loads the interop exemplar — a document with one change
// containing a variety of value types — and checks every key and value.
func TestApplyDocumentChunk(t *testing.T) {
	s := applyFile(t, "../testdata/exemplar")

	kind, ok := s.ObjType(types.RootObjectId())
	require.True(t, ok)
	require.Equal(t, types.ActionMakeMap, kind)

	checkScalar := func(key string, expected types.Action) {
		t.Helper()
		op, ok := s.MapGet(types.RootObjectId(), key)
		require.True(t, ok, "key %q not found", key)
		require.Equal(t, expected, op.Action)
	}

	checkScalar("bool", types.Action{Kind: types.ActionSet, Value: true})
	checkScalar("bytes", types.Action{Kind: types.ActionSet, Value: []byte{0x85, 0x6f, 0x4a, 0x83}})
	checkScalar("counter", types.Action{Kind: types.ActionSet, Value: types.Counter(5)})
	checkScalar("fp", types.Action{Kind: types.ActionSet, Value: float64(3.14159267)})
	checkScalar("int", types.Action{Kind: types.ActionSet, Value: int64(-4)})
	checkScalar("location", types.Action{Kind: types.ActionSet, Value: "https://automerge.org/"})
	checkScalar("timestamp", types.Action{Kind: types.ActionSet, Value: types.Timestamp(-905182979000)})
	checkScalar("title", types.Action{Kind: types.ActionSet, Value: "Hello \U0001f1ec\U0001f1e7\U0001f468\u200d\U0001f468\u200d\U0001f467\u200d\U0001f466\U0001f600"})
	checkScalar("uint", types.Action{Kind: types.ActionSet, Value: uint64(18446744073709551615)})

	// "notes" is a text object; the op creates a child object rather than storing a scalar.
	notesOp, ok := s.MapGet(types.RootObjectId(), "notes")
	require.True(t, ok)
	require.Equal(t, types.ActionMakeText, notesOp.Action.Kind)
	notesObj := types.ObjectId(notesOp.Id)
	notesKind, ok := s.ObjType(notesObj)
	require.True(t, ok)
	require.Equal(t, types.ActionMakeText, notesKind)

	// All keys must have exactly one live value (no conflicts in the exemplar).
	for _, k := range s.MapKeys(types.RootObjectId()) {
		require.Len(t, s.MapGetAll(types.RootObjectId(), k), 1, "unexpected conflict at key %q", k)
	}
}

// TestApplyChangeChunks applies two changes in order and checks the resulting state.
func TestApplyChangeChunks(t *testing.T) {
	s := applyFile(t, "../testdata/two_change_chunks.automerge")

	kind, ok := s.ObjType(types.RootObjectId())
	require.True(t, ok)
	require.Equal(t, types.ActionMakeMap, kind)

	// The file sets root["a"] to a nested map object.
	op, ok := s.MapGet(types.RootObjectId(), "a")
	require.True(t, ok)
	require.Equal(t, types.ActionMakeMap, op.Action.Kind)

	// The nested map object must also be queryable.
	nestedObj := types.ObjectId(op.Id)
	nestedKind, ok := s.ObjType(nestedObj)
	require.True(t, ok)
	require.Equal(t, types.ActionMakeMap, nestedKind)
}

// TestConflict applies two concurrent changes (no dependency between them) that
// both set the same key. MapGetAll must return both ops; MapGet must return the
// one with the higher OpId (higher actor bytes win when counters are equal).
func TestConflict(t *testing.T) {
	actorA := types.ActorId{0x00}
	actorB := types.ActorId{0xff}

	hashA := types.ChangeHash{0x01}
	hashB := types.ChangeHash{0x02}

	s := New()
	// Both changes have startOp=1, counter=1, no predecessors — a true conflict.
	require.NoError(t, s.ApplyChange(changeChunkSetKey(actorA, 1, 1, nil, hashA, "x", "hello")))
	require.NoError(t, s.ApplyChange(changeChunkSetKey(actorB, 1, 1, nil, hashB, "x", "world")))

	all := s.MapGetAll(types.RootObjectId(), "x")
	require.Len(t, all, 2, "expected conflict with two live ops")

	// Both hashes are heads (neither depends on the other).
	heads := s.Heads()
	require.Len(t, heads, 2)

	// MapGet returns the winner: actorB (0xff) > actorA (0x00).
	winner, ok := s.MapGet(types.RootObjectId(), "x")
	require.True(t, ok)
	require.Equal(t, s.Actor(winner.Id.ActorIdx), actorB)
}

// TestDelete applies a change that sets a key and then a change that deletes it.
// After the delete, MapGet must return nothing.
func TestDelete(t *testing.T) {
	actor := types.ActorId{0xab}
	hashSet := types.ChangeHash{0x01}
	hashDel := types.ChangeHash{0x02}

	s := New()
	// Change 1: set root["x"] = "hello", op counter = 1.
	require.NoError(t, s.ApplyChange(changeChunkSetKey(actor, 1, 1, nil, hashSet, "x", "hello")))

	// Verify it's visible before the delete.
	_, ok := s.MapGet(types.RootObjectId(), "x")
	require.True(t, ok)

	// Change 2: delete root["x"], predActorIdx=0 (own actor), predCounter=1.
	require.NoError(t, s.ApplyChange(changeChunkDelete(actor, 2, 2, []types.ChangeHash{hashSet}, hashDel, "x", 0, 1)))

	_, ok = s.MapGet(types.RootObjectId(), "x")
	require.False(t, ok, "key should be gone after delete")

	require.Empty(t, s.MapGetAll(types.RootObjectId(), "x"))

	// Only the delete is a head now.
	heads := s.Heads()
	require.Len(t, heads, 1)
	require.Equal(t, hashDel, heads[0])
}

// TestUnsatisfiedDependency verifies that ApplyChange rejects a change whose
// declared dependency has not been applied yet.
func TestUnsatisfiedDependency(t *testing.T) {
	actor := types.ActorId{0x01}
	hashA := types.ChangeHash{0xaa}
	hashB := types.ChangeHash{0xbb}

	s := New()
	// Change B declares a dependency on change A, but A hasn't been applied.
	cc := changeChunkSetKey(actor, 2, 2, []types.ChangeHash{hashA}, hashB, "x", "v")
	err := s.ApplyChange(cc)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsatisfied dependency")

	// After applying A first, B must succeed.
	ccA := changeChunkSetKey(actor, 1, 1, nil, hashA, "y", "v")
	require.NoError(t, s.ApplyChange(ccA))
	require.NoError(t, s.ApplyChange(cc))
}

// TestMapConflictFixture loads the map_conflict fixture generated by the Rust
// implementation: two actors concurrently set root["x"]. Verifies that both
// values are visible as a conflict and that MapGet returns the correct winner.
func TestMapConflictFixture(t *testing.T) {
	s := applyFile(t, "../testdata/map_conflict.automerge")

	all := s.MapGetAll(types.RootObjectId(), "x")
	require.Len(t, all, 2, "expected two conflicting values at root[\"x\"]")

	// Both values must be present.
	var vals []string
	for _, op := range all {
		v, ok := op.Action.Value.(string)
		require.True(t, ok, "expected string value")
		vals = append(vals, v)
	}
	require.ElementsMatch(t, []string{"from_a", "from_b"}, vals)

	// MapGet returns exactly one winner.
	winner, ok := s.MapGet(types.RootObjectId(), "x")
	require.True(t, ok)
	// Actor 0xff..ff has higher bytes than 0x00..00, so "from_b" wins.
	require.Equal(t, "from_b", winner.Action.Value)
}

// TestMapDeleteFixture loads the map_delete fixture: root["x"] is set then
// deleted. Verifies the key is absent after the delete.
func TestMapDeleteFixture(t *testing.T) {
	s := applyFile(t, "../testdata/map_delete.automerge")

	_, ok := s.MapGet(types.RootObjectId(), "x")
	require.False(t, ok, "deleted key should not be visible")

	require.Empty(t, s.MapGetAll(types.RootObjectId(), "x"))
	require.NotContains(t, s.MapKeys(types.RootObjectId()), "x")
}
