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

// -- BenchmarkGet ------------------------------------------------------------
// Hot-path read: doc.Get on an existing key, then As[T] to extract the value.

func BenchmarkGet(b *testing.B) {
	sizes := []int{100, 1_000, 10_000}

	b.Run("scalar", func(b *testing.B) {
		for _, n := range sizes {
			b.Run(fmt.Sprintf("keys=%d", n), func(b *testing.B) {
				doc := buildIncreasingPut(b, n)
				key := strconv.Itoa(n / 2) // middle key
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					v, ok := doc.Get(key)
					if !ok {
						b.Fatal("key not found")
					}
					_ = v
				}
			})
		}
	})

	// nested map: get from a MapView, one level deep
	b.Run("nested", func(b *testing.B) {
		for _, n := range sizes {
			b.Run(fmt.Sprintf("keys=%d", n), func(b *testing.B) {
				doc := NewDocument()
				require.NoError(b, doc.Change(func(txn *Txn) error {
					m := txn.Map("outer")
					for i := 0; i < n; i++ {
						m.Set(strconv.Itoa(i), int64(i))
					}
					return nil
				}))
				mv, _ := doc.Map("outer")
				key := strconv.Itoa(n / 2)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					v, ok := mv.Get(key)
					if !ok {
						b.Fatal("key not found")
					}
					_ = v
				}
			})
		}
	})
}

// -- BenchmarkAs -------------------------------------------------------------
// As[T] conversion: scalar, struct unmarshaling (reflection), and slice.

type benchStruct struct {
	X int64
	Y int64
	Z string
}

func BenchmarkAs(b *testing.B) {
	doc := NewDocument()
	require.NoError(b, doc.Change(func(txn *Txn) error {
		txn.Set("i", int64(42))
		txn.Set("s", "hello")
		m := txn.Map("obj")
		m.Set("x", int64(1))
		m.Set("y", int64(2))
		m.Set("z", "three")
		l := txn.List("list")
		for i := 0; i < 100; i++ {
			l.Append(strconv.Itoa(i))
		}
		return nil
	}))

	b.Run("int64", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = As[int64](doc.Get("i"))
		}
	})

	b.Run("string", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = As[string](doc.Get("s"))
		}
	})

	b.Run("struct", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = As[benchStruct](doc.Get("obj"))
		}
	})

	b.Run("slice_string/len=100", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, _ = As[[]string](doc.Get("list"))
		}
	})
}

// -- BenchmarkMerge ----------------------------------------------------------
// Merge two documents: each peer starts from the same snapshot and applies
// n independent changes; then one merges the other.

func BenchmarkMerge(b *testing.B) {
	sizes := []int{10, 100, 1_000}

	for _, n := range sizes {
		b.Run(fmt.Sprintf("changes=%d", n), func(b *testing.B) {
			// Build and save the shared base once.
			base := NewDocument()
			require.NoError(b, base.Change(func(txn *Txn) error {
				txn.Set("base", int64(0))
				return nil
			}))
			snap := saveDocIncremental(b, base)

			// Build peer2's changes once (they are constant across iterations).
			peer2, err := LoadDocument(bytes.NewReader(snap))
			require.NoError(b, err)
			for i := 0; i < n; i++ {
				i := i
				require.NoError(b, peer2.Change(func(txn *Txn) error {
					txn.Set("p2-"+strconv.Itoa(i), int64(i))
					return nil
				}))
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				peer1, err := LoadDocument(bytes.NewReader(snap))
				require.NoError(b, err)
				for j := 0; j < n; j++ {
					j := j
					require.NoError(b, peer1.Change(func(txn *Txn) error {
						txn.Set("p1-"+strconv.Itoa(j), int64(j))
						return nil
					}))
				}
				b.StartTimer()
				require.NoError(b, peer1.Merge(peer2))
			}
		})
	}
}

// -- BenchmarkList -----------------------------------------------------------
// List-heavy documents: build, save, load, and iterate.

func buildList(tb testing.TB, n int) *Document {
	tb.Helper()
	doc := NewDocument()
	require.NoError(tb, doc.Change(func(txn *Txn) error {
		l := txn.List("items")
		for i := 0; i < n; i++ {
			l.Append(int64(i))
		}
		return nil
	}))
	return doc
}

func BenchmarkList(b *testing.B) {
	sizes := []int{100, 1_000, 10_000}

	b.Run("build", func(b *testing.B) {
		for _, n := range sizes {
			b.Run(fmt.Sprintf("len=%d", n), func(b *testing.B) {
				b.ReportAllocs()
				for i := 0; i < b.N; i++ {
					_ = buildList(b, n)
				}
			})
		}
	})

	b.Run("save", func(b *testing.B) {
		for _, n := range sizes {
			b.Run(fmt.Sprintf("len=%d", n), func(b *testing.B) {
				b.ReportAllocs()
				var lastLen int
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					doc := buildList(b, n)
					b.StartTimer()
					var buf bytes.Buffer
					require.NoError(b, doc.Save(&buf))
					lastLen = buf.Len()
				}
				b.SetBytes(int64(lastLen))
			})
		}
	})

	b.Run("load", func(b *testing.B) {
		for _, n := range sizes {
			b.Run(fmt.Sprintf("len=%d", n), func(b *testing.B) {
				saved := saveDoc(b, buildList(b, n))
				b.SetBytes(int64(len(saved)))
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					_, err := LoadDocument(bytes.NewReader(saved))
					require.NoError(b, err)
				}
			})
		}
	})

	b.Run("iterate", func(b *testing.B) {
		for _, n := range sizes {
			b.Run(fmt.Sprintf("len=%d", n), func(b *testing.B) {
				doc := buildList(b, n)
				lv, _ := doc.List("items")
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					for range lv.Values() {
					}
				}
			})
		}
	})

	b.Run("as_slice/len=1000", func(b *testing.B) {
		doc := buildList(b, 1_000)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = As[[]int64](doc.Get("items"))
		}
	})
}

// -- BenchmarkTextRead -------------------------------------------------------
// Reading a large text document as a string. Complements BenchmarkEditTrace
// which only measures load; this measures the string-extraction hot path.

func BenchmarkTextRead(b *testing.B) {
	b.Run("poorly_simulated_typing/n=1000", func(b *testing.B) {
		doc := buildPoorlySimulatedTyping(b, 1_000)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = As[string](doc.Get("content"))
		}
	})

	b.Run("edit_trace", func(b *testing.B) {
		doc, _ := loadEditTraceDoc(b)
		if doc == nil {
			b.Skip("text-edits.amrg not found")
		}
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			tv, _ := doc.Text("text")
			_ = tv.Value()
		}
	})
}

// -- BenchmarkSaveIncremental ------------------------------------------------
// Incremental save (change chunks only) vs full save, to show the cost
// differential for hot write paths.

func BenchmarkSaveIncremental(b *testing.B) {
	const changesPerIter = 10

	b.Run("incremental", func(b *testing.B) {
		b.ReportAllocs()
		var lastLen int
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			doc := buildDeepHistory(b, changesPerIter)
			b.StartTimer()
			var buf bytes.Buffer
			require.NoError(b, doc.SaveIncremental(&buf))
			lastLen = buf.Len()
		}
		b.SetBytes(int64(lastLen))
	})

	b.Run("full", func(b *testing.B) {
		b.ReportAllocs()
		var lastLen int
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			doc := buildDeepHistory(b, changesPerIter)
			b.StartTimer()
			var buf bytes.Buffer
			require.NoError(b, doc.Save(&buf))
			lastLen = buf.Len()
		}
		b.SetBytes(int64(lastLen))
	})
}
