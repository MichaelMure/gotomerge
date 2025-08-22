package types

type OpId struct { // == lamport timestamp
	ActorIdx uint32
	Counter  uint32
}

func (oi OpId) Previous() OpId {
	return OpId{oi.ActorIdx, oi.Counter - 1}
}

func (oi OpId) Minus(n uint32) OpId {
	return OpId{oi.ActorIdx, oi.Counter - n}
}

func (oi OpId) Next() OpId {
	return OpId{oi.ActorIdx, oi.Counter + 1}
}
