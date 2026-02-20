package avl_test

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/Durelius/next-week/internal/avl"
)

const largeN = 10000

// Helper: verifies AVL theoretical height upper bound
func assertAVLHeightBound(t *testing.T, size int, height int) {
	if size == 0 {
		if height != 0 {
			t.Fatalf("expected height 0 for empty tree, got %d", height)
		}
		return
	}

	// AVL upper bound ≈ 1.44 * log2(n+2)
	maxHeight := int(1.44*math.Log2(float64(size+2))) + 2

	if height > maxHeight {
		t.Fatalf("height too large: got %d, expected <= %d (size=%d)",
			height, maxHeight, size)
	}
}

func TestEmptyTree(t *testing.T) {
	tree := avl.New[int, string]()

	if tree.Size() != 0 {
		t.Fatal("expected size 0")
	}

	if tree.Height() != 0 {
		t.Fatal("expected height 0")
	}

	if tree.Contains(10) {
		t.Fatal("empty tree should not contain key")
	}

	if _, ok := tree.Find(10); ok {
		t.Fatal("find should fail in empty tree")
	}

	if _, _, ok := tree.Min(); ok {
		t.Fatal("min should fail in empty tree")
	}
}

func TestSingleInsertDelete(t *testing.T) {
	tree := avl.New[int, string]()

	tree.Insert(42, "value")

	if !tree.Contains(42) {
		t.Fatal("tree should contain inserted key")
	}

	if tree.Size() != 1 {
		t.Fatal("size should be 1")
	}

	k, vals, ok := tree.Min()
	if !ok || k != 42 || len(vals) != 1 {
		t.Fatal("min incorrect")
	}

	tree.Delete(42)

	if tree.Size() != 0 {
		t.Fatal("tree should be empty after delete")
	}

	if tree.Contains(42) {
		t.Fatal("key should be deleted")
	}
}

func TestDuplicateKeysChatGPT(t *testing.T) {
	tree := avl.New[int, string]()

	for i := range 100 {
		tree.Insert(1, fmt.Sprintf("v%d", i))
	}

	vals, ok := tree.Find(1)
	if !ok {
		t.Fatal("expected duplicate key to exist")
	}

	if len(vals) != 100 {
		t.Fatalf("expected 100 values, got %d", len(vals))
	}

	if tree.Size() != 1 {
		t.Fatal("duplicate keys should not increase size")
	}
}

func TestAscendingInsert(t *testing.T) {
	tree := avl.New[int, int]()

	for i := range largeN {
		tree.Insert(i, i)
	}

	if tree.Size() != largeN {
		t.Fatal("size mismatch")
	}

	assertAVLHeightBound(t, tree.Size(), tree.Height())
}

func TestDescendingInsert(t *testing.T) {
	tree := avl.New[int, int]()

	for i := largeN; i >= 0; i-- {
		tree.Insert(i, i)
	}

	assertAVLHeightBound(t, tree.Size(), tree.Height())
}

func TestRandomLargeDataset(t *testing.T) {
	tree := avl.New[int, int]()
	rand.Seed(time.Now().UnixNano())

	keys := rand.Perm(largeN)

	for _, k := range keys {
		tree.Insert(k, k)
	}

	assertAVLHeightBound(t, tree.Size(), tree.Height())

	// verify all exist
	for _, k := range keys {
		if !tree.Contains(k) {
			t.Fatalf("missing key %d", k)
		}
	}
}

func TestMassDeletion(t *testing.T) {
	tree := avl.New[int, int]()

	for i := range largeN {
		tree.Insert(i, i)
	}

	for i := range largeN {
		tree.Delete(i)
		assertAVLHeightBound(t, tree.Size(), tree.Height())
	}

	if tree.Size() != 0 {
		t.Fatal("tree should be empty after mass deletion")
	}
}

func TestRepeatedDeleteSameKey(t *testing.T) {
	tree := avl.New[int, int]()

	tree.Insert(5, 1)
	tree.Delete(5)
	tree.Delete(5)
	tree.Delete(5)

	if tree.Size() != 0 {
		t.Fatal("size should remain 0")
	}
}

func TestMinAfterRandomOps(t *testing.T) {
	tree := avl.New[int, int]()
	rand.Seed(42)

	values := rand.Perm(largeN)

	for _, v := range values {
		tree.Insert(v, v)
	}

	sort.Ints(values)

	k, _, ok := tree.Min()
	if !ok || k != values[0] {
		t.Fatal("min incorrect")
	}
}

func TestInterleavedOperations(t *testing.T) {
	tree := avl.New[int, int]()
	rand.Seed(time.Now().UnixNano())

	for range 20000 {
		k := rand.Intn(5000)

		switch rand.Intn(3) {
		case 0:
			tree.Insert(k, k)
		case 1:
			tree.Delete(k)
		case 2:
			tree.Contains(k)
		}

		assertAVLHeightBound(t, tree.Size(), tree.Height())
	}
}

func TestStringStability(t *testing.T) {
	tree := avl.New[int, int]()

	for i := range 1000 {
		tree.Insert(i, i)
	}

	s1 := tree.String()
	s2 := tree.String()

	if s1 != s2 {
		t.Fatal("string representation should be deterministic")
	}
}

func TestStressMillionOps(t *testing.T) {
	tree := avl.New[int, int]()
	rand.Seed(123)

	for i := range 100000 {
		tree.Insert(rand.Intn(50000), i)
	}

	for range 100000 {
		tree.Delete(rand.Intn(50000))
	}

	assertAVLHeightBound(t, tree.Size(), tree.Height())
}
