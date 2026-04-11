package docproxy_test

import (
	"bytes"
	"fmt"
	"sort"

	"github.com/MichaelMure/gotomerge/docproxy"
)

// Example demonstrates the basic document lifecycle: create, write, read, save,
// and reload.
func Example() {
	doc := docproxy.NewDocument()

	doc.Change(func(txn *docproxy.Txn) error {
		txn.Set("title", "Hello")
		txn.Set("published", true)
		meta := txn.Map("meta") // obtain once; reuse for multiple writes
		meta.Set("author", "alice")
		meta.Set("version", int64(1))
		return nil
	})

	// Read scalars with As[T] — no type assertions needed.
	if title, ok := docproxy.As[string](doc.Get("title")); ok {
		fmt.Println(title)
	}
	if published, ok := docproxy.As[bool](doc.Get("published")); ok {
		fmt.Println(published)
	}

	// Read a nested struct in one shot.
	// Field names map to lowercase keys by default — no tags needed.
	type Meta struct {
		Author  string
		Version int64
	}
	if meta, ok := docproxy.As[Meta](doc.Get("meta")); ok {
		fmt.Println(meta.Author)
		fmt.Println(meta.Version)
	}

	// Native() lets any Value print without a cast.
	if v, ok := doc.Get("title"); ok {
		fmt.Println(v.Native())
	}

	// Save to bytes and reload: all state is preserved.
	var buf bytes.Buffer
	doc.Save(&buf)

	doc2, _ := docproxy.LoadDocument(&buf)
	if title, ok := docproxy.As[string](doc2.Get("title")); ok {
		fmt.Println(title)
	}

	// Output:
	// Hello
	// true
	// alice
	// 1
	// Hello
	// Hello
}

// ExampleDocument_Merge demonstrates concurrent editing: two peers start from
// the same state, each makes an independent change, and then the changes are
// merged. After the merge both peers' changes are visible.
func ExampleDocument_Merge() {
	base := docproxy.NewDocument()
	base.Change(func(txn *docproxy.Txn) error {
		txn.Set("shared", "base")
		return nil
	})

	// Fork: serialise base and load two independent peers from the same bytes.
	var snap bytes.Buffer
	base.SaveIncremental(&snap)
	snapBytes := snap.Bytes()

	peer1, _ := docproxy.LoadDocument(bytes.NewReader(snapBytes))
	peer2, _ := docproxy.LoadDocument(bytes.NewReader(snapBytes))

	// Each peer adds a distinct key independently.
	peer1.Change(func(txn *docproxy.Txn) error {
		txn.Set("alice", "peer1-edit")
		return nil
	})
	peer2.Change(func(txn *docproxy.Txn) error {
		txn.Set("bob", "peer2-edit")
		return nil
	})

	// Merge peer2's change into peer1.
	peer1.Merge(peer2)

	// All three keys are now visible on peer1.
	keys := peer1.Keys()
	sort.Strings(keys)
	fmt.Println(keys)

	shared, _ := docproxy.As[string](peer1.Get("shared"))
	fmt.Println(shared)

	// Output:
	// [alice bob shared]
	// base
}

// ExampleDocument_Change_text demonstrates text editing using Splice for
// precise insert/delete/replace and Update for whole-string replacement.
func ExampleDocument_Change_text() {
	doc := docproxy.NewDocument()

	// Insert initial content.
	doc.Change(func(txn *docproxy.Txn) error {
		txn.Text("body").Splice(0, 0, "Hello, world!")
		return nil
	})

	// Replace "world" (positions 7–11) with "Go".
	doc.Change(func(txn *docproxy.Txn) error {
		txn.Text("body").Splice(7, 5, "Go")
		return nil
	})

	if body, ok := docproxy.As[string](doc.Get("body")); ok {
		fmt.Println(body)
	}

	// Update replaces only the characters that differ (LCS diff internally).
	doc.Change(func(txn *docproxy.Txn) error {
		txn.Text("body").Update("Hi, Go!")
		return nil
	})

	if body, ok := docproxy.As[string](doc.Get("body")); ok {
		fmt.Println(body)
	}

	// Output:
	// Hello, Go!
	// Hi, Go!
}

// ExampleDocument_Change_list demonstrates appending to and reading from a list.
func ExampleDocument_Change_list() {
	doc := docproxy.NewDocument()

	doc.Change(func(txn *docproxy.Txn) error {
		l := txn.List("items")
		l.Append("a")
		l.Append("b")
		l.Append("c")
		return nil
	})

	list, _ := doc.List("items")
	fmt.Println(list.Len())

	// As[[]string] materialises the whole list in one call.
	if items, ok := docproxy.As[[]string](doc.Get("items")); ok {
		for i, s := range items {
			fmt.Printf("%d: %s\n", i, s)
		}
	}

	// Delete the middle element.
	doc.Change(func(txn *docproxy.Txn) error {
		txn.List("items").Delete(1)
		return nil
	})

	list, _ = doc.List("items")
	fmt.Println(list.Len())

	// Output:
	// 3
	// 0: a
	// 1: b
	// 2: c
	// 2
}
