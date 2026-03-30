package opset

import (
	"errors"
	"sort"

	"gotomerge/column"
)

const seekStride = 64

// docSeekPoint is a combined checkpoint for all document query columns.
// All readers are forked at the same op boundary (opIdx): they are positioned
// immediately before op opIdx is read.
type docSeekPoint struct {
	opIdx   uint32
	key     *column.KeyReader
	opId    *column.OpIdReader
	action  *column.ActionReader
}

// docSeekIndex is a sorted list of combined column checkpoints, spaced roughly
// seekStride operations apart. The first entry always has opIdx == 0.
type docSeekIndex []docSeekPoint

// seek returns the checkpoint at or before targetOp and the number of ops to
// skip from that checkpoint to reach targetOp. The remaining count is at most
// seekStride.
func (idx docSeekIndex) seek(targetOp uint32) (docSeekPoint, uint32) {
	if len(idx) == 0 {
		return docSeekPoint{}, targetOp
	}
	// Largest checkpoint with opIdx ≤ targetOp.
	i := sort.Search(len(idx), func(i int) bool { return idx[i].opIdx > targetOp }) - 1
	if i < 0 {
		i = 0
	}
	pt := idx[i]
	return pt, targetOp - pt.opIdx
}

// buildDocSeekIndex scans all query columns in lockstep and records a
// checkpoint every seekStride operations. The result is stored in docStore
// and used by scanDocRange to avoid full-column seeks.
func buildDocSeekIndex(ds *docStore) docSeekIndex {
	if ds.opCount == 0 {
		return nil
	}

	keyActor, e1 := column.Opt(ds.keyActorId, column.NewActorReader)
	keyCounter, e2 := column.Opt(ds.keyCounter, column.NewDeltaReader)
	keyString, e3 := column.Opt(ds.keyString, column.NewStringReader)
	opActor, e4 := column.Opt(ds.opActorId, column.NewActorReader)
	opCounter, e5 := column.Opt(ds.opCounter, column.NewDeltaReader)
	action, e6 := column.Opt(ds.action, column.NewUlebReader)
	valueMeta, e7 := column.Opt(ds.valueMetadata, column.NewValueMetadataReader)
	value, e8 := column.Opt(ds.value, column.NewValueReader)
	if err := errors.Join(e1, e2, e3, e4, e5, e6, e7, e8); err != nil {
		return nil
	}

	keyReader := column.NewKeyReader(keyActor, keyCounter, keyString)
	opIdReader := column.NewOpIdReader(opActor, opCounter)
	actionReader := column.NewActionReader(action, valueMeta, value)

	// Checkpoint at op 0: fork from the initial position.
	keyFork0, _ := keyReader.Fork()
	opIdFork0, _ := opIdReader.Fork()
	actionFork0, _ := actionReader.Fork()

	idx := docSeekIndex{docSeekPoint{
		opIdx:  0,
		key:    keyFork0,
		opId:   opIdFork0,
		action: actionFork0,
	}}

	for i := uint32(0); i < ds.opCount; i++ {
		// Advance all readers by one operation.
		keyReader.Next()   //nolint:errcheck
		opIdReader.Next()  //nolint:errcheck
		actionReader.Next() //nolint:errcheck

		// After processing op i, record a checkpoint for op i+1 if it falls
		// on a stride boundary.
		next := i + 1
		if next < ds.opCount && next%seekStride == 0 {
			keyFork, err := keyReader.Fork()
			if err != nil {
				continue
			}
			opIdFork, err := opIdReader.Fork()
			if err != nil {
				continue
			}
			actionFork, err := actionReader.Fork()
			if err != nil {
				continue
			}
			idx = append(idx, docSeekPoint{
				opIdx:  next,
				key:    keyFork,
				opId:   opIdFork,
				action: actionFork,
			})
		}
	}

	return idx
}
