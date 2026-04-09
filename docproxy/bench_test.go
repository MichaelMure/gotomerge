package docproxy

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// -- document builders -------------------------------------------------------

// buildRepeatedPut creates a doc with n puts to the same key "0".
// Equivalent to Rust's repeated_put(n).
func buildRepeatedPut(tb testing.TB, n int) *Document {
	tb.Helper()
	doc := NewDocument()
	require.NoError(tb, doc.Change(func(txn *Txn) error {
		for i := 0; i < n; i++ {
			txn.Set("0", int64(i))
		}
		return nil
	}))
	return doc
}

// buildIncreasingPut creates a doc with n puts to increasing keys "0".."n-1".
// Equivalent to Rust's increasing_put(n).
func buildIncreasingPut(tb testing.TB, n int) *Document {
	tb.Helper()
	doc := NewDocument()
	require.NoError(tb, doc.Change(func(txn *Txn) error {
		for i := 0; i < n; i++ {
			txn.Set(strconv.Itoa(i), int64(i))
		}
		return nil
	}))
	return doc
}

// buildDecreasingPut creates a doc with n puts to decreasing keys "n-1".."0".
// Equivalent to Rust's decreasing_put(n).
func buildDecreasingPut(tb testing.TB, n int) *Document {
	tb.Helper()
	doc := NewDocument()
	require.NoError(tb, doc.Change(func(txn *Txn) error {
		for i := n - 1; i >= 0; i-- {
			txn.Set(strconv.Itoa(i), int64(i))
		}
		return nil
	}))
	return doc
}

// buildMapsInMaps creates a chain of n nested maps in a single transaction.
// Equivalent to Rust's maps_in_maps_doc(n).
func buildMapsInMaps(tb testing.TB, n int) *Document {
	tb.Helper()
	doc := NewDocument()
	require.NoError(tb, doc.Change(func(txn *Txn) error {
		cur := txn.Map("0")
		for i := 1; i < n; i++ {
			cur = cur.Map(strconv.Itoa(i))
		}
		return nil
	}))
	return doc
}

// buildDeepHistory creates n transactions each setting "x" and "y".
// Equivalent to Rust's deep_history_doc(n).
func buildDeepHistory(tb testing.TB, n int) *Document {
	tb.Helper()
	doc := NewDocument()
	for i := 0; i < n; i++ {
		i := i
		require.NoError(tb, doc.Change(func(txn *Txn) error {
			txn.Set("x", int64(i))
			txn.Set("y", int64(i))
			return nil
		}))
	}
	return doc
}

// buildBigPaste creates a doc with one large string put to "content".
// Equivalent to Rust's big_paste_doc(n).
func buildBigPaste(tb testing.TB, n int) *Document {
	tb.Helper()
	doc := NewDocument()
	require.NoError(tb, doc.Change(func(txn *Txn) error {
		txn.Set("content", strings.Repeat("x", n))
		return nil
	}))
	return doc
}

// buildPoorlySimulatedTyping creates a text object and adds n characters
// one transaction at a time. Equivalent to Rust's poorly_simulated_typing_doc(n).
func buildPoorlySimulatedTyping(tb testing.TB, n int) *Document {
	tb.Helper()
	doc := NewDocument()
	// First tx: create the text object.
	require.NoError(tb, doc.Change(func(txn *Txn) error {
		txn.Text("content")
		return nil
	}))
	for i := 0; i < n; i++ {
		require.NoError(tb, doc.Change(func(txn *Txn) error {
			txn.Text("content").Splice(i, 0, "x")
			return nil
		}))
	}
	return doc
}

func saveDoc(tb testing.TB, doc *Document) []byte {
	tb.Helper()
	var buf bytes.Buffer
	require.NoError(tb, doc.Save(&buf))
	return buf.Bytes()
}

func saveDocIncremental(tb testing.TB, doc *Document) []byte {
	tb.Helper()
	var buf bytes.Buffer
	require.NoError(tb, doc.SaveIncremental(&buf))
	return buf.Bytes()
}

// -- edit trace --------------------------------------------------------------

var (
	editTraceDocOnce sync.Once
	editTraceDoc     *Document
	editTraceDocSave []byte
)

func loadEditTraceDoc(tb testing.TB) (*Document, []byte) {
	tb.Helper()
	editTraceDocOnce.Do(func() {
		data, err := os.ReadFile("../testdata/text-edits.amrg")
		if err != nil {
			tb.Logf("skipping edit trace (text-edits.amrg not found): %v", err)
			return
		}
		editTraceDoc, err = LoadDocument(bytes.NewReader(data))
		if err != nil {
			tb.Fatalf("load edit trace: %v", err)
		}
		editTraceDocSave = saveDoc(tb, editTraceDoc)
	})
	return editTraceDoc, editTraceDocSave
}

// -- BenchmarkMap ------------------------------------------------------------
// Mirrors Rust automerge/benches/map.rs: build, save, load, apply.

func BenchmarkMap(b *testing.B) {
	sizes := []int{100, 1_000, 10_000}
	type variant struct {
		name  string
		build func(testing.TB, int) *Document
	}
	variants := []variant{
		{"repeated_put", buildRepeatedPut},
		{"increasing_put", buildIncreasingPut},
		{"decreasing_put", buildDecreasingPut},
	}

	b.Run("build", func(b *testing.B) {
		for _, v := range variants {
			for _, n := range sizes {
				b.Run(fmt.Sprintf("%s/ops=%d", v.name, n), func(b *testing.B) {
					b.ReportAllocs()
					for i := 0; i < b.N; i++ {
						_ = v.build(b, n)
					}
				})
			}
		}
	})

	b.Run("save", func(b *testing.B) {
		for _, v := range variants {
			for _, n := range sizes {
				b.Run(fmt.Sprintf("%s/ops=%d", v.name, n), func(b *testing.B) {
					b.ReportAllocs()
					var lastLen int
					for i := 0; i < b.N; i++ {
						b.StopTimer()
						doc := v.build(b, n)
						b.StartTimer()
						var buf bytes.Buffer
						require.NoError(b, doc.Save(&buf))
						lastLen = buf.Len()
					}
					b.SetBytes(int64(lastLen))
				})
			}
		}
	})

	b.Run("load", func(b *testing.B) {
		for _, v := range variants {
			for _, n := range sizes {
				b.Run(fmt.Sprintf("%s/ops=%d", v.name, n), func(b *testing.B) {
					b.ReportAllocs()
					saved := saveDoc(b, v.build(b, n))
					b.SetBytes(int64(len(saved)))
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						_, err := LoadDocument(bytes.NewReader(saved))
						require.NoError(b, err)
					}
				})
			}
		}
	})

	// "apply" mirrors Rust's map apply: apply a change chunk (not a doc chunk)
	// to an empty document.
	b.Run("apply", func(b *testing.B) {
		for _, v := range variants {
			for _, n := range sizes {
				b.Run(fmt.Sprintf("%s/ops=%d", v.name, n), func(b *testing.B) {
					b.ReportAllocs()
					saved := saveDocIncremental(b, v.build(b, n))
					b.SetBytes(int64(len(saved)))
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						_, err := LoadDocument(bytes.NewReader(saved))
						require.NoError(b, err)
					}
				})
			}
		}
	})
}

// -- BenchmarkSaveLoad -------------------------------------------------------
// Mirrors Rust automerge/benches/load_save.rs: save+load round-trip.

func BenchmarkSaveLoad(b *testing.B) {
	const n = 1_000

	saveLoad := func(b *testing.B, doc *Document) {
		b.Helper()
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			saved := saveDoc(b, doc)
			_, err := LoadDocument(bytes.NewReader(saved))
			require.NoError(b, err)
		}
	}

	b.Run("big_paste", func(b *testing.B) {
		saveLoad(b, buildBigPaste(b, n))
	})

	b.Run("maps_in_maps", func(b *testing.B) {
		saveLoad(b, buildMapsInMaps(b, n))
	})

	b.Run("deep_history", func(b *testing.B) {
		saveLoad(b, buildDeepHistory(b, n))
	})

	b.Run("poorly_simulated_typing", func(b *testing.B) {
		saveLoad(b, buildPoorlySimulatedTyping(b, n))
	})
}

// -- BenchmarkEditTrace ------------------------------------------------------
// Mirrors Rust automerge/edit-trace/benches/main.rs.
//
// The Rust "replay" benchmark applies all 259,779 edits in a single
// transaction using an O(n log n) positional index inside the opset.
// Our Transaction.ListElements overlay is O(n) per splice (linear scan),
// making a full single-transaction replay O(n²) — too slow to benchmark.
// Transaction.ListElements is correct and tested (TestTextSpliceSingleTransaction);
// it just needs a rope/order-statistics tree to reach Rust's performance for
// very large single transactions.

func BenchmarkEditTrace(b *testing.B) {
	// apply: load the raw change-chunk file (259,779 individual changes).
	// Comparable to Rust's "load" benchmark on the pre-built file.
	b.Run("apply", func(b *testing.B) {
		data, err := os.ReadFile("../testdata/text-edits.amrg")
		require.NoError(b, err)
		b.SetBytes(int64(len(data)))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := LoadDocument(bytes.NewReader(data))
			require.NoError(b, err)
		}
	})

	// save: compact the doc (loaded from amrg) into a single document chunk.
	b.Run("save", func(b *testing.B) {
		doc, _ := loadEditTraceDoc(b)
		if doc == nil {
			b.Skip("text-edits.amrg not found")
		}
		b.ReportAllocs()
		var lastLen int
		for i := 0; i < b.N; i++ {
			var buf bytes.Buffer
			require.NoError(b, doc.Save(&buf))
			lastLen = buf.Len()
		}
		b.SetBytes(int64(lastLen))
	})

	// load: load from the compact saved document chunk.
	b.Run("load", func(b *testing.B) {
		_, saved := loadEditTraceDoc(b)
		if saved == nil {
			b.Skip("text-edits.amrg not found")
		}
		b.SetBytes(int64(len(saved)))
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := LoadDocument(bytes.NewReader(saved))
			require.NoError(b, err)
		}
	})
}

// -- BenchmarkRange ----------------------------------------------------------
// Mirrors Rust automerge/benches/range.rs: iterate all root values.

func BenchmarkRange(b *testing.B) {
	const n = 100_000
	doc := buildIncreasingPut(b, n)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range doc.Values() {
		}
	}
}
