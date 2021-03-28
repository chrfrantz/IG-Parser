package tree

import (
	"errors"
	"strconv"
)

const (
	ATTRIBUTES = "A"
	NAME_ATTRIBUTES = "Attributes"
	ATTRIBUTES_PROPERTY = "A,p"
	NAME_ATTRIBUTES_PROPERTY = "Attributes Property"
	DEONTIC = "D"
	NAME_DEONTIC = "Deontic"
	AIM = "I"
	NAME_AIM = "Aim"
	DIRECT_OBJECT = "Bdir"
	NAME_DIRECT_OBJECT = "Direct Object"
	DIRECT_OBJECT_PROPERTY = "Bdir,p"
	NAME_DIRECT_OBJECT_PROPERTY = "Direct Object Property"
	INDIRECT_OBJECT = "Bind"
	NAME_INDIRECT_OBJECT = "Indirect Object"
	INDIRECT_OBJECT_PROPERTY = "Bind,p"
	NAME_INDIRECT_OBJECT_PROPERTY = "Indirect Object Property"
	ACTIVATION_CONDITION = "Cac"
	NAME_ACTIVATION_CONDITION = "Activation Condition"
	EXECUTION_CONSTRAINT = "Cex"
	NAME_EXECUTION_CONSTRAINT = "Execution Constraint"
	CONSTITUTED_ENTITY = "E"
	NAME_CONSTITUTED_ENTITY = "Constituted Entity"
	CONSTITUTED_ENTITY_PROPERTY = "E,p"
	NAME_CONSTITUTED_ENTITY_PROPERTY = "Constituted Entity Property"
	MODAL = "M"
	NAME_MODAL = "Modal"
	CONSTITUTIVE_FUNCTION = "F"
	NAME_CONSTITUTIVE_FUNCTION = "Constitutive Function"
	CONSTITUTING_PROPERTIES = "P"
	NAME_CONSTITUTING_PROPERTIES = "Constituting Properties"
	CONSTITUTING_PROPERTIES_PROPERTY = "P,p"
	NAME_CONSTITUTING_PROPERTIES_PROPERTY = "Constituting Properties Properties"
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
Indicates whether a given symbol is a valid IG Component symbol
 */
func validIGComponentSymbol(symbol string) bool {
	return StringInSlice(symbol, IGComponentSymbols)
}

/*
Indicates whether a given name is a valid IG Component name
 */
func validIGComponentName(name string) bool {
	return StringInSlice(name, IGComponentNames)
}

/*
IG 2.0 Component Symbols
 */
var IGComponentSymbols = []string{
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

/*
IG 2.0 Component Symbols
*/
var IGComponentNames = []string{
	NAME_ATTRIBUTES,
	NAME_ATTRIBUTES_PROPERTY,
	NAME_DEONTIC,
	NAME_AIM,
	NAME_DIRECT_OBJECT,
	NAME_DIRECT_OBJECT_PROPERTY,
	NAME_INDIRECT_OBJECT,
	NAME_INDIRECT_OBJECT_PROPERTY,
	NAME_ACTIVATION_CONDITION,
	NAME_EXECUTION_CONSTRAINT,
	NAME_CONSTITUTED_ENTITY,
	NAME_CONSTITUTED_ENTITY_PROPERTY,
	NAME_MODAL,
	NAME_CONSTITUTIVE_FUNCTION,
	NAME_CONSTITUTING_PROPERTIES,
	NAME_CONSTITUTING_PROPERTIES_PROPERTY}

/*
Map holding mapping from IG 2.0 component symbols to proper component names
 */
var IGComponentSymbolNameMap = map[string]string{
	ATTRIBUTES: NAME_ATTRIBUTES,
	ATTRIBUTES_PROPERTY: NAME_ATTRIBUTES_PROPERTY,
	DEONTIC: NAME_DEONTIC,
	AIM: NAME_AIM,
	DIRECT_OBJECT: NAME_DIRECT_OBJECT,
	DIRECT_OBJECT_PROPERTY: NAME_DIRECT_OBJECT_PROPERTY,
	INDIRECT_OBJECT: NAME_INDIRECT_OBJECT,
	INDIRECT_OBJECT_PROPERTY: NAME_INDIRECT_OBJECT_PROPERTY,
	ACTIVATION_CONDITION: NAME_ACTIVATION_CONDITION,
	EXECUTION_CONSTRAINT: NAME_EXECUTION_CONSTRAINT,
	CONSTITUTED_ENTITY: NAME_CONSTITUTED_ENTITY,
	CONSTITUTED_ENTITY_PROPERTY: NAME_CONSTITUTED_ENTITY_PROPERTY,
	MODAL: NAME_MODAL,
	CONSTITUTIVE_FUNCTION: NAME_CONSTITUTIVE_FUNCTION,
	CONSTITUTING_PROPERTIES: NAME_CONSTITUTING_PROPERTIES,
	CONSTITUTING_PROPERTIES_PROPERTY: NAME_CONSTITUTING_PROPERTIES_PROPERTY}

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
// Indicates problems when generating logical operator references
const PARSING_ERROR_LOGICAL_EXPRESSION_GENERATION = "LOGICAL_EXPRESSION_GENERATION"

/*
Error type signaling errors during statement parsing
 */
type ParsingError struct {
	ErrorCode string
	ErrorMessage string
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
	ErrorCode string
	ErrorMessage string
	ErrorIgnoredElements []string
}

func (e *NodeError) Error() error {
	return errors.New("Node Error " + e.ErrorCode + ": " + e.ErrorMessage +
		" (Ignored elements: " + strconv.Itoa(len(e.ErrorIgnoredElements)) + ")")
}

const TREE_NO_ERROR = "NO_ERROR"
const TREE_INVALID_NODE_ADDITION = "INVALID_NODE_ADDITION"
const TREE_INVALID_NODE_SELF_LINKAGE = "INVALID_NODE_LINKAGE_TO_SELF"
const TREE_INVALID_TREE = "TREE_STRUCTURE_INVALID"
const TREE_INPUT_VALIDATION = "INPUT_VALIDATION"
//const TREE_ALREADY_VISITED = "NODE_ALREADY_VISITED"

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
		if outSlice[len(outSlice) - 1] == v {
			// and in the registered value
			if StringInSlice(v, valuesToCollapse) {
				// do nothing
			} else {
				// else append
				outSlice = append(outSlice, v)
			}
		} else {
			// otherwise simply append
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
Indicates whether a particular value is contained in a slice of strings
 */
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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
