package column

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDeltaColumn(t *testing.T) {
	buf := []byte{0x7f, 0x03, 0x03, 0x01, 0x7d, 0x03, 0x7e, 0x01}

	expected := []any{uint64(3), uint64(4), uint64(5), uint64(6), uint64(9), uint64(7), uint64(8)}
	var res []any

	for i, err := range ReadDeltaColumn(bytes.NewReader(buf)) {
		require.NoError(t, err)
		if val, valid := i.Value(); valid {
			res = append(res, val)
		} else {
			res = append(res, nil)
		}
	}

	require.Equal(t, expected, res)
}
