package gotomerge

import (
	"fmt"
	"iter"

	"gotomerge/column"
	"gotomerge/column/rle"
	"gotomerge/types"
	iterutil "gotomerge/utils/iterator"
)

// type Operation interface {
// 	Id() OperationId
// 	Actor() types.ActorId
// 	Action() types.Action
// }

func ObjectIterator(mainActor types.ActorId, otherActors []types.ActorId,
	actors column.ActorColumnIter, counters column.UlebColumnIter) iter.Seq2[types.ObjectId, error] {
	return iterutil.PullPair(actors, counters, func(actorVal rle.NullableValue[uint64], counterVal rle.NullableValue[uint64]) (types.ObjectId, error) {
		actor, ok1 := actorVal.Value()
		counter, ok2 := counterVal.Value()
		if ok1 != ok2 {
			return types.ObjectId{}, fmt.Errorf("actor and counter must be either both null or both valid")
		}
		if !ok1 {
			return types.NullObjectId(), nil
		}
		if actor == 0 {
			return types.ObjectId{Actor: mainActor, Counter: counter}, nil
		}
		if actor > uint64(len(otherActors)) {
			return types.ObjectId{}, fmt.Errorf("actor id out of range")
		}
		return types.ObjectId{Actor: otherActors[actor-1], Counter: counter}, nil
	})
}

// KeyId is either:
// - ActorId + Counter
// - String
type KeyId struct {
	Actor   types.ActorId
	Counter uint64

	String string
}

type OperationId struct { // == lamport timestamp
	Actor   types.ActorId
	Counter uint64
}

type Operation struct {
	Object     types.ObjectId
	Key        KeyId
	Id         OperationId
	Insert     bool
	Action     types.Action
	Value      any
	Successors []Operation
}
