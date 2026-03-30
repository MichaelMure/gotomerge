package rle

import (
	"github.com/jcalabro/leb128"

	ioutil "gotomerge/utils/io"
)

type Uint64Reader = Reader[uint64]

func NewUint64Reader(r ioutil.SubReader) *Uint64Reader {
	return NewReader[uint64](r, leb128.DecodeU64)
}

func NewNullableUint64(v uint64) NullableValue[uint64] { return nullable[uint64]{val: v} }
func NewNullUint64() NullableValue[uint64]              { return nullable[uint64]{null: true} }
