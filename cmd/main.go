package main

import (
	"IG-Parser/parser"
	"IG-Parser/tree"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var words = "([a-zA-Z',;]+\\s*)+"
var wordsWithParentheses = "([a-zA-Z',;()]+\\s*)+"
var logicalOperators = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"

func main() {

	//text := "National Organic Program's Program Manager, on behalf of the Secretary, may (inspect and [AND] review) (certified production and [AND] handling operations and [AND] accredited certifying agents) for compliance with the (Act or [XOR] regulations in this part)."

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(review [OR] refresh) " +
		"Bdir(certified production and [AND] handling operations and [AND] accredited certifying agents) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part) )."


	s := parseStatement(text)

	//s.String()
	fmt.Println(s.String())

}

func parseStatement(text string) tree.Statement {
	s := tree.Statement{}


	result, err := parseAttributes(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}
	s.Attributes = result

	result, err = parseDeontic(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}
	s.Deontic = result



	result, err = parseAim(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}
	s.Aim = result


	result, err = parseDirectObject(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}
	s.DirectObject = result


	result, err = parseExecutionConstraint(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}
	s.ExecutionConstraintSimple = result

	result,err = parseActivationCondition(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_ERROR_COMPONENT_NOT_FOUND {
		log.Fatal(err.Error())
	}
	s.ActivationConditionSimple = result




	// Switch on number of component patterns (not combinations)
	/*switch len(bdir) {
	case 1:		s.DirectObject = tree.ComponentLeafNode(bdir[0][0], tree.DIRECT_OBJECT)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(bdir)) + " items.")
	default: 	log.Println("No Direct Object found")
	}*/

	//fmt.Println(s.String())


	//os.Exit(0)

	return s

	/*
	bdir := parseDirectObject(text)
	switch len(bdir) {
	case 1:		s.DirectObject = igTree.ComponentLeafNode(bdir[0][0], igTree.DIRECT_OBJECT)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(bdir)) + " items.")
	default: 	log.Println("No Direct Object found")
	}

	bdirp := parseDirectObjectProperty(text)
	switch len(bdirp) {
	case 1:		s.DirectObjectProperty = igTree.ComponentLeafNode(bdirp[0][0], igTree.DIRECT_OBJECT_PROPERTY)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(bdirp)) + " items.")
	default: 	log.Println("No Direct Object Property found")
	}

	bind := parseIndirectObject(text)
	switch len(bind) {
	case 1:		s.IndirectObject = igTree.ComponentLeafNode(bind[0][0], igTree.INDIRECT_OBJECT)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(bind)) + " items.")
	default: 	log.Println("No Indirect Object found")
	}

	bindp := parseIndirectObjectProperty(text)
	switch len(bindp) {
	case 1:		s.IndirectObjectProperty = igTree.ComponentLeafNode(bindp[0][0], igTree.INDIRECT_OBJECT_PROPERTY)
	case 2: 	log.Fatal("Encountered " + strconv.Itoa(len(bindp)) + " items.")
	default: 	log.Println("No Indirect Object Property found")
	}

	return s
*/
	/*
	e := parseConstitutedEntity(text)[0]
	s.ConstitutedEntity = igTree.ComponentLeafNode(e, igTree.CONSTITUTED_ENTITY)

	ep := parseConstitutedEntityProperty(text)[0]
	s.ConstitutedEntityProperty = igTree.ComponentLeafNode(ep, igTree.CONSTITUTED_ENTITY_PROPERTY)

	f := parseConstitutingFunction(text)[0]
	s.ConstitutiveFunction = igTree.ComponentLeafNode(f, igTree.CONSTITUTIVE_FUNCTION)

	p := parseConstitutingProperties(text)[0]
	s.ConstitutingProperties = igTree.ComponentLeafNode(p, igTree.CONSTITUTING_PROPERTIES)

	pp := parseConstitutingPropertiesProperty(text)[0]
	s.ConstitutingPropertiesProperty = igTree.ComponentLeafNode(pp, igTree.CONSTITUTING_PROPERTIES_PROPERTY)

	/* TODO: Review for complex statements */
	/*cac := parseActivationCondition(text)[0]
	s.ActivationConditionSimple = igTree.ComponentLeafNode(cac, igTree.ACTIVATION_CONDITION)

	cex := parseExecutionConstraint(text)[0]
	s.ExecutionConstraintSimple = igTree.ComponentLeafNode(cex, igTree.EXECUTION_CONSTRAINT)

	return s
	*/
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

/*
func parseIndirectObjectProperty(text string) [][]string {
	return parseComponentOld(tree.INDIRECT_OBJECT_PROPERTY, text)
}

func parseConstitutedEntity(text string) [][]string {
	return parseComponentOld(tree.CONSTITUTED_ENTITY, text)
}

func parseConstitutedEntityProperty(text string) [][]string {
	return parseComponentOld(tree.CONSTITUTED_ENTITY_PROPERTY, text)
}

func parseConstitutingFunction(text string) [][]string {
	return parseComponentOld(tree.CONSTITUTIVE_FUNCTION, text)
}

func parseConstitutingProperties(text string) [][]string {
	return parseComponentOld(tree.CONSTITUTING_PROPERTIES, text)
}

func parseConstitutingPropertiesProperty(text string) [][]string {
	return parseComponentOld(tree.CONSTITUTING_PROPERTIES_PROPERTY, text)
}
*/

func parseActivationCondition(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.ACTIVATION_CONDITION, text)
}

func parseExecutionConstraint(text string) (*tree.Node, tree.ParsingError) {
	return parseComponent(tree.EXECUTION_CONSTRAINT, text)
}

func parseComponentOld(component string, text string) [][]string {
	log.Println("Invoking parsing of component " + component)

	r, _ := regexp.Compile(component + "\\(" + words + "(\\[" + logicalOperators + "\\]\\s" + words + ")*\\)")
	return r.FindAllStringSubmatch(text, -1)
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

func parseComponent(component string, text string) (*tree.Node, tree.ParsingError) {
	// Detects any content framed by component prefix
	//r, _ := regexp.Compile(component + "\\(" + words + "(\\[" + logicalOperators + "\\]\\s" + words + ")*\\)")

	/*for k,v := range r.FindAllStringSubmatch(text, -1){
		fmt.Println(k)
		fmt.Println(v[0])
	}*/

	//fmt.Println()

	/*matches := r.FindAllStringSubmatch(text, -1)
	if len(matches) == 0 {
		log.Fatal("No element found for component " + component + " found in text " + text)
	}
	if len(matches) > 1 {
		log.Fatal("More than one component " + component + " found in text " + text)
	}*/

	componentStrings := extractComponent(component, text)



	componentString := "("

	// Create root node
	node := tree.Node{}


	if len(componentStrings) > 1 {
		for i, v := range componentStrings {
			fmt.Println(strconv.Itoa(i) + ": " + v)
			componentString += v[len(component):]
			if i < len(componentStrings)-1 {
				componentString += " " + tree.AND_BRACKETS + " "
			} else {
				componentString += ")"
			}
		}
	} else if len(componentStrings) == 1 {
		// Single entry
		componentString = componentStrings[0]
	} else {
		return nil, tree.ParsingError{tree.PARSING_ERROR_COMPONENT_NOT_FOUND,
			"Component " + component + " was not found in input string"}
	}

	fmt.Println(componentString)

	//os.Exit(2)
	// Parse provided expression
	_, modifiedInput, err := parser.ParseDepth(componentString, &node)

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		err.ErrorMessage = "Error when parsing component " + component + ": " + err.ErrorMessage
	}

	// Override missing combination error, since it is not relevant at this level
	if err.ErrorCode == tree.PARSING_NO_COMBINATIONS {
		err.ErrorCode = tree.PARSING_NO_ERROR
		err.ErrorMessage = ""
	}

	fmt.Println("Modified output for " + component + ": " + modifiedInput)

	return &node, err
}

func parseCombinations(text string) [][]string {
	log.Println("Invoking parsing of generic combinations")

	//fmt.Println(text)

	text = text + " [AND] pray"

	text = "shared (inspect and [AND] ((review [XOR] muse) [AND] pray)) last"

	fmt.Println("Raw text: " + text)

	fmt.Println("Outer most:")
	// Detect most outer combinations
	r, _ := regexp.Compile("\\(" + wordsWithParentheses + "(\\[" + logicalOperators + "\\]\\s+" + wordsWithParentheses + ")+\\)")

	tp := ""
	for k,v := range r.FindAllStringSubmatch(text, -1){
		fmt.Println(k)
		fmt.Println(v[0])
		tp = v[0]
	}

	fmt.Println("Identify delta:")
	// Calculate shared elements (delta to input text; identify separation)
	fmt.Println("Prefix: " + strconv.Itoa(strings.Index(text, tp)) + ": " + text[0:strings.Index(text, tp)])
	fmt.Println("Suffix: " + strconv.Itoa(strings.LastIndex(text, tp)))

	fullTokens := strings.Split(text, " ")
	partTokens := strings.Split(tp, " ")

	//fmt.Println("Full: " + fullTokens)
	//fmt.Println("partTokens)
	delta := []string{}
	for i := range fullTokens {
		if !tree.StringInSlice(fullTokens[i], partTokens) {
			delta = append(delta, fullTokens[i])
		}
	}

	fmt.Println("Delta: " + strings.Join(delta, " "))



	fmt.Println("Next lower focus")

	r, _ = regexp.Compile("\\(" + wordsWithParentheses + "(\\[" + logicalOperators + "\\]\\s+" + wordsWithParentheses + ")+\\)")

	modifiedText := text

	fmt.Println("Starting point: " + modifiedText)
	fmt.Println(r.FindAllStringSubmatch(modifiedText, -1))

	// Outer level
	fullCombo :=  r.FindAllStringSubmatch(modifiedText, -1)
	if len(fullCombo) == 0 {
		//no combination found, i.e., break out
		log.Fatal("No combination found; properly break out")
	}



	// find left and right element
	//r, _ = regexp.Compile("\\(" + words + logicalOperators)

	//leftr.FindAllStringSubmatch(fullCombo)
	//fullCombo

	// Match elements

	input := fullCombo[0][0]

	//input = "(inspect [AND] review)"

	input = "((inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray)))"

	node,_ := ParseExpression(0, input, "left")


	fmt.Println("Close to exit: " + node.String())

	os.Exit(0)

	fmt.Println(fullCombo)

	ct := 0
	for len(r.FindAllStringSubmatch(modifiedText, -1)) > 0 {

		fmt.Println("Text before: " + modifiedText)

		modifiedText = modifiedText[1 : len(modifiedText)-1]

		fmt.Println("Text after: " + modifiedText)

		for k, v := range r.FindAllStringSubmatch(modifiedText, -1) {
			fmt.Println(k)
			fmt.Println(v[0]) // Check for multiples
			modifiedText = v[0] // Hack
		}
		fmt.Println("Round: " + strconv.Itoa(ct))
		ct += 1
	}














	//r, _ = regexp.Compile(words)

	/*for k,v := range r.FindAllStringSubmatch(text, -1){
		fmt.Println(k)
		fmt.Println(v[0])
	}*/

	fmt.Println("Inner most:")
	// Detects any inner-most generic combination in text of arbitrary number of logical operators
	r, _ = regexp.Compile("\\(" + words + "(\\[" + logicalOperators + "\\]\\s+" + words + ")+\\)")


	for k,v := range r.FindAllStringSubmatch(text, -1){
		fmt.Println(k)
		fmt.Println(v[0])
	}

	return r.FindAllStringSubmatch(text, -1)

}

func ParseExpression(level int, expression string, searchMode string) (*tree.Node, string) {

	fmt.Println("Received call on expression " + expression + " with search mode " + searchMode)

	node := tree.Node{}

	levelIndex := 0 // set index to track parsing within function
	finished := false // flag indicating whether parsing on level is finished
	for i, letter := range expression {

		if finished {
			fmt.Println("Level " + strconv.Itoa(level) + ": Constructed node: " + node.String())
			return &node, expression[levelIndex:]
		}

		//fmt.Println("Letter: " + string(letter))
		switch string(letter) {
		case "(":
			level++
			fmt.Println("Level deeper: " + strconv.Itoa(level))
			// Captures the index on which level has been entered
			levelIndex = i
			// If invoked by other call, iterate further down and search for left expression (since leading "(")
			if searchMode == "right" {
				fmt.Println("Invoking nested statement ... ")
				node.Left, expression = ParseExpression(level, expression[levelIndex+1:], "left")
				node.Right, expression = ParseExpression(level, expression, "right")
				return &node, expression
			}
		case ")":
			if searchMode == "right" {
				fmt.Println("Found end of expression")
				node.Entry = expression[levelIndex:i]
				fmt.Println("Level " + strconv.Itoa(level) + ": Constructed node: " + node.String())
				levelIndex = i+1
				//fmt.Println("Remaining string: " + expression[levelIndex])
				return &node, expression[levelIndex:]
			}
			if searchMode == "left" {
				return ParseExpression(level, expression[levelIndex+1:], "right")
			}
			level--
			fmt.Println("Level higher: " + strconv.Itoa(level))
			levelIndex = i
		case "[":
			//if searchMode == "left" {
				fmt.Println("Checking for logical operator ... " + expression[i:i+5])
				foundOperator := ""
				switch expression[i : i+5] {
				case "[AND]":
					fmt.Println("Detected [AND]")
					foundOperator = "AND"
				case "[OR] ":
					fmt.Println("Detected [OR]")
					foundOperator = "OR"
				case "[XOR]":
					fmt.Println("Detected [XOR]")
					foundOperator = "XOR"
				}
				if foundOperator != "" {
					/*if searchMode == "right" {
						centerNode := igTree.Node{}
						centerNode.LogicalOperator = foundOperator
						levelIndex = i + len(foundOperator) + 2 // Clip everything of left side and operator
						fmt.Println("Found operator right: " + foundOperator)
						return &centerNode, expression[levelIndex:]
					}*/
					if searchMode == "left" {
						// must be left leaf node
						fmt.Println("Found left leaf on level " + strconv.Itoa(level) + ": " + expression[levelIndex:i-1])
						leftLeaf := tree.Node{}
						leftLeaf.Entry = expression[levelIndex+1 : i-1]
						node.Left = &leftLeaf
						fmt.Println("Assigning operator " + foundOperator)
						node.LogicalOperator = foundOperator
						fmt.Println("Calling right-side parsing")
						levelIndex = i + len(foundOperator) + 2 // Clip everything of left side and operator
						node.Right, expression = ParseExpression(level, expression[levelIndex:], "right")
						fmt.Println("Constructed complete node: " + node.String())
						fmt.Println("Returning node to higher level ... (remaining expression: " + expression + ")")
						return &node, expression
					}
					//if level == 1 {

						//finished = true
					//}
				//}
			}
		default:
			fmt.Println("Omitting: " + string(letter))
		}
	}
	fmt.Println("Level " + strconv.Itoa(level) + ": Constructed node: " + node.String())
	return &node,expression[levelIndex:]
}