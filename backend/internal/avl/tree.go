package avl

import (
	"cmp"
	"log"
)

type Tree[K cmp.Ordered, V any] struct {
	root *node[K, V]
}

// Constructor
func New[K cmp.Ordered, V any]() *Tree[K, V] {
	return &Tree[K, V]{}
}
func (t *Tree[K, V]) Insert(key K, value V) {
	t.root = t.root.insert(key, value)
}

func (t *Tree[K, V]) Delete(key K) {
	if t.root == nil {
		return
	}
	t.root = t.root.delete(key)
}
func (t *Tree[K, V]) Height() int {
	if t.root == nil {
		return 0
	}
	return t.root.getHeight()
}
func (t *Tree[K, V]) Size() int {
	if t.root == nil {
		return 0
	}
	return t.root.getSize()
}
func (t *Tree[K, V]) Min() (K, []V, bool) {
	var zeroK K
	if t.root == nil {
		return zeroK, nil, false
	}
	k, v := t.root.minValueNode().get()
	return k, v, true
}
func (t *Tree[K, V]) Find(key K) ([]V, bool) {
	if t.root == nil {
		return nil, false
	}
	return t.root.find(key)
}
func (t *Tree[K, V]) Contains(key K) bool {
	if t.root == nil {
		return false
	}
	return t.root.contains(key)
}
func (t *Tree[K, V]) String() string {
	if t.root == nil {
		return "[]"
	}
	return t.root.String()
}
func (t *Tree[K, V]) Print() {
	if t.root == nil {
		log.Println("Tried to print empty tree")
		return
	}
	t.root.Print()
}
func (t *Tree[K, V]) SubTreeFromKey(key K) (*Tree[K, V], bool) {
	if t.root == nil {
		return nil, false
	}

	node, found := t.root.findNode(key)
	if !found || node == nil {
		return nil, false
	}
	return &Tree[K, V]{root: node}, true

}
