package column

import (
	"fmt"
)

const maxValueMetadataLength = 0x0fffffff
const maxValueMetadataType = 0x0f

type ValueType byte

const (
	ValueTypeNull      ValueType = 0
	ValueTypeFalse     ValueType = 1
	ValueTypeTrue      ValueType = 2
	ValueTypeUInt      ValueType = 3
	ValueTypeInt       ValueType = 4
	ValueTypeFloat     ValueType = 5
	ValueTypeString    ValueType = 6
	ValueTypeBytes     ValueType = 7
	ValueTypeCounter   ValueType = 8
	ValueTypeTimestamp ValueType = 9
)

type ValueMetadata uint64

func NewValueMetadata(t ValueType, length uint64) ValueMetadata {
	if length > maxValueMetadataLength {
		panic("overflow of value metadata length")
	}
	if t > maxValueMetadataType {
		panic("overflow of value metadata type")
	}
	return ValueMetadata(uint64(length)<<4 | uint64(t))
}

func (vm ValueMetadata) Type() ValueType {
	return ValueType(vm & 0x0f)
}

func (vm ValueMetadata) Length() uint64 {
	return uint64(vm >> 4)
}

var valueTypeNames = [...]string{
	ValueTypeNull:      "null",
	ValueTypeFalse:     "false",
	ValueTypeTrue:      "true",
	ValueTypeUInt:      "uint",
	ValueTypeInt:       "int",
	ValueTypeFloat:     "float",
	ValueTypeString:    "string",
	ValueTypeBytes:     "bytes",
	ValueTypeCounter:   "counter",
	ValueTypeTimestamp: "timestamp",
}

func (vm ValueMetadata) String() string {
	return fmt.Sprintf("(t=%s l=%d)", valueTypeNames[vm.Type()], vm.Length())
}
