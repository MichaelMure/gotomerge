package ioutil

import (
	"bytes"
	"compress/flate"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

// -- helpers -----------------------------------------------------------------

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

// -- Read --------------------------------------------------------------------

func TestSubReader_Read(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := NewSubReader(data)
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("ZeroLengthBuffer", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3})
		n, err := r.Read([]byte{})
		require.NoError(t, err)
		require.Equal(t, 0, n)
	})

	t.Run("EOF_OnEmpty", func(t *testing.T) {
		r := NewSubReader([]byte{})
		buf := make([]byte, 4)
		n, err := r.Read(buf)
		require.Equal(t, 0, n)
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("LargerBufferThanData", func(t *testing.T) {
		data := []byte{10, 20, 30}
		r := NewSubReader(data)
		buf := make([]byte, 100)
		n, err := r.Read(buf)
		if err != nil {
			require.ErrorIs(t, err, io.EOF)
		}
		require.Equal(t, data, buf[:n])
	})

	t.Run("MultipleChunks", func(t *testing.T) {
		data := make([]byte, 500)
		for i := range data {
			data[i] = byte(i % 251)
		}
		r := NewSubReader(data)
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

	t.Run("AfterEOF", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3})
		_, err := io.ReadAll(r)
		require.NoError(t, err)
		n, err := r.Read(make([]byte, 4))
		require.Equal(t, 0, n)
		require.ErrorIs(t, err, io.EOF)
	})
}

// -- ReadByte ----------------------------------------------------------------

func TestSubReader_ReadByte(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		data := []byte{10, 20, 30}
		r := NewSubReader(data)
		b, err := r.ReadByte()
		require.NoError(t, err)
		require.Equal(t, byte(10), b)
		b, err = r.ReadByte()
		require.NoError(t, err)
		require.Equal(t, byte(20), b)
		b, err = r.ReadByte()
		require.NoError(t, err)
		require.Equal(t, byte(30), b)
	})

	t.Run("EOF_OnEmpty", func(t *testing.T) {
		r := NewSubReader([]byte{})
		_, err := r.ReadByte()
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("EOF_AfterAllConsumed", func(t *testing.T) {
		r := NewSubReader([]byte{42})
		b, err := r.ReadByte()
		require.NoError(t, err)
		require.Equal(t, byte(42), b)
		_, err = r.ReadByte()
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("AllByteValues", func(t *testing.T) {
		data := make([]byte, 256)
		for i := range data {
			data[i] = byte(i)
		}
		r := NewSubReader(data)
		for i := range data {
			b, err := r.ReadByte()
			require.NoError(t, err)
			require.Equal(t, byte(i), b)
		}
		_, err := r.ReadByte()
		require.ErrorIs(t, err, io.EOF)
	})

	t.Run("ConsumedIncrement", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3, 4, 5})
		require.Equal(t, 0, r.Consumed())
		_, _ = r.ReadByte()
		require.Equal(t, 1, r.Consumed())
		_, _ = r.ReadByte()
		_, _ = r.ReadByte()
		require.Equal(t, 3, r.Consumed())
	})

	t.Run("MixedWithRead", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6}
		r := NewSubReader(data)
		// Read 2 bytes via Read
		buf := make([]byte, 2)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		require.Equal(t, data[:2], buf)
		// Read next byte via ReadByte
		b, err := r.ReadByte()
		require.NoError(t, err)
		require.Equal(t, byte(3), b)
		// Read remainder via Read
		rest, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[3:], rest)
	})

	t.Run("ImplementsByteReader", func(t *testing.T) {
		// *SubReader must satisfy io.ByteReader; this exercises the interface.
		data := []byte{0x01, 0x80, 0x01} // LEB128 encoding of 128
		r := NewSubReader(data)
		var br io.ByteReader = r
		b, err := br.ReadByte()
		require.NoError(t, err)
		require.Equal(t, byte(0x01), b)
	})

	t.Run("AfterSkip", func(t *testing.T) {
		data := []byte{10, 20, 30, 40, 50}
		r := NewSubReader(data)
		require.NoError(t, r.Skip(2))
		b, err := r.ReadByte()
		require.NoError(t, err)
		require.Equal(t, byte(30), b)
		require.Equal(t, 3, r.Consumed())
	})

	t.Run("SubReader_ReadByte", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6}
		r := NewSubReader(data)
		sub, err := r.SubReader(2, 3) // data[2:5] = {3, 4, 5}
		require.NoError(t, err)
		b, err := sub.ReadByte()
		require.NoError(t, err)
		require.Equal(t, byte(3), b)
		b, err = sub.ReadByte()
		require.NoError(t, err)
		require.Equal(t, byte(4), b)
		b, err = sub.ReadByte()
		require.NoError(t, err)
		require.Equal(t, byte(5), b)
		_, err = sub.ReadByte()
		require.ErrorIs(t, err, io.EOF)
	})
}

// -- SubReader ---------------------------------------------------------------

func TestSubReader_SubReader(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		data := []byte{10, 20, 30, 40, 50}
		r := NewSubReader(data)
		sub, err := r.SubReader(1, 3)
		require.NoError(t, err)
		got, err := io.ReadAll(sub)
		require.NoError(t, err)
		require.Equal(t, data[1:4], got)
	})

	t.Run("ZeroOffset", func(t *testing.T) {
		data := []byte{10, 20, 30, 40, 50}
		r := NewSubReader(data)
		sub, err := r.SubReader(0, 3)
		require.NoError(t, err)
		got, err := io.ReadAll(sub)
		require.NoError(t, err)
		require.Equal(t, data[:3], got)
	})

	t.Run("ZeroSize", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3})
		sub, err := r.SubReader(0, 0)
		require.NoError(t, err)
		got, err := io.ReadAll(sub)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("DoesNotAdvancePosition", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := NewSubReader(data)
		_, err := r.SubReader(1, 2)
		require.NoError(t, err)
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("AfterPositionAdvanced", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		r := NewSubReader(data)
		buf := make([]byte, 3)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		sub, err := r.SubReader(0, 3)
		require.NoError(t, err)
		got, err := io.ReadAll(sub)
		require.NoError(t, err)
		require.Equal(t, data[3:6], got)
	})

	t.Run("Nested", func(t *testing.T) {
		data := make([]byte, 20)
		for i := range data {
			data[i] = byte(i + 1)
		}
		r := NewSubReader(data)
		outer, err := r.SubReader(2, 15)
		require.NoError(t, err)
		inner, err := outer.SubReader(3, 5)
		require.NoError(t, err)
		got, err := io.ReadAll(inner)
		require.NoError(t, err)
		require.Equal(t, data[2+3:2+3+5], got)
	})

	t.Run("MultipleFromSameReader", func(t *testing.T) {
		data := make([]byte, 20)
		for i := range data {
			data[i] = byte(i + 1)
		}
		r := NewSubReader(data)
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

	t.Run("OutOfBounds", func(t *testing.T) {
		r := NewSubReader(generate(80).All())
		_, err := r.SubReader(70, 20) // 70 + 20 = 90 > 80
		require.Error(t, err)
		require.Contains(t, err.Error(), "offset + size > len(data)")

		// After reading some data, offset is relative to position.
		buf := make([]byte, 30)
		_, err = r.Read(buf)
		require.NoError(t, err)
		_, err = r.SubReader(30, 25) // position(30) + offset(30) + size(25) = 85 > 80
		require.Error(t, err)
		require.Contains(t, err.Error(), "offset + size > len(data)")
	})
}

// -- SubReaderOffset ---------------------------------------------------------

func TestSubReader_SubReaderOffset(t *testing.T) {
	t.Run("ZeroOffset", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := NewSubReader(data)
		fork, err := r.SubReaderOffset(0)
		require.NoError(t, err)
		got, err := io.ReadAll(fork)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("NonZero", func(t *testing.T) {
		data := []byte{10, 20, 30, 40, 50}
		r := NewSubReader(data)
		fork, err := r.SubReaderOffset(2)
		require.NoError(t, err)
		got, err := io.ReadAll(fork)
		require.NoError(t, err)
		require.Equal(t, data[2:], got)
	})

	t.Run("DoesNotAdvancePosition", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6}
		r := NewSubReader(data)
		_, err := r.SubReaderOffset(3)
		require.NoError(t, err)
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("AfterPositionAdvanced", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		r := NewSubReader(data)
		buf := make([]byte, 3)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		fork, err := r.SubReaderOffset(1)
		require.NoError(t, err)
		got, err := io.ReadAll(fork)
		require.NoError(t, err)
		require.Equal(t, data[4:], got)
	})

	t.Run("OutOfBounds", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := NewSubReader(data)
		_, err := r.SubReaderOffset(6)
		require.Error(t, err)
	})
}

// -- Skip --------------------------------------------------------------------

func TestSubReader_Skip(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		r := NewSubReader(data)
		require.NoError(t, r.Skip(3))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[3:], got)
	})

	t.Run("Zero", func(t *testing.T) {
		data := []byte{1, 2, 3}
		r := NewSubReader(data)
		require.NoError(t, r.Skip(0))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("All", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := NewSubReader(data)
		require.NoError(t, r.Skip(len(data)))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Empty(t, got)
	})

	t.Run("AfterRead", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6}
		r := NewSubReader(data)
		buf := make([]byte, 2)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		require.NoError(t, r.Skip(2))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[4:], got)
	})

	t.Run("Incremental", func(t *testing.T) {
		data := make([]byte, 100)
		for i := range data {
			data[i] = byte(i)
		}
		r := NewSubReader(data)
		require.NoError(t, r.Skip(10))
		require.NoError(t, r.Skip(20))
		require.NoError(t, r.Skip(5))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[35:], got)
	})

	t.Run("OutOfBounds", func(t *testing.T) {
		r := NewSubReader(generate(30).All())
		err := r.Skip(40)
		require.Error(t, err)
		require.Contains(t, err.Error(), "skip 40 > len(data)")
	})
}

// -- Consumed ----------------------------------------------------------------

func TestSubReader_Consumed(t *testing.T) {
	t.Run("InitialState", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3})
		require.Equal(t, 0, r.Consumed())
	})

	t.Run("AfterRead", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3, 4, 5})
		buf := make([]byte, 3)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		require.Equal(t, 3, r.Consumed())
	})

	t.Run("AfterSkip", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3, 4, 5})
		require.NoError(t, r.Skip(4))
		require.Equal(t, 4, r.Consumed())
	})

	t.Run("ReadAndSkip", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3, 4, 5, 6})
		buf := make([]byte, 2)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		require.NoError(t, r.Skip(2))
		require.Equal(t, 4, r.Consumed())
	})

	t.Run("MultipleReads", func(t *testing.T) {
		data := make([]byte, 100)
		r := NewSubReader(data)
		buf := make([]byte, 30)
		_, err := io.ReadFull(r, buf)
		require.NoError(t, err)
		_, err = io.ReadFull(r, buf)
		require.NoError(t, err)
		require.Equal(t, 60, r.Consumed())
	})
}

// -- Deflate -----------------------------------------------------------------

func TestSubReader_Deflate(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		plain := []byte("hello deflate world")
		compressed := deflateCompress(t, plain)
		r := NewSubReader(compressed)
		got, n, err := r.Deflate()
		require.NoError(t, err)
		require.Equal(t, len(plain), n)
		result, err := io.ReadAll(got)
		require.NoError(t, err)
		require.Equal(t, plain, result)
	})

	t.Run("DoesNotAdvancePosition", func(t *testing.T) {
		plain := []byte("hello deflate world")
		compressed := deflateCompress(t, plain)
		r := NewSubReader(compressed)
		_, n, err := r.Deflate()
		require.NoError(t, err)
		require.Equal(t, len(plain), n)
		remaining, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, compressed, remaining)
	})

	t.Run("Empty", func(t *testing.T) {
		compressed := deflateCompress(t, []byte{})
		r := NewSubReader(compressed)
		got, n, err := r.Deflate()
		require.NoError(t, err)
		require.Equal(t, 0, n)
		result, err := io.ReadAll(got)
		require.NoError(t, err)
		require.Empty(t, result)
	})

	t.Run("AfterPositionAdvanced", func(t *testing.T) {
		plain := []byte("hello after advance")
		compressed := deflateCompress(t, plain)
		prefix := []byte{0xAA, 0xBB, 0xCC, 0xDD}
		data := append(prefix, compressed...)
		r := NewSubReader(data)
		_, err := io.ReadFull(r, make([]byte, len(prefix)))
		require.NoError(t, err)
		got, n, err := r.Deflate()
		require.NoError(t, err)
		require.Equal(t, len(plain), n)
		result, err := io.ReadAll(got)
		require.NoError(t, err)
		require.Equal(t, plain, result)
	})

	t.Run("AfterPartialRead", func(t *testing.T) {
		plain := []byte("hello deflate world")
		compressed := deflateCompress(t, plain)
		data := append([]byte{1, 2, 3, 4}, compressed...)
		r := NewSubReader(data)
		_, err := io.ReadFull(r, make([]byte, 4))
		require.NoError(t, err)
		got, n, err := r.Deflate()
		require.NoError(t, err)
		require.Equal(t, len(plain), n)
		result, err := io.ReadAll(got)
		require.NoError(t, err)
		require.Equal(t, plain, result)
	})
}

// -- Empty -------------------------------------------------------------------

func TestSubReader_Empty(t *testing.T) {
	t.Run("InitialState_NonEmpty", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3})
		require.False(t, r.Empty())
	})

	t.Run("InitialState_Empty", func(t *testing.T) {
		r := NewSubReader([]byte{})
		require.True(t, r.Empty())
	})

	t.Run("AfterFullRead", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3})
		_, err := io.ReadAll(r)
		require.NoError(t, err)
		require.True(t, r.Empty())
	})

	t.Run("AfterSkipAll", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := NewSubReader(data)
		require.NoError(t, r.Skip(len(data)))
		require.True(t, r.Empty())
	})

	t.Run("FalseAfterPartialRead", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3, 4, 5})
		_, err := io.ReadFull(r, make([]byte, 2))
		require.NoError(t, err)
		require.False(t, r.Empty())
	})

	t.Run("FalseAfterPartialSkip", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3, 4, 5})
		require.NoError(t, r.Skip(3))
		require.False(t, r.Empty())
	})

	t.Run("AfterMixedReadAndSkip", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6}
		r := NewSubReader(data)
		_, err := io.ReadFull(r, make([]byte, 2))
		require.NoError(t, err)
		require.NoError(t, r.Skip(4))
		require.True(t, r.Empty())
	})

	t.Run("SingleByte", func(t *testing.T) {
		r := NewSubReader([]byte{42})
		require.False(t, r.Empty())
		_, err := io.ReadAll(r)
		require.NoError(t, err)
		require.True(t, r.Empty())
	})
}

// -- Peek --------------------------------------------------------------------

func TestSubReader_Peek(t *testing.T) {
	t.Run("Basic", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		r := NewSubReader(data)
		var buf bytes.Buffer
		require.NoError(t, r.Peek(&buf, 4))
		require.Equal(t, data[:4], buf.Bytes())
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("Zero", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3})
		var buf bytes.Buffer
		require.NoError(t, r.Peek(&buf, 0))
		require.Equal(t, 0, buf.Len())
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, []byte{1, 2, 3}, got)
	})

	t.Run("AfterRead", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6, 7, 8}
		r := NewSubReader(data)
		_, err := io.ReadFull(r, make([]byte, 3))
		require.NoError(t, err)
		var buf bytes.Buffer
		require.NoError(t, r.Peek(&buf, 3))
		require.Equal(t, data[3:6], buf.Bytes())
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[3:], got)
	})

	t.Run("AfterSkip", func(t *testing.T) {
		data := []byte{10, 20, 30, 40, 50, 60}
		r := NewSubReader(data)
		require.NoError(t, r.Skip(2))
		var buf bytes.Buffer
		require.NoError(t, r.Peek(&buf, 2))
		require.Equal(t, data[2:4], buf.Bytes())
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data[2:], got)
	})

	t.Run("OutOfBounds", func(t *testing.T) {
		r := NewSubReader(generate(30).All())
		var buf bytes.Buffer
		err := r.Peek(&buf, 50)
		require.Error(t, err)
		require.Contains(t, err.Error(), "peek 50 > available")
	})

	t.Run("WriterError", func(t *testing.T) {
		r := NewSubReader(generate(50).All())
		err := r.Peek(&failingWriter{failAfter: 10}, 25)
		require.Error(t, err)
		require.Contains(t, err.Error(), "write failed")
	})
}

// -- HasAtLeast --------------------------------------------------------------

func TestSubReader_HasAtLeast(t *testing.T) {
	t.Run("Zero_AlwaysTrue", func(t *testing.T) {
		require.True(t, NewSubReader([]byte{1, 2, 3}).HasAtLeast(0))
	})

	t.Run("Zero_OnEmpty", func(t *testing.T) {
		require.True(t, NewSubReader([]byte{}).HasAtLeast(0))
	})

	t.Run("ExactSize_True", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		require.True(t, NewSubReader(data).HasAtLeast(len(data)))
	})

	t.Run("OneBeyondSize_False", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		require.False(t, NewSubReader(data).HasAtLeast(len(data)+1))
	})

	t.Run("One_OnSingleByte", func(t *testing.T) {
		require.True(t, NewSubReader([]byte{42}).HasAtLeast(1))
	})

	t.Run("One_OnEmpty", func(t *testing.T) {
		require.False(t, NewSubReader([]byte{}).HasAtLeast(1))
	})

	t.Run("DoesNotAdvancePosition", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5}
		r := NewSubReader(data)
		require.True(t, r.HasAtLeast(3))
		got, err := io.ReadAll(r)
		require.NoError(t, err)
		require.Equal(t, data, got)
	})

	t.Run("AfterPartialRead", func(t *testing.T) {
		data := []byte{1, 2, 3, 4, 5, 6}
		r := NewSubReader(data)
		_, err := io.ReadFull(r, make([]byte, 3))
		require.NoError(t, err)
		require.True(t, r.HasAtLeast(3))
		require.False(t, r.HasAtLeast(4))
	})

	t.Run("AfterFullRead_False", func(t *testing.T) {
		r := NewSubReader([]byte{1, 2, 3})
		_, err := io.ReadAll(r)
		require.NoError(t, err)
		require.False(t, r.HasAtLeast(1))
	})
}

// -- Integration -------------------------------------------------------------

func TestSubReader_Integration(t *testing.T) {
	g := generate(500)
	r := NewSubReader(g.All())

	require.False(t, r.Empty())

	var peekBuf bytes.Buffer
	require.NoError(t, r.Peek(&peekBuf, 50))
	require.Equal(t, g.All()[:50], peekBuf.Bytes())

	buf := make([]byte, 50)
	n, err := r.Read(buf)
	require.NoError(t, err)
	require.Equal(t, 50, n)
	require.Equal(t, g.All()[:50], buf)

	require.NoError(t, r.Skip(100)) // skip bytes 50-149

	sub, err := r.SubReader(50, 80) // bytes 200-279 (position 150 + offset 50)
	require.NoError(t, err)
	subResult, err := io.ReadAll(sub)
	require.NoError(t, err)
	require.Equal(t, g.All()[200:280], subResult)

	remaining, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All()[150:], remaining)

	require.True(t, r.Empty())
}

func TestSubReader_FullCycle(t *testing.T) {
	g := generate(256)
	r := NewSubReader(g.All())

	var peekBuf bytes.Buffer
	require.NoError(t, r.Peek(&peekBuf, 100))
	require.Equal(t, g.All()[:100], peekBuf.Bytes())

	result, err := io.ReadAll(r)
	require.NoError(t, err)
	require.Equal(t, g.All(), result)

	for i := 0; i < 256; i++ {
		require.Equal(t, byte(i), result[i])
	}
}
