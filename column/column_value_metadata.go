package column

import (
	"fmt"
	"io"
	"iter"

	"gotomerge/column/rle"
)

type ValueMetadataColumnIter = iter.Seq2[ValueMetadata, error]

func ReadValueMetadataColumn(r io.Reader) ValueMetadataColumnIter {
	return func(yield func(ValueMetadata, error) bool) {
		for nullableUint64, err := range rle.ReadUint64RLE(r) {
			if err != nil {
				yield(0, err)
				return
			}
			val, valid := nullableUint64.Value()
			if !valid {
				// TODO: I think that's correct, need to update the spec if true
				yield(0, fmt.Errorf("null value in value metadata column"))
				return
			}
			if !yield(ValueMetadata(val), nil) {
				return
			}
		}
	}
}
