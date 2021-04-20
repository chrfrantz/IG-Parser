package exporter

import (
	"IG-Parser/tree"
	"testing"
)

/*
Tests basic flattening of two-dimensional array
 */
func TestFlatten(t *testing.T) {

	t.Fatal("Test to be completed - not yet done.")

	arr := make([][]*tree.Node, 4, 4)

	/*arr[0][0] = &tree.Node{Entry: "First"}
	arr[0][1] = &tree.Node{Entry: "Second"}
	arr[1][0] = &tree.Node{Entry: "Third"}
	arr[1][1] = &tree.Node{Entry: "Fourth"}
*/

	arr2 := flatten(arr)

	if len(arr2) != 4 ||
		arr2[0].Entry != "First" ||
		arr2[1].Entry != "Second" ||
		arr2[2].Entry != "Third" ||
		arr2[3].Entry != "Fourth" {
		t.Fatal("Array has not been correctly flattened")
	}

}

/*
Tests basic addition of element to array
 */
func TestAddingElement(t *testing.T) {

	arr := []string{"One", "Two"}

	arr = addElementIfNotExisting("Two", arr)

	if len(arr) != 2 {
		t.Error("Element should not have been added")
	}

	arr = addElementIfNotExisting("Three", arr)

	if len(arr) != 3 {
		t.Error("Element should have been added")
	}

}

/*
Test moving to last position, with and without addition during moving.
 */
func TestMoveElementToFirstPosition(t *testing.T) {

	arr := []string{"One", "Two", "Three"}

	arr = moveElementToFirstPosition("Three", arr, true)

	if len(arr) != 3 {
		t.Error("Element should not have been added")
	}

	if arr[0] != "Three" {
		t.Error("Element has not been moved to first position")
	}

	arr = moveElementToFirstPosition("Fourth", arr, true)

	if len(arr) != 4 {
		t.Error("Element should have been added")
	}

	if arr[0] != "Fourth" {
		t.Error("Element has not been moved to first position")
	}

	arr = moveElementToFirstPosition("Fifth", arr, false)

	if len(arr) != 4 {
		t.Error("Element should have been added")
	}

}

/*
Test moving to last position, with and without addition during moving.
 */
func TestMoveElementToLastPosition(t *testing.T) {

	arr := []string{"One", "Two", "Three"}

	arr = moveElementToLastPosition("Two", arr, true)

	if len(arr) != 3 {
		t.Error("Element should not have been added")
	}

	if arr[2] != "Two" {
		t.Error("Element has not been moved to last position")
	}

	arr = moveElementToLastPosition("Fourth", arr, true)

	if len(arr) != 4 {
		t.Error("Element should have been added")
	}

	if arr[3] != "Fourth" {
		t.Error("Element has not been moved to last position")
	}

	arr = moveElementToLastPosition("Two", arr, true)

	if len(arr) != 4 {
		t.Error("Element should have been added")
	}

	if arr[3] != "Two" {
		t.Error("Element has not been moved to last position")
	}

	arr = moveElementToLastPosition("Fifth", arr, false)

	if len(arr) != 4 {
		t.Error("Element should not have been added")
	}

	if arr[0] != "One" || arr[1] != "Three" || arr[2] != "Fourth" || arr[3] != "Two" {
		t.Error("Array has change, even though it shouldn't have")
	}

}

