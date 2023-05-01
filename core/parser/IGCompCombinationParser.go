package parser

import (
	"IG-Parser/core/tree"
	"fmt"
	"log"
	"strconv"
	"strings"
)

/*
Parses combinations in string. The syntactic form of input is:
"( leftSide [OPERATOR] rightSide )", where [OPERATOR] is one
of the logical operators [AND], [OR], [XOR] (including brackets),
and left and right side are either text or combinations themselves.
For all logical operators, an arbitrary number of expressions can be combined;
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
func ParseIntoNodeTree(input string, nestedNode bool, leftPar string, rightPar string) (*tree.Node, string, tree.ParsingError) {

	// Check for parentheses
	if leftPar == "" || rightPar == "" {
		return nil, "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_PARENTHESES_COMBINATION, ErrorMessage: "Missing parentheses specification when parsing into tree."}
	}

	if leftPar == LEFT_BRACE && rightPar != RIGHT_BRACE {
		return nil, "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_PARENTHESES_COMBINATION,
			ErrorMessage: "Invalid parentheses specification when parsing into tree (Left: " + leftPar + ", Right: " + rightPar + ")"}
	}

	if leftPar == LEFT_PARENTHESIS && rightPar != RIGHT_PARENTHESIS {
		return nil, "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_PARENTHESES_COMBINATION,
			ErrorMessage: "Invalid parentheses specification when parsing into tree (Left: " + leftPar + ", Right: " + rightPar + ")"}
	}

	if leftPar != LEFT_PARENTHESIS && leftPar != LEFT_BRACE {
		return nil, "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_PARENTHESES_COMBINATION,
			ErrorMessage: "Invalid parentheses specification when parsing into tree (Left: " + leftPar + ", Right: " + rightPar + ")"}
	}

	// Test content for absence of logical operators for non-component-level nested input
	if leftPar != LEFT_BRACE && rightPar != RIGHT_BRACE &&
		!strings.Contains(input, tree.AND_BRACKETS) &&
		!strings.Contains(input, tree.XOR_BRACKETS) &&
		!strings.Contains(input, tree.OR_BRACKETS) &&
		!strings.Contains(input, tree.SAND_BETWEEN_COMPONENTS_BRACKETS) {

		node := &tree.Node{Entry: strings.Trim(input, " ")}
		return node, input, tree.ParsingError{ErrorCode: tree.PARSING_NO_COMBINATIONS,
			ErrorMessage: "The input does not contain combinations"}
	}

	// Return detected combinations, alongside potentially modified input string
	combinations, _, input, err := detectCombinations(input, leftPar, rightPar)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, input, err
	}

	// If no valid combinations are left, the parsing is finished
	if len(combinations) == 0 {
		node := &tree.Node{Entry: strings.Trim(input, " ")}
		return node, strings.Trim(input, " "), tree.ParsingError{ErrorCode: tree.PARSING_NO_COMBINATIONS,
			ErrorMessage: "The input does not contain combinations"}
	}

	// Now the parsing of nested combinations starts
	Print("STARTING TREE CONSTRUCTION: Detected combinations: ")
	Println(combinations)

	// Depth first (level is 1-based, index is 0-based)
	level := 0 //len(combinations)

	// Link input node to temporary node for incremental addition of elements - needs to be copied back prior to return
	//if nodeTree == nil {
	nodeTree := &tree.Node{}
	//}
	// Map to keep order of entry for correct tree construction (retention of order)
	orderMap := make(map[int]*tree.Node)

	stop := false
	// Iterate through all levels
	for level <= len(combinations) {
		Println("Level " + strconv.Itoa(level))

		if stop {
			// Breaking out, since the higher-level combinations must necessarily include higher-level combinations during deep parsing
			Println("Breaking out after first level that contains valid combinations (stopping at Level " + strconv.Itoa(level) + ")")
			break
		}

		Print("Combinations on level " + strconv.Itoa(level) + ": ")
		Println(combinations[level], "- Count: ", len(combinations[level]))

		idx := 0
		//Iterate through all indices
		for idx < len(combinations[level]) {

			//Println("TREE BEFORE NEXT COMBINATION: " + node.String())
			Println("Parsing index " + strconv.Itoa(idx) + " on level " + strconv.Itoa(level))

			// Create node for assigning elements
			node := tree.Node{}

			// Parse complete combinations
			if combinations[level][idx].Complete {

				// Toggling break after the first complete combination on this level,
				// since all combinations at higher levels must be captured based on deep parsing.
				// If parsed explicitly again, they would appear as redundant first-order nodes.
				if !stop {
					stop = true
					Println("Signalling break out after level " + strconv.Itoa(level))
				}

				Println("Found combination to parse on level " + strconv.Itoa(level) +
					", Index: " + strconv.Itoa(idx) + ": " + fmt.Sprint(combinations[level][idx]))

				// Check for shared elements by sending ID to parse function and search boundaries for next lower level and embracing
				sharedLeft, sharedRight := extractSharedComponents(input, combinations, level, idx)

				// Assign left shared value if existing
				if sharedLeft != nil {
					Println("Assigning left shared element '" + fmt.Sprint(sharedLeft) + "'.")
					node.SharedLeft = sharedLeft
					// Reset shared
					sharedLeft = []string{}
				}
				// Assign shared value if existing
				if sharedRight != nil {
					Println("Assigning right shared element '" + fmt.Sprint(sharedRight) + "'.")
					node.SharedRight = sharedRight
					// Reset shared
					sharedRight = []string{}
				}

				Println("Input to parse over: " + input)
				// full parsing
				left := input[combinations[level][idx].Left:combinations[level][idx].Operator]
				right := input[combinations[level][idx].Operator+len(combinations[level][idx].OperatorVal)+2 : combinations[level][idx].Right]
				Println("==Raw Left value: " + left)
				Println("==Raw Operator: " + combinations[level][idx].OperatorVal)
				Println("==Raw Right value: " + right)

				// Assign logical operator
				node.LogicalOperator = combinations[level][idx].OperatorVal

				Println("Tree before deep parsing: " + node.String())

				// Left side (potentially modifying input string)
				leftCombos, leftNonShared, left, err := detectCombinations(left, leftPar, rightPar)
				if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS_DURING_NODE_PARSING {
					log.Println("Error when parsing left side: " + err.ErrorMessage)
					return &node, input, err
				}
				if err.ErrorCode == tree.PARSING_ERROR_IGNORED_ELEMENTS_DURING_NODE_PARSING {
					log.Print("Warning: Discarded elements during deep left parsing: " + tree.PrintArray(err.ErrorIgnoredElements))
				}

				if len(leftCombos) == 0 {
					// Trim content first
					left = strings.Trim(left, " ")
					if left != "" {
						// If no combinations exist, assign as left leaf
						Println("Found leaf on left side: " + left)
						res, err := node.InsertLeftLeaf(left)
						if !res {
							return nil, input, tree.ParsingError{ErrorCode: err.ErrorCode, ErrorMessage: err.ErrorMessage}
						}
					} else {
						msg := "Empty leaf value on left side: " + left +
							" (Corresponding right value and operator: " + right + "; " + node.LogicalOperator +
							");\n processed expression: " + input
						log.Println(msg)
						return nil, input, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_EMPTY_LEAF,
							ErrorMessage: msg}
					}
				} else {
					// If combinations exist (may not be complete), delegate

					// Deal with non-shared items
					if len(leftNonShared) > 0 {
						Println("NOT HANDLED - JUST INFORMATION: Non-shared elements left: " + fmt.Sprint(leftNonShared))
						// TODO: Assign the left node - order not managed
					}

					Println("Go deep on left side: " + left)
					leftNode, left, err := ParseIntoNodeTree(left, true, leftPar, rightPar)
					if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
						return nil, left, err
					}

					// Check whether returned node is empty
					if leftNode.IsEmptyOrNilNode() {
						// If so, assign original input value as leaf entry
						Println("Deep parsing: Returned node is empty; assigning complete value as left leaf")
						left = strings.Trim(left, " ")
						if left != "" {
							// If no combinations exist, assign as left leaf
							Println("Found leaf on left side: " + left)
							res, err := node.InsertLeftLeaf(left)
							if !res {
								return nil, input, tree.ParsingError{ErrorCode: err.ErrorCode, ErrorMessage: err.ErrorMessage}
							}
						}
					} else {
						// Create nested node and link with parent
						Println("Deep parsing: Returned left node is combination; assign as nested node. Node:", leftNode)

						// Link newly identified node with main node
						res, err := node.InsertLeftNode(leftNode)
						if !res {
							return nil, input, tree.ParsingError{ErrorCode: err.ErrorCode, ErrorMessage: err.ErrorMessage}
						}

						// Check for inheriting shared elements on AND nodes
						//inheritSharedElements(leftNode)
					}

					Println("Tree after processing left deep: " + node.String())
				}

				// Right side (potentially modifying input string)
				rightCombos, rightNonShared, right, err := detectCombinations(right, leftPar, rightPar)
				if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS_DURING_NODE_PARSING {
					log.Println("Error when parsing right side: " + err.ErrorMessage)
					return &node, input, err
				}
				if err.ErrorCode == tree.PARSING_ERROR_IGNORED_ELEMENTS_DURING_NODE_PARSING {
					log.Print("Warning: Discarded elements during deep right parsing: " + tree.PrintArray(err.ErrorIgnoredElements))
				}

				if len(rightCombos) == 0 {
					// Trim content first
					right = strings.Trim(right, " ")
					if right != "" {
						// If no combinations exist, assign as right leaf
						Println("Found leaf on right side: " + right)
						res, err := node.InsertRightLeaf(right)
						if !res {
							return nil, input, tree.ParsingError{ErrorCode: err.ErrorCode, ErrorMessage: err.ErrorMessage}
						}
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
						Println("NOT HANDLED - JUST INFORMATION: Non-shared elements right: " + fmt.Sprint(rightNonShared))
						// TODO: Assign to right node - order not managed
					}

					Println("Go deep on right side: " + right)
					rightNode, right, err := ParseIntoNodeTree(right, true, leftPar, rightPar)
					if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
						return nil, right, err
					}

					// Check whether returned node is empty
					if rightNode.IsEmptyOrNilNode() {
						// If so, assign original input value as leaf entry
						Println("Deep parsing: Returned node is empty; assigning complete value as right leaf")
						right = strings.Trim(right, " ")
						if right != "" {
							// If no combinations exist, assign as right leaf
							Println("Found leaf on right side: " + right)
							res, err := node.InsertRightLeaf(right)
							if !res {
								return nil, input, tree.ParsingError{ErrorCode: err.ErrorCode, ErrorMessage: err.ErrorMessage}
							}
						}
					} else {
						// Create nested node and link with parent
						Println("Deep parsing: Returned right node is combination; assign as nested node")

						// Link newly identified node with main node
						res, err := node.InsertRightNode(rightNode)
						if !res {
							return nil, input, tree.ParsingError{ErrorCode: err.ErrorCode, ErrorMessage: err.ErrorMessage}
						}

						// Check for inheriting shared elements on AND nodes
						//inheritSharedElements(rightNode)
					}

					Println("Tree after processing right deep: " + node.String())
				}

				// Store node for further preparation of return
				if !nestedNode || len(combinations[level]) > 1 {
					// Either adds non-nested, i.e. top-level, node to map with starting character index as key - for later reconstruction of tree prior to return
					// or element for nested combination that is yet to be completed and integrated into first-order expression (e.g., based on other associated statement)
					// The latter commonly applies when facing multiple implicitly linked combinations within side of combination (e.g., left, right).
					Println("Saving tree to ordermap (Level "+strconv.Itoa(level)+", Index: "+strconv.Itoa(idx)+"): ", node)
					orderMap[combinations[level][idx].Left] = &node
				} else {
					Println("Node is nested; return without further linking")
					// return node outright if nested
					return &node, input, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
				}
			}

			Println("==Finished parsing index " + strconv.Itoa(idx) + " on level " + strconv.Itoa(level))
			// Increase to explore other entries
			idx++
		}
		// break out if combination has been found and processed - must be top-level combination
		// TODO - check output
		//return node, "(" + left + " [" + node.LogicalOperator + "] " + right + ")", tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		//return node, "", tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		Println("==Finished parsing on level " + strconv.Itoa(level))

		// Increase level
		level++

	}

	Print("Combinations before return: ")
	Println(combinations)

	// Reconstructing node based on order of node input
	//Println("Node order: ")
	ct := 0

	// Default error to test for invalid combinations
	nodeCombinationError := tree.NodeError{ErrorCode: tree.PARSING_NO_ERROR}

	// If more than one node ...
	if len(orderMap) > 1 {
		for ct < len(input) {
			if _, ok := orderMap[ct]; ok {
				Println("Final tree before adding element:", ct, ":", nodeTree.String())
				// ... then synthetically link elements
				nodeTree, nodeCombinationError = tree.Combine(nodeTree, orderMap[ct], tree.SAND_WITHIN_COMPONENTS)
				// Check if combination error has been picked up - here and in the beginning of loop
				if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
					return nil, input, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
						ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
				}
				Println("Added to tree: " + fmt.Sprint(orderMap[ct]))
				Println("Final tree after adding element:", ct, ":", nodeTree.String())
			}
			ct++
		}
	} else if len(orderMap) == 1 {
		// or simply assign last node if only one has been found
		for _, v := range orderMap {
			nodeTree = v
		}
		Println("Simple assignment of single node...")
	}

	Println("RETURNING FINAL NODE: " + nodeTree.String())
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
func detectCombinations(expression string, leftPar string, rightPar string) (map[int]map[int]tree.Boundaries, []string, string, tree.ParsingError) {

	// Check for parentheses
	if leftPar == "" || rightPar == "" {
		return nil, nil, "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_PARENTHESES_COMBINATION, ErrorMessage: "Missing parentheses specification for detection of combinations."}
	}

	if leftPar == LEFT_BRACE && rightPar != RIGHT_BRACE {
		return nil, nil, "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_PARENTHESES_COMBINATION,
			ErrorMessage: "Invalid parentheses specification when parsing into tree (Left: " + leftPar + ", Right: " + rightPar + ")"}
	}

	if leftPar == LEFT_PARENTHESIS && rightPar != RIGHT_PARENTHESIS {
		return nil, nil, "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_PARENTHESES_COMBINATION,
			ErrorMessage: "Invalid parentheses specification when parsing into tree (Left: " + leftPar + ", Right: " + rightPar + ")"}
	}

	if leftPar != LEFT_PARENTHESIS && leftPar != LEFT_BRACE {
		return nil, nil, "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_PARENTHESES_COMBINATION,
			ErrorMessage: "Invalid parentheses specification when parsing into tree (Left: " + leftPar + ", Right: " + rightPar + ")"}
	}

	// Tracks current parsing level
	level := 0

	// Parentheses count to check for balance
	parCount := 0

	// Initial test run for parentheses
	for i, letter := range expression {

		switch string(letter) {
		case leftPar:
			parCount++
		case rightPar:
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

	// Parsing-independent parentheses (i.e., (, )) count for assessment parsing mode when parsing nested statements
	// Values greater 0 indicate that parsing is in component-level nesting territory; important to differentiate scope
	// for nested statement parsing (i.e., {stmt1 [AND] stmt2}, vs. {A(...) I(first [AND] second) Cac()} [AND] ...
	// to avoid picking up nested logical operator within component-level nesting scope
	generalParCount := 0

	Println("Testing expression " + expression)
	for i, letter := range expression {

		// Register parenthesis count independent from main parsing
		switch string(letter) {
		case LEFT_PARENTHESIS:
			// Increase (to signal operation inside component-level combinations)
			generalParCount++
		case RIGHT_PARENTHESIS:
			// Reduce
			generalParCount--
		}

		switch string(letter) {
		case leftPar:
			// Increase level
			level++
			Println("Expression start detected (Level " + strconv.Itoa(level) + ")")
			// Configure mode
			modeMap[level] = tree.PARSING_MODE_LEFT
			Println("Mode: " + modeMap[level])
			//Test of existing entries
			if _, ok := levelMap[level]; ok {
				Println(strconv.Itoa(len(levelMap[level])) + " key(s) already defined - increasing count")
				levelMap[level][len(levelMap[level])] = tree.Boundaries{Left: i + 1, Complete: false}
			} else {
				// Store index reference and create internal map (incremented with 1 to avoid left parenthesis - balanced)
				mp := make(map[int]tree.Boundaries)
				levelMap[level] = mp
				levelMap[level][0] = tree.Boundaries{Left: i + 1, Complete: false}
			}
			// Count parentheses to detect uneven matching
			parCount++
		case rightPar:
			Println("Expression end detected (Level " + strconv.Itoa(level) + ")")
			// Check if there are repetitions
			for k, _ := range foundOperators { // key: operator, value:
				for k2, v2 := range foundOperators[k] { // key: level, value: number of occurrences
					if v2 > 1 {
						log.Println("Found " + strconv.Itoa(v2) + " occurrences of operator " + k + " on level " + strconv.Itoa(k2) + " in map.")
					}
				}

			}
			// Store index reference
			levelIdx := len(levelMap[level]) - 1
			Println("Level before saving: " + strconv.Itoa(level) + "; Index: " + strconv.Itoa(levelIdx))
			b := levelMap[level][levelIdx]
			b.Right = i
			levelMap[level][levelIdx] = b
			Println("Level map for level " + strconv.Itoa(level) + " (after adding right value but before assessing completeness): " + b.String())
			Println("--> Content:", expression[b.Left:b.Right])
			// Test whether indices are identical or immediately following - suggesting gaps in values
			if (b.Operator + len(b.OperatorVal) + 2) == b.Right {
				msg := "Input contains invalid combination expression in the range '" + expression[b.Left:b.Right] + "'."
				return levelMap, nonSharedElements, expression, tree.ParsingError{ErrorCode: tree.PARSING_INVALID_COMBINATION,
					ErrorMessage: msg}
			}
			if b.Left != 0 && b.OperatorVal != "" {
				// Update complete marker
				//b := levelMap[level][levelIdx]
				b.Complete = true
				levelMap[level][levelIdx] = b
				Println("Expression is complete.")
			} else {
				Println("Detected end, but combination incomplete (Missing operator or left parenthesis). " +
					"Discarding combination. (Input: '" + expression + "')")
				// Delete incomplete entry - will lead to removal of shared strings...
				/*delete(levelMap[level], levelIdx)
				// Check whether all information from that level has to be deleted
				if len(levelMap[level]) == 0 {
					delete(levelMap, level)
				}*/
				Print("Map on given level after processing: ")
				Println(levelMap[level])
				// Retrieve higher level
				/*_, ok := levelMap[level+1]
				if !ok {
					Print("No nested complete combination, so deleted this incomplete one: ")
					Println(levelMap[level])
					// If no higher-level exists (within the lower level), delete this incomplete entry
					delete(levelMap, level)
				} else {
					// else retain
					Println("Nested complete combination, so incomplete higher-level combination is retained.")
				}*/
			}

			// Configure mode
			modeMap[level] = tree.PARSING_MODE_OUTSIDE_EXPRESSION
			Println("Mode: " + modeMap[level])

			// Reset operator count for given level
			for op := range foundOperators {
				Println("Resetting operator " + op + " for level " + strconv.Itoa(level))
				delete(foundOperators[op], level)
			}

			// Reduce level
			level--
			Println("Moving back to level " + strconv.Itoa(level) + ", Mode: " + modeMap[level])
			// Count parentheses to detect uneven matching
			parCount--
		case "[":
			//Println("Checking for logical operator ... " + expression[i:i+5])
			foundOperator := ""
			// Check for length of expression before testing logical operators
			if len(expression) >= i+len(tree.AND_BRACKETS) {
				switch expression[i : i+len(tree.AND_BRACKETS)] {
				case tree.AND_BRACKETS:
					Println("Detected " + tree.AND_BRACKETS)
					foundOperator = tree.AND
				case tree.XOR_BRACKETS:
					Println("Detected " + tree.XOR_BRACKETS)
					foundOperator = tree.XOR
				}
			}
			// Separately test for OR due to differing length (but remember to test for length of expression first)
			if foundOperator == "" && len(expression) >= i+len(tree.OR_BRACKETS) &&
				expression[i:i+len(tree.OR_BRACKETS)] == tree.OR_BRACKETS {
				Println("Detected " + tree.OR_BRACKETS)
				foundOperator = tree.OR
			}
			// Separately test for sAND WITHIN components due to differing length (but remember to test for length of expression first)
			if foundOperator == "" && len(expression) >= i+len(tree.SAND_WITHIN_COMPONENTS_BRACKETS) &&
				expression[i:i+len(tree.SAND_WITHIN_COMPONENTS_BRACKETS)] == tree.SAND_WITHIN_COMPONENTS_BRACKETS {
				Println("Detected " + tree.SAND_WITHIN_COMPONENTS_BRACKETS)
				foundOperator = tree.SAND_WITHIN_COMPONENTS
			}
			// Separately test for sAND BETWEEN components due to differing length (but remember to test for length of expression first)
			if foundOperator == "" && len(expression) >= i+len(tree.SAND_BETWEEN_COMPONENTS_BRACKETS) &&
				expression[i:i+len(tree.SAND_BETWEEN_COMPONENTS_BRACKETS)] == tree.SAND_BETWEEN_COMPONENTS_BRACKETS {
				Println("Detected " + tree.SAND_BETWEEN_COMPONENTS_BRACKETS)
				foundOperator = tree.SAND_BETWEEN_COMPONENTS
			}
			// If parsing for statement combinations, suppress operators if within component-level nesting scope
			if foundOperator != "" && leftPar == LEFT_BRACE && rightPar == RIGHT_BRACE && generalParCount != 0 {
				// Suppress registration of combination
				foundOperator = ""
				Println("Statement-level parsing: Suppressing nested logical operator " +
					foundOperator + " during statement-level parsing")
			}

			// Perform proper handling of operator
			if foundOperator != "" {

				Println("Found logical operator " + foundOperator + " on level " + strconv.Itoa(level))

				if modeMap[level] == tree.PARSING_MODE_OUTSIDE_EXPRESSION {
					return nil, nil, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_LOGICAL_OPERATOR_OUTSIDE_COMBINATION,
						ErrorMessage: "Logical operator (e.g., [AND], [OR], [XOR]) found outside of combination. Please check for missing parentheses in input."}
				}

				levelIdx := len(levelMap[level]) - 1
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
						Println(" -> Added. Count: " + strconv.Itoa(foundOperators[foundOperator][level]))
					} else {
						// else create new level entry with default value of 1
						foundOperators[foundOperator][level] = 1
						Println(" -> Created. Count: " + strconv.Itoa(foundOperators[foundOperator][level]))
					}
				} else {
					// if no operator entry exists, else create new operator entry with default value of 1
					foundOperators[foundOperator] = make(map[int]int)
					foundOperators[foundOperator][level] = 1
					Println(" -> Created level and value. Count: " + strconv.Itoa(foundOperators[foundOperator][level]))
				}

				// If already in right parsing mode, there should be no operator
				if modeMap[level] == tree.PARSING_MODE_RIGHT {
					log.Println("Found additional operator [" + foundOperator + "] (now " + strconv.Itoa(foundOperators[foundOperator][level]) +
						" times on level " + strconv.Itoa(level) + "), even though looking for terminating parenthesis.")
					if foundOperators[foundOperator][level] > 1 { // if AND operator and multiple on the same level
						// Consider injecting a left parenthesis before the expression and add mixfix ") " before logical operator, e.g., "( left ... [AND] right ... ) [AND] ..."
						expression = expression[:levelMap[level][levelIdx].Left] + leftPar + expression[levelMap[level][levelIdx].Left:i-1] + rightPar + " " + expression[i:]
						log.Println("Multiple [AND], [OR], or [XOR] operators found. Reconstructed nested structure by introducing parentheses, now: " + expression)
						log.Println("Rerunning all parsing on combination to capture nested AND, OR, XOR combinations")
						return detectCombinations(expression, leftPar, rightPar)
					} else {
						return levelMap, nonSharedElements, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS,
							ErrorMessage: "Error: Mix of different logical operators on level " + strconv.Itoa(level) +
								" in single expression (Expression: " + expression + "). Use parentheses to indicate precedence."}
					}
				}

				// Configure mode
				modeMap[level] = tree.PARSING_MODE_RIGHT
				Println("Mode: " + modeMap[level])

				// Store index reference and value
				o := levelMap[level][levelIdx]
				o.Operator = i
				o.OperatorVal = foundOperator
				levelMap[level][levelIdx] = o
			}
		default:
			// Note that this is called for every non-combination element; requires higher-level filtering for leaf entries
			if modeMap[level] == tree.PARSING_MODE_OUTSIDE_EXPRESSION {
				// Copy cached string to array if disruption in continuation in expression (e.g., combinations in between)
				if lastCharIndex-i > 1 && nonSharedElement != "" {
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
				Println("Prefix pos: " + strconv.Itoa(firstIdx))
			}
			if levelMap[i].Right > lastIdx {
				// Assign highest last index
				lastIdx = levelMap[i].Right
				Println("Suffix pos: " + strconv.Itoa(lastIdx))
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
		Println("Prefix: " + prefix)
		Println("Suffix: " + suffix)
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
		Println("Returning expression (ignored elements: " + errorString + "): " + expression)
		return levelMap, expression, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IGNORED_ELEMENTS_DURING_NODE_PARSING,
			ErrorMessage: "Parsing was successful, but expression parts were ignored during coding (" + errorString + "). " +
			"This commonly occurs when logical operators between simple strings and combinations are omitted " +
			"(e.g., ... some string (left [AND] right) ...) and not wrapped by parentheses to signal shared elements. " +
			"In this case, simple strings are ignored in the parsing process.",
			ErrorIgnoredElements: ignoredElements}
	}*/

	Println("Returning expression (complete parsing): " + expression)
	// if no omitted elements during parsing, regular return without error
	return levelMap, nonSharedElements, expression, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Processes potential inheritance of shared element values from parent to child nodes,
where both parent and child nodes have AND operators
*/
/*
func inheritSharedElements(node *tree.Node) {
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
		Println("Inherited shared component from parent component in mode " + SHARED_ELEMENT_INHERITANCE_MODE + ": " +
			"Left: " + fmt.Sprint(node.SharedLeft) + ", Right: " + fmt.Sprint(node.SharedRight))
	}
}
*/

/*
Symbols to be ignored when determining strings shared across combinations (e.g., 'shared' in 'Bdir(shared (left [AND] right))') - commonly parentheses and braces
*/
const ignoredSymbolsInSharedFields = "(){}"

/*
Extracts left and right shared elements of a combination (e.g., (left shared (left [AND] right) right shared))
Input: full string, all boundaries and reference to combination for which shared entries are sought
Return: string arrays for left and right side shared components
*/
func extractSharedComponents(input string, boundaries map[int]map[int]tree.Boundaries, level int, index int) ([]string, []string) {

	Println("Identifying left and right shared entries for " + input[boundaries[level][index].Left:boundaries[level][index].Right])

	// Output arrays
	sharedLeft := []string{}
	sharedRight := []string{}

	Println("=Boundaries:", boundaries)
	Println("=Input:", input)
	Println("=Level:", level)
	Println("=Index:", index)

	if boundaries[level][index].OperatorVal != "" {

		// LEFT SIDE
		if index == 0 {
			// Take everything on the left side as shared left element, unless ...
			val := input[:boundaries[level][index].Left]
			// ... lower-level boundary exists, in which case that one is used
			outerBoundary := -1
			if level > 0 && boundaries[level-1] != nil {
				for i, v := range boundaries[level-1] {
					// If lower-level boundary encompasses current one, then consider for extended boundaries
					if v.Left < boundaries[level][index].Left && v.Right > boundaries[level][index].Right && v.OperatorVal == "" {
						outerBoundary = i
						break
					}
				}
			}
			if outerBoundary != -1 {
				// If multiple levels, though, only consume up to left boundary on next lower level (and if that level does not
				// contain a valid combination)
				val = input[boundaries[level-1][outerBoundary].Left:boundaries[level][index].Left]
			}

			// Get rid of parentheses first
			val = strings.Trim(val, ignoredSymbolsInSharedFields)
			// Then get rid of spaces
			val = strings.TrimSpace(val)
			Println("Identified left shared value:", val)
			if val != "" {
				sharedLeft = append(sharedLeft, val)
			}
		} else {
			// If element is left (e.g., wAND combined on same level), only consider string to boundary as shared
			val := input[boundaries[level][index-1].Right:boundaries[level][index].Left]
			// Get rid of parentheses first
			val = strings.Trim(val, ignoredSymbolsInSharedFields)
			// Then get rid of spaces
			val = strings.TrimSpace(val)
			Println("Identified left shared value in multi-value component:", val)
			if val != "" {
				sharedLeft = append(sharedLeft, val)
			}
		}

		// RIGHT SIDE
		// Check if other combination (i.e., higher index) exists, and only consider elements to that boundary as shared right
		if value, ok := boundaries[level][index+1]; ok {
			val := input[boundaries[level][index].Right:value.Left]
			// Get rid of parentheses first
			val = strings.Trim(val, ignoredSymbolsInSharedFields)
			// Then get rid of spaces
			val = strings.TrimSpace(val)
			Println("Identified right shared value in multi-value component:", val)
			if val != "" {
				sharedRight = append(sharedRight, val)
			}
		} else {
			// Else take the entire remaining string on the right, unless ...
			val := input[boundaries[level][index].Right:]
			outerBoundary := -1
			// ... lower-level boundary exists, but also which index (in case of multiple combinations)
			// wraps the currently analyzed entry
			if level > 0 && boundaries[level-1] != nil {
				for i, v := range boundaries[level-1] {
					// If lower-level boundary encompasses current one, then consider for extended boundaries
					if v.Left < boundaries[level][index].Left && v.Right > boundaries[level][index].Right && v.OperatorVal == "" {
						outerBoundary = i
						break
					}
				}
			}
			if outerBoundary != -1 {
				// If multiple levels, though, only consume up to right boundary on next lower level (and if that level does not
				// contain a valid combination)
				val = input[boundaries[level][index].Right:boundaries[level-1][outerBoundary].Right]
			}
			// Get rid of parentheses first
			val = strings.Trim(val, ignoredSymbolsInSharedFields)
			// Then get rid of spaces
			val = strings.TrimSpace(val)
			Println("Identified right shared value:", val)
			if val != "" {
				sharedRight = append(sharedRight, val)
			}
		}
	}
	// Return nil if nothing found
	if len(sharedLeft) == 0 && len(sharedRight) == 0 {
		return nil, nil
	}
	if len(sharedLeft) == 0 && len(sharedRight) != 0 {
		return nil, sharedRight
	}
	if len(sharedLeft) != 0 && len(sharedRight) == 0 {
		return sharedLeft, nil
	}
	return sharedLeft, sharedRight
}
