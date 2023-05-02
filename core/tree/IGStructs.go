package tree

import (
	"errors"
	"strconv"
	"strings"
)

// Internal column header for components with nested statement references
const REF_SUFFIX = "-Ref"

// Column header in output for reference column
const NAME_REF_SUFFIX = " Reference"

// Brackets used for logical operators (and annotations)
const LEFT_BRACKET = "["
const RIGHT_BRACKET = "]"

// Syntax for properties in component identifiers (e.g., A,p)
const PROPERTY_SYNTAX_SUFFIX = ",p"

// Suffix for component-specific annotation column header
const ANNOTATION = " (Annotation)"

// Column header for statement annotations
const STATEMENT_ANNOTATION = "Statement Annotation"

const (
	ATTRIBUTES                 = "A"
	NAME_ATTRIBUTES            = "Attributes"
	ATTRIBUTES_ANNOTATION      = ATTRIBUTES + ANNOTATION
	NAME_ATTRIBUTES_ANNOTATION = ATTRIBUTES_ANNOTATION

	ATTRIBUTES_PROPERTY                 = ATTRIBUTES + PROPERTY_SYNTAX_SUFFIX
	NAME_ATTRIBUTES_PROPERTY            = "Attributes Property"
	ATTRIBUTES_PROPERTY_REFERENCE       = ATTRIBUTES_PROPERTY + REF_SUFFIX
	NAME_ATTRIBUTES_PROPERTY_REFERENCE  = NAME_ATTRIBUTES_PROPERTY + NAME_REF_SUFFIX
	ATTRIBUTES_PROPERTY_ANNOTATION      = ATTRIBUTES_PROPERTY + ANNOTATION
	NAME_ATTRIBUTES_PROPERTY_ANNOTATION = ATTRIBUTES_PROPERTY_ANNOTATION

	DEONTIC                 = "D"
	NAME_DEONTIC            = "Deontic"
	DEONTIC_ANNOTATION      = DEONTIC + ANNOTATION
	NAME_DEONTIC_ANNOTATION = DEONTIC_ANNOTATION

	AIM                 = "I"
	NAME_AIM            = "Aim"
	AIM_ANNOTATION      = AIM + ANNOTATION
	NAME_AIM_ANNOTATION = AIM_ANNOTATION

	DIRECT_OBJECT                 = "Bdir"
	NAME_DIRECT_OBJECT            = "Direct Object"
	DIRECT_OBJECT_REFERENCE       = DIRECT_OBJECT + REF_SUFFIX
	NAME_DIRECT_OBJECT_REFERENCE  = NAME_DIRECT_OBJECT + NAME_REF_SUFFIX
	DIRECT_OBJECT_ANNOTATION      = DIRECT_OBJECT + ANNOTATION
	NAME_DIRECT_OBJECT_ANNOTATION = DIRECT_OBJECT_ANNOTATION

	DIRECT_OBJECT_PROPERTY                 = DIRECT_OBJECT + PROPERTY_SYNTAX_SUFFIX
	NAME_DIRECT_OBJECT_PROPERTY            = "Direct Object Property"
	DIRECT_OBJECT_PROPERTY_REFERENCE       = DIRECT_OBJECT_PROPERTY + REF_SUFFIX
	NAME_DIRECT_OBJECT_PROPERTY_REFERENCE  = NAME_DIRECT_OBJECT_PROPERTY + NAME_REF_SUFFIX
	DIRECT_OBJECT_PROPERTY_ANNOTATION      = DIRECT_OBJECT_PROPERTY + ANNOTATION
	NAME_DIRECT_OBJECT_PROPERTY_ANNOTATION = DIRECT_OBJECT_PROPERTY_ANNOTATION

	INDIRECT_OBJECT                 = "Bind"
	NAME_INDIRECT_OBJECT            = "Indirect Object"
	INDIRECT_OBJECT_REFERENCE       = INDIRECT_OBJECT + REF_SUFFIX
	NAME_INDIRECT_OBJECT_REFERENCE  = NAME_INDIRECT_OBJECT + NAME_REF_SUFFIX
	INDIRECT_OBJECT_ANNOTATION      = INDIRECT_OBJECT + ANNOTATION
	NAME_INDIRECT_OBJECT_ANNOTATION = INDIRECT_OBJECT_ANNOTATION

	INDIRECT_OBJECT_PROPERTY                 = INDIRECT_OBJECT + PROPERTY_SYNTAX_SUFFIX
	NAME_INDIRECT_OBJECT_PROPERTY            = "Indirect Object Property"
	INDIRECT_OBJECT_PROPERTY_REFERENCE       = INDIRECT_OBJECT_PROPERTY + REF_SUFFIX
	NAME_INDIRECT_OBJECT_PROPERTY_REFERENCE  = NAME_INDIRECT_OBJECT_PROPERTY + NAME_REF_SUFFIX
	INDIRECT_OBJECT_PROPERTY_ANNOTATION      = INDIRECT_OBJECT_PROPERTY + ANNOTATION
	NAME_INDIRECT_OBJECT_PROPERTY_ANNOTATION = INDIRECT_OBJECT_PROPERTY_ANNOTATION

	ACTIVATION_CONDITION                 = "Cac"
	NAME_ACTIVATION_CONDITION            = "Activation Condition"
	ACTIVATION_CONDITION_REFERENCE       = ACTIVATION_CONDITION + REF_SUFFIX
	NAME_ACTIVATION_CONDITION_REFERENCE  = NAME_ACTIVATION_CONDITION + NAME_REF_SUFFIX
	ACTIVATION_CONDITION_ANNOTATION      = ACTIVATION_CONDITION + ANNOTATION
	NAME_ACTIVATION_CONDITION_ANNOTATION = ACTIVATION_CONDITION_ANNOTATION

	EXECUTION_CONSTRAINT                 = "Cex"
	NAME_EXECUTION_CONSTRAINT            = "Execution Constraint"
	EXECUTION_CONSTRAINT_REFERENCE       = EXECUTION_CONSTRAINT + REF_SUFFIX
	NAME_EXECUTION_CONSTRAINT_REFERENCE  = NAME_EXECUTION_CONSTRAINT + NAME_REF_SUFFIX
	EXECUTION_CONSTRAINT_ANNOTATION      = EXECUTION_CONSTRAINT + ANNOTATION
	NAME_EXECUTION_CONSTRAINT_ANNOTATION = EXECUTION_CONSTRAINT_ANNOTATION

	CONSTITUTED_ENTITY                 = "E"
	NAME_CONSTITUTED_ENTITY            = "Constituted Entity"
	CONSTITUTED_ENTITY_ANNOTATION      = CONSTITUTED_ENTITY + ANNOTATION
	NAME_CONSTITUTED_ENTITY_ANNOTATION = CONSTITUTED_ENTITY_ANNOTATION

	CONSTITUTED_ENTITY_PROPERTY                 = CONSTITUTED_ENTITY + PROPERTY_SYNTAX_SUFFIX
	NAME_CONSTITUTED_ENTITY_PROPERTY            = "Constituted Entity Property"
	CONSTITUTED_ENTITY_PROPERTY_REFERENCE       = CONSTITUTED_ENTITY_PROPERTY + REF_SUFFIX
	NAME_CONSTITUTED_ENTITY_PROPERTY_REFERENCE  = NAME_CONSTITUTED_ENTITY_PROPERTY + NAME_REF_SUFFIX
	CONSTITUTED_ENTITY_PROPERTY_ANNOTATION      = CONSTITUTED_ENTITY_PROPERTY + ANNOTATION
	NAME_CONSTITUTED_ENTITY_PROPERTY_ANNOTATION = CONSTITUTED_ENTITY_PROPERTY_ANNOTATION

	MODAL                 = "M"
	NAME_MODAL            = "Modal"
	MODAL_ANNOTATION      = MODAL + ANNOTATION
	NAME_MODAL_ANNOTATION = MODAL_ANNOTATION

	CONSTITUTIVE_FUNCTION                 = "F"
	NAME_CONSTITUTIVE_FUNCTION            = "Constitutive Function"
	CONSTITUTIVE_FUNCTION_ANNOTATION      = CONSTITUTIVE_FUNCTION + ANNOTATION
	NAME_CONSTITUTIVE_FUNCTION_ANNOTATION = CONSTITUTIVE_FUNCTION_ANNOTATION

	CONSTITUTING_PROPERTIES                 = "P"
	NAME_CONSTITUTING_PROPERTIES            = "Constituting Properties"
	CONSTITUTING_PROPERTIES_REFERENCE       = CONSTITUTING_PROPERTIES + REF_SUFFIX
	NAME_CONSTITUTING_PROPERTIES_REFERENCE  = NAME_CONSTITUTING_PROPERTIES + NAME_REF_SUFFIX
	CONSTITUTING_PROPERTIES_ANNOTATION      = CONSTITUTING_PROPERTIES + ANNOTATION
	NAME_CONSTITUTING_PROPERTIES_ANNOTATION = CONSTITUTING_PROPERTIES_ANNOTATION

	CONSTITUTING_PROPERTIES_PROPERTY                 = CONSTITUTING_PROPERTIES + PROPERTY_SYNTAX_SUFFIX
	NAME_CONSTITUTING_PROPERTIES_PROPERTY            = "Constituting Properties Properties"
	CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE       = CONSTITUTING_PROPERTIES_PROPERTY + REF_SUFFIX
	NAME_CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE  = NAME_CONSTITUTING_PROPERTIES_PROPERTY + NAME_REF_SUFFIX
	CONSTITUTING_PROPERTIES_PROPERTY_ANNOTATION      = CONSTITUTING_PROPERTIES_PROPERTY + ANNOTATION
	NAME_CONSTITUTING_PROPERTIES_PROPERTY_ANNOTATION = CONSTITUTING_PROPERTIES_PROPERTY_ANNOTATION

	OR_ELSE                = "O"
	NAME_OR_ELSE           = "Or Else"
	OR_ELSE_REFERENCE      = OR_ELSE + REF_SUFFIX
	NAME_OR_ELSE_REFERENCE = NAME_OR_ELSE + NAME_REF_SUFFIX

	// Made sAND a variable
	AND = "AND"
	OR  = "OR"
	XOR = "XOR"
	NOT = "NOT"

	AND_BRACKETS = LEFT_BRACKET + AND + RIGHT_BRACKET
	OR_BRACKETS  = LEFT_BRACKET + OR + RIGHT_BRACKET
	XOR_BRACKETS = LEFT_BRACKET + XOR + RIGHT_BRACKET
	NOT_BRACKETS = LEFT_BRACKET + NOT + RIGHT_BRACKET

	// Synthetic AND between components of the same type (e.g., 'I(first) I(second)' linked as
	// 'I((first [SAND_BETWEEN_COMPONENTS] second))')
	SAND_BETWEEN_COMPONENTS          = "bAND"
	SAND_BETWEEN_COMPONENTS_BRACKETS = LEFT_BRACKET + SAND_BETWEEN_COMPONENTS + RIGHT_BRACKET
	// Synthetic AND between components of the same type (e.g., 'I((first [AND] second) (third [AND] fourth))' linked as
	// 'I((first [AND] second) [SAND_WITHIN_COMPONENTS] (third [AND] fourth))'
	SAND_WITHIN_COMPONENTS          = "wAND"
	SAND_WITHIN_COMPONENTS_BRACKETS = LEFT_BRACKET + SAND_WITHIN_COMPONENTS + RIGHT_BRACKET

	PARSING_MODE_LEFT               = "PARSING_LEFT"
	PARSING_MODE_RIGHT              = "PARSING_RIGHT"
	PARSING_MODE_OUTSIDE_EXPRESSION = "PARSING_OUTSIDE"
)

/*
Indicates whether a given symbol is a valid IG Component symbol
*/
func ValidIGComponentSymbol(symbol string) bool {
	res, _ := StringInSlice(symbol, IGComponentSymbols)
	return res
}

/*
Indicates whether a given name is a valid IG Component name
*/
func validIGComponentName(name string) bool {
	res, _ := StringInSlice(name, IGComponentNames)
	return res
}

/*
IG 2.0 Component Symbols
*/
var IGComponentSymbols = []string{

	// General annotations on statement level
	STATEMENT_ANNOTATION,

	// Actual components
	ATTRIBUTES,
	ATTRIBUTES_ANNOTATION,

	ATTRIBUTES_PROPERTY,
	ATTRIBUTES_PROPERTY_REFERENCE,
	ATTRIBUTES_PROPERTY_ANNOTATION,

	DEONTIC,
	DEONTIC_ANNOTATION,

	AIM,
	AIM_ANNOTATION,

	DIRECT_OBJECT,
	DIRECT_OBJECT_REFERENCE,
	DIRECT_OBJECT_ANNOTATION,

	DIRECT_OBJECT_PROPERTY,
	DIRECT_OBJECT_PROPERTY_REFERENCE,
	DIRECT_OBJECT_PROPERTY_ANNOTATION,

	INDIRECT_OBJECT,
	INDIRECT_OBJECT_REFERENCE,
	INDIRECT_OBJECT_ANNOTATION,

	INDIRECT_OBJECT_PROPERTY,
	INDIRECT_OBJECT_PROPERTY_REFERENCE,
	INDIRECT_OBJECT_PROPERTY_ANNOTATION,

	ACTIVATION_CONDITION,
	ACTIVATION_CONDITION_REFERENCE,
	ACTIVATION_CONDITION_ANNOTATION,

	EXECUTION_CONSTRAINT,
	EXECUTION_CONSTRAINT_REFERENCE,
	EXECUTION_CONSTRAINT_ANNOTATION,

	CONSTITUTED_ENTITY,
	CONSTITUTED_ENTITY_ANNOTATION,

	CONSTITUTED_ENTITY_PROPERTY,
	CONSTITUTED_ENTITY_PROPERTY_REFERENCE,
	CONSTITUTED_ENTITY_PROPERTY_ANNOTATION,

	MODAL,
	MODAL_ANNOTATION,

	CONSTITUTIVE_FUNCTION,
	CONSTITUTIVE_FUNCTION_ANNOTATION,

	CONSTITUTING_PROPERTIES,
	CONSTITUTING_PROPERTIES_REFERENCE,
	CONSTITUTING_PROPERTIES_ANNOTATION,

	CONSTITUTING_PROPERTIES_PROPERTY,
	CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE,
	CONSTITUTING_PROPERTIES_PROPERTY_ANNOTATION,

	OR_ELSE,
	OR_ELSE_REFERENCE,
}

/*
IG 2.0 Component Symbols
*/
var IGComponentNames = []string{

	// Statement annotations
	STATEMENT_ANNOTATION,

	// Component names
	NAME_ATTRIBUTES,
	NAME_ATTRIBUTES_ANNOTATION,

	NAME_ATTRIBUTES_PROPERTY,
	NAME_ATTRIBUTES_PROPERTY_REFERENCE,
	NAME_ATTRIBUTES_PROPERTY_ANNOTATION,

	NAME_DEONTIC,
	NAME_DEONTIC_ANNOTATION,

	NAME_AIM,
	NAME_AIM_ANNOTATION,

	NAME_DIRECT_OBJECT,
	NAME_DIRECT_OBJECT_REFERENCE,
	NAME_DIRECT_OBJECT_ANNOTATION,

	NAME_DIRECT_OBJECT_PROPERTY,
	NAME_DIRECT_OBJECT_PROPERTY_REFERENCE,
	NAME_DIRECT_OBJECT_PROPERTY_ANNOTATION,

	NAME_INDIRECT_OBJECT,
	NAME_INDIRECT_OBJECT_REFERENCE,
	NAME_INDIRECT_OBJECT_ANNOTATION,

	NAME_INDIRECT_OBJECT_PROPERTY,
	NAME_INDIRECT_OBJECT_PROPERTY_REFERENCE,
	NAME_INDIRECT_OBJECT_PROPERTY_ANNOTATION,

	NAME_ACTIVATION_CONDITION,
	NAME_ACTIVATION_CONDITION_REFERENCE,
	NAME_ACTIVATION_CONDITION_ANNOTATION,

	NAME_EXECUTION_CONSTRAINT,
	NAME_EXECUTION_CONSTRAINT_REFERENCE,
	NAME_EXECUTION_CONSTRAINT_ANNOTATION,

	NAME_CONSTITUTED_ENTITY,
	NAME_CONSTITUTED_ENTITY_ANNOTATION,

	NAME_CONSTITUTED_ENTITY_PROPERTY,
	NAME_CONSTITUTED_ENTITY_PROPERTY_REFERENCE,
	NAME_CONSTITUTED_ENTITY_PROPERTY_ANNOTATION,

	NAME_MODAL,
	NAME_MODAL_ANNOTATION,

	NAME_CONSTITUTIVE_FUNCTION,
	NAME_CONSTITUTIVE_FUNCTION_ANNOTATION,

	NAME_CONSTITUTING_PROPERTIES,
	NAME_CONSTITUTING_PROPERTIES_REFERENCE,
	NAME_CONSTITUTING_PROPERTIES_ANNOTATION,

	NAME_CONSTITUTING_PROPERTIES_PROPERTY,
	NAME_CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE,
	NAME_CONSTITUTING_PROPERTIES_PROPERTY_ANNOTATION,

	NAME_OR_ELSE,
	NAME_OR_ELSE_REFERENCE,
}

/*
Map holding mapping from IG 2.0 component symbols to proper component names
*/
var IGComponentSymbolNameMap = map[string]string{

	// Statement-level annotations
	STATEMENT_ANNOTATION: STATEMENT_ANNOTATION,

	// Actual components
	ATTRIBUTES:            NAME_ATTRIBUTES,
	ATTRIBUTES_ANNOTATION: NAME_ATTRIBUTES_ANNOTATION,

	ATTRIBUTES_PROPERTY:            NAME_ATTRIBUTES_PROPERTY,
	ATTRIBUTES_PROPERTY_REFERENCE:  NAME_ATTRIBUTES_PROPERTY_REFERENCE,
	ATTRIBUTES_PROPERTY_ANNOTATION: NAME_ATTRIBUTES_PROPERTY_ANNOTATION,

	DEONTIC:            NAME_DEONTIC,
	DEONTIC_ANNOTATION: NAME_DEONTIC_ANNOTATION,

	AIM:            NAME_AIM,
	AIM_ANNOTATION: NAME_AIM_ANNOTATION,

	DIRECT_OBJECT:            NAME_DIRECT_OBJECT,
	DIRECT_OBJECT_REFERENCE:  NAME_DIRECT_OBJECT_REFERENCE,
	DIRECT_OBJECT_ANNOTATION: NAME_DIRECT_OBJECT_ANNOTATION,

	DIRECT_OBJECT_PROPERTY:            NAME_DIRECT_OBJECT_PROPERTY,
	DIRECT_OBJECT_PROPERTY_REFERENCE:  NAME_DIRECT_OBJECT_PROPERTY_REFERENCE,
	DIRECT_OBJECT_PROPERTY_ANNOTATION: NAME_DIRECT_OBJECT_PROPERTY_ANNOTATION,

	INDIRECT_OBJECT:            NAME_INDIRECT_OBJECT,
	INDIRECT_OBJECT_REFERENCE:  NAME_INDIRECT_OBJECT_REFERENCE,
	INDIRECT_OBJECT_ANNOTATION: NAME_INDIRECT_OBJECT_ANNOTATION,

	INDIRECT_OBJECT_PROPERTY:            NAME_INDIRECT_OBJECT_PROPERTY,
	INDIRECT_OBJECT_PROPERTY_REFERENCE:  NAME_INDIRECT_OBJECT_PROPERTY_REFERENCE,
	INDIRECT_OBJECT_PROPERTY_ANNOTATION: NAME_INDIRECT_OBJECT_PROPERTY_ANNOTATION,

	ACTIVATION_CONDITION:            NAME_ACTIVATION_CONDITION,
	ACTIVATION_CONDITION_REFERENCE:  NAME_ACTIVATION_CONDITION_REFERENCE,
	ACTIVATION_CONDITION_ANNOTATION: NAME_ACTIVATION_CONDITION_ANNOTATION,

	EXECUTION_CONSTRAINT:            NAME_EXECUTION_CONSTRAINT,
	EXECUTION_CONSTRAINT_REFERENCE:  NAME_EXECUTION_CONSTRAINT_REFERENCE,
	EXECUTION_CONSTRAINT_ANNOTATION: NAME_EXECUTION_CONSTRAINT_ANNOTATION,

	CONSTITUTED_ENTITY:            NAME_CONSTITUTED_ENTITY,
	CONSTITUTED_ENTITY_ANNOTATION: NAME_CONSTITUTED_ENTITY_ANNOTATION,

	CONSTITUTED_ENTITY_PROPERTY:            NAME_CONSTITUTED_ENTITY_PROPERTY,
	CONSTITUTED_ENTITY_PROPERTY_REFERENCE:  NAME_CONSTITUTED_ENTITY_PROPERTY_REFERENCE,
	CONSTITUTED_ENTITY_PROPERTY_ANNOTATION: NAME_CONSTITUTED_ENTITY_PROPERTY_ANNOTATION,

	MODAL:            NAME_MODAL,
	MODAL_ANNOTATION: NAME_MODAL_ANNOTATION,

	CONSTITUTIVE_FUNCTION:            NAME_CONSTITUTIVE_FUNCTION,
	CONSTITUTIVE_FUNCTION_ANNOTATION: NAME_CONSTITUTIVE_FUNCTION_ANNOTATION,

	CONSTITUTING_PROPERTIES:            NAME_CONSTITUTING_PROPERTIES,
	CONSTITUTING_PROPERTIES_REFERENCE:  NAME_CONSTITUTING_PROPERTIES_REFERENCE,
	CONSTITUTING_PROPERTIES_ANNOTATION: NAME_CONSTITUTING_PROPERTIES_ANNOTATION,

	CONSTITUTING_PROPERTIES_PROPERTY:            NAME_CONSTITUTING_PROPERTIES_PROPERTY,
	CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE:  NAME_CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE,
	CONSTITUTING_PROPERTIES_PROPERTY_ANNOTATION: NAME_CONSTITUTING_PROPERTIES_PROPERTY_ANNOTATION,

	OR_ELSE:           NAME_OR_ELSE,
	OR_ELSE_REFERENCE: NAME_OR_ELSE_REFERENCE,
}

type igLogicalOperator struct {
	LogicalOperatorName string
}

/*
Checks whether operator value is valid (i.e., a valid logical operator symbol).
*/
func (o *igLogicalOperator) valid() bool {
	res, _ := StringInSlice(o.LogicalOperatorName, IGLogicalOperators)
	return res
}

func (o igLogicalOperator) String() string {
	return o.LogicalOperatorName
}

/*
Valid logical operators in IG 2.0
*/
var IGLogicalOperators = []string{
	AND,
	OR,
	XOR,
	NOT,
}

/*
Signals boundary values for detected combinations
*/
type Boundaries struct {
	// Left boundary
	Left int
	// Operator position
	Operator int
	// Operator value (if combination)
	OperatorVal string
	// Right boundary
	Right int
	// Indicates whether element is combination
	Complete bool
	// Indicates whether element is shared (i.e., belongs to some combination)
	Shared bool
}

func (b *Boundaries) String() string {
	return "Boundaries{\n" +
		"  Left: " + strconv.Itoa(b.Left) + "\n" +
		"  Op Pos: " + strconv.Itoa(b.Operator) + "\n" +
		"  Operator: " + b.OperatorVal + "\n" +
		"  Right: " + strconv.Itoa(b.Right) + "\n" +
		"  Complete: " + strconv.FormatBool(b.Complete) + "\n}"
}

// Signals invalid component combinations on a given parsing level (generally non-AND components)
const PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS = "INVALID_LOGICAL_OPERATOR_COMBINATIONS"

// Signals the detection of a logical operator outside a combination (i.e., no left or right side)
const PARSING_ERROR_LOGICAL_OPERATOR_OUTSIDE_COMBINATION = "LOGICAL_OPERATOR_OUTSIDE_COMBINATION"

// Indicates that there were no issues during parsing
const PARSING_NO_ERROR = "NO_ERROR_DURING_PARSING"

// Indicates nested statements ignored during coding (relevant when parsing nested statements) - in contrast to #PARSING_ERROR_IGNORED_ELEMENTS_DURING_NODE_PARSING
const PARSING_ERROR_IGNORED_NESTED_ELEMENTS = "IGNORED_ELEMENTS_NESTED_STATEMENT_PARSING"

// Signals that no combinations were found in input
const PARSING_NO_COMBINATIONS = "NO_COMBINATIONS_IN_INPUT"

// Signals invalid combination expression (e.g., missing left, right, or operator)
const PARSING_INVALID_COMBINATION = "INVALID_COMBINATION_IN_INPUT"

// Signals empty leaf value during parsing
const PARSING_ERROR_EMPTY_LEAF = "EMPTY_LEAF_VALUE"

// Indicates that parentheses are umbalanced during parsing
const PARSING_ERROR_IMBALANCED_PARENTHESES = "IMBALANCED_PARENTHESES"

// Indicates that component was not found during preprocessing
const PARSING_ERROR_COMPONENT_NOT_FOUND = "COMPONENT_NOT_FOUND"

// Indicates that multiple component symbols have been found in component specification -- points to incorrect coding
const PARSING_ERROR_MULTIPLE_COMPONENTS_FOUND = "MULTIPLE_COMPONENTS_FOUND"

// Indicates ignored elements during parsing (applies when parsing into tree) - in contrast to #PARSING_ERROR_IGNORED_NESTED_ELEMENTS
const PARSING_ERROR_IGNORED_ELEMENTS_DURING_NODE_PARSING = "IGNORED_ELEMENTS_NODE_PARSING"

// Indicates problems when generating logical operator references
const PARSING_ERROR_LOGICAL_EXPRESSION_GENERATION = "LOGICAL_EXPRESSION_GENERATION"

// Write error
const PARSING_ERROR_WRITE = "WRITE_ERROR"

// Invalid parentheses/braces combinations
const PARSING_ERROR_INVALID_PARENTHESES_COMBINATION = "INVALID_PARENTHESES_COMBINATIONS"

// Error during regex compilation
const PARSING_ERROR_PATTERN_EXTRACTION = "PATTERN_EXTRACTION_ERROR"

// Detecting combinations of nested statements with varying component references
// (e.g., {Cac{stmt1} [AND] Cex{stmt2}}, but should be{Cac{stmt1} [AND] Cac{stmt2}})
const PARSING_ERROR_INVALID_TYPES_IN_NESTED_STATEMENT_COMBINATION = "INVALID_TYPE_COMBINATIONS_IN_NESTED_STATEMENT_COMBINATIONS"

// Indicates that operations was imposed on nil element
const PARSING_ERROR_NIL_ELEMENT = "INVALID_PARSING_OF_NIL_ELEMENT"

// Indicates an unexpected parsing error - can be diverse
const PARSING_ERROR_UNEXPECTED_ERROR = "UNEXPECTED_ERROR"

// Indicates nesting on invalid component (no component-level nesting)
const PARSING_ERROR_NESTING_ON_UNSUPPORTED_COMPONENT = "NESTING_ON_NON-NESTED_COMPONENT"

// Indicates missing separator value (for output generation)
const PARSING_ERROR_MISSING_SEPARATOR_VALUE = "MISSING_SEPARATOR"

// Indicates invalid combination of component types (e.g., Cac and Cex in combined node)
const PARSING_ERROR_INVALID_COMPONENT_TYPE_COMBINATION = "INVALID_COMPONENT_TYPE_COMBINATION"

// Indicates an embedded node error (passed via NodeError) as part of a ParsingError
const PARSING_ERROR_EMBEDDED_NODE_ERROR = "EMBEDDED_NODE_ERROR"

// Indicates invalid output type (should be one of TabularOutputGeneratorConfig #OUTPUT_TYPE_CSV or #OUTPUT_TYPE_GOOGLE_SHEETS or #OUTPUT_TYPE_NONE)
const PARSING_ERROR_INVALID_OUTPUT_TYPE = "INVALID_OUTPUT_TYPE"

// Indicates unexpected number of nodes in array
const PARSING_ERROR_TOO_MANY_NODES = "TOO_MANY_NODES"

// Indicates invalid type (i.e., no node or statement) during visual output generation
const PARSING_ERROR_INVALID_TYPE_VISUAL_OUTPUT = "INVALID TYPE FOR VISUAL OUTPUT"

/*
Error type signaling errors during statement parsing
*/
type ParsingError struct {
	ErrorCode            string
	ErrorMessage         string
	ErrorIgnoredElements []string
}

func (e *ParsingError) Error() error {
	return errors.New("Parsing Error " + e.ErrorCode + ": " + e.ErrorMessage +
		" (Ignored elements: " + strconv.Itoa(len(e.ErrorIgnoredElements)) + ")")
}

/*
Error type signaling errors during Node tree operations
*/
type NodeError struct {
	ErrorCode            string
	ErrorMessage         string
	ErrorIgnoredElements []string
}

func (e *NodeError) Error() error {
	return errors.New("Node Error " + e.ErrorCode + ": " + e.ErrorMessage +
		" (Ignored elements: " + strconv.Itoa(len(e.ErrorIgnoredElements)) + ")")
}

// Standard error type if no error was found
const TREE_NO_ERROR = "NO_ERROR"
const TREE_INVALID_NODE_ADDITION = "INVALID_NODE_ADDITION"

// Indicates that node is linked to itself as parent
const TREE_INVALID_NODE_SELF_LINKAGE = "INVALID_NODE_LINKAGE_TO_SELF"
const TREE_INVALID_NODE_REMOVAL = "INVALID_NODE_REMOVAL"
const TREE_INVALID_TREE = "TREE_STRUCTURE_INVALID"

// Indicates combination of different component types under shared component node.
const TREE_INVALID_COMPONENT_COMBINATIONS = "TREE_INVALID_COMPONENT_COMBINATIONS"
const TREE_INPUT_VALIDATION = "INPUT_VALIDATION"

// Indicates that a NodeError embeds a ParsingError
const TREE_ERROR_EMBEDDED_PARSING_ERROR = "TREE_EMBEDDED_PARSING_ERROR"

/*
StateComplexity struct containing number of individual options per component,
as well as calculated state complexity per component, as well as in total
*/
type StateComplexity struct {
	AttributesOptions                       int
	AttributesComplexity                    int
	AttributesPropertySimpleOptions         int
	AttributesPropertySimpleComplexity      int
	AttributesPropertyComplexOptions        int
	AttributesPropertyComplexComplexity     int
	DeonticOptions                          int
	DeonticComplexity                       int
	AimOptions                              int
	AimComplexity                           int
	DirectObjectSimpleOptions               int
	DirectObjectSimpleComplexity            int
	DirectObjectComplexOptions              int
	DirectObjectComplexComplexity           int
	DirectObjectPropertySimpleOptions       int
	DirectObjectPropertySimpleComplexity    int
	DirectObjectPropertyComplexOptions      int
	DirectObjectPropertyComplexComplexity   int
	IndirectObjectSimpleOptions             int
	IndirectObjectSimpleComplexity          int
	IndirectObjectComplexOptions            int
	IndirectObjectComplexComplexity         int
	IndirectObjectPropertySimpleOptions     int
	IndirectObjectPropertySimpleComplexity  int
	IndirectObjectPropertyComplexOptions    int
	IndirectObjectPropertyComplexComplexity int

	ConstitutedEntityOptions                          int
	ConstitutedEntityComplexity                       int
	ConstitutedEntityPropertySimpleOptions            int
	ConstitutedEntityPropertySimpleComplexity         int
	ConstitutedEntityPropertyComplexOptions           int
	ConstitutedEntityPropertyComplexComplexity        int
	ModalOptions                                      int
	ModalComplexity                                   int
	ConstitutiveFunctionOptions                       int
	ConstitutiveFunctionComplexity                    int
	ConstitutingPropertiesSimpleOptions               int
	ConstitutingPropertiesSimpleComplexity            int
	ConstitutingPropertiesComplexOptions              int
	ConstitutingPropertiesComplexComplexity           int
	ConstitutingPropertiesPropertiesSimpleOptions     int
	ConstitutingPropertiesPropertiesSimpleComplexity  int
	ConstitutingPropertiesPropertiesComplexOptions    int
	ConstitutingPropertiesPropertiesComplexComplexity int

	ExecutionConstraintSimpleOptions     int
	ExecutionConstraintSimpleComplexity  int
	ExecutionConstraintComplexOptions    int
	ExecutionConstraintComplexComplexity int

	ActivationConditionSimpleOptions     int
	ActivationConditionSimpleComplexity  int
	ActivationConditionComplexOptions    int
	ActivationConditionComplexComplexity int

	OrElseComplexity int

	TotalStateComplexity int
}

/*
Produces StateComplexity object as human-readable output.
*/
func (c *StateComplexity) String() string {
	str := strings.Builder{}
	str.WriteString("=== State Complexity ===\n\n")
	str.WriteString("Attributes Options: " + strconv.Itoa(c.AttributesOptions) + "\n")
	str.WriteString("Attributes Complexity: " + strconv.Itoa(c.AttributesComplexity) + "\n\n")

	str.WriteString("Attributes Properties (simple) Options: " + strconv.Itoa(c.AttributesPropertySimpleOptions) + "\n")
	str.WriteString("Attributes Properties (simple) Complexity: " + strconv.Itoa(c.AttributesPropertySimpleComplexity) + "\n")
	str.WriteString("Attributes Properties (complex) Options: " + strconv.Itoa(c.AttributesPropertyComplexOptions) + "\n")
	str.WriteString("Attributes Properties (complex) Complexity: " + strconv.Itoa(c.AttributesPropertyComplexComplexity) + "\n\n")

	str.WriteString("Deontic Options: " + strconv.Itoa(c.DeonticOptions) + "\n")
	str.WriteString("Deontic Complexity: " + strconv.Itoa(c.DeonticComplexity) + "\n\n")

	str.WriteString("Aim Options: " + strconv.Itoa(c.AimOptions) + "\n")
	str.WriteString("Aim Complexity: " + strconv.Itoa(c.AimComplexity) + "\n\n")

	str.WriteString("Direct Object (simple) Options: " + strconv.Itoa(c.DirectObjectSimpleOptions) + "\n")
	str.WriteString("Direct Object (simple) Complexity: " + strconv.Itoa(c.DirectObjectSimpleComplexity) + "\n")
	str.WriteString("Direct Object (complex) Options: " + strconv.Itoa(c.DirectObjectComplexOptions) + "\n")
	str.WriteString("Direct Object (complex) Complexity: " + strconv.Itoa(c.DirectObjectComplexComplexity) + "\n\n")

	str.WriteString("Direct Object Property (simple) Options: " + strconv.Itoa(c.DirectObjectPropertySimpleOptions) + "\n")
	str.WriteString("Direct Object Property (simple) Complexity: " + strconv.Itoa(c.DirectObjectPropertySimpleComplexity) + "\n")
	str.WriteString("Direct Object Property (complex) Options: " + strconv.Itoa(c.DirectObjectPropertyComplexOptions) + "\n")
	str.WriteString("Direct Object Property (complex) Complexity: " + strconv.Itoa(c.DirectObjectPropertyComplexComplexity) + "\n\n")

	str.WriteString("Indirect Object (simple) Options: " + strconv.Itoa(c.IndirectObjectSimpleOptions) + "\n")
	str.WriteString("Indirect Object (simple) Complexity: " + strconv.Itoa(c.IndirectObjectSimpleComplexity) + "\n")
	str.WriteString("Indirect Object (complex) Options: " + strconv.Itoa(c.IndirectObjectComplexOptions) + "\n")
	str.WriteString("Indirect Object (complex) Complexity: " + strconv.Itoa(c.IndirectObjectComplexComplexity) + "\n\n")

	str.WriteString("Indirect Object Property (simple) Options: " + strconv.Itoa(c.IndirectObjectPropertySimpleOptions) + "\n")
	str.WriteString("Indirect Object Property (simple) Complexity: " + strconv.Itoa(c.IndirectObjectPropertySimpleComplexity) + "\n")
	str.WriteString("Indirect Object Property (complex) Options: " + strconv.Itoa(c.IndirectObjectPropertyComplexOptions) + "\n")
	str.WriteString("Indirect Object Property (complex) Complexity: " + strconv.Itoa(c.IndirectObjectPropertyComplexComplexity) + "\n\n")

	str.WriteString("Constituted Entity Options: " + strconv.Itoa(c.ConstitutedEntityOptions) + "\n")
	str.WriteString("Constituted Entity Complexity: " + strconv.Itoa(c.ConstitutedEntityComplexity) + "\n\n")

	str.WriteString("Constituted Entity Property (simple) Options: " + strconv.Itoa(c.ConstitutedEntityPropertySimpleOptions) + "\n")
	str.WriteString("Constituted Entity Property (simple) Complexity: " + strconv.Itoa(c.ConstitutedEntityPropertySimpleComplexity) + "\n")
	str.WriteString("Constituted Entity Property (complex) Options: " + strconv.Itoa(c.ConstitutedEntityPropertyComplexOptions) + "\n")
	str.WriteString("Constituted Entity Property (complex) Complexity: " + strconv.Itoa(c.ConstitutedEntityPropertyComplexComplexity) + "\n\n")

	str.WriteString("Modal Options: " + strconv.Itoa(c.ModalOptions) + "\n")
	str.WriteString("Modal Complexity: " + strconv.Itoa(c.ModalComplexity) + "\n\n")

	str.WriteString("Constitutive Function Options: " + strconv.Itoa(c.ConstitutiveFunctionOptions) + "\n")
	str.WriteString("Constitutive Function Complexity: " + strconv.Itoa(c.ConstitutiveFunctionComplexity) + "\n\n")

	str.WriteString("Constituting Properties (simple) Options: " + strconv.Itoa(c.ConstitutingPropertiesSimpleOptions) + "\n")
	str.WriteString("Constituting Properties (simple) Complexity: " + strconv.Itoa(c.ConstitutingPropertiesSimpleComplexity) + "\n")
	str.WriteString("Constituting Properties (complex) Options: " + strconv.Itoa(c.ConstitutingPropertiesComplexOptions) + "\n")
	str.WriteString("Constituting Properties (complex) Complexity: " + strconv.Itoa(c.ConstitutingPropertiesComplexComplexity) + "\n\n")

	str.WriteString("Constituting Properties Properties (simple) Options: " + strconv.Itoa(c.ConstitutingPropertiesPropertiesSimpleOptions) + "\n")
	str.WriteString("Constituting Properties Properties (simple) Complexity: " + strconv.Itoa(c.ConstitutingPropertiesPropertiesSimpleComplexity) + "\n")
	str.WriteString("Constituting Properties Properties (complex) Options: " + strconv.Itoa(c.ConstitutingPropertiesPropertiesComplexOptions) + "\n")
	str.WriteString("Constituting Properties Properties (complex) Complexity: " + strconv.Itoa(c.ConstitutingPropertiesPropertiesComplexComplexity) + "\n\n")

	str.WriteString("Execution Constraints (simple) Options: " + strconv.Itoa(c.ExecutionConstraintSimpleOptions) + "\n")
	str.WriteString("Execution Constraints (simple) Complexity: " + strconv.Itoa(c.ExecutionConstraintSimpleComplexity) + "\n")
	str.WriteString("Execution Constraints (complex) Options: " + strconv.Itoa(c.ExecutionConstraintComplexOptions) + "\n")
	str.WriteString("Execution Constraints (complex) Complexity: " + strconv.Itoa(c.ExecutionConstraintComplexComplexity) + "\n\n")

	str.WriteString("Activation Conditions (simple) Options: " + strconv.Itoa(c.ActivationConditionSimpleOptions) + "\n")
	str.WriteString("Activation Conditions (simple) Complexity: " + strconv.Itoa(c.ActivationConditionSimpleComplexity) + "\n")
	str.WriteString("Activation Conditions (complex) Options: " + strconv.Itoa(c.ActivationConditionComplexOptions) + "\n")
	str.WriteString("Activation Conditions (complex) Complexity: " + strconv.Itoa(c.ActivationConditionComplexComplexity) + "\n\n")

	str.WriteString("Or Else Complexity: " + strconv.Itoa(c.OrElseComplexity) + "\n\n")

	str.WriteString("Total State Complexity: " + strconv.Itoa(c.TotalStateComplexity) + "\n")

	return str.String()
}

/*
Collapses repeated occurrences of values in a given array (e.g., [AND] [AND] becomes [AND]).
However, only applies to immediate repetition, not across the entire input
*/
func CollapseAdjacentOperators(inputArray []string, valuesToCollapse []string) []string {

	// Output structure
	outSlice := []string{}

	// Iterate over input
	for _, v := range inputArray {
		// If first round, simply append
		if len(outSlice) == 0 {
			outSlice = append(outSlice, v)
			continue
		}
		// If last value is the same as this one
		res, _ := StringInSlice(v, valuesToCollapse)
		if res {
			// do nothing
		} else {
			// else append
			outSlice = append(outSlice, v)
		}
	}
	return outSlice
}

/*
Indicates whether a given node is contained within a slice of nodes
*/
func NodeInSlice(a *Node, list []*Node) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

/*
Indicates whether a particular value is contained in a slice of strings,
and if so, indicates index in slice (else -1).
*/
func StringInSlice(a string, list []string) (bool, int) {
	for i, b := range list {
		if b == a {
			return true, i
		}
	}
	return false, -1
}

/*
Merges two slices. Values of the second slice are added after the previous shared value
in the bigger slice. Checks also for substring matches if subItemSeparator is provided
(i.e., != ""). If shared entry cannot be found, the deviating entries will
be appended at the end.
Similarity of elements is assessed based on matching on leading substring (before subItemSeparator),
e.g., substring before "_".
Input:
- Two arrays to be merged
- Separator indicating subitem prefix (e.g., before "_") to facilitate match
*/
func MergeSlices(array1 []string, array2 []string, subItemSeparator string) []string {

	Println("Slice 1 before merge: ", array1)
	Println("Slice 2 before merge: ", array2)

	result := array1
	arrayToIterate := array2

	// Iterate through smaller array
	for _, v := range arrayToIterate {
		// See if element of smaller array is already in larger array
		res, _ := StringInSlice(v, result)
		if res {
			// if so, skip addition
			continue
		}

		// Now perform fuzzy match ...

		// Check whether subitems match based on substring prior subItemSeparator (considers both input and target array entries)
		lastSimilarElement := FindLastSimilarElement(result, v, subItemSeparator)

		if lastSimilarElement != -1 {
			Println("Found last similar element for item ", v, " on position: ", lastSimilarElement)

			// Add element at position following last shared index

			// Append empty element (value does not matter)
			result = append(result, "placeholder")
			// Shift content from given position one to the right
			copy(result[lastSimilarElement+2:], result[lastSimilarElement+1:])
			// Insert new element at given position
			result[lastSimilarElement+1] = v
		} else {
			Println("No similar element found for item ", v, ", appending at the end.")
			// Append at the end of the array
			result = append(result, v)
		}
	}
	return result
}

/*
Returns index of last similar element based on prefix conventions (e.g., I_1, I_2) in input array.
Input:
- array to iterate over
- item to look up
- string indicative of subitem (indication of "similarity" with other elements)

Returns -1 if no similar item found
*/
func FindLastSimilarElement(arrayToIterate []string, itemToTest string, subItemSep string) int {

	// Prepare modified search item for similarity match
	substringedSearchItem := itemToTest
	// Check whether subitem separator (e.g., "_") is contained in search item
	substringedSearchItemIdx := -1
	if subItemSep != "" {
		substringedSearchItemIdx = strings.Index(itemToTest, subItemSep)
	}
	if substringedSearchItemIdx != -1 {
		// Remove trailing substring if match exists
		substringedSearchItem = itemToTest[:substringedSearchItemIdx]
	}
	Println("Input:", itemToTest)
	Println("Separator:", subItemSep)
	Println("Preprocessed:", substringedSearchItem)

	// Index of last similar item
	similarIndex := -1
	for i, v := range arrayToIterate {
		targetItem := v
		// Determine substring on target item
		targetItemIdx := -1
		if subItemSep != "" {
			targetItemIdx = strings.Index(v, subItemSep)
		}
		if targetItemIdx != -1 {
			// Remove trailing substring if match exists
			targetItem = v[:targetItemIdx]
			//Println("Preprocessed target item:", targetItem)
		}

		// If the current value matches search item ...
		if targetItem == substringedSearchItem {
			// ... then save the index
			similarIndex = i
			//Println("Items match")
		} else {
			//Println("Items do not match")
		}
	}

	// Return result
	return similarIndex
}

/*
Moves element in array to new position.
Input:
- index of element to be moved
- target index the element is to be moved to
- array that the operation is performed on
Output:
- array with element moved to target position
*/
func MoveElementToNewPosition(indexToTakeFrom int, indexToMoveTo int, arrayToOperateOn []string) []string {
	// Assign array to operate on
	sourceArray := arrayToOperateOn
	// Element to be moved
	val := sourceArray[indexToTakeFrom]
	// Create array consisting only of array without element to be moved
	sourceArray = append(sourceArray[:indexToTakeFrom], sourceArray[indexToTakeFrom+1:]...)
	// Create new array up to position in which element is to be inserted
	newSlice := make([]string, indexToMoveTo+1)
	// Copy old sourceArray (without element) into new one
	copy(newSlice, sourceArray[:indexToMoveTo])
	// Append element at desired target position
	newSlice[indexToMoveTo] = val
	// Fill up remaining elements
	sourceArray = append(newSlice, sourceArray[indexToMoveTo:]...)
	Println("Revised sourceArray:", sourceArray)

	return sourceArray
}

/*
Prints a given array in human-readable form (with comma separation)
*/
func PrintArray(array []string) string {
	i := 0
	out := ""
	for i < len(array) {
		out += array[i]
		if i < len(array) {
			out += ","
		}
		i++
	}
	return out
}

/*
Flatten input structure into simple array.
*/
func Flatten(input [][]*Node) []*Node {
	output := []*Node{}
	for _, v := range input {
		for _, v2 := range v {
			output = append(output, v2)
		}
	}
	return output
}
