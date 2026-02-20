package avl

import (
	"cmp"
	"fmt"
	"log"
	"strings"
)

type node[K cmp.Ordered, V any] struct {
	key    K
	value  []V
	left   *node[K, V]
	right  *node[K, V]
	height int
}

// constructor
func createNode[K cmp.Ordered, V any](key K, value V) *node[K, V] {
	return &node[K, V]{key: key, value: []V{value}, height: 1}
}

func (n *node[K, V]) insert(key K, value V) *node[K, V] {
	if n == nil {
		return createNode(key, value)
	}
	if key < n.key {
		n.left = n.left.insert(key, value)
	} else if key > n.key {
		n.right = n.right.insert(key, value)
	} else {
		n.value = append(n.value, value)
		return n
	}

	return n.balance()
}

func (n *node[K, V]) balance() *node[K, V] {
	n.updateHeight()
	balance := n.balanceFactor()

	if balance > 1 && n.left.balanceFactor() >= 0 {
		//left left
		return n.rotateRight()
	}

	if balance > 1 && n.left.balanceFactor() < 0 {
		//left right
		n.left = n.left.rotateLeft()
		return n.rotateRight()

	}
	if balance < -1 && n.right.balanceFactor() <= 0 {
		//right right
		return n.rotateLeft()
	}
	if balance < -1 && n.right.balanceFactor() > 0 {
		//right left
		n.right = n.right.rotateRight()
		return n.rotateLeft()
	}
	return n
}

func (n *node[K, V]) delete(key K) *node[K, V] {
	if n == nil {
		return n
	}
	if key < n.key {
		n.left = n.left.delete(key)
	} else if key > n.key {
		n.right = n.right.delete(key)
	} else {
		if n.left == nil {
			temp := n.right
			n = nil
			return temp
		} else if n.right == nil {
			temp := n.left
			n = nil
			return temp
		}

		temp := n.right.minValueNode()
		n.key = temp.key
		n.value = temp.value
		n.right = n.right.delete(temp.key)
	}
	return n.balance()
}
func (n *node[K, V]) minValueNode() *node[K, V] {
	current := n
	for current.left != nil {
		current = current.left
	}
	return current
}
func (n *node[K, V]) get() (K, []V) {
	return n.key, n.value
}

func (n *node[K, V]) Print() {
	if n == nil {
		return
	}
	log.Println(n.String())
}
func (n *node[K, V]) String() string {
	if n == nil {
		return ""
	}
	sb := strings.Builder{}
	sb.WriteString("[")
	n.traverseString(&sb)
	str := sb.String()
	str = strings.TrimSuffix(str, ", ")
	return str + "]"
}

func (n *node[K, V]) traverseString(sb *strings.Builder) {
	if n == nil {
		return
	}
	n.left.traverseString(sb)
	fmt.Fprintf(sb, "%v, ", n.key)
	n.right.traverseString(sb)

}
func (n *node[K, V]) find(key K) ([]V, bool) {
	if n == nil {
		return nil, false
	}
	if n.key == key {
		return n.value, true
	}
	if key < n.key {
		return n.left.find(key)
	}
	return n.right.find(key)

}
func (n *node[K, V]) findNode(key K) (*node[K, V], bool) {
	if n == nil {
		return nil, false
	}
	if n.key == key {
		return n, true
	}
	if key < n.key {
		return n.left.findNode(key)
	}
	return n.right.findNode(key)

}

func (n *node[K, V]) contains(key K) bool {
	_, found := n.find(key)
	return found
}

//rotations

func (y *node[K, V]) rotateRight() *node[K, V] {
	x := y.left
	T2 := x.right
	x.right = y
	y.left = T2
	y.updateHeight()
	x.updateHeight()
	return x
}
func (x *node[K, V]) rotateLeft() *node[K, V] {
	y := x.right
	T2 := y.left
	y.left = x
	x.right = T2
	x.updateHeight()
	y.updateHeight()
	return y

}

//helper methods

func (n *node[K, V]) getHeight() int {
	if n == nil {
		return 0
	}
	return n.height
}
func (n *node[K, V]) getSize() int {
	if n == nil {
		return 0
	}
	return 1 + n.left.getSize() + n.right.getSize()
}
func (n *node[K, V]) updateHeight() {
	if n == nil {
		log.Fatal("tried to update height on nil node")
	}
	n.height = 1 + max(n.left.getHeight(), n.right.getHeight())
}
func (n *node[K, V]) balanceFactor() int {
	if n == nil {
		return 0
	}
	return n.left.getHeight() - n.right.getHeight()
}
