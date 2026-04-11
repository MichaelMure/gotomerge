package opset

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/gotomerge/format"
	"github.com/MichaelMure/gotomerge/types"
	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

// commit is a test helper that calls Commit and fails on error.
func commit(t *testing.T, txn *Transaction) {
	t.Helper()
	require.NoError(t, txn.Commit(io.Discard))
}

func TestTxnMapSet(t *testing.T) {
	s := New()
	actor := types.NewActorId()

	txn := s.Begin(actor, 1)
	id := txn.MapSet(types.RootObjectId(), "name", "Alice")
	commit(t, txn)

	op, ok := s.MapGet(types.RootObjectId(), "name")
	require.True(t, ok)
	require.Equal(t, "Alice", op.Action.Value)
	require.Equal(t, id, op.Id)
}

func TestTxnMapSetOverwrite(t *testing.T) {
	s := New()
	actor := types.NewActorId()

	txn := s.Begin(actor, 1)
	txn.MapSet(types.RootObjectId(), "x", "first")
	commit(t, txn)

	txn2 := s.Begin(actor, 2)
	txn2.MapSet(types.RootObjectId(), "x", "second")
	commit(t, txn2)

	op, ok := s.MapGet(types.RootObjectId(), "x")
	require.True(t, ok)
	require.Equal(t, "second", op.Action.Value)
}

func TestTxnMapDelete(t *testing.T) {
	s := New()
	actor := types.NewActorId()

	txn := s.Begin(actor, 1)
	txn.MapSet(types.RootObjectId(), "x", "hello")
	commit(t, txn)

	txn2 := s.Begin(actor, 2)
	txn2.MapDelete(types.RootObjectId(), "x")
	commit(t, txn2)

	_, ok := s.MapGet(types.RootObjectId(), "x")
	require.False(t, ok)
}

func TestTxnMultipleOpsInOneTxn(t *testing.T) {
	s := New()
	actor := types.NewActorId()

	txn := s.Begin(actor, 1)
	txn.MapSet(types.RootObjectId(), "a", "1")
	txn.MapSet(types.RootObjectId(), "b", "2")
	txn.MapSet(types.RootObjectId(), "c", "3")
	commit(t, txn)

	for key, want := range map[string]string{"a": "1", "b": "2", "c": "3"} {
		op, ok := s.MapGet(types.RootObjectId(), key)
		require.True(t, ok, "key %q", key)
		require.Equal(t, want, op.Action.Value, "key %q", key)
	}
}

func TestTxnMakeNestedMap(t *testing.T) {
	s := New()
	actor := types.NewActorId()

	txn := s.Begin(actor, 1)
	nested, _ := txn.MakeMap(types.RootObjectId(), "config")
	txn.MapSet(nested, "debug", true)
	commit(t, txn)

	op, ok := s.MapGet(types.RootObjectId(), "config")
	require.True(t, ok)
	require.Equal(t, types.ActionMakeMap, op.Action.Kind)

	op2, ok := s.MapGet(nested, "debug")
	require.True(t, ok)
	require.Equal(t, true, op2.Action.Value)
}

func TestTxnListInsert(t *testing.T) {
	s := New()
	actor := types.NewActorId()

	txn := s.Begin(actor, 1)
	listObj, _ := txn.MakeList(types.RootObjectId(), "items")
	head := types.KeyOpId{} // Counter=0 = head sentinel
	id1 := txn.ListInsert(listObj, head, "a")
	id2 := txn.ListInsert(listObj, types.KeyOpId(id1), "b")
	txn.ListInsert(listObj, types.KeyOpId(id2), "c")
	commit(t, txn)

	elems := s.ListElements(listObj)
	require.Equal(t, []string{"a", "b", "c"}, elemStrings(elems))
}

func TestTxnListDelete(t *testing.T) {
	s := New()
	actor := types.NewActorId()

	txn := s.Begin(actor, 1)
	listObj, _ := txn.MakeList(types.RootObjectId(), "items")
	head := types.KeyOpId{}
	id1 := txn.ListInsert(listObj, head, "a")
	id2 := txn.ListInsert(listObj, types.KeyOpId(id1), "b")
	txn.ListInsert(listObj, types.KeyOpId(id2), "c")
	commit(t, txn)

	txn2 := s.Begin(actor, 2)
	txn2.ListDelete(listObj, id2, id2) // delete "b"
	commit(t, txn2)

	elems := s.ListElements(listObj)
	require.Equal(t, []string{"a", "c"}, elemStrings(elems))
}

func TestTxnText(t *testing.T) {
	s := New()
	actor := types.NewActorId()

	txn := s.Begin(actor, 1)
	textObj, _ := txn.MakeText(types.RootObjectId(), "doc")
	head := types.KeyOpId{}
	h := txn.ListInsert(textObj, head, "H")
	e := txn.ListInsert(textObj, types.KeyOpId(h), "e")
	l1 := txn.ListInsert(textObj, types.KeyOpId(e), "l")
	l2 := txn.ListInsert(textObj, types.KeyOpId(l1), "l")
	txn.ListInsert(textObj, types.KeyOpId(l2), "o")
	commit(t, txn)

	require.Equal(t, "Hello", s.Text(textObj))
}

// TestTxnRoundTrip confirms the full wire round-trip: commit → ReadChunk → ApplyChange on a fresh OpSet.
func TestTxnRoundTrip(t *testing.T) {
	s1 := New()
	actor := types.NewActorId()

	txn := s1.Begin(actor, 1)
	txn.MapSet(types.RootObjectId(), "greeting", "hello")

	var buf bytes.Buffer
	require.NoError(t, txn.Commit(&buf))

	// Parse the chunk back from the wire bytes.
	r := ioutil.NewSubReader(buf.Bytes())
	parsed, toSkip, err := format.ReadChunk(r)
	require.NoError(t, err)
	require.NoError(t, r.Skip(toSkip))

	parsedCC, ok := parsed.(*format.ChangeChunk)
	require.True(t, ok)

	// Apply to a fresh peer.
	s2 := New()
	require.NoError(t, s2.ApplyChange(parsedCC))

	op, ok := s2.MapGet(types.RootObjectId(), "greeting")
	require.True(t, ok)
	require.Equal(t, "hello", op.Action.Value)
}
