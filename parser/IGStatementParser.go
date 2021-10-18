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

/*
Parses statement tree from input string.
 */
func ParseStatement(text string) (tree.Statement, tree.ParsingError) {

	// Remove line breaks
	text = CleanInput(text)

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
	compAndNestedStmts, err := SeparateComponentsAndNestedStatements(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return tree.Statement{}, err
	}
	fmt.Println("Returned separated components and nested statements: " + fmt.Sprint(compAndNestedStmts))
	// Extract component-only statement and override input
	text = compAndNestedStmts[0][0]
	// Extract potential nested statements
	nestedStmts := compAndNestedStmts[1]
	if len(nestedStmts) == 0 {
		fmt.Println("No nested statements found.")
		log.Println("No nested statements found.")
	}
	nestedCombos := compAndNestedStmts[2]
	if len(nestedCombos) == 0 {
		fmt.Println("No nested statement combination candidates found.")
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

	// Reorganize tree by shifting private nodes into PrivateNode fields of components and removing them from statement tree
	ProcessPrivateComponentLinkages(&s)

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
Parses nested statements (but not combinations) and attaches those to the top-level statement
 */
func parseNestedStatements(stmtToAttachTo *tree.Statement, nestedStmts []string) (tree.ParsingError) {

	for _, v := range nestedStmts {

		log.Println("Found nested statement")
		// Extract nested statement content and parse

		component := ""

		// Identify embedded component identifier - parse properties before main components to avoid wrongful mapping
		if strings.HasPrefix(v, tree.ATTRIBUTES_PROPERTY) {
			log.Println("Identified nested attributes property")
			component = tree.ATTRIBUTES_PROPERTY
		} else if strings.HasPrefix(v, tree.DIRECT_OBJECT_PROPERTY) {
			log.Println("Identified nested direct object property")
			component = tree.DIRECT_OBJECT_PROPERTY
		} else if strings.HasPrefix(v, tree.DIRECT_OBJECT) {
			log.Println("Identified nested direct object")
			component = tree.DIRECT_OBJECT
		} else if strings.HasPrefix(v, tree.INDIRECT_OBJECT_PROPERTY) {
			log.Println("Identified nested indirect object property")
			component = tree.INDIRECT_OBJECT_PROPERTY
		} else if strings.HasPrefix(v, tree.INDIRECT_OBJECT) {
			log.Println("Identified nested indirect object")
			component = tree.INDIRECT_OBJECT
		} else if strings.HasPrefix(v, tree.ACTIVATION_CONDITION) {
			log.Println("Identified nested activation condition")
			component = tree.ACTIVATION_CONDITION
		} else if strings.HasPrefix(v, tree.EXECUTION_CONSTRAINT) {
			log.Println("Identified nested execution constraint")
			component = tree.EXECUTION_CONSTRAINT
		} else if strings.HasPrefix(v, tree.CONSTITUTED_ENTITY_PROPERTY) {
			log.Println("Identified nested constituted entity property")
			component = tree.CONSTITUTED_ENTITY_PROPERTY
		} else if strings.HasPrefix(v, tree.CONSTITUTING_PROPERTIES_PROPERTY) {
			log.Println("Identified nested constituting properties property")
			component = tree.CONSTITUTING_PROPERTIES_PROPERTY
		} else if strings.HasPrefix(v, tree.CONSTITUTING_PROPERTIES) {
			log.Println("Identified nested constituting properties")
			component = tree.CONSTITUTING_PROPERTIES
		} else if strings.HasPrefix(v, tree.OR_ELSE) {
			log.Println("Identified nested or else")
			component = tree.OR_ELSE
		}

		// Extracting suffices and annotations
		suffix, annotation, _, err := extractSuffixAndAnnotations(component, v, LEFT_BRACE, RIGHT_BRACE)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			fmt.Print("Error during extraction of suffices and annotations: " + err.ErrorCode)
			return err
		}

		fmt.Println("Nested Stmt Component Identifier:", component)
		fmt.Println("Nested Stmt Suffix:", suffix)
		fmt.Println("Nested Stmt Annotation:", annotation)
		//fmt.Println("Nested Stmt Content:", content)

		stmt, errStmt := ParseStatement(v[strings.Index(v, LEFT_BRACE)+1:strings.LastIndex(v, RIGHT_BRACE)])
		if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
			fmt.Print("Error when parsing nested statements: " + errStmt.ErrorCode)
			return errStmt
		}

		// Wrap statement into node (since individual statement)
		stmtNode := tree.Node{Entry: stmt}
		// Assign component name to parsed node
		stmtNode.ComponentType = component

		// Attach suffix if it exists
		if suffix != "" {
			stmtNode.Suffix = suffix
		}

		// Attach annotation if it exists
		if annotation != "" {
			stmtNode.Annotations = annotation
		}

		// Identify component the coded information is to be attached to
		// Checks are ordered with property variants (e.g., Bdir,p) before component variants (e.g., Bdir) to avoid wrong match
		switch component {
			case tree.ATTRIBUTES_PROPERTY:
				log.Println("Attaching nested attributes property to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.AttributesPropertyComplex = attachComplexComponent(stmtToAttachTo.AttributesPropertyComplex, &stmtNode)
			case tree.DIRECT_OBJECT_PROPERTY:
				log.Println("Attaching nested direct object property to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.DirectObjectPropertyComplex = attachComplexComponent(stmtToAttachTo.DirectObjectPropertyComplex, &stmtNode)
			case tree.DIRECT_OBJECT:
				log.Println("Attaching nested direct object to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.DirectObjectComplex = attachComplexComponent(stmtToAttachTo.DirectObjectComplex, &stmtNode)
			case tree.INDIRECT_OBJECT_PROPERTY:
				log.Println("Attaching nested indirect object property to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.IndirectObjectPropertyComplex = attachComplexComponent(stmtToAttachTo.IndirectObjectPropertyComplex, &stmtNode)
			case tree.INDIRECT_OBJECT:
				log.Println("Attaching nested indirect object to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.IndirectObjectComplex = attachComplexComponent(stmtToAttachTo.IndirectObjectComplex, &stmtNode)
			case tree.ACTIVATION_CONDITION:
				log.Println("Attaching nested activation condition to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.ActivationConditionComplex = attachComplexComponent(stmtToAttachTo.ActivationConditionComplex, &stmtNode)
			case tree.EXECUTION_CONSTRAINT:
				log.Println("Attaching nested execution constraint to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.ExecutionConstraintComplex = attachComplexComponent(stmtToAttachTo.ExecutionConstraintComplex, &stmtNode)
			case tree.CONSTITUTED_ENTITY_PROPERTY:
				log.Println("Attaching nested constituted entity property to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.ConstitutedEntityPropertyComplex = attachComplexComponent(stmtToAttachTo.ConstitutedEntityPropertyComplex, &stmtNode)
			case tree.CONSTITUTING_PROPERTIES_PROPERTY:
				log.Println("Attaching nested constituting properties property to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.ConstitutingPropertiesPropertyComplex = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesPropertyComplex, &stmtNode)
			case tree.CONSTITUTING_PROPERTIES:
				log.Println("Attaching nested constituting properties to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.ConstitutingPropertiesComplex = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesComplex, &stmtNode)
			case tree.OR_ELSE:
				log.Println("Attaching nested or else to higher-level statement")
				// Assign nested statement to higher-level statement
				stmtToAttachTo.OrElse = attachComplexComponent(stmtToAttachTo.OrElse, &stmtNode)
		}
	}
	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Parses nested statement combinations and attaches those to the top-level statement
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
		flatCombo := tree.Flatten(combo.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES))
		sharedPrefix := ""
		for _, node := range flatCombo {
			if node.Entry == nil {
				return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_NIL_ELEMENT, ErrorMessage: "Nested combination returned nil element."}
			}
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

			// Extracting suffices and annotations
			suffix, annotation, content, err := extractSuffixAndAnnotations("", oldValue, LEFT_BRACE, RIGHT_BRACE)
			if err.ErrorCode != tree.PARSING_NO_ERROR {
				fmt.Print("Error during extraction of suffices and annotations: " + err.ErrorCode)
				return tree.Statement{}, tree.ParsingError{}
			}

			fmt.Println("Nested Combo Stmt Suffix:", suffix)
			fmt.Println("Nested Combo Stmt Annotation:", annotation)
			fmt.Println("Nested Combo Stmt Content:", content)

			stmt, errStmt := ParseStatement(oldValue[strings.Index(oldValue, LEFT_BRACE)+1:strings.LastIndex(oldValue, RIGHT_BRACE)])
			if errStmt.ErrorCode != tree.PARSING_NO_ERROR{
				return stmt, errStmt
			}
			return stmt, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		})
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return err
		}

		//TODO: Check whether combinations are actually filled, or just empty nodes (e.g., { Cac{ A(), I(), Cex() } [AND] Cac{ A(), I(), Cex() } })

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


/*
Separates nested statement expressions (including component prefix)
from individual components (including combinations of components).
Returns multi-dim array, with element [0][0] containing component-only statement (no nested structure),
and element [1] containing nested statements (potentially multiple),
and element [2] containing potential statement combinations.
 */
func SeparateComponentsAndNestedStatements(statement string) ([][]string, tree.ParsingError) {

	// Prepare return structure
	ret := make([][]string, 3)

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
			// Extract statements of structure { LEFT [AND] RIGHT }
			r2, err2 := regexp.Compile(NESTED_COMBINATIONS)
			if err2 != nil {
				return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_PATTERN_EXTRACTION,
					ErrorMessage: "Error during pattern extraction in nested statement."}
			}
			result2 := r2.FindAllString(v, -1)
			if len(result2) > 0 {
				// Identified combination of component-level nested statements

				// Save for parsing as combination
				nestedCombos = append(nestedCombos, v)
				fmt.Println("Added candidate for statement combination:", v)

				// Remove nested statement combination from overall statement
				statement = strings.ReplaceAll(statement, v, "")
			} else {
				// Identified single nested statement

				// Append extracted nested statements including component prefix
				completeNestedStmts = append(completeNestedStmts, v)
				fmt.Println("Added candidate for single nested statement:", v)

				// Remove nested statement from overall statement
				statement = strings.ReplaceAll(statement, v, "")
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

	// TODO: Review
	// Extract nested statements from input string
	nestedStatements, err := ExtractComponentContent("", statement, LEFT_BRACE, RIGHT_BRACE) //, "", "")
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
		parTypeSingular = "brace"
		parTypePlural = "braces"
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
Extracts a component content from string based on component signature (e.g., A, I, etc.)
and balanced parentheses/braces. Tolerates presence of suffices and annotations
If no component content is found, an empty string is returned.
Tests against mistaken parsing of property variant of a component (e.g., A,p() instead of A()).
*/
func ExtractComponentContent(component string, input string, leftPar string, rightPar string) ([]string, tree.ParsingError) {

	// Strings for given component
	componentStrings := []string{}

	// Copy string for truncating
	processedString := input

	fmt.Println("Looking for component: " + component)

	// Assume that parentheses/braces are checked beforehand

	// Validate component identifier presence (if component is specified as part of parameter)
	if component != "" && strings.Index(processedString, component) == -1 {
		return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_COMPONENT_NOT_FOUND}
	}

	// Start position
	startPos := -1

	// Search number of entries
	//r, err := regexp.Compile(component + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "?\\" + leftPar)
	// + escapeSymbolsForRegex(input)
	//fmt.Println("Regex:", component + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)
	r, err := regexp.Compile(component + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)
	if err != nil {
		return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Error in Regular Expression compilation."}
		//log.Fatal("Error", err.Error())
	}

	for { // infinite loop - needs to break out

		//// OLD STRING-BASED PARSING OF COMPONENTS
		// Find first occurrence of signature in processedString (incrementally iterated by letter)
		/*startPos := strings.Index(processedString, component + leftPar)

		if startPos == -1 {
			// Returns component strings once opening parenthesis symbol is no longer found
			return componentStrings, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}*/

		//// NEW REGEX-BASED PARSING (TO CONSIDER ANNOTATIONS AND SUFFICES)
		//fmt.Println("String to be searched for component:", processedString)
		// Return index of found element
		result := r.FindAllStringIndex(processedString, 1)
		resultContent := r.FindString(processedString)

		//fmt.Println("Index:", result)
		//fmt.Println("Content:", resultContent)

		if len(result) > 0 {
			// Start search after potential suffix and annotation elements
			startPos = result[0][0] + len(resultContent) - len(leftPar)
			fmt.Println("Start position: ", startPos)
		} else {
			// Returns component strings once opening parenthesis symbol is no longer found
			return componentStrings, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}

		// Parentheses count to check for balance
		parCount := 0

		// Switch to stop parsing
		stop := false

		for i, letter := range processedString[startPos:] {

			switch string(letter) {
			case leftPar:
				parCount++
			case rightPar:
				parCount--
				if parCount == 0 {
					// Store candidate string before cutting off potential leading component identifier (if nested statement)
					candidateString := resultContent[:len(resultContent)-len(leftPar)] + processedString[startPos:startPos+i+1]
					if !strings.HasSuffix(component, tree.PROPERTY_SYNTAX_SUFFIX) && strings.HasPrefix(candidateString, component + tree.PROPERTY_SYNTAX_SUFFIX) {
						// Don't consider if properties component is found (e.g., A,p(...)), but main component is sought (e.g., A(...)).
						fmt.Println("Ignoring found element due to ambiguous matching with property of component (Match: " +
							component + tree.PROPERTY_SYNTAX_SUFFIX + ", Component: " + component + ")")
					} else {
						componentStrings = append(componentStrings, candidateString)
						fmt.Println("Added string " + candidateString)
					}
					// String to be processed in next round is beyond identified component
					// This includes starting position of parentheses, but moves back to include component identifier,
					// suffix, annotation, parenthesis, extracted content string, closing parenthesis
					//processedString = processedString[startPos-len(resultContent)-len(leftPar)+len(candidateString):]
					//processedString = processedString[startPos-len(component)-len(resultContent)-len(leftPar)+len(candidateString)+len(rightPar):]
					idx := strings.Index(processedString, candidateString)
					if idx == -1 {
						return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR,
							ErrorMessage: "Extracted expression cannot be found in processed string (Search string: " + candidateString + ")"}
					}
					// Cut found string and leave remainder for further processing
					processedString = processedString[idx + len(candidateString):]
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
Extracts suffix (e.g., ,p1) and annotations (e.g., [ctx=time]), and content from IG-Script-coded input.
It takes component identifier and raw coded information as input, as well as left and right parenthesis symbols (e.g., (,) or {,}).
Returns suffix as first element, annotations string as second, and component content (including identifier) as third element.
IMPORTANT:
- This function will only extract the suffix and annotation for the first element of a given component type found in the input string.
- This function will not prevent wrongful extraction of property components instead of first-order components. This is handled in #ExtractComponentContent.
TODO: Make this more efficient
 */
func extractSuffixAndAnnotations(component string, input string, leftPar string, rightPar string) (string, string, string, tree.ParsingError) {

	fmt.Println("Component:", component)
	fmt.Println("Input:", input)
	strippedInput := input // leave input unchanged
	
	// Component annotation pattern
	r, err := regexp.Compile(COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)
	// + escapeSymbolsForRegex(input)
	if err != nil {
		log.Fatal("Error", err.Error())
	}
	// Search for annotation pattern on input (without leading component identifier)
	result := r.FindAllStringSubmatch(strippedInput, 1)

	// The result will find the leftPar as a minimum (e.g., "(" or "{"). The processing needs to account for this

	if len(result) > 0 && result[0][0] != leftPar {
		// If annotations are found ...
		res := result[0][0]
		fmt.Println("Found annotation in component:", res)
		// Extract semantic annotation string
		res = res[:len(res)-1]
		//fmt.Println("Annotations:", res)
		pos := strings.Index(strippedInput, res)
		suffix := ""
		// Only attempt to extract suffix if there is actually one
		if pos > len(component) {
			// Extract component name suffix (e.g., 1), but remove component identifier
			suffix = strippedInput[len(component):pos]
			// Does not guard against mistaken choice of property variants of components (e.g., A,p instead of A) - is handled in #ExtractComponentContent.
		}
		reconstructedComponent, err := ExtractComponentContent(component, strings.ReplaceAll(strippedInput, suffix + res, ""), leftPar, rightPar)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return "", "", "", err
		}
		fmt.Println("Reconstructed statement:", reconstructedComponent)
		// Return suffix and annotations
		return suffix, res, reconstructedComponent[0], tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
	} else {
		fmt.Println("No annotations found ...")
		// ... if no annotations are found ...
		// Identifier start position for content
		contentStartPos := strings.Index(strippedInput, leftPar)
		suffix := ""
		// Component identifier is suppressed if suffix is found
		// Extract suffix (e.g., 1), but remove component identifier
		suffix = strippedInput[len(component):contentStartPos]
		// Does not guard against mistaken choice of property variants of components (e.g., A,p instead of A) - is handled in #ExtractComponentContent.
		reconstructedComponent := strings.Replace(strippedInput, suffix, "", 1)
		fmt.Println("Reconstructed statement:", reconstructedComponent)
		// Return only suffix
		return suffix, "", reconstructedComponent, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
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

/*
Generic entry point to parse individual components of a given statement.
Input is component symbol of interest, full statements, as well as delimiting parentheses signaling parsing for atomic
or nested components.
Returns the parsed node.
 */
func parseComponent(component string, text string, leftPar string, rightPar string) (*tree.Node, tree.ParsingError) {

	fmt.Println("Parsing:", component)

	// Extract component (one or multiple occurrences) from input string based on provided component identifier
	componentStrings, err := ExtractComponentContent(component, text, leftPar, rightPar)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	fmt.Println("Components (Count:", len(componentStrings), "):", fmt.Sprint(componentStrings))

	// Initialize output string for parsing
	componentString := ""

	// Node to be populated as return node
	node := &tree.Node{}

	// Synthetically linked ([sAND]) components (if multiple occur in input string)
	if len(componentStrings) > 1 {
		fmt.Println("Component combination for component", component)
		fmt.Println("Component content", componentStrings)
		r, err := regexp.Compile(COMBINATION_PATTERN_PARENTHESES)
		if err != nil {
			return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_PATTERN_EXTRACTION,
				ErrorMessage: "Error during pattern extraction in combination expression."}
		}

		for i, v := range componentStrings {
			fmt.Println("Round: " + strconv.Itoa(i) + ": " + v)

			// Extracts suffix and/or annotation for individual component instance -- must only be used with single component instance!
			componentSuffix, componentAnnotation, componentContent, err := extractSuffixAndAnnotations(component, v, leftPar, rightPar)
			if err.ErrorCode != tree.PARSING_NO_ERROR {
				return nil, err
			}

			fmt.Println("Suffix:", componentSuffix, "(Length:", len(componentSuffix), ")")
			// Store suffices
			//suffices = append(suffices, componentSuffix)
			fmt.Println("Annotations:", componentAnnotation, "(Length:", len(componentAnnotation), ")")
			// Store annotations
			//annotations = append(annotations, componentAnnotation)
			fmt.Println("Content:", componentContent)

			// Extract and concatenate individual component values but cut leading component identifier
			componentWithoutIdentifier := componentContent[len(component):]
			// Identify whether combination embedded in input string element
			result := r.FindAllStringSubmatch(componentWithoutIdentifier, -1)
			fmt.Println(result)
			fmt.Println("Length:", len(result))
			fmt.Println("Component string before:", componentWithoutIdentifier)
			if len(result) == 0 {
				// If no combination embedded in combination component, strip leading and trailing parentheses prior to combining
				componentWithoutIdentifier = componentWithoutIdentifier[1:len(componentWithoutIdentifier)-1]
			} // else don't touch, i.e., leave parentheses in string
			fmt.Println("Component string after:", componentWithoutIdentifier)

			// Parse first component into node
			if node.IsEmptyNode() {
				node1, _, err := ParseIntoNodeTree(componentWithoutIdentifier, false, leftPar, rightPar)
				if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
					log.Println("Error when parsing synthetically linked element. Error:", err)
					return nil, err
				}
				// Assign to main node if not populated and new node not nil
				if !node1.IsEmptyNode() {
					node = node1
					// Attach component name to element (will be accessible to children via GetComponentName())
					node.ComponentType = component
					// Attach node-specific suffix
					if componentSuffix != "" {
						node.Suffix = componentSuffix
					}
					// Attach node-specific annotations
					if componentAnnotation != "" {
						node.Annotations = componentAnnotation
					}
				}
			} else {
				// Parse any additional components into node and combine
				// If cached node is already populated, create separate node and link afterwards
				node2, _, err := ParseIntoNodeTree(componentWithoutIdentifier, false, leftPar, rightPar)
				if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
					log.Println("Error when parsing synthetically linked element. Error:", err)
					return nil, err
				}
				if !node2.IsEmptyNode() {
					// Attach component name to element (will be accessible to children via GetComponentName())
					node2.ComponentType = component
					// Attach node-specific suffix
					if componentSuffix != "" {
						node2.Suffix = componentSuffix
					}
					// Attach node-specific annotations
					if componentAnnotation != "" {
						node2.Annotations = componentAnnotation
					}
					// Combine existing node with newly created one based on synthetic AND
					nodeComb := tree.Combine(node, node2, tree.SAND_BETWEEN_COMPONENTS)
					// Explicitly assign component type to top-level node (for completeness)
					nodeComb.ComponentType = component
					// Assign to return node
					node = nodeComb
				}
			}
		}
	} else if len(componentStrings) == 1 {

		fmt.Println("Component strings:", componentStrings)

		// Extracts suffix and/or annotation for individual component instance -- must only be used with single component instance!
		componentSuffix, componentAnnotation, componentContent, err := extractSuffixAndAnnotations(component, componentStrings[0], leftPar, rightPar)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return nil, err
		}

		fmt.Println("Suffix:", componentSuffix, "(Length:", len(componentSuffix), ")")
		// Store suffices
		fmt.Println("Annotations:", componentAnnotation, "(Length:", len(componentAnnotation), ")")
		// Store annotations
		fmt.Println("Content:", componentContent)

		// Single entry (cut prefix)
		componentString = componentContent[len(component):]
		fmt.Println("Single component for component", component)
		fmt.Println("Component content", componentString)
		// Remove prefix including leading and trailing parenthesis (e.g., Bdir(, )) to extract inner string if not combined
		componentString = componentString[1:len(componentString)-1]

		node1, _, err := ParseIntoNodeTree(componentString, false, leftPar, rightPar)
		if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
			log.Println("Error when parsing synthetically linked element. Error:", err)
			return nil, err
		}
		// Attach component name to top-level element (will be accessible to children via GetComponentName())
		if !node1.IsEmptyNode() {
			node1.ComponentType = component
			// Attach node-specific suffix
			if componentSuffix != "" {
				node1.Suffix = componentSuffix
			}
			// Attach node-specific annotations
			if componentAnnotation != "" {
				node1.Annotations = componentAnnotation
			}
			// Overwrite main node
			node = node1
		}
	} else {
		return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_COMPONENT_NOT_FOUND,
			ErrorMessage: "Component " + component + " was not found in input string"}
	}

	fmt.Println("Component Identifier: " + component)
	fmt.Println("Full string: " + componentString)

	// Some error check and override
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
		err.ErrorMessage = "Error when parsing component " + component + ": " + err.ErrorMessage
		log.Println("Error during component parsing:", err.Error())
	}

	// Override missing combination error, since it is not relevant at this level
	if err.ErrorCode == tree.PARSING_NO_COMBINATIONS {
		err.ErrorCode = tree.PARSING_NO_ERROR
		err.ErrorMessage = ""
	}

	return node, err
}
