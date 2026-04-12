package docproxy

import (
	"github.com/MichaelMure/gotomerge/opset"
	"github.com/MichaelMure/gotomerge/types"
)

// Compile-time assertion that TextView implements Value.
var _ Value = TextView{}

// TextView is a view of a text object — a CRDT sequence of characters with its
// own [types.ObjectId]. Unlike [StringView] (which is a plain scalar), a text
// object supports character-level concurrent editing: each character is a
// separate op, so concurrent inserts from different peers can be merged without
// conflicts.
//
// At the read API level, text is surfaced as a plain Go string via [TextView.Value].
// On the write side, [TextView.Splice] and [TextView.Update] are the two entry
// points. Rich text (marks, formatting ranges) will be added later.
//
// Write methods are only available when the view was obtained in the context of
// a [Txn].
type TextView struct {
	s   *opset.OpSet
	txn *opset.Transaction // nil if read-only
	obj types.ObjectId     // ObjectId of this text object in the OpSet
}

func (TextView) isValue() {}

// Native implements [Value]. Returns the text content as a plain string.
func (tv TextView) Native() any { return tv.s.Text(tv.obj) }

// Value returns the current committed text as a plain string.
func (tv TextView) Value() string {
	return tv.s.Text(tv.obj)
}

// Len returns the number of Unicode codepoints (runes) in the text.
func (tv TextView) Len() int {
	return len(tv.s.ListElements(tv.obj))
}

// Splice replaces del Unicode codepoints starting at position pos with the
// characters in insert. To insert without deleting pass del = 0; to delete
// without inserting pass insert = "". Panics if read-only or if pos/del are
// out of bounds.
//
// pos and del count Unicode codepoints (runes), not bytes. Each rune in the
// text is stored as a separate list element in the underlying CRDT, so
// concurrent splices from different peers compose cleanly.
func (tv TextView) Splice(pos, del int, insert string) {
	tv.mustWrite()
	n := tv.txn.ListLen(tv.obj)
	if pos < 0 || pos > n {
		panic("textview: splice position out of range")
	}
	if del < 0 || pos+del > n {
		panic("textview: splice delete count out of range")
	}

	// Delete 'del' elements starting at pos. ListDelete eagerly removes each
	// node from the working rope, so we always fetch position pos — the rope
	// shifts with each removal, bringing the next target into place.
	for range del {
		op, _ := tv.txn.ListAt(tv.obj, pos)
		tv.txn.ListDelete(tv.obj, op.Id, op.Id)
	}

	// Insert each rune after position pos-1. O(log n) per rune.
	var pred types.Key
	if pos == 0 {
		pred = types.KeyOpId{} // head sentinel: insert at the beginning
	} else {
		prev, _ := tv.txn.ListAt(tv.obj, pos-1)
		pred = types.KeyOpId(prev.Id)
	}
	for _, r := range insert {
		id := tv.txn.ListInsert(tv.obj, pred, string(r))
		pred = types.KeyOpId(id)
	}
}

// Update replaces the entire text content with newText by computing the
// minimal single-splice edit (prefix + suffix kept, middle replaced). This is
// a convenience for when fine-grained cursor positions are unavailable; merge
// quality is lower than calling Splice with real edit positions.
// Panics if read-only.
func (tv TextView) Update(newText string) {
	tv.mustWrite()
	current := []rune(tv.s.Text(tv.obj))
	incoming := []rune(newText)

	// Skip common prefix.
	pre := 0
	for pre < len(current) && pre < len(incoming) && current[pre] == incoming[pre] {
		pre++
	}
	// Skip common suffix (within the differing region).
	suf := 0
	for suf < len(current)-pre && suf < len(incoming)-pre &&
		current[len(current)-1-suf] == incoming[len(incoming)-1-suf] {
		suf++
	}

	del := len(current) - pre - suf
	ins := string(incoming[pre : len(incoming)-suf])
	if del > 0 || len(ins) > 0 {
		tv.Splice(pre, del, ins)
	}
}

func (tv TextView) mustWrite() {
	if tv.txn == nil {
		panic("write operation on read-only TextView (obtain via Txn)")
	}
}
