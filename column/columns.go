package column

import (
	"fmt"

	ioutil "gotomerge/utils/io"
)

var ErrDone = fmt.Errorf("iterator done")

type ErrUnexpectedNull string

func (e ErrUnexpectedNull) Error() string {
	return fmt.Sprintf("unexpected null: %s", string(e))
}

type ErrOutOfRange string

func (e ErrOutOfRange) Error() string {
	return fmt.Sprintf("%s out of range", string(e))
}

// There are two types of columns:
// - low-level/raw: the directly stored format (bool, delta, uleb ...)
// - high level: aggregate low-level columns into higher level types (ObjectID, Operation...)
//
// While the normal push style iterator works just fine for low-level column, we need to switch
// to pull style for the high level one. This is because we need to simultaneously pull values from
// multiple columns at the same time, while not knowing beforehand how many values there are.

// Opt creates a typed column reader from an optional SubReader.
// Returns (nil, nil) if r is nil (column absent in the binary).
// Returns (nil, err) if forking fails.
func Opt[T any](r *ioutil.SubReader, ctor func(*ioutil.SubReader) *T) (*T, error) {
	if r == nil {
		return nil, nil
	}
	sub, err := r.SubReaderOffset(0)
	if err != nil {
		return nil, err
	}
	return ctor(sub), nil
}

// Req creates a typed column reader from a required SubReader.
// Returns an error if r is nil (column absent) or if forking fails.
func Req[T any](r *ioutil.SubReader, ctor func(*ioutil.SubReader) *T, name string) (*T, error) {
	if r == nil {
		return nil, fmt.Errorf("missing required column: %s", name)
	}
	sub, err := r.SubReaderOffset(0)
	if err != nil {
		return nil, fmt.Errorf("column %s: %w", name, err)
	}
	return ctor(sub), nil
}
