package column

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadata(t *testing.T) {
	meta := Metadata{
		{Spec: newSpecification(1, TypeString, true), Length: 1},
		{Spec: newSpecification(2, TypeString, false), Length: 2},
	}

	buf := bytes.Buffer{}
	err := WriteMetadata(&buf, meta)
	require.NoError(t, err)

	read, err := ReadMetadata(&buf)
	require.NoError(t, err)

	require.Equal(t, meta, read)
}

func BenchmarkReadMetadata(b *testing.B) {
	buf := bytes.Buffer{}
	meta := Metadata{
		{Spec: newSpecification(1, TypeString, true), Length: 1},
		{Spec: newSpecification(2, TypeString, false), Length: 2},
	}
	err := WriteMetadata(&buf, meta)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = ReadMetadata(bytes.NewReader(buf.Bytes()))
	}
}

func BenchmarkWriteMetadata(b *testing.B) {
	meta := Metadata{
		{Spec: newSpecification(1, TypeString, true), Length: 1},
		{Spec: newSpecification(2, TypeString, false), Length: 2},
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = WriteMetadata(io.Discard, meta)
	}
}
