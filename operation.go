package gotomerge

import "gotomerge/types"

// TODO: needs its own type?
// TODO: rename to Id ?
type OperationId struct {
	Actor   types.ActorId
	Counter uint64
}

type Operation interface {
	Id() OperationId
	Actor() types.ActorId
	Action() types.Action
}
