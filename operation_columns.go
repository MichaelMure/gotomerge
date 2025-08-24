package gotomerge

import (
	"gotomerge/column"
)

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
