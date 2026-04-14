package column

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/gotomerge/types"
)

func TestActionRoundTrip(t *testing.T) {
	cases := [][]types.Action{
		{
			{Kind: types.ActionMakeMap},
			{Kind: types.ActionSet, Value: "hello"},
			{Kind: types.ActionSet, Value: uint64(42)},
			{Kind: types.ActionDelete},
			{Kind: types.ActionMakeList},
			{Kind: types.ActionSet, Value: nil},
			{Kind: types.ActionSet, Value: true},
			{Kind: types.ActionSet, Value: int64(-3)},
		},
		repeat(100, types.Action{Kind: types.ActionDelete}),
		repeat(50, types.Action{Kind: types.ActionSet, Value: "x"}),
		repeat(200, types.Action{Kind: types.ActionMakeMap}, types.Action{Kind: types.ActionSet, Value: nil}),
	}
	for _, in := range cases {
		var kindBuf, metaBuf, valBuf bytes.Buffer
		w := NewActionWriter(&kindBuf, &metaBuf, &valBuf)
		for _, a := range in {
			w.Append(a)
		}
		require.NoError(t, w.Flush())

		r := NewActionReader(
			bytesOpt(kindBuf.Bytes(), PeekUlebReader),
			bytesOpt(metaBuf.Bytes(), PeekValueMetadataReader),
			bytesOpt(valBuf.Bytes(), PeekValueReader),
		)
		var out []types.Action
		for {
			a, err := r.Next()
			if err == ErrDone {
				break
			}
			require.NoError(t, err)
			out = append(out, a)
		}
		require.Equal(t, in, out)
	}
}

// TODO: add round-trip fuzzing for all columns
