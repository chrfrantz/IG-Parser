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

func ParseStatement(text string) tree.Statement {
	s := tree.Statement{}

	result, err := parseAttributes(text)
	handleParsingError(tree.ATTRIBUTES, err)
	s.Attributes = result

	result, err = parseAttributesProperty(text)
	handleParsingError(tree.ATTRIBUTES_PROPERTY, err)
	s.AttributesProperty = result

	result, err = parseDeontic(text)
	handleParsingError(tree.DEONTIC, err)
	s.Deontic = result

	result, err = parseAim(text)
	handleParsingError(tree.AIM, err)
	s.Aim = result

	result, err = parseDirectObject(text)
	handleParsingError(tree.DIRECT_OBJECT, err)
	s.DirectObject = result

	result, err = parseDirectObjectProperty(text)
	handleParsingError(tree.DIRECT_OBJECT_PROPERTY, err)
	s.DirectObjectProperty = result

	result, err = parseIndirectObject(text)
	handleParsingError(tree.INDIRECT_OBJECT, err)
	s.IndirectObject = result

	result, err = parseIndirectObjectProperty(text)
	handleParsingError(tree.INDIRECT_OBJECT_PROPERTY, err)
	s.IndirectObjectProperty = result

	result, err = parseActivationCondition(text)
	handleParsingError(tree.ACTIVATION_CONDITION, err)
	s.ActivationConditionSimple = result

	result, err = parseExecutionConstraint(text)
	handleParsingError(tree.EXECUTION_CONSTRAINT, err)
	s.ExecutionConstraintSimple = result

	result, err = parseConstitutedEntity(text)
	handleParsingError(tree.CONSTITUTED_ENTITY, err)
	s.ConstitutedEntity = result

	result, err = parseConstitutedEntityProperty(text)
	handleParsingError(tree.CONSTITUTED_ENTITY_PROPERTY, err)
	s.ConstitutedEntityProperty = result

	result, err = parseModal(text)
	handleParsingError(tree.MODAL, err)
	s.Modal = result

	result, err = parseConstitutingFunction(text)
	handleParsingError(tree.CONSTITUTIVE_FUNCTION, err)
	s.ConstitutiveFunction = result

	result, err = parseConstitutingProperties(text)
	handleParsingError(tree.CONSTITUTING_PROPERTIES, err)
	s.ConstitutingProperties = result

	result, err = parseConstitutingPropertiesProperty(text)
	handleParsingError(tree.CONSTITUTING_PROPERTIES_PROPERTY, err)
	s.ConstitutingPropertiesProperty = result

	return s

}

/*
Handles parsing error centrally - easier to refine.
 */
func handleParsingError(component string, err tree.ParsingError) {

	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND {
		log.Fatal("Error when parsing: ", err)
	}

}

func parseAttributes(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.ATTRIBUTES, text)
}

func parseAttributesProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.ATTRIBUTES_PROPERTY, text)
}

func parseDeontic(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.DEONTIC, text)
}

func parseAim(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.AIM, text)
}

func parseDirectObject(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.DIRECT_OBJECT, text)
}

func parseDirectObjectProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.DIRECT_OBJECT_PROPERTY, text)
}

func parseIndirectObject(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.INDIRECT_OBJECT, text)
}

func parseIndirectObjectProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.INDIRECT_OBJECT_PROPERTY, text)
}

func parseConstitutedEntity(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.CONSTITUTED_ENTITY, text)
}

func parseConstitutedEntityProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.CONSTITUTED_ENTITY_PROPERTY, text)
}

func parseModal(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.MODAL, text)
}

func parseConstitutingFunction(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.CONSTITUTIVE_FUNCTION, text)
}

func parseConstitutingProperties(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.CONSTITUTING_PROPERTIES, text)
}

func parseConstitutingPropertiesProperty(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.CONSTITUTING_PROPERTIES_PROPERTY, text)
}

func parseActivationCondition(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.ACTIVATION_CONDITION, text)
}

func parseExecutionConstraint(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.EXECUTION_CONSTRAINT, text)
}

/*
Validates input with respect to parentheses balance.
 */
func validateInput(text string) (tree.ParsingError) {
	// Validate parentheses in input
	parCount := 0
	for i, letter := range text {

		switch string(letter) {
		case "(":
			parCount++
		case ")":
			parCount--
		}
		i++
	}
	if parCount != 0 {
		msg := "Please review the parentheses in the input statement. "
		par := ""
		parCountAbs := math.Abs(float64(parCount))
		if parCount == 1 || parCount == -1 {
			msg += "There is "
			par = "parenthesis"
		} else {
			msg += "There are "
			par = "parentheses"
		}
		if parCount > 0 {
			// too many left parentheses
			msg = fmt.Sprint(msg, parCountAbs, " additional opening ", par, " ('(').")
		} else {
			// too many right parentheses
			msg = fmt.Sprint(msg, parCountAbs, " additional closing ", par, " (')').")
		}
		log.Println(msg)
		return tree.ParsingError{ErrorCode: tree.PARSING_ERROR_IMBALANCED_PARENTHESES, ErrorMessage: msg}
	}

	return tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}

/*
Extracts a component specification from string based on component signature (e.g., A, I, etc.)
and balanced parentheses.
If no component is found, an empty string is returned
*/
func extractComponent(component string, input string) ([]string, tree.ParsingError) {

	// Strings for given component
	componentStrings := []string{}

	// Copy string for truncating
	processedString := input

	fmt.Println("Looking for component: " + component)

	// Validate input string first
	err := validateInput(processedString)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	for { // infinite loop - needs to break out
		// Find first occurrence of signature
		startPos := strings.Index(processedString, component + "(")

		if startPos == -1 {
			//log.Println("Component signature " + component + " not found in input string '" + input + "'")
			return componentStrings, tree.ParsingError{ErrorCode: tree.TREE_NO_ERROR}
		}

		// Parentheses count to check for balance
		parCount := 0

		stop := false

		for i, letter := range processedString[startPos:] {

			switch string(letter) {
			case "(":
				parCount++
			case ")":
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

// Logical operators prepared for regular expression
var logicalOperators = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"
// Word pattern for regular expressions (including parentheses, spaces, square brackets, etc.)
var wordsWithParentheses = "([a-zA-Z',;()\\[\\]]+\\s*)+"
// Pattern of combinations, e.g., ( ... [AND] ... )
var combinationPattern = "\\(" + wordsWithParentheses + "(\\[" + logicalOperators + "\\]\\s" + wordsWithParentheses + ")+\\)"

func parseComponent(component string, text string) (*tree.Node, tree.ParsingError) {

	// Extract component (one or multiple occurrences) from input string based on provided component identifier
	componentStrings, err := extractComponent(component, text)
	if err.ErrorCode != tree.TREE_NO_ERROR {
		return nil, err
	}

	fmt.Println("Components: " + fmt.Sprint(componentStrings))

	// Initialize output string for parsing
	componentString := ""

	// [AND]-link different components (if multiple occur in input string)
	if len(componentStrings) > 1 {
		r, _ := regexp.Compile(combinationPattern)
		// Add leading parenthesis
		componentString = "("
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
				componentString += ")"
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

	node, modifiedInput, err := ParseIntoNodeTree(componentString, false)

	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
		err.ErrorMessage = "Error when parsing component " + component + ": " + err.ErrorMessage
		log.Fatal("Error during component parsing: ", err.Error())
	}

	// Override missing combination error, since it is not relevant at this level
	if err.ErrorCode == tree.PARSING_NO_COMBINATIONS {
		err.ErrorCode = tree.PARSING_NO_ERROR
		err.ErrorMessage = ""
	}

	fmt.Println("Modified output for " + component + ": " + modifiedInput)

	return node, err
}
