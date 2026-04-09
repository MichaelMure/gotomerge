package column

import (
	"fmt"
	"io"

	"github.com/MichaelMure/leb128"
)

const maxSpecificationId = 0x0fffffff

type Type byte

const (
	TypeGroup         Type = 0x00
	TypeActor         Type = 0x01
	TypeULEB128       Type = 0x02
	TypeDelta         Type = 0x03
	TypeBool          Type = 0x04
	TypeString        Type = 0x05
	TypeValueMetadata Type = 0x06
	TypeValue         Type = 0x07
)

// Specification describes the content of a column
type Specification uint32

func newSpecification(id uint32, _type Type, deflate bool) Specification {
	if id > maxSpecificationId {
		panic("overflow of specification ID")
	}
	s := id << 4
	if deflate {
		s |= 0b1000
	}
	s |= uint32(_type)
	return Specification(s)
}

func (s Specification) ID() uint32 {
	return uint32(s >> 4)
}

func (s Specification) Type() Type {
	return Type(s & 0b111)
}

func (s Specification) Deflate() bool {
	return (s & 0b1000) != 0
}

func (s Specification) WithoutDeflate() Specification {
	return s &^ 0b1000
}

func (s Specification) WithDeflate() Specification {
	return s | 0b1000
}

var typeNames = [...]string{
	TypeGroup:         "group",
	TypeActor:         "actor",
	TypeULEB128:       "uleb128",
	TypeDelta:         "delta",
	TypeBool:          "bool",
	TypeString:        "string",
	TypeValueMetadata: "value_metadata",
	TypeValue:         "value",
}

func (t Type) String() string {
	return typeNames[t]
}

func (s Specification) String() string {
	return fmt.Sprintf("spec(%d: id=%d, type=%s, deflate=%t)", uint64(s), s.ID(), s.Type(), s.Deflate())
}

func readSpecification(r io.Reader) (Specification, error) {
	u, err := leb128.DecodeU32(r)
	if err != nil {
		return 0, fmt.Errorf("error reading column specification: %w", err)
	}
	return Specification(u), nil
}

func writeSpecification(w io.Writer, spec Specification) error {
	_, err := w.Write(leb128.EncodeU32(uint32(spec)))
	return err
}
