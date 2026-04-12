package opset

import (
	"github.com/MichaelMure/gotomerge/types"
	"github.com/MichaelMure/gotomerge/utils/treap"
)

// MakeList creates a list at obj[key] and returns its ObjectId and OpId.
func (t *Transaction) MakeList(obj types.ObjectId, key string) (types.ObjectId, types.OpId) {
	id := t.nextOpId()
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionMakeList},
		preds:  opsToIds(t.s.MapGetAll(obj, key)),
	})
	return types.ObjectId(id), id
}

// MakeText creates a text object at obj[key] and returns its ObjectId and OpId.
func (t *Transaction) MakeText(obj types.ObjectId, key string) (types.ObjectId, types.OpId) {
	id := t.nextOpId()
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionMakeText},
		preds:  opsToIds(t.s.MapGetAll(obj, key)),
	})
	return types.ObjectId(id), id
}

// ListInsert inserts a new element after pred in list obj.
// Use types.KeyOpId{} (zero Counter = head sentinel) to insert at the front.
func (t *Transaction) ListInsert(obj types.ObjectId, pred types.Key, value any) types.OpId {
	id := t.nextOpId()
	action := types.Action{Kind: types.ActionSet, Value: value}
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    pred,
		insert: true,
		action: action,
	})
	op := Op{Id: id, Object: obj, Key: pred, Insert: true, Action: action}
	t.getOrInitList(obj).insert(id, op, pred)
	return id
}

// ListDelete marks element posId as deleted. liveOpId is the OpId of the
// current live value at that position (use posId when the position has never
// been updated by a separate op).
func (t *Transaction) ListDelete(obj types.ObjectId, posId types.OpId, liveOpId types.OpId) {
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyOpId(posId),
		action: types.Action{Kind: types.ActionDelete},
		preds:  []types.OpId{liveOpId},
	})
	t.getOrInitList(obj).remove(liveOpId)
}

// ListElements returns the live elements of a list or text object in order,
// including any operations buffered in this transaction that have not yet been
// committed. Use this inside a transaction instead of OpSet.ListElements when
// you need write-your-own-reads (e.g. multiple splices in one Change).
func (t *Transaction) ListElements(obj types.ObjectId) []Op {
	if t.lists != nil {
		if wl, ok := t.lists[obj]; ok {
			return wl.elements()
		}
	}
	// No pending list ops for this object — return committed state directly.
	return t.s.ListElements(obj)
}

// ListTail returns the OpId of the last live element of obj, or false if the
// list is empty. It is O(1) once the working list has been initialised by a
// prior ListInsert or ListDelete call; otherwise it initialises from committed
// state in O(n) and caches it for subsequent calls.
func (t *Transaction) ListTail(obj types.ObjectId) (types.OpId, bool) {
	return t.getOrInitList(obj).tailId()
}

// ListLen returns the number of live elements of obj in the working list.
func (t *Transaction) ListLen(obj types.ObjectId) int {
	return t.getOrInitList(obj).r.Len()
}

// ListAt returns the Op at 0-based position i in the working list, or false if
// i is out of range. O(log n).
func (t *Transaction) ListAt(obj types.ObjectId, i int) (Op, bool) {
	n := t.getOrInitList(obj).r.At(i)
	if n == nil {
		return Op{}, false
	}
	return n.Value(), true
}

// getOrInitList returns the working list for obj, initialising it from
// committed state the first time it is accessed.
func (t *Transaction) getOrInitList(obj types.ObjectId) *workingList {
	if t.lists == nil {
		t.lists = make(map[types.ObjectId]*workingList)
	}
	wl, ok := t.lists[obj]
	if !ok {
		wl = newWorkingList(t.s.ListElements(obj))
		t.lists[obj] = wl
	}
	return wl
}

// workingList is the live view of a list or text object during a transaction.
// It pairs a [treap.Treap] (for O(log n) positional access and ordered traversal)
// with an id→node map (for O(1) OpId lookup). Neither structure alone is
// sufficient: the treap cannot locate a node by OpId, and the map cannot
// represent order or support positional indexing.
type workingList struct {
	r    *treap.Treap[Op]
	byId map[types.OpId]*treap.Node[Op]
}

func newWorkingList(base []Op) *workingList {
	wl := &workingList{
		r:    treap.New[Op](),
		byId: make(map[types.OpId]*treap.Node[Op], len(base)),
	}
	for _, op := range base {
		wl.byId[op.Id] = wl.r.PushBack(op)
	}
	return wl
}

func (wl *workingList) insert(id types.OpId, op Op, pred types.Key) {
	predKey, isPred := pred.(types.KeyOpId)
	var n *treap.Node[Op]
	if isPred && (predKey.ActorIdx != 0 || predKey.Counter != 0) {
		if at := wl.byId[types.OpId(predKey)]; at != nil {
			n = wl.r.InsertAfter(op, at)
		} else {
			n = wl.r.PushBack(op) // pred not found → append
		}
	} else {
		n = wl.r.PushFront(op) // zero/missing pred → prepend
	}
	wl.byId[id] = n
}

func (wl *workingList) remove(id types.OpId) {
	n, ok := wl.byId[id]
	if !ok {
		return
	}
	wl.r.Remove(n)
	delete(wl.byId, id)
}

func (wl *workingList) elements() []Op {
	result := make([]Op, 0, wl.r.Len())
	for op := range wl.r.All() {
		result = append(result, op)
	}
	return result
}

func (wl *workingList) tailId() (types.OpId, bool) {
	n := wl.r.Back()
	if n == nil {
		return types.OpId{}, false
	}
	return n.Value().Id, true
}
