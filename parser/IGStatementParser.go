package parser

import (
	"IG-Parser/tree"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)


func ParseStatement(text string) (tree.Statement, tree.ParsingError) {

	// Remove line breaks
	text = cleanInput(text)

	s := tree.Statement{}

	// Validate input string first with respect to parentheses ...
	err := validateInput(text, LEFT_PARENTHESIS, RIGHT_PARENTHESIS)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return tree.Statement{}, err
	}
	// ... and braces
	err = validateInput(text, LEFT_BRACE, RIGHT_BRACE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return tree.Statement{}, err
	}

	// Now retrieve component-only and nested statements
	compAndNestedStmts, err := separateComponentsAndNestedStatements(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return tree.Statement{}, err
	}
	fmt.Println("Returned separated components and nested statements: " + fmt.Sprint(compAndNestedStmts))
	// Extract component-only statement and override input
	text = compAndNestedStmts[0][0]
	// Extract potential nested statements
	nestedStmts := compAndNestedStmts[1]
	if len(nestedStmts) == 0 {
		log.Println("No nested statements found.")
	}
	nestedCombos := compAndNestedStmts[2]
	if len(nestedCombos) == 0 {
		log.Println("No nested statement combination candidates found.")
	}

	fmt.Println("Text to be parsed: " + text)
	// Now parsing on component level

	result, err := parseAttributes(text)
	outErr := handleParsingError(tree.ATTRIBUTES, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.Attributes = result

	result, err = parseAttributesProperty(text)
	outErr = handleParsingError(tree.ATTRIBUTES_PROPERTY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.AttributesPropertySimple = result

	result, err = parseDeontic(text)
	outErr = handleParsingError(tree.DEONTIC, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.Deontic = result

	result, err = parseAim(text)
	outErr = handleParsingError(tree.AIM, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.Aim = result

	result, err = parseDirectObject(text)
	outErr = handleParsingError(tree.DIRECT_OBJECT, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.DirectObject = result

	result, err = parseDirectObjectProperty(text)
	outErr = handleParsingError(tree.DIRECT_OBJECT_PROPERTY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.DirectObjectPropertySimple = result

	result, err = parseIndirectObject(text)
	outErr = handleParsingError(tree.INDIRECT_OBJECT, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.IndirectObject = result

	result, err = parseIndirectObjectProperty(text)
	outErr = handleParsingError(tree.INDIRECT_OBJECT_PROPERTY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.IndirectObjectPropertySimple = result

	result, err = parseActivationCondition(text)
	outErr = handleParsingError(tree.ACTIVATION_CONDITION, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.ActivationConditionSimple = result

	result, err = parseExecutionConstraint(text)
	outErr = handleParsingError(tree.EXECUTION_CONSTRAINT, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.ExecutionConstraintSimple = result

	result, err = parseConstitutedEntity(text)
	outErr = handleParsingError(tree.CONSTITUTED_ENTITY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.ConstitutedEntity = result

	result, err = parseConstitutedEntityProperty(text)
	outErr = handleParsingError(tree.CONSTITUTED_ENTITY_PROPERTY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.ConstitutedEntityPropertySimple = result

	result, err = parseModal(text)
	outErr = handleParsingError(tree.MODAL, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.Modal = result

	result, err = parseConstitutingFunction(text)
	outErr = handleParsingError(tree.CONSTITUTIVE_FUNCTION, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.ConstitutiveFunction = result

	result, err = parseConstitutingProperties(text)
	outErr = handleParsingError(tree.CONSTITUTING_PROPERTIES, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.ConstitutingProperties = result

	result, err = parseConstitutingPropertiesProperty(text)
	outErr = handleParsingError(tree.CONSTITUTING_PROPERTIES_PROPERTY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		return s, outErr
	}
	s.ConstitutingPropertiesPropertySimple = result

	//fmt.Println(s.String())

	fmt.Println("Testing for nested statements in " + fmt.Sprint(nestedStmts))

	// Process nested statements
	if len(nestedStmts) > 0 {
		log.Println("Found nested statements ...")
		err = parseNestedStatements(&s, nestedStmts)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return s, err
		}
	}

	fmt.Println("Testing for nested statement combinations in " + fmt.Sprint(nestedCombos))

	// Process nested statement combinations
	if len(nestedCombos) > 0 {
		log.Println("Found nested statement combinations ...")
		err = parseNestedStatementCombinations(&s, nestedCombos)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return s, err
		}
	}

	log.Println("Statement (after assigning sub elements):\n" + s.String())

	return s, outErr
}

/*
Parses nested statements (but not combinations) and attached those to the top-level statement
 */
func parseNestedStatements(stmtToAttachTo *tree.Statement, nestedStmts []string) (tree.ParsingError) {
	for _, v := range nestedStmts {

		log.Println("Found nested statement")
		// Extract nested statement content and parse
		stmt, errStmt := ParseStatement(v[strings.Index(v, LEFT_BRACE)+1:strings.LastIndex(v, RIGHT_BRACE)])
		if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
			fmt.Print("Error when parsing nested statements: " + errStmt.ErrorCode)
			return errStmt
		}

		// Wrap statement into node (since individual statement)
		stmtNode := tree.Node{Entry: stmt}

		// Checks are ordered with property variants (e.g., Bdir,p) before component variants (e.g., Bdir) to avoid wrong match

		if strings.HasPrefix(v, tree.ATTRIBUTES_PROPERTY) {
			log.Println("Attaching nested attributes property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.AttributesPropertyComplex = attachComplexComponent(stmtToAttachTo.AttributesPropertyComplex, &stmtNode)
			continue
		}
		if strings.HasPrefix(v, tree.DIRECT_OBJECT_PROPERTY) {
			log.Println("Attaching nested direct object property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.DirectObjectPropertyComplex = attachComplexComponent(stmtToAttachTo.DirectObjectPropertyComplex, &stmtNode)
			continue
		}
		if strings.HasPrefix(v, tree.DIRECT_OBJECT) {
			log.Println("Attaching nested direct object to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.DirectObjectComplex = attachComplexComponent(stmtToAttachTo.DirectObjectComplex, &stmtNode)
			continue
		}
		if strings.HasPrefix(v, tree.INDIRECT_OBJECT_PROPERTY) {
			log.Println("Attaching nested indirect object property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.IndirectObjectPropertyComplex = attachComplexComponent(stmtToAttachTo.IndirectObjectPropertyComplex, &stmtNode)
			continue
		}
		if strings.HasPrefix(v, tree.INDIRECT_OBJECT) {
			log.Println("Attaching nested indirect object to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.IndirectObjectComplex = attachComplexComponent(stmtToAttachTo.IndirectObjectComplex, &stmtNode)
			continue
		}
		if strings.HasPrefix(v, tree.ACTIVATION_CONDITION) {
			log.Println("Attaching nested activation condition to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ActivationConditionComplex = attachComplexComponent(stmtToAttachTo.ActivationConditionComplex, &stmtNode)
			continue
		}
		if strings.HasPrefix(v, tree.EXECUTION_CONSTRAINT) {
			log.Println("Attaching nested execution constraint to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ExecutionConstraintComplex = attachComplexComponent(stmtToAttachTo.ExecutionConstraintComplex, &stmtNode)
			continue
		}
		if strings.HasPrefix(v, tree.CONSTITUTED_ENTITY_PROPERTY) {
			log.Println("Attaching nested constituted entity property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutedEntityPropertyComplex = attachComplexComponent(stmtToAttachTo.ConstitutedEntityPropertyComplex, &stmtNode)
			continue
		}
		if strings.HasPrefix(v, tree.CONSTITUTING_PROPERTIES_PROPERTY) {
			log.Println("Attaching nested constituting properties property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutingPropertiesPropertyComplex = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesPropertyComplex, &stmtNode)
			continue
		}
		if strings.HasPrefix(v, tree.CONSTITUTING_PROPERTIES) {
			log.Println("Attaching nested constituting properties to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutingPropertiesComplex = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesComplex, &stmtNode)
			continue
		}
		if strings.HasPrefix(v, tree.OR_ELSE) {
			log.Println("Attaching nested or else to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.OrElse = attachComplexComponent(stmtToAttachTo.OrElse, &stmtNode)
		}
	}
	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Parses nested statement combinations and attached those to the top-level statement
*/
func parseNestedStatementCombinations(stmtToAttachTo *tree.Statement, nestedCombos []string) (tree.ParsingError) {

	for _, v := range nestedCombos {

		log.Println("Found nested statement combination candidate", v)

		combo, _, errStmt := ParseIntoNodeTree(v, false, LEFT_BRACE, RIGHT_BRACE)
		if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
			fmt.Print("Error when parsing nested statements: " + errStmt.ErrorCode)
			return errStmt
		}

		// Check whether all leaves have the same prefix
		flatCombo := tree.Flatten(combo.GetLeafNodes())
		sharedPrefix := ""
		for _, node := range flatCombo {
			entry := node.Entry.(string)
			fmt.Println("Entry to parse of component type: " + entry)
			// Extract prefix for node
			prefix := entry[:strings.Index(entry, LEFT_BRACE)]
			if sharedPrefix == "" {
				// Cache it if not already done
				sharedPrefix = prefix
				continue
			}
			// Check if it deviates from previously cached element
			if prefix != sharedPrefix {
				return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_TYPES_IN_NESTED_STATEMENT_COMBINATION,
					ErrorMessage: "Invalid combination of component-level nested statements. Expected component: " +
						sharedPrefix + ", but found: " + prefix}
			}
		}

		// Parse all entries in tree from string to statement (walks through entire tree linked to node)
		err := combo.ParseAllEntries(func(oldValue string) (tree.Statement, tree.ParsingError) {
			stmt, errStmt := ParseStatement(oldValue[strings.Index(oldValue, LEFT_BRACE)+1:strings.LastIndex(oldValue, RIGHT_BRACE)])
			if errStmt.ErrorCode != tree.PARSING_NO_ERROR{
				return stmt, errStmt
			}
			return stmt, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		})
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return err
		}

		log.Println("Assigning nested tree structure", combo)

		// Checks are ordered with property variants (e.g., Bdir,p) before component variants (e.g., Bdir) to avoid wrong match

		if strings.HasPrefix(sharedPrefix, tree.ATTRIBUTES_PROPERTY) {
			log.Println("Attaching nested attributes property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.AttributesPropertyComplex = attachComplexComponent(stmtToAttachTo.AttributesPropertyComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.DIRECT_OBJECT_PROPERTY) {
			log.Println("Attaching nested direct object property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.DirectObjectPropertyComplex = attachComplexComponent(stmtToAttachTo.DirectObjectPropertyComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.DIRECT_OBJECT) {
			log.Println("Attaching nested direct object to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.DirectObjectComplex = attachComplexComponent(stmtToAttachTo.DirectObjectComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.INDIRECT_OBJECT_PROPERTY) {
			log.Println("Attaching nested indirect object property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.IndirectObjectPropertyComplex = attachComplexComponent(stmtToAttachTo.IndirectObjectPropertyComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.INDIRECT_OBJECT) {
			log.Println("Attaching nested indirect object to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.IndirectObjectComplex = attachComplexComponent(stmtToAttachTo.IndirectObjectComplex, combo)
			continue
		}

		if strings.HasPrefix(sharedPrefix, tree.ACTIVATION_CONDITION) {
			log.Println("Attaching nested activation condition to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ActivationConditionComplex = attachComplexComponent(stmtToAttachTo.ActivationConditionComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.EXECUTION_CONSTRAINT) {
			log.Println("Attaching nested execution constraint to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ExecutionConstraintComplex = attachComplexComponent(stmtToAttachTo.ExecutionConstraintComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.CONSTITUTED_ENTITY_PROPERTY) {
			log.Println("Attaching nested constituted entity property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutedEntityPropertyComplex = attachComplexComponent(stmtToAttachTo.ConstitutedEntityPropertyComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.CONSTITUTING_PROPERTIES_PROPERTY) {
			log.Println("Attaching nested constituting properties property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutingPropertiesPropertyComplex = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesPropertyComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.CONSTITUTING_PROPERTIES) {
			log.Println("Attaching nested constituting properties to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutingPropertiesComplex = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.OR_ELSE) {
			log.Println("Attaching nested or else to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.OrElse = attachComplexComponent(stmtToAttachTo.OrElse, combo)
		}
	}
	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Attach complex component to tree structure under consideration of existing nodes in target tree structure.
Input:
- Node of the parent tree to attach to
- Node to attach
 */
func attachComplexComponent(nodeToAttachTo *tree.Node, nodeToAttach *tree.Node) *tree.Node {
	log.Println("Attaching nested activation condition to higher-level statement")
	// Assign nested statement to higher-level statement

	// If already a statement assignment to complex element, ...
	if nodeToAttachTo != nil {
		// ... combine both
		nodeToAttachTo = tree.Combine(nodeToAttachTo, nodeToAttach, tree.AND)
	} else {
		// ... else simply assign entire subtree generated from the statement combination
		nodeToAttachTo = nodeToAttach
	}
	return nodeToAttachTo
}

/*
Handles parsing error centrally - easier to refine.
 */
func handleParsingError(component string, err tree.ParsingError) tree.ParsingError {

	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND {
		log.Println("Error when parsing component ", component, ": ", err)
		return err
	}

	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}


// Component prefix (word without spaces and parentheses, but [] brackets)
var componentPrefix = "([a-zA-Z\\[\\]]+)+"

/*
Escapes all special symbols to prepare those for input into regex expression
 */
func escapeSymbolsForRegex(text string) string {
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	text = strings.ReplaceAll(text, "(", "\\(")
	text = strings.ReplaceAll(text, ")", "\\)")
	text = strings.ReplaceAll(text, "[", "\\[")
	text = strings.ReplaceAll(text, "]", "\\]")

	return text
}

/*
Separates nested statement expressions (including component prefix)
from individual components (including combinations of components).
Returns multi-dim array, with element [0][0] containing component-only statement (no nested structure),
and element [1] containing nested statements (potentially multiple),
and element [2] containing potential statement combinations.
 */
func separateComponentsAndNestedStatements(statement string) ([][]string, tree.ParsingError) {

	// Prepare return structure
	ret := make([][]string,3)

	// Identify all nested statements
	nestedStmts, err := identifyNestedStatements(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	// Contains complete nested statements (with prefix)
	completeNestedStmts := []string{}

	// Holds candidates for nested combinations
	nestedCombos := []string{}

	if len(nestedStmts) > 0 {

		// Iterate through identified nested statements (if any) and remove those from statement
		for _, v := range nestedStmts {
			// Prepare pattern to extract nested statements including prefix from overall statement
			r, err := regexp.Compile(componentPrefix + escapeSymbolsForRegex(v))
			if err != nil {
				return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_PATTERN_EXTRACTION,
					ErrorMessage: "Error during pattern extraction in nested statement."}
			}
			// Extract nested statement including prefix embedded in overall statement
			result := r.FindAllStringSubmatch(statement, -1)
			if len(result) > 0 {
				// Append extracted nested statements including component prefix
				completeNestedStmts = append(completeNestedStmts, result[0][0])
				fmt.Println("Added candidate for single nested statement:", result[0][0])

				// Remove nested statement from overall statement
				statement = strings.ReplaceAll(statement, result[0][0], "")
			} else {
				// Save for parsing as combination
				nestedCombos = append(nestedCombos, v)
				fmt.Println("Added candidate for statement combination:", v)

				// Remove nested statement combination from overall statement
				statement = strings.ReplaceAll(statement, v, "")
				/*return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_PATTERN_EXTRACTION,
					ErrorMessage: "Unable to extract prefix of nested statement " + v}*/
			}
		}
		// Assign nested statements if found
		ret[1] = completeNestedStmts
		ret[2] = nestedCombos
		fmt.Println("Remaining non-nested input statement (without nested elements): " + statement)
	} else {
		fmt.Println("No nested statement found in input: " + statement)
	}

	// Assign component-only string
	ret[0] = []string{statement}

	fmt.Println("Array to be returned: " + fmt.Sprint(ret))

	// Return combined structure
	return ret, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Identifies nested statement patterns
 */
func identifyNestedStatements(statement string) ([]string, tree.ParsingError) {

	// Extract nested statements from input string
	nestedStatements, err := extractComponent("", statement, LEFT_BRACE, RIGHT_BRACE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	fmt.Println("Nested statements: " + fmt.Sprint(nestedStatements))

	return nestedStatements, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

func parseAttributes(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.ATTRIBUTES, text)
}

func parseAttributesProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.ATTRIBUTES_PROPERTY, text)
}

func parseDeontic(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.DEONTIC, text)
}

func parseAim(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.AIM, text)
}

func parseDirectObject(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.DIRECT_OBJECT, text)
}

func parseDirectObjectProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.DIRECT_OBJECT_PROPERTY, text)
}

func parseIndirectObject(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.INDIRECT_OBJECT, text)
}

func parseIndirectObjectProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.INDIRECT_OBJECT_PROPERTY, text)
}

func parseConstitutedEntity(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTED_ENTITY, text)
}

func parseConstitutedEntityProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTED_ENTITY_PROPERTY, text)
}

func parseModal(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.MODAL, text)
}

func parseConstitutingFunction(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTIVE_FUNCTION, text)
}

func parseConstitutingProperties(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTING_PROPERTIES, text)
}

func parseConstitutingPropertiesProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTING_PROPERTIES_PROPERTY, text)
}

func parseActivationCondition(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.ACTIVATION_CONDITION, text)
}

func parseExecutionConstraint(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.EXECUTION_CONSTRAINT, text)
}

/*
Validates input with respect to parentheses/braces balance.
Input is text to be tested, as well as left and right parenthesis/braces symbols ((,{, and ),}).
Parentheses symbols must be consistent, i.e., either both parentheses or braces.
 */
func validateInput(text string, leftPar string, rightPar string) (tree.ParsingError) {

	parTypeSingular := ""
	parTypePlural := ""

	if leftPar == LEFT_BRACE && rightPar == RIGHT_BRACE {
		parTypeSingular = "braces"
		parTypePlural = parTypeSingular
	} else if leftPar == LEFT_PARENTHESIS && rightPar == RIGHT_PARENTHESIS {
		parTypeSingular = "parenthesis"
		parTypePlural = "parentheses"
	} else {
		return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_PARENTHESES_COMBINATION,
			ErrorMessage: "Invalid combination of parentheses/braces during matching (e.g., (}, or {))"}
	}
	// Validate parentheses in input
	parCount := 0
	for i, letter := range text {

		switch string(letter) {
		case leftPar:
			parCount++
		case rightPar:
			parCount--
		}
		i++
	}
	if parCount != 0 {
		msg := "Please review the " + parTypePlural + " in the input statement. "
		par := ""
		parCountAbs := math.Abs(float64(parCount))
		if parCount == 1 || parCount == -1 {
			msg += "There is "
			par = parTypeSingular
		} else {
			msg += "There are "
			par = parTypePlural
		}
		if parCount > 0 {
			// too many left parentheses/braces
			msg = fmt.Sprint(msg, parCountAbs, " additional opening ", par, " ('" + leftPar + "').")
		} else {
			// too many right parentheses/braces
			msg = fmt.Sprint(msg, parCountAbs, " additional closing ", par, " ('" + rightPar + "').")
		}
		log.Println(msg)
		return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IMBALANCED_PARENTHESES, ErrorMessage: msg}
	}

	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Extracts a component specification from string based on component signature (e.g., A, I, etc.)
and balanced parentheses/braces.
If no component is found, an empty string is returned
*/
func extractComponent(component string, input string, leftPar string, rightPar string) ([]string, tree.ParsingError) {

	// Strings for given component
	componentStrings := []string{}

	// Copy string for truncating
	processedString := input

	fmt.Println("Looking for component: " + component)

	// Assume that parentheses/braces are checked beforehand

	for { // infinite loop - needs to break out
		// Find first occurrence of signature in processedString (incrementally iterated by letter)
		startPos := strings.Index(processedString, component + leftPar)

		if startPos == -1 {
			// Returns component strings once opening parenthesis symbol is no longer found
			return componentStrings, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}

		// Parentheses count to check for balance
		parCount := 0

		stop := false

		for i, letter := range processedString[startPos:] {

			switch string(letter) {
			case leftPar:
				parCount++
			case rightPar:
				parCount--
				if parCount == 0 {
					componentStrings = append(componentStrings, processedString[startPos:startPos+i+1])
					fmt.Println("Added string " + processedString[startPos:startPos+i+1])
					processedString = processedString[startPos+i+1:]
					stop = true
				}
			}
			if stop {
				break
			}
		}
	}
}

/*
Parses component based on surrounding parentheses.
*/
func parseComponentWithParentheses(component string, input string) (*tree.Node, tree.ParsingError) {
	return parseComponent(component, input, LEFT_PARENTHESIS, RIGHT_PARENTHESIS)
}

/*
Parses component based on surrounding braces
*/
func parseComponentWithBraces(component string, input string) (*tree.Node, tree.ParsingError) {
	return parseComponent(component, input, LEFT_BRACE, RIGHT_BRACE)
}

// Logical operators prepared for regular expression
var logicalOperators = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"
// Word pattern for regular expressions (including parentheses, spaces, square brackets, etc.)
var wordsWithParentheses = "([a-zA-Z',;()\\[\\]]+\\s*)+"
// Pattern of combinations, e.g., ( ... [AND] ... )
var combinationPattern = "\\(" + wordsWithParentheses + "(\\[" + logicalOperators + "\\]\\s" + wordsWithParentheses + ")+\\)"

func parseComponent(component string, text string, leftPar string, rightPar string) (*tree.Node, tree.ParsingError) {

	// Extract component (one or multiple occurrences) from input string based on provided component identifier
	componentStrings, err := extractComponent(component, text, leftPar, rightPar)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	fmt.Println("Components: " + fmt.Sprint(componentStrings))

	// Initialize output string for parsing
	componentString := ""

	// [AND]-link different components (if multiple occur in input string)
	if len(componentStrings) > 1 {
		r, err := regexp.Compile(combinationPattern)
		if err != nil {
			return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_PATTERN_EXTRACTION,
				ErrorMessage: "Error during pattern extraction in combination expression."}
		}
		// Add leading parenthesis
		componentString = LEFT_PARENTHESIS
		for i, v := range componentStrings {
			fmt.Println("Round: " + strconv.Itoa(i) + ": " + v)
			// Extract and concatenate individual component values but cut leading component identifier
			componentString += v[len(component):]
			// Identify whether combination embedded in input string
			result := r.FindAllStringSubmatch(componentString, -1)
			fmt.Println(result)
			if len(result) == 0 {
				// If no combination embedded in combination component, strip leading and trailing parentheses prior to combining
				componentString = componentString[1:len(componentString)-1]
			} // else don't touch, i.e., leave parentheses in string

			if i < len(componentStrings)-1 {
				// Add SAND primitive (synthetic linkage) in between if multiple component elements
				componentString += " " + tree.SAND_BRACKETS + " "
			} else {
				// Add trailing parenthesis
				componentString += RIGHT_PARENTHESIS
			}
		}
		//fmt.Println("Combination finished: " + componentString)
	} else if len(componentStrings) == 1 {
		// Single entry (cut prefix)
		componentString = componentStrings[0][len(component):]
		// Remove prefix including leading and trailing parenthesis (e.g., Bdir(, )) to extract inner string if not combined
		componentString = componentString[1:len(componentString)-1]
	} else {
		return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_COMPONENT_NOT_FOUND,
			ErrorMessage: "Component " + component + " was not found in input string"}
	}

	fmt.Println("Component Identifier: " + component)
	fmt.Println("Full string: " + componentString)

	//tree.PrintValueOrder = true

	fmt.Println("Preprocessed string: " + componentString)

	node, modifiedInput, err := ParseIntoNodeTree(componentString, false, leftPar, rightPar)

	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
		err.ErrorMessage = "Error when parsing component " + component + ": " + err.ErrorMessage
		log.Println("Error during component parsing: ", err.Error())
	}

	// Override missing combination error, since it is not relevant at this level
	if err.ErrorCode == tree.PARSING_NO_COMBINATIONS {
		err.ErrorCode = tree.PARSING_NO_ERROR
		err.ErrorMessage = ""
	}

	fmt.Println("Modified output for " + component + ": " + modifiedInput)

	return node, err
}
