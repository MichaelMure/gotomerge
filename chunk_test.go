package gotomerge

import (
	"bytes"
	"embed"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/DmitriyVTitov/size"
	"github.com/stretchr/testify/require"

	"gotomerge/types"
)

// Sources:
// - https://github.com/automerge/automerge/blob/main/interop/exemplar
// - https://github.com/automerge/automerge-perf

//go:embed testdata
var testdata embed.FS

func TestReadDocument(t *testing.T) {
	for _, tc := range []struct {
		name string
	}{
		{name: "64bit_obj_id_change.automerge"},
		{name: "64bit_obj_id_doc.automerge"},
		{name: "counter_value_is_ok.automerge"},
		{name: "exemplar"},
		{name: "two_change_chunks.automerge"},
		{name: "text-edits.amrg"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			buf, err := testdata.ReadFile("testdata/" + tc.name)
			require.NoError(t, err)
			r := bytes.NewReader(buf)

			for r.Len() > 0 {
				c, err := readChunk(r)
				require.NoError(t, err)
				fmt.Println(c)
				fmt.Println()
				fmt.Println()
			}

			require.Zero(t, r.Len()) // we should consume everything
		})
	}
}

func TestLarge(t *testing.T) {
	buf, err := testdata.ReadFile("testdata/text-edits.amrg")
	require.NoError(t, err)
	r := bytes.NewReader(buf)
	var chunks []chunk
	for r.Len() > 0 {
		c, err := readChunk(r)
		require.NoError(t, err)
		chunks = append(chunks, c)
	}
	fmt.Println(size.Of(chunks))
}

func BenchmarkReadLarge(b *testing.B) {
	buf, _ := testdata.ReadFile("testdata/text-edits.amrg")
	b.SetBytes(int64(len(buf)))

	b.ReportAllocs()
	var chunkCount int
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(buf)
		for r.Len() > 0 {
			_, _ = readChunk(r)
			chunkCount++
		}
	}
	b.ReportMetric(float64(chunkCount)/float64(b.N), "chunks")
}

func TestEmptyDocumentRead(t *testing.T) {
	buf := []byte{0x85, 0x6f, 0x4a, 0x83, 0xb8, 0x1a, 0x95, 0x44, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00}
	r := bytes.NewReader(buf)
	c, err := readChunk(r)
	require.NoError(t, err)
	require.Zero(t, r.Len()) // we should consume everything
	// TODO: add assertions once the struct is stable
	fmt.Println(c)
}

// TODO: I only have a header for compressed chunk, I need the whole thing to test
// func TestCompressedChunkRead(t *testing.T) {
// 	buf := []byte{0x85, 0x6f, 0x4a, 0x83, 0x80, 0xb7, 0x5d, 0x54, 0x02, 0xc3, 0x02}
// 	r := bytes.NewReader(buf)
// 	c, err := readChunk(r)
// 	require.NoError(t, err)
// 	fmt.Println(c)
//
// 	require.Zero(t, r.Len()) // we should consume everything
// }

func TestChangeHashesRoundTrip(t *testing.T) {
	data := []byte{0x02, 0xf9, 0x86, 0xa4, 0x31, 0x8d, 0x1f, 0x1c, 0xc0, 0xe2,
		0xe1, 0x0e, 0x42, 0x1e, 0x7a, 0x9a, 0x4c, 0xd0, 0xb7, 0x0a, 0x89, 0xda,
		0xe9, 0x8b, 0xc1, 0xd7, 0x6d, 0x78, 0x9c, 0x2b, 0xf7, 0x90, 0x4c, 0x43,
		0x55, 0xa4, 0x6b, 0x19, 0xd3, 0x48, 0xdc, 0x2f, 0x57, 0xc0, 0x46, 0xf8,
		0xef, 0x63, 0xd4, 0x53, 0x8e, 0xbb, 0x93, 0x60, 0x00, 0xf3, 0xc9, 0xee,
		0x95, 0x4a, 0x27, 0x46, 0x0d, 0xd8, 0x65}

	hashes, err := readChangeHashes(bytes.NewReader(data))
	require.NoError(t, err)
	require.Len(t, hashes, 2)

	h1, err := hex.DecodeString("f986a4318d1f1cc0e2e10e421e7a9a4cd0b70a89dae98bc1d76d789c2bf7904c")
	require.NoError(t, err)
	h2, err := hex.DecodeString("4355a46b19d348dc2f57c046f8ef63d4538ebb936000f3c9ee954a27460dd865")
	require.NoError(t, err)

	expected := []types.ChangeHash{
		types.ChangeHash(h1[:32]),
		types.ChangeHash(h2[:32]),
	}

	require.Equal(t, expected, hashes)

	buf := new(bytes.Buffer)
	err = writeChangeHashes(buf, hashes)
	require.NoError(t, err)
	require.Equal(t, data, buf.Bytes())
}

func TestActorIdsRoundTrip(t *testing.T) {
	data := []byte{0x01, 0x03, 0xab, 0xcd, 0xef}
	ids, err := readActorIds(bytes.NewReader(data))
	require.NoError(t, err)

	require.Equal(t, 1, len(ids))
	require.Equal(t, types.ActorId{0xab, 0xcd, 0xef}, ids[0])

	buf := new(bytes.Buffer)
	err = writeActorIds(buf, ids)
	require.NoError(t, err)
	require.Equal(t, data, buf.Bytes())
}

func TestHeadIndexesRoundTrip(t *testing.T) {
	data := []byte{0x01, 0x02, 0x03, 0x11}
	indexes, err := readHeadIndexes(bytes.NewReader(data), 4)
	require.NoError(t, err)

	expected := []uint64{1, 2, 3, 17}
	require.Equal(t, expected, indexes)

	buf := new(bytes.Buffer)
	err = writeHeadIndexes(buf, indexes)
	require.NoError(t, err)
	require.Equal(t, data, buf.Bytes())
}
