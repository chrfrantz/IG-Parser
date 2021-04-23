package tree

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const refSuffix = "-Ref"
const refNameSuffix = " Reference"
const leftBracket = "["
const rightBracket = "]"

const (
	ATTRIBUTES = "A"
	NAME_ATTRIBUTES = "Attributes"

	ATTRIBUTES_PROPERTY = "A,p"
	NAME_ATTRIBUTES_PROPERTY = "Attributes Property"
	ATTRIBUTES_PROPERTY_REFERENCE = ATTRIBUTES_PROPERTY + refSuffix
	NAME_ATTRIBUTES_PROPERTY_REFERENCE = NAME_ATTRIBUTES_PROPERTY + refNameSuffix

	DEONTIC = "D"
	NAME_DEONTIC = "Deontic"

	AIM = "I"
	NAME_AIM = "Aim"

	DIRECT_OBJECT = "Bdir"
	NAME_DIRECT_OBJECT = "Direct Object"
	DIRECT_OBJECT_REFERENCE = DIRECT_OBJECT + refSuffix
	NAME_DIRECT_OBJECT_REFERENCE = NAME_DIRECT_OBJECT + refNameSuffix

	DIRECT_OBJECT_PROPERTY = "Bdir,p"
	NAME_DIRECT_OBJECT_PROPERTY = "Direct Object Property"
	DIRECT_OBJECT_PROPERTY_REFERENCE = DIRECT_OBJECT_PROPERTY + refSuffix
	NAME_DIRECT_OBJECT_PROPERTY_REFERENCE = NAME_DIRECT_OBJECT_PROPERTY + refNameSuffix

	INDIRECT_OBJECT = "Bind"
	NAME_INDIRECT_OBJECT = "Indirect Object"
	INDIRECT_OBJECT_REFERENCE = INDIRECT_OBJECT + refSuffix
	NAME_INDIRECT_OBJECT_REFERENCE = NAME_INDIRECT_OBJECT + refNameSuffix

	INDIRECT_OBJECT_PROPERTY = "Bind,p"
	NAME_INDIRECT_OBJECT_PROPERTY = "Indirect Object Property"
	INDIRECT_OBJECT_PROPERTY_REFERENCE = INDIRECT_OBJECT_PROPERTY + refSuffix
	NAME_INDIRECT_OBJECT_PROPERTY_REFERENCE = NAME_INDIRECT_OBJECT_PROPERTY + refNameSuffix

	ACTIVATION_CONDITION = "Cac"
	NAME_ACTIVATION_CONDITION = "Activation Condition"
	ACTIVATION_CONDITION_REFERENCE = ACTIVATION_CONDITION + refSuffix
	NAME_ACTIVATION_CONDITION_REFERENCE = NAME_ACTIVATION_CONDITION + refNameSuffix

	EXECUTION_CONSTRAINT = "Cex"
	NAME_EXECUTION_CONSTRAINT = "Execution Constraint"
	EXECUTION_CONSTRAINT_REFERENCE = EXECUTION_CONSTRAINT + refSuffix
	NAME_EXECUTION_CONSTRAINT_REFERENCE = NAME_EXECUTION_CONSTRAINT + refNameSuffix

	CONSTITUTED_ENTITY = "E"
	NAME_CONSTITUTED_ENTITY = "Constituted Entity"

	CONSTITUTED_ENTITY_PROPERTY = "E,p"
	NAME_CONSTITUTED_ENTITY_PROPERTY = "Constituted Entity Property"
	CONSTITUTED_ENTITY_PROPERTY_REFERENCE = CONSTITUTED_ENTITY_PROPERTY + refSuffix
	NAME_CONSTITUTED_ENTITY_PROPERTY_REFERENCE = NAME_CONSTITUTED_ENTITY_PROPERTY + refNameSuffix

	MODAL = "M"
	NAME_MODAL = "Modal"

	CONSTITUTIVE_FUNCTION = "F"
	NAME_CONSTITUTIVE_FUNCTION = "Constitutive Function"

	CONSTITUTING_PROPERTIES = "P"
	NAME_CONSTITUTING_PROPERTIES = "Constituting Properties"
	CONSTITUTING_PROPERTIES_REFERENCE = CONSTITUTING_PROPERTIES + refSuffix
	NAME_CONSTITUTING_PROPERTIES_REFERENCE = NAME_CONSTITUTING_PROPERTIES + refNameSuffix

	CONSTITUTING_PROPERTIES_PROPERTY = "P,p"
	NAME_CONSTITUTING_PROPERTIES_PROPERTY = "Constituting Properties Properties"
	CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE = CONSTITUTING_PROPERTIES_PROPERTY + refSuffix
	NAME_CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE = NAME_CONSTITUTING_PROPERTIES_PROPERTY + refNameSuffix

	OR_ELSE = "O"
	NAME_OR_ELSE = "Or else"

	SAND = "sAND"
	AND = "AND"
	OR = "OR"
	XOR = "XOR"
	NOT = "NOT"

	SAND_BRACKETS = leftBracket + SAND + rightBracket
	AND_BRACKETS = leftBracket + AND + rightBracket
	OR_BRACKETS = leftBracket + OR + rightBracket
	XOR_BRACKETS = leftBracket + XOR + rightBracket
	NOT_BRACKETS = leftBracket + NOT + rightBracket

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
	ATTRIBUTES_PROPERTY_REFERENCE,
	DEONTIC,
	AIM,
	DIRECT_OBJECT,
	DIRECT_OBJECT_REFERENCE,
	DIRECT_OBJECT_PROPERTY,
	DIRECT_OBJECT_PROPERTY_REFERENCE,
	INDIRECT_OBJECT,
	INDIRECT_OBJECT_REFERENCE,
	INDIRECT_OBJECT_PROPERTY,
	INDIRECT_OBJECT_PROPERTY_REFERENCE,
	ACTIVATION_CONDITION,
	ACTIVATION_CONDITION_REFERENCE,
	EXECUTION_CONSTRAINT,
	EXECUTION_CONSTRAINT_REFERENCE,
	CONSTITUTED_ENTITY,
	CONSTITUTED_ENTITY_PROPERTY,
	CONSTITUTED_ENTITY_PROPERTY_REFERENCE,
	MODAL,
	CONSTITUTIVE_FUNCTION,
	CONSTITUTING_PROPERTIES,
	CONSTITUTING_PROPERTIES_REFERENCE,
	CONSTITUTING_PROPERTIES_PROPERTY,
	CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE,
    OR_ELSE}

/*
IG 2.0 Component Symbols
*/
var IGComponentNames = []string{
	NAME_ATTRIBUTES,
	NAME_ATTRIBUTES_PROPERTY,
	NAME_ATTRIBUTES_PROPERTY_REFERENCE,
	NAME_DEONTIC,
	NAME_AIM,
	NAME_DIRECT_OBJECT,
	NAME_DIRECT_OBJECT_REFERENCE,
	NAME_DIRECT_OBJECT_PROPERTY,
	NAME_DIRECT_OBJECT_PROPERTY_REFERENCE,
	NAME_INDIRECT_OBJECT,
	NAME_INDIRECT_OBJECT_REFERENCE,
	NAME_INDIRECT_OBJECT_PROPERTY,
	NAME_INDIRECT_OBJECT_PROPERTY_REFERENCE,
	NAME_ACTIVATION_CONDITION,
	NAME_ACTIVATION_CONDITION_REFERENCE,
	NAME_EXECUTION_CONSTRAINT,
	NAME_EXECUTION_CONSTRAINT_REFERENCE,
	NAME_CONSTITUTED_ENTITY,
	NAME_CONSTITUTED_ENTITY_PROPERTY,
	NAME_CONSTITUTED_ENTITY_PROPERTY_REFERENCE,
	NAME_MODAL,
	NAME_CONSTITUTIVE_FUNCTION,
	NAME_CONSTITUTING_PROPERTIES,
	NAME_CONSTITUTING_PROPERTIES_REFERENCE,
	NAME_CONSTITUTING_PROPERTIES_PROPERTY,
	NAME_CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE,
    NAME_OR_ELSE}

/*
Map holding mapping from IG 2.0 component symbols to proper component names
 */
var IGComponentSymbolNameMap = map[string]string{
	ATTRIBUTES: NAME_ATTRIBUTES,
	ATTRIBUTES_PROPERTY: NAME_ATTRIBUTES_PROPERTY,
	ATTRIBUTES_PROPERTY_REFERENCE: NAME_ATTRIBUTES_PROPERTY_REFERENCE,

	DEONTIC: NAME_DEONTIC,
	AIM: NAME_AIM,

	DIRECT_OBJECT: NAME_DIRECT_OBJECT,
	DIRECT_OBJECT_REFERENCE: NAME_DIRECT_OBJECT_REFERENCE,
	DIRECT_OBJECT_PROPERTY: NAME_DIRECT_OBJECT_PROPERTY,
	DIRECT_OBJECT_PROPERTY_REFERENCE: NAME_DIRECT_OBJECT_PROPERTY_REFERENCE,

	INDIRECT_OBJECT: NAME_INDIRECT_OBJECT,
	INDIRECT_OBJECT_REFERENCE: NAME_INDIRECT_OBJECT_REFERENCE,
	INDIRECT_OBJECT_PROPERTY: NAME_INDIRECT_OBJECT_PROPERTY,
	INDIRECT_OBJECT_PROPERTY_REFERENCE: NAME_INDIRECT_OBJECT_PROPERTY_REFERENCE,

	ACTIVATION_CONDITION: NAME_ACTIVATION_CONDITION,
	ACTIVATION_CONDITION_REFERENCE: NAME_ACTIVATION_CONDITION_REFERENCE,

	EXECUTION_CONSTRAINT: NAME_EXECUTION_CONSTRAINT,
	EXECUTION_CONSTRAINT_REFERENCE: NAME_EXECUTION_CONSTRAINT_REFERENCE,

	CONSTITUTED_ENTITY: NAME_CONSTITUTED_ENTITY,
	CONSTITUTED_ENTITY_PROPERTY: NAME_CONSTITUTED_ENTITY_PROPERTY,
	CONSTITUTED_ENTITY_PROPERTY_REFERENCE: NAME_CONSTITUTED_ENTITY_PROPERTY_REFERENCE,

	MODAL: NAME_MODAL,
	CONSTITUTIVE_FUNCTION: NAME_CONSTITUTIVE_FUNCTION,

	CONSTITUTING_PROPERTIES: NAME_CONSTITUTING_PROPERTIES,
	CONSTITUTING_PROPERTIES_REFERENCE: NAME_CONSTITUTING_PROPERTIES_REFERENCE,
	CONSTITUTING_PROPERTIES_PROPERTY: NAME_CONSTITUTING_PROPERTIES_PROPERTY,
	CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE: NAME_CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE,

	OR_ELSE: NAME_OR_ELSE}

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
// Detecting combinations of nested statements with varying component references
// (e.g., {Cac{stmt1} [AND] Cex{stmt2}}, but should be{Cac{stmt1} [AND] Cac{stmt2}})
const PARSING_ERROR_INVALID_TYPES_IN_NESTED_STATEMENT_COMBINATION = "INVALID_TYPE_COMBINATIONS_IN_NESTED_STATEMENT_COMBINATIONS"
// Indicates that operations was imposed on nil element
const PARSING_ERROR_NIL_ELEMENT = "INVALID_PARSING_OF_NIL_ELEMENT"

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

	fmt.Println("Slice 1 before merge: ", array1)
	fmt.Println("Slice 2 before merge: ", array2)

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
			fmt.Println("Found last similar element for item ", v, " on position: ", lastSimilarElement)

			// Add element at position following last shared index

			// Append empty element (value does not matter)
			result = append(result, "placeholder")
			// Shift content from given position one to the right
			copy(result[lastSimilarElement+2:], result[lastSimilarElement+1:])
			// Insert new element at given position
			result[lastSimilarElement+1] = v
		} else {
			fmt.Println("No similar element found for item ", v, ", appending at the end.")
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
	fmt.Println("Input:", itemToTest)
	fmt.Println("Separator:", subItemSep)
	fmt.Println("Preprocessed:", substringedSearchItem)

	// Index of last similar item
	similarIndex := -1
	for i, v :=range arrayToIterate {
		targetItem := v
		// Determine substring on target item
		targetItemIdx := -1
		if subItemSep != "" {
			targetItemIdx = strings.Index(v, subItemSep)
		}
		if targetItemIdx != -1 {
			// Remove trailing substring if match exists
			targetItem = v[:targetItemIdx]
			//fmt.Println("Preprocessed target item:", targetItem)
		}

		// If the current value matches search item ...
		if targetItem == substringedSearchItem {
			// ... then save the index
			similarIndex = i
			//fmt.Println("Items match")
		} else {
			//fmt.Println("Items do not match")
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