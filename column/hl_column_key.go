package column

import (
	"fmt"
	"iter"
	"math"

	"gotomerge/column/rle"
	"gotomerge/types"
)

type KeyColumnIter struct {
	nextActorIdx func() (rle.NullableValue[uint64], error, bool)
	stopActorIdx func()
	nextCounter  func() (rle.NullableValue[int64], error, bool)
	stopCounter  func()
	nextString   func() (rle.NullableValue[string], error, bool)
	stopString   func()
}

func KeyColumn(actorIdx ActorColumnIter, counter DeltaColumnIter, str StringColumnIter) KeyColumnIter {
	var res KeyColumnIter
	if actorIdx != nil {
		res.nextActorIdx, res.stopActorIdx = iter.Pull2(actorIdx)
	}
	if counter != nil {
		res.nextCounter, res.stopCounter = iter.Pull2(counter)
	}
	if str != nil {
		res.nextString, res.stopString = iter.Pull2(str)
	}
	return res
}

func (k KeyColumnIter) Next() (types.Key, error) {
	actorIdx, nullActorIdx, err := extract(k.nextActorIdx)
	if err != nil {
		return nil, err
	}
	counter, nullCounter, err := extract(k.nextCounter)
	if err != nil {
		return nil, err
	}
	str, nullString, err := extract(k.nextString)
	if err != nil {
		return nil, err
	}

	if (!nullActorIdx || !nullCounter) && !nullString {
		return nil, fmt.Errorf("too many values for key")
	}
	if nullActorIdx && nullCounter && nullString {
		return types.NullKey{}, nil
	}

	switch {
	case !nullString:
		return types.KeyString(str), nil
	case nullCounter:
		return nil, ErrUnexpectedNull("counter")
	case nullActorIdx && counter == 0:
		return types.KeyOpId{ActorIdx: 0, Counter: 0}, nil
	case nullActorIdx:
		return nil, ErrUnexpectedNull("actor index")
	default:
		if actorIdx < 0 || actorIdx >= math.MaxUint32 {
			return nil, ErrOutOfRange("actor index")
		}
		if counter < 0 || counter >= math.MaxUint32 {
			return nil, ErrOutOfRange("counter")
		}
		return types.KeyOpId{ActorIdx: uint32(actorIdx), Counter: uint32(counter)}, nil
	}
}

func (k KeyColumnIter) Stop() {
	if k.stopActorIdx != nil {
		k.stopActorIdx()
	}
	if k.stopCounter != nil {
		k.stopCounter()
	}
	if k.stopString != nil {
		k.stopString()
	}
}
