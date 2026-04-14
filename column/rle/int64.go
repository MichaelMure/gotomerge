package rle

import (
	"bytes"
	"io"

	"github.com/MichaelMure/leb128"

	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

type Int64Reader = Reader[int64]

func NewInt64Reader(r ioutil.SubReader) *Int64Reader {
	return NewReader[int64](r, readS64)
}

func readS64(r ioutil.ByteReader) (int64, error) { return leb128.DecodeS64(r) }

func NewNullableInt64(v int64) NullableValue[int64] { return NullableValue[int64]{val: v} }
func NewNullInt64() NullableValue[int64]            { return NullableValue[int64]{null: true} }

func NewInt64Writer(w io.Writer) *Writer[int64] {
	return NewWriter(w, leb128.AppendS64)
}

// EncodeInt64 encodes vals as an RLE int64 sequence.
func EncodeInt64(vals ...int64) []byte {
	var buf bytes.Buffer
	w := NewInt64Writer(&buf)
	for _, v := range vals {
		w.Append(NewNullableInt64(v))
	}
	_ = w.Flush()
	return buf.Bytes()
}

// EncodeInt64Delta delta-encodes a nullable int64 sequence using a Writer.
// Each non-null value is stored as its delta from the previous non-null value
// (accumulator starts at 0). Nulls do not advance the accumulator.
func EncodeInt64Delta(vals []NullableValue[int64]) []byte {
	var buf bytes.Buffer
	w := NewInt64Writer(&buf)
	var acc int64
	for _, nv := range vals {
		if v, ok := nv.Value(); ok {
			w.Append(NewNullableInt64(v - acc))
			acc = v
		} else {
			w.Append(NewNullInt64())
		}
	}
	_ = w.Flush()
	return buf.Bytes()
}
