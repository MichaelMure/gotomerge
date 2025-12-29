package docproxy

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

	require.Equal(t, map[string]any{
		"foo": true,
		"bar": map[string]any{
			"baz": map[string]any{
				"qux": true,
			},
		},
	}, doc.Value())

	err = doc.Change(func(doc *DocumentView) error {
		// Fail: "foo" is already a bool
		doc.Map("foo").Bool("bar").Set(true)
		return nil
	})
	require.Error(t, err)

	require.Equal(t, map[string]any{
		"foo": true,
		"bar": map[string]any{
			"baz": map[string]any{
				"qux": true,
			},
		},
	}, doc.Value())

	fmt.Println(doc.changes)
}
