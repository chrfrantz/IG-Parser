package tree

import (
	"strings"
)

// Entry delimiters
const TREE_PRINTER_OPEN_BRACE = "{"
const TREE_PRINTER_CLOSE_BRACE = "}"
const TREE_PRINTER_LINEBREAK = "\n"
const TREE_PRINTER_SEPARATOR = "," + TREE_PRINTER_LINEBREAK

// Key name
const TREE_PRINTER_KEY_NAME = "\"name\""
const TREE_PRINTER_KEY_COMPONENT = "\"comp\""
const TREE_PRINTER_KEY_POSITION = "\"pos\""
const TREE_PRINTER_KEY_CHILDREN = "\"children\""
const TREE_PRINTER_KEY_PROPERTIES = "\"properties\""
const TREE_PRINTER_KEY_ANNOTATIONS = "\"annotations\""

// Separator
const TREE_PRINTER_EQUALS = ": "

// Values
const TREE_PRINTER_VAL_POSITION_BELOW = "\"b\""

// Collection delimiters
const TREE_PRINTER_COLLECTION_OPEN = "["
const TREE_PRINTER_COLLECTION_CLOSE = "]"

/*
Prints JSON output format compatible with tree visualization in D3.
Allows indication for flat printing (nested property tree structure vs. flat listing of properties).
Allows indication for printing of binary trees, as opposed to tree aggregated by logical operators for given component.
This function is tested in TabularOutputGenerator_test.go, i.e., tests with focus on visual tree output.
*/
func (s Statement) PrintTree(parent *Node, printFlat bool, printBinary bool, includeAnnotations bool) string {

	// Default name if statement does not have root node
	rootName := ""

	if parent != nil {
		// if it is a nested statement, use component name it nests on as name
		rootName = parent.GetComponentName()
	}

	// Opening tree
	out := "{" + TREE_PRINTER_LINEBREAK
	// Print statement-level node
	out += TREE_PRINTER_KEY_NAME + TREE_PRINTER_EQUALS + "\"" +
		// Root node name
		rootName +
		"\"" + TREE_PRINTER_SEPARATOR

	// Indicates whether children have already been added below the top-level string
	childrenPresent := false

	// Print individual nodes
	components := []*Node{
		// Regulative Side
		s.Attributes,
		s.AttributesPropertySimple,
		s.AttributesPropertyComplex,
		s.Deontic,
		s.Aim,
		s.DirectObject,
		s.DirectObjectComplex,
		s.DirectObjectPropertySimple,
		s.DirectObjectPropertyComplex,
		s.IndirectObject,
		s.IndirectObjectComplex,
		s.IndirectObjectPropertySimple,
		s.IndirectObjectPropertyComplex,
		// Constitutive Side
		s.ConstitutedEntity,
		s.ConstitutedEntityPropertySimple,
		s.ConstitutedEntityPropertyComplex,
		s.Modal,
		s.ConstitutiveFunction,
		s.ConstitutingProperties,
		s.ConstitutingPropertiesComplex,
		s.ConstitutingPropertiesPropertySimple,
		s.ConstitutingPropertiesPropertyComplex,
		// Shared elements
		s.ActivationConditionSimple,
		s.ActivationConditionComplex,
		s.ExecutionConstraintSimple,
		s.ExecutionConstraintComplex,
		s.OrElse}

	for _, v := range components {
		// only print components that have content, and for property components (the one whose name ends on ",p"),
		// only print if they have complex children - simple values are printed alongside first-order component

		if v != nil && (!strings.HasSuffix(v.GetComponentName(), PROPERTY_SYNTAX_SUFFIX) ||
			(strings.HasSuffix(v.GetComponentName(), PROPERTY_SYNTAX_SUFFIX) && !v.HasPrimitiveEntry())) {
			prepend := ""
			if childrenPresent {
				// If in next round (not first entry), prepend separator in case output is produced
				prepend = TREE_PRINTER_SEPARATOR
			}
			// Generate actual entry
			componentString := v.PrintNodeTree(s, printFlat, printBinary, includeAnnotations)
			Println("Output for " + v.GetComponentName() + ": " + componentString)
			if !childrenPresent && componentString != "" {
				// Print children prefix if components are present
				out += TREE_PRINTER_KEY_CHILDREN + TREE_PRINTER_EQUALS + TREE_PRINTER_COLLECTION_OPEN + TREE_PRINTER_LINEBREAK
				childrenPresent = true
			}
			// Add the actual content
			out += prepend + componentString
		}
	}
	// Close children
	if childrenPresent {
		out += TREE_PRINTER_LINEBREAK + TREE_PRINTER_COLLECTION_CLOSE
	}
	// Close entire tree
	out += TREE_PRINTER_LINEBREAK + TREE_PRINTER_CLOSE_BRACE

	return out
}

/*
Returns JSON output for visual tree rendering of individual nodes using D3.
Allows indication for flat printing (nested property tree structure vs. flat listing of properties).
Allows indication for printing of binary trees, as opposed to tree aggregated by logical operators for given component.
Allows indication for including annotations in output.
*/
func (n *Node) PrintNodeTree(stmt Statement, printFlat bool, printBinary bool, includeAnnotations bool) string {
	out := ""

	if !n.IsNil() && !n.IsEmptyNode() {
		if n.HasPrimitiveEntry() || n.IsCombination() {

			// Indicates whether full entry should be printed
			printFullEntry := false

			// If the entry is not a statement but either leaf or combination
			if n.HasPrimitiveEntry() {
				// Produce output for simple entry
				out = TREE_PRINTER_OPEN_BRACE + TREE_PRINTER_KEY_NAME + TREE_PRINTER_EQUALS +
					// Actual content
					"\"" + n.Entry.(string) + "\""

				// Ensure that entry is closed
				printFullEntry = true
			} else if n.IsCombination() {
				// Produce output for combination

				if printBinary {
					// Fall back to full entry parsing either way - resolving full binary tree structure
					printFullEntry = true
				} else {
					// If non-binary, print only leaf entries linked via same logical operator on same component
					// without considering logical operators in output
					if n.Parent != nil && n.LogicalOperator == n.Parent.LogicalOperator &&
						n.GetComponentName() == n.Parent.GetComponentName() {
						// Print left side
						out += n.Left.PrintNodeTree(stmt, printFlat, printBinary, includeAnnotations)

						// Append separator to collapsed entries (i.e., on same level)
						out += ", "

						// Print right side
						out += n.Right.PrintNodeTree(stmt, printFlat, printBinary, includeAnnotations)

						// Suppress printing of closing parts of entry, since further nodes of same operator on same component may be appended
						printFullEntry = false
					} else {
						// Fall back to print full entry if logical operators or components don't match
						printFullEntry = true
					}
				}
				// Prints full entry as binary tree element (either applies if binary tree structure is activated,
				// or if no nested logical operators for given component were detected (e.g., multiple nested AND linkages)
				if printFullEntry {
					// Create logical node, two children, and delegate node entry generation to children
					out = TREE_PRINTER_OPEN_BRACE + TREE_PRINTER_KEY_NAME + TREE_PRINTER_EQUALS +
						// Logical operator
						"\"" + n.LogicalOperator + "\"" + TREE_PRINTER_SEPARATOR +
						// Children
						TREE_PRINTER_KEY_CHILDREN + TREE_PRINTER_EQUALS + TREE_PRINTER_COLLECTION_OPEN

					// Left child
					out += n.Left.PrintNodeTree(stmt, printFlat, printBinary, includeAnnotations)

					// Add separator
					out += TREE_PRINTER_SEPARATOR

					// Right child
					out += n.Right.PrintNodeTree(stmt, printFlat, printBinary, includeAnnotations)

					// Closing collection
					out += TREE_PRINTER_COLLECTION_CLOSE
				}
			}

			// Continue and close full entry (with component, property and annotation information) only if entry is complete,
			// not if branches of logical operators are collapsed
			if printFullEntry {
				// Append component name as link label for any entry
				out += ", " + TREE_PRINTER_KEY_COMPONENT + TREE_PRINTER_EQUALS + "\"" + n.GetComponentName() + "\""

				// Print private properties
				out = n.appendPropertyNodes(out, stmt, printFlat, printBinary, includeAnnotations)

				// Append annotations
				if includeAnnotations {
					out = n.appendAnnotations(out)
				}

				// Close entry
				out += TREE_PRINTER_CLOSE_BRACE
			}
		} else {
			// Produce output for nested statement
			out += n.Entry.(Statement).PrintTree(n, printFlat, printBinary, includeAnnotations)
		}
	}
	return out
}

/*
Appends shared and private nodes to D3-consumable JSON output string based on related properties, as well as own private nodes.
The shared and private property nodes are combined in the order "shared, private".
Note: In flat output mode only primitive private properties are included in the rendered output.
Flat output implies the printing of private properties as labels for component nodes, rather than an own node hierarchy.
Allows indication for printing of binary trees, as opposed to tree aggregated by logical operators for given component.
Includes indication whether annotations are to be included in output.
*/
func (n *Node) appendPropertyNodes(stringToAppendTo string, stmt Statement, printFlat bool, printBinary bool, includeAnnotations bool) string {

	// Append potential private and shared property nodes under the condition that those nodes are leaf nodes, or if flat printing is activated
	if n != nil && (n.IsLeafNode() || printFlat) &&
		// Check for shared and private properties
		len(stmt.GetPropertyComponent(n, false)) > 0 || (len(n.PrivateNodeLinks) > 0 && n.PrivateNodeLinks[0] != nil) {

		// Retrieve relevant component property and combine with private nodes before printing
		allNodes := stmt.GetPropertyComponent(n, false)
		includeAllNodes := true
		// Test whether shared property nodes exist
		if len(allNodes) == 0 || (len(allNodes) > 0 && allNodes[0] == nil) {
			includeAllNodes = false
		}

		// Check whether private nodes are populated

		if len(n.PrivateNodeLinks) > 0 && n.PrivateNodeLinks[0] != nil {

			if includeAllNodes {
				allNodes = append(allNodes, n.PrivateNodeLinks...)
			} else {
				allNodes = n.PrivateNodeLinks
			}
		}

		// Only add output if properties exist
		if len(allNodes) > 0 && allNodes[0] != nil {

			// keeps track whether any element has been extracted
			elementPrinted := false

			// Add individual items
			for _, privateNode := range allNodes {

				// Ensure that node is only added to output if it is truely primitive, not a combination
				if !privateNode.IsCombination() && privateNode.HasPrimitiveEntry() {

					if !elementPrinted {
						// Start general output for property only if nothing is printed yet
						if printFlat {
							// Initiate flat output for properties
							stringToAppendTo += ", " + TREE_PRINTER_KEY_PROPERTIES + TREE_PRINTER_EQUALS
						} else {
							// Add position information to ensure offset printing of leading node content (since it is followed by nested property structure)
							stringToAppendTo += ", " + TREE_PRINTER_KEY_POSITION + TREE_PRINTER_EQUALS + TREE_PRINTER_VAL_POSITION_BELOW
							// Initiate tree structure for tree output of properties
							stringToAppendTo += ", " + TREE_PRINTER_KEY_CHILDREN + TREE_PRINTER_EQUALS + TREE_PRINTER_COLLECTION_OPEN
						}
					}

					if elementPrinted {
						// Add separator if element has been added
						stringToAppendTo += ", "
					} else if printFlat {
						// Prepend quotation
						stringToAppendTo += "\""
					}
					// Print per-property entry
					if printFlat {
						// flat entry (and only append if content present)
						if privateNode.Entry != nil && !privateNode.IsCombination() {
							stringToAppendTo += privateNode.Entry.(string)
						}
					} else {
						// nested tree structure
						stringToAppendTo += privateNode.PrintNodeTree(stmt, printFlat, printBinary, includeAnnotations)
					}
					// Mark if initial item has been performed
					elementPrinted = true
				}
			}
			if elementPrinted {
				if printFlat {
					// Close flat entry
					stringToAppendTo += "\""
				} else {
					// Close collection
					stringToAppendTo += TREE_PRINTER_COLLECTION_CLOSE
				}
			}
		}
	}

	// Return extended string
	return stringToAppendTo
}

/*
Appends potentially existing annotations to node-specific output.
*/
func (n *Node) appendAnnotations(stringToAppendTo string) string {
	// Append potential annotations
	if n.GetAnnotations() != nil {
		stringToAppendTo += ", " + TREE_PRINTER_KEY_ANNOTATIONS + TREE_PRINTER_EQUALS
		stringToAppendTo += "\"" + n.GetAnnotations().(string) + "\""
	}
	// Return potentially extended string
	return stringToAppendTo
}
