package avl

type Node[T any] struct {
	data  T
	left  *Node[T]
	right *Node[T]
}
