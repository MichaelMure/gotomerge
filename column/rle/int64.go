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
func NewNullInt64() NullableValue[int64]             { return nullable[int64]{null: true} }
