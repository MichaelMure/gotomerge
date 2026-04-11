// Package gotomerge is a Go implementation of the Automerge CRDT document format.
//
// Create a document, write to it inside a [Change] callback, and read values
// back with [As]:
//
//	doc := gotomerge.NewDocument()
//	doc.Change(func(txn *gotomerge.Txn) error {
//	    txn.Set("title", "Hello")
//	    txn.Map("meta").Set("author", "alice")
//	    txn.List("tags").Append("go")
//	    return nil
//	})
//
//	title, _ := gotomerge.As[string](doc.Get("title"))
//
//	type Meta struct {
//	    Author string `automerge:"author"`
//	}
//	meta, _ := gotomerge.As[Meta](doc.Get("meta"))
//
// # Persistence
//
// [Document.Save] writes a compact full-state snapshot. [Document.SaveIncremental]
// writes only the changes since the last save — cheap to call after every [Change].
// Both outputs are readable by [LoadDocument].
//
// # Merging
//
// Fork a document by saving and reloading, have each peer make independent changes,
// then call [Document.Merge] to converge them. Automerge resolves all conflicts
// automatically using CRDT semantics.
package gotomerge

import (
	"io"

	"github.com/MichaelMure/gotomerge/docproxy"
)

// -- Core types --------------------------------------------------------------

// Document is an Automerge document. Safe for concurrent reads; writes
// serialise through an internal mutex.
type Document = docproxy.Document

// Txn is a write handle for a single transaction. Obtained from [Document.Begin]
// or passed into [Document.Change].
type Txn = docproxy.Txn

// Value is a read-only or read-write view of a single value inside a document.
// Use [As] to convert it to a concrete Go type without type assertions.
type Value = docproxy.Value

// ErrType is returned by [Document.Change] when a key is accessed as the wrong
// type (e.g. treating a list as a map).
type ErrType = docproxy.ErrType

// -- Container view types ----------------------------------------------------
//
// These are the return types of the named collection accessors on Document and
// Txn. Their read methods ([MapView.Get], [ListView.Get], etc.) return [Value],
// which is then converted with [As].

// MapView is a view of a map object. Returned by [Document.Map] and [Txn.Map].
type MapView = docproxy.MapView

// ListView is a view of a list object. Returned by [Document.List] and [Txn.List].
type ListView = docproxy.ListView

// TextView is a view of a collaborative text object. Returned by [Document.Text]
// and [Txn.Text].
type TextView = docproxy.TextView

// -- Constructors ------------------------------------------------------------

// NewDocument creates an empty document with a fresh random actor identity.
func NewDocument() *Document { return docproxy.NewDocument() }

// LoadDocument reads all chunks from r and returns a Document. Both full
// document snapshots and incremental change streams are accepted.
func LoadDocument(r io.Reader) (*Document, error) { return docproxy.LoadDocument(r) }

// NewDocumentFromJSON creates a new document whose root map is initialised from
// the JSON object in data. data must be a JSON object (not an array or scalar).
func NewDocumentFromJSON(data []byte) (*Document, error) {
	return docproxy.NewDocumentFromJSON(data)
}

// -- Reading -----------------------------------------------------------------

// As converts a (Value, bool) pair — as returned by [Document.Get] and similar
// methods — into a typed Go value T without explicit type assertions.
//
// Returns (zero, false) if ok is false (key absent) or if the conversion is not
// supported. Supported target types include all scalars (string, bool, int64,
// float64, []byte, time.Time), structs with automerge tags, map[string]T, []T,
// and the container view types [MapView], [ListView], [TextView].
//
// Example:
//
//	name, ok := gotomerge.As[string](doc.Get("name"))
//
//	type Config struct {
//	    Debug bool  `automerge:"debug"`
//	    Port  int64 `automerge:"port"`
//	}
//	cfg, ok := gotomerge.As[Config](doc.Get("config"))
func As[T any](v Value, ok bool) (T, bool) { return docproxy.As[T](v, ok) }
