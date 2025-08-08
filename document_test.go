package gotomerge

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDocument(t *testing.T) {
	doc := NewDocument()
	err := doc.Change(func(doc *DocumentView) error {
		doc.Bool("foo").Set(true)
		doc.Map("bar").Map("baz").Bool("qux").Set(true)
		return nil
	})
	require.NoError(t, err)

	err = doc.Change(func(doc *DocumentView) error {
		// Fail: "foo" is already a bool
		doc.Map("foo").Bool("bar").Set(true)
		return nil
	})
	require.Error(t, err)

	fmt.Println(doc.changes)
}
