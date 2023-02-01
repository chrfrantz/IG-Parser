package parser

import (
	"IG-Parser/core/tree"
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"strings"
)

/*
Parses statement tree from input string. Returns statement tree, and error.
If parsing is successful, error code tree.PARSING_NO_ERROR is returned, else
other context-specific codes are returned.
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
	Println("Returned separated components and nested statements: " + fmt.Sprint(compAndNestedStmts))
	// Extract component-only statement and override input
	text = compAndNestedStmts[0][0]
	// Extract potential nested statements
	nestedStmts := compAndNestedStmts[1]
	if len(nestedStmts) == 0 {
		Println("No nested statements found.")
		log.Println("No nested statements found.")
	}
	nestedCombos := compAndNestedStmts[2]
	if len(nestedCombos) == 0 {
		Println("No nested statement combination candidates found.")
		log.Println("No nested statement combination candidates found.")
	}

	Println("Text to be parsed: " + text)
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

	//Println(s.String())

	Println("Testing for nested statement combinations in " + fmt.Sprint(nestedCombos))

	// Process nested statement combinations
	if len(nestedCombos) > 0 {
		log.Println("Found nested statement combinations ...")
		err = parseNestedStatementCombinations(&s, nestedCombos)
		if err.ErrorCode == tree.PARSING_ERROR_NIL_ELEMENT {
			// Shift to regular nested statement if parsing as combo failed (Regex is too coarse-grained)
			nestedStmts = append(nestedStmts, nestedCombos...)
			Println("Reclassifying statement as nested statement (as opposed to nested combination) ...")
		} else if err.ErrorCode != tree.PARSING_NO_ERROR {
			return s, err
		}
	}

	Println("Testing for nested statements in " + fmt.Sprint(nestedStmts))

	// Process nested statements
	if len(nestedStmts) > 0 {
		log.Println("Found nested statements ...")
		err = parseNestedStatements(&s, nestedStmts)
		// Check whether nested statements have been ignored entirely
		if err.ErrorCode == tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS {
			return s, err
		}
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return s, err
		}
	}

	log.Println("Statement (after assigning sub elements):\n" + s.String())

	return s, outErr
}

/*
Parses nested statements (but not combinations) and attaches those to the top-level statement.
Returns an error other than tree.PARSING_NO_ERROR if no issues during parsing.
Returns err tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS if parsing did not pose problems per se,
but elements have been ignored during parsing (warranting syntax review). In this case, the statements of concern
are returned in a string array contained in the error object.
*/
func parseNestedStatements(stmtToAttachTo *tree.Statement, nestedStmts []string) tree.ParsingError {

	// Copy reference statement for comparison (to check whether modification took place based on parsed element)
	cachedStmtPriorToNestedParsing := stmtToAttachTo.String()

	// Array to keep track of ignored statements
	nestedStmtsNotConsideredDuringParsing := make([]string, 0)

	for _, v := range nestedStmts {

		log.Println("Processing nested statement: ", v)
		// Extract nested statement content and parse

		component := ""
		prefix := ""
		isProperty := false

		// Test prefix for nested statements (and remove is present before further exploring component)
		leadIdx := strings.Index(v, LEFT_BRACE)
		if leadIdx != -1 {
			prefix = v[:leadIdx]
		}

		// Identify embedded component identifier - parse properties before main components to avoid wrongful mapping
		if strings.HasPrefix(prefix, tree.ATTRIBUTES_PROPERTY) ||
			(strings.HasPrefix(prefix, tree.ATTRIBUTES) && strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX)) {
			log.Println("Identified nested Attributes Property")
			component = tree.ATTRIBUTES_PROPERTY
			isProperty = true
		} else if strings.HasPrefix(prefix, tree.DIRECT_OBJECT_PROPERTY) ||
			(strings.HasPrefix(prefix, tree.DIRECT_OBJECT) && strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX)) {
			log.Println("Identified nested Direct Object Property")
			component = tree.DIRECT_OBJECT_PROPERTY
			isProperty = true
		} else if strings.HasPrefix(prefix, tree.DIRECT_OBJECT) && !strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX) {
			log.Println("Identified nested Direct Object")
			component = tree.DIRECT_OBJECT
		} else if strings.HasPrefix(prefix, tree.INDIRECT_OBJECT_PROPERTY) ||
			(strings.HasPrefix(prefix, tree.INDIRECT_OBJECT) && strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX)) {
			log.Println("Identified nested Indirect Object Property")
			component = tree.INDIRECT_OBJECT_PROPERTY
			isProperty = true
		} else if strings.HasPrefix(prefix, tree.INDIRECT_OBJECT) && !strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX) {
			log.Println("Identified nested Indirect Object")
			component = tree.INDIRECT_OBJECT
		} else if strings.HasPrefix(prefix, tree.ACTIVATION_CONDITION) {
			log.Println("Identified nested Activation Condition")
			component = tree.ACTIVATION_CONDITION
		} else if strings.HasPrefix(prefix, tree.EXECUTION_CONSTRAINT) {
			log.Println("Identified nested Execution Constraint")
			component = tree.EXECUTION_CONSTRAINT
		} else if strings.HasPrefix(prefix, tree.CONSTITUTED_ENTITY_PROPERTY) ||
			(strings.HasPrefix(prefix, tree.CONSTITUTED_ENTITY) && strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX)) {
			log.Println("Identified nested Constituted Entity Property")
			component = tree.CONSTITUTED_ENTITY_PROPERTY
			isProperty = true
		} else if strings.HasPrefix(prefix, tree.CONSTITUTING_PROPERTIES_PROPERTY) ||
			(strings.HasPrefix(prefix, tree.CONSTITUTING_PROPERTIES) && strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX)) {
			log.Println("Identified nested Constituting Properties Property")
			component = tree.CONSTITUTING_PROPERTIES_PROPERTY
			isProperty = true
		} else if strings.HasPrefix(prefix, tree.CONSTITUTING_PROPERTIES) && !strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX) {
			log.Println("Identified nested Constituting Properties")
			component = tree.CONSTITUTING_PROPERTIES
		} else if strings.HasPrefix(prefix, tree.OR_ELSE) {
			log.Println("Identified nested Or Else")
			component = tree.OR_ELSE
		}
		// TODO: Check whether nesting on unsupported components is a challenge

		// Extracting suffices and annotations
		suffix, annotation, _, err := extractSuffixAndAnnotations(component, isProperty, v, LEFT_BRACE, RIGHT_BRACE)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			fmt.Println("Error during extraction of suffices and annotations on component '" + component + "': " + err.ErrorCode)
			return err
		}

		Println("Nested Stmt Component Identifier:", component)
		Println("Nested Stmt Suffix:", suffix)
		Println("Nested Stmt Annotation:", annotation)
		//Println("Nested Stmt Content:", content)

		stmt, errStmt := ParseStatement(v[strings.Index(v, LEFT_BRACE)+1 : strings.LastIndex(v, RIGHT_BRACE)])
		if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
			fmt.Println("Error when parsing nested statements: ", errStmt.Error())
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

		// Default error for node combination - can generally only be overridden by detected invalid component combinations
		nodeCombinationError := tree.NodeError{ErrorCode: tree.TREE_NO_ERROR}

		// Identify component the coded information is to be attached to
		// Checks are ordered with property variants (e.g., Bdir,p) before component variants (e.g., Bdir) to avoid wrong match
		switch component {
		case tree.ATTRIBUTES_PROPERTY:
			log.Println("Attaching nested attributes property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.AttributesPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.AttributesPropertyComplex, &stmtNode)
		case tree.DIRECT_OBJECT_PROPERTY:
			log.Println("Attaching nested direct object property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.DirectObjectPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.DirectObjectPropertyComplex, &stmtNode)
		case tree.DIRECT_OBJECT:
			log.Println("Attaching nested direct object to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.DirectObjectComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.DirectObjectComplex, &stmtNode)
		case tree.INDIRECT_OBJECT_PROPERTY:
			log.Println("Attaching nested indirect object property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.IndirectObjectPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.IndirectObjectPropertyComplex, &stmtNode)
		case tree.INDIRECT_OBJECT:
			log.Println("Attaching nested indirect object to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.IndirectObjectComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.IndirectObjectComplex, &stmtNode)
		case tree.ACTIVATION_CONDITION:
			log.Println("Attaching nested activation condition to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ActivationConditionComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ActivationConditionComplex, &stmtNode)
		case tree.EXECUTION_CONSTRAINT:
			log.Println("Attaching nested execution constraint to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ExecutionConstraintComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ExecutionConstraintComplex, &stmtNode)
		case tree.CONSTITUTED_ENTITY_PROPERTY:
			log.Println("Attaching nested constituted entity property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutedEntityPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutedEntityPropertyComplex, &stmtNode)
		case tree.CONSTITUTING_PROPERTIES_PROPERTY:
			log.Println("Attaching nested constituting properties property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutingPropertiesPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesPropertyComplex, &stmtNode)
		case tree.CONSTITUTING_PROPERTIES:
			log.Println("Attaching nested constituting properties to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutingPropertiesComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesComplex, &stmtNode)
		case tree.OR_ELSE:
			log.Println("Attaching nested or else to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.OrElse, nodeCombinationError = attachComplexComponent(stmtToAttachTo.OrElse, &stmtNode)
		}
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		}

		// Check if the iterated nested statement has been ignored entirely --> indicates failed detection as nested (as opposed to mere parsing problem)
		if cachedStmtPriorToNestedParsing == stmtToAttachTo.String() {
			log.Println("Nested statement has not been considered during parsing:", v)
			nestedStmtsNotConsideredDuringParsing = append(nestedStmtsNotConsideredDuringParsing, v)
		}
		// Overwrite cached content in any case to ensure capturing further statement potentially ignored during parsing
		cachedStmtPriorToNestedParsing = stmtToAttachTo.String()

	}
	if len(nestedStmtsNotConsideredDuringParsing) > 0 {
		// Indicate if elements have been ignored entirely (and return info to screen), but no further parsing errors per se
		msg := "Selected nested elements could not be properly parsed: '" + strings.Join(nestedStmtsNotConsideredDuringParsing[:], ",") +
			"'. Please review the input coding accordingly."
		return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS,
			ErrorMessage:         msg,
			ErrorIgnoredElements: nestedStmtsNotConsideredDuringParsing}
	}
	// if parsing worked out and if no elements have been ignored
	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Parses nested statement combinations and attaches those to the top-level statement
*/
func parseNestedStatementCombinations(stmtToAttachTo *tree.Statement, nestedCombos []string) tree.ParsingError {

	// Default error for node combination - can generally only be overridden by detected invalid component combinations
	nodeCombinationError := tree.NodeError{ErrorCode: tree.TREE_NO_ERROR}

	for _, v := range nestedCombos {

		// Check if combination error has been picked up - in the beginning and end of loop
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		}

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
			Println("Entry to parse for component type: " + entry)
			// Extract prefix (i.e., component type) for node, but check whether it contains nested statement
			if strings.Index(entry, LEFT_BRACE) == -1 {
				return tree.ParsingError{ErrorCode: tree.PARSING_INVALID_COMBINATION, ErrorMessage: "Element in combination of nested statement does not contain nested statement. Element of concern: " + entry}
			}
			prefix, prop, err := extractComponentType(entry[:strings.Index(entry, LEFT_BRACE)])
			if err.ErrorCode != tree.PARSING_NO_ERROR {
				// Return error and propagate error message from called function
				return tree.ParsingError{ErrorCode: err.ErrorCode, ErrorMessage: "Error when extracting component type from nested statement: " + err.ErrorMessage}
			}
			// Extract suffix and annotation
			suffix, annotation, _, err := extractSuffixAndAnnotations(prefix, prop, entry, LEFT_BRACE, RIGHT_BRACE)
			if err.ErrorCode != tree.PARSING_NO_ERROR {
				return tree.ParsingError{ErrorCode: err.ErrorCode, ErrorMessage: "Failed to extract suffix or annotation of nested statement."}
			}
			if len(suffix) > 0 {
				node.Suffix = suffix
			}
			if len(annotation) > 0 {
				node.Annotations = annotation
			}
			// Cache the component type for comparison in combinations
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

			// Check whether the combination element contains a nested structure ...
			tempComponentType := oldValue
			if strings.Contains(oldValue, LEFT_BRACE) {
				// ... and remove the nested element prior to parsing
				tempComponentType = oldValue[:strings.Index(oldValue, LEFT_BRACE)]
			}

			// Extract component type (after stripping potential nested statements)
			compType, prop, err := extractComponentType(tempComponentType)
			if err.ErrorCode != tree.PARSING_NO_ERROR {
				return tree.Statement{}, err
			}
			// Extracting suffices and annotations
			suffix, annotation, content, err := extractSuffixAndAnnotations(compType, prop, oldValue, LEFT_BRACE, RIGHT_BRACE)
			if err.ErrorCode != tree.PARSING_NO_ERROR {
				fmt.Println("Error during extraction of suffices and annotations of component '" + compType + "': " + err.ErrorCode)
				return tree.Statement{}, err
			}

			Println("Nested Combo Stmt Suffix:", suffix)
			Println("Nested Combo Stmt Annotation:", annotation)
			Println("Nested Combo Stmt Content:", content)

			stmt, errStmt := ParseStatement(oldValue[strings.Index(oldValue, LEFT_BRACE)+1 : strings.LastIndex(oldValue, RIGHT_BRACE)])
			if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
				return stmt, errStmt
			}
			return stmt, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		})
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return err
		}

		//TODO: Check whether combinations are actually filled, or just empty nodes (e.g., { Cac{ A(), I(), Cex() } [AND] Cac{ A(), I(), Cex() } })

		log.Println("Assigning nested tree structure", combo)

		// Assign component type name to combination (for proper retrieval and identification as correct type)
		combo.ComponentType = sharedPrefix
		Println("Combo component prefix:", sharedPrefix)

		// Checks are ordered with property variants (e.g., Bdir,p) before component variants (e.g., Bdir) to avoid wrong match

		if strings.HasPrefix(sharedPrefix, tree.ATTRIBUTES_PROPERTY) {
			log.Println("Attaching nested attributes property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.AttributesPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.AttributesPropertyComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.DIRECT_OBJECT_PROPERTY) {
			log.Println("Attaching nested direct object property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.DirectObjectPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.DirectObjectPropertyComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.DIRECT_OBJECT) {
			log.Println("Attaching nested direct object to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.DirectObjectComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.DirectObjectComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.INDIRECT_OBJECT_PROPERTY) {
			log.Println("Attaching nested indirect object property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.IndirectObjectPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.IndirectObjectPropertyComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.INDIRECT_OBJECT) {
			log.Println("Attaching nested indirect object to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.IndirectObjectComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.IndirectObjectComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.ACTIVATION_CONDITION) {
			log.Println("Attaching nested activation condition to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ActivationConditionComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ActivationConditionComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.EXECUTION_CONSTRAINT) {
			log.Println("Attaching nested execution constraint to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ExecutionConstraintComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ExecutionConstraintComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.CONSTITUTED_ENTITY_PROPERTY) {
			log.Println("Attaching nested constituted entity property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutedEntityPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutedEntityPropertyComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.CONSTITUTING_PROPERTIES_PROPERTY) {
			log.Println("Attaching nested constituting properties property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutingPropertiesPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesPropertyComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.CONSTITUTING_PROPERTIES) {
			log.Println("Attaching nested constituting properties to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutingPropertiesComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesComplex, combo)
			continue
		}
		if strings.HasPrefix(sharedPrefix, tree.OR_ELSE) {
			log.Println("Attaching nested or else to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.OrElse, nodeCombinationError = attachComplexComponent(stmtToAttachTo.OrElse, combo)
		}
	}
	// Check if combination error has been picked up - here and in the beginning of loop
	if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
		return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
			ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
	}
	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Attach complex component to tree structure under consideration of existing nodes in target tree structure.
Input:
- Node of the parent tree to attach to
- Node to attach
*/
func attachComplexComponent(nodeToAttachTo *tree.Node, nodeToAttach *tree.Node) (*tree.Node, tree.NodeError) {
	log.Println("Attaching nested activation condition to higher-level statement")
	// Assign nested statement to higher-level statement

	// If already a statement assignment to complex element, ...
	if nodeToAttachTo != nil {
		// ... combine both
		newNode, err := tree.Combine(nodeToAttachTo, nodeToAttach, tree.AND)
		if err.ErrorCode != tree.TREE_NO_ERROR {
			return nil, err
		}
		// Assign to input node
		nodeToAttachTo = newNode
	} else {
		// ... else simply assign entire subtree generated from the statement combination
		nodeToAttachTo = nodeToAttach
	}
	return nodeToAttachTo, tree.NodeError{ErrorCode: tree.TREE_NO_ERROR}
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
			r2, err2 := regexp.Compile(NESTED_COMBINATIONS_TERMINATED)
			if err2 != nil {
				return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_PATTERN_EXTRACTION,
					ErrorMessage: "Error during pattern extraction in nested statement."}
			}
			result2 := r2.FindAllString(v, -1)
			if len(result2) > 0 {
				// Identified combination of component-level nested statements

				// Save for parsing as combination
				nestedCombos = append(nestedCombos, v)
				Println("Added candidate for statement combination:", v)

				// Remove nested statement combination from overall statement
				statement = strings.ReplaceAll(statement, v, "")
			} else {
				// Identified single nested statement

				// Append extracted nested statements including component prefix
				completeNestedStmts = append(completeNestedStmts, v)
				Println("Added candidate for single nested statement:", v)

				// Remove nested statement from overall statement
				statement = strings.ReplaceAll(statement, v, "")
			}
		}
		// Assign nested statements if found
		ret[1] = completeNestedStmts
		ret[2] = nestedCombos
		Println("Remaining non-nested input statement (without nested elements): " + statement)
	} else {
		Println("No nested statement found in input: " + statement)
	}

	// Assign component-only string
	ret[0] = []string{statement}

	Println("Array to be returned: " + fmt.Sprint(ret))

	// Return combined structure
	return ret, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Identifies nested statement patterns
*/
func identifyNestedStatements(statement string) ([]string, tree.ParsingError) {

	// TODO: Review
	// Extract nested statements from input string
	nestedStatements, err := ExtractComponentContent("", true, statement, LEFT_BRACE, RIGHT_BRACE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	Println("Nested statements: " + fmt.Sprint(nestedStatements))

	return nestedStatements, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

func parseAttributes(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.ATTRIBUTES, false, text)
}

func parseAttributesProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.ATTRIBUTES_PROPERTY, true, text)
}

func parseDeontic(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.DEONTIC, false, text)
}

func parseAim(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.AIM, false, text)
}

func parseDirectObject(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.DIRECT_OBJECT, false, text)
}

func parseDirectObjectProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.DIRECT_OBJECT_PROPERTY, true, text)
}

func parseIndirectObject(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.INDIRECT_OBJECT, false, text)
}

func parseIndirectObjectProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.INDIRECT_OBJECT_PROPERTY, true, text)
}

func parseConstitutedEntity(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTED_ENTITY, false, text)
}

func parseConstitutedEntityProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTED_ENTITY_PROPERTY, true, text)
}

func parseModal(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.MODAL, false, text)
}

func parseConstitutingFunction(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTIVE_FUNCTION, false, text)
}

func parseConstitutingProperties(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTING_PROPERTIES, false, text)
}

func parseConstitutingPropertiesProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTING_PROPERTIES_PROPERTY, true, text)
}

func parseActivationCondition(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.ACTIVATION_CONDITION, false, text)
}

func parseExecutionConstraint(text string) (*tree.Node, tree.ParsingError) {
	return parseComponentWithParentheses(tree.EXECUTION_CONSTRAINT, false, text)
}

/*
Validates input with respect to parentheses/braces balance.
Input is text to be tested, as well as left and right parenthesis/braces symbols ((,{, and ),}).
Parentheses symbols must be consistent, i.e., either both parentheses or braces.
*/
func validateInput(text string, leftPar string, rightPar string) tree.ParsingError {

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
			msg = fmt.Sprint(msg, parCountAbs, " additional opening ", par, " ('"+leftPar+"').")
		} else {
			// too many right parentheses/braces
			msg = fmt.Sprint(msg, parCountAbs, " additional closing ", par, " ('"+rightPar+"').")
		}
		log.Println(msg)
		return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IMBALANCED_PARENTHESES, ErrorMessage: msg}
	}

	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Extracts a component content from string based on component signature (e.g., A, I, etc.)
and balanced parentheses/braces. Tolerates presence of suffices and annotations and includes those
in output (e.g., A1[type=animate](content)).
Allows for indication whether parsed component is actually a property.
If no component content is found, an empty string is returned.
Tests against mistaken parsing of property variant of a component (e.g., A,p() instead of A()).
*/
func ExtractComponentContent(component string, propertyComponent bool, input string, leftPar string, rightPar string) ([]string, tree.ParsingError) {

	// Strings for given component
	componentStrings := []string{}

	// Copy string for truncating
	processedString := input

	Println("Looking for component: ", component, "in", input, "(Property:", propertyComponent, ")")

	// Assume that parentheses/braces are checked beforehand

	// Switch indicating nested statement structure
	nestedStatement := false

	// Start position
	startPos := -1

	// Search number of entries
	//r, err := regexp.Compile(component + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "?\\" + leftPar)
	// + escapeSymbolsForRegex(input)
	//Println("Regex:", component + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)

	// General component syntax (inclusive of ,p)
	r, err := regexp.Compile(component + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)
	if err != nil {
		return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Error in Regular Expression compilation."}
		//log.Fatal("Error", err.Error())
	}
	// Component syntax to test for suffix-embedded property syntax (e.g., A1,p)
	rProp, err := regexp.Compile(component + COMPONENT_SUFFIX_SYNTAX + tree.PROPERTY_SYNTAX_SUFFIX + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)
	if err != nil {
		return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Error in Regular Expression compilation."}
		//log.Fatal("Error", err.Error())
	}

	if propertyComponent {
		Println("Identified as component", component, "as property:", propertyComponent)
		// If component is a property, extract root symbol to allow for intermediate index/suffix (e.g., A1,p)
		leadIdx := strings.Index(component, tree.PROPERTY_SYNTAX_SUFFIX)
		if leadIdx != -1 {
			// If property element is indeed found, strip it for regex generation
			componentRoot := component[:leadIdx]

			r, err = regexp.Compile(componentRoot + COMPONENT_SUFFIX_SYNTAX + tree.PROPERTY_SYNTAX_SUFFIX + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)
			if err != nil {
				return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Error in Regular Expression compilation."}
				//log.Fatal("Error", err.Error())
			}

		}
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
		//Println("String to be searched for component:", processedString)
		// Return index of found element
		result := r.FindAllStringIndex(processedString, 1)
		// Return content of found element
		resultContent := r.FindString(processedString)

		//Println("Index:", result)
		//Println("Component prefix:", resultContent)
		if nestedStatement && len(resultContent) > 0 {
			component = resultContent[:len(resultContent)-1]
			Println("Identified nested component", component)
		}

		if len(result) > 0 {
			// Start search after potential suffix and annotation elements
			startPos = result[0][0] + len(resultContent) - len(leftPar)
			Println("Start position: ", startPos)
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
					// Lead string including component identifier, suffices and annotations
					leadString := resultContent[:len(resultContent)-len(leftPar)]
					// String containing content only (including parentheses)
					contentString := processedString[startPos : startPos+i+1]
					Println(contentString)
					// Store candidate string before cutting off potential leading component identifier (if nested statement)
					candidateString := leadString + contentString
					if !strings.HasSuffix(component, tree.PROPERTY_SYNTAX_SUFFIX) && !propertyComponent &&
						// Test whether property is accidentally embedded but it is actually non-property component search
						rProp.MatchString(candidateString) {
						// Don't consider if properties component is found (e.g., A,p(...) or A1,p(...)), but main component is sought (e.g., A(...)).
						Println("Ignoring found element due to ambiguous matching with property of component (Match: " +
							component + tree.PROPERTY_SYNTAX_SUFFIX + ", Component: " + component + ")")
					} else {
						componentStrings = append(componentStrings, candidateString)
						Println("Added string " + candidateString)
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
					processedString = processedString[idx+len(candidateString):]
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
Attempts to extract the component type of a given prefix, and indicates whether it has detected a
properties component. Assumes 0 index for component type symbol.
Note: Only the prefix part of the component should be provided as input (e.g., Cac1[annotation], not Cac1[annotation]{A(actor) I(...) ...}).
Returns identified component type or error if not found.
*/
func extractComponentType(input string) (string, bool, tree.ParsingError) {

	Println("Input:", input)

	ret := ""
	prop := false

	// Filter potential annotations
	if strings.Contains(input, LEFT_BRACKET) {
		input = input[:strings.Index(input, LEFT_BRACKET)]
	}

	for _, v := range tree.IGComponentSymbols {
		// Check whether component is contained - introduces tolerance to excess text (as opposed to exact matching)
		if strings.Contains(input, v) {
			if ret != "" {
				return ret, prop, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_MULTIPLE_COMPONENTS_FOUND, ErrorMessage: "Multiple component specifications found (" + ret + " and " + v + ") " +
					"when parsing component specification '" + input + "'."}
			}
			// Assign identified label
			ret = v
			// Test whether component of concern is a property
			if strings.Contains(input, tree.PROPERTY_SYNTAX_SUFFIX) {
				ret += tree.PROPERTY_SYNTAX_SUFFIX
				prop = true
			}
			// continue iteration to check whether conflicting identification of component (i.e., multiple component labels)
		}
	}
	if ret == "" {
		return "", prop, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_COMPONENT_NOT_FOUND,
			ErrorMessage: "Component specification could not be found in input phrase '" + input + "'."}
	}

	return ret, prop, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Extracts suffix (e.g., ,p1) and annotations (e.g., [ctx=time]), and content from IG-Script-coded input.
It takes component identifier and raw coded information as input, as well as left and right parenthesis symbols (e.g., (,) or {,}).
Returns suffix as first element, annotations string as second, and component content (including identifier, but without suffix and annotations) as third element.
IMPORTANT:
- This function will only extract the suffix and annotation for the first element of a given component type found in the input string.
- This function will not prevent wrongful extraction of property components instead of first-order components. This is handled in #ExtractComponentContent.
TODO: Make this more efficient
*/
func extractSuffixAndAnnotations(component string, propertyComponent bool, input string, leftPar string, rightPar string) (string, string, string, tree.ParsingError) {

	Println("Component:", component)
	Println("Input:", input)
	Println("Property:", propertyComponent)
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
		Println("Found annotation in component:", res)
		// Extract semantic annotation string
		res = res[:len(res)-1]
		pos := strings.Index(strippedInput, res)
		suffix := ""

		if propertyComponent {
			// If component is property, find first position of property indicator
			propIdx := strings.Index(strippedInput[:pos], tree.PROPERTY_SYNTAX_SUFFIX)
			if propIdx == -1 {
				return "", "", "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Property syntax " +
					tree.PROPERTY_SYNTAX_SUFFIX + " (under consideration of potential suffix (e.g., A1,p)) could not be found in input " + strippedInput}
			}
			// Find original component identifier
			leadIdx := strings.Index(component, tree.PROPERTY_SYNTAX_SUFFIX)
			if leadIdx == -1 {
				return "", "", "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Property syntax " +
					tree.PROPERTY_SYNTAX_SUFFIX + " (e.g., A,p) could not be found in input " + strippedInput}
			}
			if propIdx > leadIdx {
				// Extract difference between index in original component and new identifier
				suffix = strippedInput[leadIdx : leadIdx+(propIdx-leadIdx)]
			}
		} else {
			// Component identifier is suppressed if suffix is found
			// Extract suffix (e.g., 1), but remove component identifier
			suffix = strippedInput[len(component):pos]
		}

		// Replace annotations
		extractedContent := strings.ReplaceAll(strippedInput, res, "")
		// Replace suffices
		extractedContent = strings.ReplaceAll(extractedContent, suffix, "")
		Println("Extracted content:", extractedContent)
		// Return suffix and annotations
		return suffix, res, extractedContent, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
	} else {
		Println("No annotations found ...")
		// ... if no annotations are found ...
		// Identifier start position for content
		contentStartPos := strings.Index(strippedInput, leftPar)
		suffix := ""

		if propertyComponent {
			// If component is property, find first position of property indicator
			propIdx := strings.Index(strippedInput[:contentStartPos], tree.PROPERTY_SYNTAX_SUFFIX)
			if propIdx == -1 {
				return "", "", "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Property syntax " +
					tree.PROPERTY_SYNTAX_SUFFIX + " (under consideration of potential suffix (e.g., A1,p)) could not be found in input " + strippedInput}
			}
			// Find original component identifier
			leadIdx := strings.Index(component, tree.PROPERTY_SYNTAX_SUFFIX)
			if leadIdx == -1 {
				return "", "", "", tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Property syntax " +
					tree.PROPERTY_SYNTAX_SUFFIX + " (e.g., A,p) could not be found in input " + strippedInput}
			}
			if propIdx > leadIdx {
				// Extract difference between index in original component and new identifier
				suffix = strippedInput[leadIdx : leadIdx+(propIdx-leadIdx)]
			}
		} else {
			// Component identifier is suppressed if suffix is found
			// Extract suffix (e.g., 1), but remove component identifier
			suffix = strippedInput[len(component):contentStartPos]
		}
		// Does not guard against mistaken choice of property variants of components (e.g., A,p instead of A) - is handled in #ExtractComponentContent.
		reconstructedComponent := strings.Replace(strippedInput, suffix, "", 1)
		Println("Reconstructed statement:", reconstructedComponent)
		// Return only suffix
		return suffix, "", reconstructedComponent, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
	}
}

/*
Parses component based on surrounding parentheses.
*/
func parseComponentWithParentheses(component string, propertyComponent bool, input string) (*tree.Node, tree.ParsingError) {
	return parseComponent(component, propertyComponent, input, LEFT_PARENTHESIS, RIGHT_PARENTHESIS)
}

/*
Parses component based on surrounding braces
*/
func parseComponentWithBraces(component string, propertyComponent bool, input string) (*tree.Node, tree.ParsingError) {
	return parseComponent(component, propertyComponent, input, LEFT_BRACE, RIGHT_BRACE)
}

/*
Generic entry point to parse individual components of a given statement.
Input is component symbol of interest, full statements, as well as delimiting parentheses signaling parsing for atomic
or nested components. Additionally, the parameter propertyComponent indicates whether the parsed component is a property
Returns the parsed node.
*/
func parseComponent(component string, propertyComponent bool, text string, leftPar string, rightPar string) (*tree.Node, tree.ParsingError) {

	Println("Parsing:", component)

	// TODO: For property variants, identify root property and search as to whether embedded midfix exists

	// Extract component (one or multiple occurrences) from input string based on provided component identifier
	componentStrings, err := ExtractComponentContent(component, propertyComponent, text, leftPar, rightPar)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	Println("Components (Count:", len(componentStrings), "):", fmt.Sprint(componentStrings))

	// Initialize output string for parsing
	componentString := ""

	// Node to be populated as return node
	node := &tree.Node{}

	// Synthetically linked ([sAND]) components (if multiple occur in input string)
	if len(componentStrings) > 1 {
		Println("Component combination for component", component)
		Println("Component content", componentStrings)
		r, err := regexp.Compile(COMBINATION_PATTERN_PARENTHESES)
		if err != nil {
			return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_PATTERN_EXTRACTION,
				ErrorMessage: "Error during pattern extraction in combination expression."}
		}

		for i, v := range componentStrings {
			Println("Round: " + strconv.Itoa(i) + ": " + v)

			// Extracts suffix and/or annotation for individual component instance -- must only be used with single component instance!
			componentSuffix, componentAnnotation, componentContent, err := extractSuffixAndAnnotations(component, propertyComponent, v, leftPar, rightPar)
			if err.ErrorCode != tree.PARSING_NO_ERROR {
				return nil, err
			}

			Println("Suffix:", componentSuffix, "(Length:", len(componentSuffix), ")")
			Println("Annotations:", componentAnnotation, "(Length:", len(componentAnnotation), ")")
			Println("Content:", componentContent)

			// Extract and concatenate individual component values but cut leading component identifier
			componentWithoutIdentifier := componentContent[len(component):]
			// Identify whether combination embedded in input string element
			result := r.FindAllStringSubmatch(componentWithoutIdentifier, -1)
			Println("Result of component match:", result)
			Println("Length:", len(result))
			Println("Component string before:", componentWithoutIdentifier)
			if len(result) == 0 {
				leadStripIdx := strings.Index(componentWithoutIdentifier, leftPar)
				if leadStripIdx != -1 {
					// If no combination embedded in combination component, strip leading and trailing parentheses prior to combining
					componentWithoutIdentifier = componentWithoutIdentifier[leadStripIdx+1 : len(componentWithoutIdentifier)-1]
				}
			} // else don't touch, i.e., leave parentheses in string
			Println("Component string after:", componentWithoutIdentifier)

			// Parse first component into node
			if node.IsEmptyOrNilNode() {
				node1, _, err := ParseIntoNodeTree(componentWithoutIdentifier, false, leftPar, rightPar)
				if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
					log.Println("Error when parsing synthetically linked element. Error:", err)
					return nil, err
				}
				// Assign to main node if not populated and new node not nil
				if !node1.IsEmptyOrNilNode() {
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
				if !node2.IsEmptyOrNilNode() {
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
					nodeComb, nodeCombinationError := tree.Combine(node, node2, tree.SAND_BETWEEN_COMPONENTS)
					// Check if combination error has been picked up
					if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
						return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
							ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
					}
					// Explicitly assign component type to top-level node (for completeness) - should be done from within combine function
					nodeComb.ComponentType = component
					// Assign to return node
					node = nodeComb
				}
			}
		}
	} else if len(componentStrings) == 1 {

		Println("Component strings:", componentStrings)

		// Extracts suffix and/or annotation for individual component instance -- must only be used with single component instance!
		componentSuffix, componentAnnotation, componentContent, err := extractSuffixAndAnnotations(component, propertyComponent, componentStrings[0], leftPar, rightPar)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return nil, err
		}

		Println("Suffix:", componentSuffix, "(Length:", len(componentSuffix), ")")
		// Store suffices
		Println("Annotations:", componentAnnotation, "(Length:", len(componentAnnotation), ")")
		// Store annotations
		Println("Content:", componentContent)

		// Single entry (cut prefix)
		componentString = componentContent[strings.Index(componentContent, leftPar):]
		Println("Single component for component", component)
		Println("Component content", componentString)
		// Remove prefix including leading and trailing parenthesis (e.g., Bdir(, )) to extract inner string if not combined
		componentString = componentString[1 : len(componentString)-1]

		node1, _, err := ParseIntoNodeTree(componentString, false, leftPar, rightPar)
		if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
			log.Println("Error when parsing synthetically linked element. Error:", err)
			return nil, err
		}
		// Attach component name to top-level element (will be accessible to children via GetComponentName())
		if !node1.IsEmptyOrNilNode() {
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

	Println("Component Identifier: " + component)
	Println("Full string: " + componentString)

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
