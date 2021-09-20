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
//const WORDS_WITH_PARENTHESES = "([a-zA-Z" + SPECIAL_SYMBOLS + "()\\[\\]]+\\s*)+"
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
//const COMPONENT_ANNOTATION_OPTIONAL = "(\\[" + COMPONENT_ANNOTATION_MAIN + "\\])*"
const COMPONENT_ANNOTATION_OPTIONAL = "(" + COMPONENT_ANNOTATION_OPTIONAL_PARENTHESES + "|" +
										COMPONENT_ANNOTATION_OPTIONAL_BRACKET + ")"

// Complete annotation syntax
const COMPONENT_ANNOTATION_SYNTAX = "(\\[(" + COMPONENT_ANNOTATION_MAIN + COMPONENT_ANNOTATION_OPTIONAL + ")+\\])?"

// Regex for component annotations (e.g., semantic labels, such as "[monitor]" in "A[monitor](Program Manager)") - no whitespace is allowed in annotations
//const COMPONENT_ANNOTATION_SYNTAX = "(\\[([0-9a-zA-Z" + SPECIAL_SYMBOLS + "{}\\[\\]\\(\\)])+\\])"

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
	"\\" + LEFT_PARENTHESIS + "\\s*" + WORDS_WITH_PARENTHESES + "\\s*\\" + RIGHT_PARENTHESIS +")"

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
		"(" + FULL_COMPONENT_SYNTAX /*+ OPTIONAL_WORDS_WITH_PARENTHESES*/ + ")+" +
		//"\\s*(" + NESTED_TERM + "\\s*)+" +
		//"\\s*" + RIGHT_BRACE + "\\s*" +
		//"\\s*" +
		"(" +
		"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
		OPTIONAL_WORDS_WITH_PARENTHESES +
		"(" + FULL_COMPONENT_SYNTAX + OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
		//"\\s+(" + NESTED_TERM + "\\s*)+" +
		")*" +
		"\\" + RIGHT_PARENTHESIS +
	// OR
	"|" +
		// combinations without parentheses
		OPTIONAL_WORDS_WITH_PARENTHESES + "(" + FULL_COMPONENT_SYNTAX +
			/*OPTIONAL_WORDS_WITH_PARENTHESES +*/ ")+" +
		//"\\s*(" + NESTED_TERM + "\\s*)+" +
		//"\\s*" + RIGHT_BRACE + "\\s*" +
		//"\\s*" +
		"(" +
		"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
		OPTIONAL_WORDS_WITH_PARENTHESES + "(" + FULL_COMPONENT_SYNTAX +
			OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
		//"\\s+(" + NESTED_TERM + "\\s*)+" +
		")*" +
	// END OF COMBINATION
	")"

// Combination of combinations to represent multi-level nesting
const NESTED_COMBINATIONS =
	COMPONENT_ANNOTATION_SYNTAX +
	"\\" + LEFT_BRACE +
	"\\s*(" + NESTED_COMBINATION + "\\s+)+" +
	"(" +
	//"\\s*" + RIGHT_BRACE + "\\s*" +
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	"\\s+(" + NESTED_COMBINATION + "\\s*)+" +
	")+" +
	"\\" + RIGHT_BRACE //+ "\\s*" +
	//"\\" + RIGHT_BRACE


const NESTED_COMBINATIONS2 =
		// leading annotation on combination
		//COMPONENT_ANNOTATION_SYNTAX +
			"\\" + LEFT_BRACE +
			".*" +
			"(" +
			"\\" + RIGHT_BRACE + "\\s*\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
			".*" +
			")+" +
			"\\" + RIGHT_BRACE + "\\" + RIGHT_BRACE


// Regex for complete component
//const COMPONENT_SYNTAX_COMPLETE = "([a-zA-Z,]+)+" + COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX

// Nested component prefix (word without spaces and parentheses, but [] brackets) TODO: Check for integration with other Suffix and Annotation regex
//const NESTED_COMPONENT_SYNTAX = COMPONENT_SUFFIX_SYNTAX + COMPONENT_ANNOTATION_SYNTAX + "?"
//"([a-zA-Z{}\\[\\]]+)+"

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
