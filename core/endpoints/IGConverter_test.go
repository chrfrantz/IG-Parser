package endpoints

import (
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/core/tree"
	"testing"
)

// GOOGLE SHEETS OUTPUT

/*
Tests basic valid statement for Google sheets output.
*/
func TestValidStatementGoogleSheets(t *testing.T) {

	originalStatement := "Original text that doesn't matter for this test."

	text := "(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	_, err := ConvertIGScriptToTabularOutput(originalStatement, text, "650", tabular.OUTPUT_TYPE_GOOGLE_SHEETS, "", true, true, tabular.ORIGINAL_STATEMENT_OUTPUT_NONE, tabular.IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Statement parsing should not fail")
	}
}

/*
Tests simple statement-level nesting on activation condition for Google Sheets output.
*/
func TestValidStatementNestingGoogleSheets(t *testing.T) {

	originalStatement := "Original text that doesn't matter for this test."

	text := "(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// This is the essential line
		"Cac{A(Program Manager) I(has gained) Bdir(competence)}"

	_, err := ConvertIGScriptToTabularOutput(originalStatement, text, "650", tabular.OUTPUT_TYPE_GOOGLE_SHEETS, "", true, true, tabular.ORIGINAL_STATEMENT_OUTPUT_NONE, tabular.IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Statement parsing should not fail")
	}
}

/*
Tests invalid attribute combinations for Google Sheets output.
*/
func TestInvalidAttributeStatementGoogleSheets(t *testing.T) {

	originalStatement := "Original text that doesn't matter for this test."

	// Statement with invalid attribute combination
	text := "A((certifying agent [AND] borrower [OR] wife)) M(may) I(investigate) " +
		"Bdir((complaints of noncompliance with the (Act or [OR] regulations of this part) " +
		"concerning " +
		"(production [operations] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
		")" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	_, err := ConvertIGScriptToTabularOutput(originalStatement, text, "650", tabular.OUTPUT_TYPE_GOOGLE_SHEETS, "", true, true, tabular.ORIGINAL_STATEMENT_OUTPUT_NONE, tabular.IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode == tree.PARSING_NO_ERROR {
		t.Fatal("Statement parsing should produce error")
	}
}

// VISUAL OUTPUT

/*
Tests basic valid statement for visual tree output
*/
func TestValidStatementVisualTree(t *testing.T) {
	text := "(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	_, err := ConvertIGScriptToVisualTree(text, "650", "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Statement parsing should not fail")
	}
}

/*
Tests simple statement-level nesting on activation condition for visual tree output
*/
func TestValidStatementNestingVisualTree(t *testing.T) {
	text := "(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// This is the essential line
		"Cac{A(Program Manager) I(has gained) Bdir(competence)}"

	_, err := ConvertIGScriptToVisualTree(text, "650", "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Statement parsing should not fail")
	}
}

/*
Tests invalid attribute combinations for visual tree output.
*/
func TestInvalidAttributeStatementVisualTree(t *testing.T) {

	// Statement with invalid attribute combination
	text := "A((certifying agent [AND] borrower [OR] wife)) M(may) I(investigate) " +
		"Bdir((complaints of noncompliance with the (Act or [OR] regulations of this part) " +
		"concerning " +
		"(production [operations] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
		")" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	_, err := ConvertIGScriptToVisualTree(text, "650", "")
	if err.ErrorCode == tree.PARSING_NO_ERROR {
		t.Fatal("Statement parsing should produce error")
	}
}

/*
Tests degree of variability for visual tree output.
*/
func TestDegreeOfVariabilityVisualTree(t *testing.T) {

	// Statement with various complexity
	text := "A(certifying agent [AND] (borrower [OR] wife)) M(may) I(investigate) " +
		"Bdir((complaints of noncompliance with the (Act [OR] regulations of this part) " +
		"concerning " +
		"(production [operations] and [AND] handling operations) as well as (shipping [XOR] packing facilities)) " +
		")" +
		"Cex(for compliance with the (Act [XOR] regulations in this part))."

	// Activate degree of variability in output
	tabular.SetIncludeDegreeOfVariability(true)

	_, err := ConvertIGScriptToVisualTree(text, "650", "")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Statement parsing should not fail")
	}
}
