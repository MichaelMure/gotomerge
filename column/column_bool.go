package column

import (
	"io"

	"github.com/MichaelMure/leb128"

	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

// BoolReader is a stateful reader for bool columns using the alternating
// run-length encoding: each LEB128 uint64 gives the count of the current
// value, then the value toggles (starting at false).
type BoolReader struct {
	r         *ioutil.SubReader
	remaining uint64
	val       bool
	done      bool
}

func NewBoolReader(r *ioutil.SubReader) *BoolReader {
	return &BoolReader{r: r}
}

func (b *BoolReader) Next() (bool, error) {
	for b.remaining == 0 {
		if b.done {
			return false, io.EOF
		}
		count, err := leb128.DecodeU64(b.r)
		if err == io.EOF {
			b.done = true
			return false, io.EOF
		}
		if err != nil {
			return false, err
		}
		b.remaining = count
		if b.remaining == 0 {
			// toggle even on zero-count run
			b.val = !b.val
		}
	}
	b.remaining--
	v := b.val
	if b.remaining == 0 {
		b.val = !b.val
	}
	return v, nil
}

func (b *BoolReader) Fork() (*BoolReader, error) {
	sub, err := b.r.SubReaderOffset(0)
	if err != nil {
		return nil, err
	}
	cp := *b
	cp.r = sub
	return &cp, nil
}

// BoolWriter is a stateful encoder for bool columns using the alternating-run
// encoding (matching what BoolReader decodes). The sequence starts at false.
type BoolWriter struct {
	w       io.Writer
	err     error
	current bool
	run     uint64
	started bool
}

func NewBoolWriter(w io.Writer) *BoolWriter { return &BoolWriter{w: w} }

func (bw *BoolWriter) Append(v bool) {
	if bw.err != nil {
		return
	}
	if !bw.started {
		bw.started = true
		// The sequence starts at false. If the first value is true, emit a
		// zero-length false run to advance the toggle to the right position.
		if v {
			_, bw.err = bw.w.Write(leb128.EncodeU64(0))
			bw.current = true
		}
		bw.run = 1
		return
	}
	if v == bw.current {
		bw.run++
	} else {
		_, bw.err = bw.w.Write(leb128.EncodeU64(bw.run))
		bw.current = v
		bw.run = 1
	}
}

// Flush writes the final run and returns any accumulated error.
func (bw *BoolWriter) Flush() error {
	if bw.err != nil {
		return bw.err
	}
	if bw.started {
		_, bw.err = bw.w.Write(leb128.EncodeU64(bw.run))
		bw.started = false
	}
	return bw.err
}
