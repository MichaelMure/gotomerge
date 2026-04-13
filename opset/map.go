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

	// mapKeys[obj][key] is a short list of op indices that ever set this key
	// (usually one, more under concurrent writes). We check succCount to skip
	// superseded ops, then decode the column data only for the candidates.
	if s.snapshot != nil {
		for _, idx := range s.snapshot.mapKeys[obj][key] {
			if s.snapshot.succCount[idx] > 0 {
				continue // superseded by a later op
			}
			s.snapshot.scanRange(opRange{start: idx, end: idx + 1}, func(_ uint32, op Op) bool {
				if winner == nil || s.opIdGreater(op.Id, winner.Id) {
					cp := op
					winner = &cp
				}
				return false
			})
		}
	}

	// Delta ops are already decoded in memory; no column seek needed.
	if s.delta != nil {
		for _, idx := range s.delta.mapKeys[obj][key] {
			op := &s.delta.ops[idx]
			if op.SuccCount > 0 {
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
	// For counter ops, fold any accumulated Inc deltas into the returned value.
	return s.applyCounterDelta(*winner), true
}

// MapGetAll returns all live values at the given key in a map object.
// Under normal operation this slice has one element. Multiple elements
// indicate a conflict: the same key was set concurrently by different peers.
// The caller can resolve it by picking one (e.g. the first, which has the
// highest OpId and matches what MapGet returns) or by presenting all to the user.
func (s *OpSet) MapGetAll(obj types.ObjectId, key string) []Op {
	var live []Op

	// Same index-driven lookup as MapGet, but collect all live candidates
	// instead of tracking a single winner. Conflicts are rare so live is
	// almost always length 0 or 1 before the append.
	if s.snapshot != nil {
		for _, idx := range s.snapshot.mapKeys[obj][key] {
			if s.snapshot.succCount[idx] > 0 {
				continue
			}
			s.snapshot.scanRange(opRange{start: idx, end: idx + 1}, func(_ uint32, op Op) bool {
				live = append(live, op)
				return false
			})
		}
	}

	if s.delta != nil {
		for _, idx := range s.delta.mapKeys[obj][key] {
			op := s.delta.ops[idx]
			if op.SuccCount > 0 {
				continue
			}
			live = append(live, op)
		}
	}

	// Sort highest OpId first so index 0 matches what MapGet would return.
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

	// We scan byObj / objRanges in document order rather than iterating the
	// mapKeys index, because Go map iteration is unordered and callers expect
	// keys in insertion order. The full scan is O(n) in ops per object, which
	// is acceptable — MapItems is the right call when you also need values.
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
		pos int // position in order, to reconstruct insertion order at the end
	}
	best := make(map[string]entry) // key → current winning op
	var order []string             // keys in first-seen (insertion) order

	// add considers one op for inclusion. It skips dead, non-value, and
	// non-string-key ops, then either seeds or updates the winner for the key.
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

	// Scan snapshot first (document order), then delta (change order).
	// Because add always keeps the highest OpId, the final winner per key is
	// correct regardless of which store it came from.
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
