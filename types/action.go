package types

import (
	"fmt"
	"math"
)

// ActionKind encodes the kind of action operated on the data model
type ActionKind byte

const (
	// Creates a new map object
	ActionMakeMap ActionKind = 0x00
	// Sets a key of a map, overwrites an item in a list, inserts an item in a list, or edits text
	ActionSet ActionKind = 0x01
	// Creates a new list object
	ActionMakeList ActionKind = 0x02
	// Unsets a key of a map, or removes an item from a list (reducing its length)
	ActionDelete ActionKind = 0x03
	// Creates a new text object
	ActionMakeText ActionKind = 0x04
	// Increments a counter stored in a map or a list
	ActionInc ActionKind = 0x05
)

type Action struct {
	Kind  ActionKind
	Value any
}

func ValidateAction(action uint64, value any) error {
	if action > math.MaxUint8 {
		return fmt.Errorf("action out of range: %d", action)
	}
	switch ActionKind(action) {
	case ActionMakeMap, ActionSet, ActionMakeList, ActionDelete, ActionMakeText:
		return nil
	case ActionInc:
		switch value.(type) {
		case int64, uint64:
			return nil
		default:
			return fmt.Errorf("invalid type for increment action: %T", value)
		}
	default:
		return fmt.Errorf("unknown action: %v", value)
	}
}
