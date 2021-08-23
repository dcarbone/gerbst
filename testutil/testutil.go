package testutil

import (
	"testing"

	"github.com/dcarbone/gerbst"
)

type CountableTree interface {
	Count() uint
	CountLeft() uint
	CountRight() uint
}

func BuildTestCounts(tree CountableTree, p bool, total, left, right uint) func(*testing.T) {
	return func(t *testing.T) {
		if p {
			t.Parallel()
		}
		if c := tree.Count(); c != total {
			t.Logf("Expected tree to have count %d, saw %d", total, c)
			t.Fail()
		}
		if c := tree.CountLeft(); c != left {
			t.Logf("Expected tree to have count left %d, saw %d", left, c)
			t.Fail()
		}
		if c := tree.CountRight(); c != right {
			t.Logf("Expected tree to have count right of %d, saw %d", right, c)
			t.Fail()
		}
	}
}

type DepthAwareTree interface {
	DepthMax() uint
	DepthMaxLeft() uint
	DepthMaxRight() uint
}

func BuildTestDepths(tree DepthAwareTree, p bool, max, left, right uint) func(*testing.T) {
	return func(t *testing.T) {
		if p {
			t.Parallel()
		}
		if d := tree.DepthMax(); d != max {
			t.Logf("Expected tree to have max depth %d, saw %d", max, d)
			t.Fail()
		}
		if d := tree.DepthMaxLeft(); d != left {
			t.Logf("Expected tree to have max depth left of %d, saw %d", left, d)
			t.Fail()
		}
		if d := tree.DepthMaxRight(); d != right {
			t.Logf("Expected tree to have max depth right of %d, saw %d", right, d)
			t.Fail()
		}
	}
}

type GetTest struct {
	Key    uint
	Exists bool
	Value  interface{}
}

type GetTests []GetTest

type GettableTree interface {
	Get(key uint) (*gerbst.Node, bool)
	GetRecurse(key uint) (*gerbst.Node, bool)
}

func BuildTestGets(tree GettableTree, p bool, gts GetTests) func(*testing.T) {
	return func(t *testing.T) {
		if p {
			t.Parallel()
		}
		for _, gt := range gts {
			gn, gok := tree.Get(gt.Key)
			grn, grok := tree.GetRecurse(gt.Key)

			if gok != grok {
				t.Logf("Expected Get and GetRecurse to agree for key=%d, saw Get=%t and GetRecurse=%t", gt.Key, gok, grok)
				t.Fail()
			}

			if gok != gt.Exists {
				t.Logf("Expected Get key %d ok=%t, saw %t", gt.Key, gt.Exists, gok)
				t.Fail()
			}
			if gok {
				if gn.Value() != gt.Value {
					t.Logf("Expected key %d value to be %T(%[2]v), saw %T(%[3]v)", gt.Key, gt.Value, gn.Value())
					t.Fail()
				}
			}

			if gn != grn {
				t.Logf("Expected Get and GetRecurse to return the same node, saw Get=%v and GetRecurse=%v", gn, grn)
				t.Fail()
			}
		}
	}
}

func GetTestsFromKeys(existsKeys, missingKeys []uint) GetTests {
	gts := make(GetTests, 0)
	for _, k := range existsKeys {
		gts = append(gts, GetTest{
			Key:    k,
			Exists: true,
			Value:  k,
		})
	}
	for _, k := range missingKeys {
		gts = append(gts, GetTest{
			Key:    k,
			Exists: false,
		})
	}
	return gts
}
