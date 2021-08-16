package gerbst_test

import (
	"testing"

	"github.com/dcarbone/gerbst"
)

func TestDoesItWorkAtAll(t *testing.T) {
	const expectedTree = `ROOT[12(12)]
└── LEFT[11(11)]
│   ├── LEFT[7(7)]
│       └── RIGHT[9(9)]
└── RIGHT[90(90)]
    └── LEFT[82(82)]
`

	input := []uint{12, 11, 90, 82, 7, 9}
	n := gerbst.NewWithKeys(input)

	if st := n.StringTree(); st != expectedTree {
		t.Log("Tree did not match expected")
		t.Logf("Expected:\n%s", expectedTree)
		t.Logf("Actual:\n%s", st)
		t.Fail()
	}

	deepest := n.DeepestNode()

	if v, ok := deepest.Value().(uint); !ok {
		t.Logf("Expected deepest value to be %d, saw %v (%T)", 9, v, v)
		t.Fail()
	}
	if d := deepest.Depth(); d != 3 {
		t.Logf("Expected deepest depth to be 3, saw %d", d)
		t.Fail()
	}

	var node11 *gerbst.Node

	searchFN := func(n *gerbst.Node) bool {
		if n.Value().(uint) == 11 {
			node11 = n
			return false
		}
		return true
	}

	n.SearchFunc(searchFN)

	if node11 == nil {
		t.Log("Unable to locate node with value of 11")
		t.Fail()
	}

	if v := n.SmallestKey(); v != 7 {
		t.Logf("Expected SmallestKey to return %d, saw %d", 7, v)
		t.Fail()
	}

	n.Put(7, 1)

	if n1, ok := n.Get(7); !ok {
		t.Logf("Unable to locate node with key %d", 7)
		t.Fail()
	} else if v := n1.Value(); v != 1 {
		t.Logf("Expected to find node key 7 with updated value of 1, saw %v (%T)", v, v)
		t.Fail()
	}
}
