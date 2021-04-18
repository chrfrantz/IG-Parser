package tree

import (
	"errors"
	"fmt"
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
	ACTIVATION_CONDITION_REFERENCE = "Cac-Ref"
	NAME_ACTIVATION_CONDITION_REFERENCE = "Activation Condition Reference"
	EXECUTION_CONSTRAINT = "Cex"
	NAME_EXECUTION_CONSTRAINT = "Execution Constraint"
	EXECUTION_CONSTRAINT_REFERENCE = "Cex-Ref"
	NAME_EXECUTION_CONSTRAINT_REFERENCE = "Execution Constraint Reference"
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
	ATTRIBUTES,
	ATTRIBUTES_PROPERTY,
	DEONTIC,
	AIM,
	DIRECT_OBJECT,
	DIRECT_OBJECT_PROPERTY,
	INDIRECT_OBJECT,
	INDIRECT_OBJECT_PROPERTY,
	ACTIVATION_CONDITION,
	ACTIVATION_CONDITION_REFERENCE,
	EXECUTION_CONSTRAINT,
	EXECUTION_CONSTRAINT_REFERENCE,
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
	NAME_ACTIVATION_CONDITION_REFERENCE,
	NAME_EXECUTION_CONSTRAINT,
	NAME_EXECUTION_CONSTRAINT_REFERENCE,
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
	ACTIVATION_CONDITION_REFERENCE: NAME_ACTIVATION_CONDITION_REFERENCE,
	EXECUTION_CONSTRAINT: NAME_EXECUTION_CONSTRAINT,
	EXECUTION_CONSTRAINT_REFERENCE: NAME_EXECUTION_CONSTRAINT_REFERENCE,
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
// Write error
const PARSING_ERROR_WRITE = "WRITE_ERROR"
// Invalid parentheses/braces combinations
const PARSING_ERROR_INVALID_PARENTHESES_COMBINATION = "INVALID_PARENTHESES_COMBINATIONS"
// Error during regex compilation
const PARSING_ERROR_PATTERN_EXTRACTION = "PATTERN_EXTRACTION_ERROR"

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
			res, _ := StringInSlice(v, valuesToCollapse)
			if res {
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
Merges two slices. Values of the smaller slice are added after the previous shared value
with the bigger slice. If no last shared entry could be found, the deviating entries will
be appended at the end.
 */
func MergeSlices(array1 []string, array2 []string) []string {
	result := []string{}
	arrayToIterate := []string{}

	// Figure out which array is larger
	if len(array1) >= len(array2) {
		result = array1
		arrayToIterate = array2
	} else {
		result = array2
		arrayToIterate = array1
	}

	// Store last match
	indexOfLastIdenticalElement := -1

	// Iterate through smaller array
	for i, v := range arrayToIterate {
		// See if element of smaller array is already in larger array
		res, idx := StringInSlice(v, result)
		if res {
			// if so, update index of last match
			indexOfLastIdenticalElement = idx
		} else {
			// if it is not first element and some shared elements have been found
			if i != 0 && indexOfLastIdenticalElement != -1 {
				// Add element at position following last shared index

				// Append empty element (value does not matter)
				result = append(result, "placeholder")
				// Shift content from given position one to the right
				copy(result[indexOfLastIdenticalElement+2:], result[indexOfLastIdenticalElement+1:])
				// Insert new element at given position
				result[indexOfLastIdenticalElement+1] = v
			} else {
				// else append at the end of result array
				result = append(result, v)
			}
		}
	}
	return result
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
	fmt.Println("Revised sourceArray:", sourceArray)

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
