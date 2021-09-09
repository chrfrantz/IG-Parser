package main

import (
	"IG-Parser/app"
	"IG-Parser/exporter"
	"IG-Parser/parser"
	"IG-Parser/tree"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

//var words = "([a-zA-Z',;]+\\s*)+"
//var wordsWithParentheses = "([a-zA-Z',;()]+\\s*)+"
//var logicalOperators = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"

func main9() {

	//text := "National Organic Program's Program Manager, on behalf of the Secretary, may (inspect and [AND] review) (certified production and [AND] handling operations and [AND] accredited certifying agents) for compliance with the (Act or [XOR] regulations in this part)."

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	text = "A((certifying agent [AND] wife)) D(may) I(investigate) " +
		"Bdir((complaints of noncompliance with the (Act or [OR] regulations of this part) " +
		"concerning " +
		"(production [operation] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
		")"
		//"fdlkgjdflg))" // certified as organic by the certifying agent))."

	//text = "A(Actor) D(must) I((comply [OR] violate)) with Bdir(provisions) Cac[dfg]{A(Actor) I(has (applied and [AND] advocated)) for Bdir(certification)}"

	text = "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{A(Programme Manager) I((suspects [AND] assesses)) Bdir(violations)}" +
		"{Cac{E(Program Manager) F(is) P(qualified)} [AND] " +
		"{Cac{E(Program Participant2) F(is2) P(employed2)} [XOR] " +
		"Cac{E(Program Participant) F(is) P(employed)}}}"

	_, err := app.ConvertIGScriptToGoogleSheets(text, "65", "fun.csv")

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}

	//fmt.Println(output)

	/*
	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}

	log.Println("Step 2: Extracting leaf arrays")
	// Retrieve leaf arrays from generated tree (alongside frequency indications for components)
	leafArrays, componentRefs := s.GenerateLeafArrays()

	fmt.Println("Printing stuff: ")

	fmt.Println(s.String())

	fmt.Println("Leaf arrays: ")
	fmt.Println(leafArrays)

	fmt.Println(fmt.Sprint(componentRefs))
	*/

}

func main0()  {

	input := []int{1,2,3,4,6,7,8}

	slc := []string{}

	for k, v := range input {
		fmt.Println("Count: ", k)
		fmt.Println("Value: ", v)
		slc = exporter.GenerateReferenceSlice(slc, v, true, true)
		fmt.Println(slc)
	}


}

func main() {
	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{A(Programme Manager) I(suspects) Bdir(violations) Cac{A(NOP Manager) I(orders) Bdir(review)}}"

	text = "A1[gsdkgjsl](National Organic Program's Program Manager), A(Banana) Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I1(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		//"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		//"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		//"Cac{A(Programme Manager) I(suspects) Bdir(violations)}"
		"Cac,p[gklsjg]{A(Programme Manager) I(suspects) Bdir(violations) Cac{A(NOP Manager) I(orders) Bdir(review)}}"

	text = "A1(content) A[annotation](content2) A2(more content) A#pd$[another| annotation](even $more| content)"
	//text = "A1[dsgjlks](A&dsisgj=), A1(sdgjlks[operato]) "

	/*
	nestedStatements, err := parser.ExtractComponentContent("", text, parser.LEFT_BRACE, parser.RIGHT_BRACE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		fmt.Println(err.Error())
	}

	fmt.Println("Nested statements: " + fmt.Sprint(nestedStatements))*/

	/*one, two, three :=  parser.extractSuffixAndAnnotations("Cac", text, "{", "}")

	fmt.Println("First:", one)
	fmt.Println("Second:", two)
	fmt.Println("Third:", three)

	 */

	s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		fmt.Println(err.Error())
	}
	fmt.Println(s.String())
}

func main100() {

	component := "Bdir,p"
	input := "Bdir,p1[ctx=time](t[sgdsgs]ext) Bdir,p{glksdjglk} A(Program Manager) A1(Nurse1,sgslk;l[sdglsj]) dsdgjoids Bdir,p1(dslkgjlkdsg)"

	fmt.Println(component)

	//,p1,p[prop=[odfs,right](sdglksdj)]

	// Remove component
	//strippedInput := strings.ReplaceAll(input, component, "")

	leftPar := "("
	rightPar := ")"

	//leftPar = "{"
	//rightPar = "}"

	processedString := input

	// Component prefix (word without spaces and parentheses, but [] brackets)
	//componentPrefix := "(\\[([a-zA-Z,=;\\w{}\\[\\]\\(\\)])+\\])?" //\\("

	/*contentSyntax := "[\\" + parser.LEFT_PARENTHESIS + "\\" + parser.LEFT_BRACE + "]" +
		".*" + //"([a-zA-Z0-9][,;-\\+]\\[\\]{}]+)" +
		"[\\" + parser.RIGHT_PARENTHESIS + "\\" + parser.RIGHT_BRACE + "]"*/
	//suffixSyntax := "[a-zA-Z,0-9]*"
	componentSyntax := component//"([a-zA-Z,]+)"

	// Strings for given component
	componentStrings := []string{}

	startPos := -1

	for { // infinite loop - needs to break out
		// Find first occurrence of signature in processedString (incrementally iterated by letter)
		//startPos := strings.Index(processedString, component + leftPar)

		// Search number of entries
		r, err := regexp.Compile(componentSyntax + parser.COMPONENT_SUFFIX_SYNTAX + parser.COMPONENT_ANNOTATION_SYNTAX + "?\\" + leftPar)
		// + escapeSymbolsForRegex(input)
		if err != nil {
			log.Fatal("Error", err.Error())
		}
		result := r.FindAllStringIndex(processedString, -1)

		fmt.Println(result)
		if len(result) > 0 {
			fmt.Println(len(result[0]))
		}
		//fmt.Println(result[0][0])
		//fmt.Println(result[0][1])
		//fmt.Println(result[0][2])

		if len(result) > 0 {
			startPos = result[0][0]
			fmt.Println("Start position: ", startPos)
		} else {
			// Returns component strings once opening parenthesis symbol is no longer found
			//return componentStrings, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
			fmt.Println("Finished: ", componentStrings)
			os.Exit(0)
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
					componentStrings = append(componentStrings, processedString[startPos:startPos+i+1])
					fmt.Println("Added string " + processedString[startPos:startPos+i+1])
					processedString = processedString[startPos+i+1:]
					fmt.Println("Remaining string ", processedString)
					stop = true
				}
			}
			if stop {
				break
			}
		}
	}


	/*

	// Search number of entries
	r, err := regexp.Compile(componentSyntax + suffixSyntax + componentPrefix + contentSyntax)
	// + escapeSymbolsForRegex(input)
	if err != nil {
		log.Fatal("Error", err.Error())
	}
	result := r.FindAllStringSubmatch(strippedInput, -1)
	//
	fmt.Println(len(result[0]))
	fmt.Println(result)
	fmt.Println(result[0][0])
	fmt.Println(result[0][1])
	fmt.Println(result[0][2])
*/

	/*
	// Property
	r, err := regexp.Compile(componentPrefix)
	// + escapeSymbolsForRegex(input)
	if err != nil {
		log.Fatal("Error", err.Error())
	}
	result := r.FindAllStringSubmatch(strippedInput, -1)
	//
	fmt.Println(len(result))

	if len(result) > 0 {
		res := result[0][0]
		res = res[:len(res)-1]
		fmt.Println(res)
		pos := strings.Index(strippedInput, res)
		suffix := strippedInput[:pos]
		fmt.Println("Suffix:", suffix)
		reconstructedComponent := component + strings.ReplaceAll(strippedInput, suffix + res, "")
		fmt.Println(reconstructedComponent)
	} else {
		// Test for suffix
		contentStartPos := strings.Index(strippedInput, "(")
		suffix := strippedInput[:contentStartPos]
		fmt.Println("Suffix:", suffix)
		reconstructedComponent := component + strings.ReplaceAll(strippedInput, suffix, "")
		fmt.Println(reconstructedComponent)
	}

*/



}

/*
Extracts potential annotations in input string, with first return element being
numeric information and second one being annotations. Returns empty for either if
empty.
 */
 func ExtractAnnotations(component string, input string) (string, string) {

	componentNamePos := strings.Index(input, component)
	if componentNamePos == -1 {
		log.Println("Couldn't find component name.")
		return "", ""
	}

	// Strip component name
	input = input[len(component):]

	annotationEndPos := strings.Index(input, "(")
	if annotationEndPos == -1 {
		log.Println("Couldn't find opening parenthesis.")
		return "", ""
	}

	prefixString := input[:annotationEndPos]

	fmt.Println(prefixString)

	// Find numeric value (preceding square brackets)
	numericAnnotationEndPos := strings.Index(prefixString, "[")

	// Return content for numeric part
	numericReturn := ""
	// Return content for annotation part
	annotationReturn := ""

	if numericAnnotationEndPos == -1 {
		log.Println("No annotation bracket found")
		// Assume whole prefix string as numeric
		numericReturn = prefixString
	} else {
		// Extract leading part until first bracket ([)
		numericReturn = input[:numericAnnotationEndPos]
		// Find last end bracket (may have embedded brackets)
		annotationPartEndPos := strings.LastIndex(prefixString, "]")
		annotationReturn = input[numericAnnotationEndPos+1:annotationPartEndPos]
	}

	return numericReturn, annotationReturn
}

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
	text = strings.ReplaceAll(text, "$", "\\$")
	text = strings.ReplaceAll(text, "+", "\\+")

	return text
}

func main00() {

	var componentPrefix = "([a-zA-Z\\[\\]]+)+"
	//v := "\\{A\\(Actor\\) I\\(has applied\\) for Bdir\\(certification\\)\\}"
	v := "{A(Actor) I(has applied) for Bdir(certification)}"

	v = strings.ReplaceAll(v, "{", "\\{")
	v = strings.ReplaceAll(v, "}", "\\}")
	v = strings.ReplaceAll(v, "(", "\\(")
	v = strings.ReplaceAll(v, ")", "\\)")
	v = strings.ReplaceAll(v, "[", "\\[")
	v = strings.ReplaceAll(v, "]", "\\]")

	statement := "A(Actor) D(must) I(comply) with Bdir(provisions) Cac[dfg]{A(Actor) I(has applied) for Bdir(certification)}"

	r, err := regexp.Compile(componentPrefix + "" + v + "")

	if err != nil {
		fmt.Println(err.Error())
	}

	result := r.FindAllStringSubmatch(statement, 1)

	fmt.Println(result[0][0])
}

func mainx() {
	text := "(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	text = "A((certifying agent [AND] borrower [OR] wife)) M(may) I(investigate) " +
		"Bdir((complaints of noncompliance with the (Act or [OR] regulations of this part) " +
		"concerning " +
		"(production [operations] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
		")"+
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	fmt.Println(text)

	//os.Exit(0)

	output, err := app.ConvertIGScriptToGoogleSheets(text, "650", "output.csv")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		fmt.Println("Top: " + fmt.Sprint(err.Error()))
	}

	fmt.Println(output)
}

func main3() {
	text := "(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	text = "Y((certifying agent [AND] borrower [AND] wife)) M(may) I(investigate) " +
	"Bdir((complaints of noncompliance with the (Act or [OR] regulations of this part) " +
	"concerning " +
	"(production [operations] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
	")"+
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	//"fdlkgjdflg))" // certified as organic by the certifying agent))."

	text = "I((drink [AND] drive [AND] drown [AND] pleasure [AND] hijack))"

	s,_ := parser.ParseStatement(text)

	fmt.Println(s.String())

	//os.Exit(0)
	leafArrays, componentRefs := s.GenerateLeafArrays()

	//fmt.Println(componentRefs)

	//os.Exit(0)

	res, _ := exporter.GenerateNodeArrayPermutations(leafArrays...)

	fmt.Println("Component references: ", componentRefs)

	fmt.Println("Input arrays: ", res)

	links := exporter.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	//os.Exit(0)

	//output,_ := exporter.GenerateGoogleSheetsOutput(res, componentRefs, links, "650")

	//fmt.Println("\n" + output)

	//exporter.WriteToFile("statement.csv", output)



	/*for i, s := range res {
		fmt.Println("Statement ", i, ": ", s)
		for v := range s {
			fmt.Println("-->", s[v])
		}
	}*/

}


