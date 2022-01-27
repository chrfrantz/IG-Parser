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
const LEFT_BRACKET = "["
const RIGHT_BRACKET = "]"

// Logical operators prepared for regular expression
const LOGICAL_OPERATORS = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"

// Special symbols supported in content, suffix and annotations
const SPECIAL_SYMBOLS = "',;+\\-*/%&=$£€¤§\"#!`\\|"

// Word pattern for regular expressions (including parentheses, spaces, square brackets, +, -, /, *, %, &, =, currency symbols, etc., BUT not braces!)
const WORDS_WITH_PARENTHESES = "[a-zA-Z,0-9" + SPECIAL_SYMBOLS + "\\(\\)\\[\\]\\s]+"

// Optional use of word pattern
const OPTIONAL_WORDS_WITH_PARENTHESES = "(" + WORDS_WITH_PARENTHESES + ")?"

// Pattern of parentheses combinations, e.g., ( ... [AND] ... )
const COMBINATION_PATTERN_PARENTHESES = "\\" + LEFT_PARENTHESIS + WORDS_WITH_PARENTHESES +
	"(\\[" + LOGICAL_OPERATORS + "\\]" + WORDS_WITH_PARENTHESES + ")+\\" + RIGHT_PARENTHESIS

// Pattern of brace combinations, e.g., { ... [AND] ... }, but it matches { ... } [LogicalOperator] ... } to ensure it captures component-level nesting
const COMBINATION_PATTERN_BRACES = "\\" + LEFT_BRACE + WORDS_WITH_PARENTHESES + "\\" + RIGHT_BRACE +
	"\\s+" + "(\\[" + LOGICAL_OPERATORS + "\\]\\s+" + WORDS_WITH_PARENTHESES + ")+\\" + RIGHT_BRACE

// Annotation syntax (e.g., [semanticAnnotations#99])
const COMPONENT_ANNOTATION_MAIN = "[a-zA-Z,0-9\\s" + SPECIAL_SYMBOLS + "]+"

// Nested annotation syntax (e.g., [first=[left,right]])
const COMPONENT_ANNOTATION_OPTIONAL_BRACKET = "(\\[" + COMPONENT_ANNOTATION_MAIN + "\\])*"

// Nested annotation syntax (e.g., [first=(left,right)])
const COMPONENT_ANNOTATION_OPTIONAL_PARENTHESES = "(\\(" + COMPONENT_ANNOTATION_MAIN + "\\))*"

// Nested annotation syntax (e.g., [first=[left,right]])
const COMPONENT_ANNOTATION_OPTIONAL = "(" + COMPONENT_ANNOTATION_OPTIONAL_PARENTHESES + "|" +
	COMPONENT_ANNOTATION_OPTIONAL_BRACKET + ")"

// Complete annotation syntax
const COMPONENT_ANNOTATION_SYNTAX = "(\\[(" + COMPONENT_ANNOTATION_MAIN + COMPONENT_ANNOTATION_OPTIONAL + ")+\\])?"

// Regex for component suffix (e.g., "1" in "A1")
const COMPONENT_SUFFIX_SYNTAX = "[a-zA-Z,0-9" + SPECIAL_SYMBOLS + "]*"

// Full component header syntax including suffix and annotations (e.g., A1[semanticAnnotation])
const COMPONENT_HEADER_SYNTAX = "(" +
	tree.ATTRIBUTES + "|" +
	tree.ATTRIBUTES_PROPERTY + "|" +
	tree.DEONTIC + "|" +
	tree.AIM + "|" +
	tree.DIRECT_OBJECT + "|" +
	tree.DIRECT_OBJECT_PROPERTY + "|" +
	tree.INDIRECT_OBJECT + "|" +
	tree.INDIRECT_OBJECT_PROPERTY + "|" +
	tree.ACTIVATION_CONDITION + "|" +
	tree.EXECUTION_CONSTRAINT + "|" +
	tree.CONSTITUTED_ENTITY + "|" +
	tree.CONSTITUTED_ENTITY_PROPERTY + "|" +
	tree.MODAL + "|" +
	tree.CONSTITUTIVE_FUNCTION + "|" +
	tree.CONSTITUTING_PROPERTIES + "|" +
	tree.CONSTITUTING_PROPERTIES_PROPERTY +
	")" +
	COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX

// Full syntax of components, including identifier, suffix, annotation and potentially nested or atomic content
const FULL_COMPONENT_SYNTAX =
// Component identifier, with suffix and annotations
COMPONENT_HEADER_SYNTAX +
	// Component-level nesting (e.g., { ... })
	"(\\" + LEFT_BRACE + "\\s*" + WORDS_WITH_PARENTHESES + "\\s*\\" + RIGHT_BRACE + "|" +
	// atomic component content (e.g., ( ... ))
	"\\" + LEFT_PARENTHESIS + "\\s*" + WORDS_WITH_PARENTHESES + "\\s*\\" + RIGHT_PARENTHESIS + ")"

// Full syntax of nested component, including identifier, suffix, annotation
const FULL_COMPONENT_SYNTAX_NESTED =
// Component identifier, with suffix and annotations
COMPONENT_HEADER_SYNTAX +
	// Component-level nesting (e.g., { ... })
	"\\" + LEFT_BRACE + "\\s\\*" + WORDS_WITH_PARENTHESES + "\\s\\*\\" + RIGHT_BRACE

// Basic combination of an arbitrary number of components, variably with or without parentheses (e.g., indication of precedence)
const NESTED_COMBINATION =
// Start of alternatives
"(" +
	// combination with parentheses
	"\\" + LEFT_PARENTHESIS +
	OPTIONAL_WORDS_WITH_PARENTHESES +
	"(" + FULL_COMPONENT_SYNTAX + ")+" +
	"(" +
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES +
	"(" + FULL_COMPONENT_SYNTAX + OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
	")*" +
	"\\" + RIGHT_PARENTHESIS +
	// OR
	"|" +
	// combinations without parentheses
	OPTIONAL_WORDS_WITH_PARENTHESES + "(" + FULL_COMPONENT_SYNTAX + ")+" +
	"(" +
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES + "(" + FULL_COMPONENT_SYNTAX +
	OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
	")*" +
	// END OF COMBINATION
	")"

// Inner part of nested combinations (i.e., without component syntax and/or termination) for flexible composition
const INNER_NESTED_COMBINATIONS = "\\" + LEFT_BRACE +
	"\\s*(" + NESTED_COMBINATION + "\\s+)+" +
	"(" +
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	"\\s+(" + NESTED_COMBINATION + "\\s*)+" +
	")+" +
	"\\" + RIGHT_BRACE

// Combination of combinations to represent multi-level nesting (does not require termination, i.e., could be embedded)
// Used in testing
const NESTED_COMBINATIONS = COMPONENT_ANNOTATION_SYNTAX +
	INNER_NESTED_COMBINATIONS

// Combinations of combinations for multi-combined component-level nesting, under consideration of termination for atomic matching
// Used in production
const NESTED_COMBINATIONS_TERMINATED = COMPONENT_ANNOTATION_SYNTAX +
	"^" + // Ensure the tested statement only contains combinations, but no leading individual component (i.e., combination embedded in nested statement)
	INNER_NESTED_COMBINATIONS +
	"$" // Ensure immediate termination of combination with additional trailing components (which would imply nested statement with embedded combination)

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
func CleanInput(input string) string {

	// Remove line breaks
	re := regexp.MustCompile(`\r?\n`)
	input = re.ReplaceAllString(input, " ")

	return input
}
