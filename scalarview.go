package gotomerge

type BoolView struct {
	doc    *Document
	parent parentView
	key    any
	data   bool // no pointer, lightweight
}

func (b BoolView) Value() bool {
	return b.data
}

func (b BoolView) Toggle() {
	b.data = !(b.data)
	b.parent.set(b.key, nil, "set", b.data)
}

func (b BoolView) Set(val bool) {
	b.data = val
	b.parent.set(b.key, nil, "set", b.data)
}

type StringView struct {
	doc    *Document
	parent parentView
	key    any
	data   string
}
