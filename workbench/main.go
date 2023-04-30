package main

import (
	"IG-Parser/core/endpoints"
	"IG-Parser/core/exporter"
	"IG-Parser/core/tree"
	"fmt"
)

/*
Helper main function as development and debugging workbench.
This is purely to support development and not included in the deployment!
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
	//text = "A(dskgjsl) Bdir((left [XOR] right)) Cac(aldgjslk)"

	// Component combination with gibberish
	//text = "dssgs Cac{ Cac{A(precond2) I(dksjld)} [XOR] Cac{A(ldkjsjg) I(sdgjls)}} sdsdg "

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

	exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	exporter.SetDynamicOutput(false)

	fmt.Println("Shared mode:")
	fmt.Println(tree.SHARED_ELEMENT_INHERITANCE_MODE)

	tree.SetFlatPrinting(true)

	/*_, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}*/

	endpoints.ConvertIGScriptToTabularOutput(text, "1.1", exporter.OUTPUT_TYPE_CSV, "example.csv", true, true)

	//fmt.Println(output)

}
