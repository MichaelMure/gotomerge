package column

import (
	"io"
	"testing"

	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

func TestReadBooleanColumn(t *testing.T) {
	buf := []byte{0x00, 0x02, 0x03}

	expected := []bool{true, true, false, false, false}
	var res []bool

	r := NewBoolReader(ioutil.NewBytesReader(buf))
	for {
		b, err := r.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		res = append(res, b)
	}

	require.Equal(t, expected, res)
}
