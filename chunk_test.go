package gotomerge

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/DmitriyVTitov/size"
	"github.com/stretchr/testify/require"

	"gotomerge/lbuf"
	"gotomerge/types"
)

// Sources:
// - https://github.com/automerge/automerge/blob/main/interop/exemplar
// - https://github.com/automerge/automerge-perf

func TestReadDocument(t *testing.T) {
	for _, tc := range []struct {
		name   string
		chunks int
	}{
		{name: "64bit_obj_id_change.automerge", chunks: 1},
		{name: "64bit_obj_id_doc.automerge", chunks: 1},
		{name: "counter_value_is_ok.automerge", chunks: 1},
		{name: "exemplar", chunks: 1},
		{name: "two_change_chunks.automerge", chunks: 2},
		{name: "two_change_chunks_compressed.automerge", chunks: 2},
		{name: "two_change_chunks_out_of_order.automerge", chunks: 2},
		{name: "text-edits.amrg", chunks: 259779},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open("testdata/" + tc.name)
			require.NoError(t, err)
			defer f.Close()

			r := lbuf.FromReader(f)
			defer r.Release()

			var chunks int
			for {
				_, err := readChunk(r)
				require.NoError(t, err)
				// fmt.Println(c)
				// fmt.Println()
				// fmt.Println()

				chunks++
				if r.Empty() {
					break
				}
			}

			require.Equal(t, tc.chunks, chunks)
			_, err = io.ReadFull(r, make([]byte, 1))
			require.ErrorIs(t, err, io.EOF) // we should consume everything
		})
	}
}

func TestLarge(t *testing.T) {
	f, err := os.Open("testdata/text-edits.amrg")
	require.NoError(t, err)
	r := lbuf.FromReader(f)
	defer r.Release()

	var chunks []chunk
	for {
		c, err := readChunk(r)
		require.NoError(t, err)
		chunks = append(chunks, c)
		if r.Empty() {
			break
		}
	}
	fmt.Println(size.Of(chunks))
}

var chunksExemplar []chunk

func BenchmarkReadExamplar(b *testing.B) {
	b.SetBytes(406)

	b.ReportAllocs()
	var chunkCount int
	for i := 0; i < b.N; i++ {
		f, _ := os.Open("testdata/exemplar")
		r := lbuf.FromReader(f)
		for {
			chunksExemplar = nil
			c, err := readChunk(r)
			if errors.Is(err, io.EOF) {
				break
			}
			chunksExemplar = append(chunksExemplar, c)
			chunkCount++
		}
		r.Release()
		_ = f.Close()
	}
	b.ReportMetric(float64(chunkCount)/float64(b.N), "chunks")
}

var chunksLarge []chunk

func BenchmarkReadLarge(b *testing.B) {
	b.SetBytes(29249554)

	b.ReportAllocs()
	var chunkCount int
	for i := 0; i < b.N; i++ {
		chunksLarge = nil
		f, _ := os.Open("testdata/text-edits.amrg")
		r := lbuf.FromReader(f)
		for {
			c, err := readChunk(r)
			if errors.Is(err, io.EOF) {
				break
			}
			chunksLarge = append(chunksLarge, c)
			chunkCount++
		}
		r.Release()
		_ = f.Close()
	}
	b.ReportMetric(float64(chunkCount)/float64(b.N), "chunks")
}

var chunkCompressed []chunk

func BenchmarkCompressed(b *testing.B) {
	b.SetBytes(192)

	b.ReportAllocs()
	var chunkCount int
	for i := 0; i < b.N; i++ {
		chunkCompressed = nil
		f, _ := os.Open("testdata/two_change_chunks_compressed.automerge")
		r := lbuf.FromReader(f)
		for {
			c, err := readChunk(r)
			if errors.Is(err, io.EOF) {
				break
			}
			chunkCompressed = append(chunkCompressed, c)
			chunkCount++
		}
		r.Release()
		_ = f.Close()
	}
	b.ReportMetric(float64(chunkCount)/float64(b.N), "chunks")
}

func TestEmptyDocumentRead(t *testing.T) {
	buf := []byte{0x85, 0x6f, 0x4a, 0x83, 0xb8, 0x1a, 0x95, 0x44, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00}
	r := lbuf.FromBytes(buf)
	c, err := readChunk(r)
	require.NoError(t, err)
	require.True(t, r.Empty()) // we should consume everything
	// TODO: add assertions once the struct is stable
	fmt.Println(c)
}

// TODO: I only have a header for compressed chunk, I need the whole thing to test
// func TestCompressedChunkRead(t *testing.T) {
// 	buf := []byte{0x85, 0x6f, 0x4a, 0x83, 0x80, 0xb7, 0x5d, 0x54, 0x01, 0xf4, 0x02}
// 	// buf := []byte{0x85, 0x6f, 0x4a, 0x83, 0x80, 0xb7, 0x5d, 0x54, 0x02, 0xc3, 0x02}
// 	r := lbuf.FromBytes(buf)
// 	c, err := readChunk(r)
// 	require.NoError(t, err)
// 	fmt.Println(c)
// }

func TestChangeHashesRoundTrip(t *testing.T) {
	data := []byte{0x02, 0xf9, 0x86, 0xa4, 0x31, 0x8d, 0x1f, 0x1c, 0xc0, 0xe2,
		0xe1, 0x0e, 0x42, 0x1e, 0x7a, 0x9a, 0x4c, 0xd0, 0xb7, 0x0a, 0x89, 0xda,
		0xe9, 0x8b, 0xc1, 0xd7, 0x6d, 0x78, 0x9c, 0x2b, 0xf7, 0x90, 0x4c, 0x43,
		0x55, 0xa4, 0x6b, 0x19, 0xd3, 0x48, 0xdc, 0x2f, 0x57, 0xc0, 0x46, 0xf8,
		0xef, 0x63, 0xd4, 0x53, 0x8e, 0xbb, 0x93, 0x60, 0x00, 0xf3, 0xc9, 0xee,
		0x95, 0x4a, 0x27, 0x46, 0x0d, 0xd8, 0x65}

	hashes, err := readChangeHashes(lbuf.FromBytes(data))
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
	ids, err := readActorIds(lbuf.FromBytes(data))
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
	indexes, err := readHeadIndexes(lbuf.FromBytes(data), 4)
	require.NoError(t, err)

	expected := []uint64{1, 2, 3, 17}
	require.Equal(t, expected, indexes)

	buf := new(bytes.Buffer)
	err = writeHeadIndexes(buf, indexes)
	require.NoError(t, err)
	require.Equal(t, data, buf.Bytes())
}
