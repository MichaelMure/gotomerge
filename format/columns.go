package format

import (
	"iter"

	ioutil "gotomerge/utils/io"
)

// OperationColumns holds SubReader references for a set of operations,
// used by both ChangeChunk and DocumentChunk.
//
// In the Automerge binary format, operations are not stored as records. Instead,
// each field of every operation is stored in its own column — a single compressed
// stream for all values of that field across all operations. To reconstruct one
// operation, you read one value from each relevant column in lockstep.
//
// Each field is nil when its column was absent from the binary data
// (meaning all values for that field default to null/zero).
//
// Object (ObjectActorId + ObjectCounter) identifies which map, list, or text
// object the operation targets. The root object is represented as (0, 0).
//
// Key identifies the position within that object:
//   - For maps: KeyString holds the property name; the actor/counter pair is unused.
//   - For lists and text: KeyActorId + KeyCounter identify the list element after
//     which this operation is inserted (i.e. the OpId of its left neighbour).
//     A null key means "insert at the head of the list."
//
// ActorId + Counter together form the operation's own OpId — its globally unique
// identity. In a ChangeChunk these columns are absent; the OpId is derived from
// the change's Actor and StartOp counter instead.
//
// Insert distinguishes an insertion from an assignment at an existing position.
// For maps, Insert is always false. For lists and text, true means a new element
// is being created; false means the operation targets an existing element.
//
// Action encodes what the operation does (set a value, delete, make a map/list/text
// object, increment a counter, etc.), together with ValueMetadata and Value which
// carry the actual scalar value when the action is a set.
//
// Predecessors (PredecessorGroup + PredecessorActorId + PredecessorCounter) list
// the operations that this operation supersedes — the previous value(s) at the
// same position. An operation with no predecessors creates a new value; one with
// predecessors overwrites or deletes an existing one.
//
// Successors (SuccessorGroup + SuccessorActorId + SuccessorCounter) are the
// inverse: operations that later superseded this one. Only present in a
// DocumentChunk — a ChangeChunk does not know its future. An operation with no
// successors is the current live value at its position.
//
// ExpandControl and Mark support rich-text mark operations (bold, italic, etc.)
// and are only relevant for text objects.
type OperationColumns struct {
	ObjectActorId ioutil.SubReader
	ObjectCounter ioutil.SubReader

	KeyActorId ioutil.SubReader
	KeyCounter ioutil.SubReader
	KeyString  ioutil.SubReader

	ActorId ioutil.SubReader
	Counter ioutil.SubReader

	Insert ioutil.SubReader

	Action        ioutil.SubReader
	ValueMetadata ioutil.SubReader
	Value         ioutil.SubReader

	PredecessorGroup   ioutil.SubReader
	PredecessorActorId ioutil.SubReader
	PredecessorCounter ioutil.SubReader

	SuccessorGroup   ioutil.SubReader
	SuccessorActorId ioutil.SubReader
	SuccessorCounter ioutil.SubReader

	ExpandControl ioutil.SubReader

	Mark ioutil.SubReader
}

// ChangeColumns holds SubReader references for the change summary table
// inside a document chunk. Each column stores one field for all changes in a
// compressed, run-length-encoded form.
//
// Dependencies are stored in a two-column group: DependenciesGroup gives the
// number of dependencies for each change, and DependenciesIndex gives the actual
// dependency indices (delta-encoded for compactness). This split is the standard
// "group column" pattern used throughout the Automerge binary format wherever a
// field has a variable-length list per row.
//
// ExtraMetadata / ExtraData are a reserved extension point in the format.
type ChangeColumns struct {
	ActorId ioutil.SubReader
	SeqNum  ioutil.SubReader

	MaxOp ioutil.SubReader

	Time ioutil.SubReader

	Message ioutil.SubReader

	DependenciesGroup ioutil.SubReader
	DependenciesIndex ioutil.SubReader

	ExtraMetadata ioutil.SubReader
	ExtraData     ioutil.SubReader
}

// errSeq returns an iterator that immediately yields err and stops.
func errSeq[T any](err error) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var zero T
		yield(zero, err)
	}
}
