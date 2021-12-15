package tree

import (
	"fmt"
	"strings"
)

const TREE_PRINTER_OPEN_BRACE = "{"
const TREE_PRINTER_CLOSE_BRACE = "}"
const TREE_PRINTER_LINEBREAK = "\n"
const TREE_PRINTER_SEPARATOR = "," + TREE_PRINTER_LINEBREAK

const TREE_PRINTER_KEY_NAME = "\"name\""
const TREE_PRINTER_KEY_COMPONENT = "\"comp\""
const TREE_PRINTER_KEY_CHILDREN = "\"children\""
const TREE_PRINTER_KEY_PROPERTIES = "\"properties\""
const TREE_PRINTER_KEY_ANNOTATIONS = "\"annotations\""
const TREE_PRINTER_EQUALS = ": "

const TREE_PRINTER_COLLECTION_OPEN = "["
const TREE_PRINTER_COLLECTION_CLOSE = "]"

/*
Prints JSON output format compatible with tree visualization in D3.
This function is tested in TabularOutputGenerator_test.go#TestStatement_PrintTree.
*/
func (s Statement) PrintTree(parent *Node, includeAnnotations bool) string {

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
			componentString := v.PrintNodeTree(s, includeAnnotations)
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
Returns visual tree output (for D3) of individual nodes.
*/
func (n *Node) PrintNodeTree(stmt Statement, includeAnnotations bool) string {
	out := ""

	// Indicates whether output is complex (tree structure), or a flat listing of properties
	printFlat := true

	if !n.IsNil() && !n.IsEmptyNode() {
		if n.HasPrimitiveEntry() || n.IsCombination() {
			// If the entry is not a statement but either leaf or combination
			if n.HasPrimitiveEntry() {
				// Produce output for simple entry
				out = TREE_PRINTER_OPEN_BRACE + TREE_PRINTER_KEY_NAME + TREE_PRINTER_EQUALS +
					// Actual content
					"\"" + n.Entry.(string) + "\""
			} else if n.IsCombination() {
				// Produce output for combination

				// Create logical node, two children, and delegate node entry generation to children
				out = TREE_PRINTER_OPEN_BRACE + TREE_PRINTER_KEY_NAME + TREE_PRINTER_EQUALS +
					// Logical operator
					"\"" + n.LogicalOperator + "\"" + TREE_PRINTER_SEPARATOR +
					// Children
					TREE_PRINTER_KEY_CHILDREN + TREE_PRINTER_EQUALS + TREE_PRINTER_COLLECTION_OPEN +
					// Left child
					n.Left.PrintNodeTree(stmt, includeAnnotations) + TREE_PRINTER_SEPARATOR +
					// Right child
					n.Right.PrintNodeTree(stmt, includeAnnotations) +
					// Closing collection
					TREE_PRINTER_COLLECTION_CLOSE
			}

			// Append component name as link label for any entry
			out += ", " + TREE_PRINTER_KEY_COMPONENT + TREE_PRINTER_EQUALS + "\"" + n.GetComponentName() + "\""

			// Print private properties
			out = n.appendPrivateNodes(out, stmt, printFlat, includeAnnotations)

			// Append annotations
			if includeAnnotations {
				out = n.appendAnnotations(out)
			}

			// Close entry
			out += TREE_PRINTER_CLOSE_BRACE

		} else {
			// Produce output for nested statement
			out += n.Entry.(Statement).PrintTree(n, includeAnnotations)
		}
	}
	return out
}

func (n *Node) appendPrivateNodes(stringToAppendTo string, stmt Statement, printFlat bool, includeAnnotations bool) string {
	// Append potential private nodes
	if len(stmt.GetPropertyComponent(n, false)) > 0 || (len(n.PrivateNodeLinks) > 0 && n.PrivateNodeLinks[0] != nil) {

		// keeps track whether any element has been extracted
		elementPrinted := false

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

			fmt.Println("Component: " + n.GetComponentName())
			fmt.Println("Nodes: ", allNodes)

			// Start general output for property
			stringToAppendTo += ", " + TREE_PRINTER_KEY_PROPERTIES + TREE_PRINTER_EQUALS

			if !printFlat {
				stringToAppendTo += "\"Properties\"" + ", " + TREE_PRINTER_KEY_CHILDREN + TREE_PRINTER_EQUALS + TREE_PRINTER_COLLECTION_OPEN
			}

			// Add individual items
			for _, privateNode := range allNodes {
				if elementPrinted {
					// Add separator if element has been added
					stringToAppendTo += ", "
				} else if printFlat {
					// Prepend quotation
					stringToAppendTo += "\""
				}
				// Print per-property entry
				if printFlat {
					fmt.Println("Printing node: ", privateNode)
					// flat entry
					stringToAppendTo += privateNode.Entry.(string)
				} else {
					// nested tree structure
					stringToAppendTo += privateNode.PrintNodeTree(stmt, includeAnnotations)
				}
				// Mark if initial item has been performed
				elementPrinted = true
			}
			if elementPrinted && printFlat {
				// Close flat entry
				stringToAppendTo += "\""
			}

			// Close collection
			if !printFlat {
				stringToAppendTo += TREE_PRINTER_COLLECTION_CLOSE
			}
		}
	}

	// Return extended string
	return stringToAppendTo
}

func (n *Node) appendAnnotations(stringToAppendTo string) string {
	// Append potential annotations
	if n.GetAnnotations() != nil {
		stringToAppendTo += ", " + TREE_PRINTER_KEY_ANNOTATIONS + TREE_PRINTER_EQUALS
		stringToAppendTo += "\"" + n.GetAnnotations().(string) + "\""
	}
	// Return potentially extended string
	return stringToAppendTo
}
