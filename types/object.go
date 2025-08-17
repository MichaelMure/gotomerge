package types

type ObjectId struct {
	Actor   ActorId // if empty, ObjectId is null
	Counter uint64
}

func NullObjectId() ObjectId {
	return ObjectId{}
}

func (o ObjectId) Null() bool {
	return len(o.Actor) == 0
}
