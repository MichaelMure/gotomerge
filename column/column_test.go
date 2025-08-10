package column

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDeltaColumn(t *testing.T) {
	buf := []byte{0x7f, 0x03, 0x03, 0x01, 0x7d, 0x03, 0x7e, 0x01}

	expected := []int64{3, 4, 5, 6, 9, 7, 8}
	var res []int64

	for i, err := range ReadDeltaColumn(bytes.NewReader(buf), uint64(len(buf))) {
		require.NoError(t, err)
		res = append(res, i)
	}

	require.Equal(t, expected, res)
}
