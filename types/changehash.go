package types

import "encoding/hex"

// ChangeHash is the 32-byte SHA256 hash of the concatenation of:
// - the chunk type (0x01)
// - chunk length
// - chunk content
type ChangeHash [32]byte

func (ch ChangeHash) String() string {
	return hex.EncodeToString(ch[:])
}
