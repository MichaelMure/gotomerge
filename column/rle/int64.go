package rle

import (
	"github.com/jcalabro/leb128"

	ioutil "gotomerge/utils/io"
)

type Int64Reader = Reader[int64]

func NewInt64Reader(r ioutil.SubReader) *Int64Reader {
	return NewReader[int64](r, leb128.DecodeS64)
}

func NewNullableInt64(v int64) NullableValue[int64] { return nullable[int64]{val: v} }
func NewNullInt64() NullableValue[int64]            { return nullable[int64]{null: true} }

// TODO: bad encoder
// EncodeInt64 encodes vals as a single literal RLE run of int64 values.
func EncodeInt64(vals ...int64) []byte {
	b := leb128.EncodeS64(int64(-len(vals)))
	for _, v := range vals {
		b = append(b, leb128.EncodeS64(v)...)
	}
	return b
}
