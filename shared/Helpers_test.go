package shared

import "testing"

/*
Tests the replacing of symbols for export.
*/
func TestReplaceSymbolsForExport(t *testing.T) {

	input := "\"left\" middle next \"center\" next \"right\""

	// Check whether escaping worked as expected
	if EscapeSymbolsForExport(input) != "'left' middle next 'center' next 'right'" {
		t.Error("Escaping of values did not result in expected outcome.")
	}
}
