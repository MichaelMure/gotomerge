package column

import (
	"io"
	"testing"

	ioutil "gotomerge/utils/io"

	"github.com/stretchr/testify/require"
)

func TestReadDeltaColumn(t *testing.T) {
	buf := []byte{0x7f, 0x03, 0x03, 0x01, 0x7d, 0x03, 0x7e, 0x01}

	expected := []any{int64(3), int64(4), int64(5), int64(6), int64(9), int64(7), int64(8)}
	var res []any

	r := NewDeltaReader(ioutil.NewBytesReader(buf))
	for {
		nv, err := r.Next()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		if val, valid := nv.Value(); valid {
			res = append(res, val)
		} else {
			res = append(res, nil)
		}
	}

	require.Equal(t, expected, res)
}
