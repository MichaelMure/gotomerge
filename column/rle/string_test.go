package rle_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gotomerge/column/rle"
	"gotomerge/lbuf"
)

func TestReadRleString(t *testing.T) {
	buf := []byte{0x7e, 0x01, 0x61, 0x00, 0x00, 0x01, 0x02, 0x03, 0x62, 0x6f, 0x6f}

	type tuple struct {
		val   string
		valid bool
	}

	expected := []tuple{
		{"a", true},
		{"", true},
		{"", false},
		{"boo", true},
		{"boo", true},
	}

	var res []tuple

	for str, err := range rle.ReadStringRLE(lbuf.FromBytes(buf)) {
		require.NoError(t, err)
		val, valid := str.Value()
		res = append(res, tuple{val: val, valid: valid})
	}

	require.Equal(t, expected, res)
}

func BenchmarkReadRleString(b *testing.B) {
	b.ReportAllocs()

	// expand to 5 values
	buf := []byte{0x7e, 0x01, 0x61, 0x00, 0x00, 0x01, 0x02, 0x03, 0x62, 0x6f, 0x6f}
	for i := 0; i < b.N; i++ {
		for _, _ = range rle.ReadStringRLE(lbuf.FromBytes(buf)) {
		}
	}
}
