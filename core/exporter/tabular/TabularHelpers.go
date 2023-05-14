package tabular

import (
	"IG-Parser/core/tree"
)

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
