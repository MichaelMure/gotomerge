package column

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/gotomerge/types"
)

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
			bytesOpt(actorBuf.Bytes(), PeekActorReader),
			bytesOpt(ctrBuf.Bytes(), PeekDeltaReader),
			bytesOpt(strBuf.Bytes(), PeekStringReader),
		)
		for i, want := range in {
			got, err := r.Next()
			require.NoError(t, err, "index %d", i)
			require.Equal(t, want, got, "index %d", i)
		}
	}
}
