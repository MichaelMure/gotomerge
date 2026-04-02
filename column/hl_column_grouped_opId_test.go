package column

import (
	"bytes"
	"testing"

	ioutil "gotomerge/utils/io"
	"gotomerge/types"

	"github.com/stretchr/testify/require"
)

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
			NewGroupReader(ioutil.NewBytesReader(grpBuf.Bytes())),
			NewActorReader(ioutil.NewBytesReader(actorBuf.Bytes())),
			NewDeltaReader(ioutil.NewBytesReader(ctrBuf.Bytes())),
		)
		for i, want := range in {
			got, err := r.Next()
			require.NoError(t, err, "index %d", i)
			require.Equal(t, want, got, "index %d", i)
		}
	}
}
