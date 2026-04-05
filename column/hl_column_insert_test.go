package column

import (
	"bytes"
	"io"
	"testing"

	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

func TestInsertHasInserts(t *testing.T) {
	t.Run("false when no entries", func(t *testing.T) {
		var buf bytes.Buffer
		w := NewInsertWriter(&buf)
		require.False(t, w.HasInserts())
	})
	t.Run("false when all false", func(t *testing.T) {
		var buf bytes.Buffer
		w := NewInsertWriter(&buf)
		for _, v := range repeat(10, false) {
			w.Append(v)
		}
		require.False(t, w.HasInserts())
	})
	t.Run("true when any true", func(t *testing.T) {
		var buf bytes.Buffer
		w := NewInsertWriter(&buf)
		w.Append(false)
		w.Append(true)
		require.True(t, w.HasInserts())
	})
}

func TestInsertRoundTrip(t *testing.T) {
	cases := [][]bool{
		{false, true, false, false, true, true, false},
		{true},
		{false},
		repeat(200, true),
		repeat(200, false),
		repeat(500, true, false),
	}
	for _, in := range cases {
		var buf bytes.Buffer
		w := NewInsertWriter(&buf)
		for _, v := range in {
			w.Append(v)
		}
		require.NoError(t, w.Flush())

		r := NewInsertReader(NewBoolReader(ioutil.NewSubReader(buf.Bytes())))
		var out []bool
		for {
			v, err := r.Next()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
			out = append(out, v)
		}
		require.Equal(t, in, out)
	}
}
