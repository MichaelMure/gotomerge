package gotomerge

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

func readActorIds(r io.Reader) ([]ActorId, error) {
	n, err := leb128.DecodeU64(r)
	if err != nil {
		return nil, fmt.Errorf("error reading actor ids length: %w", err)
	}
	// limit pre-allocation to avoid DOS
	allocate := n
	if n > 128 {
		allocate = 128
	}
	res := make([]ActorId, 0, allocate)

	for i := uint64(0); i < n; i++ {
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
			return nil, fmt.Errorf("error reading actor id: %w", err)
		}
		res = append(res, id)
	}
	return res, nil
}

func writeActorIds(w io.Writer, ids []ActorId) error {
	_, err := w.Write(leb128.EncodeU64(uint64(len(ids))))
	if err != nil {
		return err
	}
	for _, id := range ids {
		_, err = w.Write(leb128.EncodeU64(uint64(len(id))))
		if err != nil {
			return err
		}
		_, err = w.Write(id)
		if err != nil {
			return err
		}
	}
	return nil
}
