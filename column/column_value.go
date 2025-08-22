package column

import (
	"encoding/binary"
	"io"
	"iter"
	"unicode/utf8"

	"github.com/jcalabro/leb128"

	"gotomerge/lbuf"
	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

type ValueColumn struct {
	r io.Reader
}

type ValueColumnIter = iter.Seq2[any, error]
type ValueColumnIterMaker interface {
	Iter(meta iter.Seq2[ValueMetadata, error]) ValueColumnIter
}

func NewValueColumn(r io.Reader) ValueColumn {
	return ValueColumn{r: r}
}

func (vc ValueColumn) Iter(meta iter.Seq2[ValueMetadata, error]) ValueColumnIter {
	return func(yield func(any, error) bool) {
		for metadata, err := range meta {
			if err != nil {
				yield(err, nil)
				return
			}
			var val any
			switch metadata.Type() {
			case ValueTypeNull:
				val = nil
			case ValueTypeFalse:
				val = false
			case ValueTypeTrue:
				val = true
			case ValueTypeUInt:
				val, err = leb128.DecodeU64(vc.r)
			case ValueTypeInt:
				val, err = leb128.DecodeS64(vc.r)
			case ValueTypeFloat:
				var f float64
				err = binary.Read(vc.r, binary.LittleEndian, &f)
				val = f
			case ValueTypeString:
				str, err := lbuf.ReadStringLimitedPrealloc(vc.r, metadata.Length())
				if err != nil {
					yield(nil, err)
					return
				}
				if !utf8.ValidString(str) {
					// special case: if it's not valid utf8, we replace the value by
					// the unicode replacement character
					if !yield(string(utf8.RuneError), nil) {
						return
					}
					continue
				}
				val = str
			case ValueTypeBytes:
				val, err = lbuf.ReadBytesLimitedPrealloc(vc.r, metadata.Length())
			case ValueTypeCounter:
				var ctr int64
				ctr, err = leb128.DecodeS64(vc.r)
				val = types.Counter(ctr)
			case ValueTypeTimestamp:
				var ts int64
				ts, err = leb128.DecodeS64(vc.r)
				val = types.Timestamp(ts)
			default:
				data, err := ioutil.ReadBytesLimitedPrealloc(vc.r, metadata.Length())
				if err != nil {
					yield(nil, err)
					return
				}
				val = UnknownValue{Type: metadata.Type(), Bytes: data}
			}
			if err != nil {
				yield(nil, err)
				return
			}
			if !yield(val, nil) {
				return
			}
		}
	}
}

type NullValueColumn struct{}

func (NullValueColumn) Iter(meta iter.Seq2[ValueMetadata, error]) ValueColumnIter {
	return func(yield func(any, error) bool) {
		for metadata, err := range meta {
			if err != nil {
				yield(nil, err)
				return
			}
			switch metadata.Type() {
			case ValueTypeNull:
				if !yield(nil, nil) {
					return
				}
			case ValueTypeFalse:
				// TODO: what to do here?
				panic("what to do here?")
			case ValueTypeTrue:
				// TODO: what to do here?
				panic("what to do here?")
			case ValueTypeUInt:
				// TODO: what to do here?
				panic("what to do here?")
			case ValueTypeInt:
				// TODO: what to do here?
				panic("what to do here?")
			case ValueTypeFloat:
				// TODO: what to do here?
				panic("what to do here?")
			case ValueTypeString:
				// TODO: what to do here?
				panic("what to do here?")
			case ValueTypeBytes:
				// TODO: what to do here?
				panic("what to do here?")
			case ValueTypeCounter:
				// TODO: what to do here?
				panic("what to do here?")
			case ValueTypeTimestamp:
				// TODO: what to do here?
				panic("what to do here?")
			default:
				// TODO: what to do here?
				panic("what to do here?")
			}
		}
	}
}
