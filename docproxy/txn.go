package docproxy

import (
	"bytes"
	"fmt"
	"iter"

	"gotomerge/format"
	"gotomerge/opset"
	"gotomerge/types"
	ioutil "gotomerge/utils/io"
)

// Txn is a write handle for a single document transaction. It is obtained from
// [Document.Begin] or passed into the callback of [Document.Change].
//
// All mutations route through the underlying [opset.Transaction]; reads reflect
// the committed state of the document at the time Begin was called (the
// transaction is not yet applied).
//
// Call [Txn.Commit] to apply the transaction and serialise the change chunk, or
// [Txn.Rollback] to discard it. A Txn that is neither committed nor rolled back
// simply holds the document mutex open — always pair Begin with a deferred
// Rollback when using manual transactions:
//
//	txn := doc.Begin()
//	defer txn.Rollback() // no-op after Commit
//	txn.Map("x").Set("y", 42)
//	return txn.Commit()
type Txn struct {
	doc  *Document
	t    *opset.Transaction
	done bool // true after Commit or Rollback
}

// -- Write methods (root level) ----------------------------------------------

// Map returns a [MapView] for the map at key in the root. If key is absent a
// new map is created. Panics with [ErrType] if key already holds a non-map value.
func (txn *Txn) Map(key string) MapView {
	txn.mustOpen()
	obj := txn.getOrMakeMap(types.RootObjectId(), key)
	return MapView{s: txn.doc.s, txn: txn.t, obj: obj}
}

// Text returns a [TextView] for the text object at key in the root. If key is
// absent a new text object is created. Panics with [ErrType] if key already
// holds a non-text value.
func (txn *Txn) Text(key string) TextView {
	txn.mustOpen()
	obj := getOrMakeText(txn.doc.s, txn.t, types.RootObjectId(), key)
	return TextView{s: txn.doc.s, txn: txn.t, obj: obj}
}

// List returns a [ListView] for the list at key in the root. If key is absent a
// new list is created. Panics with [ErrType] if key already holds a non-list value.
func (txn *Txn) List(key string) ListView {
	txn.mustOpen()
	obj := txn.getOrMakeList(types.RootObjectId(), key)
	return ListView{s: txn.doc.s, txn: txn.t, obj: obj}
}

// Set sets key in the root map to a scalar value. value must be one of: bool,
// int64, float64, string, []byte, or nil. Panics with [ErrType] if key currently
// holds a non-scalar (map or list).
func (txn *Txn) Set(key string, value any) {
	txn.mustOpen()
	txn.t.MapSet(types.RootObjectId(), key, value)
}

// Delete removes key from the root map. No-op if key is absent.
func (txn *Txn) Delete(key string) {
	txn.mustOpen()
	txn.t.MapDelete(types.RootObjectId(), key)
}

// -- Read methods (mirrors Document) -----------------------------------------

// Keys returns all keys that have a live value in the root map.
func (txn *Txn) Keys() []string {
	return txn.doc.s.MapKeys(types.RootObjectId())
}

// Values returns an iterator over key/value pairs at the root. Views returned
// here are wired to this transaction and may be written to.
func (txn *Txn) Values() iter.Seq2[string, Value] {
	return mapValues(txn.doc.s, txn.t, types.RootObjectId())
}

// Get returns the value at key in the root map. The returned [Value] reflects
// the committed state; uncommitted operations in this transaction are not visible.
func (txn *Txn) Get(key string) (Value, bool) {
	return mapGet(txn.doc.s, txn.t, types.RootObjectId(), key)
}

// -- Lifecycle ---------------------------------------------------------------

// Commit applies the buffered operations as a single change chunk, serialises
// it to the incremental save buffer, and releases the document mutex.
// Returns an error if the transaction has no operations or serialisation fails.
// After Commit, the Txn must not be used again.
func (txn *Txn) Commit() error {
	if txn.done {
		return fmt.Errorf("Commit called on already-finished transaction")
	}
	txn.done = true
	defer txn.doc.mu.Unlock()

	if !txn.t.HasOps() {
		return nil // empty transaction; nothing to write
	}

	var buf bytes.Buffer
	if err := txn.t.Commit(&buf); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	raw := buf.Bytes()

	// Re-parse to get the ChangeChunk with its hash (set by ReadChunk after
	// checksum verification). This is a cheap round-trip over an in-memory buffer.
	sr := ioutil.NewSubReader(raw)
	chunk, toSkip, err := format.ReadChunk(sr)
	if err != nil {
		return fmt.Errorf("commit: re-parse: %w", err)
	}
	if err := sr.Skip(toSkip); err != nil {
		return fmt.Errorf("commit: re-parse skip: %w", err)
	}
	cc, ok := chunk.(*format.ChangeChunk)
	if !ok {
		return fmt.Errorf("commit: expected ChangeChunk, got %T", chunk)
	}

	txn.doc.seqNum++
	txn.doc.allChanges = append(txn.doc.allChanges, storedChange{cc: cc, raw: raw})
	txn.doc.unsaved = append(txn.doc.unsaved, raw)
	return nil
}

// Rollback discards all buffered operations and releases the document mutex.
// Safe to call after Commit (becomes a no-op). Use with defer to guarantee
// cleanup on error paths.
func (txn *Txn) Rollback() {
	if txn.done {
		return
	}
	txn.done = true
	txn.doc.mu.Unlock()
}

// -- internal ----------------------------------------------------------------

func (txn *Txn) mustOpen() {
	if txn.done {
		panic("operation on finished transaction")
	}
}

func (txn *Txn) getOrMakeMap(obj types.ObjectId, key string) types.ObjectId {
	return getOrMakeMap(txn.doc.s, txn.t, obj, key)
}

func (txn *Txn) getOrMakeList(obj types.ObjectId, key string) types.ObjectId {
	return getOrMakeList(txn.doc.s, txn.t, obj, key)
}
