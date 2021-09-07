package parser

import (
	"regexp"
	"strings"
)

// Define constants for parentheses and braces
const LEFT_PARENTHESIS = "("
const RIGHT_PARENTHESIS = ")"
const LEFT_BRACE = "{"
const RIGHT_BRACE = "}"

// Regex for component annotations (e.g., semantic labels, such as "[monitor]" in "A[monitor](Program Manager)")
const COMPONENT_ANNOTATION_SYNTAX = "(\\[([a-zA-Z,=;\\w{}\\[\\]\\(\\)])+\\])?"
// Regex for component suffix (e.g., "1" in "A1")
const COMPONENT_SUFFIX_SYNTAX = "[a-zA-Z,0-9]*"
// Regex for complete component
const COMPONENT_COMPONENT_PREFIX_SYNTAX = "([a-zA-Z,]+)+" + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX

/*
Escapes all special symbols to prepare those for input into regex expression
*/
func escapeSymbolsForRegex(text string) string {
	text = strings.ReplaceAll(text, "{", "\\{")
	text = strings.ReplaceAll(text, "}", "\\}")
	text = strings.ReplaceAll(text, "(", "\\(")
	text = strings.ReplaceAll(text, ")", "\\)")
	text = strings.ReplaceAll(text, "[", "\\[")
	text = strings.ReplaceAll(text, "]", "\\]")
	text = strings.ReplaceAll(text, "$", "\\$")
	text = strings.ReplaceAll(text, "+", "\\+")

	return text
}

/*
Generic function to clean input (e.g., strip line breaks).
*/
func cleanInput(input string) string {

	// Remove line breaks
	re := regexp.MustCompile(`\r?\n`)
	input = re.ReplaceAllString(input, " ")

	return input
}
