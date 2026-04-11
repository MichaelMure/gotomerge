package rle

import (
	"io"

	"github.com/MichaelMure/leb128"

	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

type rleState byte

const (
	rleStateIdle     rleState = 0
	rleStateRepeated rleState = 1
	rleStateNull     rleState = 2
	rleStateLiteral  rleState = 3
)

// NullableValue represents a value that may be absent (null).
type NullableValue[T any] struct {
	val  T
	null bool
}

func (n NullableValue[T]) Value() (T, bool) { return n.val, !n.null }

// Reader is a generic stateful reader for RLE-encoded columns.
// T is the element type; readFn decodes one non-null value from an io.Reader.
type Reader[T any] struct {
	r          *ioutil.SubReader
	readFn     func(io.Reader) (T, error)
	state      rleState
	remaining  int64
	cachedVal  T
	cachedNull bool
}

func NewReader[T any](r *ioutil.SubReader, readFn func(io.Reader) (T, error)) *Reader[T] {
	return &Reader[T]{r: r, readFn: readFn}
}

func (rd *Reader[T]) Next() (NullableValue[T], error) {
	var zero NullableValue[T]
	for {
		switch rd.state {
		case rleStateIdle:
			L, err := leb128.DecodeS64(rd.r)
			if err == io.EOF {
				return zero, io.EOF
			}
			if err != nil {
				return zero, err
			}
			switch {
			case L > 1:
				val, err := rd.readFn(rd.r)
				if err != nil {
					return zero, err
				}
				rd.state = rleStateRepeated
				rd.remaining = L
				rd.cachedVal = val
				rd.cachedNull = false
				continue
			case L == 0:
				count, err := leb128.DecodeU64(rd.r)
				if err != nil {
					return zero, err
				}
				rd.state = rleStateNull
				rd.remaining = int64(count)
				continue
			default: // L < 0 or L == 1
				if L == 1 {
					val, err := rd.readFn(rd.r)
					if err != nil {
						return zero, err
					}
					rd.state = rleStateRepeated
					rd.remaining = 1
					rd.cachedVal = val
					rd.cachedNull = false
					continue
				}
				rd.state = rleStateLiteral
				rd.remaining = -L
				continue
			}
		case rleStateRepeated:
			rd.remaining--
			if rd.remaining == 0 {
				rd.state = rleStateIdle
			}
			return NullableValue[T]{val: rd.cachedVal, null: rd.cachedNull}, nil
		case rleStateNull:
			rd.remaining--
			if rd.remaining == 0 {
				rd.state = rleStateIdle
			}
			return NullableValue[T]{null: true}, nil
		case rleStateLiteral:
			val, err := rd.readFn(rd.r)
			if err != nil {
				return zero, err
			}
			rd.remaining--
			if rd.remaining == 0 {
				rd.state = rleStateIdle
			}
			return NullableValue[T]{val: val}, nil
		}
	}
}

func (rd *Reader[T]) Fork() (*Reader[T], error) {
	sub, err := rd.r.SubReaderOffset(0)
	if err != nil {
		return nil, err
	}
	cp := *rd
	cp.r = sub
	return &cp, nil
}
