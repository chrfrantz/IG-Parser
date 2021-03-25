package exporter

import (
	"IG-Parser/tree"
	"fmt"
	"log"
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

		// Reset temporary node array (capturing one atomic statement each time)
		out = []*tree.Node{}

		// Adjust position references
		pos[len(nodeArrays)-1]++
	}
	return stmts
}

/*
Generates the statement IDs per component category for quick retrieval
Input is are statements composed of associated nodes. Structure: [statement ID][component nodes]
Returns an array of maps of nodes pointing to arrays of associated statement IDs.
Structure: [column ID of component]map[node reference for each value of component (e.g., Farmer, Certifier)][statement IDs
where component value apply]
 */
func GenerateLogicalOperatorLinkagePerCombination(stmts [][]*tree.Node) []map[*tree.Node][]int {

	if len(stmts) == 0 {
		log.Fatal("Empty input - no statement permutations.")
	}

	// Number of components
	componentCount := len(stmts[0])
	// Collection of logical linkages ([column ids]map[Node ref component value][row IDs where component value applies])
	compLinks := make([]map[*tree.Node][]int,0)
	// Index of component of concern in statement
	columnIdx := 0

	for columnIdx < componentCount {

		// Generate map of category values that will hold statement references
		combinationStructure := make(map[*tree.Node][]int)

		// For each statement (with id and components)
		for id, comps := range stmts {

			// Check for each component (in that statement) - column-wise
			for compIdx, compVal := range comps {

				// when index is of interest
				if compIdx == columnIdx {

					// Simply collect references per value
					nodeRefs := combinationStructure[compVal]
					if nodeRefs == nil {
						nodeRefs = []int{}
						// Go through the entire
					}

					// Add this statement's ID
					nodeRefs = append(nodeRefs, id)
					combinationStructure[compVal] = nodeRefs
				}

			}

		}
		// Now save statement references for all values of a given component to shared array
		compLinks = append(compLinks, combinationStructure)
		// Move to next column (component ID)
		columnIdx++
	}
	return compLinks
}



