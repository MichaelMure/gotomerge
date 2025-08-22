package column

import (
	"iter"
	"math"

	"gotomerge/column/rle"
	"gotomerge/types"
)

type ObjectColumnIter struct {
	nextActorIdx func() (rle.NullableValue[uint64], error, bool)
	stopActorIdx func()
	nextCounter  func() (rle.NullableValue[uint64], error, bool)
	stopCounter  func()
}

func ObjectColumn(actorIdxs ActorColumnIter, counters UlebColumnIter) ObjectColumnIter {
	var res ObjectColumnIter
	if actorIdxs != nil {
		res.nextActorIdx, res.stopActorIdx = iter.Pull2(actorIdxs)
	}
	if counters != nil {
		res.nextCounter, res.stopCounter = iter.Pull2(counters)
	}
	return res
}

func (o ObjectColumnIter) Next() (types.ObjectId, error) {
	actorIdx, nullActorIdx, err := extract(o.nextActorIdx)
	if err != nil {
		return types.ObjectId{}, err
	}

	counter, nullCounter, err := extract(o.nextCounter)
	if err != nil {
		return types.ObjectId{}, err
	}

	switch {
	case !nullActorIdx && nullCounter:
		return types.ObjectId{}, ErrUnexpectedNull("counter")
	case nullActorIdx && !nullCounter:
		return types.ObjectId{}, ErrUnexpectedNull("actor index")
	case nullActorIdx && nullCounter:
		return types.RootObjectId(), nil
	default:
		if actorIdx < 0 || actorIdx >= math.MaxUint32 {
			return types.ObjectId{}, ErrOutOfRange("actor index")
		}
		if counter < 0 || counter >= math.MaxUint32 {
			return types.ObjectId{}, ErrOutOfRange("counter")
		}
		return types.ObjectId{ActorIdx: uint32(actorIdx), Counter: uint32(counter)}, nil
	}
}

func (o ObjectColumnIter) Stop() {
	if o.stopActorIdx != nil {
		o.stopActorIdx()
	}
	if o.stopCounter != nil {
		o.stopCounter()
	}
}

// extract is a helper that pulls values from a possibly nil iterator that produces
// nullable values and errors.
// - iterator is nil --> return infinite null values
// - iterator produces an error --> return the error
// - iterator produces a value --> return the value, and isNull == false
// - iterator produces a null value --> return the zero value, and isNull == true
// - iterator is done --> return ErrDone
func extract[T any](fn func() (rle.NullableValue[T], error, bool)) (val T, isNull bool, err error) {
	if fn == nil {
		return *new(T), true, nil
	}
	nullable, err, ok := fn()
	if !ok {
		return *new(T), true, ErrDone
	}
	if err != nil {
		return *new(T), true, err
	}
	val, valid := nullable.Value()
	return val, !valid, nil
}
