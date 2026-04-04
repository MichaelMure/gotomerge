package ioutil

import (
	"bytes"
	"compress/flate"
	"errors"
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

// deflateCompress compresses data using DEFLATE (for use in tests only).
func deflateCompress(t *testing.T, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	w, err := flate.NewWriter(&buf, flate.DefaultCompression)
	require.NoError(t, err)
	_, err = w.Write(data)
	require.NoError(t, err)
	require.NoError(t, w.Close())
	return buf.Bytes()
}

func TestPagedReader_Suite(t *testing.T) {
	suiteSubReader(t, func(data []byte) SubReader {
		return NewPagedReader(bytes.NewReader(data))
	})
}

func TestPagedReader_ExtremeReallocation(t *testing.T) {
	extremeSize := 100 * pageSize // 100 pages (1.6MB), much larger than initial 8 page capacity
	g := generate(extremeSize)
	r := NewPagedReader(g)

	// Create a sub-reader that requires buffering the entire content
	sub, err := r.SubReader(0, uint64(extremeSize))
	require.NoError(t, err)

	// Verify the ring buffer has grown significantly beyond initial size
	require.True(t, len(r.(*pagedReader).pages) > 8, "Ring buffer should have grown substantially from initial 8 pages to accommodate 100 pages")

	// Read all data through the sub-reader to verify integrity after multiple reallocations
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, extremeSize, len(result))
	require.Equal(t, g.All(), result)
}

func TestPagedReader_Peek_EmptySource(t *testing.T) {
	g := generate(0)
	r := NewPagedReader(g)

	var buf bytes.Buffer
	err := r.Peek(&buf, 10)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, buf.Len())
}

func TestPagedReader_Peek_ExactAvailable(t *testing.T) {
	g := generate(1500)
	r := NewPagedReader(g)

	// Read some data to set position
	buf := make([]byte, 700)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 700, n)

	// Peek at all remaining data
	var peekBuf bytes.Buffer
	err = r.Peek(&peekBuf, 800) // Exactly remaining bytes
	require.NoError(t, err)
	require.Equal(t, g.All()[700:], peekBuf.Bytes())
}

func TestPagedReader_Peek_MoreThanAvailable(t *testing.T) {
	g := generate(1000)
	r := NewPagedReader(g)

	// Read some data
	buf := make([]byte, 600)
	_, err := r.Read(buf)
	require.NoError(t, err)

	// Try to peek more than remaining
	var peekBuf bytes.Buffer
	err = r.Peek(&peekBuf, 500) // Only 400 bytes remaining
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, peekBuf.Len())
}

func TestPagedReader_Peek_GrowBuffer(t *testing.T) {
	g := generate(100 * 1024) // Large data to trigger growth
	r := NewPagedReader(g)

	// Peek at a large amount that will require growing the buffer
	var buf bytes.Buffer
	err := r.Peek(&buf, 80*1024)
	require.NoError(t, err)
	require.Equal(t, g.All()[:80*1024], buf.Bytes())

	// Position should not have changed
	result, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All(), result)
}

func TestPagedReader_Peek_AtPageBoundary(t *testing.T) {
	g := generate(3 * pageSize) // Three full pages
	r := NewPagedReader(g)

	// Read to position exactly at page boundary
	readBuf := make([]byte, pageSize)
	n, err := r.Read(readBuf)
	require.NoError(t, err)
	require.Equal(t, pageSize, n)

	// Peek starting exactly at page boundary
	var peekBuf bytes.Buffer
	err = r.Peek(&peekBuf, pageSize)
	require.NoError(t, err)
	require.Equal(t, g.All()[pageSize:2*pageSize], peekBuf.Bytes())
}

func TestPagedReader_Peek_WriterError(t *testing.T) {
	g := generate(1000)
	r := NewPagedReader(g)

	// Create a writer that will fail after writing some bytes
	failingWriter := &failingWriter{failAfter: 50}

	err := r.Peek(failingWriter, 100)
	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

func TestPagedReader_SubReader_AfterMainReaderRead(t *testing.T) {
	g := generate(10000)
	r := NewPagedReader(g)

	// Read some data with main reader first
	buf := make([]byte, 3000)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 3000, n)

	// Now create a sub-reader for data after the consumed portion
	sub, err := r.SubReader(2000, 3000) // This should read from position 5000-8000 in original data
	require.NoError(t, err)
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[5000:8000], result)
}

func TestPagedReader_SubReaderOffset_AfterRead(t *testing.T) {
	g := generate(15000)
	r := NewPagedReader(g)

	// Read some data first
	buf := make([]byte, 4000)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 4000, n)

	// Create sub-reader from offset relative to current position
	sub, err := r.SubReaderOffset(2000)
	require.NoError(t, err)
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[6000:], result) // 4000 (read) + 2000 (offset)

	// Original reader should still be at position 4000
	remaining, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All()[4000:], remaining)
}

func TestPagedReader_SubReaderOffset_EmptyResult(t *testing.T) {
	g := generate(1000)
	r := NewPagedReader(g)

	// Try to create sub-reader from end of data
	sub, err := r.SubReaderOffset(1000)
	require.NoError(t, err)
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Empty(t, result)
}

func TestPagedReader_SubReaderOffset_OffsetBeyondData(t *testing.T) {
	g := generate(1000)
	r := NewPagedReader(g)

	// Try to create sub-reader from beyond available data
	sub, err := r.SubReaderOffset(2000)
	require.Error(t, err)
	require.Nil(t, sub)
}

func TestPagedReader_SubReaderOffset_WithSkip(t *testing.T) {
	g := generate(20000)
	r := NewPagedReader(g)

	// Skip some data
	err := r.Skip(5000)
	require.NoError(t, err)

	// Create sub-reader from offset relative to current position
	sub, err := r.SubReaderOffset(3000)
	require.NoError(t, err)
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[8000:], result) // 5000 (skipped) + 3000 (offset)

	// Original reader should still be at position 5000
	remaining, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All()[5000:], remaining)
}

func TestPagedReader_SubReaderOffset_MultipleSubReaders(t *testing.T) {
	g := generate(30000)
	r := NewPagedReader(g)

	// Create multiple sub-readers at different offsets
	sub1, err := r.SubReaderOffset(1000)
	require.NoError(t, err)

	sub2, err := r.SubReaderOffset(5000)
	require.NoError(t, err)

	sub3, err := r.SubReaderOffset(10000)
	require.NoError(t, err)

	// Read from each sub-reader
	result1, err := io.ReadAll(sub1)
	require.NoError(t, err)
	require.Equal(t, g.All()[1000:], result1)

	result2, err := io.ReadAll(sub2)
	require.NoError(t, err)
	require.Equal(t, g.All()[5000:], result2)

	result3, err := io.ReadAll(sub3)
	require.NoError(t, err)
	require.Equal(t, g.All()[10000:], result3)

	// Original reader position should not have changed
	originalResult, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All(), originalResult)
}

func TestPagedReader_SubReaderOffset_SmallReads(t *testing.T) {
	g := generate(5000)
	r := NewPagedReader(g)

	sub, err := r.SubReaderOffset(2000)
	require.NoError(t, err)

	// Read one byte at a time from the sub-reader
	var result []byte
	buf := make([]byte, 1)

	for {
		n, err := sub.Read(buf)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.Equal(t, 1, n)
		result = append(result, buf[:n]...)
	}

	require.Equal(t, g.All()[2000:], result)
}

func TestPagedReader_Skip_AfterSubReader(t *testing.T) {
	g := generate(10000)
	r := NewPagedReader(g)

	// Create and read from a sub-reader
	sub, err := r.SubReader(0, 2000)
	require.NoError(t, err)
	subData, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[:2000], subData)

	// Skip the data that was read by the sub-reader
	err = r.Skip(2000)
	require.NoError(t, err)

	// Now main reader should continue from position 2000
	buf := make([]byte, 1000)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1000, n)
	require.Equal(t, g.All()[2000:3000], buf)
}

func TestPagedReader_Skip_PartialSubReader(t *testing.T) {
	g := generate(8000)
	r := NewPagedReader(g)

	// Create sub-reader and read only part of it
	sub, err := r.SubReader(500, 3000)
	require.NoError(t, err)
	buf := make([]byte, 1500)
	n, err := sub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1500, n)
	require.Equal(t, g.All()[500:2000], buf)

	// Skip only what was actually read by the sub-reader
	err = r.Skip(2000) // 500 offset + 1500 read = 2000 total
	require.NoError(t, err)

	// Main reader should continue from position 2000
	remaining, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All()[2000:], remaining)
}

func TestPagedReader_Skip_MultipleSubReaders(t *testing.T) {
	g := generate(15000)
	r := NewPagedReader(g)

	// First sub-reader
	sub1, err := r.SubReader(0, 3000)
	require.NoError(t, err)
	data1, err := io.ReadAll(sub1)
	require.NoError(t, err)
	require.Equal(t, g.All()[:3000], data1)

	// Second sub-reader
	sub2, err := r.SubReader(3000, 2500)
	require.NoError(t, err)
	data2, err := io.ReadAll(sub2)
	require.NoError(t, err)
	require.Equal(t, g.All()[3000:5500], data2)

	// Skip both sub-readers' data
	err = r.Skip(5500)
	require.NoError(t, err)

	// Main reader should continue from position 5500
	buf := make([]byte, 2000)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 2000, n)
	require.Equal(t, g.All()[5500:7500], buf)
}

func TestPagedReader_Skip_WithOffset(t *testing.T) {
	g := generate(8000)
	r := NewPagedReader(g)

	// Sub-reader starting at an offset
	sub, err := r.SubReader(1500, 2000)
	require.NoError(t, err)
	subData, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[1500:3500], subData)

	// Skip including the offset
	err = r.Skip(3500)
	require.NoError(t, err)

	// Main reader should continue from position 3500
	buf := make([]byte, 1000)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1000, n)
	require.Equal(t, g.All()[3500:4500], buf)
}

func TestPagedReader_Skip_ExactAvailable(t *testing.T) {
	g := generate(4000)
	r := NewPagedReader(g)

	// Create sub-reader for all available data
	sub, err := r.SubReader(0, 3000)
	require.NoError(t, err)
	subData, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[:3000], subData)

	// Skip all available data
	err = r.Skip(3000)
	require.NoError(t, err)

	// Should trigger growth when trying to read more
	buf := make([]byte, 500)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 500, n)
	require.Equal(t, g.All()[3000:3500], buf)
}

func TestPagedReader_Skip_MixedWithDirectReads(t *testing.T) {
	g := generate(10000)
	r := NewPagedReader(g)

	// Direct read from main reader
	buf := make([]byte, 1000)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1000, n)
	require.Equal(t, g.All()[:1000], buf)

	// Use sub-reader on remaining data
	sub, err := r.SubReader(0, 2000) // relative to current position
	require.NoError(t, err)
	subData, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[1000:3000], subData)

	// Skip the sub-reader data
	err = r.Skip(2000)
	require.NoError(t, err)

	// Continue with direct read
	n, err = r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1000, n)
	require.Equal(t, g.All()[3000:4000], buf)
}

func TestPagedReader_Empty_AfterSubReader(t *testing.T) {
	g := generate(2000)
	r := NewPagedReader(g)

	// Use sub-readers to consume data
	sub1, err := r.SubReader(0, 800)
	require.NoError(t, err)
	_, err = io.ReadAll(sub1)
	require.NoError(t, err)

	sub2, err := r.SubReader(800, 1200)
	require.NoError(t, err)
	_, err = io.ReadAll(sub2)
	require.NoError(t, err)

	// Skip consumed data
	err = r.Skip(2000)
	require.NoError(t, err)

	// Should be empty after consuming all via sub-readers
	require.True(t, r.Empty())
}

func TestPagedReader_Empty_MixedOperations(t *testing.T) {
	g := generate(5000)
	r := NewPagedReader(g)

	require.False(t, r.Empty())

	// Read some data
	buf := make([]byte, 1000)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1000, n)
	require.False(t, r.Empty())

	// Use sub-reader
	sub, err := r.SubReader(500, 1500)
	require.NoError(t, err)
	_, err = io.ReadAll(sub)
	require.NoError(t, err)
	require.False(t, r.Empty()) // Haven't skipped yet

	// Skip the sub-reader data
	err = r.Skip(2000)
	require.NoError(t, err)
	require.False(t, r.Empty()) // Still has more data in source

	// Read remaining
	_, err = io.ReadAll(r)
	require.NoError(t, err)
	require.True(t, r.Empty())
}

func TestSubReader_Suite(t *testing.T) {
	suiteSubReader(t, func(data []byte) SubReader {
		r := NewPagedReader(bytes.NewReader(data))
		sub, err := r.SubReader(0, uint64(len(data)))
		if err != nil {
			panic(err)
		}
		return sub
	})
}

func TestSubReader_SubReader_AtEnd(t *testing.T) {
	g := generate(8000)
	r := NewPagedReader(g)

	parentSub, err := r.SubReader(2000, 4000)
	require.NoError(t, err)

	// Create child sub-reader at the very end of parent
	childSub, err := parentSub.SubReader(3500, 500)
	require.NoError(t, err)

	result, err := io.ReadAll(childSub)
	require.NoError(t, err)
	require.Equal(t, g.All()[5500:6000], result) // 2000 + 3500 = 5500, 5500 + 500 = 6000
}

func TestSubReader_SubReader_OutOfBounds(t *testing.T) {
	g := generate(5000)
	r := NewPagedReader(g)

	parentSub, err := r.SubReader(1000, 2000)
	require.NoError(t, err)

	// Try to create child sub-reader beyond parent bounds
	_, err = parentSub.SubReader(1500, 1000) // offset + size = 2500 > 2000 available
	require.Error(t, err)
	require.Contains(t, err.Error(), "offset + size > available")
}

func TestSubReader_SubReader_OffsetOutOfBounds(t *testing.T) {
	g := generate(5000)
	r := NewPagedReader(g)

	parentSub, err := r.SubReader(1000, 2000)
	require.NoError(t, err)

	// Try to create child sub-reader with offset beyond parent size
	_, err = parentSub.SubReader(2500, 100) // offset 2500 > 2000 available
	require.Error(t, err)
	require.Contains(t, err.Error(), "offset + size > available")
}

func TestSubReader_SubReaderOffset_OffsetBeyondAllowed(t *testing.T) {
	data := []byte{1, 2, 3, 4, 5}
	r := NewPagedReader(bytes.NewReader(data))
	sub, err := r.SubReader(0, 5)
	require.NoError(t, err)

	_, err = sub.SubReaderOffset(6)
	require.Error(t, err)
}

func TestSubReaderOffset_SubReaderOffset_RequiresPageLoad(t *testing.T) {
	size := pageSize + 500
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 199)
	}
	r := NewPagedReader(bytes.NewReader(data))
	sub, err := r.SubReaderOffset(0)
	require.NoError(t, err)

	// Skip into the second page — forces loadPage inside SubReaderOffset loop.
	offset := uint64(pageSize + 100)
	fork, err := sub.SubReaderOffset(offset)
	require.NoError(t, err)
	result, err := io.ReadAll(fork)
	require.NoError(t, err)
	require.Equal(t, data[offset:], result)
}

func TestSubReaderOffset_Skip_PastEOF(t *testing.T) {
	data := []byte{1, 2, 3}
	r := NewPagedReader(bytes.NewReader(data))
	sub, err := r.SubReaderOffset(0)
	require.NoError(t, err)

	// Read all data to exhaust the source, then skip — should return EOF.
	_, err = io.ReadAll(sub)
	require.NoError(t, err)

	sub2, err := r.SubReaderOffset(0)
	require.NoError(t, err)
	err = sub2.Skip(len(data) + 1)
	require.Error(t, err)
}

func TestSubReader_Skip_OutOfBounds(t *testing.T) {
	g := generate(4000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(500, 1500)
	require.NoError(t, err)

	// Try to skip more than available
	err = sub.Skip(2000) // More than 1500 available
	require.Error(t, err)
	require.Contains(t, err.Error(), "skip 2000 > available 1500")
}

func TestSubReader_Skip_MixedWithReads(t *testing.T) {
	g := generate(12000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(2000, 6000)
	require.NoError(t, err)

	// Read, skip, read pattern
	buf := make([]byte, 1000)

	// Read 1000 bytes
	n, err := sub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1000, n)
	require.Equal(t, g.All()[2000:3000], buf)

	// Skip 1500 bytes
	err = sub.Skip(1500)
	require.NoError(t, err)

	// Read 1000 bytes
	n, err = sub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1000, n)
	require.Equal(t, g.All()[4500:5500], buf) // 2000 + 1000 + 1500 = 4500

	// Skip remaining
	err = sub.Skip(2500) // 6000 - 1000 - 1500 - 1000 = 2500
	require.NoError(t, err)

	// Should be empty
	n, err = sub.Read(buf)
	require.Equal(t, 0, n)
	require.Equal(t, io.EOF, err)
}

func TestSubReader_Empty_AfterChildSubReader(t *testing.T) {
	g := generate(8000)
	r := NewPagedReader(g)

	// Create parent sub-reader
	parentSub, err := r.SubReader(1000, 4000)
	require.NoError(t, err)

	// Create child sub-reader that consumes all parent data
	childSub, err := parentSub.SubReader(0, 4000)
	require.NoError(t, err)

	// Read all child data
	_, err = io.ReadAll(childSub)
	require.NoError(t, err)

	// Parent should not be empty yet (child reading doesn't affect parent state)
	require.False(t, parentSub.Empty())

	// Read from parent should still work
	remaining, err := io.ReadAll(parentSub)
	require.NoError(t, err)
	require.Equal(t, g.All()[1000:5000], remaining)

	// Now parent should be empty
	require.True(t, parentSub.Empty())
}

func TestSubReaderOffset_Suite(t *testing.T) {
	suiteSubReader(t, func(data []byte) SubReader {
		r := NewPagedReader(bytes.NewReader(data))
		sub, err := r.SubReaderOffset(0)
		if err != nil {
			panic(err)
		}
		return sub
	})
}

func TestSubReaderOffset_Peek_FromMainReaderPosition(t *testing.T) {
	g := generate(12000)
	r := NewPagedReader(g)

	// Read some data with main reader
	buf := make([]byte, 3000)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 3000, n)

	// Create sub-reader from offset relative to current position
	sub, err := r.SubReaderOffset(2000)
	require.NoError(t, err)

	// Peek from the sub-reader
	var peekBuf bytes.Buffer
	err = sub.Peek(&peekBuf, 1000)
	require.NoError(t, err)
	require.Equal(t, g.All()[5000:6000], peekBuf.Bytes()) // 3000 + 2000 = 5000

	// Main reader position should not have changed
	remaining, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All()[3000:], remaining)
}

func TestSubReaderOffset_Peek_AtPageBoundary(t *testing.T) {
	g := generate(3 * pageSize)
	r := NewPagedReader(g)

	// Create sub-reader starting at page boundary
	sub, err := r.SubReaderOffset(pageSize)
	require.NoError(t, err)

	var buf bytes.Buffer
	err = sub.Peek(&buf, pageSize)
	require.NoError(t, err)
	require.Equal(t, g.All()[pageSize:2*pageSize], buf.Bytes())

	// Position should not have changed
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[pageSize:], result)
}

func TestSubReaderOffset_Peek_NeedToLoadPages(t *testing.T) {
	g := generate(80 * 1024) // Large data requiring page loading
	r := NewPagedReader(g)

	sub, err := r.SubReaderOffset(5 * 1024)
	require.NoError(t, err)

	// Peek at a large amount that will require loading more pages
	var buf bytes.Buffer
	err = sub.Peek(&buf, 50*1024)
	require.NoError(t, err)
	require.Equal(t, g.All()[5*1024:55*1024], buf.Bytes())

	// Position should not have changed
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[5*1024:], result)
}

func TestSubReaderOffset_Peek_WriterError(t *testing.T) {
	g := generate(2000)
	r := NewPagedReader(g)

	sub, err := r.SubReaderOffset(500)
	require.NoError(t, err)

	// Create a writer that will fail after writing some bytes
	failingWriter := &failingWriter{failAfter: 100}

	err = sub.Peek(failingWriter, 300)
	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

func TestSubReaderOffset_Peek_EmptyAtOffset(t *testing.T) {
	g := generate(1000)
	r := NewPagedReader(g)

	// Create sub-reader at the very end
	sub, err := r.SubReaderOffset(1000)
	require.NoError(t, err)

	var buf bytes.Buffer
	err = sub.Peek(&buf, 100)
	require.Equal(t, io.EOF, err)
	require.Equal(t, 0, buf.Len())
}

func TestBytesReader_Suite(t *testing.T) {
	suiteSubReader(t, func(data []byte) SubReader {
		return NewBytesReader(bytes.Clone(data))
	})
}

func TestBytesReader_SubReader_OutOfBounds(t *testing.T) {
	g := generate(80)
	r := NewBytesReader(g.All())

	// Offset + size exceeds available data
	_, err := r.SubReader(70, 20) // 70 + 20 = 90 > 80
	require.Error(t, err)
	require.Contains(t, err.Error(), "offset + size > len(data)")

	// Test after reading some data
	buf := make([]byte, 30)
	_, err = r.Read(buf)
	require.NoError(t, err)

	// Now only 50 bytes remain at position 30
	_, err = r.SubReader(30, 25) // position(30) + offset(30) + size(25) = 85 > 80
	require.Error(t, err)
	require.Contains(t, err.Error(), "offset + size > len(data)")
}

func TestBytesReader_Skip_OutOfBounds(t *testing.T) {
	g := generate(30)
	r := NewBytesReader(g.All())

	err := r.Skip(40) // More than available
	require.Error(t, err)
	require.Contains(t, err.Error(), "skip 40 > len(data)")
}

func TestBytesReader_Peek_OutOfBounds(t *testing.T) {
	g := generate(30)
	r := NewBytesReader(g.All())

	var buf bytes.Buffer
	err := r.Peek(&buf, 50) // More than available
	require.Error(t, err)
	require.Contains(t, err.Error(), "peek 50 > len(data) - position")
}

func TestBytesReader_Peek_WriterError(t *testing.T) {
	g := generate(50)
	r := NewBytesReader(g.All())

	// Create a writer that will fail
	failingWriter := &failingWriter{failAfter: 10}

	err := r.Peek(failingWriter, 25)
	require.Error(t, err)
	require.Contains(t, err.Error(), "write failed")
}

func TestBytesReader_Integration_AllMethods(t *testing.T) {
	g := generate(500)
	r := NewBytesReader(g.All())

	// Test initial state
	require.False(t, r.Empty())

	// Test peek
	var peekBuf bytes.Buffer
	err := r.Peek(&peekBuf, 50)
	require.NoError(t, err)
	require.Equal(t, g.All()[:50], peekBuf.Bytes())

	// Test read
	buf := make([]byte, 50)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 50, n)
	require.Equal(t, g.All()[:50], buf)

	// Test skip
	err = r.Skip(100) // skip bytes 50-149
	require.NoError(t, err)

	// Test sub-reader
	sub, err := r.SubReader(50, 80) // bytes 200-279 (position 150 + offset 50)
	require.NoError(t, err)

	subResult, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[200:280], subResult)

	// Continue with main reader
	remaining, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All()[150:], remaining)

	// Test empty state
	require.True(t, r.Empty())
}

func TestBytesReader_FullCycle(t *testing.T) {
	// Test with full 256-byte cycle from generator
	g := generate(256)
	r := NewBytesReader(g.All())

	// Test all methods with the complete byte range
	var peekBuf bytes.Buffer
	err := r.Peek(&peekBuf, 100)
	require.NoError(t, err)
	require.Equal(t, g.All()[:100], peekBuf.Bytes())

	// Read and verify
	result, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All(), result)

	// Verify we have all byte values 0-255
	for i := 0; i < 256; i++ {
		require.Equal(t, byte(i), result[i], "Byte at position %d should be %d", i, i)
	}
}

func TestBytesReader_Deflate_AfterPartialRead(t *testing.T) {
	// Prefix 4 bytes before the compressed payload.
	plain := []byte("hello deflate world")
	compressed := deflateCompress(t, plain)
	data := append([]byte{1, 2, 3, 4}, compressed...)

	r := NewBytesReader(data)
	// Advance past the prefix.
	_, err := io.ReadFull(r, make([]byte, 4))
	require.NoError(t, err)

	got, n, err := r.Deflate()
	require.NoError(t, err)
	require.Equal(t, len(plain), n)
	result, err := io.ReadAll(got)
	require.NoError(t, err)
	require.Equal(t, plain, result)
}

func TestPagedReader_HasAtLeast_LoadsPages(t *testing.T) {
	// pagedReader starts with nothing buffered; HasAtLeast must load pages to answer.
	data := make([]byte, 2*pageSize+100)
	for i := range data {
		data[i] = byte(i % 251)
	}
	r := NewPagedReader(bytes.NewReader(data))
	pr := r.(*pagedReader)

	// Nothing loaded yet.
	require.Equal(t, 0, pr.available)
	require.True(t, r.HasAtLeast(pageSize+1))
	// Pages were loaded as a side effect.
	require.Greater(t, pr.available, 0)

	// Position must be unchanged: we can still read all the data.
	got, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, data, got)
}

func TestSubReaderOffset_HasAtLeast_LoadsPages(t *testing.T) {
	// subReaderOffset has no upper bound, so HasAtLeast must scan/load pages.
	data := make([]byte, 2*pageSize+100)
	for i := range data {
		data[i] = byte(i % 251)
	}
	r := NewPagedReader(bytes.NewReader(data))
	sub, err := r.SubReaderOffset(0)
	require.NoError(t, err)

	require.True(t, sub.HasAtLeast(pageSize+1))
	require.True(t, sub.HasAtLeast(2*pageSize+100))
	require.False(t, sub.HasAtLeast(2*pageSize+101))

	// Position must be unchanged.
	got, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, data, got)
}

// Helper type for testing writer errors
type failingWriter struct {
	written   int
	failAfter int
}

func (fw *failingWriter) Write(p []byte) (n int, err error) {
	if fw.written+len(p) > fw.failAfter {
		return fw.failAfter - fw.written, errors.New("write failed")
	}
	fw.written += len(p)
	return len(p), nil
}
