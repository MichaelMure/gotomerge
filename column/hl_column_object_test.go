package column

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/MichaelMure/gotomerge/types"
)

// identityLocalOf maps each global actor index to itself, used across HL column round-trip tests.
var identityLocalOf = types.IdentityActorMapper(4)

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
			bytesOpt(actorBuf.Bytes(), NewActorReader),
			bytesOpt(ctrBuf.Bytes(), NewUlebReader),
		)
		for i, want := range in {
			got, err := r.Next()
			require.NoError(t, err, "index %d", i)
			require.Equal(t, want, got, "index %d", i)
		}
	}
}
