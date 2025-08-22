package types

// Key is either:
// - KeyOpId: ActorIdx + Counter
// - KeyString: String
type Key interface {
	isKey()
}

type KeyOpId OpId

func (KeyOpId) isKey() { panic("interface marker, don't call it") }

type KeyString string

func (KeyString) isKey() { panic("interface marker, don't call it") }

type NullKey struct{}

func (NullKey) isKey() { panic("interface marker, don't call it") }
