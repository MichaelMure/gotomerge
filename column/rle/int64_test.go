package rle

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gotomerge/lbuf"
)

func TestReadInt64RLE(t *testing.T) {
	buf := []byte{0x7f, 0x03, 0x03, 0x01, 0x7d, 0x03, 0x7e, 0x01}

	type tuple struct {
		val   int64
		valid bool
	}

	expected := []tuple{
		{3, true},
		{1, true},
		{1, true},
		{1, true},
		{3, true},
		{-2, true},
		{1, true},
	}

	var res []tuple

	for leb, err := range ReadInt64RLE(lbuf.FromBytes(buf)) {
		require.NoError(t, err)
		val, valid := leb.Value()
		res = append(res, tuple{val: val, valid: valid})
	}

	require.Equal(t, expected, res)
}

func BenchmarkReadInt64RLE(b *testing.B) {
	b.ReportAllocs()

	// expand to 7 values
	buf := []byte{0x7f, 0x03, 0x03, 0x01, 0x7d, 0x03, 0x7e, 0x01}
	for i := 0; i < b.N; i++ {
		for _, _ = range ReadInt64RLE(lbuf.FromBytes(buf)) {
		}
	}
}
