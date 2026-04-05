package rle

import (
	"bytes"
	"testing"

	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

// roundTripUint64 encodes values with the Writer and decodes with the Reader,
// checking that the round-trip produces the same sequence.
func roundTripUint64(t *testing.T, vals []NullableValue[uint64]) {
	t.Helper()
	var buf bytes.Buffer
	w := NewUint64Writer(&buf)
	for _, nv := range vals {
		w.Append(nv)
	}
	require.NoError(t, w.Flush())
	got := readAll(t, NewUint64Reader(ioutil.NewSubReader(buf.Bytes())))
	require.Len(t, got, len(vals))
	for i, nv := range vals {
		v, ok := nv.Value()
		require.Equal(t, !ok, got[i].null, "index %d null mismatch", i)
		if ok {
			require.Equal(t, v, got[i].val, "index %d value mismatch", i)
		}
	}
}

func TestWriterEmpty(t *testing.T) {
	var buf bytes.Buffer
	w := NewUint64Writer(&buf)
	require.NoError(t, w.Flush())
	require.Empty(t, buf.Bytes())
}

func TestWriterSingleValue(t *testing.T) {
	roundTripUint64(t, []NullableValue[uint64]{NewNullableUint64(42)})
}

func TestWriterRepeatRun(t *testing.T) {
	// 5× the same value should produce a repeat run, not a literal run.
	vals := []NullableValue[uint64]{
		NewNullableUint64(7), NewNullableUint64(7), NewNullableUint64(7),
		NewNullableUint64(7), NewNullableUint64(7),
	}
	var buf bytes.Buffer
	w := NewUint64Writer(&buf)
	for _, nv := range vals {
		w.Append(nv)
	}
	require.NoError(t, w.Flush())
	b := buf.Bytes()
	// Repeat run encoding: S64(5) + U64(7) = 2 LEB128 bytes
	// Literal run encoding: S64(-5) + 5×U64(7) = 6 bytes
	// Expect the shorter repeat encoding.
	require.Equal(t, []byte{0x05, 0x07}, b)
	roundTripUint64(t, vals)
}

func TestWriterLiteralRun(t *testing.T) {
	roundTripUint64(t, []NullableValue[uint64]{
		NewNullableUint64(1), NewNullableUint64(2), NewNullableUint64(3),
	})
}

func TestWriterNullRun(t *testing.T) {
	roundTripUint64(t, []NullableValue[uint64]{
		NewNullUint64(), NewNullUint64(), NewNullUint64(),
	})
}

func TestWriterRepeatSplitsFromLiteral(t *testing.T) {
	// [1, 2, 2, 2] — the 2,2 at the end should become a repeat run, not stay literal.
	vals := []NullableValue[uint64]{
		NewNullableUint64(1), NewNullableUint64(2), NewNullableUint64(2), NewNullableUint64(2),
	}
	var buf bytes.Buffer
	w := NewUint64Writer(&buf)
	for _, nv := range vals {
		w.Append(nv)
	}
	require.NoError(t, w.Flush())
	b := buf.Bytes()
	// Expect: literal([1]) + repeat(3×2)
	// S64(-1)=0x7f, U64(1)=0x01, S64(3)=0x03, U64(2)=0x02
	require.Equal(t, []byte{0x7f, 0x01, 0x03, 0x02}, b)
	roundTripUint64(t, vals)
}

func TestWriterNullThenValue(t *testing.T) {
	roundTripUint64(t, []NullableValue[uint64]{
		NewNullUint64(), NewNullUint64(), NewNullableUint64(5), NewNullableUint64(5),
	})
}

func TestWriterValueThenNull(t *testing.T) {
	roundTripUint64(t, []NullableValue[uint64]{
		NewNullableUint64(1), NewNullableUint64(1), NewNullUint64(), NewNullUint64(),
	})
}

func TestWriterMixed(t *testing.T) {
	// Mirrors the mixed case from TestUint64Reader: 3×0, 2 nulls, [1,2,3]
	roundTripUint64(t, []NullableValue[uint64]{
		NewNullableUint64(0), NewNullableUint64(0), NewNullableUint64(0),
		NewNullUint64(), NewNullUint64(),
		NewNullableUint64(1), NewNullableUint64(2), NewNullableUint64(3),
	})
}

func TestWriterFirstValueNull(t *testing.T) {
	roundTripUint64(t, []NullableValue[uint64]{
		NewNullUint64(), NewNullableUint64(1), NewNullableUint64(2),
	})
}

// TestWriterRepeatThenDifferent checks that two distinct repeat runs are emitted.
func TestWriterRepeatThenDifferent(t *testing.T) {
	vals := []NullableValue[uint64]{
		NewNullableUint64(1), NewNullableUint64(1), NewNullableUint64(1),
		NewNullableUint64(2), NewNullableUint64(2), NewNullableUint64(2),
	}
	var buf bytes.Buffer
	w := NewUint64Writer(&buf)
	for _, nv := range vals {
		w.Append(nv)
	}
	require.NoError(t, w.Flush())
	b := buf.Bytes()
	// S64(3)+U64(1), S64(3)+U64(2)
	require.Equal(t, []byte{0x03, 0x01, 0x03, 0x02}, b)
	roundTripUint64(t, vals)
}

// --- Int64 delta ---

func roundTripInt64Delta(t *testing.T, vals []NullableValue[int64]) {
	t.Helper()
	b := EncodeInt64Delta(vals)
	got := readAll(t, NewInt64Reader(ioutil.NewSubReader(b)))
	// The round-trip goes through the raw delta reader, so we re-accumulate manually.
	require.Len(t, got, len(vals))
	var acc int64
	for i, nv := range vals {
		v, ok := nv.Value()
		require.Equal(t, !ok, got[i].null, "index %d null mismatch", i)
		if ok {
			acc += got[i].val
			require.Equal(t, v, acc, "index %d accumulated value mismatch", i)
		}
	}
}

func TestInt64DeltaRoundTrip(t *testing.T) {
	roundTripInt64Delta(t, []NullableValue[int64]{
		NewNullableInt64(1), NewNullableInt64(3), NewNullableInt64(3),
		NewNullInt64(), NewNullableInt64(10),
	})
}

func TestInt64DeltaRepeated(t *testing.T) {
	// Same counter value repeated → deltas all zero → repeat run of 0.
	roundTripInt64Delta(t, []NullableValue[int64]{
		NewNullableInt64(5), NewNullableInt64(5), NewNullableInt64(5),
	})
}

// --- String ---

func TestStringWriterRoundTrip(t *testing.T) {
	in := []NullableValue[string]{
		NewNullableString("hello"), NewNullableString("hello"),
		NewNullString(),
		NewNullableString("world"),
	}
	var buf bytes.Buffer
	w := NewStringWriter(&buf)
	for _, nv := range in {
		w.Append(nv)
	}
	require.NoError(t, w.Flush())
	got := readAll(t, NewStringReader(ioutil.NewSubReader(buf.Bytes())))
	require.Len(t, got, len(in))
	for i, nv := range in {
		v, ok := nv.Value()
		require.Equal(t, !ok, got[i].null, "index %d", i)
		if ok {
			require.Equal(t, v, got[i].val, "index %d", i)
		}
	}
}

// --- BoolWriter (in column package, tested via EncodeBool) ---

func TestEncodeUint64RoundTrip(t *testing.T) {
	b := EncodeUint64(1, 1, 1, 2, 3)
	got := readAll(t, NewUint64Reader(ioutil.NewSubReader(b)))
	require.Equal(t, []struct {
		val  uint64
		null bool
	}{
		{1, false}, {1, false}, {1, false}, {2, false}, {3, false},
	}, got)
}

func TestEncodeInt64RoundTrip(t *testing.T) {
	b := EncodeInt64(-1, 0, 1, 1, 1)
	got := readAll(t, NewInt64Reader(ioutil.NewSubReader(b)))
	require.Equal(t, []struct {
		val  int64
		null bool
	}{
		{-1, false}, {0, false}, {1, false}, {1, false}, {1, false},
	}, got)
}

func TestEncodeStringRoundTrip(t *testing.T) {
	b := EncodeString("a", "a", "b")
	got := readAll(t, NewStringReader(ioutil.NewSubReader(b)))
	require.Equal(t, []struct {
		val  string
		null bool
	}{
		{"a", false}, {"a", false}, {"b", false},
	}, got)
}
