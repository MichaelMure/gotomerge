package types

import "fmt"

// Key is either:
// - KeyOpId: ActorIdx + Counter
// - KeyString: String
type Key interface {
	isKey()
}

type KeyOpId OpId

func (KeyOpId) isKey() { panic("interface marker, don't call it") }

func (k KeyOpId) String() string {
	return fmt.Sprintf("KeyOpId(actorIdx: %v, counter: %v)", k.ActorIdx, k.Counter)
}

type KeyString string

func (KeyString) isKey() { panic("interface marker, don't call it") }

func (k KeyString) String() string {
	return fmt.Sprintf("KeyString(%q)", string(k))
}

type NullKey struct{}

func (NullKey) isKey() { panic("interface marker, don't call it") }
