package parser

import (
	"fmt"
	"testing"
)

func TestStatementParsing(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	fmt.Println(text)

	s := ParseStatement(text)

	fmt.Println(s.String())

	/*expectedOutput := "A: Leaf entry: National Organic Program's Program Manager\n        D: Leaf entry: may\n        I: (\n        Left: Leaf entry: inspect and\n        Operator: AND\n        Right: (\n        ====Shared (left): sustain\n        ====Left: Leaf entry: review\n        ====Operator: AND\n        ====Right: (\n        ========Shared (left): sustain\n        ========Left: Leaf entry: refresh\n        ========Operator: AND\n        ========Right: Leaf entry: drink\n        ========)\n        ====)\n        )\n        Bdir: (\n        Left: (\n        ====Left: Leaf entry: certified production and\n        ====Operator: AND\n        ====Right: Leaf entry: handling operations and\n        ====)\n        Operator: AND\n        Right: Leaf entry: accredited certifying agents\n        )\n        Cex: (\n        Left: Leaf entry: on behalf of the Secretary\n        Operator: AND\n        Right: (\n        ====Shared (left): for compliance with the\n        ====Left: Leaf entry: Act or\n        ====Operator: XOR\n        ====Right: Leaf entry: regulations in this part\n        ====)\n        )\n        "

	if s.String() != expectedOutput {
		t.Fatal("Statement has not been correctly parsed. Output: '" + s.String() + "'.")
	}*/

	// Assess whether it is correctly parsed

}


func TestSyntheticComponentCombinations(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	fmt.Println(text)

	s := ParseStatement(text)

	fmt.Println(s.String())

	// Count number of elements and check correspondence to references map

	/*expectedOutput := "A: Leaf entry: National Organic Program's Program Manager\n        D: Leaf entry: may\n        I: (\n        Left: Leaf entry: inspect and\n        Operator: AND\n        Right: (\n        ====Shared (left): sustain\n        ====Left: Leaf entry: review\n        ====Operator: AND\n        ====Right: (\n        ========Shared (left): sustain\n        ========Left: Leaf entry: refresh\n        ========Operator: AND\n        ========Right: Leaf entry: drink\n        ========)\n        ====)\n        )\n        Bdir: (\n        Left: (\n        ====Left: Leaf entry: certified production and\n        ====Operator: AND\n        ====Right: Leaf entry: handling operations and\n        ====)\n        Operator: AND\n        Right: Leaf entry: accredited certifying agents\n        )\n        Cex: (\n        Left: Leaf entry: on behalf of the Secretary\n        Operator: AND\n        Right: (\n        ====Shared (left): for compliance with the\n        ====Left: Leaf entry: Act or\n        ====Operator: XOR\n        ====Right: Leaf entry: regulations in this part\n        ====)\n        )\n        "

	if s.String() != expectedOutput {
		t.Fatal("Statement has not been correctly parsed. Output: '" + s.String() + "'.")
	}*/

	// Assess whether it is correctly parsed

}