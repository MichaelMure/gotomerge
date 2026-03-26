package types

import "fmt"

type ObjectId OpId

func RootObjectId() ObjectId {
	return ObjectId{ActorIdx: 0, Counter: 0}
}

func (o ObjectId) String() string {
	return fmt.Sprintf("ObjectId(actorIdx: %v, counter: %v)", o.ActorIdx, o.Counter)
}

func (o ObjectId) IsRoot() bool {
	return o.ActorIdx == 0 && o.Counter == 0
}
