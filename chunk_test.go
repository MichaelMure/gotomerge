package gotomerge

import (
	"bytes"
	_ "embed"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// from https://github.com/automerge/automerge/blob/main/interop/exemplar
//
//go:embed testdata/exemplar
var examplar []byte

func TestExemplarRead(t *testing.T) {
	c, err := readChunk(bytes.NewReader(examplar))
	require.NoError(t, err)
	fmt.Println(c)
}
