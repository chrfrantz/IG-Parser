package main

import (
	"IG-Parser/tree"
	"fmt"
	"log"
	"strconv"
)

func main() {

	// Dummy input for convenient selection
	input := ""
	// Very simple input
	input = "(inspect and [OR] party)"
	// Simple input
	//input = "((inspect and [OR] party) [AND] sing)"
	// Proper complex input
	//input = "((inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray))"
	// Imbalanced parentheses will fail (whatever the direction)
	//input = "((inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray)"
	//input = "(inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray))"
	// Missing outer parentheses will lead to processing of right side only
	//input = "(inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray)"
	// Invalid operators lead to ignoring element in processing
	//input = "((inspect and [OR] party) [AND] ((review [XOR] muse) AND pray))"
	// Invalid operators even apply to parenthesized combinations --> FIX and make leaf node
	//input = "((inspect and OR party) [AND] ((review [XOR] muse) [AND] pray))"
	// Invalid operators without parentheses around operands will be treated as leaf node
	//input = "(inspect and OR party [AND] ((review [XOR] muse) [AND] pray))"
	// Excessive parentheses are acceptable (will be flattened in parsing process)
	input = "((((inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray))))"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	parseDepth(input, &node)
	// Print resulting tree
	fmt.Println("Final tree: \n" + node.String())
}

/*
Parses combinations in string. The syntactic form of input is:
"( leftSide [OPERATOR] rightSide )", where [OPERATOR] is one
of the logical operators [AND], [OR], [XOR] (including brackets),
and left and right side are either text or combinations themselves.
Note:
- The entire expression must be surrounded with parentheses, else only
the right-most outer combination (and combinations nested therein) is parsed.
- Parsing checks for matching parentheses and stops otherwise
- Invalid combinations (e.g., missing logical operator) are discarded
in the processing.
 */
func parseDepth(input string, node *tree.Node) *tree.Node {
	combinations := detectCombinations(input)
	if len(combinations) == 0 {
		fmt.Println("No combinations detected.")
	} else {
		ct := 0
		toBeDeleted := []int{}
		for k, v := range combinations {
			if v.Complete {
				ct++
			} else {
				toBeDeleted = append(toBeDeleted, k)
			}
		}
		errorMsg := ""
		if len(toBeDeleted) > 0 {
			errorMsg = ", with one partial " + strconv.Itoa(len(toBeDeleted)) + " to be removed"
		}
		fmt.Println(strconv.Itoa(ct) + " valid combination detected" + errorMsg + "!")
		fmt.Println(combinations)

		// Clean up invalid entries
		i := 0
		for i < len(toBeDeleted) {
			delete(combinations, toBeDeleted[i])
			i++
		}
	}

	// if no valid combinations are left, the parsing is finished
	if len(combinations) == 0 {
		return node
	}

	// Now the parsing of nested combinations starts

	// Depth first
	v := 0
	for { // infinite loop - breaks out eventually
		fmt.Println("Level " + strconv.Itoa(v))
		if _, ok := combinations[v]; ok {
			left := input[combinations[v].Left:combinations[v].Operator]
			right := input[combinations[v].Operator+len(combinations[v].OperatorVal)+2:combinations[v].Right]
			fmt.Println("==Left value: " +  left)
			fmt.Println("==Operator: " + combinations[v].OperatorVal)
			fmt.Println("==Right value: " +  right)
			// Scan through top level only and break out afterwards ...

			// Assign logical operator
			node.LogicalOperator = combinations[v].OperatorVal

			fmt.Println("Tree after adding logical operator: " + node.String())

			// Left side
			leftCombos := detectCombinations(left)
			// Assign nested nodes either way
			leftNode := tree.Node{}
			// Link both nodes
			node.Left = &leftNode
			leftNode.Parent = node

			if len(leftCombos) == 0 {
				// If no combinations exist, assign as left leaf
				fmt.Println("Found leaf on left side: " +  left)
				node.Left.Entry = left
			} else {
				// If combinations exist, delegate
				fmt.Println("Go deep on left side: " +  left)
				parseDepth(left, &leftNode)
				fmt.Println("Tree after processing left deep: " + node.String())
			}

			// Right side
			rightCombos := detectCombinations(right)
			// Assign nested nodes either way
			rightNode := tree.Node{}
			// Link both nodes
			node.Right = &rightNode
			rightNode.Parent = node
			if len(rightCombos) == 0 {
				// If no combinations exist, assign as left leaf
				fmt.Println("Found leaf on right side: " +  right)
				node.Right.Entry = right
			} else {
				// If combinations exist, delegate
				fmt.Println("Go deep on right side: " +  right)
				parseDepth(right, node.Right)
				fmt.Println("Tree after processing right deep: " + node.String())
			}
			// break out if combination has been found and processed - must be top-level combination
			return node
		} else {
			fmt.Println("==No combination for key/level " + strconv.Itoa(v))
		}
		v++
	}

	//fmt.Println("Should not really reach here; probably empty node: " + node.String())
	return node
}

/*
This function detects levels of combinations present in the expression
and returns the boundary indices as well logical operator (where present).
To signal incomplete combinations, it contains a Complete flag that signals
completeness for further postprocessing.
Note: This function does not extract all combinations present in the expression,
since combinations on the same level will not be detected, but overwritten.
In essence, the function provides the depth of the nesting in the expression.

Default syntactic form of input: "( leftSide [OPERATOR] rightSide )", where
[OPERATOR] is one of the logical operators (including brackets), and left
and right side are either text or combinations themselves.
 */
func detectCombinations(expression string) map[int]tree.Boundaries {

	level := 0
	parCount := 0
	mode := ""

	levelMap := make(map[int]tree.Boundaries)


	fmt.Println("Testing expression " + expression)
	for i, letter := range expression {

		switch string(letter) {
		case "(":
			// Increase level
			level++
			fmt.Println("Expression start detected (Level " + strconv.Itoa(level) + ")")
			// Configure mode
			mode = tree.PARSING_MODE_LEFT
			fmt.Println("Mode: " +  mode)
			//Test of existing entries
			if _, ok := levelMap[level]; ok {
				fmt.Println("Key already defined - not added")
			} else {
				// Store index reference (incremented with 1 to avoid left parenthesis - balanced)
				levelMap[level] = tree.Boundaries{Left: i+1, Complete: false}
			}
			// Count parentheses to detect uneven matching
			parCount++
		case ")":
			fmt.Println("Expression end detected (Level " + strconv.Itoa(level) + ")")
			// Store index reference
			fmt.Println("Level before saving: " + strconv.Itoa(level))
			if b, ok := levelMap[level]; ok {
				b.Right = i
				if b.Left != 0 && b.OperatorVal != "" {
					b.Complete = true
				} else {
					fmt.Println("Detected end, but combination incomplete! (Missing operator or left parenthesis)")
				}
				levelMap[level] = b
			}
			// Configure mode
			mode = tree.PARSING_MODE_OUTSIDE_EXPRESSION
			fmt.Println("Mode: " +  mode)
			// Reduce level
			level--
			// Count parentheses to detect uneven matching
			parCount--
		case "[":
			//fmt.Println("Checking for logical operator ... " + expression[i:i+5])
			foundOperator := ""
			switch expression[i : i+5] {
			case tree.AND_BRACKETS:
				fmt.Println("Detected " + tree.AND_BRACKETS)
				foundOperator = tree.AND
			case tree.OR_BRACKETS + " ": // artificial extension to match other operators' length
				fmt.Println("Detected " + tree.OR_BRACKETS)
				foundOperator = tree.OR
			case tree.XOR_BRACKETS:
				fmt.Println("Detected " +  tree.XOR_BRACKETS)
				foundOperator = tree.XOR
			}
			fmt.Println("Found logical operator " + foundOperator + " on level " + strconv.Itoa(level))
			// Configure mode
			mode = tree.PARSING_MODE_RIGHT
			fmt.Println("Mode: " +  mode)
			// Store index reference and value
			if o, ok := levelMap[level]; ok {
				o.Operator = i
				o.OperatorVal = foundOperator
				levelMap[level] = o
			}
		}
	}

	if parCount != 0 {
		log.Fatal("Uneven number of parentheses (positive --> too many left; negative --> too many right): " + strconv.Itoa(parCount))
	}

	return levelMap
}
