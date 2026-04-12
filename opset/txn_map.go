package opset

import "github.com/MichaelMure/gotomerge/types"

// MakeMap creates a nested map at obj[key] and returns its ObjectId and OpId.
func (t *Transaction) MakeMap(obj types.ObjectId, key string) (types.ObjectId, types.OpId) {
	id := t.nextOpId()
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionMakeMap},
		preds:  opsToIds(t.s.MapGetAll(obj, key)),
	})
	return types.ObjectId(id), id
}

// MapSet sets obj[key] = value, superseding any existing live values.
func (t *Transaction) MapSet(obj types.ObjectId, key string, value any) types.OpId {
	id := t.nextOpId()
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionSet, Value: value},
		preds:  opsToIds(t.s.MapGetAll(obj, key)),
	})
	return id
}

// MapIncrement adds delta to the counter at obj[key]. The key must currently
// hold a Counter value; if it does not, this is a no-op. Panics if read-only.
func (t *Transaction) MapIncrement(obj types.ObjectId, key string, delta int64) {
	preds := opsToIds(t.s.MapGetAll(obj, key))
	if len(preds) == 0 {
		return
	}
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionInc, Value: delta},
		preds:  preds,
	})
}

// MapDelete deletes obj[key] by superseding all live values there.
// No-ops if the key has no live value.
func (t *Transaction) MapDelete(obj types.ObjectId, key string) {
	preds := opsToIds(t.s.MapGetAll(obj, key))
	if len(preds) == 0 {
		return
	}
	t.ops = append(t.ops, txnOp{
		obj:    obj,
		key:    types.KeyString(key),
		action: types.Action{Kind: types.ActionDelete},
		preds:  preds,
	})
}
