package avl

import (
	"cmp"
	"log"
)

type node[T cmp.Ordered] struct {
	data   T
	left   *node[T]
	right  *node[T]
	height int
}

// constructor
func New[T cmp.Ordered](data T) *node[T] {
	return &node[T]{data: data, height: 1}
}

func (n *node[T]) Insert(pData T) *node[T] {
	if n == nil {
		return New(pData)
	}
	if pData < n.data {
		n.left = n.left.Insert(pData)
	} else if pData > n.data {
		n.right = n.right.Insert(pData)
	} else {
		log.Fatal("unhandled case for duplicate values")
	}

	n.updateHeight()
	balance := n.balanceFactor()

	if balance > 1 && n.left.balanceFactor() >= 0 {
		//left left
		return n.rotateRight()
	}

	if balance > 1 && n.left.balanceFactor() < 0 {
		//rotate left right
		n.left = n.left.rotateLeft()
		return n.rotateRight()

	}
	if balance < -1 && n.right.balanceFactor() <= 0 {
		//rotate right right
		return n.rotateLeft()
	}
	if balance < -1 && n.right.balanceFactor() > 0 {
		//rotate right left
		n.right = n.right.rotateLeft()
		return n.rotateLeft()
	}

	return n
}

//rotations

func (y *node[T]) rotateRight() *node[T] {
	x := y.left
	T2 := x.right
	x.right = y
	y.left = T2
	y.updateHeight()
	x.updateHeight()
	return x
}
func (x *node[T]) rotateLeft() *node[T] {
	y := x.right
	T2 := y.left
	y.left = x
	x.right = T2
	x.updateHeight()
	y.updateHeight()
	return x

}

//helper methods

func getHeight[T cmp.Ordered](n *node[T]) int {
	if n == nil {
		return 0
	}
	return n.height
}
func (n *node[T]) updateHeight() {
	n.height = 1 + max(getHeight(n.left), getHeight(n.right))
}
func (n *node[T]) balanceFactor() int {
	return getHeight(n.left) - getHeight(n.right)
}
