package main

import (
	"IG-Parser/core/parser"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

/*
Helper main function as development and debugging workbench.
This is purely to support development (quick access to examples statements with distinctive features)
This file is not included in the server-side deployment (i.e., the docker container does consider this file during build)!
*/
func main() {

	// complex baseline statement containing all patterns
	text := "D(deontic) Cac(atomicCondition) (lkjsdkljs) Bind(indirectobject) Cac{A(atomicnestedcondition)} " +
		"{I(maintain) Bdir((order [AND] control))  Cac{A(sharednestedcondition)} [XOR] {I(sustain) Bdir(peace) [OR] I(prevent) Bdir(war)}} " +
		" Cac{Cac{ A(leftcombo) I(leftaim) } [XOR] Cac{ Cac{ A(rightleftcombo) I(rightleftaim) } [AND] Cac{ A(rightrightcombo) I(rightrightaim) }}}"

	//text = "nadndasa {A(another1) Bdir(object1) Cac{A(dklj)} [OR] A(another2) Bdir(object2)}"

	// Component pair combos with atomic components
	//text = "D(dsgjslkj) Cac(sdjglksj) (lkjsdkljs) Bind(kdlsls)" + " Cac{Cac{Cac{A(precond2)} [XOR] Cac{A(precond3)}} [AND] Cac{A(another2)}}"

	// Component combination with atomic components
	//text = "D(dsgjslkj) Cac(sdjglksj) (lkjsdkljs) Bind(kdlsls)" + " { Cac{A(precond2)} [AND] {Cac{A(another2) I(dgjlskdjg)} [XOR] Cac{A(thirdone)}}}"

	// Nested component with atomic components
	//text = "D(dsgjslkj) Cac(sdjglksj) (lkjsdkljs) Bind(kdlsls)" + " Cac{A(precond2)}"

	// Component combination in component pair combination
	//text = "D(dsgjslkj) Cac(sdjglksj) (lkjsdkljs) Bind(kdlsls)" + " { Cac{ Cac{A(precond2) I(dksjld)} [XOR] Cac{A(ldkjsjg) I(sdgjls)}} [AND] Cac{A(another2)} }"

	// Basic pattern with within-component nested components
	text = "A(dskgjsl) Bdir((left [XOR] right)) Cac(aldgjslk)"

	// Component combination with gibberish
	//text = "dssgs Cac{ Cac{A(precond2) I(dksjld)} [XOR] Cac{A(ldkjsjg) I(sdgjls)}} sdsdg "

	// Component combination with embedded component-level nesting
	text = "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ A(actor2) I(aim2) Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } } }"
	//text = "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } } "

	// Pure component pair
	//text = "{ Cac{A(precond)} Bdir(leftbdir) I(leftact) [XOR] Bdir(rightbdir) I(rightact)}"

	// Combination of nested component and within-component nesting
	//text = "D(dlkgjslkj) I(dskgjslkgj) Cac{A(dlksgjs)} Bdir((LFT [XOR] RHT))"

	// Single component-level nesting
	//text = "D(dlkgjslkj) I(dskgjslkgj) Cac{A(dlksgjs)} Bdir(object)"

	// Component pair combination for distinction with component combinations
	//text = "{Cac{A(sdlkgjsdlk) Bdir{A(aljdgs) I(kdsjglkj) Bdir(dkslgj)}} [XOR] Cac{A(skdfjcs) Bdir{A(dlksgjie) I(dsklgjiv) Bdir(lkdsjgei)}}}"

	// Private property parsing
	//text = "{A,p(first) A(farmer) [OR] A(citizen)}"

	//text = "Bdir{A1,p(first) (A1(farmer) [OR] A2(citizen))}"
	//text = "Bdir{A2,p(first) A1(farmer) [OR] A2(citizen)}"
	//text = "{A(sdkjls) Bdir((LFT [XOR] RHT)) [OR] A(ertyu) I(dgsdkjg)}"

	// Multi-level nesting
	//text = "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}} dfghj" +
	//	" Cac{ Cac{A(precond2) I(dksjld)} [XOR] Cac{A(ldkjsjg) I(sdgjls)}} " +
	//	"{ Cac{A(precond2) I(dksjld)} [XOR] Cac{A(ldkjsjg) I(sdgjls)}}"

	// Standard combinations

	text = "A(actor0) Cac{ A(actor) I(aim) Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } I(aim2) }"

	// same as above, but proper nesting
	text = "A(actor0) Cac{ A(actor) I(aim) Cac{Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) }} Bdir(obj2) }"

	// Combinations-embedding nested component example
	//text = "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } }"

	text = "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ A(actorX) I(aim2) Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } }}"

	// Short embedded nested node
	//text = "Cac{ A(actorz) Cac{ A(actorX) Cac{ Cac{ A(actor3) I(aim3) } } [OR] Cac{ A(actor4) I(aim4) } } }"

	// Short embedded nested node (should not work)
	text = "Cac{ A(actorz) Cac{ A(actorX) Cac{ Cac{ A(actor3) I(aim3) } } [OR] Cac{ A(actor4) I(aim4) } Bdir(object2) } }"

	// embedded node (works on level 2 & 3 - only inner nested element) - correct
	text = "Cac{ A(actorz) Cac{ A(actorX) Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } Bdir(object2) } }"

	// Combinations-embedding nested component example (should not work)
	//text = "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ A(actorX) Cac{ Cac{ A(actor3) I(aim3) } } [OR] Cac{ A(actor4) I(aim4) } } }"

	// Combinations-embedding nested component example (should work including left actor1 part! & OR!)
	text = "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ A(actorX) Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } } }"

	// DEBUG Combinations-embedding nested component example (should work including left actor1 part! & OR!)
	//text = "Cac{ A(actorX) Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } }"
	// imbalanced right-double-nested
	//text = "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } } }"
	// triple right weighted
	text = "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } }"
	// quadruple right weighted
	text = "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } [OR] Cac{ A(actor5) I(aim5) } } }"
	// balanced
	text = "Cac{ Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ A(actor1) I(aim1) } } [OR] Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } }"
	// right-nested balanced
	text = "Cac{ Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ A(actor1) I(aim1) } } [OR] Cac{ Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } } }"

	// test for scale
	//text = "Cac{ Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ A(actor1) I(aim1) } } [OR] Cac{ A(skdlj) Cac{ Cac{ A(actor3) I(aim3) } [AND] Cac{ A(actor4) I(aim4) [AND] Cac{ A(actor4) I(aim4) } } } } }"

	// left narrow scope (works on 4) - matches, but does not parse
	//text = "Cac{ Cac{ A(lefty) Cac{ Cac{ A(actor1) I(aim1) }  Bdir(oxo) } [AND] Cac{ A(actor1) I(aim1) } } [OR] Cac{ A(actorX) Cac{ Cac{ A(actor3) I(aim3) } } [OR] Cac{ A(actor4) I(aim4) } } }"

	// works
	//text = "Cac{ Cac{ Cac{ A(lefty) Cac{ A(actor1) I(aim1) }  Bdir(oxo) } [AND] Cac{ A(actor1) I(aim1) } } [OR] Cac{ A(actorX) Cac{ Cac{ A(actor3) I(aim3) }  [OR] Cac{ A(actor4) I(aim4) } } } }"

	// left broad scope
	//text = "Cac{ Cac{ A(lefty) Cac{ Cac{ A(actor1) I(aim1) } [AND] Cac{ A(actor1) I(aim1) } } } [OR] Cac{ A(actorX) Cac{ Cac{ A(actor3) I(aim3) } } [OR] Cac{ A(actor4) I(aim4) } } }"

	// Combinations-embedding nested component example
	//text = "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ A(actor2) I(aim2) Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } Bdir(object2) } }"

	// Nested properties example
	//text = "Such E(notification) M(shall) F(provide): (1) A P(description of each noncompliance); (2) The P(facts upon which the notification of noncompliance is based); and (3) The P1(date) P1,p{by which the A(certified operation) D(must) {I(rebut [XOR] correct) Bdir,p(each) Bdir(noncompliance) [AND] I(submit) Bdir,p(supporting) Bdir(documentation) of Bdir,p(each such correction) Cac(when correction is possible)}} P1,p(private component)"

	// Nested properties example 2
	//text = "Such E(notification) M(shall) F(provide): (1) A P(description of each noncompliance); (2) The P(facts upon which the notification of noncompliance is based); and (3) The P1(date) P1,p{by which the A(certified operation) D(must) {I(rebut [XOR] correct) Bdir,p(each) Bdir(noncompliance) [AND] I(submit) Bdir,p(supporting) Bdir(documentation) of Bdir,p(each such correction) Cac(when correction is possible)}} P1,p(private component) P1,p{where E(date) F(is defined) in the P(Gregorian calendar)}."

	// Regex
	/*const FULL_COMPONENT_SYNTAX_WITH_NESTED_COMPONENTS_AND_NESTED_COMBINATIONS3 = "^" +
	// Component identifier, with suffix and annotations
	parser.COMPONENT_HEADER_SYNTAX +
	// component-level nesting (e.g., { ... }), including potentially embedded second-order nesting on component(s)
	"(\\" + parser.LEFT_BRACE + "\\s*" + "(" + parser.WORDS_WITH_PARENTHESES + "|" + parser.WORDS_WITH_PARENTHESES + parser.FULL_COMPONENT_SYNTAX + "|\\s*)*(" + parser.COMPONENT_HEADER_SYNTAX + parser.BRACED_7TH_ORDER_COMBINATIONS + ")+" + "\\s*\\" +
	"(" + parser.WORDS_WITH_PARENTHESES + "|" + parser.WORDS_WITH_PARENTHESES + parser.FULL_COMPONENT_SYNTAX + ")*\\s*(" + parser.RIGHT_BRACE + ")" +
	"$"*/
	fmt.Println("Candidate: " + text)

	const NESTED_COMBINATIONS_TERMINATED =
	// Component combinations need to lead with component identifier (and potential suffix and annotation)
	//"^" +
	parser.COMPONENT_HEADER_SYNTAX +
		// Ensure the tested statement only contains combinations, but no leading individual component (i.e., combination embedded in nested statement)
		//parser.BRACED_2ND_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMPONENTS //+
		parser.BRACED_3RD_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMPONENTS
	//parser.BRACED_4TH_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS
	//parser.BRACED_5TH_ORDER_COMBINATIONS
	//parser.BRACED_7TH_ORDER_COMBINATIONS +
	//"$" //Ensure immediate termination of combination with additional trailing components (which would imply nested statement with embedded combination)

	const FULL_COMPONENT_SYNTAX_WITH_NESTED_COMPONENTS_AND_NESTED_COMBINATIONS3 =
	// Component identifier, with suffix and annotations
	"^" +
		parser.COMPONENT_HEADER_SYNTAX +
		// component-level nesting (e.g., { ... }), including potentially embedded second-order nesting on component(s)
		"(\\" + parser.LEFT_BRACE +
		"\\s*" +
		"(" + parser.WORDS_WITH_PARENTHESES +
		//"|" +
		/*parser.COMPONENT_HEADER_SYNTAX + */ //"\\" + parser.LEFT_BRACE + "\\s*" + /*parser.COMPONENT_HEADER_SYNTAX + "?\\s*"*/ parser.BRACED_7TH_ORDER_COMBINATIONS + "\\s*" + "\\" + parser.RIGHT_BRACE +
		"|" +
		parser.COMPONENT_HEADER_SYNTAX + "?" + parser.BRACED_7TH_ORDER_COMBINATIONS +
		"|" +
		"\\s*" + ")*" + "\\s*\\" + parser.RIGHT_BRACE + ")" +
		"$"
	// Regex test

	//r, err := regexp.Compile(parser.FULL_COMPONENT_SYNTAX_WITH_NESTED_COMPONENTS_AND_NESTED_COMBINATIONS)

	r, err := regexp.Compile(FULL_COMPONENT_SYNTAX_WITH_NESTED_COMPONENTS_AND_NESTED_COMBINATIONS3)

	r, err = regexp.Compile(NESTED_COMBINATIONS_TERMINATED)

	//r, err = regexp.Compile(parser.FULL_COMPONENT_SYNTAX_WITH_NESTED_COMPONENTS) //parser.PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS) //parser.NESTED_COMBINATIONS)
	if err != nil {
		log.Fatal("Error during compilation:", err.Error())
	}

	res := r.FindAllString(text, -1)

	for i, val := range res {
		fmt.Println("Item " + strconv.Itoa(i) + ": " + val)
	}
	fmt.Println("Found " + strconv.Itoa(len(res)) + " items.")

	/*tabular.SetIncludeSharedElementsInTabularOutput(true)
	tabular.SetDynamicOutput(false)

	fmt.Println("Shared mode:")
	fmt.Println(tree.SHARED_ELEMENT_INHERITANCE_MODE)

	tree.SetFlatPrinting(true)*/

	/*_, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}*/

	//endpoints.ConvertIGScriptToTabularOutput("", text, "1.1", tabular.OUTPUT_TYPE_CSV, "example.csv", true, true, tabular.ORIGINAL_STATEMENT_OUTPUT_NONE, tabular.IG_SCRIPT_OUTPUT_NONE)

	//fmt.Println(output)

}
