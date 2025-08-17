package ioutil

import (
	"bytes"
	"io"
	"unsafe"
)

func ReadBytesLimitedPrealloc(r io.Reader, size uint64) ([]byte, error) {
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

func ReadStringLimitedPrealloc(r io.Reader, size uint64) (string, error) {
	strBytes, err := ReadBytesLimitedPrealloc(r, size)
	if err != nil {
		return "", err
	}
	// zero-copy cast to string
	return unsafe.String(unsafe.SliceData(strBytes), len(strBytes)), nil
}
