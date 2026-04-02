package column

import (
	"bytes"
	"testing"

	ioutil "gotomerge/utils/io"
	"gotomerge/types"

	"github.com/stretchr/testify/require"
)

func TestOpIdRoundTrip(t *testing.T) {
	cases := [][]types.OpId{
		{{ActorIdx: 0, Counter: 1}, {ActorIdx: 0, Counter: 3}, {ActorIdx: 1, Counter: 5}, {ActorIdx: 0, Counter: 10}},
		repeat(50, types.OpId{ActorIdx: 0, Counter: 7}),
		repeat(100, types.OpId{ActorIdx: 0, Counter: 1}, types.OpId{ActorIdx: 1, Counter: 1}),
	}
	for _, in := range cases {
		var actorBuf, ctrBuf bytes.Buffer
		w := NewOpIdWriter(&actorBuf, &ctrBuf)
		for _, id := range in {
			w.Append(id, identityLocalOf)
		}
		require.NoError(t, w.Flush())

		r := NewOpIdReader(
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
