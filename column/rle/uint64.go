package rle

import (
	"io"
	"iter"

	"github.com/jcalabro/leb128"
)

type NullableUint64 = NullableValue[uint64]

func ReadUint64RLE(r io.Reader) iter.Seq2[NullableUint64, error] {
	return rle(r, uint64Rig)
}

type nullableUint64 struct {
	val  uint64
	null bool
}

func (n nullableUint64) Value() (uint64, bool) {
	if n.null {
		return 0, false
	}
	return n.val, true
}

var uint64Rig = nullableRig[uint64]{
	valid: func(v uint64) bool {
		return true
	},
	null: func() NullableUint64 {
		return nullableUint64{null: true}
	},
	read: func(r io.Reader) (NullableUint64, error) {
		val, err := leb128.DecodeU64(r)
		if err != nil {
			return nullableUint64{}, err
		}
		return nullableUint64{val: val}, nil
	},
}
