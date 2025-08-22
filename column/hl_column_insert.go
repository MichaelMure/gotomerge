package column

import (
	"iter"
)

type InsertColumnIter struct {
	next func() (bool, error, bool)
	stop func()
}

func InsertColumn(col BooleanColumnIter) InsertColumnIter {
	var res InsertColumnIter
	if col != nil {
		res.next, res.stop = iter.Pull2(col)
	}
	return res
}

func (a InsertColumnIter) Next() (bool, error) {
	if a.next == nil {
		return false, nil
	}
	val, err, ok := a.next()
	if !ok {
		return false, ErrDone
	}
	if err != nil {
		return false, err
	}
	return val, nil
}

func (a InsertColumnIter) Stop() {
	if a.stop != nil {
		a.stop()
	}
}
