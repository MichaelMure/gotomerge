package docproxy

// BoolView is a proxy to a bool in a Document. It can be used to modify that boolean.
type BoolView struct {
	parent parentView
	key    any
}

// Value returns the underlying bool.
func (b BoolView) Value() bool {
	val := b.parent.get(b.key)
	cast, ok := val.(bool)
	if !ok {
		panic(ErrType{expected: cast, got: val})
	}
	return cast
}

// Toggle toggles the underlying bool.
func (b BoolView) Toggle() {
	b.parent.set(b.key, !b.Value(), "set")
}

// Set sets the underlying bool to the given value.
func (b BoolView) Set(val bool) {
	b.parent.set(b.key, val, "set")
}

type StringView struct {
	parent parentView
	key    any
}

// Value returns the underlying string.
func (sv StringView) Value() string {
	val := sv.parent.get(sv.key)
	cast, ok := val.(string)
	if !ok {
		panic(ErrType{expected: cast, got: val})
	}
	return cast
}

// Set sets the underlying string to the given value.
func (sv StringView) Set(val string) {
	sv.parent.set(sv.key, val, "set")
}
