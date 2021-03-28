package exporter

import (
	"fmt"
	"testing"
)

func TestGenerateRangesInStringSlice(t *testing.T) {
	input := []int{1,2,3,4,6,7,8,10,11}

	slc := []string{}

	for _, v := range input {
		slc = GenerateReferenceSlice(slc, v, true, false)
	}

	fmt.Println(slc)
	if len(slc) != 3 {
		t.Fatal("Resulting slice length is wrong: ", len(slc))
	}
	if slc[0] != "1-4" || slc[1] != "6-8" || slc[2] != "10-11" {
		t.Fatal("Wrong values in generated slice. Values: ", slc)
	}
}

func TestGenerateRangesInStringSliceAndIncrementValues(t *testing.T) {
	input := []int{1,2,3,4,6,7,8,10,11}

	slc := []string{}

	for _, v := range input {
		slc = GenerateReferenceSlice(slc, v, true, true)
	}

	fmt.Println(slc)
	if len(slc) != 3 {
		t.Fatal("Resulting slice length is wrong: ", len(slc))
	}
	if slc[0] != "2-5" || slc[1] != "7-9" || slc[2] != "11-12" {
		t.Fatal("Wrong values in generated slice. Values: ", slc)
	}
}

func TestGenerateRangesInStringSliceAndIncrementValuesPrepopulatedSlice(t *testing.T) {
	input := []int{1,2,3,4,6,7,8,10,11}

	slc := []string{"1","2"}

	for _, v := range input {
		slc = GenerateReferenceSlice(slc, v, true, true)
	}

	fmt.Println(slc)
	if len(slc) != 4 {
		t.Fatal("Resulting slice length is wrong: ", len(slc))
	}
	if slc[0] != "1" || slc[1] != "2-5" || slc[2] != "7-9" || slc[3] != "11-12" {
		t.Fatal("Wrong values in generated slice. Values: ", slc)
	}
}

func TestIncrementingSliceWithoutRanges(t *testing.T) {
	input := []int{1,2,3,4,6,7,8,10,11}

	slc := []string{}

	for _, v := range input {
		slc = GenerateReferenceSlice(slc, v, false, true)
	}

	fmt.Println(slc)
	if len(slc) != 9 {
		t.Fatal("Resulting slice length is wrong: ", len(slc))
	}
	if slc[0] != "2" || slc[8] != "12" {
		t.Fatal("Wrong values in generated slice. Values: ", slc)
	}
}