package rle

import (
	"io"
	"iter"

	"github.com/jcalabro/leb128"
)

type NullableValue[T any] interface {
	// Value returns:
	// - (val, true) when the value exists
	// - (_, false) when the value is null
	Value() (T, bool)
}

type nullableRig[T any] struct {
	valid func(T) bool
	null  func() NullableValue[T]
	read  func(io.Reader) (NullableValue[T], error)
}

func rle[T any](r io.Reader, rig nullableRig[T]) iter.Seq2[NullableValue[T], error] {
	return func(yield func(NullableValue[T], error) bool) {
		for {
			L, err := leb128.DecodeS64(r)
			if err == io.EOF {
				return
			}
			if err != nil {
				yield(rig.null(), err)
				return
			}

			switch {
			case L > 1: // val repeated L times
				str, err := rig.read(r)
				if err != nil {
					yield(rig.null(), err)
					return
				}
				for i := int64(0); i < L; i++ {
					if !yield(str, nil) {
						return
					}
				}

			case L == 0: // null repeated val times
				val, err := leb128.DecodeU64(r)
				if err != nil {
					yield(rig.null(), err)
					return
				}
				for i := uint64(0); i < val; i++ {
					if !yield(rig.null(), nil) {
						return
					}
				}

			case L < 1: // L values will follow
				for i := int64(0); i < -L; i++ {
					str, err := rig.read(r)
					if err != nil {
						yield(rig.null(), err)
						return
					}
					if !yield(str, nil) {
						return
					}
				}
			}
		}
	}
}
