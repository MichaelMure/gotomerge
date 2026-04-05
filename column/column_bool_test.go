package column

import (
	"bytes"
	"io"
	"testing"

	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

func TestReadBooleanColumn(t *testing.T) {
	buf := []byte{0x00, 0x02, 0x03}

	expected := []bool{true, true, false, false, false}
	var res []bool

	r := NewBoolReader(ioutil.NewSubReader(buf))
	for {
		b, err := r.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		res = append(res, b)
	}

	require.Equal(t, expected, res)
}

func repeat[T any](n int, vals ...T) []T {
	res := make([]T, n*len(vals))
	for i := 0; i < n; i++ {
		for j := range vals {
			res[i*len(vals)+j] = vals[j]
		}
	}
	return res
}

func TestBoolRoundTrip(t *testing.T) {
	cases := [][]bool{
		{false, true, true, false, false, false, true},
		{true, true, true},
		{false, false, false},
		{true},
		{false},
		{true, false, true, false},
		repeat(1000, true),
		repeat(1000, false),
		repeat(1000, true, false),
	}
	for _, in := range cases {
		var buf bytes.Buffer
		w := NewBoolWriter(&buf)
		for _, v := range in {
			w.Append(v)
		}
		require.NoError(t, w.Flush())

		r := NewBoolReader(ioutil.NewSubReader(buf.Bytes()))
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
