package rle

import (
	"bytes"
	"io"

	"github.com/MichaelMure/leb128"

	ioutil "gotomerge/utils/io"
)

type Uint64Reader = Reader[uint64]

func NewUint64Reader(r *ioutil.SubReader) *Uint64Reader {
	return NewReader[uint64](r, leb128.DecodeU64)
}

func NewNullableUint64(v uint64) NullableValue[uint64] { return NullableValue[uint64]{val: v} }
func NewNullUint64() NullableValue[uint64]             { return NullableValue[uint64]{null: true} }

func NewUint64Writer(w io.Writer) *Writer[uint64] {
	return NewWriter(w, leb128.AppendU64)
}

// EncodeUint64 encodes vals as an RLE uint64 sequence.
func EncodeUint64(vals ...uint64) []byte {
	var buf bytes.Buffer
	w := NewUint64Writer(&buf)
	for _, v := range vals {
		w.Append(NewNullableUint64(v))
	}
	_ = w.Flush()
	return buf.Bytes()
}
