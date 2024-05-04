package parser

import (
	"IG-Parser/core/tree"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

/*
This file provides the generic statement parsing functionality and is the
primary invocation point for the endpoint package.
The parsing considers all parsing features including basic component parsing,
component combination, nested statements, nested statement combinations
as well as component pairs.
- Invokes functionality contained in IGComponentCombinationParser.go and
IGComponentParser.go.
*/

/*
Parses statement tree from input string. Returns statement tree, and error.
If parsing is successful, error code tree.PARSING_NO_ERROR is returned, else
other context-specific codes are returned.
*/
func ParseStatement(text string) ([]*tree.Node, tree.ParsingError) {

	s := tree.Statement{}

	Println("INITIATING STATEMENT PARSING ...\nProcessing input statement: ", text)

	// Validate input string first with respect to parentheses ...
	err := validateInput(text, LEFT_PARENTHESIS, RIGHT_PARENTHESIS)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return []*tree.Node{}, err
	}
	// ... and braces
	err = validateInput(text, LEFT_BRACE, RIGHT_BRACE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return []*tree.Node{}, err
	}

	// Now extract component-only expressions, nested statements, statement combinations, as well as component pair combinations (Note: only processed at the end of function)
	compAndNestedStmts, err := separateComponentsNestedStatementsCombinationsAndComponentPairs(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return []*tree.Node{}, err
	}
	Println("Identified statement element patterns prior to parsing:\n" +
		" - individual atomic components (e.g., 'A(actor)') including component combinations (e.g., 'Bdir((left [AND] right))') ( --> element [0]): \n" +
		" ==> " + fmt.Sprint(compAndNestedStmts[0]) + " --> Count: " + strconv.Itoa(len(compAndNestedStmts[0])) + "\n" +
		" - nested statements/components (e.g., 'Bdir{ A(nestedA) I(nestedI) }'). Note: may also contain invalid component pairs - to be filtered later ( --> element [1]): \n" +
		" ==> " + fmt.Sprint(compAndNestedStmts[1]) + " --> Count: " + strconv.Itoa(len(compAndNestedStmts[1])) + "\n" +
		" - nested component combinations (e.g., 'Cac{ Cac{A(leftNestedA) I(leftNestedI)} [XOR] Cac{A(rightNestedA) I(rightNestedI)} }') ( --> element [2]): \n" +
		" ==> " + fmt.Sprint(compAndNestedStmts[2]) + " --> Count: " + strconv.Itoa(len(compAndNestedStmts[2])) + "\n" +
		" - component pair combinations (e.g., '{ Cac{A(leftNestedA) I(leftNestedI)} [XOR] Cac{A(rightNestedA) I(rightNestedI)} }') ( --> element [3]): \n" +
		" ==> " + fmt.Sprint(compAndNestedStmts[3]) + " --> Count: " + strconv.Itoa(len(compAndNestedStmts[3])))
	Println("Complete return structure: " + fmt.Sprint(compAndNestedStmts))

	// Extract component-only statement and override input (e.g., A(content))
	if len(compAndNestedStmts[0]) > 0 {
		// Assign array elements as text (for later string parsing)
		text = compAndNestedStmts[0][0]
	} else {
		// else just reset input text (so no basic element is parsed)
		text = ""
	}
	// Note: text variable is used to extract logical operators (if present) for multiple component-level nested statements!
	// Extract potential nested statements (e.g., Cac{ content } )
	nestedStmts := compAndNestedStmts[1]
	if len(nestedStmts) == 0 {
		Println("No nested statements found.")
	}
	// Extract potential statement combinations (e.g., Cac{ Cac{ content } [XOR] Cac{ content } })
	nestedCombos := compAndNestedStmts[2]
	if len(nestedCombos) == 0 {
		Println("No nested statement combination candidates found.")
	}

	Println("Statement before parsing: " + s.String())

	Println("Parsing basic statement ...")

	// Process basic components and component combinations
	if text != "" {
		Println("Text to be parsed: " + text)

		// Now parsing on component level
		_, remainingText, outErr := parseBasicStatement(text, &s)
		if outErr.ErrorCode != tree.PARSING_NO_ERROR {
			// Populate return structure
			ret := []*tree.Node{&tree.Node{Entry: &s}}
			return ret, outErr
		}

		// Reorganize tree by shifting private nodes into PrivateNode fields of components and removing them from statement tree
		ProcessPrivateComponentLinkages(&s, false)

		Println("Basic statement: " + s.String())

		// Substitute text with remaining parts (the one that have not been parsed as part of the basic component parsing)
		text = remainingText
		Println("Remaining text after basic component parsing: " + text)
	}

	Println("Testing for nested combinations in " + fmt.Sprint(nestedCombos))

	// Process nested statement combinations
	if len(nestedCombos) > 0 {
		Println("Found nested combination(s) ...")
		// Iterate through all combinations for fine-granular error handling
		for _, nestedCombo := range nestedCombos {
			Println("Attempting to process nested combination " + nestedCombo)
			err = parseNestedStatementCombination(&s, nestedCombo)
			if err.ErrorCode == tree.PARSING_ERROR_NIL_ELEMENT {
				// Shift to regular nested statement if parsing as combo failed (Regex is too coarse-grained and favors combinations before fine-grained parsing)
				nestedStmts = append(nestedStmts, err.ErrorIgnoredElements...)
				Println("Reclassifying statement as nested statement (as opposed to nested combination) ...")
			} else if err.ErrorCode != tree.PARSING_NO_ERROR {
				// Populate return structure
				ret := []*tree.Node{&tree.Node{Entry: &s}}
				return ret, err
			}
		}
	}

	Println("Testing for nested statements in " + fmt.Sprint(nestedStmts))

	// Process nested statements
	if len(nestedStmts) > 0 {
		Println("Found nested statements ... (Nested statements: ", nestedStmts, ")")

		// Will be populated if statement itself contains logical operator (uses 'pure' versions of logical operators and
		// defaults to tree.AND if none is found)
		detectedLogicalOperator := ""
		if len(nestedStmts) > 1 {
			// Check whether explicit logical operators are specified and in inject those.
			Println("Found more than one nested statements in ", nestedStmts, " - testing for logical linkage "+
				"of nested statements in input '"+text+"'.")
			if strings.Contains(text, tree.AND_BRACKETS) {
				detectedLogicalOperator = tree.AND
			}
			if strings.Contains(text, tree.XOR_BRACKETS) {
				if detectedLogicalOperator != "" {
					return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS,
						ErrorMessage: "Detected multiple logical operators (" + detectedLogicalOperator + " and " + tree.XOR +
							") on given parsing level. Please revise your coding with respect to indication of precedence."}
				}
				detectedLogicalOperator = tree.XOR
			}
			if strings.Contains(text, tree.OR_BRACKETS) {
				if detectedLogicalOperator != "" {
					return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS,
						ErrorMessage: "Detected multiple logical operators (" + detectedLogicalOperator + " and " + tree.OR +
							") on given parsing level. Please revise your coding with respect to indication of precedence."}
				}
				detectedLogicalOperator = tree.OR
			}
		}

		err = parseNestedStatements(&s, nestedStmts, detectedLogicalOperator)
		// Check whether nested statements have been ignored entirely
		if err.ErrorCode == tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS {
			// Populate return structure
			ret := []*tree.Node{&tree.Node{Entry: &s}}
			return ret, err
		}
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			// Populate return structure
			ret := []*tree.Node{&tree.Node{Entry: &s}}
			return ret, err
		}
		// Process potential private nodes for complex components
		ProcessPrivateComponentLinkages(&s, true)

	}

	Println("Statement (after assigning nested elements, before expanding statement with paired components):\n" + s.String())

	// Process component pair combinations and extrapolate into multiple statements
	if len(compAndNestedStmts[3]) > 0 {
		if len(compAndNestedStmts[3]) > 1 {
			return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_MULTIPLE_COMPONENT_PAIRS_ON_SAME_LEVEL,
				ErrorMessage: "The coding contains multiple component pairs on the same nesting level. " +
					"There should only be one component pair expression on a given level (e.g., { left options [OR] right options } for simple cases, " +
					"or for complex cases with more options, e.g., { first option [AND] { second option [OR] third option } }, etc. " +
					"Consider getting in touch with the maintainer in case you believe that you encounter a case that requires multiple separate component pairs on a given nesting level. " +
					"The expressions of concern are: '" + strings.Join(compAndNestedStmts[3], " , ") + "'",
				ErrorIgnoredElements: compAndNestedStmts[3]}
		} else {
			// If one component pair combination on a given nesting level, extrapolate (may contain nested pair combination (e.g., { left [AND] { right [XOR] alsoRight }})
			extrapolatedStmts, err2 := extrapolateStatementWithPairedComponents(&s, compAndNestedStmts[3])
			if err2.ErrorCode != tree.PARSING_NO_ERROR {
				return extrapolatedStmts, err2
			}
			Println("Final statements (with extrapolation): " + tree.PrintNodes(extrapolatedStmts))

			return extrapolatedStmts, err2
		}
	} else {
		Println("No expansion of statement necessary (no component pair combinations in input)")
	}

	// Return error indicating that statement is empty
	if s.IsEmpty() {
		Println("Statement is empty - returning error " + tree.PARSING_ERROR_EMPTY_STATEMENT)
		return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_EMPTY_STATEMENT}
	}

	// Else return wrapped statement (no extrapolation included)
	return []*tree.Node{&tree.Node{Entry: &s}}, err
}

/*
Separates nested statement expressions (including component prefix)
from individual components (including combinations of components).
Returns multi-dim array, with element [0] containing component-only statement (no nested structure),
and element [1] containing nested statements (potentially multiple),
and element [2] containing potential component-level statement combinations,
and element [3] containing component pairs (that need to be extrapolated into entire separate statements).
*/
func separateComponentsNestedStatementsCombinationsAndComponentPairs(statement string) ([][]string, tree.ParsingError) {

	// Prepare return structure
	ret := make([][]string, 4)

	// Identify all component pair combinations (e.g., {I(something) Bdir(something) [XOR] I(something else) Bdir(something else)})
	pairCombos, err := identifyComponentPairCombinations(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	// Identify all nested statements in a very coarse-grained match (this may overlap with already identified component pairs due to overlapping syntax
	// Component nesting syntax: Cac{ A(actor) I(action) }
	// Component combination syntax: Cac{ Cac{A(leftNestedA) I(leftNestedI)} [XOR] Cac{A(rightNestedA) I(rightNestedI)} }
	// Component pair combination syntax: { Cac{A(leftNestedA) I(leftNestedI)} [XOR] Cac{A(rightNestedA) I(rightNestedI)} }
	nestedStmts, err := identifyNestedStatements(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	// Contains complete nested statements (with prefix)
	completeNestedStmts := []string{}

	// Holds candidates for nested combinations
	nestedCombos := []string{}

	// keeps track of filtered pair candidates
	skippedPairs := []string{}
	acceptedPairs := []string{}

	if len(nestedStmts) > 0 {

		// Iterate through identified nested statements (if any) and remove those from statement
		for _, v := range nestedStmts {
			// Extract statements of structure (e.g., Cac{ LEFT [AND] RIGHT }) -
			// Note: component prefix is necessary for combinations and single nested statements; not allowed in component pair combinations
			// Use of terminated statements is important to capture complete nested statements (prefiltering before guarantees nested structures)
			r2, err2 := regexp.Compile(NESTED_COMBINATIONS_TERMINATED)
			if err2 != nil {
				Println("Error in regex compilation: ", err2.Error())
				return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Error in Regular Expression compilation. Error: " + err2.Error()}
			}
			nestedCombinationCandidates := r2.FindAllString(v, -1)
			if len(nestedCombinationCandidates) > 0 {

				if len(pairCombos) > 0 {
				INNER:
					for _, pair := range pairCombos {

						// Check whether pair to be checked has already been discarded
						for _, skippedPair := range skippedPairs {
							if skippedPair == pair {
								continue INNER
							}
						}

						// If identified pair candidate is longer than nested combo ...
						// If same length, it must have more content, since combination is always prefixed
						if len(pair) >= len(v) {
							// ... then test whether combo is substring ...
							if strings.Contains(pair, v) {
								// ... and if this is the case take candidate pair (i.e., remove from statement) ...
								statement = strings.ReplaceAll(statement, pair, "")
								Println("Confirmed candidate for component pair combination:", pair)
							}
						} else {
							// ... else if combo is equal or longer

							// ... then test whether pair is substring ...
							if strings.Contains(v, pair) {
								// ... and if this is the case take combo (i.e., remove from statement) ...
								statement = strings.ReplaceAll(statement, v, "")
								// ... save combination for parsing ...
								// ... but test against potential duplicate in combos (repeated detection or reclassification)
								found := false
								for _, elem := range nestedCombos {
									if elem == v {
										found = true
									}
								}
								if !found {
									nestedCombos = append(nestedCombos, v)
								}
								Println("Confirmed candidate for statement combination:", v)
								// ... and skip pair
								skippedPairs = append(skippedPairs, pair)
							}
						}

					}
				} else {
					// no pairs found - only consider combinations

					// Identified combination of component-level nested statements
					// Save for parsing as combination
					nestedCombos = append(nestedCombos, v)
					Println("Added candidate for statement combination:", v)

					// Remove nested statement combination from overall statement
					statement = strings.ReplaceAll(statement, v, "")
				}
			} else {
				// Check for pairs combo

				found := false

				if len(pairCombos) > 0 {

					// ... If no component combination candidates, but component pair combinations found ...
				PAIRSONLY:
					for _, pair := range pairCombos {

						found = false
						// Check if currently processed nested candidate matches pair
						if pair == v {
							found = true
						}
						if found {
							// Remove component pair combination from overall statement (if present)
							statement = strings.ReplaceAll(statement, v, "")
							// saving for downstream processing is done below (every string present in originally detected collection is deemed found if not filtered) ...
							break PAIRSONLY
						}

					}
				}

				// Check for simple nesting
				if !found {
					// ... else deem it single nested statement

					// Append extracted nested statements including component prefix (e.g., Cac{ A(stuff) ... })
					// Note: may be incorrect, since no checking for leading component - will be caught during deep parsing
					completeNestedStmts = append(completeNestedStmts, v)
					Println("Added candidate for single nested statement (to be checked during deep parsing):", v)

					// Remove nested statement from overall statement
					statement = strings.ReplaceAll(statement, v, "")

					// Remove pairs that are contained in fragments added as individual nested components
					for _, v2 := range pairCombos {
						if strings.Contains(v, v2) {
							// Annotate as skipped pairs
							skippedPairs = append(skippedPairs, v2)
						}
					}

				}
			}
		}
		// Assign nested statements if found
		ret[1] = completeNestedStmts
		ret[2] = nestedCombos

		// Handling of component pair combinations

		// Filter pairs from pairCombos
		for _, originalPair := range pairCombos {
			found := false
			for _, removedPair := range skippedPairs {
				if removedPair == originalPair {
					found = true
				}
			}
			if !found {
				acceptedPairs = append(acceptedPairs, originalPair)
			}
		}

		// Save remaining pair combinations
		ret[3] = acceptedPairs
		// Doublecheck that all identified component pair combinations are removed from input statement
		for _, v := range acceptedPairs {
			statement = strings.ReplaceAll(statement, v, "")
		}
		Println("Remaining non-nested input statement (without nested elements): " + statement)
	} else {
		Println("No nested statement found in input: " + statement)

		// Assign potential pair combinations in case no other nested statements were found
		ret[3] = pairCombos
		// Remove from input statement
		for _, v := range pairCombos {
			statement = strings.ReplaceAll(statement, v, "")
		}
		if len(pairCombos) > 0 {
			Println("Remaining input statement (after removal of pair combinations): " + statement)
		}
	}

	// Assign (remaining) component-only string to first element of returned array (if nothing left, first element will be nil)
	if statement != "" {
		ret[0] = []string{statement}
	}

	// Return combined structure
	return ret, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Generically identifies any nested statements for subsequent deep detection (component-level, combinations, combination pairs).

Used by #separateComponentsNestedStatementsCombinationsAndComponentPairs.
*/
func identifyNestedStatements(statement string) ([]string, tree.ParsingError) {

	// Extract any nested statements from input string
	nestedStatements, err := extractComponentContent("", true, statement, LEFT_BRACE, RIGHT_BRACE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	Println("Found nested statements: " + fmt.Sprint(nestedStatements))

	return nestedStatements, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Handles parsing errors centrally - easier to refine. Used by #parseBasicStatement.
*/
func handleParsingError(component string, err tree.ParsingError) tree.ParsingError {

	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND {
		Println("Error when parsing component ", component, ": ", err)
		return err
	}

	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Parses nested statements (but not combinations) and attaches those to the top-level statement.
Uses given logical operator to link to existing statements (takes only tree.OR, tree.AND, and tree.XOR - no brackets).
Returns an error other than tree.PARSING_NO_ERROR if issues during parsing.
Specific errors:
Returns err tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS if parsing did not pose problems per se,
but elements have been ignored during parsing (warranting syntax review). In this case, the statements of concern
are returned in a string array contained in the error object.
*/
func parseNestedStatements(stmtToAttachTo *tree.Statement, nestedStmts []string, logicalOperator string) tree.ParsingError {

	// Copy reference statement for comparison (to check whether modification took place based on parsed element)
	cachedStmtPriorToNestedParsing := stmtToAttachTo.String()

	// Array to keep track of ignored statements
	nestedStmtsNotConsideredDuringParsing := make([]string, 0)

	for _, v := range nestedStmts {

		Println("Processing nested statement: ", v)
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
			Println("Identified nested Attributes Property")
			component = tree.ATTRIBUTES_PROPERTY
			isProperty = true
		} else if strings.HasPrefix(prefix, tree.DIRECT_OBJECT_PROPERTY) ||
			(strings.HasPrefix(prefix, tree.DIRECT_OBJECT) && strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX)) {
			Println("Identified nested Direct Object Property")
			component = tree.DIRECT_OBJECT_PROPERTY
			isProperty = true
		} else if strings.HasPrefix(prefix, tree.DIRECT_OBJECT) && !strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX) {
			Println("Identified nested Direct Object")
			component = tree.DIRECT_OBJECT
		} else if strings.HasPrefix(prefix, tree.INDIRECT_OBJECT_PROPERTY) ||
			(strings.HasPrefix(prefix, tree.INDIRECT_OBJECT) && strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX)) {
			Println("Identified nested Indirect Object Property")
			component = tree.INDIRECT_OBJECT_PROPERTY
			isProperty = true
		} else if strings.HasPrefix(prefix, tree.INDIRECT_OBJECT) && !strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX) {
			Println("Identified nested Indirect Object")
			component = tree.INDIRECT_OBJECT
		} else if strings.HasPrefix(prefix, tree.ACTIVATION_CONDITION) {
			Println("Identified nested Activation Condition")
			component = tree.ACTIVATION_CONDITION
		} else if strings.HasPrefix(prefix, tree.EXECUTION_CONSTRAINT) {
			Println("Identified nested Execution Constraint")
			component = tree.EXECUTION_CONSTRAINT
		} else if strings.HasPrefix(prefix, tree.CONSTITUTED_ENTITY_PROPERTY) ||
			(strings.HasPrefix(prefix, tree.CONSTITUTED_ENTITY) && strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX)) {
			Println("Identified nested Constituted Entity Property")
			component = tree.CONSTITUTED_ENTITY_PROPERTY
			isProperty = true
		} else if strings.HasPrefix(prefix, tree.CONSTITUTING_PROPERTIES_PROPERTY) ||
			(strings.HasPrefix(prefix, tree.CONSTITUTING_PROPERTIES) && strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX)) {
			Println("Identified nested Constituting Properties Property")
			component = tree.CONSTITUTING_PROPERTIES_PROPERTY
			isProperty = true
		} else if strings.HasPrefix(prefix, tree.CONSTITUTING_PROPERTIES) && !strings.Contains(prefix, tree.PROPERTY_SYNTAX_SUFFIX) {
			Println("Identified nested Constituting Properties")
			component = tree.CONSTITUTING_PROPERTIES
		} else if strings.HasPrefix(prefix, tree.OR_ELSE) {
			Println("Identified nested Or Else")
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

		// Parse actual content wrapped in nested component (e.g., content inside Cac{ ... })
		stmt, errStmt := ParseStatement(v[strings.Index(v, LEFT_BRACE)+1 : strings.LastIndex(v, RIGHT_BRACE)])
		if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
			fmt.Println("Error when parsing nested statements: ", errStmt.Error())
			return errStmt
		}
		if len(stmt) > 1 {
			fmt.Println("Unhandled case: Multiple decomposed statements in nested component ...", stmt)
			return errStmt
		}

		// Return first statement (wrapped into node) - (since individual statement)
		stmtNode := stmt[0]
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
			Println("Attaching nested attributes property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.AttributesPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.AttributesPropertyComplex, stmtNode, logicalOperator)
		case tree.DIRECT_OBJECT_PROPERTY:
			Println("Attaching nested direct object property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.DirectObjectPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.DirectObjectPropertyComplex, stmtNode, logicalOperator)
		case tree.DIRECT_OBJECT:
			Println("Attaching nested direct object to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.DirectObjectComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.DirectObjectComplex, stmtNode, logicalOperator)
		case tree.INDIRECT_OBJECT_PROPERTY:
			Println("Attaching nested indirect object property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.IndirectObjectPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.IndirectObjectPropertyComplex, stmtNode, logicalOperator)
		case tree.INDIRECT_OBJECT:
			Println("Attaching nested indirect object to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.IndirectObjectComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.IndirectObjectComplex, stmtNode, logicalOperator)
		case tree.ACTIVATION_CONDITION:
			Println("Attaching nested activation condition to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ActivationConditionComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ActivationConditionComplex, stmtNode, logicalOperator)
		case tree.EXECUTION_CONSTRAINT:
			Println("Attaching nested execution constraint to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ExecutionConstraintComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ExecutionConstraintComplex, stmtNode, logicalOperator)
		case tree.CONSTITUTED_ENTITY_PROPERTY:
			Println("Attaching nested constituted entity property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutedEntityPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutedEntityPropertyComplex, stmtNode, logicalOperator)
		case tree.CONSTITUTING_PROPERTIES_PROPERTY:
			Println("Attaching nested constituting properties property to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutingPropertiesPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesPropertyComplex, stmtNode, logicalOperator)
		case tree.CONSTITUTING_PROPERTIES:
			Println("Attaching nested constituting properties to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.ConstitutingPropertiesComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesComplex, stmtNode, logicalOperator)
		case tree.OR_ELSE:
			Println("Attaching nested or else to higher-level statement")
			// Assign nested statement to higher-level statement
			stmtToAttachTo.OrElse, nodeCombinationError = attachComplexComponent(stmtToAttachTo.OrElse, stmtNode, logicalOperator)
		}
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: nodeCombinationError.ErrorCode, ErrorMessage: "Error when merging substatements into statement. Error: " +
				nodeCombinationError.ErrorMessage}
		}

		// Check if the iterated nested statement has been ignored entirely --> indicates failed detection as nested (as opposed to mere parsing problem)
		if cachedStmtPriorToNestedParsing == stmtToAttachTo.String() {
			Println("Nested statement has not been considered during parsing:", v)
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
Parses a given nested statement combination and attaches it to the corresponding component of the provided top-level statement.
Returns error code tree.PARSING_NO_ERROR if no problem occurs.
The error code tree.PARSING_ERROR_NIL_ELEMENT indicates inability to extract the combination structure for a given input.
In this case the violating statement is passed in the string array provided as part of the error object.
Error tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION and tree.PARSING_ERROR_INVALID_TYPES_IN_NESTED_STATEMENT_COMBINATION
point to an invalid combination of different component types (e.g., Cac and Bdir)
Error tree.PARSING_ERROR_INVALID_COMBINATION points to syntactic issues during combination parsing.
*/
func parseNestedStatementCombination(stmtToAttachTo *tree.Statement, nestedCombo string) tree.ParsingError {

	// Default error for node combination - can generally only be overridden by detected invalid component combinations
	nodeCombinationError := tree.NodeError{ErrorCode: tree.TREE_NO_ERROR}

	Println("Found nested statement combination candidate", nestedCombo)

	combo, _, errStmt := ParseIntoNodeTree(nestedCombo, false, LEFT_BRACE, RIGHT_BRACE)
	if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
		fmt.Print("Error when parsing nested statements: " + errStmt.ErrorCode)
		return errStmt
	}

	// Check whether all leaves have the same prefix
	flatCombo := tree.Flatten(combo.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES))
	sharedPrefix := ""
	for _, node := range flatCombo {
		if node.Entry == nil {
			// Parsing did not work (incomplete combination (e.g., embedded logical operator, but on wrong level)); simply return error and violating entry
			// (i.e., statement that has not been parseable as combination -- e.g., for retrying as a regular nested statement)
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_NIL_ELEMENT, ErrorMessage: "Nested combination returned nil element.",
				ErrorIgnoredElements: []string{nestedCombo}}
		}
		entry := node.Entry.(string)
		Println("Entry to parse for component type: " + entry)
		// Extract prefix (i.e., component type) for node, but check whether it contains nested statement
		if strings.Index(entry, LEFT_BRACE) == -1 {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMBINATION, ErrorMessage: "Element in combination of nested statement does not contain nested statement. Element of concern: " + entry}
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
			//continue
		}
		// Check if it deviates from previously cached element
		if prefix != sharedPrefix {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_TYPES_IN_NESTED_STATEMENT_COMBINATION,
				ErrorMessage: "Invalid combination of component-level nested statements. Expected component: " +
					sharedPrefix + ", but found: " + prefix}
		}
	}

	// Parse all entries in tree from string to statement (walks through entire tree linked to node)
	err := combo.ParseAllEntries(func(oldValue string) (*tree.Statement, tree.ParsingError) {

		// Check whether the combination element contains a nested structure ...
		tempComponentType := oldValue
		if strings.Contains(oldValue, LEFT_BRACE) {
			// ... and remove the nested element prior to parsing
			tempComponentType = oldValue[:strings.Index(oldValue, LEFT_BRACE)]
		}

		// Extract component type (after stripping potential nested statements)
		compType, prop, err := extractComponentType(tempComponentType)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return &tree.Statement{}, err
		}
		// Extracting suffices and annotations
		suffix, annotation, content, err := extractSuffixAndAnnotations(compType, prop, oldValue, LEFT_BRACE, RIGHT_BRACE)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			fmt.Println("Error during extraction of suffices and annotations of component '" + compType + "': " + err.ErrorCode)
			return &tree.Statement{}, err
		}

		Println("Nested Combo Stmt Suffix:", suffix)
		Println("Nested Combo Stmt Annotation:", annotation)
		Println("Nested Combo Stmt Content:", content)

		stmt, errStmt := ParseStatement(oldValue[strings.Index(oldValue, LEFT_BRACE)+1 : strings.LastIndex(oldValue, RIGHT_BRACE)])
		if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
			return stmt[0].Entry.(*tree.Statement), errStmt
		}
		return stmt[0].Entry.(*tree.Statement), tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
	})
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return err
	}

	//TODO: Check whether combinations are actually filled, or just empty nodes (e.g., { Cac{ A(), I(), Cex() } [AND] Cac{ A(), I(), Cex() } })

	Println("Assigning nested tree structure", combo.String())

	// Assign component type name to combination (for proper retrieval and identification as correct type)
	combo.ComponentType = sharedPrefix
	Println("Combo component prefix:", sharedPrefix)

	// Checks are ordered with property variants (e.g., Bdir,p) before component variants (e.g., Bdir) to avoid wrong match

	if strings.HasPrefix(sharedPrefix, tree.ATTRIBUTES_PROPERTY) {
		Println("Attaching nested attributes property to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.AttributesPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.AttributesPropertyComplex, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}
	if strings.HasPrefix(sharedPrefix, tree.DIRECT_OBJECT_PROPERTY) {
		Println("Attaching nested direct object property to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.DirectObjectPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.DirectObjectPropertyComplex, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}
	if strings.HasPrefix(sharedPrefix, tree.DIRECT_OBJECT) {
		Println("Attaching nested direct object to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.DirectObjectComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.DirectObjectComplex, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}
	if strings.HasPrefix(sharedPrefix, tree.INDIRECT_OBJECT_PROPERTY) {
		Println("Attaching nested indirect object property to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.IndirectObjectPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.IndirectObjectPropertyComplex, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}
	if strings.HasPrefix(sharedPrefix, tree.INDIRECT_OBJECT) {
		Println("Attaching nested indirect object to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.IndirectObjectComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.IndirectObjectComplex, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}
	if strings.HasPrefix(sharedPrefix, tree.ACTIVATION_CONDITION) {
		Println("Attaching nested activation condition to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.ActivationConditionComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ActivationConditionComplex, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}
	if strings.HasPrefix(sharedPrefix, tree.EXECUTION_CONSTRAINT) {
		Println("Attaching nested execution constraint to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.ExecutionConstraintComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ExecutionConstraintComplex, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}
	if strings.HasPrefix(sharedPrefix, tree.CONSTITUTED_ENTITY_PROPERTY) {
		Println("Attaching nested constituted entity property to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.ConstitutedEntityPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutedEntityPropertyComplex, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}
	if strings.HasPrefix(sharedPrefix, tree.CONSTITUTING_PROPERTIES_PROPERTY) {
		Println("Attaching nested constituting properties property to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.ConstitutingPropertiesPropertyComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesPropertyComplex, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}
	if strings.HasPrefix(sharedPrefix, tree.CONSTITUTING_PROPERTIES) {
		Println("Attaching nested constituting properties to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.ConstitutingPropertiesComplex, nodeCombinationError = attachComplexComponent(stmtToAttachTo.ConstitutingPropertiesComplex, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}
	if strings.HasPrefix(sharedPrefix, tree.OR_ELSE) {
		Println("Attaching nested or else to higher-level statement")
		// Assign nested statement to higher-level statement
		stmtToAttachTo.OrElse, nodeCombinationError = attachComplexComponent(stmtToAttachTo.OrElse, combo, "")
		// Process error and return
		if nodeCombinationError.ErrorCode != tree.TREE_NO_ERROR {
			return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
				ErrorMessage: "Invalid combination of component types of different kinds. Error: " + nodeCombinationError.ErrorMessage}
		} else {
			return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}
	}

	// Should not occur. Assignment should work out - in which case this error is not met
	return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Nested combination has not been attached to statement for unknown reasons. Node:" + combo.String()}
}

/*
Attach complex component to tree structure under consideration of existing nodes in target tree structure.
Input:
- Node of the parent tree to attach to
- Node to attach
- Logical operator with which node should be added if a node already exists. Only takes tree.AND, tree.XOR and tree.OR (no brackets).

Used by #parseNestedStatementCombination.
*/
func attachComplexComponent(nodeToAttachTo *tree.Node, nodeToAttach *tree.Node, logicalOperator string) (*tree.Node, tree.NodeError) {

	Println("Attaching nested complex component to higher-level statement (with logical linkage '" + logicalOperator + "')")

	// Identify correct logical operator and check for proper version
	if logicalOperator == "" {
		logicalOperator = tree.AND
	} else {
		switch logicalOperator {
		case tree.AND:
			logicalOperator = tree.AND
		case tree.XOR:
			logicalOperator = tree.XOR
		case tree.OR:
			logicalOperator = tree.OR
		default:
			return nil, tree.NodeError{ErrorCode: tree.PARSING_ERROR_UNKNOWN_LOGICAL_OPERATOR, ErrorMessage: "Detected unknown logical operator during processing: " + logicalOperator +
				" - please review your coding accordingly. Note that the use of the bracket versions (" + tree.AND_BRACKETS + ", " + tree.XOR_BRACKETS + ", " + tree.OR_BRACKETS + " is not supported)."}
		}
	}

	// Assign nested statement to higher-level statement

	// If already a statement assignment to complex element, ...
	if nodeToAttachTo != nil {
		// ... combine both
		newNode, err := tree.Combine(nodeToAttachTo, nodeToAttach, logicalOperator)
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
Identifies component pair combinations in input text (e.g., '{I(monitor) Bdir(one thing) [XOR] I(enforce) Bdir(the other)}').
Returns identified component pair combinations as string array.
*/
func identifyComponentPairCombinations(statement string) ([]string, tree.ParsingError) {

	r, err := regexp.Compile(COMPONENT_PAIR_COMBINATIONS)
	if err != nil {
		Println("Error in regex compilation: ", err.Error())
		return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Error in Regular Expression compilation. Error: " + err.Error()}
	}

	// Extract all matches as string array
	pairs := r.FindAllString(statement, -1)

	Println("Found", len(pairs), "component pair combination/s: ", pairs)

	return pairs, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Process pair combinations and extrapolate individual statements and populate with content from atomic input statement.
*/
func extrapolateStatementWithPairedComponents(s *tree.Statement, pairs []string) ([]*tree.Node, tree.ParsingError) {

	// Parse all elements of tree structure
	extrapolatedPairStmts := []*tree.Node{}
	for k, v := range pairs {
		Println("Extrapolation Iteration ", k)

		// Convert individual pair into node structure
		idvStmt, _, err := ParseIntoNodeTree(v, true, LEFT_BRACE, RIGHT_BRACE)
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			return nil, err
		}

		// Extract leaves for conversion
		leaves := idvStmt.GetLeafNodes(true)
		// Test for unsuccessful extraction of leaves - this results in empty array with "nil" entry.
		// This is caused by lack of logical operator on component pair level, but logical operators contained in components in either pair
		// Example: '{ I(action1 [XOR] action2) Bdir(object1) and I(action3) Bdir(object2) }' <-- note the missing logical operator between pairs
		if len(leaves) == 0 || (len(leaves) == 1 && len(leaves[0]) == 1 && leaves[0][0].Entry == nil) {
			return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_PAIR, ErrorMessage: "Invalid component pair found (Missing logical operator?). " +
				"Please review expression '" + v + "'."}
		}
		// More leaf elements in component pair than expected
		if len(leaves) > 1 {
			return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_TOO_MANY_NODES, ErrorMessage: "Too many nodes generated for atomic statement."}
		}

		// Parse leaves into statement structure individually (i.e., atomic statements, i.e., {I(enforce] [XOR] I(monitor)} --> one for enforce, one for monitor),
		// complete with original atomic structure components (the ones that were atomic in the first place (e.g., A(enforcer)) and reattach to extrapolated statement
		for _, v2 := range leaves[0] {

			// Parse content of tree
			tpNode, err := ParseStatement(v2.Entry.(string))
			if err.ErrorCode != tree.PARSING_NO_ERROR {
				Println("Error when parsing statement: ", err, "; expression for which parsing failed:", v2.Entry)
				return nil, err
			}

			if len(tpNode) > 1 {
				return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_TOO_MANY_NODES, ErrorMessage: "Expecting single node/statement, as opposed to multiple. Aborting extrapolation of statements"}
			}

			// Assign node embedding statement to higher-level node containing statement collection (the one with logical operator)
			tpNode[0].Parent = v2.Parent

			// Complete decomposed partial statement with parsed linear statement (can only be one statement in decomposed pair combinations)
			tpNode[0].Entry = tree.CopyComponentsFromStatement(tpNode[0].Entry.(*tree.Statement), s)

			// Assign statement to statement tree (top-level extrapolated structure)
			v2.Parent = idvStmt

			// Replace Entry content
			v2.Entry = tpNode
		}
		// Store parsed statement to return structure
		extrapolatedPairStmts = append(extrapolatedPairStmts, idvStmt)
	}

	Println("Number of elements: ", len(extrapolatedPairStmts))

	return extrapolatedPairStmts, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
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
		Println(msg)
		return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IMBALANCED_PARENTHESES, ErrorMessage: msg}
	}

	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}
