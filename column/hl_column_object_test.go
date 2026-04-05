package column

import (
	"bytes"
	"testing"

	"gotomerge/types"
	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

// identityLocalOf maps each global actor index to itself, used across HL column round-trip tests.
var identityLocalOf = map[uint32]uint32{0: 0, 1: 1, 2: 2, 3: 3}

func TestObjectHasNonRoot(t *testing.T) {
	t.Run("false when no entries", func(t *testing.T) {
		var a, c bytes.Buffer
		w := NewObjectWriter(&a, &c)
		require.False(t, w.HasNonRoot())
	})
	t.Run("false when all root", func(t *testing.T) {
		var a, c bytes.Buffer
		w := NewObjectWriter(&a, &c)
		for range 5 {
			w.Append(types.RootObjectId(), identityLocalOf)
		}
		require.False(t, w.HasNonRoot())
	})
	t.Run("true when any non-root", func(t *testing.T) {
		var a, c bytes.Buffer
		w := NewObjectWriter(&a, &c)
		w.Append(types.RootObjectId(), identityLocalOf)
		w.Append(types.ObjectId{ActorIdx: 1, Counter: 3}, identityLocalOf)
		require.True(t, w.HasNonRoot())
	})
}

func TestObjectRoundTrip(t *testing.T) {
	cases := [][]types.ObjectId{
		{types.RootObjectId(), {ActorIdx: 1, Counter: 5}, {ActorIdx: 2, Counter: 10}, types.RootObjectId()},
		repeat(100, types.RootObjectId()),
		repeat(50, types.ObjectId{ActorIdx: 1, Counter: 3}),
		repeat(200, types.RootObjectId(), types.ObjectId{ActorIdx: 1, Counter: 1}),
	}
	for _, in := range cases {
		var actorBuf, ctrBuf bytes.Buffer
		w := NewObjectWriter(&actorBuf, &ctrBuf)
		for _, obj := range in {
			w.Append(obj, identityLocalOf)
		}
		require.NoError(t, w.Flush())

		r := NewObjectReader(
			NewActorReader(ioutil.NewSubReader(actorBuf.Bytes())),
			NewUlebReader(ioutil.NewSubReader(ctrBuf.Bytes())),
		)
		for i, want := range in {
			got, err := r.Next()
			require.NoError(t, err, "index %d", i)
			require.Equal(t, want, got, "index %d", i)
		}
	}
}
