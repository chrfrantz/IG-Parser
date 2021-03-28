package exporter

import (
	"IG-Parser/tree"
	"fmt"
	"log"
	"strconv"
	"strings"
)

/*
Generates all permutations of a given set of input arrays, representing an
institutional statement alongside its components per entry.
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
Inputs are statements composed of associated nodes (Structure: [statement][Node references for components].
In addition, the output can be returned in terms of individual references (e.g., 1,2,3,4, etc.), or as aggregate ranges
(e.g., 1-4,6,9-13). The latter is usually useful if producing human-readable output.
If indicated, references can be incremented during aggregation (e.g., 1 is entered as 2, 2 as 3, etc.) This may be
particularly useful if converting from 0- to 1-based indices.
Returns an array of maps of nodes pointing to arrays of associated statement IDs, or if generateRanges is activated,
ranges of statement references, alongside potential incrementing of reference values.
Structure: [column ID of component]map[node reference for each value of component (e.g., Farmer, Certifier)][statement IDs
where component value apply]
 */
func GenerateLogicalOperatorLinkagePerCombination(stmts [][]*tree.Node, generateRanges bool, incrementReferences bool) []map[*tree.Node][]string {

	if len(stmts) == 0 {
		log.Fatal("Empty input - no statement permutations.")
	}

	// Number of components
	componentCount := len(stmts[0])
	// Collection of logical linkages ([column ids]map[Node ref component value][row IDs where component value applies])
	compLinks := make([]map[*tree.Node][]string,0)
	// Index of component of concern in statement
	columnIdx := 0

	for columnIdx < componentCount {

		// Generate map of category values that will hold statement references
		combinationStructure := make(map[*tree.Node][]string)

		// For each statement (with id and components)
		for id, comps := range stmts {

			// Check for each component (in that statement) - column-wise
			for compIdx, compVal := range comps {

				// When index is of interest ...
				if compIdx == columnIdx {

					// ... simply collect references per value
					nodeRefs := combinationStructure[compVal]
					// Add id to references, alongside potential range generation and incrementing
					nodeRefs = GenerateReferenceSlice(nodeRefs, id, generateRanges, incrementReferences)
					// Adding updated reference list to combination
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

// Symbol used to represent ranges (e.g., 5-7)
const rangeSeparator = "-"

/*
Add value (id) to existing slice (nodeRefs), either by simple appending, or generation of ranges (e.g., 5-7 instead of
5,6,7), and optional incrementing of references during generation (e.g., for mapping zero-based input to 1-based output)
 */
func GenerateReferenceSlice(nodeRefs []string, id int, generateRanges bool, incrementReferences bool) []string{
	if nodeRefs == nil {
		// by retrieving potentially existing references and ...
		nodeRefs = []string{}
	}

	addedId := id
	// Increment reference id if indicated
	if incrementReferences {
		addedId += 1
	}
	fmt.Println("Testing value ", addedId)
	// If generation of ranges is indicated and previous entries exist ...
	if generateRanges && len(nodeRefs) > 0 {
		// Retrieve previously added element
		val := nodeRefs[len(nodeRefs) - 1]
		if strings.Contains(val, rangeSeparator) {
			fmt.Println("Detected existing range ", val)
			// If range separator symbol is contained, extract prefix
			firstEndIndex := strings.Index(val, rangeSeparator)
			// Extract last value in range expression
			lastValue := val[firstEndIndex+1:len(val)]
			// Convert to int
			intVal, err := strconv.Atoi(lastValue)
			fmt.Println("Last value in range: ", intVal)
			// Prepare for error
			stopRangeCheck := false
			if err != nil {
				log.Println("Extraction of integer from ", val , " did not work.")
				stopRangeCheck = true
			}
			if !stopRangeCheck && intVal == (addedId - 1) {
				// Generate range structure first-last
				valueToAdd := val[:firstEndIndex] + rangeSeparator + strconv.Itoa(addedId)
				// Overwrite previous entry
				nodeRefs[len(nodeRefs)-1] = valueToAdd
				fmt.Println("Replaced previous entry with new value ", valueToAdd)
			} else {
				// Create new entry
				nodeRefs = append(nodeRefs, strconv.Itoa(addedId))
				fmt.Println("Created new entry ", addedId)
			}
		} else {
			// else simply check if previous value is decrement of to-be-added value
			stopRangeCheck := false
			intVal, err := strconv.Atoi(val)
			if err != nil {
				log.Println("Extraction of integer from ", val , " did not work.")
				stopRangeCheck = true
			}
			fmt.Println("Checking for previous non-range value ", val)
			if !stopRangeCheck && intVal == (addedId - 1) {
				// Overwrite previous entry
				nodeRefs[len(nodeRefs)-1] = val + rangeSeparator + strconv.Itoa(addedId)
				fmt.Println("Overwrite existing entry with ", (val + rangeSeparator + strconv.Itoa(addedId)))
			} else if intVal != addedId {
				// Create new entry
				nodeRefs = append(nodeRefs, strconv.Itoa(addedId))
				fmt.Println("Created new entry ", addedId)
 			} else {
 				// if old value and current are the same, simply omit if creating ranges, i.e., aggregating
 				fmt.Println("Values has not been added, since identical to last value ", intVal)
			}
		}
	} else {
		// ... adding this statement's ID
		nodeRefs = append(nodeRefs, strconv.Itoa(addedId))
	}
	return nodeRefs
}


