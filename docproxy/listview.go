package docproxy

import (
	"iter"

	"gotomerge/opset"
	"gotomerge/types"
)

// ListView is a view of a list object inside a document. It is returned by
// [Document.List], [Txn.List], or via a type assertion on a [Value] from [Document.Values].
//
// Write methods are only available when the view was obtained from a [Txn].
type ListView struct {
	s   *opset.OpSet
	txn *opset.Transaction // nil if read-only
	obj types.ObjectId
}

func (ListView) isValue() {}

// Len returns the number of live (non-deleted) elements.
func (lv ListView) Len() int {
	return len(lv.s.ListElements(lv.obj))
}

// Get returns the element at index idx. Returns false if idx is out of range.
// The concrete type of the returned [Value] reflects the actual data kind.
func (lv ListView) Get(idx int) (Value, bool) {
	op, ok := lv.s.ListGet(lv.obj, idx)
	if !ok {
		return nil, false
	}
	return wrapOp(lv.s, lv.txn, op), true
}

// Values returns an iterator over the live elements in order.
func (lv ListView) Values() iter.Seq2[int, Value] {
	return func(yield func(int, Value) bool) {
		ops := lv.s.ListElements(lv.obj)
		for i, op := range ops {
			if !yield(i, wrapOp(lv.s, lv.txn, op)) {
				return
			}
		}
	}
}

// Append inserts value at the end of the list. value must be a scalar:
// bool, int64, float64, string, []byte, or nil. Panics if read-only.
func (lv ListView) Append(value any) {
	lv.mustWrite()
	ops := lv.txn.ListElements(lv.obj)
	var pred types.Key
	if len(ops) == 0 {
		pred = types.KeyOpId{} // head sentinel
	} else {
		last := ops[len(ops)-1]
		pred = types.KeyOpId(last.Id)
	}
	lv.txn.ListInsert(lv.obj, pred, value)
}

// Delete removes the element at idx. Panics if read-only or idx is out of range.
func (lv ListView) Delete(idx int) {
	lv.mustWrite()
	ops := lv.txn.ListElements(lv.obj)
	if idx < 0 || idx >= len(ops) {
		panic("listview: index out of range")
	}
	op := ops[idx]
	lv.txn.ListDelete(lv.obj, op.Id, op.Id)
}

func (lv ListView) mustWrite() {
	if lv.txn == nil {
		panic("write operation on read-only ListView (obtain via Txn.List)")
	}
}
