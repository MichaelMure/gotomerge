package rle

import (
	"io"
	"iter"
	"strconv"

	"github.com/jcalabro/leb128"
)

func ReadUint64RLE(r io.Reader) iter.Seq2[NullableValue[uint64], error] {
	return rle(r, uint64Rig)
}

type NullableUint64 struct {
	val  uint64
	null bool
}

func NewNullableUint64(val uint64) NullableUint64 {
	return NullableUint64{val: val}
}

func NewNullUint64() NullableUint64 {
	return NullableUint64{null: true}
}

func (n NullableUint64) Value() (uint64, bool) {
	return n.val, !n.null
}

func (n NullableUint64) String() string {
	if n.null {
		return "<null>"
	}
	return strconv.FormatUint(n.val, 10)
}

var uint64Rig = nullableRig[uint64]{
	valid: func(v uint64) bool {
		return true
	},
	null: func() NullableValue[uint64] {
		return NullableUint64{null: true}
	},
	read: func(r io.Reader) (NullableValue[uint64], error) {
		val, err := leb128.DecodeU64(r)
		if err != nil {
			return NullableUint64{null: true}, err
		}
		return NullableUint64{val: val, null: false}, nil
	},
}
