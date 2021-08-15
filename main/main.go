package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/dcarbone/gerbst"
)

func main() {

	rand.Seed(time.Now().UnixNano())

	input := []uint{12, 11, 90, 82, 7, 9}

	//input := []uint{55, 123, 53134, 56, 33, 11}
	//for i := 0; i < 20; i++ {
	//	input = append(input, uint(rand.Uint32()))
	//}

	n := gerbst.NewWithKeys(input)

	fmt.Println(n.StringTree())

	deepest := n.DeepestNode()

	var node11 *gerbst.Node

	searchFN := func(n *gerbst.Node) bool {
		if n.Value().(uint) == 11 {
			node11 = n
			return false
		}
		return true
	}

	n.SearchFunc(searchFN)

	fmt.Println(deepest, deepest.Depth(), node11)
}
