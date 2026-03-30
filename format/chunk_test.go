package format

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

	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

// Sources:
// - https://github.com/automerge/automerge/blob/main/interop/exemplar
// - https://github.com/automerge/automerge-perf

func TestReadDocument(t *testing.T) {
	for _, tc := range []struct {
		name   string
		chunks int
	}{
		// {name: "64bit_obj_id_change.automerge", chunks: 1}, // those are actually invalid, counter overflows
		// {name: "64bit_obj_id_doc.automerge", chunks: 1},
		{name: "counter_value_is_ok.automerge", chunks: 1},
		{name: "exemplar", chunks: 1},
		{name: "two_change_chunks.automerge", chunks: 2},
		{name: "two_change_chunks_compressed.automerge", chunks: 2},
		{name: "two_change_chunks_out_of_order.automerge", chunks: 2},
		{name: "text-edits.amrg", chunks: 259779},
		{name: "list_sequential.automerge", chunks: 1},
		{name: "list_concurrent_inserts.automerge", chunks: 1},
		{name: "list_with_delete.automerge", chunks: 1},
		{name: "list_insert_after_deleted.automerge", chunks: 1},
		{name: "text_sequential.automerge", chunks: 1},
		{name: "map_conflict.automerge", chunks: 1},
		{name: "map_delete.automerge", chunks: 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, err := os.Open("../testdata/" + tc.name)
			require.NoError(t, err)
			defer f.Close()

			r := ioutil.NewPagedReader(f)

			var chunks int
			for {
				_, toSkip, err := ReadChunk(r)
				require.NoError(t, err)
				// fmt.Println(c)
				// fmt.Println()
				// fmt.Println()

				err = r.Skip(toSkip)
				require.NoError(t, err)

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
	f, err := os.Open("../testdata/text-edits.amrg")
	require.NoError(t, err)
	r := ioutil.NewPagedReader(f)

	var chunks []Chunk
	for {
		c, toSkip, err := ReadChunk(r)
		require.NoError(t, err)
		chunks = append(chunks, c)
		err = r.Skip(toSkip)
		require.NoError(t, err)
		if r.Empty() {
			break
		}
	}
	fmt.Println(size.Of(chunks))
}

var chunksExemplar []Chunk

func BenchmarkReadExamplar(b *testing.B) {
	b.SetBytes(406)

	b.ReportAllocs()
	var chunkCount int
	for i := 0; i < b.N; i++ {
		f, _ := os.Open("../testdata/exemplar")
		r := ioutil.NewPagedReader(f)
		for {
			chunksExemplar = nil
			c, toSkip, err := ReadChunk(r)
			if errors.Is(err, io.EOF) {
				break
			}
			chunksExemplar = append(chunksExemplar, c)
			chunkCount++
			_ = r.Skip(toSkip)
		}
		_ = f.Close()
	}
	b.ReportMetric(float64(chunkCount)/float64(b.N), "chunks")
}

var chunksLarge []Chunk

func BenchmarkReadLarge(b *testing.B) {
	b.SetBytes(29249554)

	b.ReportAllocs()
	var chunkCount int
	for i := 0; i < b.N; i++ {
		chunksLarge = nil
		f, _ := os.Open("../testdata/text-edits.amrg")
		r := ioutil.NewPagedReader(f)
		for {
			c, toSkip, err := ReadChunk(r)
			if errors.Is(err, io.EOF) {
				break
			}
			chunksLarge = append(chunksLarge, c)
			chunkCount++
			_ = r.Skip(toSkip)
		}
		_ = f.Close()
	}
	b.ReportMetric(float64(chunkCount)/float64(b.N), "chunks")
}

var chunkCompressed []Chunk

func BenchmarkCompressed(b *testing.B) {
	b.SetBytes(192)

	b.ReportAllocs()
	var chunkCount int
	for i := 0; i < b.N; i++ {
		chunkCompressed = nil
		f, _ := os.Open("../testdata/two_change_chunks_compressed.automerge")
		r := ioutil.NewPagedReader(f)
		for {
			c, toSkip, err := ReadChunk(r)
			if errors.Is(err, io.EOF) {
				break
			}
			chunkCompressed = append(chunkCompressed, c)
			chunkCount++
			_ = r.Skip(toSkip)
		}
		_ = f.Close()
	}
	b.ReportMetric(float64(chunkCount)/float64(b.N), "chunks")
}

// TestInvalidChunks checks that malformed Automerge files are rejected.
func TestInvalidChunks(t *testing.T) {
	// counter_value_is_overlong uses non-minimal LEB128 encoding, which the spec forbids.
	// The error surfaces during column iteration, not chunk parsing.
	t.Run("counter_value_is_overlong", func(t *testing.T) {
		f, err := os.Open("../testdata/counter_value_is_overlong.automerge")
		require.NoError(t, err)
		defer f.Close()

		c, _, err := ReadChunk(ioutil.NewPagedReader(f))
		require.NoError(t, err)
		cc := c.(*ChangeChunk)

		var iterErr error
		for _, e := range cc.Operations() {
			if e != nil {
				iterErr = e
				break
			}
		}
		require.Error(t, iterErr, "overlong LEB128 should be rejected during iteration")
	})

	// counter_value_has_incorrect_meta parses today without error because we don't yet
	// validate counter metadata consistency (the Rust impl does). This test documents the
	// current behaviour so we notice if it accidentally starts failing.
	// TODO: add metadata validation and flip this to require.Error.
	t.Run("counter_value_has_incorrect_meta", func(t *testing.T) {
		f, err := os.Open("../testdata/counter_value_has_incorrect_meta.automerge")
		require.NoError(t, err)
		defer f.Close()

		c, _, err := ReadChunk(ioutil.NewPagedReader(f))
		require.NoError(t, err)
		cc := c.(*ChangeChunk)

		var iterErr error
		for _, e := range cc.Operations() {
			if e != nil {
				iterErr = e
				break
			}
		}
		require.NoError(t, iterErr, "metadata validation not yet implemented; update this test when it is")
	})
}

// TestDocumentChunkChanges verifies that Changes() on a document chunk correctly
// yields change metadata with resolved ActorIds, seqnums, and dep indices.
func TestDocumentChunkChanges(t *testing.T) {
	f, err := os.Open("../testdata/exemplar")
	require.NoError(t, err)
	defer f.Close()

	c, _, err := ReadChunk(ioutil.NewPagedReader(f))
	require.NoError(t, err)
	doc := c.(*DocumentChunk)

	var changes []types.DocChange
	for ch, err := range doc.Changes() {
		require.NoError(t, err)
		changes = append(changes, ch)
	}

	require.NotEmpty(t, changes, "exemplar should have at least one change")
	// The exemplar has one change from one actor
	require.Len(t, changes, 1)
	require.NotEmpty(t, changes[0].ActorId, "actor ID must be resolved")
	require.Equal(t, uint64(1), changes[0].SeqNum)
	require.Greater(t, changes[0].MaxOp, uint64(0))
}

// TestChangeChunkOperationIds verifies that Operations() on a change chunk sets
// correct Id values: ActorIdx=0 (change's own actor) and Counter=startOp+i.
//
// Operations must be iterated before calling r.Skip(), because column iterators
// hold lazy references into the paged reader's ring buffer.
func TestChangeChunkOperationIds(t *testing.T) {
	f, err := os.Open("../testdata/two_change_chunks.automerge")
	require.NoError(t, err)
	defer f.Close()

	r := ioutil.NewPagedReader(f)

	var totalChunks int
	for !r.Empty() {
		c, toSkip, err := ReadChunk(r)
		require.NoError(t, err)

		cc := c.(*ChangeChunk)
		var opIdx uint64
		for op, err := range cc.Operations() {
			require.NoError(t, err)
			require.Equal(t, uint32(0), op.Id.ActorIdx, "ActorIdx should be 0 (local actor) in change chunk")
			require.Equal(t, uint32(cc.StartOp+opIdx), op.Id.Counter,
				"Counter should be startOp+i for op %d", opIdx)
			opIdx++
		}
		require.Greater(t, opIdx, uint64(0), "change should have at least one operation")

		require.NoError(t, r.Skip(toSkip))
		totalChunks++
	}
	require.Equal(t, 2, totalChunks)
}

func TestEmptyDocumentRead(t *testing.T) {
	buf := []byte{0x85, 0x6f, 0x4a, 0x83, 0xb8, 0x1a, 0x95, 0x44, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00}
	r := ioutil.NewBytesReader(buf)
	c, toSkip, err := ReadChunk(r)
	require.NoError(t, err)
	require.NoError(t, r.Skip(toSkip))
	require.True(t, r.Empty()) // we should consume everything
	// TODO: add assertions once the struct is stable
	fmt.Println(c)
}

// TODO: I only have a header for compressed chunk, I need the whole thing to test
// func TestCompressedChunkRead(t *testing.T) {
// 	buf := []byte{0x85, 0x6f, 0x4a, 0x83, 0x80, 0xb7, 0x5d, 0x54, 0x01, 0xf4, 0x02}
// 	// buf := []byte{0x85, 0x6f, 0x4a, 0x83, 0x80, 0xb7, 0x5d, 0x54, 0x02, 0xc3, 0x02}
// 	r := lbuf.FromBytes(buf)
// 	c, err := ReadChunk(r)
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

// TestChangeChunkHashes verifies that ReadChunk correctly populates ChangeChunk.Hash:
//   - the hash is non-zero for both plain and compressed chunks
//   - hashing is deterministic (re-reading produces the same value)
//   - a compressed change and its uncompressed equivalent hash to the same value
//     (the spec defines the hash over the uncompressed content in both cases)
func TestChangeChunkHashes(t *testing.T) {
	readHashes := func(path string) []types.ChangeHash {
		f, err := os.Open(path)
		require.NoError(t, err)
		defer f.Close()
		var hashes []types.ChangeHash
		r := ioutil.NewPagedReader(f)
		for !r.Empty() {
			c, toSkip, err := ReadChunk(r)
			require.NoError(t, err)
			hashes = append(hashes, c.(*ChangeChunk).Hash)
			require.NoError(t, r.Skip(toSkip))
		}
		return hashes
	}

	plain := readHashes("../testdata/two_change_chunks.automerge")
	compressed := readHashes("../testdata/two_change_chunks_compressed.automerge")

	require.NotEmpty(t, plain)
	for _, h := range plain {
		require.NotEqual(t, types.ChangeHash{}, h)
	}

	// Deterministic: reading the same file twice gives the same hashes.
	require.Equal(t, plain, readHashes("../testdata/two_change_chunks.automerge"))

	// Compressed and uncompressed forms of the same changes must hash identically.
	require.Equal(t, plain, compressed)
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
