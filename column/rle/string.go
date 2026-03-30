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
func NewNullString() NullableValue[string]             { return nullable[string]{null: true} }

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

// TODO: bad encoder
// EncodeString encodes vals as a single literal RLE run of strings.
func EncodeString(vals ...string) []byte {
	b := leb128.EncodeS64(int64(-len(vals)))
	for _, s := range vals {
		b = append(b, leb128.EncodeU64(uint64(len(s)))...)
		b = append(b, s...)
	}
	return b
}
