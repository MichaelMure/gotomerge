package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/MichaelMure/leb128"
)

// ActorId is a random byte string that uniquely identifies a peer (an instance
// that creates changes). It is conventionally 16 bytes of cryptographic random,
// but the format allows any length up to 32 bytes. It is not a UUID — it carries
// no version/variant bits and has no structured encoding; it is just an opaque
// random identifier that is unlikely to collide across peers.
type ActorId []byte

func NewActorId() ActorId {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return b
}

func (a ActorId) String() string {
	return hex.EncodeToString(a)
}

func ReadLengthEncodedActorId(r io.Reader) (ActorId, error) {
	l, err := leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading actor id length: %w", err)
	}
	if l > 32 {
		return nil, fmt.Errorf("unexpectedly large actor id length")
	}
	id := make([]byte, l)
	_, err = io.ReadFull(r, id[:])
	if err != nil {
		return nil, err
	}
	return id, nil
}
