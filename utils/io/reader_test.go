package ioutil

import (
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

func TestBufferedRead(t *testing.T) {
	g := generate(200 * 1024)
	r := NewPagedReader(g)

	res, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All(), res)
}

func TestPagedReader_EmptyRead(t *testing.T) {
	g := generate(0)
	r := NewPagedReader(g)

	buf := make([]byte, 10)
	n, err := r.Read(buf)
	require.Equal(t, 0, n)
	require.Equal(t, io.EOF, err)
}

func TestPagedReader_ZeroLengthRead(t *testing.T) {
	g := generate(100)
	r := NewPagedReader(g)

	n, err := r.Read(nil)
	require.Equal(t, 0, n)
	require.NoError(t, err)

	n, err = r.Read([]byte{})
	require.Equal(t, 0, n)
	require.NoError(t, err)
}

func TestPagedReader_SmallReads(t *testing.T) {
	g := generate(1000)
	r := NewPagedReader(g)

	var result []byte
	buf := make([]byte, 1) // Read one byte at a time

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.Equal(t, 1, n)
		result = append(result, buf[:n]...)
	}

	require.Equal(t, g.All(), result)
}

func TestPagedReader_LargeReads(t *testing.T) {
	// Test with data larger than page size
	g := generate(100 * 1024) // 100KB
	r := NewPagedReader(g)

	// Try to read more than available in one go
	buf := make([]byte, 200*1024) // 200KB buffer
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 100*1024, n)
	require.Equal(t, g.All(), buf[:n])
}

func TestPagedReader_MultipleReads(t *testing.T) {
	g := generate(50 * 1024) // 50KB
	r := NewPagedReader(g)

	var result []byte
	buf := make([]byte, 7*1024) // Read in 7KB chunks

	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		result = append(result, buf[:n]...)
	}

	require.Equal(t, g.All(), result)
}

func TestPagedReader_Skip(t *testing.T) {
	g := generate(10000)
	r := NewPagedReader(g)

	// Read some data to populate pages
	buf := make([]byte, 5000)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 5000, n)

	// Skip part of the read data
	err = r.Skip(2000)
	require.NoError(t, err)

	// Read remaining data
	remaining, err := io.ReadAll(r)
	require.NoError(t, err)

	// Verify we got the correct remaining data
	require.Equal(t, g.All()[7000:], remaining)
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

func TestPagedReader_SubReader(t *testing.T) {
	g := generate(10000)
	r := NewPagedReader(g)

	// Test reading a sub-section from the beginning
	sub1, err := r.SubReader(0, 1000)
	require.NoError(t, err)
	result1, err := io.ReadAll(sub1)
	require.NoError(t, err)
	require.Equal(t, g.All()[:1000], result1)

	// Test reading a sub-section from the middle
	sub2, err := r.SubReader(2000, 1500)
	require.NoError(t, err)
	result2, err := io.ReadAll(sub2)
	require.NoError(t, err)
	require.Equal(t, g.All()[2000:3500], result2)

	// Test reading a sub-section near the end
	sub3, err := r.SubReader(8500, 1000)
	require.NoError(t, err)
	result3, err := io.ReadAll(sub3)
	require.NoError(t, err)
	require.Equal(t, g.All()[8500:9500], result3)
}

func TestPagedReader_SubReader_MultiplePages(t *testing.T) {
	// Test with data that spans multiple pages
	g := generate(50 * 1024) // 50KB across multiple pages
	r := NewPagedReader(g)

	// Read a sub-section that spans multiple pages
	sub, err := r.SubReader(10*1024, 30*1024)
	require.NoError(t, err)
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[10*1024:40*1024], result)
}

func TestPagedReader_SubReader_SmallChunks(t *testing.T) {
	g := generate(5000)
	r := NewPagedReader(g)

	// Test reading sub-reader in small chunks
	sub, err := r.SubReader(1000, 2000)
	require.NoError(t, err)
	var result []byte
	buf := make([]byte, 100) // Read in 100-byte chunks

	for {
		n, err := sub.Read(buf)
		result = append(result, buf[:n]...)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}

	require.Equal(t, g.All()[1000:3000], result)
}

func TestPagedReader_SubReader_ZeroSize(t *testing.T) {
	g := generate(1000)
	r := NewPagedReader(g)

	// Test zero-size sub-reader
	sub, err := r.SubReader(500, 0)
	require.NoError(t, err)
	buf := make([]byte, 100)
	n, err := sub.Read(buf)
	require.Equal(t, 0, n)
	require.Equal(t, io.EOF, err)
}

func TestPagedReader_SubReader_ZeroLengthRead(t *testing.T) {
	g := generate(1000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(100, 500)
	require.NoError(t, err)

	// Test zero-length read
	n, err := sub.Read(nil)
	require.Equal(t, 0, n)
	require.NoError(t, err)

	n, err = sub.Read([]byte{})
	require.Equal(t, 0, n)
	require.NoError(t, err)
}

func TestPagedReader_SubReader_ConcurrentReads(t *testing.T) {
	g := generate(20000)
	r := NewPagedReader(g)

	// Create multiple sub-readers for different sections
	sub1, err := r.SubReader(0, 5000)
	require.NoError(t, err)
	sub2, err := r.SubReader(5000, 7000)
	require.NoError(t, err)
	sub3, err := r.SubReader(12000, 3000)
	require.NoError(t, err)

	// Read from all sub-readers
	result1, err := io.ReadAll(sub1)
	require.NoError(t, err)
	require.Equal(t, g.All()[:5000], result1)

	result2, err := io.ReadAll(sub2)
	require.NoError(t, err)
	require.Equal(t, g.All()[5000:12000], result2)

	result3, err := io.ReadAll(sub3)
	require.NoError(t, err)
	require.Equal(t, g.All()[12000:15000], result3)
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

func TestPagedReader_SubReader_ExactPageBoundary(t *testing.T) {
	g := generate(32 * 1024) // Exactly 2 pages
	r := NewPagedReader(g)

	// Test sub-reader that starts exactly at page boundary
	sub, err := r.SubReader(16*1024, 8*1024)
	require.NoError(t, err)
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[16*1024:24*1024], result)
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

func TestPagedReader_Skip_CrossPageBoundary(t *testing.T) {
	g := generate(40 * 1024) // 40KB
	r := NewPagedReader(g)

	// Sub-reader that spans multiple pages
	sub, err := r.SubReader(5*1024, 25*1024)
	require.NoError(t, err)
	subData, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[5*1024:30*1024], subData)

	// Skip the data crossing page boundaries
	err = r.Skip(30 * 1024)
	require.NoError(t, err)

	// Main reader should continue from correct position
	remaining, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All()[30*1024:], remaining)
}

func TestPagedReader_Skip_ZeroBytes(t *testing.T) {
	g := generate(2000)
	r := NewPagedReader(g)

	// Skip zero bytes - should be no-op
	err := r.Skip(0)
	require.NoError(t, err)

	// Should be able to read all data
	data, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All(), data)
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

func TestPagedReader_Empty_InitialState(t *testing.T) {
	g := generate(1000)
	r := NewPagedReader(g)

	// Initially should not be empty (source has data)
	require.False(t, r.Empty())
}

func TestPagedReader_Empty_EmptySource(t *testing.T) {
	g := generate(0)
	r := NewPagedReader(g)

	// Should be empty immediately with empty source
	require.True(t, r.Empty())
}

func TestPagedReader_Empty_AfterFullConsumption(t *testing.T) {
	g := generate(1000)
	r := NewPagedReader(g)

	// Read all data
	data, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All(), data)

	// Should be empty after consuming all data
	require.True(t, r.Empty())
}

func TestPagedReader_Empty_WithBufferedData(t *testing.T) {
	g := generate(2000)
	r := NewPagedReader(g)

	// Should not be empty with buffered data
	require.False(t, r.Empty())

	// Read part of the buffered data
	buf := make([]byte, 500)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 500, n)

	// Should still not be empty
	require.False(t, r.Empty())

	// Read all data
	_, err = io.ReadAll(r)
	require.NoError(t, err)

	// Should be empty after consuming all data
	require.True(t, r.Empty())
}

func TestPagedReader_Empty_AfterSkip(t *testing.T) {
	g := generate(1500)
	r := NewPagedReader(g)

	// Skip all buffered data
	err := r.Skip(1500)
	require.NoError(t, err)

	// Should be empty after skipping all data
	require.True(t, r.Empty())
}

func TestPagedReader_Empty_PartialSkip(t *testing.T) {
	g := generate(2000)
	r := NewPagedReader(g)

	// Skip part of the data
	err := r.Skip(1000)
	require.NoError(t, err)

	// Should not be empty - still has data in source and buffer
	require.False(t, r.Empty())
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

func TestPagedReader_Empty_LargeData(t *testing.T) {
	g := generate(100 * 1024) // 100KB
	r := NewPagedReader(g)

	require.False(t, r.Empty())

	// Read in chunks
	buf := make([]byte, 16*1024)
	var read int
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		read += n
		if read >= 100*1024 {
			// if we fall exactly at the end: err==nil, yet it's empty
			break
		}
		require.NoError(t, err)
		require.False(t, r.Empty())
	}

	// Should be empty after consuming everything
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

func TestPagedReader_Empty_EdgeCaseSingleByte(t *testing.T) {
	g := generate(1)
	r := NewPagedReader(g)

	require.False(t, r.Empty())

	// Read the single byte
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1, n)

	require.True(t, r.Empty())
}

func TestPagedReader_Empty_CrossPageBoundaries(t *testing.T) {
	// Create data that will span multiple pages
	g := generate(50 * 1024) // 50KB across multiple pages
	r := NewPagedReader(g)

	require.False(t, r.Empty())

	// Read page by page
	buf := make([]byte, pageSize)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.True(t, n > 0)
	}

	require.True(t, r.Empty())
}

func TestSubReader_SubReader_Basic(t *testing.T) {
	g := generate(10000)
	r := NewPagedReader(g)

	// Create a parent sub-reader
	parentSub, err := r.SubReader(1000, 5000)
	require.NoError(t, err)

	// Create a child sub-reader within the parent
	childSub, err := parentSub.SubReader(500, 2000)
	require.NoError(t, err)

	// Read from the child sub-reader
	result, err := io.ReadAll(childSub)
	require.NoError(t, err)
	require.Equal(t, g.All()[1500:3500], result) // 1000 + 500 = 1500, 1500 + 2000 = 3500
}

func TestSubReader_SubReader_ZeroOffset(t *testing.T) {
	g := generate(8000)
	r := NewPagedReader(g)

	// Create parent sub-reader
	parentSub, err := r.SubReader(2000, 4000)
	require.NoError(t, err)

	// Create child sub-reader with zero offset
	childSub, err := parentSub.SubReader(0, 1500)
	require.NoError(t, err)

	result, err := io.ReadAll(childSub)
	require.NoError(t, err)
	require.Equal(t, g.All()[2000:3500], result) // parent offset + child size
}

func TestSubReader_SubReader_ZeroSize(t *testing.T) {
	g := generate(5000)
	r := NewPagedReader(g)

	parentSub, err := r.SubReader(1000, 3000)
	require.NoError(t, err)

	// Create zero-size child sub-reader
	childSub, err := parentSub.SubReader(500, 0)
	require.NoError(t, err)

	buf := make([]byte, 100)
	n, err := childSub.Read(buf)
	require.Equal(t, 0, n)
	require.Equal(t, io.EOF, err)
}

func TestSubReader_SubReader_ExactBounds(t *testing.T) {
	g := generate(6000)
	r := NewPagedReader(g)

	// Create parent sub-reader
	parentSub, err := r.SubReader(1000, 3000)
	require.NoError(t, err)

	// Create child sub-reader that uses the entire parent range
	childSub, err := parentSub.SubReader(0, 3000)
	require.NoError(t, err)

	result, err := io.ReadAll(childSub)
	require.NoError(t, err)
	require.Equal(t, g.All()[1000:4000], result)
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

func TestSubReader_SubReader_MultipleLevels(t *testing.T) {
	g := generate(15000)
	r := NewPagedReader(g)

	// Level 1: Create parent sub-reader
	level1, err := r.SubReader(2000, 8000)
	require.NoError(t, err)

	// Level 2: Create child sub-reader
	level2, err := level1.SubReader(1000, 5000)
	require.NoError(t, err)

	// Level 3: Create grandchild sub-reader
	level3, err := level2.SubReader(500, 2000)
	require.NoError(t, err)

	// Read from the deepest level
	result, err := io.ReadAll(level3)
	require.NoError(t, err)
	// Final offset: 2000 + 1000 + 500 = 3500
	// Final range: 3500 to 5500
	require.Equal(t, g.All()[3500:5500], result)
}

func TestSubReader_SubReader_CrossPageBoundary(t *testing.T) {
	g := generate(50 * 1024) // 50KB across multiple pages
	r := NewPagedReader(g)

	// Create parent sub-reader that spans multiple pages
	parentSub, err := r.SubReader(10*1024, 30*1024)
	require.NoError(t, err)

	// Create child sub-reader that also spans multiple pages
	childSub, err := parentSub.SubReader(5*1024, 20*1024)
	require.NoError(t, err)

	result, err := io.ReadAll(childSub)
	require.NoError(t, err)
	// Parent starts at 10KB, child starts at +5KB = 15KB total offset
	// Child size is 20KB, so range is 15KB to 35KB
	require.Equal(t, g.All()[15*1024:35*1024], result)
}

func TestSubReader_SubReader_AfterPartialRead(t *testing.T) {
	g := generate(8000)
	r := NewPagedReader(g)

	parentSub, err := r.SubReader(1000, 5000)
	require.NoError(t, err)

	// Read part of the parent first
	buf := make([]byte, 1500)
	n, err := parentSub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1500, n)
	require.Equal(t, g.All()[1000:2500], buf)

	// Now create child sub-reader from remaining data
	childSub, err := parentSub.SubReader(500, 2000)
	require.NoError(t, err)

	// The child should read from the current position + offset
	// Parent current position is at 1500 (after reading 1500 bytes)
	// Child offset is 500, so child starts at position 1500 + 500 = 2000 in parent's coordinate system
	// In global coordinates: 1000 (parent start) + 2000 = 3000
	result, err := io.ReadAll(childSub)
	require.NoError(t, err)
	require.Equal(t, g.All()[3000:5000], result)
}

func TestSubReader_SubReader_SmallChunks(t *testing.T) {
	g := generate(6000)
	r := NewPagedReader(g)

	parentSub, err := r.SubReader(500, 4000)
	require.NoError(t, err)

	childSub, err := parentSub.SubReader(1000, 2000)
	require.NoError(t, err)

	// Read child in small chunks
	var result []byte
	buf := make([]byte, 200) // Small chunks

	for {
		n, err := childSub.Read(buf)
		result = append(result, buf[:n]...)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}

	require.Equal(t, g.All()[1500:3500], result) // 500 + 1000 = 1500, 1500 + 2000 = 3500
}

func TestSubReader_SubReader_ConcurrentChildren(t *testing.T) {
	g := generate(12000)
	r := NewPagedReader(g)

	// Create parent sub-reader
	parentSub, err := r.SubReader(2000, 8000)
	require.NoError(t, err)

	// Create multiple non-overlapping child sub-readers
	child1, err := parentSub.SubReader(500, 1500)
	require.NoError(t, err)

	child2, err := parentSub.SubReader(2500, 2000)
	require.NoError(t, err)

	child3, err := parentSub.SubReader(5000, 1000)
	require.NoError(t, err)

	// Read from all children
	result1, err := io.ReadAll(child1)
	require.NoError(t, err)
	require.Equal(t, g.All()[2500:4000], result1) // 2000 + 500 = 2500, 2500 + 1500 = 4000

	result2, err := io.ReadAll(child2)
	require.NoError(t, err)
	require.Equal(t, g.All()[4500:6500], result2) // 2000 + 2500 = 4500, 4500 + 2000 = 6500

	result3, err := io.ReadAll(child3)
	require.NoError(t, err)
	require.Equal(t, g.All()[7000:8000], result3) // 2000 + 5000 = 7000, 7000 + 1000 = 8000
}

func TestSubReader_SubReader_ZeroLengthRead(t *testing.T) {
	g := generate(4000)
	r := NewPagedReader(g)

	parentSub, err := r.SubReader(1000, 2000)
	require.NoError(t, err)

	childSub, err := parentSub.SubReader(500, 1000)
	require.NoError(t, err)

	// Test zero-length reads
	n, err := childSub.Read(nil)
	require.Equal(t, 0, n)
	require.NoError(t, err)

	n, err = childSub.Read([]byte{})
	require.Equal(t, 0, n)
	require.NoError(t, err)

	// Verify normal read still works
	buf := make([]byte, 500)
	n, err = childSub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 500, n)
	require.Equal(t, g.All()[1500:2000], buf) // 1000 + 500 = 1500
}

func TestSubReader_SubReader_EdgeCases(t *testing.T) {
	g := generate(3000)
	r := NewPagedReader(g)

	// Create parent sub-reader
	parentSub, err := r.SubReader(1000, 1000)
	require.NoError(t, err)

	// Edge case: child at very beginning
	child1, err := parentSub.SubReader(0, 1)
	require.NoError(t, err)
	result1, err := io.ReadAll(child1)
	require.NoError(t, err)
	require.Equal(t, g.All()[1000:1001], result1)

	// Edge case: child at very end
	child2, err := parentSub.SubReader(999, 1)
	require.NoError(t, err)
	result2, err := io.ReadAll(child2)
	require.NoError(t, err)
	require.Equal(t, g.All()[1999:2000], result2) // 1000 + 999 = 1999
}

func TestSubReader_SubReader_LargeData(t *testing.T) {
	g := generate(100 * 1024) // 100KB
	r := NewPagedReader(g)

	// Create large parent sub-reader
	parentSub, err := r.SubReader(10*1024, 80*1024)
	require.NoError(t, err)

	// Create large child sub-reader
	childSub, err := parentSub.SubReader(20*1024, 40*1024)
	require.NoError(t, err)

	// Read in page-sized chunks
	var result []byte
	buf := make([]byte, pageSize)

	for {
		n, err := childSub.Read(buf)
		result = append(result, buf[:n]...)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
	}

	// Verify the large read
	expectedStart := 10*1024 + 20*1024     // 30KB
	expectedEnd := expectedStart + 40*1024 // 70KB
	require.Equal(t, g.All()[expectedStart:expectedEnd], result)
}

func TestSubReader_Skip_Basic(t *testing.T) {
	g := generate(8000)
	r := NewPagedReader(g)

	// Create a sub-reader
	sub, err := r.SubReader(1000, 4000)
	require.NoError(t, err)

	// Read some data first
	buf := make([]byte, 1500)
	n, err := sub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1500, n)
	require.Equal(t, g.All()[1000:2500], buf)

	// Skip some bytes
	err = sub.Skip(1000)
	require.NoError(t, err)

	// Read remaining data
	remaining, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[3500:5000], remaining) // 1000 + 1500 + 1000 = 3500, 1000 + 4000 = 5000
}

func TestSubReader_Skip_ZeroBytes(t *testing.T) {
	g := generate(5000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(500, 2000)
	require.NoError(t, err)

	// Skip zero bytes - should be no-op
	err = sub.Skip(0)
	require.NoError(t, err)

	// Should be able to read all sub-reader data
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[500:2500], result)
}

func TestSubReader_Skip_ExactAvailable(t *testing.T) {
	g := generate(6000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(1000, 3000)
	require.NoError(t, err)

	// Skip all available bytes
	err = sub.Skip(3000)
	require.NoError(t, err)

	// Should be empty now
	buf := make([]byte, 100)
	n, err := sub.Read(buf)
	require.Equal(t, 0, n)
	require.Equal(t, io.EOF, err)
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

func TestSubReader_Skip_AfterPartialRead(t *testing.T) {
	g := generate(10000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(2000, 5000)
	require.NoError(t, err)

	// Read part of the data
	buf := make([]byte, 2000)
	n, err := sub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 2000, n)
	require.Equal(t, g.All()[2000:4000], buf)

	// Skip remaining data
	err = sub.Skip(3000) // 5000 - 2000 = 3000 remaining
	require.NoError(t, err)

	// Should be empty now
	n, err = sub.Read(buf)
	require.Equal(t, 0, n)
	require.Equal(t, io.EOF, err)
}

func TestSubReader_Skip_CrossPageBoundary(t *testing.T) {
	g := generate(50 * 1024) // 50KB across multiple pages
	r := NewPagedReader(g)

	// Create sub-reader that spans multiple pages
	sub, err := r.SubReader(5*1024, 30*1024)
	require.NoError(t, err)

	// Read some data first
	buf := make([]byte, 8*1024)
	n, err := sub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 8*1024, n)

	// Skip across page boundaries
	err = sub.Skip(15 * 1024)
	require.NoError(t, err)

	// Read remaining data
	remaining, err := io.ReadAll(sub)
	require.NoError(t, err)
	// Should have 30KB - 8KB - 15KB = 7KB remaining
	require.Equal(t, 7*1024, len(remaining))
	require.Equal(t, g.All()[28*1024:35*1024], remaining)
}

func TestSubReader_Skip_IncrementalSkips(t *testing.T) {
	g := generate(6000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(1000, 3000)
	require.NoError(t, err)

	// Skip in small increments
	err = sub.Skip(500)
	require.NoError(t, err)

	err = sub.Skip(800)
	require.NoError(t, err)

	err = sub.Skip(700)
	require.NoError(t, err)

	// Read remaining data (3000 - 500 - 800 - 700 = 1000 remaining)
	remaining, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[3000:4000], remaining) // 1000 + 2000 = 3000, 3000 + 1000 = 4000
}

func TestSubReader_Skip_WithinSinglePage(t *testing.T) {
	g := generate(8000)
	r := NewPagedReader(g)

	// Create small sub-reader within single page
	sub, err := r.SubReader(1000, 2000)
	require.NoError(t, err)

	// Skip part of it
	err = sub.Skip(500)
	require.NoError(t, err)

	// Read remaining
	remaining, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[1500:3000], remaining)
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

func TestSubReader_Skip_AtPageBoundary(t *testing.T) {
	g := generate(40 * 1024) // 40KB
	r := NewPagedReader(g)

	// Create sub-reader that starts at page boundary
	sub, err := r.SubReader(16*1024, 20*1024)
	require.NoError(t, err)

	// Skip exactly one page worth
	err = sub.Skip(16 * 1024)
	require.NoError(t, err)

	// Read remaining
	remaining, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, 4*1024, len(remaining)) // 20KB - 16KB = 4KB
	require.Equal(t, g.All()[32*1024:36*1024], remaining)
}

func TestSubReader_Empty_InitialState(t *testing.T) {
	g := generate(5000)
	r := NewPagedReader(g)

	// Non-empty sub-reader
	sub, err := r.SubReader(1000, 2000)
	require.NoError(t, err)
	require.False(t, sub.Empty())

	// Zero-size sub-reader
	emptySub, err := r.SubReader(500, 0)
	require.NoError(t, err)
	require.True(t, emptySub.Empty())
}

func TestSubReader_Empty_AfterFullRead(t *testing.T) {
	g := generate(6000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(1500, 2000)
	require.NoError(t, err)

	// Initially not empty
	require.False(t, sub.Empty())

	// Read all data
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[1500:3500], result)

	// Should be empty now
	require.True(t, sub.Empty())
}

func TestSubReader_Empty_AfterPartialRead(t *testing.T) {
	g := generate(8000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(2000, 3000)
	require.NoError(t, err)

	// Read part of the data
	buf := make([]byte, 1000)
	n, err := sub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1000, n)

	// Should not be empty yet
	require.False(t, sub.Empty())

	// Read remaining data
	_, err = io.ReadAll(sub)
	require.NoError(t, err)

	// Should be empty now
	require.True(t, sub.Empty())
}

func TestSubReader_Empty_AfterFullSkip(t *testing.T) {
	g := generate(5000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(1000, 2500)
	require.NoError(t, err)

	// Initially not empty
	require.False(t, sub.Empty())

	// Skip all data
	err = sub.Skip(2500)
	require.NoError(t, err)

	// Should be empty now
	require.True(t, sub.Empty())
}

func TestSubReader_Empty_AfterPartialSkip(t *testing.T) {
	g := generate(7000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(500, 3000)
	require.NoError(t, err)

	// Skip part of the data
	err = sub.Skip(1500)
	require.NoError(t, err)

	// Should not be empty yet
	require.False(t, sub.Empty())

	// Skip remaining data
	err = sub.Skip(1500)
	require.NoError(t, err)

	// Should be empty now
	require.True(t, sub.Empty())
}

func TestSubReader_Empty_MixedOperations(t *testing.T) {
	g := generate(10000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(2000, 4000)
	require.NoError(t, err)

	require.False(t, sub.Empty())

	// Read some data
	buf := make([]byte, 1000)
	n, err := sub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1000, n)
	require.False(t, sub.Empty())

	// Skip some data
	err = sub.Skip(1500)
	require.NoError(t, err)
	require.False(t, sub.Empty())

	// Read remaining
	remaining, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, 1500, len(remaining)) // 4000 - 1000 - 1500 = 1500
	require.True(t, sub.Empty())
}

func TestSubReader_Empty_ZeroSizeSubReader(t *testing.T) {
	g := generate(3000)
	r := NewPagedReader(g)

	// Create zero-size sub-reader
	sub, err := r.SubReader(1000, 0)
	require.NoError(t, err)

	// Should be empty immediately
	require.True(t, sub.Empty())

	// Any read should return EOF
	buf := make([]byte, 100)
	n, err := sub.Read(buf)
	require.Equal(t, 0, n)
	require.Equal(t, io.EOF, err)

	// Should still be empty
	require.True(t, sub.Empty())
}

func TestSubReader_Empty_SingleByteSubReader(t *testing.T) {
	g := generate(4000)
	r := NewPagedReader(g)

	// Create single-byte sub-reader
	sub, err := r.SubReader(1500, 1)
	require.NoError(t, err)

	require.False(t, sub.Empty())

	// Read the single byte
	buf := make([]byte, 1)
	n, err := sub.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 1, n)
	require.Equal(t, g.All()[1500:1501], buf)

	// Should be empty now
	require.True(t, sub.Empty())
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

func TestSubReader_Empty_CrossPageBoundaries(t *testing.T) {
	g := generate(60 * 1024) // 60KB across multiple pages
	r := NewPagedReader(g)

	// Create sub-reader spanning multiple pages
	sub, err := r.SubReader(10*1024, 40*1024)
	require.NoError(t, err)

	require.False(t, sub.Empty())

	// Read in page-sized chunks
	buf := make([]byte, pageSize)
	for {
		n, err := sub.Read(buf)
		if err == io.EOF {
			break
		}
		require.NoError(t, err)
		require.True(t, n > 0)
		// Should not be empty until the last read
		if n == pageSize {
			require.False(t, sub.Empty())
		}
	}

	// Should be empty after consuming all data
	require.True(t, sub.Empty())
}

func TestSubReader_Empty_ConsistentAfterEOF(t *testing.T) {
	g := generate(3000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(500, 1000)
	require.NoError(t, err)

	// Read all data
	result, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[500:1500], result)

	// Should be consistently empty on multiple calls
	require.True(t, sub.Empty())
	require.True(t, sub.Empty())
	require.True(t, sub.Empty())

	// Any further reads should return EOF
	buf := make([]byte, 100)
	n, err := sub.Read(buf)
	require.Equal(t, 0, n)
	require.Equal(t, io.EOF, err)

	// Should still be empty
	require.True(t, sub.Empty())
}

func TestSubReader_Empty_AfterSkipAndRead(t *testing.T) {
	g := generate(6000)
	r := NewPagedReader(g)

	sub, err := r.SubReader(1000, 3000)
	require.NoError(t, err)

	// Skip most of the data
	err = sub.Skip(2800)
	require.NoError(t, err)
	require.False(t, sub.Empty()) // Still has 200 bytes

	// Read remaining bytes
	buf := make([]byte, 200)
	n, _ := sub.Read(buf)
	require.Equal(t, 200, n)
	require.Equal(t, g.All()[3800:4000], buf)

	// Should be empty now
	require.True(t, sub.Empty())
}
