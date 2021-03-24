package tree

import (
	"fmt"
	"strconv"
)

const (
	ATTRIBUTES = "A"
	ATTRIBUTES_PROPERTY = "A,p"
	DEONTIC = "D"
	AIM = "I"
	DIRECT_OBJECT = "Bdir"
	DIRECT_OBJECT_PROPERTY = "Bdir,p"
	INDIRECT_OBJECT = "Bind"
	INDIRECT_OBJECT_PROPERTY = "Bind,p"
	ACTIVATION_CONDITION = "Cac"
	EXECUTION_CONSTRAINT = "Cex"
	CONSTITUTED_ENTITY = "E"
	CONSTITUTED_ENTITY_PROPERTY = "E,p"
	MODAL = "M"
	CONSTITUTIVE_FUNCTION = "F"
	CONSTITUTING_PROPERTIES = "P"
	CONSTITUTING_PROPERTIES_PROPERTY = "P,p"
	SAND = "sAND"
	AND = "AND"
	OR = "OR"
	XOR = "XOR"
	NOT = "NOT"
	SAND_BRACKETS = "[" + SAND + "]"
	AND_BRACKETS = "[" + AND + "]"
	OR_BRACKETS = "[" + OR + "]"
	XOR_BRACKETS = "[" + XOR + "]"
	NOT_BRACKETS = "[" + NOT + "]"
	PARSING_MODE_LEFT = "PARSING_LEFT"
	PARSING_MODE_RIGHT = "PARSING_RIGHT"
	PARSING_MODE_OUTSIDE_EXPRESSION = "PARSING_OUTSIDE"
)

/*
type igComponent struct {
	ComponentName string
}

func IGComponent(componentName string) *igComponent {
	i := igComponent{}
	i.ComponentName = componentName
	return &i
}

func (i igComponent) String() string {
	return i.ComponentName
}
/*
Checks whether component value is valid (i.e., a valid IG Component symbol).
 */
/*func (c *igComponent) valid() bool {
	return StringInSlice(c.ComponentName, igComponents)
}*/

/*
IG 2.0 Component Symbols
 */
var IGComponents = []string{
	ATTRIBUTES,
	ATTRIBUTES_PROPERTY,
	DEONTIC,
	AIM,
	DIRECT_OBJECT,
	DIRECT_OBJECT_PROPERTY,
	INDIRECT_OBJECT,
	INDIRECT_OBJECT_PROPERTY,
	ACTIVATION_CONDITION,
	EXECUTION_CONSTRAINT,
	CONSTITUTED_ENTITY,
	CONSTITUTED_ENTITY_PROPERTY,
	MODAL,
	CONSTITUTIVE_FUNCTION,
	CONSTITUTING_PROPERTIES,
	CONSTITUTING_PROPERTIES_PROPERTY}

type igLogicalOperator struct {
	LogicalOperatorName string
}

/*
Checks whether operator value is valid (i.e., a valid logical operator symbol).
*/
func (o *igLogicalOperator) valid() bool {
	return StringInSlice(o.LogicalOperatorName, IGLogicalOperators)
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

/**
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
	// Signals whether a boundary value has already been added to the final output
	//AlreadyAdded bool
}

func (b *Boundaries) String() string {
	return "Boundaries{\n"+
		"  Left: " + strconv.Itoa(b.Left) + "\n" +
		"  Op Pos: " + strconv.Itoa(b.Operator) + "\n" +
		"  Operator: " + b.OperatorVal + "\n" +
		"  Right: " + strconv.Itoa(b.Right) + "\n" +
		"  Complete: " +  strconv.FormatBool(b.Complete) + "\n}"
}


// Signals invalid component combinations on a given parsing level (generally non-AND components)
const PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS = "INVALID_LOGICAL_OPERATOR_COMBINATIONS"
// Signals the detection of a logical operator outside a combination (i.e., no left or right side)
const PARSING_ERROR_LOGICAL_OPERATOR_OUTSIDE_COMBINATION = "LOGICAL_OPERATOR_OUTSIDE_COMBINATION"
// Indicates that there were no issues during parsing
const PARSING_NO_ERROR = "NO_ERROR_DURING_PARSING"
// Signals that no combinations were found in input
const PARSING_NO_COMBINATIONS = "NO_COMBINATIONS_IN_INPUT"
// Signals invalid combination expression (e.g., missing left, right, or operator)
const PARSING_INVALID_COMBINATION = "INVALID_COMBINATION_IN_INPUT"
// Signals empty leaf value during parsing
const PARSING_ERROR_EMPTY_LEAF = "EMPTY_LEAF_VALUE"
// Indicates that parentheses are umbalanced during parsing
const PARSING_ERROR_IMBALANCED_PARENTHESES = "UNBALANCED_PARENTHESES"
// Indicates that component was not found during preprocessing
const PARSING_ERROR_COMPONENT_NOT_FOUND = "COMPONENT_NOT_FOUND"
// Indicates ignored elements during parsing
const PARSING_ERROR_IGNORED_ELEMENTS = "IGNORED_ELEMENTS"

type ParsingError struct {
	ErrorCode string
	ErrorMessage string
	ErrorIgnoredElements []string
}

func (e *ParsingError) Error() string {
	return fmt.Sprint("Parsing Error " + e.ErrorCode + ": " + e.ErrorMessage +
		" (Ignored elements: " + strconv.Itoa(len(e.ErrorIgnoredElements)) + ")")
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

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
