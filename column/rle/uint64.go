package rle

import (
	"iter"
	"strconv"

	"github.com/jcalabro/leb128"

	"gotomerge/lbuf"
)

type NullableUint64 = NullableValue[uint64]

func ReadUint64RLE(r *lbuf.Reader) iter.Seq2[NullableUint64, error] {
	return rle(r, uint64Rig)
}

type nullableUint64 struct {
	val  uint64
	null bool
}

func (n nullableUint64) Value() (uint64, bool) {
	return n.val, !n.null
}

func (n nullableUint64) String() string {
	if n.null {
		return "null"
	}
	return strconv.FormatUint(n.val, 10)
}

var uint64Rig = nullableRig[uint64]{
	valid: func(v uint64) bool {
		return true
	},
	null: func() NullableUint64 {
		return nullableUint64{null: true}
	},
	read: func(r *lbuf.Reader) (NullableUint64, error) {
		val, err := leb128.DecodeU64(r)
		if err != nil {
			return nullableUint64{null: true}, err
		}
		return nullableUint64{val: val, null: false}, nil
	},
}
