package lbuf

import (
	"bytes"
	"compress/flate"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

type byteGenerator struct {
	to      int
	counter int
}

func generate(n int) *byteGenerator {
	return &byteGenerator{to: n, counter: 0}
}

func (b *byteGenerator) Read(p []byte) (n int, err error) {
	for i := 0; i < len(p); i++ {
		if b.to == b.counter {
			return i, io.EOF
		}
		p[i] = byte(b.counter % 256)
		b.counter++
	}
	return len(p), nil
}

func (b *byteGenerator) All() []byte {
	res := make([]byte, b.to)
	for i := 0; i < b.to; i++ {
		res[i] = byte(i % 256)
	}
	return res
}

func TestReader(t *testing.T) {
	g := generate(10000)
	r := FromReader(g)

	require.False(t, r.Empty())
	res, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All(), res)
	require.True(t, r.Empty())
}

func TestLimit(t *testing.T) {
	g := generate(10000)
	r := FromReader(g)
	r = r.Limit(100)

	require.False(t, r.Empty())
	res, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, generate(100).All(), res)
	require.True(t, r.Empty())
}

func TestProcessor(t *testing.T) {
	g := generate(10000)

	buf := &bytes.Buffer{}
	w, err := flate.NewWriter(buf, 2)
	require.NoError(t, err)
	_, err = w.Write(g.All())
	require.NoError(t, err)
	require.NoError(t, w.Close())
	compressed := buf.Bytes()

	// ReadAll
	r := FromBytes(compressed)
	r = r.Deflate()
	require.False(t, r.Empty())
	res, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All(), res)
	require.True(t, r.Empty())

	// only add a limiter after the processor, the other way doesn't make sense
	r = FromReader(bytes.NewReader(compressed))
	r = r.Deflate()
	r = r.Limit(100)
	res, err = io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, generate(100).All(), res)
	require.True(t, r.Empty())
}
