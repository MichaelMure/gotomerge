// Package rope implements a generic order-statistic tree (implicit treap).
//
// Despite the name, this is not a string rope — it is a general sequence
// structure that stores one element per node, using subtree sizes as implicit
// keys. It provides O(log n) expected time for positional access ([Rope.At]),
// insert ([Rope.InsertAfter], [Rope.PushBack], [Rope.PushFront]), and delete
// ([Rope.Remove]).
//
// Callers may hold [*Node] pointers for direct O(log n) access to a specific
// element without a position scan — for example, an external map from some
// domain key to *Node enables O(1) lookup followed by O(log n)
// InsertAfter/Remove. This is the primary use case the package is designed
// for.
package rope

import (
	"iter"
	"math/rand/v2"
)

// Node is a node in a [Rope]. Callers may hold *Node pointers across
// mutations: a node remains valid until it is passed to [Rope.Remove].
type Node[T any] struct {
	val      T
	priority uint32
	size     int     // subtree size including self
	left     *Node[T]
	right    *Node[T]
	parent   *Node[T]
}

// Value returns the value stored in n.
func (n *Node[T]) Value() T { return n.val }

// Rope is an implicit treap: a randomised binary search tree whose ordering
// key is the implicit in-order position of each node (derived from subtree
// sizes). All operations are O(log n) expected time.
type Rope[T any] struct {
	root *Node[T]
	len  int
}

// New returns an empty Rope.
func New[T any]() *Rope[T] { return &Rope[T]{} }

// Len returns the number of elements.
func (r *Rope[T]) Len() int { return r.len }

// PushBack appends v at the end and returns its Node.
func (r *Rope[T]) PushBack(v T) *Node[T] {
	n := newNode[T](v)
	r.root = merge(r.root, n)
	r.len++
	return n
}

// PushFront prepends v at the front and returns its Node.
func (r *Rope[T]) PushFront(v T) *Node[T] {
	n := newNode[T](v)
	r.root = merge(n, r.root)
	r.len++
	return n
}

// InsertAfter inserts v immediately after at and returns its Node.
func (r *Rope[T]) InsertAfter(v T, at *Node[T]) *Node[T] {
	n := newNode[T](v)
	l, right := split(r.root, rank(at)+1)
	r.root = merge(merge(l, n), right)
	r.len++
	return n
}

// Remove removes n from the rope. n must belong to this rope.
func (r *Rope[T]) Remove(n *Node[T]) {
	l, rest := split(r.root, rank(n))
	_, right := split(rest, 1)
	r.root = merge(l, right)
	r.len--
}

// At returns the node at 0-based position i, or nil if i is out of range.
func (r *Rope[T]) At(i int) *Node[T] {
	if i < 0 || i >= r.len {
		return nil
	}
	return nodeAt(r.root, i)
}

// Front returns the first node, or nil if the rope is empty.
func (r *Rope[T]) Front() *Node[T] {
	if r.root == nil {
		return nil
	}
	n := r.root
	for n.left != nil {
		n = n.left
	}
	return n
}

// Back returns the last node, or nil if the rope is empty.
func (r *Rope[T]) Back() *Node[T] {
	if r.root == nil {
		return nil
	}
	n := r.root
	for n.right != nil {
		n = n.right
	}
	return n
}

// All returns an iterator over values in order.
func (r *Rope[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		inorder(r.root, yield)
	}
}

// --- internal helpers -------------------------------------------------------

func newNode[T any](v T) *Node[T] {
	return &Node[T]{val: v, priority: rand.Uint32(), size: 1}
}

func sz[T any](n *Node[T]) int {
	if n == nil {
		return 0
	}
	return n.size
}

func pull[T any](n *Node[T]) {
	n.size = sz(n.left) + 1 + sz(n.right)
}

// rank returns the 0-based position of n in its rope by climbing to the root.
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

// nodeAt returns the node at 0-based position i within the subtree rooted at n.
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

// split splits the subtree rooted at n into left (first k elements) and right
// (the rest). The parent pointers of both returned roots are nil.
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
	return inorder(n.left, yield) && yield(n.val) && inorder(n.right, yield)
}
