package column

import (
	"io"
	"math"

	"github.com/MichaelMure/gotomerge/column/rle"
	"github.com/MichaelMure/gotomerge/types"
)

// OpIdReader is a stateful reader for operation ID columns.
type OpIdReader struct {
	actor   *ActorReader
	counter *DeltaReader
}

func NewOpIdReader(actor *ActorReader, counter *DeltaReader) *OpIdReader {
	return &OpIdReader{actor: actor, counter: counter}
}

func (o *OpIdReader) Next() (types.OpId, error) {
	var actorIdx uint64
	var nullActorIdx bool
	var counter int64
	var nullCounter bool

	if o.actor == nil {
		nullActorIdx = true
	} else {
		nv, err := o.actor.Next()
		if err != nil {
			return types.OpId{}, err
		}
		if v, valid := nv.Value(); valid {
			actorIdx = v
		} else {
			nullActorIdx = true
		}
	}

	if o.counter == nil {
		nullCounter = true
	} else {
		nv, err := o.counter.Next()
		if err != nil {
			return types.OpId{}, err
		}
		if v, valid := nv.Value(); valid {
			counter = v
		} else {
			nullCounter = true
		}
	}

	switch {
	case nullActorIdx:
		return types.OpId{}, ErrUnexpectedNull("counter")
	case nullCounter:
		return types.OpId{}, ErrUnexpectedNull("actor index")
	default:
		if actorIdx >= math.MaxUint32 {
			return types.OpId{}, ErrOutOfRange("actor index")
		}
		if counter < 0 || counter >= math.MaxInt64 {
			return types.OpId{}, ErrOutOfRange("counter")
		}
		return types.OpId{ActorIdx: uint32(actorIdx), Counter: uint32(counter)}, nil
	}
}

func (o *OpIdReader) Fork() (*OpIdReader, error) {
	var actor *ActorReader
	var counter *DeltaReader
	var err error

	if o.actor != nil {
		actor, err = o.actor.Fork()
		if err != nil {
			return nil, err
		}
	}
	if o.counter != nil {
		counter, err = o.counter.Fork()
		if err != nil {
			return nil, err
		}
	}
	return &OpIdReader{actor: actor, counter: counter}, nil
}

// OpIdWriter is a stateful encoder for operation ID columns.
type OpIdWriter struct {
	actor *ActorWriter
	ctr   *DeltaWriter
}

func NewOpIdWriter(actor, ctr io.Writer) *OpIdWriter {
	return &OpIdWriter{actor: NewActorWriter(actor), ctr: NewDeltaWriter(ctr)}
}

func (o *OpIdWriter) Append(id types.OpId, mapper types.ActorMapper) {
	o.actor.Append(rle.NewNullableUint64(uint64(mapper.Map(id.ActorIdx))))
	o.ctr.Append(rle.NewNullableInt64(int64(id.Counter)))
}

func (o *OpIdWriter) Flush() error {
	if err := o.actor.Flush(); err != nil {
		return err
	}
	return o.ctr.Flush()
}
