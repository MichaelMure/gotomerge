package docproxy

var _ parentView = MapView{}

// MapView is a proxy to a key/value map. It can be used to modify that map.
// Keys are always strings, values can be any type supported by automerge (map, list, text, null, bool, ...).
type MapView struct {
	parent parentView
	key    any
}

// Value returns the underlying map.
func (mv MapView) Value() map[string]any {
	val := mv.parent.get(mv.key)
	cast, ok := val.(map[string]any)
	if !ok {
		panic(ErrType{expected: cast, got: val})
	}
	return cast
}

// Bool returns a proxy to a boolean value.
func (mv MapView) Bool(key string) BoolView {
	mv.init(key, false)
	return BoolView{
		parent: mv,
		key:    key,
	}
}

// Map returns a proxy to a map value.
func (mv MapView) Map(key string) MapView {
	mv.init(key, make(map[string]any))
	return MapView{
		parent: mv,
		key:    key,
	}
}

func (mv MapView) notify(key any, subkeys []string, changetype string, value any) {
	mv.parent.notify(mv.key, append(subkeys, key.(string)), changetype, value)
}

func (mv MapView) init(key any, value any) {
	m := mv.Value()
	if v, ok := m[key.(string)]; !ok {
		m[key.(string)] = value
	} else if _, ok := v.(bool); !ok {
		panic(ErrType{expected: true, got: v})
	}
}

func (mv MapView) get(key any) any {
	return mv.Value()[key.(string)]
}

func (mv MapView) set(key any, value any, changetype string) {
	mv.Value()[key.(string)] = value
	mv.parent.notify(mv.key, []string{key.(string)}, changetype, value)
}
