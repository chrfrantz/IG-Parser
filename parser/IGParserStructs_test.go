package parser

import "testing"

/*
Tests the cleaning of input string from line breaks
 */
func TestCleanInput(t *testing.T) {

	input := "Program Manager\n has objectives \r\n and we have variable input \n\n    to clean.\n\n\n"

	input = CleanInput(input)

	expectedOutput := "Program Manager  has objectives   and we have variable input       to clean.   "

	if input != expectedOutput {
		t.Fatal("Preprocessing did not work according to expectation.")
	}
}