package avl_test

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"testing"

	"github.com/Durelius/next-week/internal/avl"
)

// ─────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────

func newIntTree() *avl.Tree[int, int] {
	return avl.New[int, int]()
}

func newStringTree() *avl.Tree[string, string] {
	return avl.New[string, string]()
}

// insertAll inserts key→value pairs where value == key.
func insertAll(t *avl.Tree[int, int], keys []int) {
	for _, k := range keys {
		t.Insert(k, k)
	}
}

// sortedUnique returns the sorted, de-duplicated slice.
func sortedUnique(keys []int) []int {
	seen := map[int]struct{}{}
	out := []int{}
	for _, k := range keys {
		if _, ok := seen[k]; !ok {
			seen[k] = struct{}{}
			out = append(out, k)
		}
	}
	sort.Ints(out)
	return out
}

// ─────────────────────────────────────────────
// 1. Empty Tree
// ─────────────────────────────────────────────

func TestEmptyTreeClaude(t *testing.T) {
	tr := newIntTree()

	if got := tr.Size(); got != 0 {
		t.Errorf("empty Size: want 0, got %d", got)
	}
	if got := tr.Height(); got != 0 {
		t.Errorf("empty Height: want 0, got %d", got)
	}
	if _, _, ok := tr.Min(); ok {
		t.Error("empty Min: want ok=false")
	}
	if _, ok := tr.Find(0); ok {
		t.Error("empty Find: want ok=false")
	}
	if tr.Contains(0) {
		t.Error("empty Contains: want false")
	}
	// String / Print must not panic on empty Tree
	_ = tr.String()
	tr.Print()
}

// ─────────────────────────────────────────────
// 2. Single element
// ─────────────────────────────────────────────

func TestSingleElement(t *testing.T) {
	tr := newIntTree()
	tr.Insert(42, 100)

	if tr.Size() != 1 {
		t.Errorf("size after 1 insert: want 1, got %d", tr.Size())
	}
	if tr.Height() != 1 {
		t.Errorf("height after 1 insert: want 1, got %d", tr.Height())
	}

	k, vals, ok := tr.Min()
	if !ok || k != 42 {
		t.Errorf("Min: want (42,true), got (%d,%v)", k, ok)
	}
	if len(vals) != 1 || vals[0] != 100 {
		t.Errorf("Min values: want [100], got %v", vals)
	}

	vals, ok = tr.Find(42)
	if !ok || len(vals) == 0 || vals[0] != 100 {
		t.Errorf("Find(42): want ([100],true), got (%v,%v)", vals, ok)
	}
	if !tr.Contains(42) {
		t.Error("Contains(42): want true")
	}
	if tr.Contains(0) {
		t.Error("Contains(0) on single-node Tree: want false")
	}
}

// ─────────────────────────────────────────────
// 3. Insert / Find / Contains – basic
// ─────────────────────────────────────────────

func TestInsertAndFind(t *testing.T) {
	tr := newIntTree()
	n := 1000
	for i := range n {
		tr.Insert(i, i*10)
	}

	if tr.Size() != n {
		t.Errorf("Size: want %d, got %d", n, tr.Size())
	}
	for i := range n {
		vals, ok := tr.Find(i)
		if !ok {
			t.Errorf("Find(%d): want ok=true", i)
			continue
		}
		found := false
		for _, v := range vals {
			if v == i*10 {
				found = true
			}
		}
		if !found {
			t.Errorf("Find(%d): value %d not present in %v", i, i*10, vals)
		}
	}
}

// ─────────────────────────────────────────────
// 4. Duplicate keys – multiple values per key
// ─────────────────────────────────────────────

func TestDuplicateKeys(t *testing.T) {
	tr := newIntTree()
	const key = 7
	const copies = 50

	for i := range copies {
		tr.Insert(key, i)
	}

	vals, ok := tr.Find(key)
	if !ok {
		t.Fatal("Find after duplicate inserts: want ok=true")
	}
	if len(vals) != copies {
		t.Errorf("expected %d values for duplicate key, got %d", copies, len(vals))
	}
	// Size may count distinct keys or total inserts depending on implementation.
	// We only assert it is > 0.
	if tr.Size() == 0 {
		t.Error("Size must be > 0 after inserts")
	}
}

// ─────────────────────────────────────────────
// 5. Delete – basic
// ─────────────────────────────────────────────

func TestDeleteBasic(t *testing.T) {
	tr := newIntTree()
	keys := []int{10, 5, 15, 3, 7, 12, 20}
	insertAll(tr, keys)

	for _, k := range keys {
		tr.Delete(k)
		if tr.Contains(k) {
			t.Errorf("after Delete(%d): Contains still true", k)
		}
	}
	if tr.Size() != 0 {
		t.Errorf("after deleting all keys: want Size=0, got %d", tr.Size())
	}
}

// ─────────────────────────────────────────────
// 6. Delete non-existent key (must not panic)
// ─────────────────────────────────────────────

func TestDeleteNonExistent(t *testing.T) {
	tr := newIntTree()
	insertAll(tr, []int{1, 2, 3})
	sizeBefore := tr.Size()

	// Should not panic
	tr.Delete(999)
	tr.Delete(-1)
	tr.Delete(0)

	if tr.Size() != sizeBefore {
		t.Errorf("Size changed after deleting non-existent keys: was %d, now %d", sizeBefore, tr.Size())
	}
}

// ─────────────────────────────────────────────
// 7. Delete on empty Tree (must not panic)
// ─────────────────────────────────────────────

func TestDeleteEmptyTree(t *testing.T) {
	tr := newIntTree()
	tr.Delete(1) // should not panic
}

// ─────────────────────────────────────────────
// 8. Min after operations
// ─────────────────────────────────────────────

func TestMin(t *testing.T) {
	tr := newIntTree()
	keys := []int{50, 30, 70, 10, 40, 60, 80}
	insertAll(tr, keys)

	k, _, ok := tr.Min()
	if !ok || k != 10 {
		t.Errorf("Min: want (10, true), got (%d, %v)", k, ok)
	}

	// Delete the minimum; new min should be 30
	tr.Delete(10)
	k, _, ok = tr.Min()
	if !ok || k != 30 {
		t.Errorf("Min after deleting 10: want (30, true), got (%d, %v)", k, ok)
	}
}

// ─────────────────────────────────────────────
// 9. Height bounds (AVL invariant)
//    For n nodes, height <= 1.44 * log2(n+2)
// ─────────────────────────────────────────────

func TestHeightAVLBound(t *testing.T) {
	tr := newIntTree()
	n := 10_000
	for i := 1; i <= n; i++ {
		tr.Insert(i, i)
	}
	h := tr.Height()
	// Upper bound: 1.44 * log2(n+2) + 1  (generous)
	logBound := int(math.Ceil(1.44*math.Log2(float64(n+2)))) + 2
	if h > logBound {
		t.Errorf("Height %d exceeds AVL bound ~%d for n=%d", h, logBound, n)
	}
	if h < 1 {
		t.Errorf("Height must be >= 1 for non-empty Tree, got %d", h)
	}
}

// ─────────────────────────────────────────────
// 10. Sequential (worst-case for BST) insertion
// ─────────────────────────────────────────────

func TestSequentialInsert(t *testing.T) {
	tr := newIntTree()
	n := 5_000
	for i := range n {
		tr.Insert(i, i)
	}
	if tr.Size() != n {
		t.Errorf("sequential insert Size: want %d, got %d", n, tr.Size())
	}
	// AVL should keep height well below n
	if tr.Height() >= n {
		t.Errorf("sequential insert height %d: Tree is not balanced", tr.Height())
	}
	k, _, ok := tr.Min()
	if !ok || k != 0 {
		t.Errorf("sequential Min: want (0,true), got (%d,%v)", k, ok)
	}
}

// ─────────────────────────────────────────────
// 11. Reverse sequential insertion
// ─────────────────────────────────────────────

func TestReverseSequentialInsert(t *testing.T) {
	tr := newIntTree()
	n := 5_000
	for i := n - 1; i >= 0; i-- {
		tr.Insert(i, i)
	}
	if tr.Size() != n {
		t.Errorf("reverse insert Size: want %d, got %d", n, tr.Size())
	}
	if tr.Height() >= n {
		t.Errorf("reverse insert height %d: Tree is not balanced", tr.Height())
	}
}

// ─────────────────────────────────────────────
// 12. Random insertions and deletions
// ─────────────────────────────────────────────

func TestRandomInsertDelete(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	tr := newIntTree()
	present := map[int]bool{}

	// Insert 2000 random keys
	for range 2000 {
		k := rng.Intn(500) // intentional collisions
		tr.Insert(k, k)
		present[k] = true
	}

	// Verify all inserted keys are found
	for k := range present {
		if !tr.Contains(k) {
			t.Errorf("Contains(%d): want true after insert", k)
		}
	}

	// Delete half the keys
	deleted := []int{}
	i := 0
	for k := range present {
		if i%2 == 0 {
			tr.Delete(k)
			deleted = append(deleted, k)
		}
		i++
	}
	for _, k := range deleted {
		delete(present, k)
	}

	// Verify deleted are gone, remaining are present
	for _, k := range deleted {
		if _, ok := tr.Find(k); ok {
			t.Errorf("Find(%d): should be gone after delete", k)
		}
	}
	for k := range present {
		if !tr.Contains(k) {
			t.Errorf("Contains(%d): want true after partial delete", k)
		}
	}
}

// ─────────────────────────────────────────────
// 13. Insert → delete all → re-insert
// ─────────────────────────────────────────────

func TestReinsertAfterFullDelete(t *testing.T) {
	tr := newIntTree()
	keys := make([]int, 200)
	for i := range keys {
		keys[i] = i
	}
	insertAll(tr, keys)

	for _, k := range keys {
		tr.Delete(k)
	}
	if tr.Size() != 0 {
		t.Errorf("after full delete: Size want 0, got %d", tr.Size())
	}

	// Re-insert
	insertAll(tr, keys)
	if tr.Size() != len(keys) {
		t.Errorf("after re-insert: Size want %d, got %d", len(keys), tr.Size())
	}
	for _, k := range keys {
		if !tr.Contains(k) {
			t.Errorf("Contains(%d) after re-insert: want true", k)
		}
	}
}

// ─────────────────────────────────────────────
// 14. String keys
// ─────────────────────────────────────────────

func TestStringKeys(t *testing.T) {
	tr := newStringTree()
	words := []string{
		"banana", "apple", "cherry", "date", "elderberry",
		"fig", "grape", "honeydew", "kiwi", "lemon",
	}
	for _, w := range words {
		tr.Insert(w, w+"_val")
	}

	if tr.Size() != len(words) {
		t.Errorf("string Tree Size: want %d, got %d", len(words), tr.Size())
	}

	k, _, ok := tr.Min()
	if !ok || k != "apple" {
		t.Errorf("string Min: want apple, got %s (ok=%v)", k, ok)
	}

	vals, ok := tr.Find("grape")
	if !ok || len(vals) == 0 {
		t.Errorf("Find(grape): want ok=true with values")
	}

	tr.Delete("banana")
	if tr.Contains("banana") {
		t.Error("after deleting 'banana': Contains should be false")
	}
}

// ─────────────────────────────────────────────
// 15. Negative keys
// ─────────────────────────────────────────────

func TestNegativeKeys(t *testing.T) {
	tr := newIntTree()
	for i := -500; i <= 500; i++ {
		tr.Insert(i, i)
	}

	k, _, ok := tr.Min()
	if !ok || k != -500 {
		t.Errorf("Min with negatives: want -500, got %d (ok=%v)", k, ok)
	}
	if tr.Size() != 1001 {
		t.Errorf("Size with negatives: want 1001, got %d", tr.Size())
	}

	// Delete all negatives
	for i := -500; i < 0; i++ {
		tr.Delete(i)
	}
	k, _, ok = tr.Min()
	if !ok || k != 0 {
		t.Errorf("Min after deleting negatives: want 0, got %d (ok=%v)", k, ok)
	}
}

// ─────────────────────────────────────────────
// 16. Find returns all values for a key
// ─────────────────────────────────────────────

func TestFindAllValues(t *testing.T) {
	tr := newIntTree()
	const key = 99
	expected := []int{1, 2, 3, 4, 5, 100, 200, 300}
	for _, v := range expected {
		tr.Insert(key, v)
	}

	vals, ok := tr.Find(key)
	if !ok {
		t.Fatal("Find: want ok=true")
	}
	if len(vals) != len(expected) {
		t.Fatalf("Find: want %d values, got %d: %v", len(expected), len(vals), vals)
	}

	got := make([]int, len(vals))
	copy(got, vals)
	sort.Ints(got)
	sort.Ints(expected)
	for i := range expected {
		if got[i] != expected[i] {
			t.Errorf("values mismatch at %d: want %d, got %d", i, expected[i], got[i])
		}
	}
}

// ─────────────────────────────────────────────
// 17. Large-scale stress test
// ─────────────────────────────────────────────

func TestLargeScaleStress(t *testing.T) {
	rng := rand.New(rand.NewSource(1337))
	tr := newIntTree()
	truth := map[int]bool{}

	const ops = 50_000
	for range ops {
		k := rng.Intn(10_000)
		op := rng.Intn(3)
		switch op {
		case 0, 1: // insert (2/3 of ops)
			tr.Insert(k, k)
			truth[k] = true
		case 2: // delete
			tr.Delete(k)
			delete(truth, k)
		}
	}

	for k := range truth {
		if !tr.Contains(k) {
			t.Errorf("stress: Contains(%d) should be true", k)
		}
	}
	if tr.Size() != len(truth) {
		t.Errorf("stress: Size want %d, got %d", len(truth), tr.Size())
	}
}

// ─────────────────────────────────────────────
// 18. Size tracks correctly through inserts/deletes
// ─────────────────────────────────────────────

func TestSizeTracking(t *testing.T) {
	tr := newIntTree()
	for i := range 100 {
		tr.Insert(i, i)
		if tr.Size() != i+1 {
			t.Errorf("after inserting %d elements: want Size=%d, got %d", i+1, i+1, tr.Size())
		}
	}
	for i := 99; i >= 0; i-- {
		tr.Delete(i)
		if tr.Size() != i {
			t.Errorf("after deleting down to %d: want Size=%d, got %d", i, i, tr.Size())
		}
	}
}

// ─────────────────────────────────────────────
// 19. Height is 0 on empty, 1 on single node
// ─────────────────────────────────────────────

func TestHeightEdgeCases(t *testing.T) {
	tr := newIntTree()
	if tr.Height() != 0 {
		t.Errorf("empty Tree height: want 0, got %d", tr.Height())
	}
	tr.Insert(1, 1)
	if tr.Height() != 1 {
		t.Errorf("single-node height: want 1, got %d", tr.Height())
	}
	tr.Delete(1)
	if tr.Height() != 0 {
		t.Errorf("height after deleting only node: want 0, got %d", tr.Height())
	}
}

// ─────────────────────────────────────────────
// 20. Min reflects the actual minimum key
// ─────────────────────────────────────────────

func TestMinAlwaysCorrect(t *testing.T) {
	rng := rand.New(rand.NewSource(77))
	tr := newIntTree()
	nums := make([]int, 300)
	for i := range nums {
		nums[i] = rng.Intn(10_000) - 5_000
		tr.Insert(nums[i], nums[i])
	}

	unique := sortedUnique(nums)
	k, _, ok := tr.Min()
	if !ok {
		t.Fatal("Min: want ok=true")
	}
	if k != unique[0] {
		t.Errorf("Min: want %d, got %d", unique[0], k)
	}
}

// ─────────────────────────────────────────────
// 21. Delete root repeatedly
// ─────────────────────────────────────────────

func TestDeleteRootRepeatedly(t *testing.T) {
	tr := newIntTree()
	insertAll(tr, []int{5, 3, 7, 1, 4, 6, 8})

	// Keep deleting the minimum (which in a sorted AVL will often be the root
	// after rotations) to stress root-deletion logic.
	expected := []int{1, 3, 4, 5, 6, 7, 8}
	for _, want := range expected {
		k, _, ok := tr.Min()
		if !ok {
			t.Fatalf("Min: want ok=true, Tree empty too early")
		}
		if k != want {
			t.Errorf("Min: want %d, got %d", want, k)
		}
		tr.Delete(k)
	}
	if tr.Size() != 0 {
		t.Errorf("after deleting all: Size want 0, got %d", tr.Size())
	}
}

// ─────────────────────────────────────────────
// 22. Float64 keys (via cmp.Ordered)
// ─────────────────────────────────────────────

func TestFloat64Keys(t *testing.T) {
	tr := avl.New[float64, string]()
	vals := []float64{3.14, 2.71, 1.41, 0.577, 1.732, -1.0, 100.0}
	for _, v := range vals {
		tr.Insert(v, fmt.Sprintf("%.3f", v))
	}

	if tr.Size() != len(vals) {
		t.Errorf("float64 Size: want %d, got %d", len(vals), tr.Size())
	}

	k, _, ok := tr.Min()
	if !ok || k != -1.0 {
		t.Errorf("float64 Min: want -1.0, got %v (ok=%v)", k, ok)
	}

	_, ok = tr.Find(3.14)
	if !ok {
		t.Error("Find(3.14): want ok=true")
	}
}

// ─────────────────────────────────────────────
// 23. Contains vs Find consistency
// ─────────────────────────────────────────────

func TestContainsFindConsistency(t *testing.T) {
	tr := newIntTree()
	rng := rand.New(rand.NewSource(55))

	inserted := map[int]bool{}
	for range 500 {
		k := rng.Intn(200)
		tr.Insert(k, k)
		inserted[k] = true
	}

	for i := range 200 {
		inTree := inserted[i]
		c := tr.Contains(i)
		_, f := tr.Find(i)
		if c != inTree {
			t.Errorf("Contains(%d): want %v, got %v", i, inTree, c)
		}
		if f != inTree {
			t.Errorf("Find(%d) ok: want %v, got %v", i, inTree, f)
		}
	}
}

// ─────────────────────────────────────────────
// 24. String() and Print() don't panic and are non-empty for non-empty Trees
// ─────────────────────────────────────────────

func TestStringAndPrint(t *testing.T) {
	tr := newIntTree()
	insertAll(tr, []int{1, 2, 3, 4, 5})

	s := tr.String()
	if s == "" {
		t.Error("String() on non-empty Tree returned empty string")
	}
	// Print should not panic
	tr.Print()
}

// ─────────────────────────────────────────────
// 25. Insert identical key-value pairs many times
// ─────────────────────────────────────────────

func TestInsertSameKeyValueMany(t *testing.T) {
	tr := newIntTree()
	for range 1000 {
		tr.Insert(1, 1)
	}
	vals, ok := tr.Find(1)
	if !ok {
		t.Fatal("Find after 1000 identical inserts: want ok=true")
	}
	// All 1000 values should be stored (or at minimum the key is present)
	if len(vals) == 0 {
		t.Error("Find: expected at least one value")
	}
}

// ─────────────────────────────────────────────
// 26. Interleaved insert/delete/find
// ─────────────────────────────────────────────

func TestInterleavedOps(t *testing.T) {
	tr := newIntTree()
	present := map[int]bool{}

	seq := []struct {
		op  string
		key int
	}{
		{"insert", 10}, {"insert", 20}, {"insert", 5}, {"find", 10},
		{"delete", 10}, {"find", 10}, {"insert", 10}, {"insert", 15},
		{"delete", 5}, {"find", 5}, {"insert", 5}, {"delete", 20},
		{"insert", 25}, {"delete", 25}, {"insert", 30}, {"find", 30},
	}

	for _, step := range seq {
		switch step.op {
		case "insert":
			tr.Insert(step.key, step.key)
			present[step.key] = true
		case "delete":
			tr.Delete(step.key)
			delete(present, step.key)
		case "find":
			_, got := tr.Find(step.key)
			want := present[step.key]
			if got != want {
				t.Errorf("Find(%d) after %s: want %v, got %v", step.key, step.op, want, got)
			}
		}
	}
}

// ─────────────────────────────────────────────
// 27. All keys in a range are present
// ─────────────────────────────────────────────

func TestContiguousRange(t *testing.T) {
	tr := newIntTree()
	lo, hi := -250, 750
	for i := lo; i <= hi; i++ {
		tr.Insert(i, i*2)
	}
	n := hi - lo + 1
	if tr.Size() != n {
		t.Errorf("contiguous range Size: want %d, got %d", n, tr.Size())
	}
	for i := lo; i <= hi; i++ {
		if !tr.Contains(i) {
			t.Errorf("Contains(%d) in range [%d,%d]: want true", i, lo, hi)
		}
	}
}

// ─────────────────────────────────────────────
// 28. Tree remains correct after large delete wave
// ─────────────────────────────────────────────

func TestLargeDeleteWave(t *testing.T) {
	tr := newIntTree()
	n := 2000
	for i := range n {
		tr.Insert(i, i)
	}
	// Delete even numbers
	for i := 0; i < n; i += 2 {
		tr.Delete(i)
	}
	// Odd numbers must still be there
	for i := 1; i < n; i += 2 {
		if !tr.Contains(i) {
			t.Errorf("Contains(%d) after even-delete wave: want true", i)
		}
	}
	// Even numbers must be gone
	for i := 0; i < n; i += 2 {
		if tr.Contains(i) {
			t.Errorf("Contains(%d) after delete: want false", i)
		}
	}
}

// ─────────────────────────────────────────────
// 29. String representation changes after mutations
// ─────────────────────────────────────────────

func TestStringChangesAfterMutation(t *testing.T) {
	tr := newIntTree()
	tr.Insert(1, 1)
	s1 := tr.String()
	tr.Insert(2, 2)
	s2 := tr.String()
	if s1 == s2 {
		t.Log("String() did not change after insert — may be fine if implementation ignores it")
	}
	tr.Delete(1)
	s3 := tr.String()
	_ = s3 // just ensure no panic
}

// ─────────────────────────────────────────────
// 30. Benchmark: sequential insert 100k
// ─────────────────────────────────────────────

func BenchmarkInsertSequential(b *testing.B) {
	for b.Loop() {
		tr := newIntTree()
		for i := range 100_000 {
			tr.Insert(i, i)
		}
	}
}

// ─────────────────────────────────────────────
// 31. Benchmark: random insert 100k
// ─────────────────────────────────────────────

func BenchmarkInsertRandom(b *testing.B) {
	rng := rand.New(rand.NewSource(42))

	for b.Loop() {
		tr := newIntTree()
		for i := range 100_000 {
			tr.Insert(rng.Int(), i)
		}
	}
}

// ─────────────────────────────────────────────
// 32. Benchmark: Find in large Tree
// ─────────────────────────────────────────────

func BenchmarkFind(b *testing.B) {
	tr := newIntTree()
	for i := range 100_000 {
		tr.Insert(i, i)
	}

	for n := 0; b.Loop(); n++ {
		tr.Find(n % 100_000)
	}
}

// ─────────────────────────────────────────────
// 33. Verify AVL height is strictly logarithmic
// ─────────────────────────────────────────────

func TestHeightLogarithmic(t *testing.T) {
	sizes := []int{1, 10, 100, 1_000, 10_000, 100_000}
	for _, n := range sizes {
		tr := newIntTree()
		for i := range n {
			tr.Insert(i, i)
		}
		h := tr.Height()
		// Very loose upper bound: 2 * ceil(log2(n+1)) + 2
		bound := 2
		for tmp := n; tmp > 0; tmp >>= 1 {
			bound++
		}
		bound *= 2
		if h > bound {
			t.Errorf("n=%d: height %d exceeds generous bound %d", n, h, bound)
		}
	}
}

// ─────────────────────────────────────────────
// 34. Min key has correct values (multi-value)
// ─────────────────────────────────────────────

func TestMinWithMultipleValues(t *testing.T) {
	tr := newIntTree()
	// Insert two values for the minimum key
	tr.Insert(1, 100)
	tr.Insert(1, 200)
	tr.Insert(5, 500)
	tr.Insert(3, 300)

	k, vals, ok := tr.Min()
	if !ok || k != 1 {
		t.Fatalf("Min: want (1,true), got (%d,%v)", k, ok)
	}
	if len(vals) < 2 {
		t.Errorf("Min values: want at least 2, got %d: %v", len(vals), vals)
	}
	sum := 0
	for _, v := range vals {
		sum += v
	}
	if sum != 300 { // 100+200
		t.Errorf("Min value sum: want 300, got %d", sum)
	}
}

// ─────────────────────────────────────────────
// 35. Long strings as keys
// ─────────────────────────────────────────────

func TestLongStringKeys(t *testing.T) {
	tr := avl.New[string, int]()
	n := 200
	keys := make([]string, n)
	for i := range keys {
		keys[i] = strings.Repeat("x", i+1) + fmt.Sprintf("%05d", i)
		tr.Insert(keys[i], i)
	}
	if tr.Size() != n {
		t.Errorf("long-string Size: want %d, got %d", n, tr.Size())
	}
	for _, k := range keys {
		if !tr.Contains(k) {
			t.Errorf("long key %q not found", k[:10]+"...")
		}
	}
}
