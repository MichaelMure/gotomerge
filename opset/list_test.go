package opset

import (
	"testing"

	"github.com/MichaelMure/gotomerge/types"

	"github.com/stretchr/testify/require"
)

func listObj(t *testing.T, s *OpSet) types.ObjectId {
	t.Helper()
	op, ok := s.MapGet(types.RootObjectId(), "list")
	require.True(t, ok, "root[\"list\"] not found")
	return types.ObjectId(op.Id)
}

func textObj(t *testing.T, s *OpSet) types.ObjectId {
	t.Helper()
	op, ok := s.MapGet(types.RootObjectId(), "text")
	require.True(t, ok, "root[\"text\"] not found")
	return types.ObjectId(op.Id)
}

func elemStrings(elems []Op) []string {
	var out []string
	for _, op := range elems {
		if v, ok := op.Action.Value.(string); ok {
			out = append(out, v)
		}
	}
	return out
}

func TestListSequential(t *testing.T) {
	s := applyFile(t, "../testdata/list_sequential.automerge")
	obj := listObj(t, s)

	elems := s.ListElements(obj)
	require.Equal(t, []string{"a", "b", "c"}, elemStrings(elems))

	op, ok := s.ListGet(obj, 1)
	require.True(t, ok)
	require.Equal(t, "b", op.Action.Value)

	_, ok = s.ListGet(obj, 3)
	require.False(t, ok)
}

func TestListWithDelete(t *testing.T) {
	s := applyFile(t, "../testdata/list_with_delete.automerge")
	obj := listObj(t, s)

	require.Equal(t, []string{"a", "c"}, elemStrings(s.ListElements(obj)))
}

func TestListInsertAfterDeleted(t *testing.T) {
	// ["a", "b", "c"] → delete "b" → insert "x" after deleted slot → ["a", "x", "c"]
	s := applyFile(t, "../testdata/list_insert_after_deleted.automerge")
	obj := listObj(t, s)

	require.Equal(t, []string{"a", "x", "c"}, elemStrings(s.ListElements(obj)))
}

func TestListConcurrentInserts(t *testing.T) {
	// Actor A (0x00...) and actor B (0xff...) both insert at head concurrently.
	// B has higher actor bytes → "from_b" comes first.
	s := applyFile(t, "../testdata/list_concurrent_inserts.automerge")
	obj := listObj(t, s)

	elems := elemStrings(s.ListElements(obj))
	require.Equal(t, 2, len(elems))
	require.Equal(t, "from_b", elems[0])
	require.Equal(t, "from_a", elems[1])
}

func TestTextSequential(t *testing.T) {
	s := applyFile(t, "../testdata/text_sequential.automerge")
	obj := textObj(t, s)

	require.Equal(t, "Hello, World!", s.Text(obj))
}

// BenchmarkTextEdits measures Text() on a large editing history (~29 MB).
// The result is O(N) in the number of live characters; this benchmark exists
// to catch regressions before the planned incremental/rope optimisation.
func BenchmarkTextEdits(b *testing.B) {
	s := applyFile(b, "../testdata/text-edits.amrg")
	obj, ok := s.MapGet(types.RootObjectId(), "text")
	if !ok {
		b.Fatal("root[\"text\"] not found")
	}
	textObjId := types.ObjectId(obj.Id)

	b.ResetTimer()
	var n int
	for range b.N {
		n = len(s.Text(textObjId))
	}
	b.ReportMetric(float64(n), "chars")
}
