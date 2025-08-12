package types

import (
	"crypto/rand"
	"encoding/hex"
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
