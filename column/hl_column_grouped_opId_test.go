package column

import (
	"bytes"
	"testing"

	"gotomerge/types"
	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

func TestGroupedOpIdHasPreds(t *testing.T) {
	newWriter := func() *GroupedOpIdWriter {
		var g, a, c bytes.Buffer
		return NewGroupedOpIdWriter(&g, &a, &c)
	}

	t.Run("false when no entries", func(t *testing.T) {
		require.False(t, newWriter().HasPreds())
	})
	t.Run("false when all groups empty", func(t *testing.T) {
		w := newWriter()
		for range 5 {
			w.Append([]types.OpId{}, identityLocalOf)
		}
		require.False(t, w.HasPreds())
	})
	t.Run("true when any group non-empty", func(t *testing.T) {
		w := newWriter()
		w.Append([]types.OpId{}, identityLocalOf)
		w.Append([]types.OpId{{ActorIdx: 0, Counter: 1}}, identityLocalOf)
		require.True(t, w.HasPreds())
	})
}

func TestGroupedOpIdRoundTrip(t *testing.T) {
	cases := [][][]types.OpId{
		{
			{},
			{{ActorIdx: 0, Counter: 1}},
			{{ActorIdx: 0, Counter: 2}, {ActorIdx: 1, Counter: 5}},
			{},
		},
		repeat(100, []types.OpId{}),
		repeat(50, []types.OpId{{ActorIdx: 0, Counter: 1}}),
		repeat(100, []types.OpId{}, []types.OpId{{ActorIdx: 0, Counter: 3}}),
	}
	for _, in := range cases {
		var grpBuf, actorBuf, ctrBuf bytes.Buffer
		w := NewGroupedOpIdWriter(&grpBuf, &actorBuf, &ctrBuf)
		for _, ids := range in {
			w.Append(ids, identityLocalOf)
		}
		require.NoError(t, w.Flush())

		r := NewGroupedOpIdReader("preds",
			NewGroupReader(ioutil.NewSubReader(grpBuf.Bytes())),
			NewActorReader(ioutil.NewSubReader(actorBuf.Bytes())),
			NewDeltaReader(ioutil.NewSubReader(ctrBuf.Bytes())),
		)
		for i, want := range in {
			got, err := r.Next()
			require.NoError(t, err, "index %d", i)
			require.Equal(t, want, got, "index %d", i)
		}
	}
}
