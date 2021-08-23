package gerbst_test

import (
	"testing"

	"github.com/dcarbone/gerbst"
	"github.com/dcarbone/gerbst/testutil"
)

func TestLockingTree(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		getTests := testutil.GetTests{
			{
				Key:    0,
				Exists: false,
			},
			{
				Key:    1,
				Exists: false,
			},
		}

		lt := gerbst.NewLockingTree()

		t.Run("counts", testutil.BuildTestCounts(lt, true, 0, 0, 0))
		t.Run("depths", testutil.BuildTestDepths(lt, true, 0, 0, 0))
		t.Run("gets", testutil.BuildTestGets(lt, true, getTests))
	})

	t.Run("new_keys", func(t *testing.T) {
		keys := []uint{12, 11, 90, 82, 7, 9}
		getTests := testutil.GetTestsFromKeys(keys, []uint{0, 83, 100, 55})

		lt := gerbst.NewLockingTreeWithKeys(keys)

		t.Run("counts", testutil.BuildTestCounts(lt, true, 6, 3, 2))
		t.Run("depths", testutil.BuildTestDepths(lt, true, 4, 4, 3))
		t.Run("gets", testutil.BuildTestGets(lt, true, getTests))
	})
}

func TestDoesItWorkAtAll(t *testing.T) {
	const expectedTree = `ROOT[12(12)]
└── LEFT[11(11)]
│   ├── LEFT[7(7)]
│       └── RIGHT[9(9)]
└── RIGHT[90(90)]
    └── LEFT[82(82)]
`

	input := []uint{12, 11, 90, 82, 7, 9}
	n := gerbst.NewLockingTreeWithKeys(input)

	if st := n.StringTree(); st != expectedTree {
		t.Log("Tree did not match expected")
		t.Logf("Expected:\n%s", expectedTree)
		t.Logf("Actual:\n%s", st)
		t.Fail()
	}

	if v := n.LowestKey(); v != 7 {
		t.Logf("Expected LowestKey to return %d, saw %d", 7, v)
		t.Fail()
	}

	n.PutRecurse(7, 1)

	if n1, ok := n.GetRecurse(7); !ok {
		t.Logf("Unable to locate node with key %d", 7)
		t.Fail()
	} else if v := n1.Value(); v != 1 {
		t.Logf("Expected to find node key 7 with updated value of 1, saw %v (%T)", v, v)
		t.Fail()
	}
}
