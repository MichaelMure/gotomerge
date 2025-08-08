package gotomerge

import (
	"fmt"
	"strings"
	"sync"
)

// stopgap
type Change struct {
	Type  string `json:"type"`
	Path  string `json:"path"`
	Value any    `json:"value"`
}

type Document struct {
	mu      sync.Mutex
	data    map[string]any
	changes []Change
}

func NewDocument() *Document {
	return &Document{
		data: make(map[string]any),
	}
}

type parentView interface {
	set(key any, subkeys []string, changetype string, value any)
}

type ErrType struct {
	expected any
	got      any
}

func (e ErrType) Error() string {
	return fmt.Sprintf("value already assigned to %T, not %T", e.got, e.expected)
}

// Change executes a callback function with a Document proxy to modify the Document
func (doc *Document) Change(fn func(doc *DocumentView) error) (err error) {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	defer func() {
		if r := recover(); r != nil {
			switch r := r.(type) {
			case ErrType:
				err = r
				return
			default:
				panic(r)
			}
		}
	}()

	err = fn(&DocumentView{doc: doc})

	// TODO: auto-commit

	return err
}

var _ parentView = DocumentView{}

type DocumentView struct {
	MapView
	doc *Document
}

func (dv DocumentView) set(key any, subkeys []string, changetype string, value any) {
	dv.doc.data[key.(string)] = value
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

func (dv DocumentView) Bool(key string) BoolView {
	if v, ok := dv.doc.data[key]; !ok {
		dv.doc.data[key] = false
	} else if _, ok := v.(bool); !ok {
		panic(ErrType{expected: true, got: v})
	}
	return BoolView{
		doc:    dv.doc,
		parent: dv,
		key:    key,
		data:   dv.doc.data[key].(bool),
	}
}

func (dv DocumentView) Map(key string) MapView {
	if v, ok := dv.doc.data[key]; !ok {
		dv.doc.data[key] = make(map[string]any)
	} else if _, ok := v.(map[string]any); !ok {
		panic(ErrType{expected: map[string]any{}, got: v})
	}
	return MapView{
		doc:    dv.doc,
		parent: dv,
		key:    key,
		data:   dv.doc.data[key].(map[string]any),
	}
}
