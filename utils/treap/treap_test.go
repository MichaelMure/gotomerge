package treap

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// collect drains All() into a slice.
func collect[T any](r *Treap[T]) []T {
	var out []T
	for v := range r.All() {
		out = append(out, v)
	}
	return out
}

// checkInvariants verifies size metadata and parent pointers throughout the tree.
func checkInvariants[T any](t *testing.T, r *Treap[T]) {
	t.Helper()
	var check func(n, parent *Node[T]) int
	check = func(n, parent *Node[T]) int {
		if n == nil {
			return 0
		}
		assert.Equal(t, parent, n.parent, "wrong parent pointer")
		ls := check(n.left, n)
		rs := check(n.right, n)
		assert.Equal(t, ls+1+rs, n.size, "wrong subtree size")
		if n.left != nil {
			assert.GreaterOrEqual(t, n.priority, n.left.priority, "heap violation (left)")
		}
		if n.right != nil {
			assert.GreaterOrEqual(t, n.priority, n.right.priority, "heap violation (right)")
		}
		return n.size
	}
	total := check(r.root, nil)
	assert.Equal(t, r.len, total, "Len mismatch")
}

func TestPushBack(t *testing.T) {
	r := New[int]()
	for i := range 5 {
		r.PushBack(i)
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4}, collect(r))
	assert.Equal(t, 5, r.Len())
	checkInvariants(t, r)
}

func TestPushFront(t *testing.T) {
	r := New[int]()
	for i := range 5 {
		r.PushFront(i)
	}
	assert.Equal(t, []int{4, 3, 2, 1, 0}, collect(r))
	checkInvariants(t, r)
}

func TestInsertAfter(t *testing.T) {
	r := New[int]()
	n0 := r.PushBack(0)
	n2 := r.PushBack(2)
	r.InsertAfter(1, n0)  // [0, 1, 2]
	r.InsertAfter(3, n2)  // [0, 1, 2, 3]
	assert.Equal(t, []int{0, 1, 2, 3}, collect(r))
	checkInvariants(t, r)
}

func TestInsertAfterLast(t *testing.T) {
	r := New[int]()
	n0 := r.PushBack(0)
	r.InsertAfter(1, n0)
	r.InsertAfter(2, r.Back())
	assert.Equal(t, []int{0, 1, 2}, collect(r))
	checkInvariants(t, r)
}

func TestRemoveMiddle(t *testing.T) {
	r := New[int]()
	r.PushBack(0)
	n1 := r.PushBack(1)
	r.PushBack(2)
	r.Remove(n1)
	assert.Equal(t, []int{0, 2}, collect(r))
	assert.Equal(t, 2, r.Len())
	checkInvariants(t, r)
}

func TestRemoveFront(t *testing.T) {
	r := New[int]()
	n0 := r.PushBack(0)
	r.PushBack(1)
	r.Remove(n0)
	assert.Equal(t, []int{1}, collect(r))
	checkInvariants(t, r)
}

func TestRemoveBack(t *testing.T) {
	r := New[int]()
	r.PushBack(0)
	n1 := r.PushBack(1)
	r.Remove(n1)
	assert.Equal(t, []int{0}, collect(r))
	checkInvariants(t, r)
}

func TestRemoveOnly(t *testing.T) {
	r := New[int]()
	n := r.PushBack(42)
	r.Remove(n)
	assert.Equal(t, 0, r.Len())
	assert.Nil(t, r.Front())
	assert.Nil(t, r.Back())
	checkInvariants(t, r)
}

func TestAt(t *testing.T) {
	r := New[int]()
	for i := range 10 {
		r.PushBack(i)
	}
	for i := range 10 {
		n := r.At(i)
		require.NotNil(t, n)
		assert.Equal(t, i, n.Value())
	}
	assert.Nil(t, r.At(-1))
	assert.Nil(t, r.At(10))
}

func TestAtAfterMutations(t *testing.T) {
	r := New[int]()
	nodes := make([]*Node[int], 5)
	for i := range 5 {
		nodes[i] = r.PushBack(i) // [0,1,2,3,4]
	}
	r.Remove(nodes[2]) // [0,1,3,4]
	r.InsertAfter(9, nodes[1]) // [0,1,9,3,4]

	assert.Equal(t, []int{0, 1, 9, 3, 4}, collect(r))
	assert.Equal(t, 9, r.At(2).Value())
	checkInvariants(t, r)
}

func TestFrontBack(t *testing.T) {
	r := New[int]()
	assert.Nil(t, r.Front())
	assert.Nil(t, r.Back())

	r.PushBack(1)
	r.PushBack(2)
	r.PushBack(3)
	assert.Equal(t, 1, r.Front().Value())
	assert.Equal(t, 3, r.Back().Value())
}

func TestRankConsistency(t *testing.T) {
	r := New[int]()
	nodes := make([]*Node[int], 20)
	for i := range 20 {
		nodes[i] = r.PushBack(i)
	}
	checkInvariants(t, r)
	for i, n := range nodes {
		assert.Equal(t, i, rank(n), "rank mismatch for element %d", i)
	}
}

func TestNodePointerStability(t *testing.T) {
	// Inserting and removing other elements must not invalidate held *Node pointers.
	r := New[int]()
	n5 := r.PushBack(5)

	for i := range 100 {
		r.PushBack(i)
	}
	for i := range 50 {
		r.InsertAfter(i*-1, n5)
	}

	assert.Equal(t, 5, n5.Value())
	checkInvariants(t, r)
}

func TestLargeSequential(t *testing.T) {
	const n = 10_000
	r := New[int]()
	for i := range n {
		r.PushBack(i)
	}
	assert.Equal(t, n, r.Len())
	checkInvariants(t, r)

	got := collect(r)
	want := make([]int, n)
	for i := range n {
		want[i] = i
	}
	assert.Equal(t, want, got)
}

func TestAllEarlyExit(t *testing.T) {
	r := New[int]()
	for i := range 10 {
		r.PushBack(i)
	}
	var got []int
	for v := range r.All() {
		got = append(got, v)
		if v == 4 {
			break
		}
	}
	assert.Equal(t, []int{0, 1, 2, 3, 4}, got)
}

func TestExternalMap(t *testing.T) {
	// Simulate the workingList use case: external map from key → *Node for
	// O(1) predecessor lookup, then O(log n) InsertAfter.
	r := New[string]()
	byKey := make(map[int]*Node[string])

	byKey[0] = r.PushBack("a")
	byKey[1] = r.PushBack("c")
	byKey[2] = r.InsertAfter("b", byKey[0]) // insert "b" after "a"

	assert.Equal(t, []string{"a", "b", "c"}, collect(r))

	r.Remove(byKey[2])
	assert.Equal(t, []string{"a", "c"}, collect(r))
	checkInvariants(t, r)
}

func TestInsertAfterEveryPosition(t *testing.T) {
	// Build [0..n-1] by always inserting after the previous node.
	const n = 500
	r := New[int]()
	prev := r.PushBack(0)
	for i := 1; i < n; i++ {
		prev = r.InsertAfter(i, prev)
	}
	assert.Equal(t, n, r.Len())
	got := collect(r)
	assert.True(t, slices.Equal(got, func() []int {
		s := make([]int, n)
		for i := range n {
			s[i] = i
		}
		return s
	}()))
	checkInvariants(t, r)
}
