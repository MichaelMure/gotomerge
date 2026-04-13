package opset

import (
	"strings"

	"github.com/MichaelMure/gotomerge/types"
)

// ListElements returns the live elements of a list or text object in order.
// Each returned Op carries the winning live value at that position.
//
// A position is live if at least one op at that position (the insert op or a
// later update) has SuccCount==0 and action != Delete. Dead positions
// (tombstones) are skipped in the result but still traversed, so elements
// inserted after a deleted element remain in the list.
//
// The result is backed by a lazily built listTreap that is cached on the
// OpSet and invalidated when new ops arrive for this object, so repeated calls
// pay only an O(n) treap scan rather than a full map-and-DFS reconstruction.
func (s *OpSet) ListElements(obj types.ObjectId) []Op {
	lt := s.getOrBuildListTreap(obj)
	result := make([]Op, 0, lt.r.Len())
	for op := range lt.r.All() {
		result = append(result, op)
	}
	return result
}

// ListGet returns the element at the given zero-based index.
// Returns false if index is out of range.
func (s *OpSet) ListGet(obj types.ObjectId, index int) (Op, bool) {
	elems := s.ListElements(obj)
	if index < 0 || index >= len(elems) {
		return Op{}, false
	}
	return elems[index], true
}

// Text concatenates the string values of all live elements of a text object
// in order. It reads directly from the listTreap cache, avoiding both the
// intermediate []Op slice and the map-and-DFS reconstruction on repeated calls.
func (s *OpSet) Text(obj types.ObjectId) string {
	lt := s.getOrBuildListTreap(obj)
	var b strings.Builder
	b.Grow(lt.r.Len()) // at least one byte per character
	for op := range lt.r.All() {
		if v, ok := op.Action.Value.(string); ok {
			b.WriteString(v)
		}
	}
	return b.String()
}
