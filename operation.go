package gotomerge

type Action byte

const (
	// Creates a new map object
	ActionMakeMap Action = 0x00
	// Sets a key of a map, overwrites an item in a list, inserts an item in a list, or edits text
	ActionSet Action = 0x01
	// Creates a new list object
	ActionMakeList Action = 0x02
	// Unsets a key of a map, or removes an item from a list (reducing its length)
	ActionDelete Action = 0x03
	// Creates a new text object
	ActionMakeText Action = 0x04
	// Increments a counter stored in a map or a list
	ActionInc Action = 0x05
)

// TODO: needs its own type?
// TODO: rename to Id ?
type OperationId struct {
	Actor   ActorId
	Counter uint64
}

type Operation interface {
	Id() OperationId
	Actor() ActorId
	Action() Action
}
