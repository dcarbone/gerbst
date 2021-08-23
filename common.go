package gerbst

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
