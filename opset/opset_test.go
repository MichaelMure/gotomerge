package opset

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"gotomerge/format"
	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

func applyFile(t *testing.T, path string) *OpSet {
	t.Helper()

	// Load the entire file into memory so that SubReaders can outlive Skip calls.
	// bytesReader.Skip only advances a position offset; the underlying []byte is
	// never freed. SubReaders created from this reader remain valid for the
	// lifetime of the OpSet.
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	s := New()
	r := ioutil.NewBytesReader(data)
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
