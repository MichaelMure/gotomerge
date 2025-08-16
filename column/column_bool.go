package column

import (
	"bytes"
	"io"
	"iter"

	"github.com/jcalabro/leb128"
)

type BooleanColumn []byte

type BooleanColumnIter = iter.Seq2[bool, error]

func BooleanColumnFromBytes(b []byte) BooleanColumn {
	return BooleanColumn(b)
}

func (bc BooleanColumn) Iter() BooleanColumnIter {
	return func(yield func(bool, error) bool) {
		var val bool
		for {
			count, err := leb128.DecodeU64(bytes.NewReader(bc))
			if err == io.EOF {
				return
			}
			if err != nil {
				yield(false, err)
				return
			}
			for i := uint64(0); i < count; i++ {
				if !yield(val, nil) {
					return
				}
			}
			val = !val
		}
	}
}
