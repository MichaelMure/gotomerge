package column

import (
	"iter"
	"math"

	"gotomerge/column/rle"
	"gotomerge/types"
)

type OperationIdColumnIter struct {
	nextActorIdx func() (rle.NullableValue[uint64], error, bool)
	stopActorIdx func()
	nextCounter  func() (rle.NullableValue[int64], error, bool)
	stopCounter  func()
}

func OperationIdColumn(actorIdxs ActorColumnIter, counters DeltaColumnIter) OperationIdColumnIter {
	var res OperationIdColumnIter
	if actorIdxs != nil {
		res.nextActorIdx, res.stopActorIdx = iter.Pull2(actorIdxs)
	}
	if counters != nil {
		res.nextCounter, res.stopCounter = iter.Pull2(counters)
	}
	return res
}

func (o OperationIdColumnIter) Next() (types.OpId, error) {
	actorIdx, nullActorIdx, err := extract(o.nextActorIdx)
	if err != nil {
		return types.OpId{}, err
	}
	counter, nullCounter, err := extract(o.nextCounter)
	if err != nil {
		return types.OpId{}, err
	}

	// TODO: a bit of guessing below, is that correct?
	switch {
	case !nullActorIdx:
		return types.OpId{}, ErrUnexpectedNull("counter")
	case !nullCounter:
		return types.OpId{}, ErrUnexpectedNull("actor index")
	default:
		if actorIdx < 0 || actorIdx >= math.MaxUint32 {
			return types.OpId{}, ErrOutOfRange("actor index")
		}
		if counter < 0 || counter >= math.MaxUint32 {
			return types.OpId{}, ErrOutOfRange("counter")
		}
		return types.OpId{ActorIdx: uint32(actorIdx), Counter: uint32(counter)}, nil
	}
}

func (o OperationIdColumnIter) Stop() {
	if o.stopActorIdx != nil {
		o.stopActorIdx()
	}
	if o.stopCounter != nil {
		o.stopCounter()
	}
}
