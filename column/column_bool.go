package column

import (
	"io"

	"github.com/jcalabro/leb128"

	ioutil "gotomerge/utils/io"
)

// BoolReader is a stateful reader for bool columns using the alternating
// run-length encoding: each LEB128 uint64 gives the count of the current
// value, then the value toggles (starting at false).
type BoolReader struct {
	r         ioutil.SubReader
	remaining uint64
	val       bool
	done      bool
}

func NewBoolReader(r ioutil.SubReader) *BoolReader {
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
