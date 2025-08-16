package column

import (
	"bytes"
	"fmt"
	"iter"

	"gotomerge/column/rle"
)

type ValueMetadataColumn []byte

type ValueMetadataColumnIter = iter.Seq2[ValueMetadata, error]

func ValueMetadataColumnFromBytes(b []byte) ValueMetadataColumn {
	return ValueMetadataColumn(b)
}

func (vmc ValueMetadataColumn) Iter() ValueMetadataColumnIter {
	return func(yield func(ValueMetadata, error) bool) {
		for nullableUint64, err := range rle.ReadUint64RLE(bytes.NewReader(vmc)) {
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
