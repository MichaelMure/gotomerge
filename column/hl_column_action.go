package column

import (
	"fmt"
	"io"

	"github.com/MichaelMure/gotomerge/column/rle"
	"github.com/MichaelMure/gotomerge/types"
)

// ActionReader is a stateful reader for action columns.
type ActionReader struct {
	kind  *UlebReader
	meta  *ValueMetadataReader
	value *ValueReader
}

func NewActionReader(kind *UlebReader, meta *ValueMetadataReader, value *ValueReader) *ActionReader {
	return &ActionReader{kind: kind, meta: meta, value: value}
}

func (a *ActionReader) Next() (types.Action, error) {
	if a.kind == nil {
		return types.Action{}, ErrDone
	}

	nv, err := a.kind.Next()
	if err == io.EOF {
		// the action column has a special treatment. As all operations need an action,
		// we consider that we can't have nil columns (== implied null), and therefore
		// we consider the iteration done.
		return types.Action{}, ErrDone
	}
	if err != nil {
		return types.Action{}, err
	}
	kindVal, valid := nv.Value()
	if !valid {
		return types.Action{}, ErrUnexpectedNull("action kind")
	}

	var meta ValueMetadata
	if a.meta != nil {
		meta, err = a.meta.Next()
		if err != nil && err != io.EOF {
			return types.Action{}, err
		}
	}

	var value any
	if a.value != nil {
		value, err = a.value.Next(meta)
		if err != nil {
			return types.Action{}, err
		}
	} else {
		// absent value column: Null/False/True are encoded in the metadata type itself.
		switch meta.Type() {
		case ValueTypeNull:
			value = nil
		case ValueTypeFalse:
			value = false
		case ValueTypeTrue:
			value = true
		default:
			return types.Action{}, fmt.Errorf("value column absent for value type %v", meta.Type())
		}
	}

	err = types.ValidateAction(kindVal, value)
	if err != nil {
		return types.Action{}, err
	}

	return types.Action{
		Kind:  types.ActionKind(kindVal),
		Value: value,
	}, nil
}

func (a *ActionReader) Fork() (*ActionReader, error) {
	var kind *UlebReader
	var meta *ValueMetadataReader
	var value *ValueReader
	var err error

	if a.kind != nil {
		kind, err = a.kind.Fork()
		if err != nil {
			return nil, err
		}
	}
	if a.meta != nil {
		meta, err = a.meta.Fork()
		if err != nil {
			return nil, err
		}
	}
	if a.value != nil {
		value, err = a.value.Fork()
		if err != nil {
			return nil, err
		}
	}
	return &ActionReader{kind: kind, meta: meta, value: value}, nil
}

// ActionWriter is a stateful encoder for action columns (kind + value metadata + value bytes).
type ActionWriter struct {
	kind  *UlebWriter
	value *ValueWriter
}

func NewActionWriter(kind, meta, val io.Writer) *ActionWriter {
	return &ActionWriter{
		kind:  NewUlebWriter(kind),
		value: NewValueWriter(meta, val),
	}
}

func (a *ActionWriter) Append(action types.Action) {
	a.kind.Append(rle.NewNullableUint64(uint64(action.Kind)))
	a.value.Append(action)
}

func (a *ActionWriter) Flush() error {
	if err := a.kind.Flush(); err != nil {
		return err
	}
	return a.value.Flush()
}
