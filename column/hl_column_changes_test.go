package column

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

// bytesOpt creates a typed reader from a byte slice.
// Returns nil when data is empty — an all-null column produces zero bytes
// (matching Rust's InitialNullRun behaviour), and nil readers are handled as
// all-null by the high-level column readers.
func bytesOpt[T any](data []byte, ctor func(ioutil.SubReader) *T) *T {
	if len(data) == 0 {
		return nil
	}
	return ctor(*ioutil.NewSubReader(data))
}

func TestChangesRoundTrip(t *testing.T) {
	t0 := int64(1000)
	msg0 := "first change"
	cases := [][]RawChangeMeta{
		{
			{ActorIdx: 0, SeqNum: 1, MaxOp: 5, Time: &t0, Message: &msg0, Deps: []uint64{}},
			{ActorIdx: 1, SeqNum: 1, MaxOp: 10, Time: nil, Message: nil, Deps: []uint64{0}},
			{ActorIdx: 0, SeqNum: 2, MaxOp: 15, Time: nil, Message: nil, Deps: []uint64{0, 1}},
		},
		// same metadata repeated: stresses delta repeat-run encoding
		repeat(50, RawChangeMeta{ActorIdx: 0, SeqNum: 1, MaxOp: 5, Deps: []uint64{}}),
		// alternating actors
		repeat(100,
			RawChangeMeta{ActorIdx: 0, SeqNum: 1, MaxOp: 1, Deps: []uint64{}},
			RawChangeMeta{ActorIdx: 1, SeqNum: 1, MaxOp: 1, Deps: []uint64{}},
		),
	}
	for _, in := range cases {
		var actorBuf, seqBuf, maxOpBuf, timeBuf, msgBuf, grpBuf, idxBuf bytes.Buffer
		w := NewChangesWriter(&actorBuf, &seqBuf, &maxOpBuf, &timeBuf, &msgBuf, &grpBuf, &idxBuf)
		for _, m := range in {
			w.Append(m)
		}
		require.NoError(t, w.Flush())

		r := NewChangesReader(
			bytesOpt(actorBuf.Bytes(), PeekActorReader),
			bytesOpt(seqBuf.Bytes(), PeekDeltaReader),
			bytesOpt(maxOpBuf.Bytes(), PeekDeltaReader),
			bytesOpt(timeBuf.Bytes(), PeekDeltaReader),
			bytesOpt(msgBuf.Bytes(), PeekStringReader),
			bytesOpt(grpBuf.Bytes(), PeekGroupReader),
			bytesOpt(idxBuf.Bytes(), PeekDeltaReader),
		)
		for i, want := range in {
			got, err := r.Next()
			require.NoError(t, err, "index %d", i)
			require.Equal(t, want, got, "index %d", i)
		}
	}
}
