package docproxy

import (
	"iter"

	"github.com/MichaelMure/gotomerge/opset"
	"github.com/MichaelMure/gotomerge/types"
)

// MapView is a view of a map object inside a document. It is returned by
// [Document.Map], [Txn.Map], or via a type assertion on a [Value] from [Document.Values].
//
// When obtained from a [Txn] (txn != nil) it is in write mode: [MapView.Map],
// [MapView.List], [MapView.Set], and [MapView.Delete] are available.
// When obtained from [Document] directly (txn == nil) it is read-only:
// write methods panic with a clear message.
type MapView struct {
	s   *opset.OpSet
	txn *opset.Transaction // nil if read-only
	obj types.ObjectId
}

func (MapView) isValue() {}

// -- Read methods ------------------------------------------------------------

// Keys returns all keys that have a live value in this map.
func (mv MapView) Keys() []string {
	return mv.s.MapKeys(mv.obj)
}

// Values returns an iterator over all live key/value pairs in this map.
func (mv MapView) Values() iter.Seq2[string, Value] {
	return mapValues(mv.s, mv.txn, mv.obj)
}

// Get returns the value at key. Returns false if key is absent.
// The concrete type of the returned [Value] reflects the actual data kind.
func (mv MapView) Get(key string) (Value, bool) {
	return mapGet(mv.s, mv.txn, mv.obj, key)
}

// Native implements [Value]. Returns a map[string]any with all live key/value
// pairs recursively materialised.
func (mv MapView) Native() any {
	return materializeObj(mv.s, mv.obj, types.ActionMakeMap)
}

// -- Write-chain methods (create-or-get) -------------------------------------
//
// These methods are intended for write chaining within a [Txn] callback:
//
//	txn.Map("config").Map("nested").Set("key", value)
//
// When txn is nil (read-only context) they panic if the key is absent or holds
// a wrong type. For safe reads use [MapView.Get] or [MapView.Values] instead.

// Map returns a [MapView] for the map at key. In write mode, creates the map if
// absent. Panics with [ErrType] on a type mismatch.
func (mv MapView) Map(key string) MapView {
	if mv.txn != nil {
		obj := getOrMakeMap(mv.s, mv.txn, mv.obj, key)
		return MapView{s: mv.s, txn: mv.txn, obj: obj}
	}
	op, ok := mv.s.MapGet(mv.obj, key)
	if !ok || op.Action.Kind != types.ActionMakeMap {
		panic(ErrType{expected: MapView{}, got: nil})
	}
	return MapView{s: mv.s, obj: types.ObjectId(op.Id)}
}

// Text returns a [TextView] for the text object at key. In write mode, creates
// the text object if absent. Panics with [ErrType] on a type mismatch.
func (mv MapView) Text(key string) TextView {
	if mv.txn != nil {
		obj := getOrMakeText(mv.s, mv.txn, mv.obj, key)
		return TextView{s: mv.s, txn: mv.txn, obj: obj}
	}
	op, ok := mv.s.MapGet(mv.obj, key)
	if !ok || op.Action.Kind != types.ActionMakeText {
		panic(ErrType{expected: TextView{}, got: nil})
	}
	return TextView{s: mv.s, obj: types.ObjectId(op.Id)}
}

// List returns a [ListView] for the list at key. In write mode, creates the list
// if absent. Panics with [ErrType] on a type mismatch.
func (mv MapView) List(key string) ListView {
	if mv.txn != nil {
		obj := getOrMakeList(mv.s, mv.txn, mv.obj, key)
		return ListView{s: mv.s, txn: mv.txn, obj: obj}
	}
	op, ok := mv.s.MapGet(mv.obj, key)
	if !ok || op.Action.Kind != types.ActionMakeList {
		panic(ErrType{expected: ListView{}, got: nil})
	}
	return ListView{s: mv.s, obj: types.ObjectId(op.Id)}
}

// -- Write-only methods ------------------------------------------------------

// Set sets key to a scalar value. value must be one of: bool, int64, float64,
// string, []byte, or nil. Panics if this view is read-only.
func (mv MapView) Set(key string, value any) {
	mv.mustWrite()
	mv.txn.MapSet(mv.obj, key, value)
}

// Delete removes key from this map. No-op if key is absent.
// Panics if this view is read-only.
func (mv MapView) Delete(key string) {
	mv.mustWrite()
	mv.txn.MapDelete(mv.obj, key)
}

// Increment adds delta to the counter at key in this map.
// The key must already hold a counter value. No-op if the key is absent.
// Panics if this view is read-only.
func (mv MapView) Increment(key string, delta int64) {
	mv.mustWrite()
	mv.txn.MapIncrement(mv.obj, key, delta)
}

func (mv MapView) mustWrite() {
	if mv.txn == nil {
		panic("write operation on read-only MapView (obtain via Txn.Map)")
	}
}

// -- shared helpers ----------------------------------------------------------

func getOrMakeMap(s *opset.OpSet, txn *opset.Transaction, obj types.ObjectId, key string) types.ObjectId {
	if op, ok := s.MapGet(obj, key); ok {
		if op.Action.Kind != types.ActionMakeMap {
			panic(ErrType{expected: MapView{}, got: op.Action.Value})
		}
		return types.ObjectId(op.Id)
	}
	newObj, _ := txn.MakeMap(obj, key)
	return newObj
}

func getOrMakeList(s *opset.OpSet, txn *opset.Transaction, obj types.ObjectId, key string) types.ObjectId {
	if op, ok := s.MapGet(obj, key); ok {
		if op.Action.Kind != types.ActionMakeList {
			panic(ErrType{expected: ListView{}, got: op.Action.Value})
		}
		return types.ObjectId(op.Id)
	}
	newObj, _ := txn.MakeList(obj, key)
	return newObj
}

func getOrMakeText(s *opset.OpSet, txn *opset.Transaction, obj types.ObjectId, key string) types.ObjectId {
	if op, ok := s.MapGet(obj, key); ok {
		if op.Action.Kind != types.ActionMakeText {
			panic(ErrType{expected: TextView{}, got: op.Action.Value})
		}
		return types.ObjectId(op.Id)
	}
	newObj, _ := txn.MakeText(obj, key)
	return newObj
}
