package column

import (
	"errors"
	"fmt"
	"io"
	"iter"
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

// isDone reports whether err signals iterator exhaustion.
func isDone(err error) bool {
	return errors.Is(err, ErrDone) || errors.Is(err, io.EOF)
}

// errSeq returns an iterator that immediately yields err and stops.
func errSeq[T any](err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var zero T
		yield(zero, err)
	}
}
