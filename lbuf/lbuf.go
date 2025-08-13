package lbuf

import (
	"bufio"
	"bytes"
	"io"
	"math"
	"sync"
	"unsafe"
)

var _ io.Reader = &Reader{}
var _ io.ByteReader = &Reader{}

// Reader is a specialized stream reader combining multiple capabilities:
// - buffered (through bufio.Reader) to batch and reduce the number of reading syscall
// - support recursively limiting the number of bytes available (similar as io.LimitedReader)
// - support efficient reads of size-delimited []byte or string, while limiting the pre-allocation to avoid DOS
// - use a sync.Pool to reuse reading buffers
type Reader struct {
	buf *bufio.Reader // buffered reader
	N   int64         // max bytes remaining
}

func FromReader(r io.Reader) *Reader {
	rr := pool.Get().(*bufio.Reader)
	rr.Reset(r)
	return &Reader{buf: rr, N: math.MaxInt64}
}

func FromBytes(b []byte) *Reader {
	rr := pool.Get().(*bufio.Reader)
	rr.Reset(bytes.NewReader(b))
	return &Reader{buf: rr, N: int64(len(b))}
}

func (r *Reader) Release() {
	pool.Put(r.buf)
	r.buf = nil
}

// Limit returns a Reader with the same underlying buffered reader,
// but limited to reading only n bytes. After those, io.EOF is returned.
func (r *Reader) Limit(n int64) *Reader {
	return &Reader{buf: r.buf, N: n}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	if r.N <= 0 {
		return 0, io.EOF
	}
	if int64(len(p)) > r.N {
		p = p[0:r.N]
	}
	n, err = r.buf.Read(p)
	r.N -= int64(n)
	return
}

func (r *Reader) ReadByte() (byte, error) {
	if r.N <= 0 {
		return 0, io.EOF
	}
	b, err := r.buf.ReadByte()
	if err != nil {
		return 0, err
	}
	r.N--
	return b, nil
}

func (r *Reader) Empty() bool {
	if r.N <= 0 {
		return true
	}
	_, err := r.buf.Peek(1)
	return err == io.EOF
}

func (r *Reader) ReadBytesLimitedPrealloc(size uint64) ([]byte, error) {
	if size > math.MaxInt64 || int64(size) > r.N {
		return nil, io.EOF
	}

	if size <= bytes.MinRead {
		// reading with this size would force the buffer to grow and realloc with bytes.Buffer.ReadFrom()
		buf := make([]byte, size)
		n, err := io.ReadFull(r, buf)
		r.N -= int64(n)
		return buf, err
	}

	prealloc := size
	if prealloc > 10_000 {
		// limit the pre-allocation to 10kB to avoid DOS
		// the buffer will grow if there is actually more data to read
		// the downside is that it can create reallocation and copy for larger value.
		prealloc = 10_000
	}
	buf := bytes.NewBuffer(make([]byte, 0, prealloc))
	n, err := buf.ReadFrom(io.LimitReader(r, int64(size)))
	if err != nil {
		return nil, err
	}
	r.N -= n
	return buf.Bytes(), nil
}

func (r *Reader) ReadStringLimitedPrealloc(size uint64) (string, error) {
	strBytes, err := r.ReadBytesLimitedPrealloc(size)
	if err != nil {
		return "", err
	}
	// zero-copy cast to string
	return unsafe.String(unsafe.SliceData(strBytes), len(strBytes)), nil
}

var pool = sync.Pool{
	New: func() interface{} {
		// size from a quick bench tuning
		return bufio.NewReaderSize(nil, 16*1024) // 16kB
	},
}
