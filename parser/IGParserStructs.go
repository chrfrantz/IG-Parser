package parser

import "regexp"

// Define constants for parentheses and braces
const LEFT_PARENTHESIS = "("
const RIGHT_PARENTHESIS = ")"
const LEFT_BRACE = "{"
const RIGHT_BRACE = "}"


/*
Generic function to clean input (e.g., strip line breaks).
*/
func cleanInput(input string) string {

	// Remove line breaks
	re := regexp.MustCompile(`\r?\n`)
	input = re.ReplaceAllString(input, " ")

	return input
}
