package tree

/*
This file contains configuration information for IGTreePrinter.go
*/

// Indicates whether output for properties is complex (as a tree structure nested on the corresponding component),
// or a flat listing of properties attached as labels to the component.
var print_FLAT = false

/*
Enables flat printing of properties in generated output (as opposed to tree structure).
*/
func SetFlatPrinting(printFlat bool) {
	print_FLAT = printFlat
}

/*
Indicates whether flat printing is enabled.
*/
func FlatPrinting() bool {
	return print_FLAT
}

// Indicates whether visual output is a binary tree structure, or collapses identical logical operators for given component types

var print_BINARY = false

/*
Enables binary printing parsed tree output (as opposed to collapsing entries linked via the same logical operators for a given component).
*/
func SetBinaryPrinting(printBinary bool) {
	print_BINARY = printBinary
}

/*
Indicates whether binary printing (i.e., production of binary tree as output) is enabled.
*/
func BinaryPrinting() bool {
	return print_BINARY
}

/*
Indicates whether activation conditions are moved to the top of the output
*/
var moveActivationConditionsToFront = false

/*
Defines whether activation conditions are moved to the top of the visual tree output.
*/
func SetMoveActivationConditionsToFront(moveToFront bool) {
	moveActivationConditionsToFront = moveToFront
}

/*
Indicates whether activation condition output should be moved to the top in the visual tree output.
*/
func MoveActivationConditionsToFront() bool {
	return moveActivationConditionsToFront
}

/*
Indicates whether shared elements (i.e., preceding or following core entry in the case of scoped combinations)
are included in output. Example: Bdir(sharedLeft (leftOption [XOR] rightOption) rightShared)
*/
var includeSharedElementsInVisualOutput = true

/*
Indicates whether shared elements are included in visual output.
*/
func IncludeSharedElementsInVisualOutput() bool {
	return includeSharedElementsInVisualOutput
}

/*
Sets visual output parameter indicating whether shared elements are included in output.
*/
func SetIncludeSharedElementsInVisualOutput(includeSharedElements bool) {
	includeSharedElementsInVisualOutput = includeSharedElements
}
