package app

import (
	"IG-Parser/tree"
	"testing"
)

func TestValidStatement(t *testing.T) {
	text := "(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	_, err := ConvertIGScriptToGoogleSheets(text, "650", "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Statement parsing should not fail")
	}
}

func TestValidStatementNesting(t *testing.T) {
	text := "(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{A(Program Manager) I(has gained) Bdir(competence)}"

	_, err := ConvertIGScriptToGoogleSheets(text, "650", "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Statement parsing should not fail")
	}
}

func TestInvalidAttributeStatement(t *testing.T) {

	// Statement with invalid attribute combination
	text := "A((certifying agent [AND] borrower [OR] wife)) M(may) I(investigate) " +
		"Bdir((complaints of noncompliance with the (Act or [OR] regulations of this part) " +
		"concerning " +
		"(production [operations] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
		")"+
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	_, err := ConvertIGScriptToGoogleSheets(text, "650", "")
	if err.ErrorCode == tree.PARSING_NO_ERROR {
		t.Fatal("Statement parsing should produce error")
	}
}