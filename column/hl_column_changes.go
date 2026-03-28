package column

import (
	"fmt"
	"iter"

	"gotomerge/column/rle"
)

// RawChangeMeta is a single change's metadata decoded from change columns.
// ActorIdx is an index into the document's actor array.
// Deps contains indices into the document's changes array (not hashes; resolved at a higher layer).
type RawChangeMeta struct {
	ActorIdx uint64
	SeqNum   uint64
	MaxOp    uint64
	Time     *int64  // nil if not present in the column
	Message  *string // nil if not present
	Deps     []uint64
}

// ChangesIter reads change metadata from document chunk change columns simultaneously.
// Follows the same pull-iterator pattern as OperationIdColumnIter, KeyColumnIter, etc.
type ChangesIter struct {
	nextActorIdx func() (rle.NullableValue[uint64], error, bool)
	stopActorIdx func()
	nextSeqNum   func() (rle.NullableValue[int64], error, bool)
	stopSeqNum   func()
	nextMaxOp    func() (rle.NullableValue[int64], error, bool)
	stopMaxOp    func()
	nextTime     func() (rle.NullableValue[int64], error, bool)
	stopTime     func()
	nextMessage  func() (rle.NullableValue[string], error, bool)
	stopMessage  func()
	nextDepGroup func() (rle.NullableValue[uint64], error, bool)
	stopDepGroup func()
	nextDepIdx   func() (rle.NullableValue[int64], error, bool)
	stopDepIdx   func()
}

// NewChangesIter creates a ChangesIter from change column iterators.
// Any iterator may be nil, in which case null/zero values are used for that column.
func NewChangesIter(
	actorIdx ActorColumnIter,
	seqNum DeltaColumnIter,
	maxOp DeltaColumnIter,
	time DeltaColumnIter,
	message StringColumnIter,
	depGroup GroupColumnIter,
	depIdx DeltaColumnIter,
) *ChangesIter {
	c := &ChangesIter{}
	if actorIdx != nil {
		c.nextActorIdx, c.stopActorIdx = iter.Pull2(actorIdx)
	}
	if seqNum != nil {
		c.nextSeqNum, c.stopSeqNum = iter.Pull2(seqNum)
	}
	if maxOp != nil {
		c.nextMaxOp, c.stopMaxOp = iter.Pull2(maxOp)
	}
	if time != nil {
		c.nextTime, c.stopTime = iter.Pull2(time)
	}
	if message != nil {
		c.nextMessage, c.stopMessage = iter.Pull2(message)
	}
	if depGroup != nil {
		c.nextDepGroup, c.stopDepGroup = iter.Pull2(depGroup)
	}
	if depIdx != nil {
		c.nextDepIdx, c.stopDepIdx = iter.Pull2(depIdx)
	}
	return c
}

// Next returns the next change's metadata. Returns ErrDone when iteration is complete.
func (c *ChangesIter) Next() (RawChangeMeta, error) {
	// ActorIdx drives iteration: ErrDone here means no more changes.
	actorIdx, isNull, err := extract(c.nextActorIdx)
	if err != nil {
		return RawChangeMeta{}, err
	}
	if isNull {
		return RawChangeMeta{}, ErrUnexpectedNull("change actor")
	}

	seqNum, isNull, err := extract(c.nextSeqNum)
	if err != nil {
		return RawChangeMeta{}, err
	}
	if isNull {
		return RawChangeMeta{}, ErrUnexpectedNull("sequence number")
	}
	if seqNum < 0 {
		return RawChangeMeta{}, fmt.Errorf("negative sequence number: %d", seqNum)
	}

	maxOp, isNull, err := extract(c.nextMaxOp)
	if err != nil {
		return RawChangeMeta{}, err
	}
	if isNull {
		return RawChangeMeta{}, ErrUnexpectedNull("max op")
	}
	if maxOp < 0 {
		return RawChangeMeta{}, fmt.Errorf("negative max op: %d", maxOp)
	}

	rawTime, timeNull, err := extract(c.nextTime)
	if err != nil {
		return RawChangeMeta{}, err
	}
	var timePtr *int64
	if !timeNull {
		timePtr = &rawTime
	}

	message, msgNull, err := extract(c.nextMessage)
	if err != nil {
		return RawChangeMeta{}, err
	}
	var msgPtr *string
	if !msgNull {
		msgPtr = &message
	}

	depCount, isNull, err := extract(c.nextDepGroup)
	if err != nil {
		return RawChangeMeta{}, err
	}
	if isNull {
		depCount = 0
	}

	deps := make([]uint64, 0, min(depCount, 64))
	for i := uint64(0); i < depCount; i++ {
		depIdx, isNull, err := extract(c.nextDepIdx)
		if err != nil {
			return RawChangeMeta{}, fmt.Errorf("dep index: %w", err)
		}
		if isNull {
			return RawChangeMeta{}, ErrUnexpectedNull("dep index")
		}
		if depIdx < 0 {
			return RawChangeMeta{}, fmt.Errorf("negative dep index: %d", depIdx)
		}
		deps = append(deps, uint64(depIdx))
	}

	return RawChangeMeta{
		ActorIdx: actorIdx,
		SeqNum:   uint64(seqNum),
		MaxOp:    uint64(maxOp),
		Time:     timePtr,
		Message:  msgPtr,
		Deps:     deps,
	}, nil
}

// Stop releases resources held by the iterator.
func (c *ChangesIter) Stop() {
	for _, stop := range []func(){
		c.stopActorIdx, c.stopSeqNum, c.stopMaxOp, c.stopTime,
		c.stopMessage, c.stopDepGroup, c.stopDepIdx,
	} {
		if stop != nil {
			stop()
		}
	}
}
