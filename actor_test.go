package gotomerge

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActorIdsRoundtrip(t *testing.T) {
	data := []byte{0x01, 0x03, 0xab, 0xcd, 0xef}
	ids, err := readActorIds(bytes.NewReader(data))
	require.NoError(t, err)

	require.Equal(t, 1, len(ids))
	require.Equal(t, ActorId{0xab, 0xcd, 0xef}, ids[0])

	buf := new(bytes.Buffer)
	err = writeActorIds(buf, ids)
	require.NoError(t, err)
	require.Equal(t, data, buf.Bytes())
}
