package column

import (
	"iter"

	"gotomerge/column/rle"
	"gotomerge/types"
)

type ActionColumnIter struct {
	nextKind  func() (rle.NullableValue[uint64], error, bool)
	stopKind  func()
	nextValue func() (any, error, bool)
	stopValue func()
}

func ActionColumn(kind UlebColumnIter, meta ValueMetadataColumnIter, values ValueColumnIterMaker) ActionColumnIter {
	var res ActionColumnIter
	if kind != nil {
		res.nextKind, res.stopKind = iter.Pull2(kind)
	}
	if values == nil {
		// if we don't have a value column, a column of null values is implied
		values = NullValueColumn{}
	}
	if meta != nil {
		res.nextValue, res.stopValue = iter.Pull2(values.Iter(meta))
	}
	return res
}

func (a ActionColumnIter) Next() (types.Action, error) {
	if a.nextKind == nil {
		// the action column has a special treatment. As all operations need an action,
		// we consider that we can't have nil columns (== implied null), and therefore
		// we consider the iteration done.
		return types.Action{}, ErrDone
	}

	kind, nullKind, err := extract(a.nextKind)
	if err != nil {
		return types.Action{}, err
	}
	if nullKind {
		return types.Action{}, ErrUnexpectedNull("action kind")
	}

	value, err, ok := a.nextValue()
	if err != nil {
		return types.Action{}, err
	}
	if !ok {
		return types.Action{}, ErrDone
	}

	err = types.ValidateAction(kind, value)
	if err != nil {
		return types.Action{}, err
	}

	return types.Action{
		Kind:  types.ActionKind(kind),
		Value: value,
	}, nil
}

func (a ActionColumnIter) Stop() {
	if a.stopKind != nil {
		a.stopKind()
	}
	if a.stopValue != nil {
		a.stopValue()
	}
}
