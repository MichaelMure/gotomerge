package column

import (
	"fmt"
	"io"
	"math"

	"gotomerge/column/rle"
	ioutil "gotomerge/utils/io"
)

// DeltaReader is a stateful reader for delta-encoded int64 columns.
type DeltaReader struct {
	r   *rle.Int64Reader
	acc int64
}

func NewDeltaReader(r *ioutil.SubReader) *DeltaReader {
	return &DeltaReader{r: rle.NewInt64Reader(r)}
}

func (dr *DeltaReader) Next() (rle.NullableValue[int64], error) {
	nv, err := dr.r.Next()
	if err != nil {
		return rle.NullableValue[int64]{}, err
	}
	val, valid := nv.Value()
	if !valid {
		// null: acc stays unchanged, return null
		return rle.NewNullInt64(), nil
	}
	var zero rle.NullableValue[int64]
	switch {
	case val == 0:
		// no change to acc
	case val > 0:
		if dr.acc > math.MaxInt64-val {
			return zero, fmt.Errorf("overflow in delta column")
		}
		dr.acc += val
	case val < 0:
		if dr.acc < math.MinInt64-val {
			return zero, fmt.Errorf("underflow in delta column")
		}
		dr.acc += val
	}
	return rle.NewNullableInt64(dr.acc), nil
}

func (dr *DeltaReader) Fork() (*DeltaReader, error) {
	forkedR, err := dr.r.Fork()
	if err != nil {
		return nil, err
	}
	return &DeltaReader{r: forkedR, acc: dr.acc}, nil
}

// DeltaWriter is a stateful encoder for delta-encoded int64 columns (matching
// what DeltaReader decodes). Values are stored as deltas from the previous
// non-null value; the accumulator starts at 0.
type DeltaWriter struct {
	w   *rle.Writer[int64]
	acc int64
}

func NewDeltaWriter(w io.Writer) *DeltaWriter {
	return &DeltaWriter{w: rle.NewInt64Writer(w)}
}

func (dw *DeltaWriter) Append(nv rle.NullableValue[int64]) {
	if v, ok := nv.Value(); ok {
		dw.w.Append(rle.NewNullableInt64(v - dw.acc))
		dw.acc = v
	} else {
		dw.w.Append(rle.NewNullInt64())
	}
}

// Flush writes the final run and returns any accumulated error.
func (dw *DeltaWriter) Flush() error { return dw.w.Flush() }


