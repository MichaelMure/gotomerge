// Package treap implements a generic implicit treap (order-statistic tree).
//
// # What is a treap?
//
// A treap is a randomised binary search tree that simultaneously satisfies two
// invariants: the binary search tree (BST) property on keys, and the heap
// property on random priorities. Because priorities are random the tree stays
// balanced in expectation, giving O(log n) height with high probability.
//
// An "implicit" treap uses each node's subtree size as its implicit key rather
// than storing an explicit key. This makes the data structure a general-purpose
// ordered sequence: you can insert, delete, and access by 0-based position all
// in O(log n), with no key domain to worry about.
//
// # What this package provides
//
// [Treap] is a sequence (not a map or set). It stores one element per node in
// insertion order, with subtree sizes as the implicit in-order keys.
//
// All operations run in O(log n) expected time:
//   - Positional access: [Treap.At], [Treap.Front], [Treap.Back]
//   - Insert: [Treap.InsertAfter], [Treap.InsertBefore], [Treap.PushBack], [Treap.PushFront]
//   - Delete: [Treap.Remove]
//   - Traversal: [Treap.All]
//
// Callers may retain [*Node] pointers across mutations. A node pointer remains
// valid until it is passed to [Treap.Remove]. This enables O(1) external lookup
// (e.g. via a map keyed on some domain ID) followed by O(log n) InsertAfter or
// Remove — the primary use case this package is designed for.
//
// # Soft deletion (tombstones)
//
// [Treap.MarkDead] marks a node as dead without removing it from the sequence.
// Dead nodes remain in the treap as position markers — useful when other nodes
// reference them as predecessors — but are excluded from [Treap.Len],
// [Treap.LiveAt], [Treap.LiveRank], and [Treap.All]. A separate total-element
// count is available via [Treap.TotalLen] and positional access including dead
// nodes via [Treap.At].
package treap

import (
	"iter"
	"math/rand/v2"
)

// Node is a node in a [Treap]. Callers may hold *Node pointers across
// mutations: a node remains valid until it is passed to [Treap.Remove].
type Node[T any] struct {
	val      T
	priority uint32
	size     int // subtree size including self (live + dead)
	liveSize int // subtree size counting only alive nodes
	alive    bool
	left     *Node[T]
	right    *Node[T]
	parent   *Node[T]
}

// Value returns the value stored in n.
func (n *Node[T]) Value() T { return n.val }

// Alive reports whether n has not been marked dead via [Treap.MarkDead].
func (n *Node[T]) Alive() bool { return n.alive }

// Treap is an implicit treap: a randomised binary search tree whose ordering
// key is the implicit in-order position of each node (derived from subtree
// sizes). All operations are O(log n) expected time.
type Treap[T any] struct {
	root     *Node[T]
	totalLen int // all nodes, including dead
}

// New returns an empty Treap.
func New[T any]() *Treap[T] { return &Treap[T]{} }

// Len returns the number of live (non-dead) elements.
func (r *Treap[T]) Len() int { return liveSz(r.root) }

// TotalLen returns the total number of nodes, including dead ones.
func (r *Treap[T]) TotalLen() int { return r.totalLen }

// PushBack appends v at the end and returns its Node.
func (r *Treap[T]) PushBack(v T) *Node[T] {
	n := newNode[T](v)
	r.root = merge(r.root, n)
	r.totalLen++
	return n
}

// PushFront prepends v at the front and returns its Node.
func (r *Treap[T]) PushFront(v T) *Node[T] {
	n := newNode[T](v)
	r.root = merge(n, r.root)
	r.totalLen++
	return n
}

// InsertAfter inserts v immediately after at and returns its Node.
func (r *Treap[T]) InsertAfter(v T, at *Node[T]) *Node[T] {
	n := newNode[T](v)
	l, right := split(r.root, rank(at)+1)
	r.root = merge(merge(l, n), right)
	r.totalLen++
	return n
}

// InsertBefore inserts v immediately before at and returns its Node.
func (r *Treap[T]) InsertBefore(v T, at *Node[T]) *Node[T] {
	n := newNode[T](v)
	l, right := split(r.root, rank(at))
	r.root = merge(merge(l, n), right)
	r.totalLen++
	return n
}

// Remove removes n from the treap. n must belong to this treap.
func (r *Treap[T]) Remove(n *Node[T]) {
	l, rest := split(r.root, rank(n))
	_, right := split(rest, 1)
	r.root = merge(l, right)
	r.totalLen--
}

// MarkDead marks n as dead: it remains in the treap as a position marker but
// is excluded from Len, LiveAt, LiveRank, and All. Calling MarkDead on an
// already-dead node is a no-op. O(log n) due to the parent-pointer walk that
// updates liveSize up to the root.
func (r *Treap[T]) MarkDead(n *Node[T]) {
	if !n.alive {
		return
	}
	n.alive = false
	for cur := n; cur != nil; cur = cur.parent {
		cur.liveSize--
	}
}

// At returns the node at 0-based position i among all nodes (including dead),
// or nil if i is out of range.
func (r *Treap[T]) At(i int) *Node[T] {
	if i < 0 || i >= r.totalLen {
		return nil
	}
	return nodeAt(r.root, i)
}

// LiveAt returns the node at 0-based position i among live nodes only,
// or nil if i is out of range.
func (r *Treap[T]) LiveAt(i int) *Node[T] {
	if i < 0 || i >= liveSz(r.root) {
		return nil
	}
	return liveNodeAt(r.root, i)
}

// LiveRank returns the 0-based position of n among live nodes. If n is dead
// the result is the position it would occupy if it were alive (i.e. the number
// of live nodes before n).
func (r *Treap[T]) LiveRank(n *Node[T]) int {
	return liveRank(n)
}

// Next returns the next node in sequence order (including dead nodes), or nil
// if n is the last node. O(log n) amortized.
// TODO: make this a method on Node?
func Next[T any](n *Node[T]) *Node[T] {
	if n.right != nil {
		n = n.right
		for n.left != nil {
			n = n.left
		}
		return n
	}
	for n.parent != nil {
		if n == n.parent.left {
			return n.parent
		}
		n = n.parent
	}
	return nil
}

// Front returns the first live node, or nil if the treap is empty.
func (r *Treap[T]) Front() *Node[T] {
	if r.root == nil {
		return nil
	}
	n := r.root
	for n.left != nil {
		n = n.left
	}
	return n
}

// Back returns the last live node, or nil if the treap is empty.
func (r *Treap[T]) Back() *Node[T] {
	if r.root == nil {
		return nil
	}
	n := r.root
	for n.right != nil {
		n = n.right
	}
	return n
}

// All returns an iterator over live values in order, skipping dead nodes.
func (r *Treap[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		inorder(r.root, yield)
	}
}

// --- internal helpers -------------------------------------------------------

func newNode[T any](v T) *Node[T] {
	return &Node[T]{val: v, priority: rand.Uint32(), size: 1, liveSize: 1, alive: true}
}

func sz[T any](n *Node[T]) int {
	if n == nil {
		return 0
	}
	return n.size
}

func liveSz[T any](n *Node[T]) int {
	if n == nil {
		return 0
	}
	return n.liveSize
}

func pull[T any](n *Node[T]) {
	lw := 0
	if n.alive {
		lw = 1
	}
	n.size = sz(n.left) + 1 + sz(n.right)
	n.liveSize = liveSz(n.left) + lw + liveSz(n.right)
}

// rank returns the 0-based total position of n (counting dead nodes).
func rank[T any](n *Node[T]) int {
	r := sz(n.left)
	cur := n
	for cur.parent != nil {
		if cur == cur.parent.right {
			r += sz(cur.parent.left) + 1
		}
		cur = cur.parent
	}
	return r
}

// liveRank returns the 0-based live position of n (skipping dead nodes).
func liveRank[T any](n *Node[T]) int {
	r := liveSz(n.left)
	cur := n
	for cur.parent != nil {
		if cur == cur.parent.right {
			lw := 0
			if cur.parent.alive {
				lw = 1
			}
			r += liveSz(cur.parent.left) + lw
		}
		cur = cur.parent
	}
	return r
}

// nodeAt returns the node at 0-based total position i within the subtree rooted at n.
func nodeAt[T any](n *Node[T], i int) *Node[T] {
	for {
		ls := sz(n.left)
		switch {
		case i < ls:
			n = n.left
		case i == ls:
			return n
		default:
			i -= ls + 1
			n = n.right
		}
	}
}

// liveNodeAt returns the node at 0-based live position i within the subtree rooted at n.
func liveNodeAt[T any](n *Node[T], i int) *Node[T] {
	for n != nil {
		ls := liveSz(n.left)
		lw := 0
		if n.alive {
			lw = 1
		}
		switch {
		case i < ls:
			n = n.left
		case i < ls+lw:
			return n // i == ls and n is alive
		default:
			i -= ls + lw
			n = n.right
		}
	}
	return nil
}

// split splits the subtree rooted at n into left (first k total elements) and
// right (the rest). The parent pointers of both returned roots are nil.
func split[T any](n *Node[T], k int) (*Node[T], *Node[T]) {
	if n == nil {
		return nil, nil
	}
	n.parent = nil
	ls := sz(n.left)
	if k <= ls {
		l, r := split(n.left, k)
		n.left = r
		if r != nil {
			r.parent = n
		}
		pull(n)
		return l, n
	}
	l, r := split(n.right, k-ls-1)
	n.right = l
	if l != nil {
		l.parent = n
	}
	pull(n)
	return n, r
}

// merge merges two treaps where all elements in l precede r.
// The parent pointer of the returned root is nil.
func merge[T any](l, r *Node[T]) *Node[T] {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	if l.priority >= r.priority {
		l.parent = nil
		l.right = merge(l.right, r)
		if l.right != nil {
			l.right.parent = l
		}
		pull(l)
		return l
	}
	r.parent = nil
	r.left = merge(l, r.left)
	if r.left != nil {
		r.left.parent = r
	}
	pull(r)
	return r
}

func inorder[T any](n *Node[T], yield func(T) bool) bool {
	if n == nil {
		return true
	}
	if !inorder(n.left, yield) {
		return false
	}
	if n.alive && !yield(n.val) {
		return false
	}
	return inorder(n.right, yield)
}
