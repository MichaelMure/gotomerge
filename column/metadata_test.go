package column

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetadataRoundTrip(t *testing.T) {
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

func TestMetadataErrors(t *testing.T) {
	for _, tc := range []struct {
		name string
		bad  Metadata
	}{
		{
			name: "non ordered",
			bad: Metadata{
				{Spec: newSpecification(2, TypeString, false), Length: 2},
				{Spec: newSpecification(1, TypeString, true), Length: 1},
				{Spec: newSpecification(3, TypeString, false), Length: 2},
			},
		},
		{
			name: "value without a preceding metadata 1",
			bad: Metadata{
				{Spec: newSpecification(1, TypeValue, true), Length: 1},
			},
		},
		{
			name: "value without a preceding metadata 2",
			bad: Metadata{
				{Spec: newSpecification(1, TypeString, true), Length: 1},
				{Spec: newSpecification(1, TypeValue, true), Length: 1},
			},
		},
		{
			name: "value without a preceding metadata 3",
			bad: Metadata{
				{Spec: newSpecification(1, TypeValueMetadata, true), Length: 1},
				{Spec: newSpecification(2, TypeValue, true), Length: 1},
			},
		},
		{
			name: "multiple metadata with the same id",
			bad: Metadata{
				{Spec: newSpecification(1, TypeValueMetadata, true), Length: 1},
				{Spec: newSpecification(1, TypeString, false), Length: 1},
				{Spec: newSpecification(1, TypeValueMetadata, true), Length: 1},
				{Spec: newSpecification(1, TypeString, false), Length: 1},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			buf := bytes.Buffer{}
			err := WriteMetadata(&buf, tc.bad)
			require.NoError(t, err)

			_, err = ReadMetadata(&buf)
			require.Error(t, err)
		})
	}
}

func TestEmptyMetadata(t *testing.T) {
	data := []byte{0x00}
	buf := bytes.NewBuffer(data)
	meta, err := ReadMetadata(buf)
	require.NoError(t, err)
	require.Empty(t, meta)
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
