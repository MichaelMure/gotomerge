package column

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"iter"
	"unicode/utf8"
	"unsafe"

	"github.com/jcalabro/leb128"

	"gotomerge/types"
)

func ReadValueColumn(r io.Reader, meta []ValueMetadata) iter.Seq2[any, error] {
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
				val, err = leb128.DecodeU64(r)
			case ValueTypeInt:
				val, err = leb128.DecodeS64(r)
			case ValueTypeFloat:
				var f float64
				// TODO: little endian? complete the spec
				err = binary.Read(r, binary.LittleEndian, &f)
				val = f
			case ValueTypeString:
				strBytes, err := readWithLimitedPrealloc(r, metadata.Length())
				if err != nil {
					yield(nil, err)
					return
				}
				if !utf8.Valid(strBytes) {
					// special case: if it's not valid utf8, we replace the value by
					// the unicode replacement character
					if !yield(string(utf8.RuneError), nil) {
						return
					}
					continue
				}
				// zero-copy cast to string
				val = unsafe.String(unsafe.SliceData(strBytes), len(strBytes))
			case ValueTypeBytes:
				val, err = readWithLimitedPrealloc(r, metadata.Length())
			case ValueTypeCounter:
				var ctr int64
				ctr, err = leb128.DecodeS64(r)
				val = types.Counter(ctr)
			case ValueTypeTimestamp:
				var ts int64
				ts, err = leb128.DecodeS64(r)
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

func readWithLimitedPrealloc(r io.Reader, hintSize uint64) ([]byte, error) {
	prealloc := hintSize
	if prealloc > 10_000 {
		// limit the pre-allocation to 10kB to avoid DOS
		// the buffer will grow if there is actually more data to read
		// the downside is that it can create reallocation and copy for larger value.
		prealloc = 10_000
	}
	buf := bytes.NewBuffer(make([]byte, prealloc))
	_, err := buf.ReadFrom(io.LimitReader(r, int64(hintSize)))
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
