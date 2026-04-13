// Package docproxy provides the public API for reading and writing Automerge documents.
//
// # Overview
//
// A [Document] is the top-level object. It wraps an [opset.OpSet] (the CRDT state) and an
// actor identity, and exposes four concerns:
//
//   - Construction: [NewDocument] creates an empty document; [LoadDocument] and
//     [NewDocumentFromJSON] deserialise existing data.
//   - Reading: [Document.Get], [Document.Map], [Document.List], [Document.Values],
//     [Document.Keys] traverse the current committed state. Text objects are accessed
//     via [Document.Get] with a type assertion to [TextView].
//   - Writing: [Document.Change] runs a callback with a [Txn] write handle; [Document.Begin]
//     gives a manual [Txn] for explicit Commit / Rollback control. Every non-empty transaction
//     produces a causally-linked, hashable change record applied to the document.
//     [Document.Merge] applies changes from another document.
//   - Persistence: [Document.Save] / [Document.SaveIncremental] serialise the document to
//     bytes. [Document.MarshalJSON] renders the current state as JSON.
//
// # Transactions
//
// Every [Document.Change] (or [Txn.Commit]) call that contains at least one operation
// produces a change record: a named, causally-linked, hashable unit of history. Change
// records carry authorship (actor + sequence number), causal dependencies (the current
// heads at Begin time), and a SHA-256 hash used by peers for synchronisation and conflict
// detection.
//
// # Persistence model
//
// [Document.Save] writes a compact document chunk representing the full current state.
// After a successful Save the incremental buffer is cleared, because the document chunk
// already subsumes all prior changes.
//
// [Document.SaveIncremental] writes only the change chunks produced since the last Save or
// SaveIncremental call. It is cheap to call after every Change.
//
// Typical usage — new document:
//
//	doc := docproxy.NewDocument()
//	doc.Change(func(txn *docproxy.Txn) error {
//	    txn.Map("config").Set("debug", true)
//	    return nil
//	})
//	doc.Save(w)
//
// Typical usage — load and modify:
//
//	doc, _ := docproxy.LoadDocument(f)
//	doc.Change(func(txn *docproxy.Txn) error {
//	    txn.Map("config").Set("debug", true)
//	    return nil
//	})
//	doc.SaveIncremental(w) // append only the new change chunk
//
// TODO: generate a diagram of the major code path: views, opset, index, ...
package docproxy

import (
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"sync"

	"github.com/MichaelMure/gotomerge/format"
	"github.com/MichaelMure/gotomerge/opset"
	"github.com/MichaelMure/gotomerge/types"
	ioutil "github.com/MichaelMure/gotomerge/utils/io"
)

// storedChange holds the hash and raw serialised bytes of an applied change.
// The hash (set by ReadChunk after checksum verification) is the change's
// globally unique identity used for deduplication in Merge. The raw bytes are
// kept because the change hash is defined over those exact bytes — re-encoding
// would produce a different hash and break the dependency graph.
// The decoded *ChangeChunk is intentionally not retained: keeping 259k decoded
// chunks alive simultaneously causes severe GC pressure for large documents.
// Callers that need the decoded form (Save, Merge) re-parse from raw on demand,
// keeping only one *ChangeChunk live at a time.
type storedChange struct {
	hash types.ChangeHash
	raw  []byte
}

// Document is an Automerge document. It is safe for concurrent reads, but
// [Document.Change], [Document.Begin], [Document.Save], [Document.SaveIncremental],
// and [Document.Merge] serialise through an internal mutex.
type Document struct {
	mu     sync.Mutex
	s      *opset.OpSet
	actor  types.ActorId
	seqNum uint64 // next sequence number for this actor; incremented on each Commit

	// allChanges holds every ChangeChunk that has been applied, in dependency
	// order. It is the source of truth for Merge and Save.
	allChanges []storedChange

	// unsaved holds raw change chunk bytes produced since the last Save.
	// SaveIncremental drains this slice; Save clears it after writing a full
	// document chunk (which subsumes all prior changes).
	unsaved [][]byte
}

// NewDocument creates an empty document with a fresh random actor identity.
func NewDocument() *Document {
	return &Document{
		s:      opset.New(),
		actor:  types.NewActorId(),
		seqNum: 1,
	}
}

// LoadDocument reads all chunks from r and returns a Document backed by the
// resulting OpSet. Both document chunks and change chunks are accepted and
// applied in order. The returned document is assigned a fresh random actor
// identity so that subsequent Change calls do not collide with the loaded history.
func LoadDocument(r io.Reader) (*Document, error) {
	// Read everything into memory so we can slice raw bytes per chunk by position.
	all, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("load document: read: %w", err)
	}

	s := opset.New()
	var allChanges []storedChange

	sr := ioutil.NewSubReader(all)
	for {
		start := sr.Consumed()
		chunk, toSkip, err := format.ReadChunk(sr)
		if err == io.EOF {
			break
		}
		if skipErr := sr.Skip(toSkip); skipErr != nil {
			return nil, fmt.Errorf("load document: skip: %w", skipErr)
		}
		if err != nil {
			return nil, fmt.Errorf("load document: %w", err)
		}
		raw := all[start:sr.Consumed()]

		switch c := chunk.(type) {
		case *format.DocumentChunk:
			if err := s.ApplyDocument(c); err != nil {
				return nil, fmt.Errorf("load document: apply document chunk: %w", err)
			}
		case *format.ChangeChunk:
			if err := s.ApplyChange(c); err != nil {
				return nil, fmt.Errorf("load document: apply change chunk: %w", err)
			}
			allChanges = append(allChanges, storedChange{hash: c.Hash, raw: raw})
		}
		if sr.Empty() {
			break
		}
	}
	return &Document{
		s:          s,
		actor:      types.NewActorId(),
		seqNum:     1,
		allChanges: allChanges,
	}, nil
}

// NewDocumentFromJSON creates a new empty document and applies a single Change
// that sets the root map from the JSON object in data. data must be a JSON object
// (not an array or scalar). The resulting document has a single change chunk in its
// history representing the initial state.
func NewDocumentFromJSON(data []byte) (*Document, error) {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("NewDocumentFromJSON: %w", err)
	}
	doc := NewDocument()
	err := doc.Change(func(txn *Txn) error {
		setMapFromJSON(txn.t, types.RootObjectId(), m)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// -- Read methods -------------------------------------------------------------

// Keys returns all keys that have a live value in the root map.
func (doc *Document) Keys() []string {
	return doc.s.MapKeys(types.RootObjectId())
}

// Values returns an iterator over key/value pairs at the root. Each value is
// a [MapView], [ListView], [BoolView], [StringView], or other concrete [Value]
// type matching the actual data. The iteration order is insertion order within
// each source (snapshot then delta).
func (doc *Document) Values() iter.Seq2[string, Value] {
	return mapValues(doc.s, nil, types.RootObjectId())
}

// Get returns the value at key in the root map. The returned [Value] is a
// concrete view type that can be type-asserted: [MapView], [ListView],
// [BoolView], [StringView], etc.
func (doc *Document) Get(key string) (Value, bool) {
	return mapGet(doc.s, nil, types.RootObjectId(), key)
}

// Map returns a read-only [MapView] for the map object at key.
// Returns false if the key is absent or holds a non-map value.
func (doc *Document) Map(key string) (MapView, bool) {
	v, ok := doc.Get(key)
	if !ok {
		return MapView{}, false
	}
	mv, ok := v.(MapView)
	return mv, ok
}

// List returns a read-only [ListView] for the list object at key.
// Returns false if the key is absent or holds a non-list value.
func (doc *Document) List(key string) (ListView, bool) {
	v, ok := doc.Get(key)
	if !ok {
		return ListView{}, false
	}
	lv, ok := v.(ListView)
	return lv, ok
}

// Text returns a read-only [TextView] for the text object at key.
// Returns false if the key is absent or holds a non-text value.
func (doc *Document) Text(key string) (TextView, bool) {
	v, ok := doc.Get(key)
	if !ok {
		return TextView{}, false
	}
	tv, ok := v.(TextView)
	return tv, ok
}

// -- Write methods ------------------------------------------------------------

// Begin starts a new manual transaction. The caller must call [Txn.Commit] or
// [Txn.Rollback] when done. If Commit is not called the transaction is silently
// discarded (no panic, no change applied).
//
// Only one transaction may be active per document at a time; Begin acquires the
// document mutex and holds it until Commit or Rollback.
func (doc *Document) Begin() *Txn {
	doc.mu.Lock()
	t := doc.s.Begin(doc.actor, doc.seqNum)
	return &Txn{doc: doc, t: t}
}

// Change executes fn with a [Txn] write handle. If fn returns nil the
// transaction is committed; otherwise it is rolled back and the error is
// returned. Panics inside fn that carry an [ErrType] are caught and returned
// as errors; all other panics propagate normally.
func (doc *Document) Change(fn func(*Txn) error) (err error) {
	txn := doc.Begin()

	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(ErrType); ok {
				txn.Rollback()
				err = e
			} else {
				txn.Rollback()
				panic(r)
			}
		}
	}()

	if err = fn(txn); err != nil {
		txn.Rollback()
		return err
	}
	return txn.Commit()
}

// -- Persistence -------------------------------------------------------------

// Save writes a full document chunk to w. The document chunk is a compacted
// representation of the entire current OpSet state. After a successful Save,
// the incremental buffer is cleared because the document chunk subsumes all
// prior change chunks.
//
// Use [Document.SaveIncremental] for cheap incremental persistence between full saves.
func (doc *Document) Save(w io.Writer) error {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	if err := doc.s.ExportDocument(w); err != nil {
		return fmt.Errorf("Save: %w", err)
	}
	doc.unsaved = doc.unsaved[:0]
	return nil
}

// SaveIncremental writes only the change chunks produced since the last Save or
// SaveIncremental call to w. If no new changes exist, w is not written to.
//
// The output is a valid Automerge byte stream that can be appended to a file or
// fed to [LoadDocument] together with a preceding full save.
func (doc *Document) SaveIncremental(w io.Writer) error {
	doc.mu.Lock()
	defer doc.mu.Unlock()

	for _, chunk := range doc.unsaved {
		if _, err := w.Write(chunk); err != nil {
			return fmt.Errorf("SaveIncremental: %w", err)
		}
	}
	doc.unsaved = doc.unsaved[:0]
	return nil
}

// Merge applies all changes from other that are not already present in doc.
// After Merge, doc contains the union of both documents' histories. The merged
// changes are added to the incremental save buffer so a subsequent
// [Document.SaveIncremental] will include them.
//
// Both documents must share a common origin (same initial state); merging
// completely unrelated documents produces undefined results.
func (doc *Document) Merge(other *Document) error {
	other.mu.Lock()
	changes := make([]storedChange, len(other.allChanges))
	copy(changes, other.allChanges)
	other.mu.Unlock()

	doc.mu.Lock()
	defer doc.mu.Unlock()

	for _, sc := range changes {
		if _, already := doc.s.AppliedHashes()[sc.hash]; already {
			continue
		}
		sr := ioutil.NewSubReader(sc.raw)
		chunk, _, err := format.ReadChunk(sr)
		if err != nil {
			return fmt.Errorf("merge: re-parse change: %w", err)
		}
		cc, ok := chunk.(*format.ChangeChunk)
		if !ok {
			return fmt.Errorf("merge: expected ChangeChunk, got %T", chunk)
		}
		if err := doc.s.ApplyChange(cc); err != nil {
			return fmt.Errorf("merge: %w", err)
		}
		doc.allChanges = append(doc.allChanges, sc)
		doc.unsaved = append(doc.unsaved, sc.raw)
	}
	return nil
}

// -- JSON --------------------------------------------------------------------

// MarshalJSON serialises the current document state as a JSON object.
// Conflicts (multiple live values at a key) are resolved by picking the
// CRDT winner (highest OpId). Nested maps and lists are serialised
// recursively. Text objects are serialised as plain strings.
func (doc *Document) MarshalJSON() ([]byte, error) {
	v := materializeObj(doc.s, types.RootObjectId(), types.ActionMakeMap)
	return json.Marshal(v)
}
