package ioutil

import (
	"compress/flate"
	"fmt"
	"io"
)

// SubReader is a zero-copy reader over a []byte. All sub-reader views share
// the same backing slice; no copies are made on fork or sub-view creation.
type SubReader struct {
	data     []byte
	position int
	consumed int // total bytes advanced via Read() or Skip()
}

// NewSubReader creates a SubReader over data. data is not copied.
func NewSubReader(data []byte) *SubReader {
	return &SubReader{data: data}
}

// ReadFrom reads all bytes from r and returns a SubReader over them.
func ReadFrom(r io.Reader) (*SubReader, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewSubReader(data), nil
}

func (b *SubReader) Read(p []byte) (n int, err error) {
	var toRead int
	if len(p) <= len(b.data)-b.position {
		toRead = len(p)
	} else {
		toRead = len(b.data) - b.position
		err = io.EOF
	}
	copy(p, b.data[b.position:b.position+toRead])
	b.position += toRead
	b.consumed += toRead
	return toRead, err
}

// ReadByte implements io.ByteReader. Returns io.EOF if no bytes remain.
func (b *SubReader) ReadByte() (byte, error) {
	if b.position >= len(b.data) {
		return 0, io.EOF
	}
	v := b.data[b.position]
	b.position++
	b.consumed++
	return v, nil
}

// SubReader returns a new SubReader covering [current+offset, current+offset+size).
// This does NOT advance the current position.
func (b *SubReader) SubReader(offset uint64, size uint64) (*SubReader, error) {
	if uint64(b.position)+offset+size > uint64(len(b.data)) {
		return nil, fmt.Errorf("SubReader: offset + size > len(data)")
	}
	start := uint64(b.position) + offset
	return &SubReader{data: b.data[start : start+size]}, nil
}

// SubReaderOffset returns a new SubReader starting at current+offset extending
// to the end of this reader's data. This does NOT advance the current position.
func (b *SubReader) SubReaderOffset(offset uint64) (*SubReader, error) {
	if uint64(b.position)+offset > uint64(len(b.data)) {
		return nil, fmt.Errorf("SubReader: offset > len(data)")
	}
	return &SubReader{data: b.data[uint64(b.position)+offset:]}, nil
}

// Skip advances the position by n bytes.
func (b *SubReader) Skip(n int) error {
	if b.position+n > len(b.data) {
		return fmt.Errorf("SubReader: skip %d > len(data)", n)
	}
	b.position += n
	b.consumed += n
	return nil
}

// Empty reports whether all bytes have been consumed.
func (b *SubReader) Empty() bool { return b.position >= len(b.data) }

// Consumed returns the total bytes advanced via Read or Skip since creation.
func (b *SubReader) Consumed() int { return b.consumed }

// Available returns the number of bytes remaining to be read.
func (b *SubReader) Available() int { return len(b.data) - b.position }

// HasAtLeast reports whether at least n bytes remain.
func (b *SubReader) HasAtLeast(n int) bool { return len(b.data)-b.position >= n }

// Peek writes the next n bytes to w without advancing the position.
func (b *SubReader) Peek(w io.Writer, n int) error {
	if n > len(b.data)-b.position {
		return fmt.Errorf("SubReader: peek %d > available", n)
	}
	_, err := w.Write(b.data[b.position : b.position+n])
	return err
}

// Deflate decompresses the remaining bytes using DEFLATE and returns a new
// SubReader over the decompressed data and the decompressed size.
func (b *SubReader) Deflate() (*SubReader, int, error) {
	fork, err := b.SubReaderOffset(0)
	if err != nil {
		return nil, 0, fmt.Errorf("deflate: fork: %w", err)
	}
	fr := flate.NewReader(fork)
	defer fr.Close()
	decompressed, err := io.ReadAll(fr)
	if err != nil {
		return nil, 0, fmt.Errorf("deflate: %w", err)
	}
	return NewSubReader(decompressed), len(decompressed), nil
}
