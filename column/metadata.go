package column

import (
	"fmt"
	"io"

	"github.com/MichaelMure/leb128"
)

type Metadata []struct {
	Spec   Specification
	Length uint64 // in bytes
}

func ReadMetadata(r io.Reader) (Metadata, error) {
	n, err := leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading column metadata length: %w", err)
	}
	if n == 0 {
		// special case for empty metadata
		return nil, nil
	}

	// limit pre-allocation to avoid DOS
	allocate := n
	// TODO: adjust with reasonable value
	if n > 128 {
		allocate = 128
	}

	var prevSpec Specification
	res := make(Metadata, 0, allocate)
	for i := uint64(0); i < n; i++ {
		spec, err := readSpecification(r)
		if err != nil {
			return nil, fmt.Errorf("error reading column metadata spec: %w", err)
		}
		if spec <= prevSpec {
			return nil, fmt.Errorf("column metadata IDs must be sorted and unique")
		}
		if i != 0 && uint32(spec)|0b1000 == uint32(prevSpec)|0b1000 {
			return nil, fmt.Errorf("both uncompressed and compressed columns present with the same ID")
		}
		if spec.Type() == TypeValue &&
			((prevSpec.Type() != TypeValueMetadata) || spec.ID() != prevSpec.ID()) {
			return nil, fmt.Errorf("value column must be preceded by a value metadata column with the same ID")
		}

		// Bellow test is in the spec, but it's actually already impossible with the constraints above
		// if spec.Type() == TypeValueMetadata {
		// 	for j := uint64(0); j < i-1; j++ {
		// 		if res[j].Spec.ID() == spec.ID() {
		// 			return nil, fmt.Errorf("value metadata ID must be unique")
		// 		}
		// 	}
		// }

		prevSpec = spec
		length, err := leb128.DecodeU64(r)
		if err != nil {
			return nil, fmt.Errorf("error reading column metadata length: %w", err)
		}
		res = append(res, struct {
			Spec   Specification
			Length uint64
		}{Spec: spec, Length: length})
	}

	return res, nil
}

func WriteMetadata(w io.Writer, metadata Metadata) error {
	_, err := w.Write(leb128.EncodeU64(uint64(len(metadata))))
	if err != nil {
		return err
	}
	for _, m := range metadata {
		err = writeSpecification(w, m.Spec)
		if err != nil {
			return err
		}
		_, err = w.Write(leb128.EncodeU64(m.Length))
		if err != nil {
			return err
		}
	}
	return nil
}
