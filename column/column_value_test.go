package column

import (
	"bytes"
	"io"
	"testing"

	"gotomerge/types"
	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

func TestValueMetadataRoundTrip(t *testing.T) {
	in := []ValueMetadata{
		NewValueMetadata(ValueTypeNull, 0),
		NewValueMetadata(ValueTypeTrue, 0),
		NewValueMetadata(ValueTypeFalse, 0),
		NewValueMetadata(ValueTypeString, 5),
		NewValueMetadata(ValueTypeUInt, 2),
		NewValueMetadata(ValueTypeInt, 1),
		NewValueMetadata(ValueTypeNull, 0),
	}
	var buf bytes.Buffer
	w := NewValueMetadataWriter(&buf)
	for _, m := range in {
		w.Append(m)
	}
	require.NoError(t, w.Flush())

	r := NewValueMetadataReader(ioutil.NewSubReader(buf.Bytes()))
	for i, want := range in {
		got, err := r.Next()
		require.NoError(t, err, "index %d", i)
		require.Equal(t, want, got, "index %d", i)
	}
	_, err := r.Next()
	require.ErrorIs(t, err, io.EOF)
}

func TestValueRoundTrip(t *testing.T) {
	in := []types.Action{
		{Kind: types.ActionSet, Value: nil},
		{Kind: types.ActionSet, Value: true},
		{Kind: types.ActionSet, Value: false},
		{Kind: types.ActionSet, Value: "hello"},
		{Kind: types.ActionSet, Value: uint64(42)},
		{Kind: types.ActionSet, Value: int64(-7)},
		{Kind: types.ActionMakeMap},
		{Kind: types.ActionDelete},
		{Kind: types.ActionSet, Value: ""},
	}

	var metaBuf, valBuf bytes.Buffer
	w := NewValueWriter(&metaBuf, &valBuf)
	for _, a := range in {
		w.Append(a)
	}
	require.NoError(t, w.Flush())

	mr := NewValueMetadataReader(ioutil.NewSubReader(metaBuf.Bytes()))
	vr := NewValueReader(ioutil.NewSubReader(valBuf.Bytes()))

	for i, want := range in {
		meta, err := mr.Next()
		require.NoError(t, err, "index %d meta", i)
		got, err := vr.Next(meta)
		require.NoError(t, err, "index %d value", i)
		if HasScalarValue(want) {
			require.Equal(t, want.Value, got, "index %d", i)
		} else {
			require.Nil(t, got, "index %d", i)
		}
	}
}
