package types

// type ObjectId struct {
// 	Actor   ActorId // if empty, ObjectId is null
// 	Counter uint64
// }
//
// func NullObjectId() ObjectId {
// 	return ObjectId{}
// }
//
// func (o ObjectId) Null() bool {
// 	return len(o.Actor) == 0
// }

type ObjectId OpId

func RootObjectId() ObjectId {
	return ObjectId{ActorIdx: 0, Counter: 0}
}

func (o ObjectId) IsRoot() bool {
	return o.ActorIdx == 0 && o.Counter == 0
}

// counter==0 --> root

// methods:
// prev: counter--
// next: counter++
// minus: counter-=n
