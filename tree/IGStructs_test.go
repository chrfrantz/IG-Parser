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

	res := validIGComponentSymbol("A")
	if !res {
		t.Fatal("Could not resolve valid component name")
	}

}

func TestMergeSlices(t *testing.T) {

	// Array2 has the same size
	array1 := []string{"First", "Second", "Third", "Fourth", "Fifth"}
	array2 := []string{"First", "Second", "Second-and-half", "Fourth", "Sixth"}

	res := MergeSlices(array1, array2)

	expected := []string{"First", "Second", "Second-and-half", "Third", "Fourth", "Sixth", "Fifth"}

	if fmt.Sprint(res) != fmt.Sprint(expected) {
		t.Error("Generated array does not reflect expected output (Expected: " +
			fmt.Sprint(expected) + ", actual: " + fmt.Sprint(res))
	}

	// Array2 is larger
	array1 = []string{"First", "Second", "Third", "Fourth", "Fifth"}
	array2 = []string{"First", "Second", "Second-and-half", "Fourth", "Sixth", "Seventh"}

	res = MergeSlices(array1, array2)

	expected = []string{"First", "Second", "Third", "Second-and-half", "Fourth", "Fifth", "Sixth", "Seventh"}

	if fmt.Sprint(res) != fmt.Sprint(expected) {
		t.Error("Generated array does not reflect expected output (Expected: " +
			fmt.Sprint(expected) + ", actual: " + fmt.Sprint(res))
	}

}