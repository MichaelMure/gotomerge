package column

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	ioutil "gotomerge/utils/io"
)

func TestSpecification(t *testing.T) {
	for _, tc := range []struct {
		name    string
		id      uint32
		_type   Type
		deflate bool
	}{
		{
			name:    "basic",
			id:      1,
			_type:   TypeString,
			deflate: true,
		},
		{
			name:    "min ID",
			id:      0,
			_type:   TypeGroup,
			deflate: true,
		},
		{
			name:    "max ID",
			id:      maxSpecificationId,
			_type:   TypeValue,
			deflate: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)

			spec := newSpecification(tc.id, tc._type, tc.deflate)
			err := writeSpecification(buf, spec)
			require.NoError(t, err)

			read, err := readSpecification(buf)
			require.NoError(t, err)

			require.Equal(t, spec, read)
			require.Equal(t, tc.id, read.ID())
			require.Equal(t, tc._type, read.Type())
			require.Equal(t, tc.deflate, read.Deflate())
		})
	}
}

func BenchmarkReadSpecification(b *testing.B) {
	buf := &bytes.Buffer{}
	spec := newSpecification(maxSpecificationId, TypeString, true)
	err := writeSpecification(buf, spec)
	require.NoError(b, err)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = readSpecification(buf)
	}
}

func BenchmarkWriteSpecification(b *testing.B) {
	spec := newSpecification(maxSpecificationId, TypeString, true)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = writeSpecification(io.Discard, spec)
	}
}

func FuzzSpecificationRoundTrip(f *testing.F) {
	f.Fuzz(func(t *testing.T, data []byte) {
		r := ioutil.NewBytesReader(data)
		spec, err := readSpecification(r)
		if err != nil {
			return // ignore invalid input
		}
		if !r.Empty() {
			return // ignore input with extra bytes
		}
		buf := new(bytes.Buffer)
		err = writeSpecification(buf, spec)
		if err != nil {
			t.Errorf("error writing specification: %v", err)
		}
		require.Equal(t, data, buf.Bytes())
	})
}
