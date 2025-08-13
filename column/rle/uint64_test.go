package rle

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"gotomerge/lbuf"
)

func TestReadUint64RLE(t *testing.T) {
	buf := []byte{0x03, 0x00, 0x00, 0x02, 0x7d, 0x01, 0x02, 0x03}

	type tuple struct {
		val   uint64
		valid bool
	}

	expected := []tuple{
		{0, true},
		{0, true},
		{0, true},
		{0, false},
		{0, false},
		{1, true},
		{2, true},
		{3, true},
	}

	var res []tuple

	for uleb, err := range ReadUint64RLE(lbuf.FromBytes(buf)) {
		require.NoError(t, err)
		val, valid := uleb.Value()
		res = append(res, tuple{val: val, valid: valid})
	}

	require.Equal(t, expected, res)
}

func BenchmarkReadUint64RLE(b *testing.B) {
	b.ReportAllocs()

	// expand to 8 values
	buf := []byte{0x03, 0x00, 0x00, 0x02, 0x7d, 0x01, 0x02, 0x03}
	for i := 0; i < b.N; i++ {
		for _, _ = range ReadUint64RLE(lbuf.FromBytes(buf)) {
		}
	}
}

func TestName(t *testing.T) {
	var a uint32 = math.MaxUint32
	var b uint32 = 1
	fmt.Println(a + b)
}
