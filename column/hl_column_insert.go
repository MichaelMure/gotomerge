package column

// InsertReader is a stateful reader for insert (bool) columns.
type InsertReader struct {
	r *BoolReader
}

func NewInsertReader(r *BoolReader) *InsertReader {
	return &InsertReader{r: r}
}

func (i *InsertReader) Next() (bool, error) {
	if i.r == nil {
		return false, nil
	}
	return i.r.Next()
}

func (i *InsertReader) Fork() (*InsertReader, error) {
	if i.r == nil {
		return &InsertReader{}, nil
	}
	forked, err := i.r.Fork()
	if err != nil {
		return nil, err
	}
	return &InsertReader{r: forked}, nil
}
