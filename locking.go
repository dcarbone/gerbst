package gerbst

import (
	"sync"
)

// LockingNodeSearchFunc is used in conjunction with LockingTree.SearchFunc to recurse through all nodes present in the
// tree, halting when "false" is returned for "continue_"
type LockingNodeSearchFunc = func(node *LockingTree) (continue_ bool)

// LockingTree represents a singular position at any point within the tree.
type LockingTree struct {
	mu sync.RWMutex

	root *treeNode
}

// NewLockingTree constructs a new root node.  Value is optional, if left blank will be set to value of key.
func NewLockingTree() *LockingTree {
	lt := new(LockingTree)
	return lt
}

// NewLockingTreeWithKeys populates the tree using a list of keys.  The value of each node will be that of the key of
// that node.
func NewLockingTreeWithKeys(keys []uint) *LockingTree {
	lt := NewLockingTree()
	for _, k := range keys {
		lt.Put(k, k)
	}
	return lt
}

// Count returns the total number of nodes within this tree
func (n *LockingTree) Count() uint {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.root == nil {
		return 0
	}
	return n.root.count
}

// CountLeft returns the total number of nodes on the left side of this tree
func (n *LockingTree) CountLeft() uint {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.root == nil {
		return 0
	}
	return n.root.countLeft
}

// CountRight returns the total number of nodes on the right side of this tree
func (n *LockingTree) CountRight() uint {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.root == nil {
		return 0
	}
	return n.root.countRight
}

// LowestKey returns the smallest key in this node's subtree
func (n *LockingTree) LowestKey() uint {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.root == nil {
		return 0
	}
	return n.root.loKey
}

// HighestKey returns the highest key in this node's subtree
func (n *LockingTree) HighestKey() uint {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.root == nil {
		return 0
	}
	return n.root.hiKey
}

// DepthMax returns the absolute deepest a branch goes
func (n *LockingTree) DepthMax() uint {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.root == nil {
		return 0
	}
	return n.root.depthMax
}

// DepthMaxLeft returns the maximum depth of the left branch
func (n *LockingTree) DepthMaxLeft() uint {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.root == nil {
		return 0
	}
	return n.root.depthMaxLeft
}

// DepthMaxRight returns the maximum depth of the right branch
func (n *LockingTree) DepthMaxRight() uint {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.root == nil {
		return 0
	}
	return n.root.depthMaxRight
}

// Get attempts to retrieve a node by value
func (n *LockingTree) Get(key uint) (*Node, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	// fast fail if this tree is empty or if the requested key is beyond our bounds
	if n.root == nil || key < n.root.loKey || key > n.root.hiKey {
		return nil, false
	}
	return n.root.Get(key)
}

// GetRecurse attempts to retrieve a node by key using recursion
func (n *LockingTree) GetRecurse(key uint) (*Node, bool) {
	n.mu.RLock()
	defer n.mu.RUnlock()
	// fast fail if this tree is empty or if the requested key is beyond our bounds
	if n.root == nil || key < n.root.loKey || key > n.root.hiKey {
		return nil, false
	}
	return n.root.GetRecurse(key)
}

// Put inserts a new node or updates the value of an existing node
func (n *LockingTree) Put(key uint, value interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.put(key, value, false)
}

// PutRecurse inserts a new node or updates the value of an existing node using recursion
func (n *LockingTree) PutRecurse(key uint, value interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.put(key, value, true)
}

func (n *LockingTree) put(key uint, value interface{}, recurse bool) {
	if n.root == nil {
		n.root = newTreeNode(key, value, 1, NodeSideRoot, nil, nil, nil)
		return
	}
	if recurse {
		n.root.PutRecurse(key, value)
	} else {
		n.root.Put(key, value)
	}
}

// StringTree returns a string representation of the tree meant for printing
func (n *LockingTree) StringTree() string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	if n.root == nil {
		return ""
	}
	tree := n.root.buildTreePrinter()
	return tree.Print()
}
