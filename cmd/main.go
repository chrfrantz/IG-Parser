package main

import (
	"IG-Parser/exporter"
	"IG-Parser/parser"
	"fmt"
)

//var words = "([a-zA-Z',;]+\\s*)+"
//var wordsWithParentheses = "([a-zA-Z',;()]+\\s*)+"
//var logicalOperators = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"

func main2() {

	//text := "National Organic Program's Program Manager, on behalf of the Secretary, may (inspect and [AND] review) (certified production and [AND] handling operations and [AND] accredited certifying agents) for compliance with the (Act or [XOR] regulations in this part)."

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	/*text = "A((certifying agent [AND] wife)) D(may) I(investigate) " +
		"Bdir((complaints of noncompliance with the (Act or [OR] regulations of this part) " +
		"concerning " +
		"(production [operation] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
		")"*/
		//"fdlkgjdflg))" // certified as organic by the certifying agent))."

	s := parser.ParseStatement(text)

	fmt.Println("Printing stuff: ")

	fmt.Println(s.String())

}

func main() {
	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	text = "A((certifying agent [AND] borrower [AND] wife)) D(may) I(investigate) " +
	"Bdir((complaints of noncompliance with the (Act or [OR] regulations of this part) " +
	"concerning " +
	"(production [operations] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
	")"+
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	//"fdlkgjdflg))" // certified as organic by the certifying agent))."

	s := parser.ParseStatement(text)

	fmt.Println(s.String())

	//os.Exit(0)
	leafArrays, componentRefs := s.GenerateLeafArrays()

	//fmt.Println(componentRefs)

	//os.Exit(0)

	res := exporter.GenerateNodeArrayPermutations(leafArrays...)

	fmt.Println("Component references: ", componentRefs)

	links := exporter.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println(links)

	//os.Exit(0)

	output := exporter.GenerateGoogleSheetsOutput(res, componentRefs, links, "650")

	//fmt.Println("\n" + output)

	exporter.WriteToFile("statement.csv", output)



	/*for i, s := range res {
		fmt.Println("Statement ", i, ": ", s)
		for v := range s {
			fmt.Println("-->", s[v])
		}
	}*/

}


