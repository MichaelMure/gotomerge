package gotomerge_test

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"github.com/MichaelMure/gotomerge"
	"github.com/MichaelMure/gotomerge/types"
)

// Example_quickstart shows the minimal end-to-end lifecycle: create a
// document, write some data, read it back, save and reload.
func Example_quickstart() {
	// Create an empty document.
	doc := gotomerge.NewDocument()

	// All writes happen inside a Change callback. The callback receives a Txn
	// and must return nil to commit or an error to roll back.
	doc.Change(func(txn *gotomerge.Txn) error {
		txn.Set("title", "My note")
		txn.Set("draft", true)
		txn.Set("views", int64(0))
		return nil
	})

	// Read scalars with As[T] — no type assertions, no intermediate types.
	title, _ := gotomerge.As[string](doc.Get("title"))
	draft, _ := gotomerge.As[bool](doc.Get("draft"))
	views, _ := gotomerge.As[int64](doc.Get("views"))
	fmt.Println(title, draft, views)

	// Save to bytes and reload. The full document state is preserved.
	var buf bytes.Buffer
	doc.Save(&buf)

	doc2, _ := gotomerge.LoadDocument(&buf)
	title2, _ := gotomerge.As[string](doc2.Get("title"))
	fmt.Println(title2)

	// Output:
	// My note true 0
	// My note
}

// Example_nestedData shows how to work with nested maps and read them
// as plain Go structs using struct tags.
func Example_nestedData() {
	doc := gotomerge.NewDocument()

	doc.Change(func(txn *gotomerge.Txn) error {
		// Obtain a MapView for the nested map and write to it.
		// Calling txn.Map("address") creates the map if it does not exist.
		addr := txn.Map("address")
		addr.Set("city", "London")
		addr.Set("postcode", "EC1A 1BB")

		// Nest as deep as needed.
		txn.Map("settings").Map("display").Set("theme", "dark")
		return nil
	})

	// Read nested fields one at a time...
	if addr, ok := doc.Map("address"); ok {
		city, _ := gotomerge.As[string](addr.Get("city"))
		fmt.Println(city)
	}

	// ...or unmarshal a whole map into a struct in one call.
	type Address struct {
		City     string `automerge:"city"`
		Postcode string `automerge:"postcode"`
	}
	if addr, ok := gotomerge.As[Address](doc.Get("address")); ok {
		fmt.Println(addr.City, addr.Postcode)
	}

	if settings, ok := doc.Map("settings"); ok {
		theme, _ := gotomerge.As[string](settings.Map("display").Get("theme"))
		fmt.Println(theme)
	}

	// Output:
	// London
	// London EC1A 1BB
	// dark
}

// Example_lists shows how to build and read lists.
func Example_lists() {
	doc := gotomerge.NewDocument()

	doc.Change(func(txn *gotomerge.Txn) error {
		l := txn.List("tags")
		l.Append("go")
		l.Append("crdt")
		l.Append("automerge")
		return nil
	})

	// Materialise the whole list into a Go slice in one call.
	tags, _ := gotomerge.As[[]string](doc.Get("tags"))
	fmt.Println(tags)

	// Or access individual elements via ListView.
	if lv, ok := doc.List("tags"); ok {
		first, _ := gotomerge.As[string](lv.Get(0))
		fmt.Println("first:", first)
		fmt.Println("len:", lv.Len())
	}

	// Delete an element and confirm the list shrinks.
	doc.Change(func(txn *gotomerge.Txn) error {
		txn.List("tags").Delete(1) // remove "crdt"
		return nil
	})

	remaining, _ := gotomerge.As[[]string](doc.Get("tags"))
	fmt.Println(remaining)

	// Output:
	// [go crdt automerge]
	// first: go
	// len: 3
	// [go automerge]
}

// Example_text shows collaborative text editing: Splice for precise
// insert/delete at a position, Update for whole-string replacement.
func Example_text() {
	doc := gotomerge.NewDocument()

	doc.Change(func(txn *gotomerge.Txn) error {
		txn.Text("body").Splice(0, 0, "Hello, world!")
		return nil
	})

	// Replace "world" (positions 7–11) with "Go".
	doc.Change(func(txn *gotomerge.Txn) error {
		txn.Text("body").Splice(7, 5, "Go")
		return nil
	})

	body, _ := gotomerge.As[string](doc.Get("body"))
	fmt.Println(body)

	// Update computes the minimal diff and applies it as a single splice.
	// Use it when you have a new desired string but no cursor positions.
	doc.Change(func(txn *gotomerge.Txn) error {
		txn.Text("body").Update("Hi, Go!")
		return nil
	})

	body, _ = gotomerge.As[string](doc.Get("body"))
	fmt.Println(body)

	// TextView gives direct access to length (in runes) and the string value.
	if tv, ok := doc.Text("body"); ok {
		fmt.Println(tv.Len(), tv.Value())
	}

	// Output:
	// Hello, Go!
	// Hi, Go!
	// 7 Hi, Go!
}

// Example_counters shows CRDT counters: unlike plain integers, concurrent
// increments from different peers add together instead of overwriting each other.
func Example_counters() {
	doc := gotomerge.NewDocument()

	// Initialise a counter with types.Counter (distinguishes it from a plain int).
	doc.Change(func(txn *gotomerge.Txn) error {
		txn.Set("likes", types.Counter(0))
		return nil
	})

	// Increment and decrement via Txn.Increment — no type assertions needed.
	doc.Change(func(txn *gotomerge.Txn) error {
		txn.Increment("likes", 3)
		return nil
	})
	doc.Change(func(txn *gotomerge.Txn) error {
		txn.Increment("likes", -1)
		return nil
	})

	likes, _ := gotomerge.As[int64](doc.Get("likes"))
	fmt.Println(likes)

	// Output:
	// 2
}

// Example_merge shows the core Automerge value proposition: two peers start
// from the same snapshot, each makes independent changes, then converge by
// merging. All changes are preserved; no data is lost.
func Example_merge() {
	// Shared starting point.
	base := gotomerge.NewDocument()
	base.Change(func(txn *gotomerge.Txn) error {
		txn.Set("status", "draft")
		txn.Set("likes", types.Counter(0))
		return nil
	})

	// Fork: save to bytes and load two independent copies.
	var snap bytes.Buffer
	base.SaveIncremental(&snap)
	b := snap.Bytes()

	peer1, _ := gotomerge.LoadDocument(bytes.NewReader(b))
	peer2, _ := gotomerge.LoadDocument(bytes.NewReader(b))

	// Each peer makes independent edits.
	peer1.Change(func(txn *gotomerge.Txn) error {
		txn.Set("status", "published") // peer1 publishes
		txn.Increment("likes", 10)
		return nil
	})
	peer2.Change(func(txn *gotomerge.Txn) error {
		txn.Set("author", "alice") // peer2 adds authorship
		txn.Increment("likes", 5)
		return nil
	})

	// Merge peer2 into peer1. Both peers' changes are now visible on peer1.
	peer1.Merge(peer2)

	keys := peer1.Keys()
	sort.Strings(keys)
	fmt.Println(keys)

	status, _ := gotomerge.As[string](peer1.Get("status"))
	fmt.Println(status)

	author, _ := gotomerge.As[string](peer1.Get("author"))
	fmt.Println(author)

	// Both increments add together: 10 + 5 = 15.
	likes, _ := gotomerge.As[int64](peer1.Get("likes"))
	fmt.Println(likes)

	// Output:
	// [author likes status]
	// published
	// alice
	// 15
}

// Example_timestamps shows how to store and retrieve time values.
func Example_timestamps() {
	t0 := time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC)
	doc := gotomerge.NewDocument()

	doc.Change(func(txn *gotomerge.Txn) error {
		txn.Set("created_at", types.FromTime(t0))
		return nil
	})

	ts, _ := gotomerge.As[time.Time](doc.Get("created_at"))
	fmt.Println(ts.UTC().Format(time.RFC3339))

	// Output:
	// 2024-01-15T09:00:00Z
}

// Example_incrementalSave shows how to persist changes cheaply: SaveIncremental
// writes only the new change chunks since the last save, which can then be
// appended to a file or streamed to peers.
func Example_incrementalSave() {
	doc := gotomerge.NewDocument()
	doc.Change(func(txn *gotomerge.Txn) error {
		txn.Set("x", int64(1))
		return nil
	})

	// Full snapshot — use once at creation or after compaction.
	var full bytes.Buffer
	doc.Save(&full)

	// Apply more changes and save only the delta.
	doc.Change(func(txn *gotomerge.Txn) error {
		txn.Set("y", int64(2))
		return nil
	})

	var delta bytes.Buffer
	doc.SaveIncremental(&delta)

	// A reader can load snapshot + delta together to get the full state.
	combined := append(full.Bytes(), delta.Bytes()...)
	doc2, _ := gotomerge.LoadDocument(bytes.NewReader(combined))

	x, _ := gotomerge.As[int64](doc2.Get("x"))
	y, _ := gotomerge.As[int64](doc2.Get("y"))
	fmt.Println(x, y)

	// Output:
	// 1 2
}
