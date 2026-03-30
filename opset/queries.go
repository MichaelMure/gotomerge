package opset

import (
	"bytes"
	"io"

	"gotomerge/types"
)

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

	if s.doc != nil {
		r, ok := s.doc.objRanges[obj]
		if ok {
			scanDocRange(s.doc, r, func(idx uint32, op Op) bool {
				if s.doc.succCnt[idx] > 0 || op.Insert {
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

	for _, idx := range s.changesByObj[obj] {
		op := &s.changes[idx]
		if op.SuccCnt > 0 || op.Insert {
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

	if winner == nil {
		return Op{}, false
	}
	return *winner, true
}

// MapGetAll returns all live values at the given key in a map object.
// Under normal operation this slice has one element. Multiple elements
// indicate a conflict: the same key was set concurrently by different peers.
// The caller can resolve it by picking one (e.g. the first, which has the
// highest OpId and matches what MapGet returns) or by presenting all to the user.
func (s *OpSet) MapGetAll(obj types.ObjectId, key string) []Op {
	var live []Op

	if s.doc != nil {
		r, ok := s.doc.objRanges[obj]
		if ok {
			scanDocRange(s.doc, r, func(idx uint32, op Op) bool {
				if s.doc.succCnt[idx] > 0 || op.Insert {
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

	for _, idx := range s.changesByObj[obj] {
		op := s.changes[idx]
		if op.SuccCnt > 0 || op.Insert {
			continue
		}
		k, strKey := op.Key.(types.KeyString)
		if !strKey || string(k) != key {
			continue
		}
		live = append(live, op)
	}

	sortOpsDesc(live, s)
	return live
}

// MapKeys returns all map keys that have at least one live value in obj.
func (s *OpSet) MapKeys(obj types.ObjectId) []string {
	seen := make(map[string]struct{})
	var keys []string

	if s.doc != nil {
		r, ok := s.doc.objRanges[obj]
		if ok {
			scanDocRange(s.doc, r, func(idx uint32, op Op) bool {
				if s.doc.succCnt[idx] > 0 || op.Insert {
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

	for _, idx := range s.changesByObj[obj] {
		op := &s.changes[idx]
		if op.SuccCnt > 0 || op.Insert {
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

	return keys
}

// ObjType returns the ActionKind of the operation that created obj
// (ActionMakeMap, ActionMakeList, or ActionMakeText). The root object
// always returns ActionMakeMap. Returns false if the object is unknown.
func (s *OpSet) ObjType(obj types.ObjectId) (types.ActionKind, bool) {
	if obj.IsRoot() {
		return types.ActionMakeMap, true
	}
	if s.doc != nil {
		if kind, ok := s.doc.objCreators[obj]; ok {
			return kind, true
		}
	}
	if idx, ok := s.changesById[types.OpId(obj)]; ok {
		return s.changes[idx].Action.Kind, true
	}
	return 0, false
}

// scanDocRange iterates over operations [r.start, r.end) in the docStore,
// calling fn for each with its index and decoded Op. fn returns true to
// continue, false to stop early. Decode errors abort the scan silently.
//
// The seek index lets us jump to within seekStride ops of r.start in O(1)
// and then advance at most seekStride ops to reach the target, instead of
// scanning from the beginning of each column.
func scanDocRange(ds *docStore, r opRange, fn func(idx uint32, op Op) bool) {
	if r.start >= r.end {
		return
	}

	pt, skip := ds.seekIdx.seek(r.start)

	keyIter, err := pt.key.Fork()
	if err != nil {
		return
	}
	opIdIter, err := pt.opId.Fork()
	if err != nil {
		return
	}
	actionIter, err := pt.action.Fork()
	if err != nil {
		return
	}

	// Advance from the checkpoint to r.start (at most seekStride ops).
	for i := uint32(0); i < skip; i++ {
		if _, err := keyIter.Next(); err != nil && err != io.EOF {
			return
		}
		if _, err := opIdIter.Next(); err != nil && err != io.EOF {
			return
		}
		if _, err := actionIter.Next(); err != nil && err != io.EOF {
			return
		}
	}

	for idx := r.start; idx < r.end; idx++ {
		key, err := keyIter.Next()
		if err != nil {
			return
		}
		id, err := opIdIter.Next()
		if err != nil {
			return
		}
		action, err := actionIter.Next()
		if err != nil {
			return
		}
		insert := idx < uint32(len(ds.insert)) && ds.insert[idx]

		op := Op{
			Id:     id,
			Key:    key,
			Insert: insert,
			Action: action,
		}
		if !fn(idx, op) {
			return
		}
	}
}

// opIdGreater reports whether a is greater than b.
// Higher Counter wins; equal counters are broken by actor bytes (higher wins).
func (s *OpSet) opIdGreater(a, b types.OpId) bool {
	if a.Counter != b.Counter {
		return a.Counter > b.Counter
	}
	return bytes.Compare(s.actors[a.ActorIdx], s.actors[b.ActorIdx]) > 0
}

// sortOpsDesc sorts ops in descending OpId order (highest first).
func sortOpsDesc(ops []Op, s *OpSet) {
	for i := 1; i < len(ops); i++ {
		for j := i; j > 0 && s.opIdGreater(ops[j].Id, ops[j-1].Id); j-- {
			ops[j], ops[j-1] = ops[j-1], ops[j]
		}
	}
}

