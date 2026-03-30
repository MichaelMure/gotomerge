package column

import (
	"encoding/binary"
	"unicode/utf8"

	"github.com/jcalabro/leb128"

	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

// ValueReader is a stateful reader for value columns.
type ValueReader struct {
	r ioutil.SubReader
}

func NewValueReader(r ioutil.SubReader) *ValueReader {
	return &ValueReader{r: r}
}

func (vr *ValueReader) Next(meta ValueMetadata) (any, error) {
	var val any
	var err error
	switch meta.Type() {
	case ValueTypeNull:
		val = nil
	case ValueTypeFalse:
		val = false
	case ValueTypeTrue:
		val = true
	case ValueTypeUInt:
		val, err = leb128.DecodeU64(vr.r)
	case ValueTypeInt:
		val, err = leb128.DecodeS64(vr.r)
	case ValueTypeFloat:
		var f float64
		err = binary.Read(vr.r, binary.LittleEndian, &f)
		val = f
	case ValueTypeString:
		str, rerr := ioutil.ReadStringLimitedPrealloc(vr.r, meta.Length())
		if rerr != nil {
			return nil, rerr
		}
		if !utf8.ValidString(str) {
			return string(utf8.RuneError), nil
		}
		val = str
	case ValueTypeBytes:
		val, err = ioutil.ReadBytesLimitedPrealloc(vr.r, meta.Length())
	case ValueTypeCounter:
		var ctr int64
		ctr, err = leb128.DecodeS64(vr.r)
		val = types.Counter(ctr)
	case ValueTypeTimestamp:
		var ts int64
		ts, err = leb128.DecodeS64(vr.r)
		val = types.Timestamp(ts)
	default:
		data, rerr := ioutil.ReadBytesLimitedPrealloc(vr.r, meta.Length())
		if rerr != nil {
			return nil, rerr
		}
		val = UnknownValue{Type: meta.Type(), Bytes: data}
	}
	if err != nil {
		return nil, err
	}
	return val, nil
}

func (vr *ValueReader) Fork() (*ValueReader, error) {
	sub, err := vr.r.SubReaderOffset(0)
	if err != nil {
		return nil, err
	}
	return &ValueReader{r: sub}, nil
}
