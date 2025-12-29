package types

type DocChange struct {
	ActorId   uint64
	SeqNum    uint64
	MaxOp     uint64
	Time      Timestamp
	Message   *string // optional
	Deps      []int64
	ExtraData any
}
