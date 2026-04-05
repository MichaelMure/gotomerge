package column

import (
	"bytes"
	"testing"

	"gotomerge/types"
	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

func TestActionHasValues(t *testing.T) {
	newWriter := func() *ActionWriter {
		var k, m, v bytes.Buffer
		return NewActionWriter(&k, &m, &v)
	}

	t.Run("false when no entries", func(t *testing.T) {
		require.False(t, newWriter().HasValues())
	})
	t.Run("false for non-scalar actions only", func(t *testing.T) {
		w := newWriter()
		w.Append(types.Action{Kind: types.ActionMakeMap})
		w.Append(types.Action{Kind: types.ActionDelete})
		w.Append(types.Action{Kind: types.ActionMakeList})
		require.False(t, w.HasValues())
	})
	t.Run("true for Set", func(t *testing.T) {
		w := newWriter()
		w.Append(types.Action{Kind: types.ActionSet, Value: "x"})
		require.True(t, w.HasValues())
	})
	t.Run("true for Inc", func(t *testing.T) {
		w := newWriter()
		w.Append(types.Action{Kind: types.ActionInc, Value: int64(1)})
		require.True(t, w.HasValues())
	})
	t.Run("true when mixed with non-scalar", func(t *testing.T) {
		w := newWriter()
		w.Append(types.Action{Kind: types.ActionDelete})
		w.Append(types.Action{Kind: types.ActionSet, Value: nil})
		require.True(t, w.HasValues())
	})
}

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
			NewUlebReader(ioutil.NewSubReader(kindBuf.Bytes())),
			NewValueMetadataReader(ioutil.NewSubReader(metaBuf.Bytes())),
			NewValueReader(ioutil.NewSubReader(valBuf.Bytes())),
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
