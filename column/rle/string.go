package rle

import (
	"fmt"
	"io"
	"iter"

	"github.com/jcalabro/leb128"

	"gotomerge/lbuf"
)

type NullableString = NullableValue[string]

func ReadStringRLE(r *lbuf.Reader) iter.Seq2[NullableString, error] {
	return rle(r, stringRig)
}

// nullString is a marker for a null value of type string
// It's a 1-byte non-utf8 string: continuation byte without a preceding starter.
// This is used to have zero memory overhead in the normal case and a tiny value for null.
const nullString = "\x80"

type nullableString string

func (s nullableString) Value() (string, bool) {
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
		return nullableString(nullString)
	},
	read: func(r *lbuf.Reader) (NullableValue[string], error) {
		strLen, err := leb128.DecodeU64(r)
		if err != nil {
			return nullableString(""), err
		}
		buf := make([]byte, strLen)
		_, err = io.ReadFull(r, buf)
		if err != nil {
			return nullableString(""), err
		}
		str := string(buf)
		if str == nullString {
			return nullableString(""), fmt.Errorf("invalid string")
		}
		return nullableString(str), nil
	},
}
