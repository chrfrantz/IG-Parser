package exporter

import (
	"IG-Parser/tree"
	"fmt"
	"strconv"
)

/*
Generates all permutations of a given set of input arrays, representing an
institutional statement alongside its components per entry
Output array structure is [statement][component instances of statement]
 */
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
				//fmt.Println("Append " + ar[p].String())
				out = append(out, ar[p])
			}
		}
		fmt.Println("Wrote atomic statement ", ct)

		// Assign generated statement to return data structure
		stmts[ct] = out
		ct++

		//fmt.Println("Full combo: ", out)

		// Reset temporary node array (capturing one atomic statement each time)
		out = []*tree.Node{}

		// Adjust position references
		pos[len(nodeArrays)-1]++
	}
	return stmts
}

func GenerateLogicalOperatorLinkagePerCombination(stmts [][]*tree.Node) map[*tree.Node]map[string][]int {

	// Structure: Node instance (comparison source), logical operator
	//compLinks := make(map[*tree.Node]map[string][]int)

	//compLinks := []LogicalOperatorLinkage{}


	searchIdx := 0 // Index of component of concern in statement
	//for com
	compLink := LogicalOperatorLinkage{Index: searchIdx}
	// For each statement (with id and components)
	for id, comps := range stmts {

		// Check for each component
		for compIdx, compVal := range comps {
			// when index is of interest
			if compIdx == searchIdx {
				// Assign component reference upon first iteration
				if compLink.Component == nil {
					// Link source component
					compLink.Component = compVal
				// Search for target component in same column (i.e., same index in any combination)
				} else if compVal != compLink.Component {

					if compLink.LinkedStatements[compVal] == nil {
						fmt.Println("Found different component value ", compVal)
						res, ops := tree.FindLogicalLinkage(compLink.Component, compVal, nil)
						fmt.Println("Return value: ", res)
						if res {
							fmt.Println("Found logical linkage between ", compLink.Component, " and ", compVal, ": ", ops, " in statement, ", id)

							// Create new collection for this particular combination
							newLink := []int{id}
							linked := compLink.LinkedStatements
							if linked == nil {
								linked = make(map[*tree.Node][]int)
							}
							// Add entry for target component
							linked[compVal] = newLink
							// Store into main structure for source component
							compLink.LinkedStatements = linked

							// Store logical operator
							operatorLink := compLink.LinkedComponentOperator
							// If operator links are empty
							if operatorLink == nil {
								// Create new link
								operatorLink = make(map[*tree.Node][]string)
							}
							// Save new elements
							operatorLink[compVal] = ops
							// Store into main structure
							compLink.LinkedComponentOperator = operatorLink
						}
					} else {
						// Add ID of statement to reference list
						arr := compLink.LinkedStatements[compVal]
						arr = append(arr, id)
						compLink.LinkedStatements[compVal] = arr
						fmt.Println("Component value ", compVal, " found previously. Stored statement id ", id)
					}


				}


			}
		}
	}

	fmt.Println(compLink)

	return nil
}



