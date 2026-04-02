package column

import (
	"bytes"
	"testing"

	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

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
			NewActorReader(ioutil.NewBytesReader(actorBuf.Bytes())),
			NewDeltaReader(ioutil.NewBytesReader(seqBuf.Bytes())),
			NewDeltaReader(ioutil.NewBytesReader(maxOpBuf.Bytes())),
			NewDeltaReader(ioutil.NewBytesReader(timeBuf.Bytes())),
			NewStringReader(ioutil.NewBytesReader(msgBuf.Bytes())),
			NewGroupReader(ioutil.NewBytesReader(grpBuf.Bytes())),
			NewDeltaReader(ioutil.NewBytesReader(idxBuf.Bytes())),
		)
		for i, want := range in {
			got, err := r.Next()
			require.NoError(t, err, "index %d", i)
			require.Equal(t, want, got, "index %d", i)
		}
	}
}
