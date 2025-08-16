package docproxy

var _ parentView = MapView{}

type MapView struct {
	doc    *Document
	parent parentView
	key    any
	data   map[string]any // no pointer, only copy the map header
}

func (mv MapView) set(key any, subkeys []string, changetype string, value any) {
	mv.data[key.(string)] = value
	mv.parent.set(mv.key, append(subkeys, key.(string)), changetype, value)
}

func (mv MapView) Value() map[string]any {
	return mv.data
}

func (mv MapView) Bool(key string) BoolView {
	if v, ok := mv.doc.data[key]; !ok {
		mv.doc.data[key] = false
	} else if _, ok := v.(bool); !ok {
		panic(ErrType{expected: true, got: v})
	}
	return BoolView{
		doc:    mv.doc,
		parent: mv,
		key:    key,
		data:   mv.doc.data[key].(bool),
	}
}

func (mv MapView) Map(key string) MapView {
	if v, ok := mv.doc.data[key]; !ok {
		mv.doc.data[key] = make(map[string]any)
	} else if _, ok := v.(map[string]any); !ok {
		panic(ErrType{expected: map[string]any{}, got: v})
	}
	return MapView{
		doc:    mv.doc,
		parent: mv,
		key:    key,
		data:   mv.doc.data[key].(map[string]any),
	}
}
