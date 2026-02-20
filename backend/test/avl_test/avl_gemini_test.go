package avl_test

import (
	"math"
	"math/rand"
	"sort"
	"testing"

	"github.com/Durelius/next-week/internal/avl"
)

// Helper to check if the AVL height is within theoretical limits
func isBalanced(n int, height int) bool {
	if n == 0 {
		return height == 0
	}
	// AVL Max height formula: h <= 1.44 * log2(n + 2) - 0.328
	maxHeight := 1.44*math.Log2(float64(n)+2) - 0.328
	return float64(height) <= math.Max(float64(height), maxHeight)
}

func TestAVL_Extensive(t *testing.T) {
	t.Run("EdgeCases", func(t *testing.T) {
		tr := avl.New[int, string]()

		// Empty tree checks
		if tr.Size() != 0 {
			t.Errorf("Expected size 0, got %d", tr.Size())
		}
		if _, _, found := tr.Min(); found {
			t.Error("Min() on empty tree should be false")
		}

		// Single node
		tr.Insert(10, "ten")
		if tr.Height() != 1 {
			t.Errorf("Expected height 1, got %d", tr.Height())
		}

		// Delete non-existent
		tr.Delete(99)
		if tr.Size() != 1 {
			t.Error("Size changed after deleting non-existent key")
		}
	})

	t.Run("BulkInsertSequential", func(t *testing.T) {
		// Sequential insertion forces maximum rotations in a standard BST
		// but an AVL should handle it gracefully.
		tr := avl.New[int, int]()
		count := 10000
		for i := range count {
			tr.Insert(i, i*2)
		}

		if tr.Size() != count {
			t.Fatalf("Expected size %d, got %d", count, tr.Size())
		}

		if !isBalanced(tr.Size(), tr.Height()) {
			t.Errorf("Tree out of balance! Height: %d, Size: %d", tr.Height(), tr.Size())
		}
	})

	t.Run("BulkRandomChurn", func(t *testing.T) {
		tr := avl.New[int, int]()
		data := make(map[int]int)
		keys := []int{}

		count := 50000
		// 1. Insert 50k random elements
		for range count {
			key := rand.Intn(1000000)
			val := rand.Intn(100)
			tr.Insert(key, val)
			data[key] = val // Track expected state
			keys = append(keys, key)
		}

		// 2. Verify all exist
		for k := range data {
			if !tr.Contains(k) {
				t.Errorf("Tree missing key %d", k)
			}
		}

		// 3. Verify Min
		sort.Ints(keys)
		minK, _, _ := tr.Min()
		if minK != keys[0] {
			t.Errorf("Min mismatch: expected %d, got %d", keys[0], minK)
		}

		// 4. Delete half the elements and check balance
		for i := 0; i < count/2; i++ {
			tr.Delete(keys[i])
		}

		if !isBalanced(tr.Size(), tr.Height()) {
			t.Errorf("Unbalanced after deletions. Height: %d, Size: %d", tr.Height(), tr.Size())
		}
	})
}

func BenchmarkTree_Insert(b *testing.B) {
	tr := avl.New[int, int]()
	for i := 0; b.Loop(); i++ {
		tr.Insert(i, i)
	}
}
