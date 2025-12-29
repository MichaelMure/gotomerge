package docproxy

import (
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

// Change executes a callback function with a Document proxy to modify the Document.
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

func (doc *Document) Value() map[string]any {
	return doc.data
}
