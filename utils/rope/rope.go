// Package rope implements a chunked rope for efficiently maintaining and
// reading rendered text as a flat list of character chunks.
package rope

import "strings"

// chunkCap is the target maximum number of characters per chunk.
// Chunks above this size are split in two; chunks below it are left as-is
// (no merging) to keep the implementation simple.
const chunkCap = 128

// Rope maintains text as a flat list of character chunks. Each chunk holds up
// to chunkCap characters (one string per insert op). This avoids the O(n)
// pointer-chasing cost of treap traversal when reading text: Text()
// concatenates n/chunkCap chunks instead of visiting n individual treap nodes.
//
// Insert and delete are O(n/K + K): a linear scan to locate the target chunk
// plus a slice operation within the chunk. For K=128 and n=100k, that is
// about 800+128 ≈ 930 simple slice operations.
type Rope struct {
	chunks []chunk
	total  int // total number of live characters
}

type chunk struct {
	chars []string // one element per live character (one insert op's value)
}

// PushBack appends a character at the end of the rope.
// Used only during list treap construction where elements arrive in order.
func (r *Rope) PushBack(char string) {
	if len(r.chunks) == 0 || len(r.chunks[len(r.chunks)-1].chars) >= chunkCap {
		r.chunks = append(r.chunks, chunk{chars: make([]string, 0, chunkCap)})
	}
	last := len(r.chunks) - 1
	r.chunks[last].chars = append(r.chunks[last].chars, char)
	r.total++
}

// InsertAt inserts char at live position pos (0-based).
// pos == r.total appends at the end.
func (r *Rope) InsertAt(pos int, char string) {
	r.total++

	if len(r.chunks) == 0 {
		r.chunks = []chunk{{chars: []string{char}}}
		return
	}

	ci, off := r.locate(pos)

	if ci == len(r.chunks) {
		// Append to last chunk.
		ci = len(r.chunks) - 1
		off = len(r.chunks[ci].chars)
	}

	ch := &r.chunks[ci]
	ch.chars = append(ch.chars, "")
	copy(ch.chars[off+1:], ch.chars[off:])
	ch.chars[off] = char

	r.maybeSplit(ci)
}

// DeleteAt removes the character at live position pos (0-based).
func (r *Rope) DeleteAt(pos int) {
	if len(r.chunks) == 0 || r.total == 0 {
		return
	}
	r.total--

	ci, off := r.locate(pos)
	if ci >= len(r.chunks) {
		return
	}

	ch := &r.chunks[ci]
	ch.chars = append(ch.chars[:off], ch.chars[off+1:]...)

	if len(ch.chars) == 0 {
		r.chunks = append(r.chunks[:ci], r.chunks[ci+1:]...)
	}
}

// Text returns the full text by concatenating all chunks.
// O(n/K) concatenations of K-char segments.
func (r *Rope) Text() string {
	if r.total == 0 {
		return ""
	}
	var b strings.Builder
	b.Grow(r.total)
	for i := range r.chunks {
		for _, ch := range r.chunks[i].chars {
			b.WriteString(ch)
		}
	}
	return b.String()
}

// locate returns the chunk index and offset within that chunk for live
// position pos. Returns (len(chunks), 0) when pos == total (append-at-end).
func (r *Rope) locate(pos int) (ci int, off int) {
	for i := range r.chunks {
		n := len(r.chunks[i].chars)
		if pos < n {
			return i, pos
		}
		pos -= n
	}
	return len(r.chunks), 0
}

// maybeSplit splits chunk ci in two if it exceeds chunkCap.
func (r *Rope) maybeSplit(ci int) {
	if len(r.chunks[ci].chars) <= chunkCap {
		return
	}
	mid := len(r.chunks[ci].chars) / 2
	left := r.chunks[ci].chars[:mid:mid]
	right := make([]string, len(r.chunks[ci].chars)-mid)
	copy(right, r.chunks[ci].chars[mid:])

	r.chunks[ci].chars = left

	// Insert a new chunk after ci.
	r.chunks = append(r.chunks, chunk{})
	copy(r.chunks[ci+2:], r.chunks[ci+1:])
	r.chunks[ci+1] = chunk{chars: right}
}
