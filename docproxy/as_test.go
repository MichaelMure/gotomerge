package docproxy

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/gotomerge/types"
)

// -- Native() -----------------------------------------------------------------

func TestNativeScalars(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("s", "hello")
		txn.Set("b", true)
		txn.Set("i", int64(-7))
		txn.Set("u", uint64(9))
		txn.Set("f", float64(2.5))
		txn.Set("n", nil)
		return nil
	}))

	cases := []struct {
		key  string
		want any
	}{
		{"s", "hello"},
		{"b", true},
		{"i", int64(-7)},
		{"u", uint64(9)},
		{"f", float64(2.5)},
		{"n", nil},
	}
	for _, c := range cases {
		v, ok := doc.Get(c.key)
		require.True(t, ok, c.key)
		require.Equal(t, c.want, v.Native(), c.key)
	}
}

func TestNativeMap(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		m := txn.Map("m")
		m.Set("x", "hello")
		m.Set("y", int64(1))
		return nil
	}))
	v, ok := doc.Get("m")
	require.True(t, ok)
	m, ok := v.Native().(map[string]any)
	require.True(t, ok)
	require.Equal(t, "hello", m["x"])
	require.Equal(t, int64(1), m["y"])
}

func TestNativeList(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		l := txn.List("l")
		l.Append("a")
		l.Append("b")
		return nil
	}))
	v, ok := doc.Get("l")
	require.True(t, ok)
	s, ok := v.Native().([]any)
	require.True(t, ok)
	require.Equal(t, []any{"a", "b"}, s)
}

func TestNativeText(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Text("t").Splice(0, 0, "hi")
		return nil
	}))
	v, ok := doc.Get("t")
	require.True(t, ok)
	require.Equal(t, "hi", v.Native())
}

// -- As[T] scalars ------------------------------------------------------------

func TestAsString(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("s", "world")
		txn.Text("t").Splice(0, 0, "text")
		return nil
	}))

	s, ok := As[string](doc.Get("s"))
	require.True(t, ok)
	require.Equal(t, "world", s)

	s, ok = As[string](doc.Get("t"))
	require.True(t, ok)
	require.Equal(t, "text", s)

	_, ok = As[string](doc.Get("missing"))
	require.False(t, ok)
}

func TestAsBool(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("b", true)
		return nil
	}))

	b, ok := As[bool](doc.Get("b"))
	require.True(t, ok)
	require.True(t, b)
}

func TestAsInt64(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("i", int64(42))
		txn.Set("u", uint64(10))
		txn.Set("f", float64(3.0))
		return nil
	}))

	i, ok := As[int64](doc.Get("i"))
	require.True(t, ok)
	require.Equal(t, int64(42), i)

	// uint64 → int64
	i, ok = As[int64](doc.Get("u"))
	require.True(t, ok)
	require.Equal(t, int64(10), i)

	// float64 → int64
	i, ok = As[int64](doc.Get("f"))
	require.True(t, ok)
	require.Equal(t, int64(3), i)
}

func TestAsIntNarrowing(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("small", int64(100))
		txn.Set("big", int64(1000))
		return nil
	}))

	// int64(100) fits in int8 (-128..127)
	v, ok := As[int8](doc.Get("small"))
	require.True(t, ok)
	require.Equal(t, int8(100), v)

	// int64(1000) does not fit in int8
	_, ok = As[int8](doc.Get("big"))
	require.False(t, ok)
}

func TestAsUintNegativeRejected(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("neg", int64(-1))
		return nil
	}))
	_, ok := As[uint64](doc.Get("neg"))
	require.False(t, ok)
}

func TestAsFloat64(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("f", float64(1.5))
		txn.Set("i", int64(2))
		return nil
	}))

	f, ok := As[float64](doc.Get("f"))
	require.True(t, ok)
	require.InDelta(t, 1.5, f, 1e-9)

	f, ok = As[float64](doc.Get("i"))
	require.True(t, ok)
	require.InDelta(t, 2.0, f, 1e-9)
}

func TestAsTimestamp(t *testing.T) {
	now := time.Now().Truncate(time.Millisecond)
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("ts", types.FromTime(now)) // store as types.Timestamp
		return nil
	}))
	ts, ok := As[time.Time](doc.Get("ts"))
	require.True(t, ok)
	require.Equal(t, now, ts)
}

// -- As[T] null ---------------------------------------------------------------

func TestAsNull(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("n", nil)
		return nil
	}))

	// NullView → concrete target: false
	_, ok := As[string](doc.Get("n"))
	require.False(t, ok)

	// NullView → pointer target: nil, true
	ptr, ok := As[*string](doc.Get("n"))
	require.True(t, ok)
	require.Nil(t, ptr)
}

// -- As[T] struct unmarshaling ------------------------------------------------

type asTestConfig struct {
	Debug   bool   `automerge:"debug"`
	Version int64  `automerge:"version"`
	Label   string `automerge:"label"`
	Ignored string `automerge:"-"`
	hidden  string // nolint unexported
}

func TestAsStruct(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		m := txn.Map("cfg")
		m.Set("debug", true)
		m.Set("version", int64(3))
		m.Set("label", "prod")
		m.Set("Ignored", "should not appear") // key matches field name, not tag
		return nil
	}))

	cfg, ok := As[asTestConfig](doc.Get("cfg"))
	require.True(t, ok)
	require.True(t, cfg.Debug)
	require.Equal(t, int64(3), cfg.Version)
	require.Equal(t, "prod", cfg.Label)
	require.Empty(t, cfg.Ignored) // tag "-" skips it
}

func TestAsStructPointer(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		m := txn.Map("cfg")
		m.Set("debug", false)
		m.Set("version", int64(1))
		m.Set("label", "dev")
		return nil
	}))

	cfg, ok := As[*asTestConfig](doc.Get("cfg"))
	require.True(t, ok)
	require.NotNil(t, cfg)
	require.Equal(t, int64(1), cfg.Version)

	// Missing key → nil pointer
	ptr, ok := As[*asTestConfig](doc.Get("missing"))
	require.False(t, ok)
	require.Nil(t, ptr)
}

func TestAsStructMissingField(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Map("cfg").Set("debug", true)
		// version and label absent
		return nil
	}))

	cfg, ok := As[asTestConfig](doc.Get("cfg"))
	require.True(t, ok)
	require.True(t, cfg.Debug)
	require.Zero(t, cfg.Version) // absent → zero
	require.Empty(t, cfg.Label)  // absent → zero
}

// -- As[T] slice and map ------------------------------------------------------

func TestAsSliceString(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		l := txn.List("tags")
		l.Append("a")
		l.Append("b")
		l.Append("c")
		return nil
	}))

	tags, ok := As[[]string](doc.Get("tags"))
	require.True(t, ok)
	require.Equal(t, []string{"a", "b", "c"}, tags)
}

func TestAsSliceAny(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		l := txn.List("items")
		l.Append("x")
		l.Append(int64(1))
		return nil
	}))

	items, ok := As[[]any](doc.Get("items"))
	require.True(t, ok)
	require.Equal(t, []any{"x", int64(1)}, items)
}

func TestAsMapStringString(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		m := txn.Map("labels")
		m.Set("env", "prod")
		m.Set("region", "us-east")
		return nil
	}))

	labels, ok := As[map[string]string](doc.Get("labels"))
	require.True(t, ok)
	require.Equal(t, "prod", labels["env"])
	require.Equal(t, "us-east", labels["region"])
}

func TestAsMapStringAny(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		m := txn.Map("meta")
		m.Set("name", "alice")
		m.Set("count", int64(5))
		return nil
	}))

	meta, ok := As[map[string]any](doc.Get("meta"))
	require.True(t, ok)
	require.Equal(t, "alice", meta["name"])
	require.Equal(t, int64(5), meta["count"])
}

// -- Counter ------------------------------------------------------------------

func TestCounterIncrement(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("c", types.Counter(10))
		return nil
	}))

	// Initial value via As[int64].
	v, ok := As[int64](doc.Get("c"))
	require.True(t, ok)
	require.Equal(t, int64(10), v)

	// Increment.
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Increment("c", 5)
		return nil
	}))

	v, ok = As[int64](doc.Get("c"))
	require.True(t, ok)
	require.Equal(t, int64(15), v)

	// Decrement.
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Increment("c", -3)
		return nil
	}))

	v, ok = As[int64](doc.Get("c"))
	require.True(t, ok)
	require.Equal(t, int64(12), v)
}

func TestCounterSaveReload(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("c", types.Counter(0))
		return nil
	}))
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Increment("c", 7)
		return nil
	}))

	// Save full snapshot and reload.
	var buf bytes.Buffer
	require.NoError(t, doc.Save(&buf))

	doc2, err := LoadDocument(&buf)
	require.NoError(t, err)

	v, ok := As[int64](doc2.Get("c"))
	require.True(t, ok)
	require.Equal(t, int64(7), v)
}

// -- As[any] ------------------------------------------------------------------

func TestAsAny(t *testing.T) {
	doc := NewDocument()
	require.NoError(t, doc.Change(func(txn *Txn) error {
		txn.Set("s", "hello")
		return nil
	}))

	v, ok := As[any](doc.Get("s"))
	require.True(t, ok)
	require.Equal(t, "hello", v)
}
