package rle

import (
	"io"

	"github.com/MichaelMure/leb128"
)

// Writer is a stateful RLE encoder that streams encoded runs to an io.Writer
// as they complete. Only the current literal run needs buffering; completed
// repeat and null runs are written immediately.
//
// Errors from the underlying io.Writer are sticky: the first write error is
// stored and all subsequent operations are no-ops. Call Flush() at the end to
// finalise the last run and retrieve any accumulated error.
type Writer[T comparable] struct {
	w        io.Writer
	err      error
	encodeFn func(T) []byte

	state       rleState
	nullCount   uint64
	repeatVal   T
	repeatCount int64
	literals    []T
}

func NewWriter[T comparable](w io.Writer, encodeFn func(T) []byte) *Writer[T] {
	return &Writer[T]{w: w, encodeFn: encodeFn}
}

// Append adds one nullable value to the encoder.
func (w *Writer[T]) Append(nv NullableValue[T]) {
	if w.err != nil {
		return
	}
	if v, ok := nv.Value(); ok {
		w.appendVal(v)
	} else {
		w.appendNull()
	}
}

// Flush writes the final buffered run and returns any accumulated error.
func (w *Writer[T]) Flush() error {
	if w.err != nil {
		return w.err
	}
	switch w.state {
	case rleStateNull:
		w.flushNull()
	case rleStateRepeated:
		if w.repeatCount == 1 {
			// Match Rust: a lone trailing value is flushed as a literal run (-1),
			// not a repeat run (+1). Both decode identically, but only the literal
			// form produces stable hashes that interoperate with Rust.
			w.literals = append(w.literals, w.repeatVal)
			w.flushLiteral()
		} else {
			w.flushRepeat()
		}
	case rleStateLiteral:
		w.flushLiteral()
	}
	return w.err
}

func (w *Writer[T]) appendNull() {
	switch w.state {
	case rleStateIdle:
		w.state = rleStateNull
		w.nullCount = 1
	case rleStateNull:
		w.nullCount++
	case rleStateRepeated:
		w.flushRepeat()
		w.state = rleStateNull
		w.nullCount = 1
	case rleStateLiteral:
		w.flushLiteral()
		w.state = rleStateNull
		w.nullCount = 1
	}
}

func (w *Writer[T]) appendVal(v T) {
	switch w.state {
	case rleStateIdle:
		w.state = rleStateRepeated
		w.repeatVal = v
		w.repeatCount = 1

	case rleStateNull:
		w.flushNull()
		w.state = rleStateRepeated
		w.repeatVal = v
		w.repeatCount = 1

	case rleStateRepeated:
		if v == w.repeatVal {
			w.repeatCount++
		} else if w.repeatCount >= 2 {
			w.flushRepeat()
			w.state = rleStateRepeated
			w.repeatVal = v
			w.repeatCount = 1
		} else {
			// Single value so far: not worth a repeat run; fold into a literal.
			w.state = rleStateLiteral
			w.literals = append(w.literals, w.repeatVal, v)
			w.repeatCount = 0
		}

	case rleStateLiteral:
		if len(w.literals) > 0 && v == w.literals[len(w.literals)-1] {
			// The last two values match — split a repeat run out of the literal.
			prev := w.literals[len(w.literals)-1]
			w.literals = w.literals[:len(w.literals)-1]
			if len(w.literals) > 0 {
				w.flushLiteral()
			}
			w.state = rleStateRepeated
			w.repeatVal = prev
			w.repeatCount = 2
		} else {
			w.literals = append(w.literals, v)
		}
	}
}

func (w *Writer[T]) write(b []byte) {
	if w.err != nil {
		return
	}
	_, w.err = w.w.Write(b)
}

func (w *Writer[T]) flushNull() {
	w.write(leb128.EncodeS64(0))
	w.write(leb128.EncodeU64(w.nullCount))
	w.nullCount = 0
	w.state = rleStateIdle
}

func (w *Writer[T]) flushRepeat() {
	w.write(leb128.EncodeS64(w.repeatCount))
	w.write(w.encodeFn(w.repeatVal))
	w.repeatCount = 0
	w.state = rleStateIdle
}

func (w *Writer[T]) flushLiteral() {
	if len(w.literals) == 0 {
		return
	}
	w.write(leb128.EncodeS64(int64(-len(w.literals))))
	for _, v := range w.literals {
		w.write(w.encodeFn(v))
	}
	w.literals = w.literals[:0]
	w.state = rleStateIdle
}
