package shared

import (
	"strings"
)

/*
This file provides generic helper functions relevant across packages.
*/

/*
Replaces/escapes selected symbols in as far as relevant for export (e.g., quotation marks).
*/
func EscapeSymbolsForExport(rawValue string) string {
	// Replace quotation marks with single quotes
	return strings.ReplaceAll(rawValue, "\"", "'")
}

/*
Aggregates values in a given array that exceed a given threshold value.
Returns sum, if it exceeds a given default value. Otherwise it returns the given default value.
*/
func AggregateIfGreaterThan(arr []int, threshold int, defaultValue int) int {
	sum := 0
	for i := 0; i < len(arr); i++ {
		if arr[i] > threshold {
			// Only aggregate values above threshold
			sum += arr[i]
		}
	}
	if sum > defaultValue {
		// Only return sum if greater than default value
		return sum
	} else {
		return defaultValue
	}
}

/*
Return the maximum value for all values in a given array.
If none of the values is higher, it returns a given default value.
*/
func FindMaxValue(arr []int, defaultValue int) int {
	max := 0
	for i := 0; i < len(arr); i++ {
		if arr[i] > max {
			max += arr[i]
		}
	}
	if max > defaultValue {
		return max
	} else {
		return defaultValue
	}
}

/*
Stringifies slices in whitespace-separated string. Does not add whitespace at beginning or end.
*/
func StringifySlices(elements []string) string {
	outString := ""
	for i, v := range elements {
		outString += v
		if i < len(elements)-1 {
			outString += " "
		}
	}
	return outString
}

/*
Detects the presence of duplicate components (i.e., component type and content).
Returns empty string ("") if no duplicate component is found, else the first identified duplicate entry.
*/
func DuplicateElement(array []string) string {
	// Iterate through array and compare each element with each other
	for i := 0; i < len(array); i++ {
		// Start at element i + 1
		for j := i + 1; j < len(array); j++ {
			// if element at index i is equal to element at index j ...
			if array[i] == array[j] {
				// ... then return duplicate element
				return array[i]
			}
		}
	}
	// Return empty string - i.e., no duplicate found
	return ""
}
