package parser

import (
	"IG-Parser/tree"
	"regexp"
	"strings"
)

// Define constants for parentheses and braces
const LEFT_PARENTHESIS = "("
const RIGHT_PARENTHESIS = ")"
const LEFT_BRACE = "{"
const RIGHT_BRACE = "}"

// Logical operators prepared for regular expression
const logicalOperators = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"
// Special symbols supported in content, suffix and annotations
const specialSymbols = "',;+\\-*/%&=$£€¤§\"#!`\\|"
// Word pattern for regular expressions (including parentheses, spaces, square brackets, +, -, /, *, %, &, =, currency symbols, etc.)
const wordsWithParentheses = "([a-zA-Z" + specialSymbols + "(){}\\[\\]]+\\s*)+"
// Pattern of combinations, e.g., ( ... [AND] ... )
const combinationPattern = "\\(" + wordsWithParentheses + "(\\[" + logicalOperators + "\\]\\s" + wordsWithParentheses + ")+\\)"

// Regex for component annotations (e.g., semantic labels, such as "[monitor]" in "A[monitor](Program Manager)")
const COMPONENT_ANNOTATION_SYNTAX = "(\\[([0-9a-zA-Z" + specialSymbols + "\\s{}\\[\\]\\(\\)])+\\])"
// Regex for component suffix (e.g., "1" in "A1")
const COMPONENT_SUFFIX_SYNTAX = "[a-zA-Z,0-9" + specialSymbols + "]*"
// Regex for complete component
//const COMPONENT_SYNTAX_COMPLETE = "([a-zA-Z,]+)+" + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX

// Nested component prefix (word without spaces and parentheses, but [] brackets) TODO: Check for integration with other Suffix and Annotation regex
var NESTED_COMPONENT_SYNTAX = "([a-zA-Z{}\\[\\]]+)+"

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
	text = strings.ReplaceAll(text, "|", "\\|")

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
