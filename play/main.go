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

	text := "Bdir{A(actor2) I(aim2) (((Cac{A(actor3) I(aim3) Bdir(something)} [OR] Cac{A(actor4) I(aim4) Bdir(something else)}) [AND] Cac{A(actor5)}))}"

	exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	exporter.SetDynamicOutput(false)

	fmt.Println("Shared mode:")
	fmt.Println(tree.SHARED_ELEMENT_INHERITANCE_MODE)

	tree.SetFlatPrinting(true)

	output, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		log.Fatal(err.Error())
	}

	fmt.Println(output.String())

}
