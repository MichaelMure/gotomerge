package opset

import "github.com/MichaelMure/gotomerge/types"

// MapGet returns the winning live op at the given key in a map object.
// When multiple peers set the same key concurrently (a conflict), the op
// with the highest OpId wins: highest Counter first, then actor bytes as a
// tiebreaker. Returns false if the key has no live value.
//
// For scalar values, the result is in op.Action.Value. For nested objects
// (MakeMap, MakeList, MakeText), op.Id is the ObjectId of the created object —
// cast it with types.ObjectId(op.Id) to query it further.
func (s *OpSet) MapGet(obj types.ObjectId, key string) (Op, bool) {
	var winner *Op

	if s.snapshot != nil {
		if r, ok := s.snapshot.objRanges[obj]; ok {
			s.snapshot.scanRange(r, func(idx uint32, op Op) bool {
				if s.snapshot.succCount[idx] > 0 || op.Insert ||
					op.Action.Kind == types.ActionDelete || op.Action.Kind == types.ActionInc {
					return true
				}
				k, strKey := op.Key.(types.KeyString)
				if !strKey || string(k) != key {
					return true
				}
				if winner == nil || s.opIdGreater(op.Id, winner.Id) {
					winner = &op
				}
				return true
			})
		}
	}

	if s.delta != nil {
		for _, idx := range s.delta.byObj[obj] {
			op := &s.delta.ops[idx]
			if op.SuccCount > 0 || op.Insert ||
				op.Action.Kind == types.ActionDelete || op.Action.Kind == types.ActionInc {
				continue
			}
			k, strKey := op.Key.(types.KeyString)
			if !strKey || string(k) != key {
				continue
			}
			if winner == nil || s.opIdGreater(op.Id, winner.Id) {
				winner = op
			}
		}
	}

	if winner == nil {
		return Op{}, false
	}
	return s.applyCounterDelta(*winner), true
}

// MapGetAll returns all live values at the given key in a map object.
// Under normal operation this slice has one element. Multiple elements
// indicate a conflict: the same key was set concurrently by different peers.
// The caller can resolve it by picking one (e.g. the first, which has the
// highest OpId and matches what MapGet returns) or by presenting all to the user.
func (s *OpSet) MapGetAll(obj types.ObjectId, key string) []Op {
	var live []Op

	if s.snapshot != nil {
		if r, ok := s.snapshot.objRanges[obj]; ok {
			s.snapshot.scanRange(r, func(idx uint32, op Op) bool {
				if s.snapshot.succCount[idx] > 0 || op.Insert ||
					op.Action.Kind == types.ActionDelete || op.Action.Kind == types.ActionInc {
					return true
				}
				k, strKey := op.Key.(types.KeyString)
				if !strKey || string(k) != key {
					return true
				}
				live = append(live, op)
				return true
			})
		}
	}

	if s.delta != nil {
		for _, idx := range s.delta.byObj[obj] {
			op := s.delta.ops[idx]
			if op.SuccCount > 0 || op.Insert ||
				op.Action.Kind == types.ActionDelete || op.Action.Kind == types.ActionInc {
				continue
			}
			k, strKey := op.Key.(types.KeyString)
			if !strKey || string(k) != key {
				continue
			}
			live = append(live, op)
		}
	}

	sortOpsDesc(live, s)
	for i := range live {
		live[i] = s.applyCounterDelta(live[i])
	}
	return live
}

// MapKeys returns all map keys that have at least one live value in obj.
func (s *OpSet) MapKeys(obj types.ObjectId) []string {
	seen := make(map[string]struct{})
	var keys []string

	if s.snapshot != nil {
		if r, ok := s.snapshot.objRanges[obj]; ok {
			s.snapshot.scanRange(r, func(idx uint32, op Op) bool {
				if s.snapshot.succCount[idx] > 0 || op.Insert ||
					op.Action.Kind == types.ActionDelete || op.Action.Kind == types.ActionInc {
					return true
				}
				k, strKey := op.Key.(types.KeyString)
				if !strKey {
					return true
				}
				if _, exists := seen[string(k)]; !exists {
					seen[string(k)] = struct{}{}
					keys = append(keys, string(k))
				}
				return true
			})
		}
	}

	if s.delta != nil {
		for _, idx := range s.delta.byObj[obj] {
			op := &s.delta.ops[idx]
			if op.SuccCount > 0 || op.Insert ||
				op.Action.Kind == types.ActionDelete || op.Action.Kind == types.ActionInc {
				continue
			}
			k, strKey := op.Key.(types.KeyString)
			if !strKey {
				continue
			}
			if _, exists := seen[string(k)]; !exists {
				seen[string(k)] = struct{}{}
				keys = append(keys, string(k))
			}
		}
	}

	return keys
}

// MapItems returns one winning Op per live key in obj, in insertion order.
// When multiple peers set the same key concurrently (a conflict), the op with
// the highest OpId wins — matching MapGet's conflict resolution.
// This is a single O(n) pass; prefer it over MapKeys + MapGet when iterating
// all entries.
func (s *OpSet) MapItems(obj types.ObjectId) []Op {
	type entry struct {
		op  Op
		pos int
	}
	best := make(map[string]entry)
	var order []string

	add := func(op Op, succCount uint32) {
		if succCount > 0 || op.Insert ||
			op.Action.Kind == types.ActionDelete || op.Action.Kind == types.ActionInc {
			return
		}
		k, strKey := op.Key.(types.KeyString)
		if !strKey {
			return
		}
		key := string(k)
		if e, exists := best[key]; !exists {
			best[key] = entry{op: op, pos: len(order)}
			order = append(order, key)
		} else if s.opIdGreater(op.Id, e.op.Id) {
			e.op = op
			best[key] = e
		}
	}

	if s.snapshot != nil {
		if r, ok := s.snapshot.objRanges[obj]; ok {
			s.snapshot.scanRange(r, func(idx uint32, op Op) bool {
				add(op, s.snapshot.succCount[idx])
				return true
			})
		}
	}

	if s.delta != nil {
		for _, idx := range s.delta.byObj[obj] {
			op := s.delta.ops[idx]
			add(op, op.SuccCount)
		}
	}

	result := make([]Op, len(order))
	for i, key := range order {
		result[i] = s.applyCounterDelta(best[key].op)
	}
	return result
}
