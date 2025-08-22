package column

import "fmt"

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
