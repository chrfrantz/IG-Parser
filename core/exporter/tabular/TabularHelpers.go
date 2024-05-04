package tabular

import (
	"IG-Parser/core/shared"
	"IG-Parser/core/tree"
	"regexp"
)

/*
This file contains helper functions for the tabular output generation.
*/

/*
Adds element to (the end of) array if not existing.
*/
func addElementIfNotExisting(elementToAdd string, arrayToAddTo []string) []string {
	res, _ := tree.StringInSlice(elementToAdd, arrayToAddTo)
	if !res {
		// Append logical operator header to symbols output
		arrayToAddTo = append(arrayToAddTo, elementToAdd)
	}
	return arrayToAddTo
}

/*
Moves element to first position in given array and returns resulting array. Optionally adds element to array prior to operation.
*/
func moveElementToFirstPosition(elementToMove string, arrayToModify []string, addIfNotExisting bool) []string {
	res, pos := tree.StringInSlice(elementToMove, arrayToModify)
	if !res {
		if addIfNotExisting {
			// Simply add ...
			arrayToModify = append(arrayToModify, elementToMove)
			// ... and call function again
			return moveElementToFirstPosition(elementToMove, arrayToModify, false)
		}
	} else if pos != len(arrayToModify) {
		// Move element to first position
		arrayToModify = tree.MoveElementToNewPosition(pos, 0, arrayToModify)
	}
	return arrayToModify
}

/*
Moves element to last position in given array and returns resulting array. Optionally adds element to array prior to operation.
*/
func moveElementToLastPosition(elementToMove string, arrayToModify []string, addIfNotExisting bool) []string {
	res, pos := tree.StringInSlice(elementToMove, arrayToModify)
	if !res {
		if addIfNotExisting {
			// Simply add ...
			arrayToModify = append(arrayToModify, elementToMove)
			// ... and call function again
			return moveElementToLastPosition(elementToMove, arrayToModify, false)
		}
	} else if pos != len(arrayToModify) {
		// Move element to last position
		arrayToModify = tree.MoveElementToNewPosition(pos, len(arrayToModify)-1, arrayToModify)
	}
	return arrayToModify
}

/*
Generic function to clean input in preparation for tabular output
(substituting line breaks, removing cell separator symbols).
*/
func CleanInput(input string, separator string) string {

	// Remove line breaks
	re := regexp.MustCompile(`\r?\n`)
	input = re.ReplaceAllString(input, " ")
	// Remove separator symbol used in output
	re2 := regexp.MustCompile(`\` + separator)
	input = re2.ReplaceAllString(input, "")

	return input
}

/*
Performs output-specific modification of value prior to inclusion in output.
Returns modified value to be used in corresponding output.
Takes input value (value), as well as indication of output type (outputType)
to administer output-specific adjustments. If outputType is empty (""), no
output-specific modification (other than general substitutions) are performed.
*/
func performOutputSpecificAdjustments(value string, outputType string) string {

	// Substitute specific symbols relevant for output generation (e.g. quotation marks).
	value = shared.EscapeSymbolsForExport(value)

	// Perform Google-specific adjustment
	if outputType == OUTPUT_TYPE_GOOGLE_SHEETS {
		// Duplicate leading ' for proper Google Sheets parsing
		if len(value) > 0 && value[0:1] == "'" {
			value = "'" + value
		}
	}

	return value
}
