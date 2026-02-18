package avl

import (
	"cmp"
	"math"
	"math/rand"
	"testing"
	"time"
)

//
// ==========================
// Helpers
// ==========================
//

// BST validation
func isBST[T cmp.Ordered](n *node[T], min *T, maxVal *T) bool {
	if n == nil {
		return true
	}

	if min != nil && n.data <= *min {
		return false
	}
	if maxVal != nil && n.data >= *maxVal {
		return false
	}

	return isBST(n.left, min, &n.data) &&
		isBST(n.right, &n.data, maxVal)
}

// AVL validation (checks balance + height correctness)
func validateAVL[T cmp.Ordered](t *testing.T, n *node[T]) int {
	if n == nil {
		return 0
	}

	leftHeight := validateAVL(t, n.left)
	rightHeight := validateAVL(t, n.right)

	// Check balance property
	if math.Abs(float64(leftHeight-rightHeight)) > 1 {
		t.Fatalf("AVL balance violated at node %v", n.data)
	}

	// Check stored height correctness
	expectedHeight := 1 + max(leftHeight, rightHeight)
	if n.height != expectedHeight {
		t.Fatalf("Height mismatch at node %v: got %d expected %d",
			n.data, n.height, expectedHeight)
	}

	return expectedHeight
}

//
// ==========================
// Tests
// ==========================
//

func TestEmptyTree(t *testing.T) {
	var root *node[int]

	if root != nil {
		t.Fatal("expected nil root")
	}
}

func TestSingleInsert(t *testing.T) {
	var root *node[int]
	root = root.Insert(10)

	if root == nil {
		t.Fatal("root should not be nil")
	}

	if root.data != 10 {
		t.Fatal("incorrect root value")
	}

	validateAVL(t, root)
}

func TestLLRotation(t *testing.T) {
	var root *node[int]

	root = root.Insert(30)
	root = root.Insert(20)
	root = root.Insert(10)

	if root.data != 20 {
		t.Fatalf("expected root 20 after LL rotation, got %d", root.data)
	}

	validateAVL(t, root)
}

func TestRRRotation(t *testing.T) {
	var root *node[int]

	root = root.Insert(10)
	root = root.Insert(20)
	root = root.Insert(30)

	if root.data != 20 {
		t.Fatalf("expected root 20 after RR rotation, got %d", root.data)
	}

	validateAVL(t, root)
}

func TestLRRotation(t *testing.T) {
	var root *node[int]

	root = root.Insert(30)
	root = root.Insert(10)
	root = root.Insert(20)

	if root.data != 20 {
		t.Fatalf("expected root 20 after LR rotation, got %d", root.data)
	}

	validateAVL(t, root)
}

func TestRLRotation(t *testing.T) {
	var root *node[int]

	root = root.Insert(10)
	root = root.Insert(30)
	root = root.Insert(20)

	if root.data != 20 {
		t.Fatalf("expected root 20 after RL rotation, got %d", root.data)
	}

	validateAVL(t, root)
}

func TestSequentialInsert(t *testing.T) {
	var root *node[int]

	for i := 1; i <= 1000; i++ {
		root = root.Insert(i)
	}

	validateAVL(t, root)

	// Height must remain O(log n)
	maxHeight := int(1.45 * math.Log2(1000))
	if root.height > maxHeight {
		t.Fatalf("tree too tall: height=%d", root.height)
	}
}

func TestRandomInsert(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	var root *node[int]
	values := rand.Perm(2000)

	for _, v := range values {
		root = root.Insert(v)
	}

	validateAVL(t, root)

	if !isBST(root, nil, nil) {
		t.Fatal("BST property violated")
	}
}

func TestHugeInsert(t *testing.T) {
	var root *node[int]

	for i := 0; i < 10000; i++ {
		root = root.Insert(i)
	}

	validateAVL(t, root)

	maxHeight := int(1.45 * math.Log2(10000))
	if root.height > maxHeight {
		t.Fatalf("tree too tall for 10k elements: height=%d", root.height)
	}
}
