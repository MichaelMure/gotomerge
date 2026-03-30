package column

import (
	"fmt"

	"gotomerge/column/rle"
	ioutil "gotomerge/utils/io"
)

// ValueMetadataReader is a stateful reader for value metadata columns.
type ValueMetadataReader struct {
	r *rle.Uint64Reader
}

func NewValueMetadataReader(r ioutil.SubReader) *ValueMetadataReader {
	return &ValueMetadataReader{r: rle.NewUint64Reader(r)}
}

func (vr *ValueMetadataReader) Next() (ValueMetadata, error) {
	nv, err := vr.r.Next()
	if err != nil {
		return 0, err
	}
	val, valid := nv.Value()
	if !valid {
		return 0, fmt.Errorf("null value in value metadata column")
	}
	return ValueMetadata(val), nil
}

func (vr *ValueMetadataReader) Fork() (*ValueMetadataReader, error) {
	forkedR, err := vr.r.Fork()
	if err != nil {
		return nil, err
	}
	return &ValueMetadataReader{r: forkedR}, nil
}
