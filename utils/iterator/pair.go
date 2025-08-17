package iterutil

import (
	"fmt"
	"iter"
)

// PullPair pulls from two iterators simultaneously and applies a transformation function.
// It stops when either iterator is exhausted or encounters an error.
func PullPair[T1, T2, R any](
	iter1 iter.Seq2[T1, error],
	iter2 iter.Seq2[T2, error],
	transform func(T1, T2) (R, error),
) iter.Seq2[R, error] {
	return func(yield func(R, error) bool) {
		next1, stop1 := iter.Pull2(iter1)
		defer stop1()
		next2, stop2 := iter.Pull2(iter2)
		defer stop2()

		for {
			val1, err1, ok1 := next1()
			if err1 != nil {
				var zero R
				if !yield(zero, err1) {
					return
				}
				return
			}

			val2, err2, ok2 := next2()
			if err2 != nil {
				var zero R
				if !yield(zero, err2) {
					return
				}
				return
			}

			// Check if iterators have different lengths
			if ok1 != ok2 {
				var zero R
				if !yield(zero, fmt.Errorf("iterators have different lengths")) {
					return
				}
				return
			}

			// Both iterators are exhausted
			if !ok1 && !ok2 {
				return
			}

			// Transform the pair of values
			result, err := transform(val1, val2)
			if err != nil {
				if !yield(result, err) {
					return
				}
				return
			}

			if !yield(result, nil) {
				return
			}
		}
	}
}
