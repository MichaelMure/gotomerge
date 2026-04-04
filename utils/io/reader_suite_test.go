package ioutil

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// suiteSubReader runs the common SubReader contract tests against any
// implementation. newReader(data) must return a fresh reader positioned at
// the start of data with no bytes yet consumed.
func suiteSubReader(t *testing.T, newReader func(data []byte) SubReader) {
	t.Helper()

	// --- Read ---

	t.Run("Read_Basic", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := newReader(data)
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("Read_ZeroLengthBuffer", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3})
		n, err := r.Read([]byte{})
		require.NoError(t, err)
		require.Equal(t, 0, n)
	})

	t.Run("Read_EOF_OnEmpty", func(t *testing.T) {
		r := newReader([]byte{})
		buf := make([]byte, 4)
		n, err := r.Read(buf)
		require.Equal(t, 0, n)
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("Read_LargerBufferThanData", func(t *testing.T) {
		data := []byte{10, 20, 30}
		r := newReader(data)
		buf := make([]byte, 100)
		n, err := r.Read(buf)
		// Implementations may return (n, nil) or (n, io.EOF) on the last read.
		if err != nil {
			require.ErrorIs(t, err, io.EOF)
		}
		require.Equal(t, data, buf[:n])
	})

	t.Run("Read_MultipleChunks", func(t *testing.T) {
		data := make([]byte, 500)
		for i := range data {
			data[i] = byte(i % 251)
		}
		r := newReader(data)
		var got []byte
		buf := make([]byte, 70)
		for {
			n, err := r.Read(buf)
			got = append(got, buf[:n]...)
			if err == io.EOF {
				break
			}
			require.NoError(t, err)
		}
		require.Equal(t, data, got)
	})

	t.Run("Read_CrossPageBoundary", func(t *testing.T) {
		data := make([]byte, pageSize+100)
		for i := range data {
			data[i] = byte(i % 251)
		}
		r := newReader(data)
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("Read_AfterEOF", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3})
		_, err := io.ReadAll(r)
		require.NoError(t, err)
		// Subsequent reads must consistently return (0, io.EOF).
		n, err := r.Read(make([]byte, 4))
		require.Equal(t, 0, n)
		require.ErrorIs(t, err, io.EOF)
	})

	// --- SubReader ---

	t.Run("SubReader_Basic", func(t *testing.T) {
		data := []byte{10, 20, 30, 40, 50}
		r := newReader(data)
		sub, err := r.SubReader(1, 3)
		require.NoError(t, err)
		got, err := io.ReadAll(sub)
		require.NoError(t, err)
		require.Equal(t, data[1:4], got)
	})

	t.Run("SubReader_ZeroOffset", func(t *testing.T) {
		data := []byte{10, 20, 30, 40, 50}
		r := newReader(data)
		sub, err := r.SubReader(0, 3)
		require.NoError(t, err)
		got, err := io.ReadAll(sub)
		require.NoError(t, err)
		require.Equal(t, data[:3], got)
	})

	t.Run("SubReader_DoesNotAdvancePosition", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := newReader(data)
		_, err := r.SubReader(1, 2)
		require.NoError(t, err)
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("SubReader_ZeroSize", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3})
		sub, err := r.SubReader(0, 0)
		require.NoError(t, err)
		got, err := io.ReadAll(sub)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("SubReader_CrossPageBoundary", func(t *testing.T) {
		data := make([]byte, pageSize+200)
		for i := range data {
			data[i] = byte(i % 251)
		}
		r := newReader(data)
		sub, err := r.SubReader(0, uint64(len(data)))
		require.NoError(t, err)
		got, err := io.ReadAll(sub)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("SubReader_AfterPositionAdvanced", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		r := newReader(data)
		buf := make([]byte, 3)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		// SubReader offset is relative to current position.
		sub, err := r.SubReader(0, 3)
		require.NoError(t, err)
		got, err := io.ReadAll(sub)
		require.NoError(t, err)
		require.Equal(t, data[3:6], got)
	})

	t.Run("SubReader_Nested", func(t *testing.T) {
		data := make([]byte, 20)
		for i := range data {
			data[i] = byte(i + 1)
		}
		r := newReader(data)
		outer, err := r.SubReader(2, 15)
		require.NoError(t, err)
		inner, err := outer.SubReader(3, 5)
		require.NoError(t, err)
		got, err := io.ReadAll(inner)
		require.NoError(t, err)
		require.Equal(t, data[2+3:2+3+5], got)
	})

	t.Run("SubReader_MultipleFromSameReader", func(t *testing.T) {
		data := make([]byte, 20)
		for i := range data {
			data[i] = byte(i + 1)
		}
		r := newReader(data)
		s1, err := r.SubReader(0, 5)
		require.NoError(t, err)
		s2, err := r.SubReader(5, 5)
		require.NoError(t, err)
		got1, err := io.ReadAll(s1)
		require.NoError(t, err)
		require.Equal(t, data[0:5], got1)
		got2, err := io.ReadAll(s2)
		require.NoError(t, err)
		require.Equal(t, data[5:10], got2)
	})

	// --- SubReaderOffset ---

	t.Run("SubReaderOffset_ZeroOffset", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := newReader(data)
		fork, err := r.SubReaderOffset(0)
		require.NoError(t, err)
		got, err := io.ReadAll(fork)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("SubReaderOffset_NonZero", func(t *testing.T) {
		data := []byte{10, 20, 30, 40, 50}
		r := newReader(data)
		fork, err := r.SubReaderOffset(2)
		require.NoError(t, err)
		got, err := io.ReadAll(fork)
		require.NoError(t, err)
		require.Equal(t, data[2:], got)
	})

	t.Run("SubReaderOffset_DoesNotAdvancePosition", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6}
		r := newReader(data)
		_, err := r.SubReaderOffset(3)
		require.NoError(t, err)
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("SubReaderOffset_CrossPageBoundary", func(t *testing.T) {
		data := make([]byte, pageSize+200)
		for i := range data {
			data[i] = byte(i % 251)
		}
		r := newReader(data)
		offset := uint64(pageSize + 50)
		fork, err := r.SubReaderOffset(offset)
		require.NoError(t, err)
		got, err := io.ReadAll(fork)
		require.NoError(t, err)
		require.Equal(t, data[offset:], got)
	})

	t.Run("SubReaderOffset_AfterPositionAdvanced", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		r := newReader(data)
		buf := make([]byte, 3)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		// Fork relative to current position with an additional offset.
		fork, err := r.SubReaderOffset(1)
		require.NoError(t, err)
		got, err := io.ReadAll(fork)
		require.NoError(t, err)
		require.Equal(t, data[4:], got)
	})

	// --- Skip ---

	t.Run("Skip_Basic", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		r := newReader(data)
		require.NoError(t, r.Skip(3))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[3:], got)
	})

	t.Run("Skip_Zero", func(t *testing.T) {
		data := []byte{1, 2, 3}
		r := newReader(data)
		require.NoError(t, r.Skip(0))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("Skip_All", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := newReader(data)
		require.NoError(t, r.Skip(len(data)))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("Skip_CrossPageBoundary", func(t *testing.T) {
		data := make([]byte, pageSize+300)
		for i := range data {
			data[i] = byte(i % 251)
		}
		r := newReader(data)
		skip := pageSize + 50
		require.NoError(t, r.Skip(skip))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[skip:], got)
	})

	t.Run("Skip_ExactPageBoundary", func(t *testing.T) {
		data := make([]byte, pageSize*2)
		for i := range data {
			data[i] = byte(i % 251)
		}
		r := newReader(data)
		require.NoError(t, r.Skip(pageSize))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[pageSize:], got)
	})

	t.Run("Skip_AfterRead", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6}
		r := newReader(data)
		buf := make([]byte, 2)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		require.NoError(t, r.Skip(2))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[4:], got)
	})

	t.Run("Skip_Incremental", func(t *testing.T) {
		data := make([]byte, 100)
		for i := range data {
			data[i] = byte(i)
		}
		r := newReader(data)
		require.NoError(t, r.Skip(10))
		require.NoError(t, r.Skip(20))
		require.NoError(t, r.Skip(5))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[35:], got)
	})

	// --- Consumed ---

	t.Run("Consumed_InitialState", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3})
		require.Equal(t, 0, r.Consumed())
	})

	t.Run("Consumed_AfterRead", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3, 4, 5})
		buf := make([]byte, 3)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		require.Equal(t, 3, r.Consumed())
	})

	t.Run("Consumed_AfterSkip", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3, 4, 5})
		require.NoError(t, r.Skip(4))
		require.Equal(t, 4, r.Consumed())
	})

	t.Run("Consumed_ReadAndSkip", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3, 4, 5, 6})
		buf := make([]byte, 2)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		require.NoError(t, r.Skip(2))
		require.Equal(t, 4, r.Consumed())
	})

	t.Run("Consumed_MultipleReads", func(t *testing.T) {
		data := make([]byte, 100)
		r := newReader(data)
		buf := make([]byte, 30)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		_, err = io.ReadFull(r, buf)
		require.NoError(t, err)
		require.Equal(t, 60, r.Consumed())
	})

	// --- Deflate ---

	t.Run("Deflate_Basic", func(t *testing.T) {
		plain := []byte("hello deflate world — suite")
		compressed := deflateCompress(t, plain)
		r := newReader(compressed)
		got, n, err := r.Deflate()
		require.NoError(t, err)
		require.Equal(t, len(plain), n)
		result, err := io.ReadAll(got)
		require.NoError(t, err)
		require.Equal(t, plain, result)
	})

	t.Run("Deflate_DoesNotAdvancePosition", func(t *testing.T) {
		plain := []byte("hello deflate world — suite")
		compressed := deflateCompress(t, plain)
		r := newReader(compressed)
		_, n, err := r.Deflate()
		require.NoError(t, err)
		require.Equal(t, len(plain), n)
		remaining, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, compressed, remaining)
	})

	t.Run("Deflate_Empty", func(t *testing.T) {
		compressed := deflateCompress(t, []byte{})
		r := newReader(compressed)
		got, n, err := r.Deflate()
		require.NoError(t, err)
		require.Equal(t, 0, n)
		result, err := io.ReadAll(got)
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("Deflate_AfterPositionAdvanced", func(t *testing.T) {
		plain := []byte("hello after advance")
		compressed := deflateCompress(t, plain)
		prefix := []byte{0xAA, 0xBB, 0xCC, 0xDD}
		data := append(prefix, compressed...)
		r := newReader(data)
		buf := make([]byte, len(prefix))
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		got, n, err := r.Deflate()
		require.NoError(t, err)
		require.Equal(t, len(plain), n)
		result, err := io.ReadAll(got)
		require.NoError(t, err)
		require.Equal(t, plain, result)
	})

	// --- Empty ---

	t.Run("Empty_InitialState_NonEmpty", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3})
		require.False(t, r.Empty())
	})

	t.Run("Empty_InitialState_Empty", func(t *testing.T) {
		r := newReader([]byte{})
		require.True(t, r.Empty())
	})

	t.Run("Empty_AfterFullRead", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3})
		_, err := io.ReadAll(r)
		require.NoError(t, err)
		require.True(t, r.Empty())
	})

	t.Run("Empty_AfterSkipAll", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := newReader(data)
		require.NoError(t, r.Skip(len(data)))
		require.True(t, r.Empty())
	})

	t.Run("Empty_FalseAfterPartialRead", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3, 4, 5})
		buf := make([]byte, 2)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		require.False(t, r.Empty())
	})

	t.Run("Empty_CrossPageBoundary", func(t *testing.T) {
		data := make([]byte, pageSize+50)
		for i := range data {
			data[i] = byte(i % 251)
		}
		r := newReader(data)
		require.False(t, r.Empty())
		_, err := io.ReadAll(r)
		require.NoError(t, err)
		require.True(t, r.Empty())
	})

	t.Run("Empty_FalseAfterPartialSkip", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3, 4, 5})
		require.NoError(t, r.Skip(3))
		require.False(t, r.Empty())
	})

	t.Run("Empty_AfterMixedReadAndSkip", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6}
		r := newReader(data)
		_, err := io.ReadFull(r, make([]byte, 2))
		require.NoError(t, err)
		require.NoError(t, r.Skip(4))
		require.True(t, r.Empty())
	})

	t.Run("Empty_SingleByte", func(t *testing.T) {
		r := newReader([]byte{42})
		require.False(t, r.Empty())
		_, err := io.ReadAll(r)
		require.NoError(t, err)
		require.True(t, r.Empty())
	})

	// --- Peek ---

	t.Run("Peek_Basic", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		r := newReader(data)
		var buf bytes.Buffer
		require.NoError(t, r.Peek(&buf, 4))
		require.Equal(t, data[:4], buf.Bytes())
		// position unchanged
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("Peek_Zero", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3})
		var buf bytes.Buffer
		require.NoError(t, r.Peek(&buf, 0))
		require.Equal(t, 0, buf.Len())
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, []byte{1, 2, 3}, got)
	})

	t.Run("Peek_AfterRead", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		r := newReader(data)
		_, err := io.ReadFull(r, make([]byte, 3))
		require.NoError(t, err)
		var buf bytes.Buffer
		require.NoError(t, r.Peek(&buf, 3))
		require.Equal(t, data[3:6], buf.Bytes())
		// position still at 3
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[3:], got)
	})

	t.Run("Peek_AfterSkip", func(t *testing.T) {
		data := []byte{10, 20, 30, 40, 50, 60}
		r := newReader(data)
		require.NoError(t, r.Skip(2))
		var buf bytes.Buffer
		require.NoError(t, r.Peek(&buf, 2))
		require.Equal(t, data[2:4], buf.Bytes())
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[2:], got)
	})

	t.Run("Peek_CrossPageBoundary", func(t *testing.T) {
		data := make([]byte, pageSize+100)
		for i := range data {
			data[i] = byte(i % 251)
		}
		r := newReader(data)
		var buf bytes.Buffer
		require.NoError(t, r.Peek(&buf, pageSize+50))
		require.Equal(t, data[:pageSize+50], buf.Bytes())
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	// --- HasAtLeast ---

	t.Run("HasAtLeast_Zero_AlwaysTrue", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3})
		require.True(t, r.HasAtLeast(0))
	})

	t.Run("HasAtLeast_Zero_OnEmpty", func(t *testing.T) {
		r := newReader([]byte{})
		require.True(t, r.HasAtLeast(0))
	})

	t.Run("HasAtLeast_ExactSize_True", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := newReader(data)
		require.True(t, r.HasAtLeast(len(data)))
	})

	t.Run("HasAtLeast_OneBeyondSize_False", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := newReader(data)
		require.False(t, r.HasAtLeast(len(data)+1))
	})

	t.Run("HasAtLeast_One_OnSingleByte", func(t *testing.T) {
		r := newReader([]byte{42})
		require.True(t, r.HasAtLeast(1))
	})

	t.Run("HasAtLeast_One_OnEmpty", func(t *testing.T) {
		r := newReader([]byte{})
		require.False(t, r.HasAtLeast(1))
	})

	t.Run("HasAtLeast_DoesNotAdvancePosition", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := newReader(data)
		require.True(t, r.HasAtLeast(3))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("HasAtLeast_AfterPartialRead", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6}
		r := newReader(data)
		_, err := io.ReadFull(r, make([]byte, 3))
		require.NoError(t, err)
		require.True(t, r.HasAtLeast(3))
		require.False(t, r.HasAtLeast(4))
	})

	t.Run("HasAtLeast_AfterFullRead_False", func(t *testing.T) {
		r := newReader([]byte{1, 2, 3})
		_, err := io.ReadAll(r)
		require.NoError(t, err)
		require.False(t, r.HasAtLeast(1))
	})

	t.Run("HasAtLeast_CrossPageBoundary", func(t *testing.T) {
		data := make([]byte, pageSize+50)
		for i := range data {
			data[i] = byte(i % 251)
		}
		r := newReader(data)
		require.True(t, r.HasAtLeast(pageSize+50))
		require.False(t, r.HasAtLeast(pageSize+51))
		// position must be unchanged
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})
}
