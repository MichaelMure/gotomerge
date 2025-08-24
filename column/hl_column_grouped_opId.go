package column

import (
	"iter"
	"math"

	"gotomerge/column/rle"
	"gotomerge/types"
)

type GroupedOperationIdColumnIter struct {
	column       string
	nextGroup    func() (rle.NullableValue[uint64], error, bool)
	stopGroup    func()
	nextActorIdx func() (rle.NullableValue[uint64], error, bool)
	stopActorIdx func()
	nextCounter  func() (rle.NullableValue[int64], error, bool)
	stopCounter  func()
}

func GroupedOperationIdColumn(column string, group GroupColumnIter, actorIdxs ActorColumnIter, counters DeltaColumnIter) GroupedOperationIdColumnIter {
	var res GroupedOperationIdColumnIter
	res.column = column
	if group != nil {
		res.nextGroup, res.stopGroup = iter.Pull2(group)
	}
	if actorIdxs != nil {
		res.nextActorIdx, res.stopActorIdx = iter.Pull2(actorIdxs)
	}
	if counters != nil {
		res.nextCounter, res.stopCounter = iter.Pull2(counters)
	}
	return res
}

func (o GroupedOperationIdColumnIter) Next() ([]types.OpId, error) {
	count, nullCount, err := extract(o.nextGroup)
	if err != nil {
		// this includes ErrDone
		return nil, err
	}
	if nullCount {
		return nil, ErrUnexpectedNull(o.column + " group")
	}

	// prealloc, but limit the size as we can't trust count
	res := make([]types.OpId, 0, min(count, 100))

	for i := uint64(0); i < count; i++ {
		actorIdx, nullActorIdx, err := extract(o.nextActorIdx)
		if err != nil {
			return nil, err
		}
		counter, nullCounter, err := extract(o.nextCounter)
		if err != nil {
			return nil, err
		}

		switch {
		case nullActorIdx:
			return nil, ErrUnexpectedNull("counter")
		case nullCounter:
			return nil, ErrUnexpectedNull("actor index")
		default:
			if actorIdx < 0 || actorIdx >= math.MaxUint32 {
				return nil, ErrOutOfRange("actor index")
			}
			if counter < 0 || counter >= math.MaxUint32 {
				return nil, ErrOutOfRange("counter")
			}
			res = append(res, types.OpId{ActorIdx: uint32(actorIdx), Counter: uint32(counter)})
		}
	}

	return res, nil
}
