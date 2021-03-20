package parser

import (
	"IG-Parser/tree"
	"fmt"
	"github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"strconv"
	"strings"
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
	// Non-combinations are ignored (i.e., parentheses within components)
	input = "(French (and) [AND] (accredited certifying agents))"
	// Repeated AND combinations (e.g., "expr1 AND expr2 AND expr3") are collapsed into nested structures (e.g., "(expr1 AND expr2) AND expr3")
	input = "(French (and) [AND] (certified production and [XOR] handling operations) and [AND] (accredited certifying agents))"
	// Repeated logical operators will break (empty leaf value in the AND case)
	//input = "(French (and) [AND] [AND] (certified production and [XOR] handling operations) and [AND] (accredited certifying agents))"
	// Repeated logical operators will break (multiple non-AND operators (or mix thereof))
	//input = "(French (and) [AND] [OR] (certified production and [XOR] handling operations) and [AND] (accredited certifying agents))"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, modifiedInput, _ := ParseDepth(input, &node)
	// Print resulting tree
	fmt.Println("Final tree: \n" + node.String())
	fmt.Println("Corresponding (potentially modified) input string: " + modifiedInput)

	fmt.Println(node.Stringify())
}

/*
Parses combinations in string. The syntactic form of input is:
"( leftSide [OPERATOR] rightSide )", where [OPERATOR] is one
of the logical operators [AND], [OR], [XOR] (including brackets),
and left and right side are either text or combinations themselves.
For [AND] operators, an arbitrary number of expressions can be combined;
in this case the function will decompose those into nested structures
(e.g., expanding "( expr1 [AND] expr2 [AND] expr3 )" into
"(( expr1 [AND] expr2 ) [AND] expr3)"), with precedence for left combinations.
Note that expressions are trimmed prior storing in tree structure.
The parsing further supports shared values outside of the combination (e.g.,
'(shared left value (left element [AND] right element) shared right value)',
and returns those as part of the node that holds the logical operator.

Hint: Call Stringify() on the returned node to reconstruct string

The function returns
- a node tree of the structure, as well as
- the potentially modified input string corresponding to the node tree
  Note: Shared elements are stripped from the modified output string (but
  included in the node instance

Note:
- The entire expression must be surrounded with parentheses, else only
the right-most outer combination (and combinations nested therein) is parsed.
- Parsing checks for matching parentheses and stops otherwise
- Invalid combinations (e.g., missing logical operator) are discarded
in the processing.
 */
func ParseDepth(input string, nodeTree *tree.Node) (*tree.Node, string, tree.ParsingError) {
	// Return detected combinations, alongside potentially modified input string
	combinations, _, input, err := detectCombinations(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nodeTree, input, err
	}

	/*if len(combinations) == 0 {
		fmt.Println("No combinations detected.")
	} else {
		//TODO to be reviewed if issues arise with unwanted elements; else remove
		/*ct := 0
		toBeDeleted := []int{}
		for k, level := range combinations {
			if level.Complete {
				ct++
			} else {
				toBeDeleted = append(toBeDeleted, k)
			}
		}
		errorMsg := ""
		if len(toBeDeleted) > 0 {
			errorMsg = ", with one partial " + strconv.Itoa(len(toBeDeleted)) + " to be removed"
		}
		fmt.Println(strconv.Itoa(ct) + " valid combination detected" + errorMsg + "!")*/


		// Clean up invalid entries
		/*i := 0
		for i < len(toBeDeleted) {
			delete(combinations, toBeDeleted[i])
			i++
		}
	}*/

	// if no valid combinations are left, the parsing is finished
	if len(combinations) == 0 {
		nodeTree.Entry = input
		return nodeTree, input, tree.ParsingError{ErrorCode: tree.PARSING_NO_COMBINATIONS,
			ErrorMessage: "The input does not contain combinations"}
	}

	// Now the parsing of nested combinations starts
	fmt.Print("STARTING TREE CONSTRUCTION: Detected combinations: ")
	fmt.Println(combinations)


	// Depth first
	level := 0

	// Link input node to temporary node for incremental addition of elements - needs to be copied back prior to return
	node := nodeTree
	// Map to keep order of entry for correct tree construction (retention of order)
	orderMap := make(map[int]*tree.Node)

	// Iterate through all levels
	for level <= len(combinations) {
		fmt.Println("Level " + strconv.Itoa(level))
		fmt.Print("Combinations on level " + strconv.Itoa(level) + ": ")
		fmt.Println(combinations[level])

		idx := 0
		//Iterate through all indices
		for idx < len(combinations[level]) {

			fmt.Println("TREE BEFORE NEXT COMBINATION: " + node.String())

			// Parse complete combinations, and combinations can extract their respective shared elements
			if combinations[level][idx].Complete {

				fmt.Println("Found combination to parse on level " + strconv.Itoa(level) +
					", Index: " + strconv.Itoa(idx))

				node := &tree.Node{}

				// full parsing
				left := input[combinations[level][idx].Left:combinations[level][idx].Operator]
				right := input[combinations[level][idx].Operator+len(combinations[level][idx].OperatorVal)+2 : combinations[level][idx].Right]
				fmt.Println("==Raw Left value: " + left)
				fmt.Println("==Raw Operator: " + combinations[level][idx].OperatorVal)
				fmt.Println("==Raw Right value: " + right)

				// Check for shared elements by sending ID to parse function and search boundaries for next lower level and embracing
				sharedLeft, sharedRight := extractSharedComponents(input, combinations, level, idx)

				// Assign left shared value if existing
				if sharedLeft != nil {
					fmt.Println("Assigning left shared element '" + fmt.Sprint(sharedLeft) + "'.")
					node.SharedLeft = sharedLeft
					// Reset shared
					sharedLeft = []string{}
				}
				// Assign shared value if existing
				if sharedRight != nil {
					fmt.Println("Assigning right shared element '" + fmt.Sprint(sharedRight) + "'.")
					node.SharedRight = sharedRight
					// Reset shared
					sharedRight = []string{}
				}

				// Assign logical operator
				node.LogicalOperator = combinations[level][idx].OperatorVal

				fmt.Println("Tree before deep parsing: " + node.String())

				// Left side (potentially modifying input string)
				leftCombos, leftNonShared, left, err := detectCombinations(left)
				if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
					log.Println("Error when parsing left side: " + err.ErrorMessage)
					return node, input, err
				}
				if err.ErrorCode == tree.PARSING_ERROR_IGNORED_ELEMENTS {
					log.Print("Warning: Discarded elements during deep left parsing: " + tree.PrintArray(err.ErrorIgnoredElements))
				}

				// Assign nested nodes either way
				leftNode := tree.Node{}
				// Link both nodes
				node.Left = &leftNode
				leftNode.Parent = node

				if len(leftCombos) == 0 {
					// Trim content first
					left = strings.Trim(left, " ")
					if left != "" {
						// If no combinations exist, assign as left leaf
						fmt.Println("Found leaf on left side: " + left)
						node.Left.Entry = left
					} else {
						msg := "Empty leaf value on left side: " + left +
							" (Corresponding right value and operator: " + right + "; " + node.LogicalOperator +
							");\n processed expression: " + input
						log.Println(msg)
						return nil, input, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_EMPTY_LEAF,
							ErrorMessage: msg}
					}
				} else {
					// If combinations exist, delegate

					// Deal with non-shared items
					if len(leftNonShared) > 0 {
						fmt.Println("NOT HANDLED - JUST INFORMATION: Non-shared elements left: " + fmt.Sprint(leftNonShared))
						// TODO: Assign the left node - order not managed
					}

					fmt.Println("Go deep on left side: " + left)
					_, left, err := ParseDepth(left, &leftNode)
					if err.ErrorCode != tree.PARSING_NO_ERROR {
						return nil, left, err
					}

					// Check for inheriting shared elements on AND nodes
					inheritSharedElements(&leftNode)

					fmt.Println("Tree after processing left deep: " + node.String())
				}

				// Right side (potentially modifying input string)
				rightCombos, rightNonShared, right, err := detectCombinations(right)
				if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
					log.Println("Error when parsing right side: " + err.ErrorMessage)
					return node, input, err
				}
				if err.ErrorCode == tree.PARSING_ERROR_IGNORED_ELEMENTS {
					log.Print("Warning: Discarded elements during deep right parsing: " + tree.PrintArray(err.ErrorIgnoredElements))
				}

				// Assign nested nodes either way
				rightNode := tree.Node{}
				// Link both nodes
				node.Right = &rightNode
				rightNode.Parent = node

				if len(rightCombos) == 0 {
					// Trim content first
					right = strings.Trim(right, " ")
					if right != "" {
						// If no combinations exist, assign as left leaf
						fmt.Println("Found leaf on right side: " + right)
						node.Right.Entry = right
					} else {
						msg := "Empty leaf value on right side: " + right + " (Corresponding left value and operator: " + left + "; " + node.LogicalOperator +
							");\n processed expression: " + input
						log.Println(msg)
						return nil, input, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_EMPTY_LEAF,
							ErrorMessage: msg}
					}
				} else {
					// If combinations exist, delegate

					// Deal with non-shared items
					if len(rightNonShared) > 0 {
						fmt.Println("NOT HANDLED - JUST INFORMATION: Non-shared elements right: " + fmt.Sprint(rightNonShared))
						// TODO: Assign to right node - order not managed
					}

					fmt.Println("Go deep on right side: " + right)
					_, right, err := ParseDepth(right, node.Right)
					if err.ErrorCode != tree.PARSING_NO_ERROR {
						return nil, right, err
					}

					// Check for inheriting shared elements on AND nodes
					inheritSharedElements(&rightNode)

					fmt.Println("Tree after processing right deep: " + node.String())
				}

				//Adds node to map with starting character index as key - for later reconstruction of tree prior to return
				orderMap[combinations[level][idx].Left] = node

			}
			fmt.Println("==Finished parsing index " + strconv.Itoa(idx) + " on level " + strconv.Itoa(level))
			// Increase to explore other entries
			idx++
		}
		// break out if combination has been found and processed - must be top-level combination
		// TODO - check output
		//return node, "(" + left + " [" + node.LogicalOperator + "] " + right + ")", tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		//return node, "", tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		fmt.Println("==Finished parsing on level " + strconv.Itoa(level))

		// Increase level
		level++

	}

	fmt.Print("Combinations before return: ")
	fmt.Println(combinations)

	// Reconstructing node based on order of node input
	//fmt.Println("Node order: ")
	ct := 0
	for ct < len(input) {
		if _, ok := orderMap[ct]; ok {
			nodeTree.Insert(orderMap[ct])
			//fmt.Println("Added to tree: " + fmt.Sprint(orderMap[ct]))
		}
		ct++
	}

	fmt.Println("RETURNING FINAL NODE: " + nodeTree.String())
	return nodeTree, input, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
This function detects levels of combinations present in the expression
and returns the boundary indices as well logical operator (where present).
It further returns an array of text elements that exist outside of combinations,
and the input string in potentially modified form to reflect
changes performed during processing.
To signal incomplete combinations, it contains a Complete flag that signals
completeness for further postprocessing.
Note: This function does not extract all combinations present in the expression,
since combinations on the same level will not be detected, but overwritten.
In essence, the function provides the depth of the nesting in the expression.

Default syntactic form of input: "( leftSide [OPERATOR] rightSide )", where
[OPERATOR] is one of the logical operators (including brackets), and left
and right side are either text or combinations themselves.
 */
func detectCombinations(expression string) (map[int]map[int]tree.Boundaries, []string, string, tree.ParsingError) {

	// Tracks current parsing level
	level := 0

	// Parentheses count to check for balance
	parCount := 0

	// Initial test run for parentheses
	for i, letter := range expression {

		switch string(letter) {
		case "(":
			parCount++
		case ")":
			parCount--
		}
		i++
	}
	if parCount != 0 {
		msg := "Uneven number of parentheses (positive --> too many left; negative --> too many right): " + strconv.Itoa(parCount)
		log.Println(msg)
		return nil, nil, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IMBALANCED_PARENTHESES, ErrorMessage: msg}
	}
	// Passed parentheses count

	// Map of mode states across levels (to recover during parsing)
	modeMap := make(map[int]string)
	// Default mode for level 0
	modeMap[level] = tree.PARSING_MODE_OUTSIDE_EXPRESSION

	// Maintain map of boundaries across different levels (with level as key, followed by index (if multiple entries) and corresponding value)
	levelMap := make(map[int]map[int]tree.Boundaries)

	// Collection of found operators (with operator as key, followed by level, and count as value)
	foundOperators := make(map[string]map[int]int)

	// Holds the current non-shared element (to be stored into nonSharedElements upon completion)
	nonSharedElement := ""
	// Stores reference to last character index in which an non-shared element was identified (break in between suggests new token)
	lastCharIndex := 0
	// Holds all non-shared elements across entire expression (to be returned with identified combinations)
	nonSharedElements := []string{}

	fmt.Println("Testing expression " + expression)
	for i, letter := range expression {

		switch string(letter) {
		case "(":
			// Increase level
			level++
			fmt.Println("Expression start detected (Level " + strconv.Itoa(level) + ")")
			// Configure mode
			modeMap[level] = tree.PARSING_MODE_LEFT
			fmt.Println("Mode: " + modeMap[level])
			//Test of existing entries
			if _, ok := levelMap[level]; ok {
				fmt.Println(strconv.Itoa(len(levelMap[level])) + " key(s) already defined - increasing count")
				levelMap[level][len(levelMap[level])] = tree.Boundaries{Left: i + 1, Complete: false}
			} else {
				// Store index reference and create internal map (incremented with 1 to avoid left parenthesis - balanced)
				mp := make(map[int]tree.Boundaries)
				levelMap[level] = mp
				levelMap[level][0] = tree.Boundaries{Left: i + 1, Complete: false}
			}
			// Count parentheses to detect uneven matching
			parCount++
		case ")":
			fmt.Println("Expression end detected (Level " + strconv.Itoa(level) + ")")
			// Check if there are repetitions
			for k, _ := range foundOperators { // key: operator, value:
				for k2, v2 := range foundOperators[k] { // key: level, value: number of occurrences
					if v2 > 1 {
						log.Println("Found " + strconv.Itoa(v2) + " occurrences of operator " + k + " on level " + strconv.Itoa(k2) + " in map.")
					}
				}

			}
			// Store index reference
			levelIdx := len(levelMap[level])-1
			fmt.Println("Level before saving: " + strconv.Itoa(level) + "; Index: " + strconv.Itoa(levelIdx))
			//if b, ok := levelMap[level][levelIdx]; ok {
			b := levelMap[level][levelIdx]
				b.Right = i
				levelMap[level][levelIdx] = b
				fmt.Println("Level map for level " + strconv.Itoa(level) + " (after adding right value but before assessing completeness): " + b.String())
				// Test whether indices are identical or immediately following - suggesting gaps in values
				if ((b.Operator + len(b.OperatorVal) + 2) == b.Right) {
					msg := "Input contains invalid combination expression in the range '" + expression[b.Left:b.Right] + "'."
					return levelMap, nonSharedElements, expression, tree.ParsingError{ErrorCode: tree.PARSING_INVALID_COMBINATION,
						ErrorMessage: msg}
				}
				if b.Left != 0 && b.OperatorVal != "" {
					// Update complete marker
					//b := levelMap[level][levelIdx]
					b.Complete = true
					levelMap[level][levelIdx] = b
					fmt.Println("Expression is complete.")
				} else {
					fmt.Println("Detected end, but combination incomplete (Missing operator or left parenthesis). " +
						"Discarding combination. (Input: '" + expression + "')")
					fmt.Print("Map on given level after processing: ")
					fmt.Println(levelMap[level])
					// Retrieve higher level
					/*_, ok := levelMap[level+1]
					if !ok {
						fmt.Print("No nested complete combination, so deleted this incomplete one: ")
						fmt.Println(levelMap[level])
						// If no higher-level exists (within the lower level), delete this incomplete entry
						delete(levelMap, level)
					} else {
						// else retain
						fmt.Println("Nested complete combination, so incomplete higher-level combination is retained.")
					}*/
				}
			//}

			// Configure mode
			modeMap[level] = tree.PARSING_MODE_OUTSIDE_EXPRESSION
			fmt.Println("Mode: " + modeMap[level])

			// Reset operator count for given level
			for op := range foundOperators {
				fmt.Println("Deleting operator " + op + " for level " + strconv.Itoa(level))
				delete(foundOperators[op], level)
			}

			// Reduce level
			level--
			fmt.Println("Moving back to level " + strconv.Itoa(level) + ", Mode: " + modeMap[level])
			// Count parentheses to detect uneven matching
			parCount--
		case "[":
			//fmt.Println("Checking for logical operator ... " + expression[i:i+5])
			foundOperator := ""
			switch expression[i:i+5] {
			case tree.AND_BRACKETS:
				fmt.Println("Detected " + tree.AND_BRACKETS)
				foundOperator = tree.AND
			case tree.XOR_BRACKETS:
				fmt.Println("Detected " + tree.XOR_BRACKETS)
				foundOperator = tree.XOR
			}
			// Separately test for OR due to differing length
			if foundOperator == "" && expression[i:i+4] == tree.OR_BRACKETS {
				fmt.Println("Detected " + tree.OR_BRACKETS)
				foundOperator = tree.OR
			}
			if foundOperator != "" {

				fmt.Println("Found logical operator " + foundOperator + " on level " + strconv.Itoa(level))

				levelIdx := len(levelMap[level])-1
				// Check whether the logical operator is immediately adjacent to left parenthesis (e.g., ... ([AND] ... - invalid combination
				if levelMap[level][levelIdx].Left == i {
					msg := "Input contains invalid combination expression in the range '" + expression[levelMap[level][levelIdx].Left:] + "'."
					log.Println(msg)
					return levelMap, nonSharedElements, expression, tree.ParsingError{ErrorCode: tree.PARSING_INVALID_COMBINATION,
						ErrorMessage: msg}
				}

				// Store operators
				if _, ok := foundOperators[foundOperator]; ok {
					if _, ok2 := foundOperators[foundOperator][level]; ok2 {
						// if entry exists, increment
						foundOperators[foundOperator][level] = foundOperators[foundOperator][level] + 1
						fmt.Println(" -> Added. Count: " + strconv.Itoa(foundOperators[foundOperator][level]))
					} else {
						// else create new level entry with default value of 1
						foundOperators[foundOperator][level] = 1
						fmt.Println(" -> Created. Count: " + strconv.Itoa(foundOperators[foundOperator][level]))
					}
				} else {
					// if no operator entry exists, else create new operator entry with default value of 1
					foundOperators[foundOperator] = make(map[int]int)
					foundOperators[foundOperator][level] = 1
					fmt.Println(" -> Created level and value. Count: " + strconv.Itoa(foundOperators[foundOperator][level]))
				}

				// If already in right parsing mode, there should be no operator
				if modeMap[level] == tree.PARSING_MODE_RIGHT {
					log.Println("Found additional operator [" + foundOperator + "] (now " + strconv.Itoa(foundOperators[foundOperator][level]) +
						" times on level " + strconv.Itoa(level) + "), even though looking for terminating parenthesis.")
					if foundOperator == tree.AND && foundOperators[foundOperator][level] > 1 { // if AND operator and multiple on the same level
						// Consider injecting a left parenthesis before the expression and add mixfix ") " before logical operator, e.g., "( left ... [AND] right ... ) [AND] ..."
						expression = expression[:levelMap[level][levelIdx].Left] + "(" + expression[levelMap[level][levelIdx].Left:i-1] + ")" + expression[i:]
						log.Println("Multiple [AND] operators found. Reconstructed nested structure by introducing parentheses, now: " + expression)
						log.Println("Rerunning all parsing on combination to capture nested AND combinations")
						return detectCombinations(expression)
					} else {
						return levelMap, nonSharedElements, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS,
							ErrorMessage: "Error: Duplicate non-[AND] operators (or mix of [AND] and non-[AND] operators) on level " + strconv.Itoa(level) +
								" in single expression (Expression: " + expression + ")"}
					}
				}

				// Configure mode
				modeMap[level] = tree.PARSING_MODE_RIGHT
				fmt.Println("Mode: " + modeMap[level])

				// Store index reference and value
				o := levelMap[level][levelIdx]
				o.Operator = i
				o.OperatorVal = foundOperator
				levelMap[level][levelIdx] = o
				/*if o := levelMap[level][levelIdx]; o {
					o.Operator = i
					o.OperatorVal = foundOperator
					levelMap[level] = o
				}*/
			}
		default:
			// Note that this is called for every non-combination element; requires higher-level filtering for leaf entries
			if modeMap[level] == tree.PARSING_MODE_OUTSIDE_EXPRESSION {
				// Copy cached string to array if disruption in continuation in expression (e.g., combinations in between)
				if lastCharIndex - i > 1 && nonSharedElement != "" {
					nonSharedElements = append(nonSharedElements, nonSharedElement)
					nonSharedElement = ""
				}
				// Append current letter to cached string
				nonSharedElement += string(letter)
				// Update index of last letter identified as non-shared
				lastCharIndex = i
			}
		}
	}

	// Copy outstanding cached non-shared element into elements array before returning
	if nonSharedElement != "" {
		nonSharedElements = append(nonSharedElements, nonSharedElement)
	}

	if parCount != 0 {
		msg := "Uneven number of parentheses (positive --> too many left; negative --> too many right): " + strconv.Itoa(parCount)
		log.Println(msg)
		return nil, nil, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IMBALANCED_PARENTHESES, ErrorMessage: msg}
	}

	// Check for non-parsed prefix or suffix of input string
	/*i := 0
	firstIdx := -1
	lastIdx := -1
	for i < len(levelMap) {
		if _, ok := levelMap[i]; ok {
			if firstIdx == -1 {
				// Assign first value
				firstIdx = levelMap[i].Left
				fmt.Println("Prefix pos: " + strconv.Itoa(firstIdx))
			}
			if levelMap[i].Right > lastIdx {
				// Assign highest last index
				lastIdx = levelMap[i].Right
				fmt.Println("Suffix pos: " + strconv.Itoa(lastIdx))
			}
		}
		i++
	}
	prefix := ""
	suffix := ""
	if firstIdx > 0 {
		prefix = strings.Trim(expression[:firstIdx], " (")
	}
	if lastIdx != -1 {
		suffix = strings.Trim(expression[lastIdx+1:], ") ")
	}
	if prefix != "" || suffix != "" {
		fmt.Println("Prefix: " + prefix)
		fmt.Println("Suffix: " + suffix)
		ignoredElements := []string{}
		errorString := ""
		if prefix != "" {
			ignoredElements = append(ignoredElements, prefix)
			errorString += prefix
		}
		if suffix != "" {
			ignoredElements = append(ignoredElements, suffix)
			if errorString != "" {
				errorString += ", "
			}
			errorString += suffix
		}
		fmt.Println("Returning expression (ignored elements: " + errorString + "): " + expression)
		return levelMap, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IGNORED_ELEMENTS,
			ErrorMessage: "Parsing was successful, but expression parts were ignored during coding (" + errorString + "). " +
			"This commonly occurs when logical operators between simple strings and combinations are omitted " +
			"(e.g., ... some string (left [AND] right) ...) and not wrapped by parentheses to signal shared elements. " +
			"In this case, simple strings are ignored in the parsing process.",
			ErrorIgnoredElements: ignoredElements}
	}*/


	fmt.Println("Returning expression (complete parsing): " + expression)
	// if no omitted elements during parsing, regular return without error
	return levelMap, nonSharedElements, expression, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Processes potential inheritance of shared element values from parent to child nodes,
where both parent and child nodes have AND operators
 */
func inheritSharedElements(node *tree.Node) {
	if node.LogicalOperator == tree.AND &&
		node.Parent.LogicalOperator == tree.AND &&
		SHARED_ELEMENT_INHERITANCE_MODE != SHARED_ELEMENT_INHERIT_NOTHING {

		switch SHARED_ELEMENT_INHERITANCE_MODE {
		case SHARED_ELEMENT_INHERIT_OVERRIDE:
			// Overwrite child with parent shared element values
			if node.Parent.SharedLeft != nil {
				node.SharedLeft = node.Parent.SharedLeft
			}
			if node.Parent.SharedRight != nil {
				node.SharedRight = node.Parent.SharedRight
			}
		case SHARED_ELEMENT_INHERIT_APPEND:
			if node.Parent.SharedLeft != nil && node.SharedLeft != nil {
				// Append child to parent values and assign to child
				node.SharedLeft = append(node.Parent.SharedLeft,node.SharedLeft...)
			} else if node.Parent.SharedLeft != nil {
				//if child is empty, just overwrite
				node.SharedLeft = node.Parent.SharedLeft
			}
			if node.Parent.SharedRight != nil && node.SharedRight != nil {
				// Append child to parent values and assign to child
				node.SharedRight = append(node.Parent.SharedRight, node.SharedRight...)
			} else if node.Parent.SharedRight != nil {
				//if child is empty, just overwrite
				node.SharedRight = node.Parent.SharedRight
			}
		}
		fmt.Println("Inherited shared component from parent component in mode " + SHARED_ELEMENT_INHERITANCE_MODE + ": " +
			"Left: " + fmt.Sprint(node.SharedLeft) + ", Right: " + fmt.Sprint(node.SharedRight))
	}
}

/*
Extracts left and right shared elements of a combination (e.g., (left shared (left [AND] right) right shared))
Input: full string, all boundaries and reference to combination for which shared entries are sought
Return: string arrays for left and right side shared components
 */
func extractSharedComponents(input string, boundaries map[int]map[int]tree.Boundaries, level int, index int) ([]string, []string) {

	// shared components can only be on first level onwards, i.e., if combination is on second level
	if len(boundaries) < 2 {
		return nil, nil
	}

	fmt.Println("IDENTIFYING left and right shared entries for " + input[boundaries[level][index].Left:boundaries[level][index].Right])

	sharedLeft := []string{}
	sharedRight := []string{}
	combination := boundaries[level][index]

	idx := 0
	for idx < len(boundaries[level-1]) {

		// shared entry candidate
		entry := boundaries[level-1][idx]


		// Tested entry must not be combination itself and must frame combination sent as input
		if !entry.Complete &&
			(entry.Left < combination.Left &&
			 entry.Right > combination.Right) {

			fmt.Println("Testing left side")

			if entry.Left < combination.Left {

				// parse from left until combination boundary and check for other combinations in between
				idx1 := 0
				
				matcher := diffmatchpatch.New()

				for idx1 < len(boundaries[level])+1 {

					// Don't check for own index
					if index != idx1 {
						elementToTest := boundaries[level][idx1]
						// check whether another element is in between and a combination
						if elementToTest.Left > entry.Left &&
							elementToTest.Right < combination.Left &&
							elementToTest.Complete {
							// then remove overlapping string in between

							diff := matcher.DiffMain(input[entry.Left:combination.Left-1], input[elementToTest.Left-1:elementToTest.Right+1], false)

							fmt.Print("Found overlaps: ")
							fmt.Println(diff)

							fidx := 0
							for fidx < len(diff) {
								if diff[fidx].Type.String() == "Delete" {
									sharedLeft = append(sharedLeft, diff[fidx].Text)
								}
								fidx++
							}
						}
					}
					idx1++
				}

				//Left shared content here
				if len(sharedLeft) != 0 {
					ct := 0
					for ct < len(sharedLeft) {
						sharedLeft[ct] = strings.Trim(sharedLeft[ct], " ")
						ct++
					}
				} else {
					sharedLeft = append(sharedLeft, strings.Trim(input[entry.Left:combination.Left-1], " "))
				}
				fmt.Println("Found left shared content: " + fmt.Sprint(sharedLeft))
				// Set flag to signal that this entry is shared for later processing of non-shared content
				entry.Shared = true
				boundaries[level-1][idx] = entry
			}

			fmt.Println("Testing right side")

			if entry.Right > combination.Right+1 {

				// parse from left until combination boundary and check for other combinations in between
				idx1 := 0
				
				matcher := diffmatchpatch.New()

				elementsToExclude := []tree.Boundaries{}

				for idx1 < len(boundaries[level])+1 {
					fmt.Println("Index: " + strconv.Itoa(idx1))
					// Don't check for own index
					if index != idx1 {
						fmt.Print("Testing element ")
						fmt.Println(boundaries[level][idx1])
						elementToTest := boundaries[level][idx1]
						// check whether another element is in between and a combination
						if elementToTest.Right < entry.Right &&
							elementToTest.Left > combination.Right &&
							elementToTest.Complete {
							// collect elements that have overlaps
							elementsToExclude = append(elementsToExclude, elementToTest)
							fmt.Print("Planning to exclude element: ")
							fmt.Println(elementToTest)
						}
					}
					idx1++
				}

				if len(elementsToExclude) > 0 {
					// Perform actual filtering
					diff := matcher.DiffMain(input[combination.Right+1:entry.Right],
						input[elementsToExclude[0].Left-1:elementsToExclude[len(elementsToExclude)-1].Right+1], false)
					fmt.Print("Found overlaps: ")
					fmt.Println(diff)

					fidx := 0
					for fidx < len(diff) {
						if diff[fidx].Type.String() == "Delete" {
							sharedRight = append(sharedRight, diff[fidx].Text)
						}
						fidx++
					}
				}

				//Right shared content here
				if len(sharedRight) != 0 {
					ct := 0
					for ct < len(sharedLeft) {
						sharedRight[ct] = strings.Trim(sharedRight[ct], " ")
						ct++
					}
				} else {
					sharedRight = append(sharedRight, strings.Trim(input[combination.Right+1:entry.Right], " "))
				}

				//sharedRight = strings.Trim(input[entry.Left:entry.Right], " ")
				//sharedRight = strings.Trim(input[combination.Right+1:entry.Right], " ")
				fmt.Println("Found right shared content: " + fmt.Sprint(sharedRight))
				// Set flag to signal that this entry is shared for later processing of non-shared content
				entry.Shared = true
				boundaries[level-1][idx] = entry
			}
			return sharedLeft, sharedRight
		}
		idx++
	}
	// no shared content found
	return nil, nil
}