package rle

// RLE encoding quick reference:
//   Repeated run:  s64(L > 1), then one value  → L copies of that value
//   Single value:  s64(1),     then one value  → 1 copy (same path as repeated)
//   Null run:      s64(0),     then u64(count) → count null values
//   Literal run:   s64(L < 0), then -L values  → those values verbatim

import (
	"io"
	"testing"

	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

// value is shorthand for a non-null expected value.
func value[T any](v T) struct {
	val  T
	null bool
} {
	return struct {
		val  T
		null bool
	}{val: v}
}

// null is shorthand for a null expected value.
func null[T any]() struct {
	val  T
	null bool
} {
	return struct {
		val  T
		null bool
	}{null: true}
}

// readAll drains a reader into a slice of (val, null) pairs.
func readAll[T any](t *testing.T, r *Reader[T]) []struct {
	val  T
	null bool
} {
	t.Helper()
	var res []struct {
		val  T
		null bool
	}
	for {
		nv, err := r.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		v, valid := nv.Value()
		res = append(res, struct {
			val  T
			null bool
		}{val: v, null: !valid})
	}
	return res
}

// suiteRLEReader applies the full set of contract tests to a reader type.
// newReader must return a fresh reader over the given raw encoded bytes.
// Each case provides encoded bytes and the expected decoded sequence.
func suiteRLEReader[T comparable](
	t *testing.T,
	newReader func([]byte) *Reader[T],
	encodedEmpty []byte,
	encodedRepeated []byte, repeatedExpected []struct {
		val  T
		null bool
	},
	encodedSingle []byte, singleExpected []struct {
		val  T
		null bool
	},
	encodedNull []byte, nullExpected []struct {
		val  T
		null bool
	},
	encodedLiteral []byte, literalExpected []struct {
		val  T
		null bool
	},
	encodedMixed []byte, mixedExpected []struct {
		val  T
		null bool
	},
) {
	t.Helper()

	t.Run("Empty_ReturnsEOF", func(t *testing.T) {
		r := newReader(encodedEmpty)
		_, err := r.Next()
		require.Equal(t, io.EOF, err)
	})

	t.Run("Repeated", func(t *testing.T) {
		r := newReader(encodedRepeated)
		require.Equal(t, repeatedExpected, readAll(t, r))
	})

	t.Run("Single", func(t *testing.T) {
		r := newReader(encodedSingle)
		require.Equal(t, singleExpected, readAll(t, r))
	})

	t.Run("NullRun", func(t *testing.T) {
		r := newReader(encodedNull)
		require.Equal(t, nullExpected, readAll(t, r))
	})

	t.Run("Literal", func(t *testing.T) {
		r := newReader(encodedLiteral)
		require.Equal(t, literalExpected, readAll(t, r))
	})

	t.Run("Mixed", func(t *testing.T) {
		r := newReader(encodedMixed)
		require.Equal(t, mixedExpected, readAll(t, r))
	})

	t.Run("Fork_AtStart", func(t *testing.T) {
		r1 := newReader(encodedMixed)
		r2, err := r1.Fork()
		require.NoError(t, err)
		// both should yield the full sequence independently
		got1 := readAll(t, r1)
		got2 := readAll(t, r2)
		require.Equal(t, mixedExpected, got1)
		require.Equal(t, mixedExpected, got2)
	})

	t.Run("Fork_MidStream", func(t *testing.T) {
		if len(mixedExpected) < 2 {
			t.Skip("need at least 2 elements")
		}
		r1 := newReader(encodedMixed)

		// advance r1 by one element
		first, err := r1.Next()
		require.NoError(t, err)

		// fork at current position
		r2, err := r1.Fork()
		require.NoError(t, err)

		// r1 and r2 should both yield the remaining elements
		tail := mixedExpected[1:]
		got1 := readAll(t, r1)
		got2 := readAll(t, r2)
		require.Equal(t, tail, got1)
		require.Equal(t, tail, got2)

		// verify the first element was correct
		v, valid := first.Value()
		require.Equal(t, mixedExpected[0].val, v)
		require.Equal(t, !mixedExpected[0].null, valid)
	})

	t.Run("Fork_Independent", func(t *testing.T) {
		if len(mixedExpected) < 2 {
			t.Skip("need at least 2 elements")
		}
		r1 := newReader(encodedMixed)
		r2, err := r1.Fork()
		require.NoError(t, err)

		// drain r1 completely
		readAll(t, r1)

		// r2 must be unaffected
		require.Equal(t, mixedExpected, readAll(t, r2))
	})

	t.Run("Fork_AfterEOF", func(t *testing.T) {
		r1 := newReader(encodedMixed)
		readAll(t, r1)

		r2, err := r1.Fork()
		require.NoError(t, err)
		_, err = r2.Next()
		require.Equal(t, io.EOF, err)
	})

	t.Run("Truncated_ReturnsError", func(t *testing.T) {
		// A single byte "s64(3)" claims a 3x repeated run, but no value follows.
		r := newReader([]byte{0x03})
		_, err := r.Next()
		require.Error(t, err)
	})
}

// --- uint64 ---

// Uint64 encoding helpers (all LEB128):
//   s64(3)=0x03  u64(42)=0x2a  s64(0)=0x00  u64(2)=0x02  s64(-3)=0x7d

func newUint64Reader(data []byte) *Reader[uint64] {
	return NewUint64Reader(ioutil.NewSubReader(data))
}

func TestUint64Reader(t *testing.T) {
	suiteRLEReader(t, newUint64Reader,
		// empty
		[]byte{},
		// repeated: 3× 42
		[]byte{0x03, 0x2a}, []struct {
			val  uint64
			null bool
		}{value[uint64](42), value[uint64](42), value[uint64](42)},
		// single (L=1): 1× 7
		[]byte{0x01, 0x07}, []struct {
			val  uint64
			null bool
		}{value[uint64](7)},
		// null run: 2 nulls
		[]byte{0x00, 0x02}, []struct {
			val  uint64
			null bool
		}{null[uint64](), null[uint64]()},
		// literal: [1, 2, 3]
		[]byte{0x7d, 0x01, 0x02, 0x03}, []struct {
			val  uint64
			null bool
		}{value[uint64](1), value[uint64](2), value[uint64](3)},
		// mixed: 3×0, 2 nulls, [1,2,3]
		[]byte{0x03, 0x00, 0x00, 0x02, 0x7d, 0x01, 0x02, 0x03}, []struct {
			val  uint64
			null bool
		}{value[uint64](0), value[uint64](0), value[uint64](0), null[uint64](), null[uint64](), value[uint64](1), value[uint64](2), value[uint64](3)},
	)
}

func BenchmarkUint64Reader(b *testing.B) {
	buf := []byte{0x03, 0x00, 0x00, 0x02, 0x7d, 0x01, 0x02, 0x03}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NewUint64Reader(ioutil.NewSubReader(buf))
		for {
			if _, err := r.Next(); err != nil {
				break
			}
		}
	}
}

// --- int64 ---

// Int64 encoding helpers:
//   s64(3)=0x03  s64(5)=0x05  s64(1)=0x01  s64(-1)=0x7f
//   s64(0)=0x00  u64(2)=0x02  s64(-3)=0x7d  s64(-2)=0x7e

func newInt64Reader(data []byte) *Reader[int64] {
	return NewInt64Reader(ioutil.NewSubReader(data))
}

func TestInt64Reader(t *testing.T) {
	suiteRLEReader(t, newInt64Reader,
		// empty
		[]byte{},
		// repeated: 3× 5
		[]byte{0x03, 0x05}, []struct {
			val  int64
			null bool
		}{value[int64](5), value[int64](5), value[int64](5)},
		// single (L=1): 1× -1
		[]byte{0x01, 0x7f}, []struct {
			val  int64
			null bool
		}{value[int64](-1)},
		// null run: 2 nulls
		[]byte{0x00, 0x02}, []struct {
			val  int64
			null bool
		}{null[int64](), null[int64]()},
		// literal: [3, -2, 1]
		[]byte{0x7d, 0x03, 0x7e, 0x01}, []struct {
			val  int64
			null bool
		}{value[int64](3), value[int64](-2), value[int64](1)},
		// mixed: literal [3], 3×1, literal [3,-2,1]
		[]byte{0x7f, 0x03, 0x03, 0x01, 0x7d, 0x03, 0x7e, 0x01}, []struct {
			val  int64
			null bool
		}{value[int64](3), value[int64](1), value[int64](1), value[int64](1), value[int64](3), value[int64](-2), value[int64](1)},
	)
}

func BenchmarkInt64Reader(b *testing.B) {
	buf := []byte{0x7f, 0x03, 0x03, 0x01, 0x7d, 0x03, 0x7e, 0x01}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NewInt64Reader(ioutil.NewSubReader(buf))
		for {
			if _, err := r.Next(); err != nil {
				break
			}
		}
	}
}

// --- string ---

// String encoding: s64(L), then for each value: u64(len), bytes.
//   repeated 2× "hi": [0x02, 0x02, 'h', 'i']
//   single "x":        [0x01, 0x01, 'x']
//   null run 3:        [0x00, 0x03]
//   literal ["a","bb"]: [0x7e, 0x01,'a', 0x02,'b','b']

func newStringReader(data []byte) *Reader[string] {
	return NewStringReader(ioutil.NewSubReader(data))
}

func TestStringReader(t *testing.T) {
	suiteRLEReader(t, newStringReader,
		// empty
		[]byte{},
		// repeated: 2× "hi"
		[]byte{0x02, 0x02, 'h', 'i'}, []struct {
			val  string
			null bool
		}{value("hi"), value("hi")},
		// single (L=1): "x"
		[]byte{0x01, 0x01, 'x'}, []struct {
			val  string
			null bool
		}{value("x")},
		// null run: 3 nulls
		[]byte{0x00, 0x03}, []struct {
			val  string
			null bool
		}{null[string](), null[string](), null[string]()},
		// literal: ["a", "", "boo"]
		[]byte{0x7d, 0x01, 'a', 0x00, 0x03, 'b', 'o', 'o'}, []struct {
			val  string
			null bool
		}{value("a"), value(""), value("boo")},
		// mixed: literal ["a",""], null×1, repeated 2× "boo"
		[]byte{0x7e, 0x01, 'a', 0x00, 0x00, 0x01, 0x02, 0x03, 'b', 'o', 'o'}, []struct {
			val  string
			null bool
		}{value("a"), value(""), null[string](), value("boo"), value("boo")},
	)
}

func BenchmarkStringReader(b *testing.B) {
	buf := []byte{0x7e, 0x01, 'a', 0x00, 0x00, 0x01, 0x02, 0x03, 'b', 'o', 'o'}
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		r := NewStringReader(ioutil.NewSubReader(buf))
		for {
			if _, err := r.Next(); err != nil {
				break
			}
		}
	}
}
