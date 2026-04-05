package column

import (
	"encoding/binary"
	"io"
	"math"
	"unicode/utf8"

	"github.com/MichaelMure/leb128"

	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

// ValueReader is a stateful reader for value columns.
type ValueReader struct {
	r *ioutil.SubReader
}

func NewValueReader(r *ioutil.SubReader) *ValueReader {
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

// ValueWriter is a stateful encoder for value columns. It writes value metadata
// (as RLE uint64) and raw value bytes to separate io.Writers, matching what
// ValueReader decodes.
type ValueWriter struct {
	meta      *ValueMetadataWriter
	val       io.Writer
	err       error
	hasValues bool
}

func NewValueWriter(meta io.Writer, val io.Writer) *ValueWriter {
	return &ValueWriter{meta: NewValueMetadataWriter(meta), val: val}
}

func (vw *ValueWriter) Append(action types.Action) {
	if vw.err != nil {
		return
	}
	if HasScalarValue(action) {
		vw.hasValues = true
		m, b := EncodeValue(action)
		vw.meta.Append(m)
		if len(b) > 0 {
			_, vw.err = vw.val.Write(b)
		}
	} else {
		vw.meta.Append(NewValueMetadata(ValueTypeNull, 0))
	}
}

// HasValues reports whether any appended action had a scalar value.
func (vw *ValueWriter) HasValues() bool { return vw.hasValues }

// Flush writes the final metadata run and returns any accumulated error.
func (vw *ValueWriter) Flush() error {
	if vw.err != nil {
		return vw.err
	}
	return vw.meta.Flush()
}

// HasScalarValue reports whether the action carries a scalar value requiring
// entries in the value metadata and value columns.
func HasScalarValue(a types.Action) bool {
	return a.Kind == types.ActionSet || a.Kind == types.ActionInc
}

// EncodeValue returns the ValueMetadata and raw value bytes for a Set or Inc
// action. For actions without a scalar value, returns (Null metadata, nil).
func EncodeValue(a types.Action) (ValueMetadata, []byte) {
	if a.Kind != types.ActionSet && a.Kind != types.ActionInc {
		return NewValueMetadata(ValueTypeNull, 0), nil
	}
	switch v := a.Value.(type) {
	case nil:
		return NewValueMetadata(ValueTypeNull, 0), nil
	case bool:
		if v {
			return NewValueMetadata(ValueTypeTrue, 0), nil
		}
		return NewValueMetadata(ValueTypeFalse, 0), nil
	case string:
		b := []byte(v)
		return NewValueMetadata(ValueTypeString, uint64(len(b))), b
	case uint64:
		b := leb128.EncodeU64(v)
		return NewValueMetadata(ValueTypeUInt, uint64(len(b))), b
	case int64:
		b := leb128.EncodeS64(v)
		return NewValueMetadata(ValueTypeInt, uint64(len(b))), b
	case float64:
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, math.Float64bits(v))
		return NewValueMetadata(ValueTypeFloat, 8), b
	case []byte:
		return NewValueMetadata(ValueTypeBytes, uint64(len(v))), v
	case types.Counter:
		b := leb128.EncodeS64(int64(v))
		return NewValueMetadata(ValueTypeCounter, uint64(len(b))), b
	case types.Timestamp:
		b := leb128.EncodeS64(int64(v))
		return NewValueMetadata(ValueTypeTimestamp, uint64(len(b))), b
	default:
		return NewValueMetadata(ValueTypeNull, 0), nil
	}
}
