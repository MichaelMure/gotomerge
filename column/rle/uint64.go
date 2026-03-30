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
func NewNullUint64() NullableValue[uint64]             { return nullable[uint64]{null: true} }

// TODO: bad encoder
// EncodeUint64 encodes vals as a single literal RLE run of uint64 values.
func EncodeUint64(vals ...uint64) []byte {
	b := leb128.EncodeS64(int64(-len(vals)))
	for _, v := range vals {
		b = append(b, leb128.EncodeU64(v)...)
	}
	return b
}
