package column

import (
	"fmt"
	"io"

	"github.com/jcalabro/leb128"
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

// specification describes the content of a column
type specification uint32

func newSpecification(id uint32, _type Type, deflate bool) specification {
	s := id << 4
	if deflate {
		s |= 0b1000
	}
	s |= uint32(_type)
	return specification(s)
}

func (s specification) ID() uint32 {
	return uint32(s >> 4)
}

func (s specification) Type() Type {
	return Type(s & 0b111)
}

func (s specification) Deflate() bool {
	return (s & 0b1000) != 0
}

func readSpecification(r io.Reader) (specification, error) {
	u, err := leb128.DecodeU32(r)
	if err != nil {
		return 0, fmt.Errorf("error reading column specification: %w", err)
	}
	return specification(u), nil
}

func writeSpecification(w io.Writer, spec specification) error {
	_, err := w.Write(leb128.EncodeU32(uint32(spec)))
	return err
}
