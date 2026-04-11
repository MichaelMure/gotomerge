package docproxy

import (
	"time"

	"github.com/MichaelMure/gotomerge/opset"
	"github.com/MichaelMure/gotomerge/types"
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
func (b BoolView) Value() bool { return b.scalar }

// Native implements [Value].
func (b BoolView) Native() any { return b.scalar }

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
func (sv StringView) Value() string { return sv.scalar }

// Native implements [Value].
func (sv StringView) Native() any { return sv.scalar }

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

// Int64View is a view of an integer value (int64).
type Int64View struct {
	s      *opset.OpSet
	txn    *opset.Transaction
	obj    types.ObjectId
	key    string
	scalar int64
}

func (Int64View) isValue() {}

// Value returns the current committed integer value.
func (v Int64View) Value() int64 { return v.scalar }

// Native implements [Value].
func (v Int64View) Native() any { return v.scalar }

// Set sets the integer to val. Panics if read-only.
func (v Int64View) Set(val int64) {
	v.mustWrite()
	v.txn.MapSet(v.obj, v.key, val)
}

func (v Int64View) mustWrite() {
	if v.txn == nil {
		panic("write operation on read-only Int64View (obtain via Txn)")
	}
}

// Uint64View is a view of an unsigned integer value (uint64).
// These arise when loading documents produced by other Automerge implementations
// that stored unsigned integers.
type Uint64View struct {
	s      *opset.OpSet
	txn    *opset.Transaction
	obj    types.ObjectId
	key    string
	scalar uint64
}

func (Uint64View) isValue() {}

// Value returns the current committed unsigned integer value.
func (v Uint64View) Value() uint64 { return v.scalar }

// Native implements [Value].
func (v Uint64View) Native() any { return v.scalar }

// Set sets the value to val. Panics if read-only.
func (v Uint64View) Set(val uint64) {
	v.mustWrite()
	v.txn.MapSet(v.obj, v.key, val)
}

func (v Uint64View) mustWrite() {
	if v.txn == nil {
		panic("write operation on read-only Uint64View (obtain via Txn)")
	}
}

// Float64View is a view of a floating-point value (float64).
type Float64View struct {
	s      *opset.OpSet
	txn    *opset.Transaction
	obj    types.ObjectId
	key    string
	scalar float64
}

func (Float64View) isValue() {}

// Value returns the current committed float value.
func (v Float64View) Value() float64 { return v.scalar }

// Native implements [Value].
func (v Float64View) Native() any { return v.scalar }

// Set sets the float to val. Panics if read-only.
func (v Float64View) Set(val float64) {
	v.mustWrite()
	v.txn.MapSet(v.obj, v.key, val)
}

func (v Float64View) mustWrite() {
	if v.txn == nil {
		panic("write operation on read-only Float64View (obtain via Txn)")
	}
}

// BytesView is a view of a raw byte-string value ([]byte).
type BytesView struct {
	s      *opset.OpSet
	txn    *opset.Transaction
	obj    types.ObjectId
	key    string
	scalar []byte
}

func (BytesView) isValue() {}

// Value returns the current committed byte value. The returned slice must not
// be modified.
func (v BytesView) Value() []byte { return v.scalar }

// Native implements [Value].
func (v BytesView) Native() any { return v.scalar }

// Set sets the bytes to val. Panics if read-only.
func (v BytesView) Set(val []byte) {
	v.mustWrite()
	v.txn.MapSet(v.obj, v.key, val)
}

func (v BytesView) mustWrite() {
	if v.txn == nil {
		panic("write operation on read-only BytesView (obtain via Txn)")
	}
}

// NullView is a view of a null (nil) value.
type NullView struct{}

func (NullView) isValue()    {}
func (NullView) Native() any { return nil }

// TimestampView is a view of a timestamp value.
type TimestampView struct {
	s      *opset.OpSet
	txn    *opset.Transaction
	obj    types.ObjectId
	key    string
	scalar types.Timestamp
}

func (TimestampView) isValue() {}

// Value returns the current committed timestamp as a [time.Time].
func (v TimestampView) Value() time.Time { return v.scalar.Time() }

// Native implements [Value].
func (v TimestampView) Native() any { return v.scalar.Time() }

// Set sets the timestamp. Panics if read-only.
func (v TimestampView) Set(val time.Time) {
	v.mustWrite()
	v.txn.MapSet(v.obj, v.key, types.FromTime(val))
}

func (v TimestampView) mustWrite() {
	if v.txn == nil {
		panic("write operation on read-only TimestampView (obtain via Txn)")
	}
}

// CounterView is a view of a counter value — an integer that supports
// concurrent increment/decrement from multiple peers without conflicts.
//
// Write methods ([CounterView.Increment]) are only available when the view was
// obtained in the context of a [Txn].
type CounterView struct {
	s      *opset.OpSet
	txn    *opset.Transaction // nil if read-only
	obj    types.ObjectId
	key    string
	scalar types.Counter
}

func (CounterView) isValue() {}

// Value returns the current counter value (base + all applied increments).
func (v CounterView) Value() int64 { return int64(v.scalar) }

// Native implements [Value].
func (v CounterView) Native() any { return int64(v.scalar) }

// Increment adds delta to the counter. Panics if read-only.
func (v CounterView) Increment(delta int64) {
	v.mustWrite()
	v.txn.MapIncrement(v.obj, v.key, delta)
}

func (v CounterView) mustWrite() {
	if v.txn == nil {
		panic("write operation on read-only CounterView (obtain via Txn)")
	}
}
