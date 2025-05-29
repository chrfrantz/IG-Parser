package parser

import (
	"IG-Parser/core/shared"
	"IG-Parser/core/tree"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

/*
This file includes all functionality relevant to parsing of individual components.
- Invoked by parser.IGStatementParser.go, and implicitly tested via
parser.IGStatementParser_test.go.
*/

/*
Parse basic statements identified as part of #separateComponentsNestedStatementsCombinationsAndComponentPairs.
Takes plain string and statement reference as input. Generates statement if no statement reference (i.e., nil) is passed.
Returns statement embedded in node, as well as remaining part of original input string that has not been parsed.
Note: If returning an error, the identified part of the text associated with the last parsed component is returned (to simplify diagnostics).
*/
func parseBasicStatement(text string, s *tree.Statement) ([]tree.Node, string, tree.ParsingError) {

	// Check whether statement is passed, else create new one
	if s == nil {
		s = &tree.Statement{}
	}
	// Initial full string content: keeps track of remaining string content
	remainingString := text

	result, componentText, err := parseAttributes(text)
	outErr := handleParsingError(tree.ATTRIBUTES, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.Attributes = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseAttributesProperty(text)
	outErr = handleParsingError(tree.ATTRIBUTES_PROPERTY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.AttributesPropertySimple = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseDeontic(text)
	outErr = handleParsingError(tree.DEONTIC, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.Deontic = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseAim(text)
	outErr = handleParsingError(tree.AIM, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.Aim = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseDirectObject(text)
	outErr = handleParsingError(tree.DIRECT_OBJECT, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.DirectObject = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseDirectObjectProperty(text)
	outErr = handleParsingError(tree.DIRECT_OBJECT_PROPERTY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.DirectObjectPropertySimple = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseIndirectObject(text)
	outErr = handleParsingError(tree.INDIRECT_OBJECT, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.IndirectObject = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseIndirectObjectProperty(text)
	outErr = handleParsingError(tree.INDIRECT_OBJECT_PROPERTY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.IndirectObjectPropertySimple = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseActivationCondition(text)
	outErr = handleParsingError(tree.ACTIVATION_CONDITION, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.ActivationConditionSimple = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseExecutionConstraint(text)
	outErr = handleParsingError(tree.EXECUTION_CONSTRAINT, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.ExecutionConstraintSimple = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseConstitutedEntity(text)
	outErr = handleParsingError(tree.CONSTITUTED_ENTITY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.ConstitutedEntity = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseConstitutedEntityProperty(text)
	outErr = handleParsingError(tree.CONSTITUTED_ENTITY_PROPERTY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.ConstitutedEntityPropertySimple = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseModal(text)
	outErr = handleParsingError(tree.MODAL, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.Modal = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseConstitutingFunction(text)
	outErr = handleParsingError(tree.CONSTITUTIVE_FUNCTION, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.ConstitutiveFunction = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseConstitutingProperties(text)
	outErr = handleParsingError(tree.CONSTITUTING_PROPERTIES, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.ConstitutingProperties = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	result, componentText, err = parseConstitutingPropertiesProperty(text)
	outErr = handleParsingError(tree.CONSTITUTING_PROPERTIES_PROPERTY, err)
	if outErr.ErrorCode != tree.PARSING_NO_ERROR {
		// Populate return structure
		ret := []tree.Node{tree.Node{Entry: &s}}
		return ret, strings.Join(componentText, " "), outErr
	}
	s.ConstitutingPropertiesPropertySimple = result
	// Remove elements parsed as part of the component parsing
	for _, v := range componentText {
		remainingString = strings.ReplaceAll(remainingString, v, "")
	}

	Println("Basic statement: " + s.String())
	if !s.IsEmpty() {
		// Returns generated node, alongside remaining input string content for further parsing
		return []tree.Node{tree.Node{Entry: s}}, remainingString, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
	}
	// else return empty array
	return []tree.Node{}, remainingString, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

func parseAttributes(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.ATTRIBUTES, false, text)
}

func parseAttributesProperty(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.ATTRIBUTES_PROPERTY, true, text)
}

func parseDeontic(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.DEONTIC, false, text)
}

func parseAim(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.AIM, false, text)
}

func parseDirectObject(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.DIRECT_OBJECT, false, text)
}

func parseDirectObjectProperty(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.DIRECT_OBJECT_PROPERTY, true, text)
}

func parseIndirectObject(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.INDIRECT_OBJECT, false, text)
}

func parseIndirectObjectProperty(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.INDIRECT_OBJECT_PROPERTY, true, text)
}

func parseConstitutedEntity(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTED_ENTITY, false, text)
}

func parseConstitutedEntityProperty(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTED_ENTITY_PROPERTY, true, text)
}

func parseModal(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.MODAL, false, text)
}

func parseConstitutingFunction(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTIVE_FUNCTION, false, text)
}

func parseConstitutingProperties(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTING_PROPERTIES, false, text)
}

func parseConstitutingPropertiesProperty(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.CONSTITUTING_PROPERTIES_PROPERTY, true, text)
}

func parseActivationCondition(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.ACTIVATION_CONDITION, false, text)
}

func parseExecutionConstraint(text string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponentWithParentheses(tree.EXECUTION_CONSTRAINT, false, text)
}

/*
Extracts a component content from string based on component signature (e.g., A, I, etc.)
and balanced parentheses/braces. Tolerates presence of suffices and annotations and includes those
in output (e.g., A1[type=animate](content)).
Allows for indication whether parsed component is actually a property.
If no component content is found, an empty string is returned.
Tests against mistaken parsing of property variant of a component (e.g., A,p() instead of A()).
*/
func extractComponentContent(component string, propertyComponent bool, input string, leftPar string, rightPar string) ([]string, tree.ParsingError) {

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

	// General component syntax (inclusive of ,p)
	r, err := regexp.Compile(component + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)
	if err != nil {
		Println("Error in regex compilation: ", err.Error())
		return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Error in Regular Expression compilation. Error: " + err.Error()}
	}
	// Component syntax to test for suffix-embedded property syntax (e.g., A1,p)
	rProp, err := regexp.Compile(component + COMPONENT_SUFFIX_SYNTAX + tree.PROPERTY_SYNTAX_SUFFIX + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)
	if err != nil {
		Println("Error in regex compilation: ", err.Error())
		return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNEXPECTED_ERROR, ErrorMessage: "Error in Regular Expression compilation. Error: " + err.Error()}
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
			}

		}
	}

	for { // infinite loop - needs to break out

		// Return index of found element
		result := r.FindAllStringIndex(processedString, 1)
		// Return content of found element
		resultContent := r.FindString(processedString)

		if nestedStatement && len(resultContent) > 0 {
			component = resultContent[:len(resultContent)-1]
			Println("Identified nested component", component)
		}

		if len(result) > 0 {
			// Start search after potential suffix and annotation elements
			startPos = result[0][0] + len(resultContent) - len(leftPar)
			Println("Component: ", resultContent)
			Println("Search start position: ", startPos)
		} else {
			// Returns component strings once opening parenthesis symbol is no longer found
			return componentStrings, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
		}

		// Parentheses count to check for balance
		parCount := 0

		// Switch to stop parsing
		stop := false

		for i, letter := range processedString[startPos:] {

			//Println("Letter to iterate over in parentheses/braces search: " + string(letter))

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
					Println("Identified component content: " + contentString)
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
		if !stop {
			// Could not find terminating parenthesis/brace if stop is not set but input string exhausted
			// Common issue: they may have passed parentheses count as part of initial validation, but may not be in correct order.
			// Example: Bdir,p(left [AND] right)) Bdir((left [AND] right)
			return nil, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_UNABLE_TO_EXTRACT_COMPONENT_CONTENT,
				ErrorMessage: "Could not determine component content of component '" + resultContent + "'. " +
					"Please review parentheses/braces in input '" + processedString + "'."}
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
- This function will not prevent wrongful extraction of property components instead of first-order components. This is handled in #extractComponentContent.
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
		// Does not guard against mistaken choice of property variants of components (e.g., A,p instead of A) - is handled in #extractComponentContent.
		reconstructedComponent := strings.Replace(strippedInput, suffix, "", 1)
		Println("Reconstructed statement:", reconstructedComponent)
		// Return only suffix
		return suffix, "", reconstructedComponent, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
	}
}

/*
Extracts statement-level annotation(s), i.e., not to be called on component, but extracting statement-level annotations.
Statement-level annotation example: "A(actor) I(aim) [annotation]". Operates only on a given nesting level.
This extraction does not support nested brackets (e.g., [ldjlkdsjg[dskljgsl]jlksg]) or nested parentheses (e.g., [ldfjs(ljdfs)sdflkj]).
In the case of brackets, only the inner bracket is extracted; in the case of parentheses, the entire annotation is ignored.
Calls log.Fatal if regex compilation fails (should never happen, since regex is static, and to be caught by runtime environment).
Returns annotations as a string array, with values including original brackets (e.g., "[annotation1][annotation2]").
The second return value is the remaining string, which is the input string with the annotations removed.
*/
func parseStatementLevelAnnotations(input string) ([]string, string) {
	// Component annotation pattern (note: does not support embedded brackets or parentheses)
	r, err := regexp.Compile("\\[" + COMPONENT_ANNOTATION_MAIN + "\\]")
	if err != nil {
		log.Fatal("Error when parsing statement-level annotations:", err.Error())
	}
	// Search for annotation patterns on input (pure annotation on statement level)
	result := r.FindAllStringSubmatch(input, -1)

	// Copy of input as remainingOutput, which will contain the remaining elements following the extraction of the annotations
	remainingOutput := input
	// Array containing individual annotations
	annotationArray := []string{}
	// Check all annotations for logical operators and strip from input
	for _, element := range result {
		for _, element2 := range element {

			// Filter for potential presence of logical operators in input
			skip := false
			for _, logOp := range tree.IGLogicalOperators {
				if "["+logOp+"]" == element2 {
					skip = true
					break
				}
			}
			if !skip {
				// Remove detected annotation from remainingOutput
				remainingOutput = strings.ReplaceAll(remainingOutput, element2, "")
				// Add annotations to output array
				annotationArray = append(annotationArray, element2)
			}
		}
	}
	Println("Number of statement-level annotations: ", len(annotationArray))
	Println("Statement annotation: ", annotationArray)
	Println("Remaining string: ", remainingOutput)
	if len(annotationArray) == 0 {
		return []string{}, remainingOutput
	}
	return annotationArray, remainingOutput
}

/*
Parses component based on surrounding parentheses.
Returns parsed node, as well as substring of input text identified as component content (including annotation and suffix).
*/
func parseComponentWithParentheses(component string, propertyComponent bool, input string) (*tree.Node, []string, tree.ParsingError) {
	return parseComponent(component, propertyComponent, input, LEFT_PARENTHESIS, RIGHT_PARENTHESIS)
}

/*
Generic entry point to parse individual components of a given statement.
Input is component symbol of interest, full statements, as well as delimiting parentheses signaling parsing for atomic
or nested components. Additionally, the parameter propertyComponent indicates whether the parsed component is a property
Returns the parsed node as well as a string array containing all relevant substrings of the input text that have been parsed.
*/
func parseComponent(component string, propertyComponent bool, text string, leftPar string, rightPar string) (*tree.Node, []string, tree.ParsingError) {

	Println("Parsing:", component)

	// TODO: For property variants, identify root property and search as to whether embedded midfix exists

	// Extract component (one or multiple occurrences) from input string based on provided component identifier
	componentStrings, err := extractComponentContent(component, propertyComponent, text, leftPar, rightPar)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, componentStrings, err
	}

	Println("Components (Count:", len(componentStrings), "):", fmt.Sprint(componentStrings))

	// Check for duplicates in individual atomic components
	duplicateChk := shared.DuplicateElement(componentStrings)
	if duplicateChk != "" {
		// If duplicate is found, return contextualized error
		return nil, componentStrings, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_DUPLICATE_COMPONENT_ENTRIES,
			ErrorMessage: "Duplicate component entry '" + duplicateChk + "' found in input statement. " +
				"Check for duplicate components and statements.",
			ErrorIgnoredElements: []string{duplicateChk}}
	}

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
			return nil, componentStrings, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_PATTERN_EXTRACTION,
				ErrorMessage: "Error during pattern extraction in combination expression."}
		}

		for i, v := range componentStrings {
			Println("Round: " + strconv.Itoa(i) + ": " + v)

			// Extracts suffix and/or annotation for individual component instance -- must only be used with single component instance!
			componentSuffix, componentAnnotation, componentContent, err := extractSuffixAndAnnotations(component, propertyComponent, v, leftPar, rightPar)
			if err.ErrorCode != tree.PARSING_NO_ERROR {
				return nil, componentStrings, err
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
				if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_NO_COMBINATIONS {
					log.Println("Error when parsing synthetically linked element. Error:", err)
					return nil, componentStrings, err
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
				if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_NO_COMBINATIONS {
					log.Println("Error when parsing synthetically linked element. Error:", err)
					return nil, componentStrings, err
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
						return nil, componentStrings, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION,
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
			return nil, componentStrings, err
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

		node1, _, err2 := ParseIntoNodeTree(componentString, false, leftPar, rightPar)
		if err2.ErrorCode == tree.PARSING_ERROR_LOGICAL_OPERATOR_OUTSIDE_COMBINATION {
			// Means that there is logical operator, but probably missing outer parentheses.
			// Try again by augmenting with parentheses, before return error if it still fails.
			node1, _, err2 = ParseIntoNodeTree(leftPar+componentString+rightPar, false, leftPar, rightPar)
		}
		if err2.ErrorCode != tree.PARSING_NO_ERROR && err2.ErrorCode != tree.PARSING_ERROR_NO_COMBINATIONS {
			log.Println("Error when parsing synthetically linked element. Error:", err2)
			return nil, componentStrings, err2
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
		return nil, componentStrings, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_COMPONENT_NOT_FOUND,
			ErrorMessage: "Component " + component + " was not found in input string"}
	}

	Println("Component Identifier: " + component)
	Println("Full string: " + componentString)

	// Some error check and override
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_NO_COMBINATIONS {
		err.ErrorMessage = "Error when parsing component " + component + ": " + err.ErrorMessage
		log.Println("Error during component parsing:", err.Error())
	}

	// Override missing combination error, since it is not relevant at this level
	if err.ErrorCode == tree.PARSING_ERROR_NO_COMBINATIONS {
		err.ErrorCode = tree.PARSING_NO_ERROR
		err.ErrorMessage = ""
	}

	return node, componentStrings, err
}
