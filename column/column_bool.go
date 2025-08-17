package column

import (
	"io"
	"iter"

	"github.com/jcalabro/leb128"
)

type BooleanColumnIter = iter.Seq2[bool, error]

func ReadBooleanColumn(r io.Reader) BooleanColumnIter {
	return func(yield func(bool, error) bool) {
		var val bool
		for {
			count, err := leb128.DecodeU64(r)
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
