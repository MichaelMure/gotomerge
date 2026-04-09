package docproxy

import (
	"gotomerge/opset"
	"gotomerge/types"
)

// BoolView is a view of a boolean value. Obtained via a type assertion on a
// [Value] from [MapView.Get] or [Document.Values].
//
// Write methods ([BoolView.Set], [BoolView.Toggle]) are only available when the
// view was obtained in the context of a [Txn].
type BoolView struct {
	s      *opset.OpSet
	txn    *opset.Transaction // nil if read-only
	obj    types.ObjectId     // containing map
	key    string
	scalar bool // cached value from the OpSet at the time of creation
}

func (BoolView) isValue() {}

// Value returns the current committed boolean value.
func (b BoolView) Value() bool {
	return b.scalar
}

// Set sets the boolean to val. Panics if read-only.
func (b BoolView) Set(val bool) {
	b.mustWrite()
	b.txn.MapSet(b.obj, b.key, val)
}

// Toggle flips the current boolean value. Panics if read-only.
func (b BoolView) Toggle() {
	b.Set(!b.scalar)
}

func (b BoolView) mustWrite() {
	if b.txn == nil {
		panic("write operation on read-only BoolView (obtain via Txn)")
	}
}

// StringView is a view of a string value. Obtained via a type assertion on a
// [Value] from [MapView.Get] or [Document.Values].
//
// Also used to represent text objects (which are serialised as plain strings
// at the document API level).
//
// Write methods are only available when the view was obtained in the context
// of a [Txn].
type StringView struct {
	s      *opset.OpSet
	txn    *opset.Transaction // nil if read-only
	obj    types.ObjectId     // containing map (zero for text objects)
	key    string
	scalar string // cached value
}

func (StringView) isValue() {}

// Value returns the current committed string value.
func (sv StringView) Value() string {
	return sv.scalar
}

// Set sets the string to val. Panics if read-only or if the key currently holds
// a text object (text objects require splice operations, not a plain set).
func (sv StringView) Set(val string) {
	sv.mustWrite()
	sv.txn.MapSet(sv.obj, sv.key, val)
}

func (sv StringView) mustWrite() {
	if sv.txn == nil {
		panic("write operation on read-only StringView (obtain via Txn)")
	}
}
