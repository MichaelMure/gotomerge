package column

import (
	"fmt"
	"math"

	"gotomerge/column/rle"
	ioutil "gotomerge/utils/io"
)

// DeltaReader is a stateful reader for delta-encoded int64 columns.
type DeltaReader struct {
	r   *rle.Int64Reader
	acc int64
}

func NewDeltaReader(r ioutil.SubReader) *DeltaReader {
	return &DeltaReader{r: rle.NewInt64Reader(r)}
}

func (dr *DeltaReader) Next() (rle.NullableValue[int64], error) {
	nv, err := dr.r.Next()
	if err != nil {
		return nil, err
	}
	val, valid := nv.Value()
	if !valid {
		// null: acc stays unchanged, return null
		return rle.NewNullInt64(), nil
	}
	switch {
	case val == 0:
		// no change to acc
	case val > 0:
		if dr.acc > math.MaxInt64-val {
			return nil, fmt.Errorf("overflow in delta column")
		}
		dr.acc += val
	case val < 0:
		if dr.acc < math.MinInt64-val {
			return nil, fmt.Errorf("underflow in delta column")
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
