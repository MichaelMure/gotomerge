package gotomerge_test

import (
	"bytes"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/gotomerge"
	"github.com/MichaelMure/gotomerge/types"
)

func TestRootAPI_BasicReadWrite(t *testing.T) {
	doc := gotomerge.NewDocument()

	require.NoError(t, doc.Change(func(txn *gotomerge.Txn) error {
		txn.Set("title", "Hello")
		txn.Set("published", true)
		txn.Set("score", int64(42))
		m := txn.Map("meta")
		m.Set("author", "alice")
		m.Set("version", int64(1))
		return nil
	}))

	title, ok := gotomerge.As[string](doc.Get("title"))
	require.True(t, ok)
	require.Equal(t, "Hello", title)

	pub, ok := gotomerge.As[bool](doc.Get("published"))
	require.True(t, ok)
	require.True(t, pub)

	type Meta struct {
		Author  string
		Version int64
	}
	meta, ok := gotomerge.As[Meta](doc.Get("meta"))
	require.True(t, ok)
	require.Equal(t, "alice", meta.Author)
	require.Equal(t, int64(1), meta.Version)
}

func TestRootAPI_ListAndText(t *testing.T) {
	doc := gotomerge.NewDocument()

	require.NoError(t, doc.Change(func(txn *gotomerge.Txn) error {
		l := txn.List("tags")
		l.Append("go")
		l.Append("crdt")
		txn.Text("body").Splice(0, 0, "Hello, world!")
		return nil
	}))

	tags, ok := gotomerge.As[[]string](doc.Get("tags"))
	require.True(t, ok)
	require.Equal(t, []string{"go", "crdt"}, tags)

	body, ok := gotomerge.As[string](doc.Get("body"))
	require.True(t, ok)
	require.Equal(t, "Hello, world!", body)
}

func TestRootAPI_Counter(t *testing.T) {
	doc := gotomerge.NewDocument()

	require.NoError(t, doc.Change(func(txn *gotomerge.Txn) error {
		txn.Set("hits", types.Counter(0))
		return nil
	}))
	require.NoError(t, doc.Change(func(txn *gotomerge.Txn) error {
		txn.Increment("hits", 5)
		return nil
	}))

	hits, ok := gotomerge.As[int64](doc.Get("hits"))
	require.True(t, ok)
	require.Equal(t, int64(5), hits)
}

func TestRootAPI_SaveLoadMerge(t *testing.T) {
	base := gotomerge.NewDocument()
	require.NoError(t, base.Change(func(txn *gotomerge.Txn) error {
		txn.Set("shared", "base")
		return nil
	}))

	var snap bytes.Buffer
	base.SaveIncremental(&snap)
	snapBytes := snap.Bytes()

	peer1, _ := gotomerge.LoadDocument(bytes.NewReader(snapBytes))
	peer2, _ := gotomerge.LoadDocument(bytes.NewReader(snapBytes))

	require.NoError(t, peer1.Change(func(txn *gotomerge.Txn) error {
		txn.Set("alice", "peer1")
		return nil
	}))
	require.NoError(t, peer2.Change(func(txn *gotomerge.Txn) error {
		txn.Set("bob", "peer2")
		return nil
	}))

	require.NoError(t, peer1.Merge(peer2))

	keys := peer1.Keys()
	sort.Strings(keys)
	require.Equal(t, []string{"alice", "bob", "shared"}, keys)
}

func TestRootAPI_Timestamp(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	doc := gotomerge.NewDocument()
	require.NoError(t, doc.Change(func(txn *gotomerge.Txn) error {
		txn.Set("ts", types.FromTime(now))
		return nil
	}))
	ts, ok := gotomerge.As[time.Time](doc.Get("ts"))
	require.True(t, ok)
	require.Equal(t, now, ts)
}

func TestRootAPI_JSON(t *testing.T) {
	doc, err := gotomerge.NewDocumentFromJSON([]byte(`{"name":"bob","active":true}`))
	require.NoError(t, err)

	name, ok := gotomerge.As[string](doc.Get("name"))
	require.True(t, ok)
	require.Equal(t, "bob", name)
}

func TestRootAPI_MapViewOperations(t *testing.T) {
	doc := gotomerge.NewDocument()
	require.NoError(t, doc.Change(func(txn *gotomerge.Txn) error {
		m := txn.Map("cfg")
		m.Set("debug", true)
		m.Set("port", int64(8080))
		return nil
	}))

	mv, ok := doc.Map("cfg")
	require.True(t, ok)

	debug, ok := gotomerge.As[bool](mv.Get("debug"))
	require.True(t, ok)
	require.True(t, debug)

	got := make(map[string]any)
	for k, v := range mv.Values() {
		got[k] = v.Native()
	}
	require.Equal(t, true, got["debug"])
	require.Equal(t, int64(8080), got["port"])
}

func TestRootAPI_ListViewOperations(t *testing.T) {
	doc := gotomerge.NewDocument()
	require.NoError(t, doc.Change(func(txn *gotomerge.Txn) error {
		l := txn.List("items")
		l.Append("a")
		l.Append("b")
		l.Append("c")
		return nil
	}))

	lv, ok := doc.List("items")
	require.True(t, ok)
	require.Equal(t, 3, lv.Len())

	first, ok := gotomerge.As[string](lv.Get(0))
	require.True(t, ok)
	require.Equal(t, "a", first)
}
