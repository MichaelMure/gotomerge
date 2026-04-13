package opset

import (
	"github.com/MichaelMure/gotomerge/types"
)

// ListElements returns the live elements of a list or text object in order.
// Each returned Op carries the winning live value at that position.
//
// The result is backed by a lazily built listTreap that is cached on the
// OpSet and updated incrementally as new ops arrive, so repeated calls pay
// only an O(n) treap scan.
func (s *OpSet) ListElements(obj types.ObjectId) []Op {
	lt := s.getOrBuildListTreap(obj)
	result := make([]Op, 0, lt.r.Len())
	for op := range lt.r.All() {
		// Treap nodes store the insert op. For list objects, an update op may be
		// the winner; retrieve it from liveValues if present.
		if lt.liveValues != nil {
			if winner, ok := lt.liveValues[op.Id]; ok {
				result = append(result, winner)
				continue
			}
		}
		result = append(result, op)
	}
	return result
}

// ListLen returns the number of live elements in a list or text object.
// O(1) — reads directly from the cached treap without materializing a slice.
func (s *OpSet) ListLen(obj types.ObjectId) int {
	return s.getOrBuildListTreap(obj).r.Len()
}

// ListGet returns the element at the given zero-based index.
// Returns false if index is out of range.
func (s *OpSet) ListGet(obj types.ObjectId, index int) (Op, bool) {
	lt := s.getOrBuildListTreap(obj)
	node := lt.r.LiveAt(index)
	if node == nil {
		return Op{}, false
	}
	op := node.Value()
	if lt.liveValues != nil {
		if winner, ok := lt.liveValues[op.Id]; ok {
			return winner, true
		}
	}
	return op, true
}

// Text returns the full text of a text object. For text objects the result
// is read from the cached rope in O(n/K) chunk concatenations rather
// than an O(n) treap traversal, where K=128 is the rope chunk size.
func (s *OpSet) Text(obj types.ObjectId) string {
	lt := s.getOrBuildListTreap(obj)
	if lt.rope == nil {
		panic("Text called on non-text object")
	}
	return lt.rope.Text()
}
