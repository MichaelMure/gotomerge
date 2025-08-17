package types

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/jcalabro/leb128"
)

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
