package parser

import (
	"IG-Parser/tree"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func ParseStatement(text string) tree.Statement {
	s := tree.Statement{}

	result, err := parseAttributes(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND &&
		err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
		log.Fatal(err.Error())
	}
	s.Attributes = result

	result, err = parseDeontic(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND &&
		err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
		log.Fatal(err.Error())
	}
	s.Deontic = result

	result, err = parseAim(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND &&
		err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
		log.Fatal(err.Error())
	}
	s.Aim = result

	result, err = parseDirectObject(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND { //&&
		//err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS

		log.Fatal(err.Error())
	}
	s.DirectObject = result

	result, err = parseDirectObjectProperty(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND { //&&
		//err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS

		log.Fatal(err.Error())
	}
	s.DirectObjectProperty = result

	result, err = parseIndirectObject(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND { //&&
		//err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS

		log.Fatal(err.Error())
	}
	s.IndirectObject = result

	result, err = parseIndirectObjectProperty(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND { //&&
		//err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS

		log.Fatal(err.Error())
	}
	s.IndirectObjectProperty = result

	result, err = parseActivationCondition(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND &&
		err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
		log.Fatal(err.Error())
	}
	s.ActivationConditionSimple = result

	result, err = parseExecutionConstraint(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND &&
		err.ErrorCode != tree.PARSING_ERROR_IGNORED_ELEMENTS {
		log.Fatal(err.Error())
	}
	s.ExecutionConstraintSimple = result

	// Switch on number of component patterns (not combinations)
	/*switch len(bdir) {
	case 1:		s.DirectObject = tree.ComponentLeafNode(bdir[0][0], tree.DIRECT_OBJECT)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(bdir)) + " items.")
	default: 	log.Println("No Direct Object found")
	}*/

	//fmt.Println(s.String())

	//os.Exit(0)

	return s

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
Extracts a component specification from string based on component signature (e.g., A, I, etc.)
and balanced parentheses.
If no component is found, an empty string is returned
*/
func extractComponent(component string, input string) []string {

	// Strings for given component
	componentStrings := []string{}

	// Copy string for truncating
	processedString := input

	fmt.Println("Looking for component: " + component)

	for { // infinite loop - needs to break out
		// Find first occurrence of signature
		startPos := strings.Index(processedString, component + "(")

		if startPos == -1 {
			//log.Println("Component signature " + component + " not found in input string '" + input + "'")
			return componentStrings
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
	componentStrings := extractComponent(component, text)

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

	node, modifiedInput, err := ParseDepth(componentString, false)

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
