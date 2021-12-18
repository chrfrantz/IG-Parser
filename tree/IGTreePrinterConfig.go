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
