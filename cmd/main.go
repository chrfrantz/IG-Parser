package main

import (
	"IG-Parser/app"
	"IG-Parser/exporter"
	"IG-Parser/parser"
	"IG-Parser/tree"
	"fmt"
	"log"
	"regexp"
	"strings"
)

//var words = "([a-zA-Z',;]+\\s*)+"
//var wordsWithParentheses = "([a-zA-Z',;()]+\\s*)+"
//var logicalOperators = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"

func main() {

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


