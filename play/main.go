package main

import (
	"IG-Parser/core/exporter"
	"IG-Parser/core/parser"
	"IG-Parser/core/tree"
	"fmt"
	"log"
)

/*
Helper main function for flexible adaptation during development.
This is purely to support development and not included in the deployment.
*/
func main() {

	text := "D(dsgjslkj) Cac(sdjglksj) (lkjsdkljs) Bind(kdlsls)" +
		"{I(maintain) Bdir((order [AND] control))  Cac{A(precond)} [XOR] {I(sustain) Bdir(peace) [OR] I(prevent) Bdir(war)}} " +
		"nadndasa {A(another1) Bdir(object1) [OR] A(another2) Bdir(object2)}"

	//text = "{ Cac{A(precond)} Bdir(leftbdir) I(leftact) [XOR] Bdir(rightbdir) I(rightact)}"
	//text := "D(dlkgjslkj) I(dskgjslkgj) Cac{A(dlksgjs)} Bdir((LFT [XOR] RHT))"
	//text := "{A(sdkjls) Bdir((LFT [XOR] RHT)) [OR] A(sdkjfs)}"

	exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	exporter.SetDynamicOutput(false)

	fmt.Println("Shared mode:")
	fmt.Println(tree.SHARED_ELEMENT_INHERITANCE_MODE)

	tree.SetFlatPrinting(true)

	_, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}

	//fmt.Println(output)

}
