package types

type ObjectId OpId

func RootObjectId() ObjectId {
	return ObjectId{ActorIdx: 0, Counter: 0}
}

func (o ObjectId) IsRoot() bool {
	return o.ActorIdx == 0 && o.Counter == 0
}
