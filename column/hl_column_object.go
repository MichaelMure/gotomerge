package column

import (
	"io"
	"math"

	"github.com/MichaelMure/gotomerge/column/rle"
	"github.com/MichaelMure/gotomerge/types"
)

// ObjectReader is a stateful reader for object ID columns.
type ObjectReader struct {
	actor   *ActorReader
	counter *UlebReader
}

func NewObjectReader(actor *ActorReader, counter *UlebReader) *ObjectReader {
	return &ObjectReader{actor: actor, counter: counter}
}

func (o *ObjectReader) Next() (types.ObjectId, error) {
	var actorIdx uint64
	var nullActorIdx bool
	var counter uint64
	var nullCounter bool

	if o.actor == nil {
		nullActorIdx = true
	} else {
		nv, err := o.actor.Next()
		if err != nil {
			return types.ObjectId{}, err
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
			return types.ObjectId{}, err
		}
		if v, valid := nv.Value(); valid {
			counter = v
		} else {
			nullCounter = true
		}
	}

	switch {
	case !nullActorIdx && nullCounter:
		return types.ObjectId{}, ErrUnexpectedNull("counter")
	case nullActorIdx && !nullCounter:
		return types.ObjectId{}, ErrUnexpectedNull("actor index")
	case nullActorIdx && nullCounter:
		return types.RootObjectId(), nil
	default:
		if actorIdx >= math.MaxUint32 {
			return types.ObjectId{}, ErrOutOfRange("actor index")
		}
		if counter >= math.MaxUint32 {
			return types.ObjectId{}, ErrOutOfRange("counter")
		}
		return types.ObjectId{ActorIdx: uint32(actorIdx), Counter: uint32(counter)}, nil
	}
}

func (o *ObjectReader) Fork() (*ObjectReader, error) {
	var actor *ActorReader
	var counter *UlebReader
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
	return &ObjectReader{actor: actor, counter: counter}, nil
}

// ObjectWriter is a stateful encoder for object ID columns.
type ObjectWriter struct {
	actor *ActorWriter
	ctr   *UlebWriter
}

func NewObjectWriter(actor, ctr io.Writer) *ObjectWriter {
	return &ObjectWriter{actor: NewActorWriter(actor), ctr: NewUlebWriter(ctr)}
}

func (o *ObjectWriter) Append(obj types.ObjectId, mapper types.ActorMapper) {
	if obj.IsRoot() {
		o.actor.Append(rle.NewNullUint64())
		o.ctr.Append(rle.NewNullUint64())
	} else {
		o.actor.Append(rle.NewNullableUint64(uint64(mapper.Map(obj.ActorIdx))))
		o.ctr.Append(rle.NewNullableUint64(uint64(obj.Counter)))
	}
}

func (o *ObjectWriter) Flush() error {
	if err := o.actor.Flush(); err != nil {
		return err
	}
	return o.ctr.Flush()
}
