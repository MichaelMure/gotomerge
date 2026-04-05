package column

import (
	"bytes"
	"testing"

	"gotomerge/types"
	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

func TestKeyHasFlags(t *testing.T) {
	newWriter := func() (*KeyWriter, bytes.Buffer, bytes.Buffer, bytes.Buffer) {
		var a, c, s bytes.Buffer
		return NewKeyWriter(&a, &c, &s), a, c, s
	}

	t.Run("all false when no entries", func(t *testing.T) {
		w, _, _, _ := newWriter()
		require.False(t, w.HasOpId())
		require.False(t, w.HasString())
		require.False(t, w.HasNonNullActor())
	})
	t.Run("HasString only for string keys", func(t *testing.T) {
		w, _, _, _ := newWriter()
		w.Append(types.KeyString("x"), identityLocalOf)
		require.False(t, w.HasOpId())
		require.True(t, w.HasString())
		require.False(t, w.HasNonNullActor())
	})
	t.Run("HasOpId without HasNonNullActor for head sentinel (counter=0)", func(t *testing.T) {
		w, _, _, _ := newWriter()
		w.Append(types.KeyOpId{ActorIdx: 0, Counter: 0}, identityLocalOf)
		require.True(t, w.HasOpId())
		require.False(t, w.HasString())
		require.False(t, w.HasNonNullActor())
	})
	t.Run("HasOpId and HasNonNullActor for opId with non-zero counter", func(t *testing.T) {
		w, _, _, _ := newWriter()
		w.Append(types.KeyOpId{ActorIdx: 1, Counter: 5}, identityLocalOf)
		require.True(t, w.HasOpId())
		require.False(t, w.HasString())
		require.True(t, w.HasNonNullActor())
	})
	t.Run("all true for mixed keys", func(t *testing.T) {
		w, _, _, _ := newWriter()
		w.Append(types.KeyString("a"), identityLocalOf)
		w.Append(types.KeyOpId{ActorIdx: 0, Counter: 3}, identityLocalOf)
		require.True(t, w.HasOpId())
		require.True(t, w.HasString())
		require.True(t, w.HasNonNullActor())
	})
}

func TestKeyRoundTrip(t *testing.T) {
	cases := [][]types.Key{
		{
			types.KeyString("foo"), types.KeyString("bar"),
			types.KeyOpId{ActorIdx: 0, Counter: 3}, types.KeyOpId{ActorIdx: 1, Counter: 7},
			types.KeyString("baz"),
		},
		repeat[types.Key](100, types.KeyString("name")),
		repeat[types.Key](50, types.KeyOpId{ActorIdx: 0, Counter: 5}),
		repeat[types.Key](200, types.KeyString("x"), types.KeyOpId{ActorIdx: 0, Counter: 1}),
	}
	for _, in := range cases {
		var actorBuf, ctrBuf, strBuf bytes.Buffer
		w := NewKeyWriter(&actorBuf, &ctrBuf, &strBuf)
		for _, k := range in {
			w.Append(k, identityLocalOf)
		}
		require.NoError(t, w.Flush())

		r := NewKeyReader(
			NewActorReader(ioutil.NewSubReader(actorBuf.Bytes())),
			NewDeltaReader(ioutil.NewSubReader(ctrBuf.Bytes())),
			NewStringReader(ioutil.NewSubReader(strBuf.Bytes())),
		)
		for i, want := range in {
			got, err := r.Next()
			require.NoError(t, err, "index %d", i)
			require.Equal(t, want, got, "index %d", i)
		}
	}
}
