// Package bitset provides a compact bit array backed by []uint64.
package bitset

// Bitset is a compact bit array (64 bits per word).
// The zero value is an empty, ready-to-use Bitset.
type Bitset struct {
	words []uint64
}

func (b *Bitset) Set(i uint32) {
	word := i / 64
	for word >= uint32(len(b.words)) {
		b.words = append(b.words, 0)
	}
	b.words[word] |= 1 << (i % 64)
}

func (b *Bitset) Get(i uint32) bool {
	word := i / 64
	if word >= uint32(len(b.words)) {
		return false
	}
	return b.words[word]>>(i%64)&1 != 0
}
