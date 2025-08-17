package iterutil

import (
	"fmt"
	"iter"
)

// PullTriple pulls from three iterators simultaneously and applies a transformation function.
// It stops when any iterator is exhausted or encounters an error.
func PullTriple[T1, T2, T3, R any](
	iter1 iter.Seq2[T1, error],
	iter2 iter.Seq2[T2, error],
	iter3 iter.Seq2[T3, error],
	transform func(T1, T2, T3) (R, error),
) iter.Seq2[R, error] {
	return func(yield func(R, error) bool) {
		next1, stop1 := iter.Pull2(iter1)
		defer stop1()
		next2, stop2 := iter.Pull2(iter2)
		defer stop2()
		next3, stop3 := iter.Pull2(iter3)
		defer stop3()

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

			val3, err3, ok3 := next3()
			if err3 != nil {
				var zero R
				if !yield(zero, err3) {
					return
				}
				return
			}

			// Check if all iterators have the same length
			if (ok1 || ok2 || ok3) != (ok1 && ok2 && ok3) {
				var zero R
				if !yield(zero, fmt.Errorf("iterators have different lengths")) {
					return
				}
				return
			}

			// All iterators are exhausted
			if !ok1 {
				return
			}

			// Transform the triple of values
			result, err := transform(val1, val2, val3)
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
