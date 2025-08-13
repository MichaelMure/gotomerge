package rle

import (
	"fmt"
	"io"
	"iter"

	"github.com/jcalabro/leb128"

	"gotomerge/lbuf"
)

func ReadStringRLE(r *lbuf.Reader) iter.Seq2[NullableValue[string], error] {
	return rle(r, stringRig)
}

// nullString is a marker for a null value of type string
// It's a 1-byte non-utf8 string: continuation byte without a preceding starter.
// This is used to have zero memory overhead in the normal case and a tiny value for null.
const nullString = "\x80"

type NullableString string

func NewNullableString(s string) NullableString {
	return NullableString(s)
}

func NewNullString() NullableString {
	return NullableString(nullString)
}

func (s NullableString) Value() (string, bool) {
	if s == nullString {
		return "", false
	}
	return string(s), true
}

var stringRig = nullableRig[string]{
	valid: func(s string) bool {
		return s != nullString
	},
	null: func() NullableValue[string] {
		return NullableString(nullString)
	},
	read: func(r *lbuf.Reader) (NullableValue[string], error) {
		strLen, err := leb128.DecodeU64(r)
		if err != nil {
			return NullableString(""), err
		}
		buf := make([]byte, strLen)
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return NullableString(""), err
		}
		str := string(buf)
		if str == nullString {
			return NullableString(""), fmt.Errorf("invalid string")
		}
		return NullableString(str), nil
	},
}
