package column

import (
	"bytes"
	"io"
	"testing"

	"gotomerge/column/rle"
	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

func TestReadDeltaColumn(t *testing.T) {
	buf := []byte{0x7f, 0x03, 0x03, 0x01, 0x7d, 0x03, 0x7e, 0x01}

	expected := []any{int64(3), int64(4), int64(5), int64(6), int64(9), int64(7), int64(8)}
	var res []any

	r := NewDeltaReader(ioutil.NewBytesReader(buf))
	for {
		nv, err := r.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		if val, valid := nv.Value(); valid {
			res = append(res, val)
		} else {
			res = append(res, nil)
		}
	}

	require.Equal(t, expected, res)
}

func TestDeltaRoundTrip(t *testing.T) {
	cases := []struct {
		name string
		in   []rle.NullableValue[int64]
	}{
		{"ascending", []rle.NullableValue[int64]{
			rle.NewNullableInt64(3), rle.NewNullableInt64(7), rle.NewNullableInt64(7), rle.NewNullableInt64(10),
		}},
		{"with nulls", []rle.NullableValue[int64]{
			rle.NewNullableInt64(5), rle.NewNullInt64(), rle.NewNullableInt64(8),
		}},
		{"single", []rle.NullableValue[int64]{rle.NewNullableInt64(42)}},
		{"descending", []rle.NullableValue[int64]{
			rle.NewNullableInt64(100), rle.NewNullableInt64(50), rle.NewNullableInt64(0),
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			w := NewDeltaWriter(&buf)
			for _, nv := range tc.in {
				w.Append(nv)
			}
			require.NoError(t, w.Flush())

			r := NewDeltaReader(ioutil.NewBytesReader(buf.Bytes()))
			for i, want := range tc.in {
				got, err := r.Next()
				require.NoError(t, err, "index %d", i)
				wantV, wantOk := want.Value()
				gotV, gotOk := got.Value()
				require.Equal(t, wantOk, gotOk, "index %d null", i)
				if wantOk {
					require.Equal(t, wantV, gotV, "index %d value", i)
				}
			}
			_, err := r.Next()
			require.ErrorIs(t, err, io.EOF)
		})
	}
}
