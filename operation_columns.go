package gotomerge

import (
	"errors"
	"iter"

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

func (oc OperationColumns) Operations() iter.Seq2[Operation, error] {
	objIter := column.ObjectColumn(oc.ObjectActorId, oc.ObjectCounter)
	keyIter := column.KeyColumn(oc.KeyActorId, oc.KeyCounter, oc.KeyString)
	// only for documents
	OpIdIter := column.OperationIdColumn(oc.ActorId, oc.Counter)
	insertIter := column.InsertColumn(oc.Insert)
	actionIter := column.ActionColumn(oc.Action, oc.ValueMetadata, oc.Value)

	// TODO: predecessor
	// TODO: successor
	// TODO: text formatting

	return func(yield func(Operation, error) bool) {
		defer objIter.Stop()
		defer keyIter.Stop()
		defer OpIdIter.Stop()
		defer insertIter.Stop()
		defer actionIter.Stop()

		for {
			action, errAction := actionIter.Next()
			if errAction != nil && !errors.Is(errAction, column.ErrDone) {
				yield(Operation{}, errAction)
				return
			}

			// Action act as the marker for how long we should iterate
			// (from the rust codebase)
			if errors.Is(errAction, column.ErrDone) {
				return
			}

			obj, err := objIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(Operation{}, err)
				return
			}

			key, err := keyIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(Operation{}, err)
				return
			}

			insert, err := insertIter.Next()
			if err != nil && !errors.Is(err, column.ErrDone) {
				yield(Operation{}, err)
				return
			}

			if !yield(Operation{
				Object: obj,
				Key:    key,
				// Id:         types.ObjectId{},
				Insert:     insert,
				Action:     action,
				Successors: nil,
			}, nil) {
				return
			}
		}
	}
}
