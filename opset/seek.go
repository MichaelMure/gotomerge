package opset

import (
	"errors"
	"sort"

	"gotomerge/column"
)

const seekStride = 64

// seekPoint is a combined checkpoint for all snapshot query columns,
// positioned immediately before op opIdx is read.
type seekPoint struct {
	opIdx  uint32
	key    *column.KeyReader
	opId   *column.OpIdReader
	action *column.ActionReader
}

// seekIndex is a sorted list of combined column checkpoints, spaced roughly
// seekStride operations apart. The first entry always has opIdx == 0.
type seekIndex []seekPoint

// seek returns the checkpoint at or before targetOp and the number of ops to
// skip from that checkpoint to reach targetOp. The remaining count is at most
// seekStride.
func (idx seekIndex) seek(targetOp uint32) (seekPoint, uint32) {
	if len(idx) == 0 {
		return seekPoint{}, targetOp
	}
	i := sort.Search(len(idx), func(i int) bool { return idx[i].opIdx > targetOp }) - 1
	if i < 0 {
		i = 0
	}
	pt := idx[i]
	return pt, targetOp - pt.opIdx
}

// buildSeekIndex scans all query columns in lockstep and records a checkpoint
// every seekStride operations.
func buildSeekIndex(ss *snapshotStore) seekIndex {
	if ss.opCount == 0 {
		return nil
	}

	keyActor, e1 := column.Opt(ss.keyActorId, column.NewActorReader)
	keyCounter, e2 := column.Opt(ss.keyCounter, column.NewDeltaReader)
	keyString, e3 := column.Opt(ss.keyString, column.NewStringReader)
	opActor, e4 := column.Opt(ss.opActorId, column.NewActorReader)
	opCounter, e5 := column.Opt(ss.opCounter, column.NewDeltaReader)
	action, e6 := column.Opt(ss.action, column.NewUlebReader)
	valueMeta, e7 := column.Opt(ss.valueMetadata, column.NewValueMetadataReader)
	value, e8 := column.Opt(ss.value, column.NewValueReader)
	if err := errors.Join(e1, e2, e3, e4, e5, e6, e7, e8); err != nil {
		return nil
	}

	keyReader := column.NewKeyReader(keyActor, keyCounter, keyString)
	opIdReader := column.NewOpIdReader(opActor, opCounter)
	actionReader := column.NewActionReader(action, valueMeta, value)

	// Checkpoint at op 0.
	keyFork0, _ := keyReader.Fork()
	opIdFork0, _ := opIdReader.Fork()
	actionFork0, _ := actionReader.Fork()

	idx := seekIndex{seekPoint{
		opIdx:  0,
		key:    keyFork0,
		opId:   opIdFork0,
		action: actionFork0,
	}}

	for i := uint32(0); i < ss.opCount; i++ {
		keyReader.Next()    //nolint:errcheck
		opIdReader.Next()   //nolint:errcheck
		actionReader.Next() //nolint:errcheck

		next := i + 1
		if next < ss.opCount && next%seekStride == 0 {
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
			idx = append(idx, seekPoint{
				opIdx:  next,
				key:    keyFork,
				opId:   opIdFork,
				action: actionFork,
			})
		}
	}

	return idx
}
