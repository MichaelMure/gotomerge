package ioutil

import "io"

var _ io.Reader = &teeReader{}
var _ io.ByteReader = &teeReader{}

type teeReader struct {
	r ByteReader
	w io.Writer
}

// TeeReader returns a [Reader] that writes to w what it reads from r.
// All reads from r performed through it are matched with
// corresponding writes to w. There is no internal buffering -
// the write must complete before the read completes.
// Any error encountered while writing is reported as a read error.
func TeeReader(r ByteReader, w io.Writer) ByteReader {
	return &teeReader{r, w}
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}

func (t *teeReader) ReadByte() (b byte, err error) {
	b, err = t.r.ReadByte()
	if err == nil {
		if _, err := t.w.Write([]byte{b}); err != nil {
			return 0, err
		}
	}
	return b, err
}

type ByteReader interface {
	io.Reader
	io.ByteReader
}
