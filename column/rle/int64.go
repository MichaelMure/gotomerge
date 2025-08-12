package rle

import (
	"io"
	"iter"
	"strconv"

	"github.com/jcalabro/leb128"
)

type NullableInt64 = NullableValue[int64]

func ReadInt64RLE(r io.Reader) iter.Seq2[NullableInt64, error] {
	return rle(r, int64Rig)
}

type nullableInt64 struct {
	val  int64
	null bool
}

func (n nullableInt64) Value() (int64, bool) {
	return n.val, !n.null
}

func (n nullableInt64) String() string {
	if n.null {
		return "null"
	}
	return strconv.FormatInt(n.val, 10)
}

var int64Rig = nullableRig[int64]{
	valid: func(v int64) bool {
		return true
	},
	null: func() NullableInt64 {
		return nullableInt64{null: true}
	},
	read: func(r io.Reader) (NullableInt64, error) {
		val, err := leb128.DecodeS64(r)
		if err != nil {
			return nullableInt64{}, err
		}
		return nullableInt64{val: val}, nil
	},
}
