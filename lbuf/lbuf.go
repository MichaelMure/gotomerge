package lbuf

import (
	"bufio"
	"bytes"
	"compress/flate"
	"io"
	"sync"
	"unsafe"
)

var _ io.Reader = &Reader{}

// Reader is a specialized stream reader combining multiple capabilities:
// - buffered (through bufio.Reader) to batch and reduce the number of reading syscall
// - support recursively limiting the number of bytes available (similar as io.LimitedReader)
// - support efficient reads of size-delimited []byte or string, while limiting the pre-allocation to avoid DOS
// - support chaining a "processor" on the output, like flate.Reader
// - use a sync.Pool to reuse reading buffers
type Reader struct {
	r       io.Reader
	buf     *bufio.Reader // buffered reader
	limiter *io.LimitedReader
}

func FromReader(r io.Reader) *Reader {
	rr := pool.Get().(*bufio.Reader)
	rr.Reset(r)
	return &Reader{r: rr, buf: rr}
}

func FromBytes(b []byte) *Reader {
	// we add a limiter so that Empty() can work ok
	l := io.LimitReader(bytes.NewReader(b), int64(len(b))).(*io.LimitedReader)
	return &Reader{r: l, limiter: l}
}

func (r *Reader) Release() {
	if r.buf == nil {
		return
	}
	pool.Put(r.buf)
	r.buf = nil
	r.r = nil
	r.limiter = nil
}

// Limit returns a Reader with the same underlying buffered reader,
// but limited to reading only n bytes. After those, io.EOF is returned.
func (r *Reader) Limit(n int64) *Reader {
	l := io.LimitReader(r.r, n).(*io.LimitedReader)
	return &Reader{r: l, limiter: l}
}

// AddProcessor adds an io.Reader wrapper (like flate.NewReader) on the output of Reader.
func (r *Reader) AddProcessor(fn func(io.Reader) io.Reader) *Reader {
	return &Reader{r: fn(r.r), limiter: r.limiter}
}

func (r *Reader) Deflate() *Reader {
	return &Reader{r: flate.NewReader(r.r), limiter: r.limiter}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	return r.r.Read(p)
}

func (r *Reader) Empty() bool {
	if r.limiter != nil {
		return r.limiter.N <= 0
	}
	if rr, ok := r.r.(*bytes.Reader); ok {
		return rr.Len() <= 0
	}
	if r.buf != nil {
		_, err := r.buf.Peek(1)
		return err == io.EOF
	}
	// Note: this is not entirely correct, in particular with deflate that may or
	// may not have buffered data, but we have no way to check that.
	return false
}

func (r *Reader) ReadBytesLimitedPrealloc(size uint64) ([]byte, error) {
	if size <= bytes.MinRead {
		// reading with this size would force the buffer to grow and realloc with bytes.Buffer.ReadFrom()
		buf := make([]byte, size)
		_, err := io.ReadFull(r, buf)
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
	_, err := buf.ReadFrom(io.LimitReader(r, int64(size)))
	if err != nil {
		return nil, err
	}
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
