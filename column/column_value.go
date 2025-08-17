package column

import (
	"encoding/binary"
	"fmt"
	"io"
	"iter"
	"unicode/utf8"

	"github.com/jcalabro/leb128"

	"gotomerge/lbuf"
	"gotomerge/types"
)

type ValueColumn struct {
	r io.Reader
}

type ValueColumnIter = iter.Seq2[any, error]

func NewValueColumn(r io.Reader) ValueColumn {
	return ValueColumn{r: r}
}

func (vc ValueColumn) Iter(meta []ValueMetadata) ValueColumnIter {
	return func(yield func(any, error) bool) {
		for _, metadata := range meta {
			var val any
			var err error
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
				// TODO: little endian? complete the spec
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
				// TODO: should we read to write back later?
				yield(nil, fmt.Errorf("unknown value type %d", metadata.Type()))
				return
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
