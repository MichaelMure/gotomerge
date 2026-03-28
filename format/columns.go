package format

import (
	"gotomerge/column"
)

// OperationColumns holds the decoded column iterators for a set of operations,
// used by both ChangeChunk and DocumentChunk.
//
// In the Automerge binary format, operations are not stored as records. Instead,
// each field of every operation is stored in its own column — a single compressed
// stream for all values of that field across all operations. To reconstruct one
// operation, you read one value from each relevant column in lockstep.
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
	ObjectActorId column.ActorColumnIter
	ObjectCounter column.UlebColumnIter

	KeyActorId column.ActorColumnIter
	KeyCounter column.DeltaColumnIter
	KeyString  column.StringColumnIter

	ActorId column.ActorColumnIter
	Counter column.DeltaColumnIter

	Insert column.BooleanColumnIter

	Action        column.UlebColumnIter
	ValueMetadata column.ValueMetadataColumnIter
	Value         column.ValueColumnIterMaker

	PredecessorGroup   column.GroupColumnIter
	PredecessorActorId column.ActorColumnIter
	PredecessorCounter column.DeltaColumnIter

	SuccessorGroup   column.GroupColumnIter
	SuccessorActorId column.ActorColumnIter
	SuccessorCounter column.DeltaColumnIter

	ExpandControl column.BooleanColumnIter

	Mark column.StringColumnIter
}

// ChangeColumns holds the decoded column iterators for the change summary table
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
	ActorId column.ActorColumnIter
	SeqNum  column.DeltaColumnIter

	MaxOp column.DeltaColumnIter

	Time column.DeltaColumnIter

	Message column.StringColumnIter

	DependenciesGroup column.GroupColumnIter
	DependenciesIndex column.DeltaColumnIter

	ExtraMetadata column.ValueMetadataColumnIter
	ExtraData     column.ValueColumn
}
