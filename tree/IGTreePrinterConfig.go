package tree

/*
This file contains configuration information for IGTreePrinter.go
*/

// Indicates whether output for properties is complex (as a tree structure nested on the corresponding component),
// or a flat listing of properties attached as labels to the component.
var PRINT_FLAT = false

/*
Enables flat printing of properties in generated output (as opposed to tree structure).
*/
func SetFlatPrinting(printFlat bool) {
	PRINT_FLAT = printFlat
}

/*
Indicates whether flat printing is enabled.
*/
func FlatPrinting() bool {
	return PRINT_FLAT
}
