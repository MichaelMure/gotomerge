package opset

import "github.com/MichaelMure/gotomerge/types"

// Op is a single operation stored in the OpSet. All actor indices in OpIds
// are relative to the OpSet's own actor table (OpSet.Actor), not to the
// actor table of whichever chunk the op was read from.
//
// SuccCount counts how many successor operations have been applied on top of
// this one. An op is the current live value at its position when SuccCount == 0.
// When SuccCount > 0 the op has been overwritten or deleted.
//
// Insert distinguishes creating a new list element (true) from targeting an
// existing position (false). It is always false for map operations.
type Op struct {
	Id        types.OpId
	Object    types.ObjectId
	Key       types.Key
	Insert    bool
	Action    types.Action
	SuccCount uint32
}
