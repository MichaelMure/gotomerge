package rle

import (
	"io"
	"iter"
	"strconv"

	"github.com/jcalabro/leb128"
)

func ReadInt64RLE(r io.Reader) iter.Seq2[NullableValue[int64], error] {
	return rle(r, int64Rig)
}

type NullableInt64 struct {
	val  int64
	null bool
}

func NewNullableInt64(val int64) NullableInt64 {
	return NullableInt64{val: val}
}

func NewNullInt64() NullableInt64 {
	return NullableInt64{null: true}
}

func (n NullableInt64) Value() (int64, bool) {
	return n.val, !n.null
}

func (n NullableInt64) String() string {
	if n.null {
		return "<null>"
	}
	return strconv.FormatInt(n.val, 10)
}

var int64Rig = nullableRig[int64]{
	valid: func(v int64) bool {
		return true
	},
	null: func() NullableValue[int64] {
		return NullableInt64{null: true}
	},
	read: func(r io.Reader) (NullableValue[int64], error) {
		val, err := leb128.DecodeS64(r)
		if err != nil {
			return NullableInt64{}, err
		}
		return NullableInt64{val: val}, nil
	},
}
