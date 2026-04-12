package docproxy

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/gotomerge/types"
)

func TestLoadDocument(t *testing.T) {
	f, err := os.Open("../testdata/exemplar")
	require.NoError(t, err)
	defer f.Close()

	doc, err := LoadDocument(f)
	require.NoError(t, err)

	keys := doc.Keys()
	require.NotEmpty(t, keys)

	jsonStr, err := json.MarshalIndent(doc, "", "  ")
	require.NoError(t, err)
	fmt.Println(string(jsonStr))
}

func TestLoadDocumentTwoChangeChunks(t *testing.T) {
	f, err := os.Open("../testdata/two_change_chunks.automerge")
	require.NoError(t, err)
	defer f.Close()

	doc, err := LoadDocument(f)
	require.NoError(t, err)

	keys := doc.Keys()
	require.NotEmpty(t, keys)
	for _, k := range keys {
		v, ok := doc.Get(k)
		require.True(t, ok)
		require.NotNil(t, v)
	}
}

func TestChange(t *testing.T) {
	doc := NewDocument()

	err := doc.Change(func(txn *Txn) error {
		txn.Set("foo", true)
		txn.Map("bar").Map("baz").Set("qux", true)
		return nil
	})
	require.NoError(t, err)

	// read back scalars
	v, ok := doc.Get("foo")
	require.True(t, ok)
	require.Equal(t, true, v.(BoolView).Value())

	// read back nested maps via Get + type assertion
	barV, ok := doc.Get("bar")
	require.True(t, ok)
	bar := barV.(MapView)
	bazV, ok := bar.Get("baz")
	require.True(t, ok)
	baz := bazV.(MapView)
	qux, ok := baz.Get("qux")
	require.True(t, ok)
	require.Equal(t, true, qux.(BoolView).Value())
}

func TestScalarViews(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("i", int64(-42))
		txn.Set("u", uint64(99))
		txn.Set("f", float64(3.14))
		txn.Set("b", []byte{1, 2, 3})
		txn.Set("n", nil)
		return nil
	}))

	v, ok := doc.Get("i")
	require.True(t, ok)
	require.Equal(t, int64(-42), v.(Int64View).Value())

	v, ok = doc.Get("u")
	require.True(t, ok)
	require.Equal(t, uint64(99), v.(Uint64View).Value())

	v, ok = doc.Get("f")
	require.True(t, ok)
	require.InDelta(t, 3.14, v.(Float64View).Value(), 1e-9)

	v, ok = doc.Get("b")
	require.True(t, ok)
	require.Equal(t, []byte{1, 2, 3}, v.(BytesView).Value())

	v, ok = doc.Get("n")
	require.True(t, ok)
	require.IsType(t, NullView{}, v)
}

func TestTextLen(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Text("t").Splice(0, 0, "héllo") // 5 runes, 6 bytes
		return nil
	}))
	tv, ok := doc.Text("t")
	require.True(t, ok)
	require.Equal(t, 5, tv.Len())
	require.Equal(t, "héllo", tv.Value())
}

func TestDocumentText(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Text("body").Splice(0, 0, "hello")
		return nil
	}))

	tv, ok := doc.Text("body")
	require.True(t, ok)
	require.Equal(t, "hello", tv.Value())

	_, ok = doc.Text("missing")
	require.False(t, ok)

	// Non-text key returns false.
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("scalar", "x")
		return nil
	}))
	_, ok = doc.Text("scalar")
	require.False(t, ok)
}

func TestChangeTypeConflict(t *testing.T) {
	doc := NewDocument()

	// set foo as a bool
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("foo", true)
		return nil
	}))

	// try to treat foo as a map — should return ErrType
	err := doc.Change(func(txn *Txn) error {
		txn.Map("foo").Set("bar", 42) // foo is bool, not map
		return nil
	})
	require.Error(t, err)
	require.IsType(t, ErrType{}, err)
}

func TestMarshalJSON(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("name", "alice")
		txn.Set("active", true)
		txn.Map("meta").Set("version", int64(1))
		return nil
	}))

	data, err := json.Marshal(doc)
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(data, &got))
	require.Equal(t, "alice", got["name"])
	require.Equal(t, true, got["active"])
}

func TestSaveIncremental(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("x", int64(1))
		return nil
	}))
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("y", int64(2))
		return nil
	}))

	// two changes should have produced two unsaved chunks
	require.Len(t, doc.unsaved, 2)

	var buf []byte
	w := &appendWriter{&buf}
	require.NoError(t, doc.SaveIncremental(w))
	require.NotEmpty(t, buf)
	require.Empty(t, doc.unsaved)

	// loading the incremental chunks on top of a fresh doc should reproduce state
	doc2, err := LoadDocument(&byteReader{buf, 0})
	require.NoError(t, err)
	v, ok := doc2.Get("x")
	require.True(t, ok)
	_ = v
}

func TestSave(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("name", "alice")
		txn.Set("active", true)
		txn.Map("meta").Set("version", int64(1))
		return nil
	}))
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("score", int64(42))
		return nil
	}))

	var buf []byte
	w := &appendWriter{&buf}
	require.NoError(t, doc.Save(w))
	require.NotEmpty(t, buf)
	require.Empty(t, doc.unsaved) // Save clears the incremental buffer

	// Reload from the saved document chunk and verify state is preserved.
	doc2, err := LoadDocument(&byteReader{buf, 0})
	require.NoError(t, err)

	v, ok := doc2.Get("name")
	require.True(t, ok)
	require.Equal(t, "alice", v.(StringView).Value())

	v, ok = doc2.Get("active")
	require.True(t, ok)
	require.Equal(t, true, v.(BoolView).Value())

	v, ok = doc2.Get("score")
	require.True(t, ok)

	meta, ok := doc2.Map("meta")
	require.True(t, ok)
	v, ok = meta.Get("version")
	require.True(t, ok)
	_ = v
}

func TestSaveLoadedDocument(t *testing.T) {
	// Load a pre-existing document chunk, apply a change, save, reload, verify.
	f, err := os.Open("../testdata/two_change_chunks.automerge")
	require.NoError(t, err)
	defer f.Close()

	doc, err := LoadDocument(f)
	require.NoError(t, err)

	// Apply a new change on top.
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("savedKey", "savedValue")
		return nil
	}))

	var buf []byte
	require.NoError(t, doc.Save(&appendWriter{&buf}))
	require.NotEmpty(t, buf)

	doc2, err := LoadDocument(&byteReader{buf, 0})
	require.NoError(t, err)

	v, ok := doc2.Get("savedKey")
	require.True(t, ok)
	require.Equal(t, "savedValue", v.(StringView).Value())
}

func TestSaveRoundTripFurtherChanges(t *testing.T) {
	// Verify that a document saved via Save() can receive further changes after reload.
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("x", int64(1))
		return nil
	}))

	var buf []byte
	require.NoError(t, doc.Save(&appendWriter{&buf}))

	doc2, err := LoadDocument(&byteReader{buf, 0})
	require.NoError(t, err)

	// Apply a new change after loading from saved document chunk.
	require.NoError(t, doc2.Change(func(txn *Txn) error {
		txn.Set("y", int64(2))
		return nil
	}))

	v, ok := doc2.Get("x")
	require.True(t, ok)
	_ = v
	v, ok = doc2.Get("y")
	require.True(t, ok)
	_ = v

	// Merge back into the original doc and verify y is visible.
	require.NoError(t, doc.Merge(doc2))
	v, ok = doc.Get("y")
	require.True(t, ok)
	_ = v
}

func TestTextSpliceSingleTransaction(t *testing.T) {
	// Multiple splices within a single transaction must see each other's ops
	// (write-your-own-reads via Transaction.ListElements).
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		tv := txn.Text("body")
		tv.Splice(0, 0, "Hello")
		tv.Splice(5, 0, ", world!")
		tv.Splice(7, 5, "Go") // replace "world" with "Go"
		return nil
	}))
	v, ok := doc.Get("body")
	require.True(t, ok)
	require.Equal(t, "Hello, Go!", v.(TextView).Value())
}

func TestTextSplice(t *testing.T) {
	doc := NewDocument()

	// Create a text object and insert initial content.
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Text("body").Splice(0, 0, "Hello, world!")
		return nil
	}))

	tv, ok := doc.Get("body")
	require.True(t, ok)
	require.Equal(t, "Hello, world!", tv.(TextView).Value())

	// Replace "world" with "Go".
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Text("body").Splice(7, 5, "Go")
		return nil
	}))
	tv, _ = doc.Get("body")
	require.Equal(t, "Hello, Go!", tv.(TextView).Value())

	// Delete the comma and space (positions 5 and 6), leaving "HelloGo!".
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Text("body").Splice(5, 2, "")
		return nil
	}))
	tv, _ = doc.Get("body")
	require.Equal(t, "HelloGo!", tv.(TextView).Value())
}

func TestTextUpdate(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Text("t").Splice(0, 0, "abcdef")
		return nil
	}))

	// Update: change middle, keep prefix "ab" and suffix "f".
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Text("t").Update("abXYZf")
		return nil
	}))
	v, _ := doc.Get("t")
	require.Equal(t, "abXYZf", v.(TextView).Value())

	// Update to identical string: no-op.
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Text("t").Update("abXYZf")
		return nil
	}))
	v, _ = doc.Get("t")
	require.Equal(t, "abXYZf", v.(TextView).Value())
}

func TestTextMergeConcurrent(t *testing.T) {
	// Two peers insert at different positions concurrently; after merge both
	// edits should be visible.
	base := NewDocument()
	require.NoError(t, base.Change(func(txn *Txn) error {
		txn.Text("t").Splice(0, 0, "ac")
		return nil
	}))

	// Fork.
	var buf []byte
	require.NoError(t, base.SaveIncremental(&appendWriter{&buf}))
	peer1, err := LoadDocument(&byteReader{buf, 0})
	require.NoError(t, err)
	peer2, err := LoadDocument(&byteReader{buf, 0})
	require.NoError(t, err)

	// peer1 inserts 'b' at position 1 (between 'a' and 'c').
	require.NoError(t, peer1.Change(func(txn *Txn) error {
		txn.Text("t").Splice(1, 0, "b")
		return nil
	}))

	// peer2 appends 'd'.
	require.NoError(t, peer2.Change(func(txn *Txn) error {
		txn.Text("t").Splice(2, 0, "d")
		return nil
	}))

	// Merge peer2 into peer1.
	require.NoError(t, peer1.Merge(peer2))

	v, _ := peer1.Get("t")
	// Both insertions should be present; exact order depends on actor IDs.
	text := v.(TextView).Value()
	require.Len(t, []rune(text), 4)
	require.Contains(t, text, "a")
	require.Contains(t, text, "b")
	require.Contains(t, text, "c")
	require.Contains(t, text, "d")
}

func TestNewDocumentFromJSON(t *testing.T) {
	doc, err := NewDocumentFromJSON([]byte(`{"name":"bob","score":42}`))
	require.NoError(t, err)

	v, ok := doc.Get("name")
	require.True(t, ok)
	require.Equal(t, "bob", v.(StringView).Value())
}

// -- Values iterators --------------------------------------------------------

func TestDocumentValuesIterator(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("a", "alpha")
		txn.Set("b", int64(2))
		txn.Map("m").Set("x", true)
		return nil
	}))

	got := make(map[string]Value)
	for k, v := range doc.Values() {
		got[k] = v
	}
	require.Len(t, got, 3)
	require.Equal(t, "alpha", got["a"].(StringView).Value())
	require.Equal(t, int64(2), got["b"].(Int64View).Value())
	_, isMap := got["m"].(MapView)
	require.True(t, isMap)
}

func TestMapViewValuesIterator(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		m := txn.Map("config")
		m.Set("debug", true)
		m.Set("level", int64(3))
		return nil
	}))

	mv, ok := doc.Map("config")
	require.True(t, ok)

	got := make(map[string]Value)
	for k, v := range mv.Values() {
		got[k] = v
	}
	require.Len(t, got, 2)
	require.Equal(t, true, got["debug"].(BoolView).Value())
	require.Equal(t, int64(3), got["level"].(Int64View).Value())
}

func TestMapViewKeys(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		m := txn.Map("ns")
		m.Set("z", "last")
		m.Set("a", "first")
		return nil
	}))

	mv, ok := doc.Map("ns")
	require.True(t, ok)
	keys := mv.Keys()
	require.ElementsMatch(t, []string{"a", "z"}, keys)
}

func TestTxnValuesAndKeys(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("x", int64(1))
		txn.Set("y", int64(2))
		return nil
	}))

	require.NoError(t, doc.Change(func(txn *Txn) error {
		// Keys reflects committed state inside the transaction.
		keys := txn.Keys()
		require.ElementsMatch(t, []string{"x", "y"}, keys)

		// Values reflects committed state too.
		got := make(map[string]Value)
		for k, v := range txn.Values() {
			got[k] = v
		}
		require.Len(t, got, 2)
		require.Equal(t, int64(1), got["x"].(Int64View).Value())
		return nil
	}))
}

// -- ListView access ---------------------------------------------------------

func TestListViewGetAndValues(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		l := txn.List("items")
		l.Append("x")
		l.Append(int64(42))
		l.Append(true)
		return nil
	}))

	lv, ok := doc.List("items")
	require.True(t, ok)
	require.Equal(t, 3, lv.Len())

	// Get by index.
	v, ok := lv.Get(0)
	require.True(t, ok)
	require.Equal(t, "x", v.(StringView).Value())

	v, ok = lv.Get(1)
	require.True(t, ok)
	require.Equal(t, int64(42), v.(Int64View).Value())

	v, ok = lv.Get(2)
	require.True(t, ok)
	require.Equal(t, true, v.(BoolView).Value())

	_, ok = lv.Get(3)
	require.False(t, ok)

	// Values iterator.
	var vals []Value
	for _, v := range lv.Values() {
		vals = append(vals, v)
	}
	require.Len(t, vals, 3)
	require.Equal(t, "x", vals[0].(StringView).Value())
	require.Equal(t, int64(42), vals[1].(Int64View).Value())
}

// -- Manual Begin / Rollback -------------------------------------------------

func TestBeginRollback(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("committed", "yes")
		return nil
	}))

	txn := doc.Begin()
	txn.Set("uncommitted", "never")
	txn.Rollback()

	// Rollback must be idempotent.
	txn.Rollback()

	_, ok := doc.Get("uncommitted")
	require.False(t, ok, "rolled-back key must not be visible")

	v, ok := doc.Get("committed")
	require.True(t, ok)
	require.Equal(t, "yes", v.(StringView).Value())
}

func TestEmptyChange(t *testing.T) {
	doc := NewDocument()
	before := len(doc.allChanges)
	require.NoError(t, doc.Change(func(txn *Txn) error {
		// no ops
		return nil
	}))
	// An empty transaction must not produce a change record.
	require.Equal(t, before, len(doc.allChanges))
}

// -- MapView write methods ---------------------------------------------------

func TestMapViewDelete(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		m := txn.Map("cfg")
		m.Set("keep", "yes")
		m.Set("remove", "no")
		return nil
	}))

	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Map("cfg").Delete("remove")
		return nil
	}))

	mv, ok := doc.Map("cfg")
	require.True(t, ok)
	_, ok = mv.Get("remove")
	require.False(t, ok)
	v, ok := mv.Get("keep")
	require.True(t, ok)
	require.Equal(t, "yes", v.(StringView).Value())
}

// -- Counter CRDT ------------------------------------------------------------

func TestCounterMerge(t *testing.T) {
	// Two peers start from the same state, each increment the counter
	// independently, then merge. The result must be the sum of both deltas.
	base := NewDocument()
	require.NoError(t, base.Change(func(txn *Txn) error {
		txn.Set("hits", types.Counter(0))
		return nil
	}))

	var snap []byte
	require.NoError(t, base.SaveIncremental(&appendWriter{&snap}))

	peer1, err := LoadDocument(&byteReader{snap, 0})
	require.NoError(t, err)
	peer2, err := LoadDocument(&byteReader{snap, 0})
	require.NoError(t, err)

	require.NoError(t, peer1.Change(func(txn *Txn) error {
		txn.Increment("hits", 3)
		return nil
	}))

	require.NoError(t, peer2.Change(func(txn *Txn) error {
		txn.Increment("hits", 5)
		return nil
	}))

	require.NoError(t, peer1.Merge(peer2))

	v, ok := As[int64](peer1.Get("hits"))
	require.True(t, ok)
	require.Equal(t, int64(8), v)
}

func TestCounterIncrementalSave(t *testing.T) {
	// Save counter state via change chunks (not a full document snapshot),
	// reload, and verify the value is correct.
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("c", types.Counter(10))
		return nil
	}))
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Increment("c", 4)
		return nil
	}))

	var buf []byte
	require.NoError(t, doc.SaveIncremental(&appendWriter{&buf}))

	doc2, err := LoadDocument(&byteReader{buf, 0})
	require.NoError(t, err)

	v, ok := As[int64](doc2.Get("c"))
	require.True(t, ok)
	require.Equal(t, int64(14), v)
}

// -- helpers -----------------------------------------------------------------

type appendWriter struct{ b *[]byte }

func (w *appendWriter) Write(p []byte) (int, error) {
	*w.b = append(*w.b, p...)
	return len(p), nil
}

type byteReader struct {
	b   []byte
	pos int
}

func (r *byteReader) Read(p []byte) (int, error) {
	if r.pos >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.pos:])
	r.pos += n
	return n, nil
}
