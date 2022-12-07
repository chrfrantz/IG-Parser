package shared

import "testing"

/*
Tests the replacing of symbols for export.
*/
func TestReplaceSymbolsForExport(t *testing.T) {

	input := "\"left\" middle next \"center\" next \"right\""

	// Check whether escaping worked as expected
	if EscapeSymbolsForExport(input) != "'left' middle next 'center' next 'right'" {
		t.Error("Escaping of values did not result in expected outcome.")
	}
}

/*
Tests the identification of maximum value in a given array, with a minimum threshold.
*/
func TestFindMaxValue(t *testing.T) {

	// Highest value is default value
	vals := []int{27, 25, 12, 2}
	defaultVal := 30

	res1 := FindMaxValue(vals, defaultVal)

	if res1 != defaultVal {
		t.Error("Returned value is not default value. Return value:", res1)
	}

	// Highest value in array
	vals = []int{27, 25, 12, 2}
	defaultVal = 20

	res2 := FindMaxValue(vals, defaultVal)

	if res2 != 27 {
		t.Error("Returned value is not default value. Return value:", res1)
	}

	// Test with negative values
	vals = []int{-27, 25, 12, 2}
	defaultVal = 20

	res3 := FindMaxValue(vals, defaultVal)

	if res3 != 25 {
		t.Error("Returned value is not default value. Return value:", res1)
	}

}

/*
Tests aggregation function, which operates on values beyond a given threshold.
*/
func TestAggregateIfGreaterThan(t *testing.T) {

	// Test reporting of default value
	vals := []int{-27, 25, 12, 2}
	threshold := 3
	defaultVal := 50

	res1 := AggregateIfGreaterThan(vals, threshold, defaultVal)
	if res1 != defaultVal {
		t.Error("Wrong value reported in test:", res1)
	}

	// Test reporting of sum
	vals = []int{27, 25, 12, 2}

	res2 := AggregateIfGreaterThan(vals, threshold, defaultVal)
	if res2 != 64 {
		t.Error("Wrong value reported in test:", res2)
	}

	// Test negative threshold
	vals = []int{27, 25, 12, 2}
	threshold = -3

	res3 := AggregateIfGreaterThan(vals, threshold, defaultVal)
	if res3 != 66 {
		t.Error("Wrong value reported in test:", res3)
	}

}
