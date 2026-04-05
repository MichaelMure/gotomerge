package rle

import (
	"bytes"
	"io"

	"github.com/MichaelMure/leb128"

	ioutil "gotomerge/utils/io"
)

type StringReader = Reader[string]

func NewStringReader(r *ioutil.SubReader) *StringReader {
	return NewReader[string](r, readStringValue)
}

func NewNullableString(s string) NullableValue[string] { return NullableValue[string]{val: s} }
func NewNullString() NullableValue[string]             { return NullableValue[string]{null: true} }

func NewStringWriter(w io.Writer) *Writer[string] {
	return NewWriter(w, func(buf []byte, s string) []byte {
		// len || str
		return append(leb128.AppendU64(buf, uint64(len(s))), s...)
	})
}

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

// EncodeString encodes vals as an RLE string sequence.
func EncodeString(vals ...string) []byte {
	var buf bytes.Buffer
	w := NewStringWriter(&buf)
	for _, s := range vals {
		w.Append(NewNullableString(s))
	}
	_ = w.Flush()
	return buf.Bytes()
}
