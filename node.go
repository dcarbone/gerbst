package gerbst

import (
	"fmt"

	"github.com/disiqueira/gotree"
)

// Node represents the exportable representation of a given node within a tree
type Node struct {
	key   uint
	value interface{}
	depth uint
	side  NodeSide
}

// newNode constructs the actual node instance
func newNode(key uint, value interface{}, depth uint, side NodeSide) *Node {
	n := new(Node)
	n.key = key
	n.value = value
	n.depth = depth
	n.side = side
	return n
}

// Key returns this node's key
func (n *Node) Key() uint {
	return n.key
}

// Value returns this node's value
func (n *Node) Value() interface{} {
	return n.value
}

// Depth returns the depth of the current node from root
func (n *Node) Depth() uint {
	return n.depth
}

// Side returns the position of this node relative to its parent, or ROOT if it is the root node.
func (n *Node) Side() NodeSide {
	return n.side
}

type treeNode struct {
	*Node

	parent *treeNode
	left   *treeNode
	right  *treeNode

	loKey uint
	hiKey uint

	count      uint // count is 1 (self) + countLeft + countRight
	countLeft  uint
	countRight uint

	depthMax      uint
	depthMaxLeft  uint
	depthMaxRight uint
}

func newTreeNode(key uint, value interface{}, depth uint, side NodeSide, parent, left, right *treeNode) *treeNode {
	tn := new(treeNode)
	tn.Node = newNode(key, value, depth, side)

	// set nodes
	tn.parent = parent
	tn.left = left
	tn.right = right

	// set base meta values
	tn.count = 1
	tn.depthMax = tn.depth
	tn.loKey = tn.key
	tn.hiKey = tn.key

	return tn
}

// Left returns the left branch of this tree, if there is one
func (tn *treeNode) Left() *treeNode {
	return tn.left
}

// Right returns the right branch of this tree, if there is one
func (tn *treeNode) Right() *treeNode {
	return tn.right
}

func (tn *treeNode) Get(key uint) (*Node, bool) {
	n := tn

	// execute walk
	for n != nil {
		if n.key == key {
			break
		} else if n.key > key && n.left != nil {
			n = n.left
		} else if n.key < key && n.right != nil {
			n = n.right
		} else {
			n = nil
			break
		}
	}

	// handle response
	if n == nil {
		return nil, false
	}

	return n.Node, true
}

func (tn *treeNode) GetRecurse(key uint) (*Node, bool) {
	if tn.key == key {
		return tn.Node, true
	} else if tn.key > key && tn.left != nil {
		if ln, ok := tn.left.GetRecurse(key); ok {
			return ln, ok
		}
	} else if tn.key < key && tn.right != nil {
		if rn, ok := tn.right.GetRecurse(key); ok {
			return rn, ok
		}
	}
	return nil, false
}

func (tn *treeNode) Put(key uint, value interface{}) {
	n := tn
	for n != nil {
		// if we need to update the existing node
		if n.key == key {
			n.Node = newNode(key, value, tn.depth, tn.side)
			return
		} else if n.key > key {
			if n.left == nil {
				// if we get here, key is lower than local and we have no left node, so create one
				// and move on.
				n.left = newTreeNode(key, value, n.depth+1, NodeSideLeft, n, nil, nil)
				updateMeta(n.left)
				return
			} else {
				// set parent to local and update local to left side of local
				n = n.left
			}
		} else if n.right == nil {
			// if we get here, key is higher than local and we have no right node, so create one
			// and move on.
			n.right = newTreeNode(key, value, n.depth+1, NodeSideRight, n, nil, nil)
			updateMeta(n.right)
			return
		} else {
			// update parent to n and update local to right side of local
			n = n.right
		}
	}
}

func (tn *treeNode) PutRecurse(key uint, value interface{}) {
	if tn.key == key {
		tn.Node = newNode(key, value, tn.depth, tn.side)
	} else if tn.key > key {
		if tn.left == nil {
			tn.left = newTreeNode(key, value, tn.depth+1, NodeSideLeft, tn, nil, nil)
			updateMeta(tn.left)
		} else {
			tn.left.PutRecurse(key, value)
		}
	} else if tn.right == nil {
		tn.right = newTreeNode(key, value, tn.depth+1, NodeSideRight, tn, nil, nil)
		updateMeta(tn.right)
	} else {
		tn.right.PutRecurse(key, value)
	}
}

func (tn *treeNode) metaString() string {
	return fmt.Sprintf(
		"node=%p; parent=%p; side=%q, count=%d; countLeft=%d; countRight=%d; depth=%d; depthMax=%d; depthMaxLeft=%d; depthMaxRight=%d",
		tn,
		tn.parent,
		tn.side,
		tn.count,
		tn.countLeft,
		tn.countRight,
		tn.depth,
		tn.depthMax,
		tn.depthMaxLeft,
		tn.depthMaxRight)
}

// String returns a printable sum of this node in the format of SIDE[KEY(VALUE)]
func (tn *treeNode) String() string {
	return fmt.Sprintf("%s[%d(%v)]", tn.side, tn.key, tn.value)
}

// buildTreePrinter recursively builds our tree printer for us.  This was included so I can be lazy and not
// write my own visual inspector
func (tn *treeNode) buildTreePrinter() gotree.Tree {
	// construct new tree
	root := gotree.New(tn.String())

	// add left branch
	if tn.left != nil {
		root.AddTree(tn.left.buildTreePrinter())
	}

	// add right branch
	if tn.right != nil {
		root.AddTree(tn.right.buildTreePrinter())
	}

	// we did it.
	return root
}

func updateMeta(src *treeNode) {
	srcDepth := src.depth
	srcKey := src.key

	local := src
	parent := src.parent

	for parent != nil {
		// increment overall count
		parent.count++

		// side-specific logic
		switch local.side {
		case NodeSideLeft:
			parent.countLeft++
			if parent.depthMaxLeft < srcDepth {
				parent.depthMaxLeft = srcDepth
			}
		case NodeSideRight:
			parent.countRight++
			if parent.depthMaxRight < srcDepth {
				parent.depthMaxRight = srcDepth
			}
		}

		// update parent max depth
		if parent.depthMax < srcDepth {
			parent.depthMax = srcDepth
		}

		// update parent high or low key
		if parent.loKey > srcKey {
			parent.loKey = srcKey
		} else if parent.hiKey < srcKey {
			parent.hiKey = srcKey
		}

		// continue loop
		local = parent
		parent = parent.parent
	}
}
