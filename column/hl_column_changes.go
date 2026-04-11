package column

import (
	"fmt"
	"io"

	"github.com/MichaelMure/gotomerge/column/rle"
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

// ChangesReader is a stateful reader for change columns.
type ChangesReader struct {
	actor     *ActorReader
	seqNum    *DeltaReader
	maxOp     *DeltaReader
	time      *DeltaReader
	message   *StringReader
	depsGroup *GroupReader
	depsIndex *DeltaReader
}

func NewChangesReader(
	actor *ActorReader,
	seqNum *DeltaReader,
	maxOp *DeltaReader,
	time *DeltaReader,
	message *StringReader,
	depsGroup *GroupReader,
	depsIndex *DeltaReader,
) *ChangesReader {
	return &ChangesReader{
		actor:     actor,
		seqNum:    seqNum,
		maxOp:     maxOp,
		time:      time,
		message:   message,
		depsGroup: depsGroup,
		depsIndex: depsIndex,
	}
}

func (c *ChangesReader) nextUint64(r *ActorReader) (uint64, bool, error) {
	if r == nil {
		return 0, true, nil
	}
	nv, err := r.Next()
	if err != nil {
		return 0, true, err
	}
	v, valid := nv.Value()
	return v, !valid, nil
}

func (c *ChangesReader) nextInt64(r *DeltaReader) (int64, bool, error) {
	if r == nil {
		return 0, true, nil
	}
	nv, err := r.Next()
	if err != nil {
		return 0, true, err
	}
	v, valid := nv.Value()
	return v, !valid, nil
}

func (c *ChangesReader) nextString(r *StringReader) (string, bool, error) {
	if r == nil {
		return "", true, nil
	}
	nv, err := r.Next()
	if err != nil {
		return "", true, err
	}
	v, valid := nv.Value()
	return v, !valid, nil
}

func (c *ChangesReader) nextGroup(r *GroupReader) (uint64, bool, error) {
	if r == nil {
		return 0, true, nil
	}
	nv, err := r.Next()
	if err != nil {
		return 0, true, err
	}
	v, valid := nv.Value()
	return v, !valid, nil
}

func (c *ChangesReader) Next() (RawChangeMeta, error) {
	actorIdx, isNull, err := c.nextUint64(c.actor)
	if err != nil {
		return RawChangeMeta{}, err
	}
	if isNull {
		return RawChangeMeta{}, ErrUnexpectedNull("change actor")
	}

	seqNum, isNull, err := c.nextInt64(c.seqNum)
	if err != nil {
		return RawChangeMeta{}, err
	}
	if isNull {
		return RawChangeMeta{}, ErrUnexpectedNull("sequence number")
	}
	if seqNum < 0 {
		return RawChangeMeta{}, fmt.Errorf("negative sequence number: %d", seqNum)
	}

	maxOp, isNull, err := c.nextInt64(c.maxOp)
	if err != nil {
		return RawChangeMeta{}, err
	}
	if isNull {
		return RawChangeMeta{}, ErrUnexpectedNull("max op")
	}
	if maxOp < 0 {
		return RawChangeMeta{}, fmt.Errorf("negative max op: %d", maxOp)
	}

	rawTime, timeNull, err := c.nextInt64(c.time)
	if err != nil {
		return RawChangeMeta{}, err
	}
	var timePtr *int64
	if !timeNull {
		timePtr = &rawTime
	}

	message, msgNull, err := c.nextString(c.message)
	if err != nil {
		return RawChangeMeta{}, err
	}
	var msgPtr *string
	if !msgNull {
		msgPtr = &message
	}

	depCount, isNull, err := c.nextGroup(c.depsGroup)
	if err != nil {
		return RawChangeMeta{}, err
	}
	if isNull {
		depCount = 0
	}

	deps := make([]uint64, 0, min(depCount, 64))
	for i := uint64(0); i < depCount; i++ {
		depIdx, isNull, e := c.nextInt64(c.depsIndex)
		if e != nil {
			return RawChangeMeta{}, fmt.Errorf("dep index: %w", e)
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

func (c *ChangesReader) Fork() (*ChangesReader, error) {
	fork := &ChangesReader{}
	var err error

	if c.actor != nil {
		fork.actor, err = c.actor.Fork()
		if err != nil {
			return nil, err
		}
	}
	if c.seqNum != nil {
		fork.seqNum, err = c.seqNum.Fork()
		if err != nil {
			return nil, err
		}
	}
	if c.maxOp != nil {
		fork.maxOp, err = c.maxOp.Fork()
		if err != nil {
			return nil, err
		}
	}
	if c.time != nil {
		fork.time, err = c.time.Fork()
		if err != nil {
			return nil, err
		}
	}
	if c.message != nil {
		fork.message, err = c.message.Fork()
		if err != nil {
			return nil, err
		}
	}
	if c.depsGroup != nil {
		fork.depsGroup, err = c.depsGroup.Fork()
		if err != nil {
			return nil, err
		}
	}
	if c.depsIndex != nil {
		fork.depsIndex, err = c.depsIndex.Fork()
		if err != nil {
			return nil, err
		}
	}
	return fork, nil
}

// ChangesWriter is a stateful encoder for change metadata columns in document chunks.
type ChangesWriter struct {
	actor     *ActorWriter
	seqNum    *DeltaWriter
	maxOp     *DeltaWriter
	time      *DeltaWriter
	message   *StringWriter
	depsGroup *GroupWriter
	depsIndex *DeltaWriter
}

func NewChangesWriter(actor, seqNum, maxOp, time, message, depsGroup, depsIndex io.Writer) *ChangesWriter {
	return &ChangesWriter{
		actor:     NewActorWriter(actor),
		seqNum:    NewDeltaWriter(seqNum),
		maxOp:     NewDeltaWriter(maxOp),
		time:      NewDeltaWriter(time),
		message:   NewStringWriter(message),
		depsGroup: NewGroupWriter(depsGroup),
		depsIndex: NewDeltaWriter(depsIndex),
	}
}

func (c *ChangesWriter) Append(m RawChangeMeta) {
	c.actor.Append(rle.NewNullableUint64(m.ActorIdx))
	c.seqNum.Append(rle.NewNullableInt64(int64(m.SeqNum)))
	c.maxOp.Append(rle.NewNullableInt64(int64(m.MaxOp)))
	if m.Time != nil {
		c.time.Append(rle.NewNullableInt64(*m.Time))
	} else {
		c.time.Append(rle.NewNullInt64())
	}
	if m.Message != nil {
		c.message.Append(rle.NewNullableString(*m.Message))
	} else {
		c.message.Append(rle.NewNullString())
	}
	c.depsGroup.Append(rle.NewNullableUint64(uint64(len(m.Deps))))
	for _, dep := range m.Deps {
		c.depsIndex.Append(rle.NewNullableInt64(int64(dep)))
	}
}

func (c *ChangesWriter) Flush() error {
	for _, f := range []interface{ Flush() error }{
		c.actor, c.seqNum, c.maxOp, c.time, c.message, c.depsGroup, c.depsIndex,
	} {
		if err := f.Flush(); err != nil {
			return err
		}
	}
	return nil
}
