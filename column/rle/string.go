package rle

import (
	"io"

	"github.com/jcalabro/leb128"

	ioutil "gotomerge/utils/io"
)

type StringReader = Reader[string]

func NewStringReader(r ioutil.SubReader) *StringReader {
	return NewReader[string](r, readStringValue)
}

func NewNullableString(s string) NullableValue[string] { return nullable[string]{val: s} }
func NewNullString() NullableValue[string]              { return nullable[string]{null: true} }

func readStringValue(r io.Reader) (string, error) {
	strLen, err := leb128.DecodeU64(r)
	if err != nil {
		return "", err
	}
	buf := make([]byte, strLen)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
