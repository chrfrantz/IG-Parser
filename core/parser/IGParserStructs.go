package parser

import (
	"IG-Parser/core/tree"
	"strings"
)

/*
This file contains the main regular expressions used in the statement
parsing process.
*/

// Define constants for parentheses and braces
const LEFT_PARENTHESIS = "("
const RIGHT_PARENTHESIS = ")"
const LEFT_BRACE = "{"
const RIGHT_BRACE = "}"
const LEFT_BRACKET = "["
const RIGHT_BRACKET = "]"

// Logical operators prepared for regular expression
const LOGICAL_OPERATORS = "(" + tree.AND + "|" + tree.OR + "|" + tree.XOR + ")"

// Special symbols supported in content, suffix and annotations (wide range of special characters, as well as +, -, /, *, %, &, =, currency symbols, periods (.), relative operators (<,>), etc., BUT not braces!)
const SPECIAL_SYMBOLS = "'’,;.<>+~:\\-*/%&=@$£€¤§\"#!`\\|"

// Word pattern for regular expressions (including parentheses, spaces, square brackets, and all symbols contained in SPECIAL_SYMBOLS).
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

// Regex for component identifier
const COMPONENT_IDENTIFIER = "(" +
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
	tree.CONSTITUTING_PROPERTIES_PROPERTY + "|" +
	tree.OR_ELSE +
	")"

// Full component header syntax including identifier, suffix and annotations (e.g., A1[semanticAnnotation]), where suffix and annotation are optional
const COMPONENT_HEADER_SYNTAX =
// Component identifier, ...
COMPONENT_IDENTIFIER +
	// ... followed by optional suffix ...
	COMPONENT_SUFFIX_SYNTAX +
	// ... followed by optional annotation
	COMPONENT_ANNOTATION_SYNTAX

// Full syntax of components, including identifier, suffix, annotation and potentially nested or atomic content (but without consideration of embedded component-level nesting)
const FULL_COMPONENT_SYNTAX =
// Component identifier, with suffix and annotations
COMPONENT_HEADER_SYNTAX +
	// component-level nesting (e.g., { ... })
	"(\\" + LEFT_BRACE + "\\s*" + WORDS_WITH_PARENTHESES + "\\s*\\" + RIGHT_BRACE + "|" +
	// atomic component content (e.g., ( ... ))
	"\\" + LEFT_PARENTHESIS + "\\s*" + WORDS_WITH_PARENTHESES + "\\s*\\" + RIGHT_PARENTHESIS + ")"

// Full syntax of components, including identifier, suffix, annotation and potentially multi-level nested or atomic content, with consideration of component-level nesting embedded within expression
const FULL_COMPONENT_SYNTAX_WITH_NESTED_COMPONENTS =
// Component identifier, with suffix and annotations
COMPONENT_HEADER_SYNTAX +
	// component-level nesting (e.g., { ... }), including potentially embedded second-order nesting on component(s)
	"(\\" + LEFT_BRACE + "\\s*" + "(" + WORDS_WITH_PARENTHESES + "|" + WORDS_WITH_PARENTHESES + FULL_COMPONENT_SYNTAX + ")" + "\\s*\\" + RIGHT_BRACE + "|" +
	// atomic component content (e.g., ( ... )), including potentially embedded second-order nesting on component(s)
	"\\" + LEFT_PARENTHESIS + "\\s*" + "(" + WORDS_WITH_PARENTHESES + "|" + WORDS_WITH_PARENTHESES + FULL_COMPONENT_SYNTAX + ")" + "\\s*\\" + RIGHT_PARENTHESIS + ")"

// Full syntax of nested component, including identifier, suffix, annotation
const FULL_COMPONENT_SYNTAX_NESTED =
// Component identifier, with suffix and annotations
COMPONENT_HEADER_SYNTAX +
	// component-level nesting (e.g., { ... })
	"\\" + LEFT_BRACE + "\\s\\*" + WORDS_WITH_PARENTHESES + "\\s\\*\\" + RIGHT_BRACE

// Basic combination of an arbitrary number of components, variably with or without parentheses (e.g., indication of precedence)
const PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS =
// Start of alternatives
"(" +
	// combination with surrounding parentheses, e.g., '( some words Cac{ ... } [AND] Cac{ ... } [AND] Cac{ ... } ... )', or variably containing Cac ( ... ) for each element
	"\\" + LEFT_PARENTHESIS +
	OPTIONAL_WORDS_WITH_PARENTHESES + "(" + FULL_COMPONENT_SYNTAX_WITH_NESTED_COMPONENTS + OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
	"(" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"(" + FULL_COMPONENT_SYNTAX_WITH_NESTED_COMPONENTS + OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
	OPTIONAL_WORDS_WITH_PARENTHESES +
	")*" +
	"\\" + RIGHT_PARENTHESIS +
	// OR
	"|" +
	// combinations without surrounding parentheses, e.g., some words 'Cac{ ... } [AND] Cac{ ... } ...' (arbitrary length, but no closing parentheses)
	OPTIONAL_WORDS_WITH_PARENTHESES + "(" + FULL_COMPONENT_SYNTAX_WITH_NESTED_COMPONENTS + ")+" +
	"(" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"(" + FULL_COMPONENT_SYNTAX_WITH_NESTED_COMPONENTS + OPTIONAL_WORDS_WITH_PARENTHESES + ")+" +
	")*" +
	// END OF COMBINATION
	")"

// TODO NOTE: From hereon it gets hacky, since multi-level nesting is represented iteratively to establish higher-order nesting - no recursion in regex; needs to be expanded per demand

// 2nd order braced combinations of combinations
// (Inner part of nested combinations, including single combination or multiple combination pairs on either side
// (e.g., { {Cac{ ... } [AND] Cac{ ... } } [XOR] { {Cac{ ... } [AND] Cac{ ... } }}),
// but without leading component syntax and/or termination for flexible composition)
const BRACED_2ND_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMPONENTS =
// Optional leading component identifier
COMPONENT_IDENTIFIER + "?" +
	// Leading brace
	"\\" + LEFT_BRACE +
	// Testing of potential excessive words preceding component specification is captured in left component matching
	"\\s*(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "\\s*)+" + // tolerate presence or absence of separating spaces
	"(" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\s*(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "\\s*)+" +
	")+" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before closing brace
	"\\" + RIGHT_BRACE

// 3rd order combinations of parenthesized or braced combinations, including combinations of combinations as components
const BRACED_3RD_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMPONENTS =
// Optional leading component identifier
COMPONENT_IDENTIFIER + "?" +
	"(\\" + LEFT_BRACE +
	// Testing of potential excessive words preceding component specification is captured in left component matching
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_2ND_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMPONENTS +
	")" + OPTIONAL_WORDS_WITH_PARENTHESES + // random words following combination element and logical operator
	"\\s*)+" + // tolerate presence or absence of separating spaces
	"(" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_2ND_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMPONENTS +
	")" +
	"\\s*)+" + // tolerate presence or absence of separating spaces
	")+" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before closing brace
	"\\" + RIGHT_BRACE + ")"

// 4th order combinations of combinations of combinations of parenthesized or braced combinations, including combinations of combinations as components
const BRACED_4TH_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS =
// Optional leading component identifier
COMPONENT_IDENTIFIER + "?" +
	"(\\" + LEFT_BRACE +
	// Testing of potential excessive words preceding component specification is captured in left component matching
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_3RD_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMPONENTS +
	")" + OPTIONAL_WORDS_WITH_PARENTHESES + // random words following combination element and logical operator
	"\\s*)+" + // tolerate presence or absence of separating spaces
	"(" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_3RD_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMPONENTS +
	")" +
	"\\s*)+" + // tolerate presence or absence of separating spaces
	")+" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before closing brace
	"\\" + RIGHT_BRACE + ")"

// 5th order combinations of combinations of combinations of parenthesized or braced combinations, including combinations of combinations as components
const BRACED_5TH_ORDER_COMBINATIONS =
// Optional leading component identifier
COMPONENT_IDENTIFIER + "?" +
	"(\\" + LEFT_BRACE +
	// Testing of potential excessive words preceding component specification is captured in left component matching
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_4TH_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS +
	")" + OPTIONAL_WORDS_WITH_PARENTHESES + // random words following combination element and logical operator
	"\\s*)+" + // tolerate presence or absence of separating spaces
	"(" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_4TH_ORDER_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS_OF_COMBINATIONS +
	")" +
	"\\s*)+" + // tolerate presence or absence of separating spaces
	")+" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before closing brace
	"\\" + RIGHT_BRACE + ")"

// 6th order combinations of combinations of combinations of parenthesized or braced combinations, including combinations of combinations as components
const BRACED_6TH_ORDER_COMBINATIONS =
// Optional leading component identifier
//COMPONENT_IDENTIFIER + "?" +
"(\\" + LEFT_BRACE +
	// Testing of potential excessive words preceding component specification is captured in left component matching
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_5TH_ORDER_COMBINATIONS +
	")" + OPTIONAL_WORDS_WITH_PARENTHESES + // random words following combination element and logical operator
	"\\s*)+" + // tolerate presence or absence of separating spaces
	"(" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_5TH_ORDER_COMBINATIONS +
	")" +
	"\\s*)+" + // tolerate presence or absence of separating spaces
	")+" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before closing brace
	"\\" + RIGHT_BRACE + ")"

// 7th order combinations of combinations of combinations of parenthesized or braced combinations, including combinations of combinations as components
const BRACED_7TH_ORDER_COMBINATIONS =
// Optional leading component identifier
COMPONENT_IDENTIFIER + "?" +
	"(\\" + LEFT_BRACE +
	// Testing of potential excessive words preceding component specification is captured in left component matching
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_6TH_ORDER_COMBINATIONS +
	")" + OPTIONAL_WORDS_WITH_PARENTHESES + // random words following combination element and logical operator
	"\\s*)+" + // tolerate presence or absence of separating spaces
	"(" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_6TH_ORDER_COMBINATIONS +
	")" +
	"\\s*)+" + // tolerate presence or absence of separating spaces
	")+" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before closing brace
	"\\" + RIGHT_BRACE + ")"

// 8th order combinations of combinations of combinations of parenthesized or braced combinations, including combinations of combinations as components
const BRACED_8TH_ORDER_COMBINATIONS =
// Optional leading component identifier
//COMPONENT_IDENTIFIER + "?" +
"(\\" + LEFT_BRACE +
	// Testing of potential excessive words preceding component specification is captured in left component matching
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_7TH_ORDER_COMBINATIONS +
	")" + OPTIONAL_WORDS_WITH_PARENTHESES + // random words following combination element and logical operator
	"\\s*)+" + // tolerate presence or absence of separating spaces
	"(" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\" + LEFT_BRACKET + LOGICAL_OPERATORS + "\\" + RIGHT_BRACKET +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before/after logical operator
	"\\s*(" + "(" + PARENTHESIZED_OR_NON_PARENTHESIZED_COMBINATION_OF_COMPONENTS + "|" +
	BRACED_7TH_ORDER_COMBINATIONS +
	")" +
	"\\s*)+" + // tolerate presence or absence of separating spaces
	")+" +
	OPTIONAL_WORDS_WITH_PARENTHESES + // random words before closing brace
	"\\" + RIGHT_BRACE + ")"

/*
General note on regular expressions for nested combinations/component pairs:
Currently, the nesting level is limited for performance reasons. If higher
nesting levels are required, adjust the BRACED_6TH_ORDER_COMBINATIONS
reference to higher order (e.g., BRACED_8TH_ORDER_COMBINATIONS).
If adjusting the nesting level, ensure that the used level suppresses
the need for a component identifier (COMPONENT_IDENTIFIER) - as shown for
BRACED_6TH_ORDER_COMBINATIONS.
*/

// Expression to filter combinations of combinations to distinguish from component-level component-level nesting
// Termination prevents partial matches (which leads to unintended omission of parts)
// (e.g., embedded combinations such as 'Cac{Cac{A(dlsjg) I(lkdjsg)} ...' in nested statements such as
// 'Bdir{A(dlsjg) I(klsdjg) Cac{Cac{A(dlsjg) I(lkdjsg)} [AND] Cac{A(dslgkj) I(jdlskgj)}}')
const NESTED_COMBINATIONS_TERMINATED =
// Component combinations need to lead with component identifier (and potential suffix and annotation)
"^" +
	COMPONENT_HEADER_SYNTAX +
	// Ensure the tested statement only contains combinations, but no leading individual component (i.e., combination embedded in nested statement)
	BRACED_6TH_ORDER_COMBINATIONS +
	"$" // Ensure immediate termination of combination with additional trailing components (which would imply nested statement with embedded combination)

// Combination of combinations to represent multi-level nesting (does not require termination, i.e., could be embedded)
// Example: 'Cac{ Cac{ I(leftact) Bdir(object1) } [XOR] Cac{ I(rightact) Bdir(object2) }}')
// Not to be used if filtering against component-level nesting (since combinations may be contained as subset)
const NESTED_COMBINATIONS =
// Component combinations need to lead with component identifier (and potential suffix and annotation), e.g., 'Cac1[annotation]{ ... }')
COMPONENT_HEADER_SYNTAX +
	BRACED_6TH_ORDER_COMBINATIONS

// Component combination pairs to be extrapolated into separate statements complemented with basic components (may contain leading annotation,
// but no leading component identifier)
// Example: '{ Cac{ I(leftact) Bdir(object1) } [XOR] Cac{ I(rightact) Bdir(object2) }}')
const COMPONENT_PAIR_COMBINATIONS =
// Component pairs can contain statement-level annotations (e.g., '[boundaryStmt]{ ... }'), but not component identifier (which would make it component combination)
COMPONENT_ANNOTATION_SYNTAX +
	BRACED_6TH_ORDER_COMBINATIONS

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
