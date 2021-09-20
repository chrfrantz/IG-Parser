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

func main9999() {

	leftPar := "("

	r, err := regexp.Compile("A" + parser.COMPONENT_SUFFIX_SYNTAX +
		parser.COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)
	if err != nil {
		log.Fatal("Error", err.Error())
	}

	text := "A1[annotation=(left,right)](con( )tent)"

	fmt.Println("A" + parser.COMPONENT_SUFFIX_SYNTAX +
		parser.COMPONENT_ANNOTATION_SYNTAX + "\\" + leftPar)

	fmt.Println(text)

	res := r.FindAllString(text, -1)

	fmt.Println(res)

}

func main() {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		// This is another nested component - should have implicit link to regular Bdir
		"Bdir{A(farmers) that I((apply [OR] plan to apply)) for Bdir(organic farming status)}" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		//"{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		"{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		//"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"


	sep, _ := parser.SeparateComponentsAndNestedStatements(text)

	fmt.Println(sep)


}

func main10000() {
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

	//"A1[role=origin](actor)
	text = "Cac1[role=recipient]{A2(actor2) I(helps) Bdir1(someone)}"
	text = "Cac1{A2(actor2) I(helps) Bdir1(someone)}"
	//text = "Cac((Cac1[role=recipient]{A2(actor2) I(helps) Bdir1(someone)} [AND] Cac2[gov=monitor]{A(actor3) I(aims) at Bdir(support)}))"

	text = "Cac1[leftAnno]{A1[annotation=(left,right)](content) A2[annot](content2)} Cac2[rightAnno]{A5(actor)}"

	//text = "A1[annotation1](content1) A2[annotation2](content2) A3(content3)"

	//text = "A5(content) I2(aim) Cex1[glskdgjlks](constraint)"


	//A1((attr1 [AND] attr2))
	text = "A1[glksdjgl](attr) A,p(general prop), A,p1(specific prop) I(action) Bdir,p(klgjdsklg) Bdir(dsgskjg) Bdir,p1(dslkgjslkg) Bdir1(dsgjls)"

	text = "{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	/*
	text = "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// This is the tricky lines, specifically the second Cac{}
		"Cac{E(Program Manager) F(is) P((approved [AND] committed)) Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}}"
	*/

	text = "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		// This is another nested component - should have implicit link to regular Bdir
		"Bdir{A(farmers) that I((apply [OR] plan to apply)) for Bdir(organic farming status)}" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	fmt.Println(text)

	/*
	text = "Bdir{A(farmers) that I((apply [OR] plan to apply)) for Bdir(organic farming status)}"
	text = "{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		"A{Ce(sdkgljds) I(ldkgsl)} [AND] Bdir(dslkgjsd)"

		// non-linked additional activation condition (should be linked by implicit AND)
	//text = "Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	 */

	// Remove line breaks
	text = parser.CleanInput(text)

	/*combinationPatternBraces = "\\" + LEFT_BRACE + wordsWithParentheses + "\\" + RIGHT_BRACE +
		"\\s+" + "(\\[" + logicalOperators + "\\]\\s+" + wordsWithParentheses + ")+\\" + RIGHT_BRACE

	*/

	//t1 := "E(Program Manager)"

	// WORDS

	//const SPECIAL_SYMBOLS = "',;+\\-*/%&=$£€¤§\"#!`\\|"

	WORDS_WITH_PARENTHESES := "([a-zA-Z(){}\\[\\]]+\\s*)+"
	//WORDS_WITH_PARENTHESES := "([a-zA-Z\\[\\]]+\\s*)+"

	NESTED_TERM1 := //parser.NESTED_COMPONENT_SYNTAX +
	 "(\\{" + WORDS_WITH_PARENTHESES + "\\}|\\(" + WORDS_WITH_PARENTHESES + "\\))"

	r, err := regexp.Compile(NESTED_TERM1)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	text = "(Aklsdjglksgv sklvjds) {[]jdskgl ds()}"

	res := r.FindAllString(text, -1)

	fmt.Println("Matching word (including parentheses/braces)")
	fmt.Println(res)
	fmt.Println("Count:", len(res))

	// Special characters

	const SPECIAL_SYMBOLS = "',;+\\-*/%&=$£€¤§\"#!`\\|"

	WORDS_WITH_PARENTHESES = "([a-zA-Z" + SPECIAL_SYMBOLS + "()\\[\\]]+\\s*)+"
	//WORDS_WITH_PARENTHESES := "([a-zA-Z\\[\\]]+\\s*)+"

	NESTED_TERM2 := //parser.NESTED_COMPONENT_SYNTAX +
		"(\\{" + WORDS_WITH_PARENTHESES + "\\}|\\(" + WORDS_WITH_PARENTHESES + "\\))"

	r, err = regexp.Compile(NESTED_TERM2)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	text = "(Aklsdjgl#k{sgv sk}lvjds) {[]jdskgl ds()} "

	res = r.FindAllString(text, -1)

	fmt.Println("Matching word (including parentheses/braces and special characters)")
	fmt.Println(res)
	fmt.Println("Count:", len(res))


	// Component structure

	//const SPECIAL_SYMBOLS = "',;+\\-*/%&=$£€¤§\"#!`\\|"

	WORDS_WITH_PARENTHESES = "([a-zA-Z0-9" + SPECIAL_SYMBOLS + "()\\[\\]]+\\s*)+"
	//WORDS_WITH_PARENTHESES := "([a-zA-Z\\[\\]]+\\s*)+"

	OPTIONAL_WORDS_WITH_PARENTHESES := "(" + WORDS_WITH_PARENTHESES + ")?"

	COMPONENT_SUFFIX_SYNTAX := "[a-zA-Z,0-9" + SPECIAL_SYMBOLS + "]*"

	COMPONENT_ANNOTATION_MAIN := "[a-zA-Z,0-9\\s" + SPECIAL_SYMBOLS + "]+"

	COMPONENT_ANNOTATION_OPTIONAL := "(\\[" + COMPONENT_ANNOTATION_MAIN + "\\])*"

	COMPONENT_ANNOTATION_SYNTAX := "(\\[(" + COMPONENT_ANNOTATION_MAIN + COMPONENT_ANNOTATION_OPTIONAL + ")+\\])?"
	//COMPONENT_ANNOTATION_SYNTAX := "(\\[([0-9a-zA-Z" + SPECIAL_SYMBOLS + "{}\\[\\]\\(\\)])+\\])"

	NESTED_COMPONENT_SYNTAX := "(A|D|I|Bdir|Bind|Cac|Cex|E|M|F|P)" + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX

	NESTED_TERM := //parser.NESTED_COMPONENT_SYNTAX +
		NESTED_COMPONENT_SYNTAX + "(\\{\\s*" + WORDS_WITH_PARENTHESES + "\\s*\\}|" +
			"\\(\\s*" + WORDS_WITH_PARENTHESES + "\\s*\\))"

	r, err = regexp.Compile(NESTED_TERM)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	text = "(Aklsdjgl#k{sgv sk}lvjds) {[]jdskgl ds()} Bdir,p1[ruler=governor](jglkdsjgsiovs) Cac[left=right[anotherLeft,anotherRight],right=[left,right], key=values]{A(actor) I(aim)}"

	res = r.FindAllString(text, -1)

	fmt.Println("Matching component structure (primitive and nested)")
	fmt.Println(res)
	fmt.Println("Count:", len(res))

	// Combinations

	//TODO Review for complex combinations and make reliable (multiple elements, variable presence parentheses/braces, variable use of annotations)
	NESTED_COMBINATIONS := "\\" + parser.LEFT_BRACE + "\\s*(" + NESTED_TERM + "\\s+)+" + "\\" + parser.LEFT_BRACKET +
		parser.LOGICAL_OPERATORS + "\\" + parser.RIGHT_BRACKET + "\\s+(" + NESTED_TERM + "\\s*)+" + "\\" + parser.RIGHT_BRACE

	fmt.Println(NESTED_COMBINATIONS)

	r, err = regexp.Compile(NESTED_COMBINATIONS)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	text = "(Aklsdjgl#k{sgv sk}lvjds) {[]jdskgl ds()} Bdir,p1[ruler=governor](jglkdsjgsiovs) Cac[left=right[anotherLeft,anotherRight],right=[left,right], key=values]{A(actor) I(aim)}" +
		"{A(dlkgjsg) I[dgisg](kjsdglkds) [AND] Bdir(djglksjdgkd) Cex(sdlgjlskd)}"

	res = r.FindAllString(text, -1)

	fmt.Println("Matching combinations")
	fmt.Println(res)
	fmt.Println("Count:", len(res))

	// Fixing combination refinements

	//TODO Review for complex combinations and make reliable (multiple elements, variable presence parentheses/braces, variable use of annotations)
	//NESTED_COMBINATIONS = "\\" + parser.LEFT_BRACE + "\\s*(" + NESTED_TERM + "\\s+)+" +
	//	"(" + "\\" + parser.LEFT_BRACKET + parser.LOGICAL_OPERATORS + "\\" + parser.RIGHT_BRACKET +
	//	"\\s+(" + NESTED_TERM + "\\s*)+" + ")+" + "\\" + parser.RIGHT_BRACE

	NESTED_COMBINATION :=

		// Start of alternatives
		"(" +
		// combination with parentheses
		"\\" + parser.LEFT_PARENTHESIS +
		OPTIONAL_WORDS_WITH_PARENTHESES +
		"(" + NESTED_TERM + OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
		//"\\s*(" + NESTED_TERM + "\\s*)+" +
		"(" +
		"\\" + parser.LEFT_BRACKET + parser.LOGICAL_OPERATORS + "\\" + parser.RIGHT_BRACKET +
		OPTIONAL_WORDS_WITH_PARENTHESES +
		"(" + NESTED_TERM + OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
		//"\\s+(" + NESTED_TERM + "\\s*)+" +
		")*" +
		"\\" + parser.RIGHT_PARENTHESIS +
		// OR
		"|" +
		// combinations without parentheses
		OPTIONAL_WORDS_WITH_PARENTHESES + "(" + NESTED_TERM + OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
		//"\\s*(" + NESTED_TERM + "\\s*)+" +
		"(" +
		"\\" + parser.LEFT_BRACKET + parser.LOGICAL_OPERATORS + "\\" + parser.RIGHT_BRACKET +
		OPTIONAL_WORDS_WITH_PARENTHESES + "(" + NESTED_TERM + OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
		//"\\s+(" + NESTED_TERM + "\\s*)+" +
		")*" +
		// END OF COMBINATION
		")"

	NESTED_COMBINATIONS = "\\" + parser.LEFT_BRACE +
		"\\s*(" + NESTED_COMBINATION + "\\s+)+" +
		"(" +
		"\\" + parser.LEFT_BRACKET + parser.LOGICAL_OPERATORS + "\\" + parser.RIGHT_BRACKET +
		"\\s+(" + NESTED_COMBINATION + "\\s*)+" +
		")+" +
		"\\" + parser.RIGHT_BRACE

	fmt.Println(NESTED_COMBINATIONS)

	r, err = regexp.Compile(NESTED_COMBINATIONS)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	text = "(Aklsdjgl#k{sgv sk}lvjds) {[]jdskgl ds()} Bdir,p1[ruler=governor](jglkdsjgsiovs) Cac[left=right[anotherLeft,anotherRight],right=[left,right], key=values]{A(actor) I(aim)}" +
		"{A(dlkgjsg) I[dgisg](kjsdglkds) [AND] (Bdir{djglksjdgkd} Cex(A(sdlgjlskd)) [XOR] A(dsgjslkj) E(gklsjgls))}" +
		"{Cac{ A(actor) I(fjhgjh) Bdir(rtyui)} [XOR] Cac{A(ertyui) I(dfghj)}}" +
		"{Cac{ A(as(dslks)a) I(adgklsjlg)} [XOR] Cac(asas) [AND] Cac12[kgkg]{lkdjgdls} [OR] A(dslgkjds)}" +
		"{Cac(andsdjsglk) [AND] A(sdjlgsl) Bdir(jslkgsjlkgds)}" +
		"{Cac(andsdjsglk) [AND] ( A(sdjlgsl) [XOR] (A(sdoidjs) [OR] A(sdjglksj)))}" +
		"((dglkdsjg [AND] jdlgksjlkgd))"


	res = r.FindAllString(text, -1)

	fmt.Println("Refined matching combinations")
	fmt.Println(res)
	fmt.Println("Count:", len(res))





	os.Exit(0)

	const LOGICAL_OPERATORS = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"




	regex := parser.LEFT_BRACE + NESTED_TERM + parser.RIGHT_BRACE +
		"\\s+" + "\\[" + LOGICAL_OPERATORS + "\\]\\s+" //+ "\\" +
		//parser.LEFT_BRACE + NESTED_TERM + "\\" + parser.RIGHT_BRACE



	//const COMPONENT_ANNOTATION_SYNTAX = "(\\[([0-9a-zA-Z" + SPECIAL_SYMBOLS + "{}\\[\\]\\(\\)])+\\])"

	//const NESTED_COMPONENT_SYNTAX = "(Bdir|Bind|Cac|Cex|E|F|P)" + parser.COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "?"

	//regex = NESTED_COMPONENT_SYNTAX + "\\" + parser.LEFT_BRACE + "(" + NESTED_TERM + ")+" + "\\" + parser.RIGHT_BRACE

	//regex = "\\[" + parser.LOGICAL_OPERATORS + "\\]"

	fmt.Println(regex)

	//fmt.Println(NESTED_TERM)

	r, err = regexp.Compile(regex)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	//text = "{Cac{E(Program Manager) F(is) P(approved)} [XOR] Cac{A(NOP Official) I(accepts) Bdir(Program Manager)}}"

	//text = "A{gldksjgk} Bdir(dgkjslkg) Cex(sdglkdsjg)"

	text = "{{A} [AND] {A}}"

	//text = "(A) [AND] {A}"

	res = r.FindAllString(text, -1)

	fmt.Println(res)

	fmt.Println(res[0][0])
	//fmt.Println(res[0][0])

	/*
	res1, err := parser.SeparateComponentsAndNestedStatements(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		fmt.Println("Error: ", err.Error())
	}

	fmt.Println("Results")
	fmt.Println(res1)
	*/

	// No shared elements
	//exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = false

	/*s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		fmt.Errorf("%v", "Error during parsing of statement", err.Error())
	}

	fmt.Println(s.String())*/

	//linkMap := parser.ExtractLinkBetweenProperties(s)

	//fmt.Println("Identified links:", linkMap)

/*
	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays()

	fmt.Println("Component References:", componentRefs)

	res, err := exporter.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		fmt.Errorf("%v", "Unexpected error during array generation.")
	}

	fmt.Println("Results:", res)
*/

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

	/*s, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		fmt.Println(err.Error())
	}
	fmt.Println(s.String())*/
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


