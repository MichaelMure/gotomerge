package docproxy

import "fmt"

type parentView interface {
	// notify is called by a child proxy when a value is changed.
	// This cascades up to the document root, which can record changes.
	notify(key any, subkeys []string, changetype string, value any)

	// init is called by a child to ensure that the key exist with the given type.
	// This will panic with ErrType if the key already exists with a different type.
	init(key any, value any)

	// get returns the value for a given key.
	// This will panic with ErrType if the actual value in the document doesn't
	// match the expected type for the proxy.
	get(key any) any

	// set is called by a child proxy when a value is set.
	// This also act as a notify() to the parent, which can cascade as more notify().
	set(key any, value any, changetype string)
}

type ErrType struct {
	expected any
	got      any
}

func (e ErrType) Error() string {
	return fmt.Sprintf("value already assigned to %T, not %T", e.got, e.expected)
}
