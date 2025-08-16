package column

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadBooleanColumn(t *testing.T) {
	buf := []byte{0x00, 0x02, 0x03}

	expected := []bool{true, true, false, false, false}
	var res []bool

	for b, err := range BooleanColumnFromBytes(buf).Iter() {
		require.NoError(t, err)
		res = append(res, b)
	}

	require.Equal(t, expected, res)
}
