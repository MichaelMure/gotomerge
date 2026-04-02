package column

import (
	"fmt"
	"io"

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

// ValueMetadataWriter is a stateful encoder for value metadata columns.
type ValueMetadataWriter struct {
	w *rle.Writer[uint64]
}

func NewValueMetadataWriter(w io.Writer) *ValueMetadataWriter {
	return &ValueMetadataWriter{w: rle.NewUint64Writer(w)}
}

func (vw *ValueMetadataWriter) Append(m ValueMetadata) {
	vw.w.Append(rle.NewNullableUint64(uint64(m)))
}

func (vw *ValueMetadataWriter) Flush() error { return vw.w.Flush() }

func (vr *ValueMetadataReader) Fork() (*ValueMetadataReader, error) {
	forkedR, err := vr.r.Fork()
	if err != nil {
		return nil, err
	}
	return &ValueMetadataReader{r: forkedR}, nil
}
