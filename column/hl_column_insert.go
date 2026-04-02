package column

import "io"

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

// InsertWriter is a stateful encoder for insert (bool) columns.
type InsertWriter struct {
	w          *BoolWriter
	hasInserts bool
}

func NewInsertWriter(w io.Writer) *InsertWriter {
	return &InsertWriter{w: NewBoolWriter(w)}
}

func (i *InsertWriter) Append(insert bool) {
	i.w.Append(insert)
	if insert {
		i.hasInserts = true
	}
}

func (i *InsertWriter) HasInserts() bool { return i.hasInserts }

func (i *InsertWriter) Flush() error { return i.w.Flush() }
