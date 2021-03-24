package exporter

import (
	"IG-Parser/tree"
	"fmt"
	"strconv"
)

func GenerateNodeArrayPermutations(nodeArrays ...[]*tree.Node) (stmts [][]*tree.Node) {

	n := 1
	// Determine output parameters
	// Iterate over number of input arrays (i.e., components)
	for _, array := range nodeArrays {
		fmt.Println("Entries in array: " + strconv.Itoa(len(array)))
		// Tolerate empty nodeArrays
		if len(array) != 0 {
			n *= len(array)
		}
	}
	fmt.Println("Number of anticipated atomic statements: " + strconv.Itoa(n))

	// Generate arrays of node arrays representing the different atomic statements
	stmts = make([][]*tree.Node, n)

	// Array of position references
	pos := make([]int, len(nodeArrays))
	// Statement counter
	ct := 0
	// Basic atomic statement structure
	out := []*tree.Node{}

	// Named loop
	loop:
	for {
		// Shift in-array position counters
		for i := len(nodeArrays) - 1; i >= 0; i-- {
			if pos[i] > 0 && pos[i] >= len(nodeArrays[i]) {
				if i == 0 || (i == 1 && pos[i-1] == len(nodeArrays[0])-1) {
					break loop
				}
				pos[i] = 0
				pos[i-1]++
			}
		}
		// Construct atomic statement
		for i, ar := range nodeArrays {
			var p = pos[i]
			if p >= 0 && p < len(ar) {
				fmt.Println("Append " + ar[p].String())
				out = append(out, ar[p])
			}
		}
		fmt.Println("Wrote atomic statement ", ct)

		// Assign generated statement to return data structure
		stmts[ct] = out
		ct++

		fmt.Println("Full combo: ", out)

		// Reset temporary node array (capturing one atomic statement each time)
		out = []*tree.Node{}

		// Adjust position references
		pos[len(nodeArrays)-1]++
	}
	return stmts
}

