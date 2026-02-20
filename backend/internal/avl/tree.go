package avl

import (
	"cmp"
)

type tree[K cmp.Ordered, V any] struct {
	root *node[K, V]
}

// Constructor
func New[K cmp.Ordered, V any](key K, value V) *tree[K, V] {
	return &tree[K, V]{root: createNode(key, value)}
}
func (t *tree[K, V]) Insert(key K, value V) {
	t.root = t.root.insert(key, value)
}

func (t *tree[K, V]) Get() (K, []V) {
	return t.root.get()
}
func (t *tree[K, V]) Key() K {
	key, _ := t.root.get()
	return key
}
func (t *tree[K, V]) Value() []V {
	_, value := t.root.get()
	return value
}
func (t *tree[K, V]) Delete(key K) {
	t.root = t.root.delete(key)
}
func (t *tree[K, V]) Height() int {
	return t.root.getHeight()
}
func (t *tree[K, V]) Size() int {
	return t.root.getSize()
}
