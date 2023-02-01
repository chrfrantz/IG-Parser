package tree

import (
	"fmt"
	"testing"
)

func TestValidIGComponentName(t *testing.T) {

	res := validIGComponentName("Attributes")
	if !res {
		t.Fatal("Could not resolve valid component name")
	}

}

func TestValidIGSymbol(t *testing.T) {

	res := ValidIGComponentSymbol("A")
	if !res {
		t.Fatal("Could not resolve valid component name")
	}

}

/*
Tests merging of arrays without consideration of subitems
*/
func TestMergeSlices(t *testing.T) {

	// Array2 has the same size
	array1 := []string{"First", "Second", "Third", "Fourth", "Fifth"}
	array2 := []string{"First", "Second", "Second-and-half", "Fourth", "Sixth"}

	res := MergeSlices(array1, array2, "")

	expected := []string{"First", "Second", "Third", "Fourth", "Fifth", "Second-and-half", "Sixth"}

	if fmt.Sprint(res) != fmt.Sprint(expected) {
		t.Error("Generated array does not reflect expected output (Expected: " +
			fmt.Sprint(expected) + ", actual: " + fmt.Sprint(res))
	}

	// Array2 is larger
	array1 = []string{"First", "Second", "Third", "Fourth", "Fifth"}
	array2 = []string{"First", "Second", "Second-and-half", "Fourth", "Sixth", "Seventh"}

	res = MergeSlices(array1, array2, "")

	expected = []string{"First", "Second", "Third", "Fourth", "Fifth", "Second-and-half", "Sixth", "Seventh"}

	if fmt.Sprint(res) != fmt.Sprint(expected) {
		t.Error("Generated array does not reflect expected output (Expected: " +
			fmt.Sprint(expected) + ", actual: " + fmt.Sprint(res))
	}

}

/*
Tests the merging of slices under consideration of matching based on sub item separator.
*/
func TestMergeSlicesWithMatching(t *testing.T) {

	array1 := []string{"First", "Second", "Third", "Fourth", "Fifth"}
	array2 := []string{"First", "Second", "Second-and-half", "Fourth", "Sixth"}

	res := MergeSlices(array1, array2, "-")

	expected := []string{"First", "Second", "Second-and-half", "Third", "Fourth", "Fifth", "Sixth"}

	if fmt.Sprint(res) != fmt.Sprint(expected) {
		t.Error("Generated array does not reflect expected output (Expected: " +
			fmt.Sprint(expected) + ", actual: " + fmt.Sprint(res))
	}
}

func TestFindLastSimilarElement(t *testing.T) {

	array1 := []string{"First", "Second", "Third", "Fourth", "Fifth"}
	item := "Second_1"

	idx := FindLastSimilarElement(array1, item, "_")

	fmt.Println(idx)

	if idx != 1 {
		t.Error("Did not correctly identify similar element. Expected: 1, Actual", idx)
	}

	// Test on subitem in input array
	array1 = []string{"First", "Second", "Second_2", "Third", "Fourth", "Fifth"}
	item = "Second_1"

	idx = FindLastSimilarElement(array1, item, "_")

	fmt.Println(idx)

	if idx != 2 {
		t.Error("Did not correctly identify similar element. Expected: 1, Actual", idx)
	}

	// Test on identical subitem in input array should yield same result
	array1 = []string{"First", "Second", "Second_1", "Third", "Fourth", "Fifth"}
	item = "Second_1"

	idx = FindLastSimilarElement(array1, item, "_")

	fmt.Println(idx)

	if idx != 2 {
		t.Error("Did not correctly identify similar element. Expected: 1, Actual", idx)
	}

}

func TestMergeSlicesSubItemMatch(t *testing.T) {

	array1 := []string{"First", "Second", "Third", "Fourth", "Fifth"}
	array2 := []string{"First", "Second_1", "Second_2", "Fourth", "Sixth"}

	res := MergeSlices(array1, array2, "_")

	expected := []string{"First", "Second", "Second_1", "Second_2", "Third", "Fourth", "Fifth", "Sixth"}

	if fmt.Sprint(res) != fmt.Sprint(expected) {
		t.Error("Generated array does not reflect expected output (Expected: " +
			fmt.Sprint(expected) + ", actual: " + fmt.Sprint(res))
	}

}

/*
Tests collapsing adjacent operators.
*/
func TestCollapseAdjacentOperators(t *testing.T) {

	input := []string{"AND", "AND", "AND", "sAND", "XOR"}
	collapseValues := []string{AND, SAND_WITHIN_COMPONENTS, SAND_BETWEEN_COMPONENTS}

	result := CollapseAdjacentOperators(input, collapseValues)
	if len(result) != 3 {
		t.Fatal("Did not identify correct number of operators. Found ", len(result), "in", result)
	}

}
