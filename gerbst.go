package gerbst

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/disiqueira/gotree"
)

// NodeSide represents the position of the node relatives to its parent
type NodeSide uint

const (
	NodeSideRoot NodeSide = iota + 1
	NodeSideLeft
	NodeSideRight
)

// String returns a printable representation of this node's location relative to its parent
func (ns NodeSide) String() string {
	switch ns {
	case NodeSideRoot:
		return "ROOT"
	case NodeSideLeft:
		return "LEFT"
	case NodeSideRight:
		return "RIGHT"

	default:
		return "UNKNOWN"
	}
}

// IsRoot will return true if the node that returned this has no parents
func (ns NodeSide) IsRoot() bool {
	return ns == NodeSideRoot
}

// IsLeft will return true if the node that returned this is on the left side of its immediate parent
func (ns NodeSide) IsLeft() bool {
	return ns == NodeSideLeft
}

// IsRight will return true if the node that returned this is on the right side of its immediate parent
func (ns NodeSide) IsRight() bool {
	return ns == NodeSideRight
}

// NodeSearchFunc is used in conjunction with Node.SearchFunc to recurse through all nodes present in the tree, halting
// when "false" is returned for "continue_"
type NodeSearchFunc = func(node *Node) (continue_ bool)

// Node represents a singular position at any point within the tree.
type Node struct {
	mu    sync.Mutex
	key   uint
	value interface{}
	depth uint
	side  NodeSide
	left  *Node
	right *Node
}

// New constructs a new root node.  Value is optional, if left blank will be set to value of key.
func New(key uint, value interface{}) *Node {
	var v interface{}
	if value == nil {
		v = key
	} else {
		v = value
	}
	return newNode(key, v, 0, NodeSideRoot)
}

// NewWithKeys populates the tree using a list of keys.  The value of each node will be that of the key of that node.
func NewWithKeys(keys []uint) *Node {
	if len(keys) == 0 {
		return New(0, nil)
	}
	n := New(keys[0], nil)
	for _, k := range keys[1:] {
		n.put(k, k)
	}
	return n
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

// Left returns the left branch of this tree, if there is one
func (n *Node) Left() *Node {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.left
}

// Right returns the right branch of this tree, if there is one
func (n *Node) Right() *Node {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.right
}

// SearchFunc will recurse through both branches calling the provided function at each node and each  child node.
// From the root, each immediate child branch is called in a separate go routine.
//
// To halt recursion, return false from your provided func.
//
// This method acquires a lock on the internal mutex of each node.  This means that while you may call exported methods
// that don't require a lock (Key, Value, Depth, Side, and String, as examples), exported methods that do acquire a lock
// may result in deadlock.  Be mindful of this in your function definition.
func (n *Node) SearchFunc(fn NodeSearchFunc) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if !fn(n) {
		return
	}

	stop := new(uint32)
	*stop = 0
	wg := new(sync.WaitGroup)

	// if we have a left branch, recurse through it
	if n.left != nil {
		wg.Add(1)
		go func() {
			n.left.searchFunc(fn, stop)
			wg.Done()
		}()
	}

	// if we have a right branch, recurse through it
	if n.right != nil {
		wg.Add(1)
		go func() {
			n.right.searchFunc(fn, stop)
			wg.Done()
		}()
	}

	// wait until left and right branch recursion has finished
	wg.Wait()
}

func (n *Node) searchFunc(fn NodeSearchFunc, stop *uint32) {
	// immediately test before acquiring lock and spinning up defer
	if 1 == atomic.LoadUint32(stop) {
		return
	}

	// acquire lock
	n.mu.Lock()
	defer n.mu.Unlock()

	// test ourselves
	if !fn(n) {
		// if recursion is halted, update stop pointer and return
		atomic.StoreUint32(stop, 1)
		return
	}

	// search through the left branch
	if n.left != nil {
		n.left.searchFunc(fn, stop)
	}

	// search through the right branch
	if n.right != nil {
		n.right.searchFunc(fn, stop)
	}
}

// DeepestNode returns the leafiest node there is.  It searches both branches at once in a separate routine.  This could
// eventually be improved in a few possible ways:
// 1. allowing the caller to set the maximum depth routines may be spun up
// 2. keeping track of remaining depth per branch per node and using routines up until the remaining depth is such that
// a routine would degrade, rather than improve, performance.
//
// Using routines offers no benefit with smaller trees, so a more sophisticated implementation would be able to provide
// a better balance of performance between nodes with a low amount of remaining depth vs nodes with a high amount.
func (n *Node) DeepestNode() *Node {
	n.mu.Lock()
	defer n.mu.Unlock()

	// if we have no branches, return ourselves
	if n.left == nil && n.right == nil {
		return n
	}

	// used in below routines
	var (
		ln *Node
		rn *Node

		wg = new(sync.WaitGroup)
	)

	// look at left node subtree in separate routine
	if n.left != nil {
		wg.Add(1)
		go func() {
			ln = n.left.deepestNode()
			wg.Done()
		}()
	}

	// look at right node subtree in separate routine
	if n.right != nil {
		wg.Add(1)
		go func() {
			rn = n.right.deepestNode()
			wg.Done()
		}()
	}

	// wait for subtree recursion to complete
	wg.Wait()

	// determine leafiest

	if rn == nil {
		return ln
	} else if ln == nil {
		return rn
	} else if ln.depth > rn.depth {
		return ln
	} else {
		return rn
	}
}

// deepestNode does, in effect, what DeepestNode does but without spinning up separate routines.  It is for internal
// use only.
func (n *Node) deepestNode() *Node {
	n.mu.Lock()
	defer n.mu.Unlock()

	// if we have no branches, return ourselves
	if n.left == nil && n.right == nil {
		return n
	}

	// will eventually contain the leafiest node from each branch
	var (
		ln *Node
		rn *Node
	)

	// do we have a left branch?
	if n.left != nil {
		ln = n.left.deepestNode()
	}
	// do we have a right branch?
	if n.right != nil {
		rn = n.right.deepestNode()
	}

	// return leafiest.
	if rn == nil {
		return ln
	} else if ln == nil {
		return rn
	} else if ln.depth > rn.depth {
		return ln
	} else {
		return rn
	}
}

// Put inserts a new key into the tree, if it doesn't already exist
func (n *Node) Put(key uint, value interface{}) {
	var v interface{}
	if value == nil {
		v = key
	}
	n.put(key, v)
}

// put performs the actual act of creating a new Node, if necessary.  It is separated out to prevent repeatedly testing
// value for nil-ness
func (n *Node) put(key uint, value interface{}) {
	n.mu.Lock()
	defer n.mu.Unlock()

	if n.key > key {
		if n.left == nil {
			n.left = newNode(key, value, n.depth+1, NodeSideLeft)
		} else {
			n.left.put(key, value)
		}
	} else if n.key < key {
		if n.right == nil {
			n.right = newNode(key, value, n.depth+1, NodeSideRight)
		} else {
			n.right.put(key, value)
		}
	}
}

// String returns a printable sum of this node in the format of SIDE[KEY(VALUE)]
func (n *Node) String() string {
	return fmt.Sprintf("%s[%d(%v)]", n.side, n.key, n.value)
}

// StringTree returns a string representation of the tree meant for printing
func (n *Node) StringTree() string {
	tree := n.buildTreePrinter()
	return tree.Print()
}

// buildTreePrinter recursively builds our tree printer for us.  This was included so I can be lazy and not
// write my own visual inspector
func (n *Node) buildTreePrinter() gotree.Tree {
	n.mu.Lock()
	defer n.mu.Unlock()

	// construct new tree
	root := gotree.New(n.String())

	// add left branch
	if n.left != nil {
		root.AddTree(n.left.buildTreePrinter())
	}

	// add right branch
	if n.right != nil {
		root.AddTree(n.right.buildTreePrinter())
	}

	// we did it.
	return root
}
