package docproxy

import (
	"fmt"
	"iter"

	"gotomerge/opset"
	"gotomerge/types"
)

// materializeObj recursively converts an object in the OpSet to a plain Go
// value suitable for json.Marshal: map[string]any, []any, string, or scalar.
func materializeObj(s *opset.OpSet, obj types.ObjectId, kind types.ActionKind) any {
	switch kind {
	case types.ActionMakeMap:
		m := make(map[string]any)
		for _, key := range s.MapKeys(obj) {
			op, ok := s.MapGet(obj, key)
			if !ok {
				continue
			}
			m[key] = materializeOp(s, op)
		}
		return m
	case types.ActionMakeList:
		ops := s.ListElements(obj)
		out := make([]any, len(ops))
		for i, op := range ops {
			out[i] = materializeOp(s, op)
		}
		return out
	case types.ActionMakeText:
		return s.Text(obj)
	default:
		return nil
	}
}

func materializeOp(s *opset.OpSet, op opset.Op) any {
	switch op.Action.Kind {
	case types.ActionMakeMap, types.ActionMakeList, types.ActionMakeText:
		childKind, _ := s.ObjType(types.ObjectId(op.Id))
		return materializeObj(s, types.ObjectId(op.Id), childKind)
	default:
		return op.Action.Value
	}
}

// setMapFromJSON recursively walks a map[string]any decoded from JSON and
// issues the corresponding MapSet / MakeMap / MakeList ops on txn.
func setMapFromJSON(txn *opset.Transaction, obj types.ObjectId, m map[string]any) {
	for k, v := range m {
		setValueFromJSON(txn, obj, k, v)
	}
}

func setValueFromJSON(txn *opset.Transaction, obj types.ObjectId, key string, v any) {
	switch v := v.(type) {
	case map[string]any:
		child, _ := txn.MakeMap(obj, key)
		setMapFromJSON(txn, child, v)
	case []any:
		child, _ := txn.MakeList(obj, key)
		setListFromJSON(txn, child, v)
	default:
		txn.MapSet(obj, key, v)
	}
}

func setListFromJSON(txn *opset.Transaction, obj types.ObjectId, items []any) {
	pred := types.Key(types.KeyOpId{}) // head sentinel
	for _, item := range items {
		switch v := item.(type) {
		case map[string]any:
			// TODO: nested objects in lists require MakeMap inside a list —
			// needs ListInsert variant for objects, deferred.
			_ = v
		default:
			id := txn.ListInsert(obj, pred, item)
			pred = types.KeyOpId(id)
		}
	}
}

// Value is a read-only or read-write view of a single value inside a document.
// Value is a sealed interface: only types defined in this package implement it.
type Value interface {
	isValue()
}

// ErrType is returned (or panicked) when a view method is called on a key
// that already holds a value of an incompatible type. [Document.Change] catches
// this panic and returns it as an error.
type ErrType struct {
	expected any
	got      any
}

func (e ErrType) Error() string {
	return fmt.Sprintf("type conflict: expected %T, got %T", e.expected, e.got)
}

// -- internal helpers --------------------------------------------------------

// wrapOp converts a live opset.Op into the appropriate Value view.
func wrapOp(s *opset.OpSet, txn *opset.Transaction, op opset.Op) Value {
	obj := types.ObjectId(op.Id)
	switch op.Action.Kind {
	case types.ActionMakeMap:
		return MapView{s: s, txn: txn, obj: obj}
	case types.ActionMakeList:
		return ListView{s: s, txn: txn, obj: obj}
	case types.ActionMakeText:
		return TextView{s: s, txn: txn, obj: obj}
	default:
		switch v := op.Action.Value.(type) {
		case bool:
			key, _ := op.Key.(types.KeyString)
			return BoolView{s: s, txn: txn, obj: op.Object, key: string(key), scalar: v}
		case string:
			key, _ := op.Key.(types.KeyString)
			return StringView{s: s, txn: txn, obj: op.Object, key: string(key), scalar: v}
		default:
			// TODO: int64, float64, []byte, null, counter, timestamp — add typed views
			key, _ := op.Key.(types.KeyString)
			return StringView{s: s, txn: txn, obj: op.Object, key: string(key), scalar: fmt.Sprint(v)}
		}
	}
}

// mapGet returns the live value at key in obj, wrapping it as a Value.
func mapGet(s *opset.OpSet, txn *opset.Transaction, obj types.ObjectId, key string) (Value, bool) {
	op, ok := s.MapGet(obj, key)
	if !ok {
		return nil, false
	}
	return wrapOp(s, txn, op), true
}

// mapValues returns an iter.Seq2 over all live key/value pairs in obj.
func mapValues(s *opset.OpSet, txn *opset.Transaction, obj types.ObjectId) iter.Seq2[string, Value] {
	return func(yield func(string, Value) bool) {
		for _, op := range s.MapItems(obj) {
			key, _ := op.Key.(types.KeyString)
			if !yield(string(key), wrapOp(s, txn, op)) {
				return
			}
		}
	}
}
