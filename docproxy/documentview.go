package docproxy

import "strings"

var _ parentView = DocumentView{}

// DocumentView is a proxy to a Document. It can be used to modify the Document.
// It only maps to the root level of the Document.
type DocumentView struct {
	// TODO: add mutex
	doc *Document
}

func (dv DocumentView) Bool(key string) BoolView {
	dv.init(key, false)
	return BoolView{
		parent: dv,
		key:    key,
	}
}

func (dv DocumentView) Map(key string) MapView {
	dv.init(key, make(map[string]any))
	return MapView{
		parent: dv,
		key:    key,
	}
}

func (dv DocumentView) notify(key any, subkeys []string, changetype string, value any) {
	var fullKey strings.Builder
	fullKey.WriteString(".")
	fullKey.WriteString(key.(string))
	for i := len(subkeys) - 1; i >= 0; i-- {
		fullKey.WriteString(".")
		fullKey.WriteString(subkeys[i])
	}
	dv.doc.changes = append(dv.doc.changes, Change{
		Type:  changetype,
		Path:  fullKey.String(),
		Value: value,
	})
}

func (dv DocumentView) init(key any, value any) {
	if v, ok := dv.doc.data[key.(string)]; !ok {
		dv.doc.data[key.(string)] = value
	} else if _, ok := v.(bool); !ok {
		panic(ErrType{expected: true, got: v})
	}
}

func (dv DocumentView) get(key any) any {
	return dv.doc.data[key.(string)]
}

func (dv DocumentView) set(key any, value any, changetype string) {
	dv.doc.data[key.(string)] = value
	dv.notify(key, []string{}, changetype, value)
}
