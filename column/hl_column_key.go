package column

import (
	"fmt"
	"io"
	"math"

	"github.com/MichaelMure/gotomerge/column/rle"
	"github.com/MichaelMure/gotomerge/types"
)

// KeyReader is a stateful reader for key columns.
type KeyReader struct {
	actor   *ActorReader
	counter *DeltaReader
	str     *StringReader
}

func NewKeyReader(actor *ActorReader, counter *DeltaReader, str *StringReader) *KeyReader {
	return &KeyReader{actor: actor, counter: counter, str: str}
}

func (k *KeyReader) Next() (types.Key, error) {
	var actorIdx uint64
	var nullActorIdx bool
	var counter int64
	var nullCounter bool
	var str string
	var nullString bool

	if k.actor == nil {
		nullActorIdx = true
	} else {
		nv, e := k.actor.Next()
		if e != nil {
			return nil, e
		}
		if v, valid := nv.Value(); valid {
			actorIdx = v
		} else {
			nullActorIdx = true
		}
	}

	if k.counter == nil {
		nullCounter = true
	} else {
		nv, e := k.counter.Next()
		if e != nil {
			return nil, e
		}
		if v, valid := nv.Value(); valid {
			counter = v
		} else {
			nullCounter = true
		}
	}

	if k.str == nil {
		nullString = true
	} else {
		nv, e := k.str.Next()
		if e != nil {
			return nil, e
		}
		if v, valid := nv.Value(); valid {
			str = v
		} else {
			nullString = true
		}
	}

	if (!nullActorIdx || !nullCounter) && !nullString {
		return nil, fmt.Errorf("too many values for key")
	}
	if nullActorIdx && nullCounter && nullString {
		return types.NullKey{}, nil
	}

	switch {
	case !nullString:
		return types.KeyString(str), nil
	case nullCounter:
		return nil, ErrUnexpectedNull("counter")
	case nullActorIdx && counter == 0:
		return types.KeyOpId{ActorIdx: 0, Counter: 0}, nil
	case nullActorIdx:
		return nil, ErrUnexpectedNull("actor index")
	default:
		if actorIdx >= math.MaxUint32 {
			return nil, ErrOutOfRange("actor index")
		}
		if counter < 0 || counter >= math.MaxInt64 {
			return nil, ErrOutOfRange("counter")
		}
		return types.KeyOpId{ActorIdx: uint32(actorIdx), Counter: uint32(counter)}, nil
	}
}

func (k *KeyReader) Fork() (*KeyReader, error) {
	var actor *ActorReader
	var counter *DeltaReader
	var str *StringReader
	var err error

	if k.actor != nil {
		actor, err = k.actor.Fork()
		if err != nil {
			return nil, err
		}
	}
	if k.counter != nil {
		counter, err = k.counter.Fork()
		if err != nil {
			return nil, err
		}
	}
	if k.str != nil {
		str, err = k.str.Fork()
		if err != nil {
			return nil, err
		}
	}
	return &KeyReader{actor: actor, counter: counter, str: str}, nil
}

// KeyWriter is a stateful encoder for key columns.
type KeyWriter struct {
	actor *ActorWriter
	ctr   *DeltaWriter
	str   *StringWriter
}

func NewKeyWriter(actor, ctr, str io.Writer) *KeyWriter {
	return &KeyWriter{
		actor: NewActorWriter(actor),
		ctr:   NewDeltaWriter(ctr),
		str:   NewStringWriter(str),
	}
}

func (k *KeyWriter) Append(key types.Key, mapper types.ActorMapper) {
	switch v := key.(type) {
	case types.KeyString:
		k.actor.Append(rle.NewNullUint64())
		k.ctr.Append(rle.NewNullInt64())
		k.str.Append(rle.NewNullableString(string(v)))
	case types.KeyOpId:
		k.str.Append(rle.NewNullString())
		if v.Counter == 0 {
			// head sentinel: null actor, counter = 0 (non-null)
			k.actor.Append(rle.NewNullUint64())
			k.ctr.Append(rle.NewNullableInt64(0))
		} else {
			k.actor.Append(rle.NewNullableUint64(uint64(mapper.Map(v.ActorIdx))))
			k.ctr.Append(rle.NewNullableInt64(int64(v.Counter)))
		}
	default:
		k.actor.Append(rle.NewNullUint64())
		k.ctr.Append(rle.NewNullInt64())
		k.str.Append(rle.NewNullString())
	}
}

func (k *KeyWriter) Flush() error {
	if err := k.actor.Flush(); err != nil {
		return err
	}
	if err := k.ctr.Flush(); err != nil {
		return err
	}
	return k.str.Flush()
}
