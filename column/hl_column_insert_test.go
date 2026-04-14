package column

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

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

		r := NewInsertReader(
			bytesOpt(buf.Bytes(), PeekBoolReader),
		)
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
