package column

import (
	"math"

	"gotomerge/types"
)

// GroupedOpIdReader is a stateful reader for grouped operation ID columns.
type GroupedOpIdReader struct {
	name    string
	group   *GroupReader
	actor   *ActorReader
	counter *DeltaReader
}

func NewGroupedOpIdReader(name string, group *GroupReader, actor *ActorReader, counter *DeltaReader) *GroupedOpIdReader {
	return &GroupedOpIdReader{name: name, group: group, actor: actor, counter: counter}
}

func (g *GroupedOpIdReader) Next() ([]types.OpId, error) {
	var count uint64
	if g.group == nil {
		return nil, ErrDone
	}
	nv, err := g.group.Next()
	if err != nil {
		return nil, err
	}
	v, valid := nv.Value()
	if !valid {
		return nil, ErrUnexpectedNull(g.name + " group")
	}
	count = v

	res := make([]types.OpId, 0, min(count, 100))

	for i := uint64(0); i < count; i++ {
		var actorIdx uint64
		var nullActorIdx bool
		var counter int64
		var nullCounter bool

		if g.actor == nil {
			nullActorIdx = true
		} else {
			anv, e := g.actor.Next()
			if e != nil {
				return nil, e
			}
			if av, aValid := anv.Value(); aValid {
				actorIdx = av
			} else {
				nullActorIdx = true
			}
		}

		if g.counter == nil {
			nullCounter = true
		} else {
			cnv, e := g.counter.Next()
			if e != nil {
				return nil, e
			}
			if cv, cValid := cnv.Value(); cValid {
				counter = cv
			} else {
				nullCounter = true
			}
		}

		switch {
		case nullActorIdx:
			return nil, ErrUnexpectedNull("counter")
		case nullCounter:
			return nil, ErrUnexpectedNull("actor index")
		default:
			if actorIdx >= math.MaxUint32 {
				return nil, ErrOutOfRange("actor index")
			}
			if counter < 0 || counter >= math.MaxInt64 {
				return nil, ErrOutOfRange("counter")
			}
			res = append(res, types.OpId{ActorIdx: uint32(actorIdx), Counter: uint32(counter)})
		}
	}

	return res, nil
}

func (g *GroupedOpIdReader) Fork() (*GroupedOpIdReader, error) {
	var group *GroupReader
	var actor *ActorReader
	var counter *DeltaReader
	var err error

	if g.group != nil {
		group, err = g.group.Fork()
		if err != nil {
			return nil, err
		}
	}
	if g.actor != nil {
		actor, err = g.actor.Fork()
		if err != nil {
			return nil, err
		}
	}
	if g.counter != nil {
		counter, err = g.counter.Fork()
		if err != nil {
			return nil, err
		}
	}
	return &GroupedOpIdReader{name: g.name, group: group, actor: actor, counter: counter}, nil
}
