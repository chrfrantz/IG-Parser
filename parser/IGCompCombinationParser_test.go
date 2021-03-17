package parser

import (
	"IG-Parser/tree"
	"fmt"
	"strconv"
	"testing"
)


func TestBasicExpression(t *testing.T) {

	input := "(inspect and [OR] party)"

	// Create root node
	node := tree.Node{}

	// Sanity check on depth and breadth
	if node.CalculateDepth() != 0 {
		t.Error("Node depth is incorrect (empty node)")
	}
	if node.CountLeaves() != 0 {
		t.Error("Node leaf count is incorrect (empty node)")
	}

	// Parse provided expression
	parseDepth(input, &node)

	if node.IsLeafNode() {
		t.Fatal("Node should not be leaf node")
	}

	if !node.Left.IsLeafNode() {
		t.Error("Left leaf node not detected.")
	}
	if !node.Right.IsLeafNode() {
		t.Error("Right leaf node not detected.")
	}

	if node.Left.Entry != "inspect and" {
		t.Error("Left leaf node has wrong value.")
	}

	if node.Right.Entry != "party" {
		t.Error("Right leaf node has wrong value.")
	}

	if node.CalculateDepth() != 1 {
		t.Error("Tree depth is wrong: " + strconv.Itoa(node.CalculateDepth()))
	}

	if node.CountLeaves() != 2 {
		t.Error("Tree leaf count is wrong: " + strconv.Itoa(node.CountLeaves()))
	}

	if node.Stringify() != input {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}

}

func TestTwoLevelTree(t *testing.T) {
	input := "((inspect and [OR] party) [AND] sing)"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	parseDepth(input, &node)

	if node.IsLeafNode() {
		t.Fatal("Node should not be leaf node")
	}

	if !(node.LogicalOperator == "AND" &&
		node.Left.LogicalOperator == "OR" &&
		node.Left.Left.Entry == "inspect and" &&
		node.Left.Right.Entry == "party" &&
		node.Right.Entry == "sing") {
		t.Fatal("Tree content is not correct.")
	}

	if node.CalculateDepth() != 2 {
		t.Error("Tree depth is wrong: " + strconv.Itoa(node.CalculateDepth()))
	}

	if node.CountLeaves() != 3 {
		t.Error("Tree leaf count is wrong: " + strconv.Itoa(node.CountLeaves()))
	}
}

func TestComplexExpression(t *testing.T) {
	input := "((((inspect and [OR] party) [AND] ((review [XOR] muse) [AND] pray))))"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	parseDepth(input, &node)

	if !(node.LogicalOperator == "AND" &&
		node.Left.Left.Entry == "inspect and" &&
		node.Left.LogicalOperator == "OR" &&
		node.Left.Right.Entry == "party" &&
		node.Right.Left.Left.Entry == "review" &&
		node.Right.Left.LogicalOperator == "XOR" &&
		node.Right.Left.Right.Entry == "muse" &&
		node.Right.LogicalOperator == "AND" &&
		node.Right.Right.Entry == "pray") {
		t.Fatal("Tree content is not correct.")
	}

	if node.CalculateDepth() != 3 {
		t.Error("Tree depth is wrong: " + strconv.Itoa(node.CalculateDepth()))
	}

	if node.CountLeaves() != 5 {
		t.Error("Tree leaf count is wrong: " + strconv.Itoa(node.CountLeaves()))
	}
}

func TestDeepExpression(t *testing.T) {
	input := "((((inspect and [OR] (party [OR] hoard)) [AND] (((review [AND] (establish [XOR] (identify [AND] detect something))) [XOR] muse) [AND] pray))))"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	parseDepth(input, &node)

	if node.Stringify() != input[2:len(input)-2] {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}

	if node.CalculateDepth() != 6 {
		t.Error("Tree depth is wrong: " + strconv.Itoa(node.CalculateDepth()))
	}

	if node.CountLeaves() != 9 {
		t.Error("Tree leaf count is wrong: " + strconv.Itoa(node.CountLeaves()))
	}

}

/*
Tests the parser's ability to handle multiple AND operators on the same level
 */
func TestAutomatedAndExpansion(t *testing.T) {
	input := "(Left side information [AND] middle information and [AND] right-side)"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, modified, err := parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing throws error where there should be none.")
	}

	// Test return information from parsing
	if modified != "((Left side information [AND] middle information and) [AND] right-side)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	// Test reconstruction from tree
	if node.Stringify() != "((Left side information [AND] middle information and) [AND] right-side)" {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}
}

func TestNonCombinationParentheses(t *testing.T) {
	input := "(Left side information (source) [AND] middle information and [AND] right-side)"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, modified, err := parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Inline parentheses without embedded combinations (e.g., (dgjslkgjsø)) should not produce error")
	}

	// Test return information from parsing
	if modified != "((Left side information (source) [AND] middle information and) [AND] right-side)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	// Test reconstruction from tree
	if node.Stringify() != "((Left side information (source) [AND] middle information and) [AND] right-side)" {
		t.Fatal("Stringified output does not correspond to input (Output: " + node.Stringify() + "')")
	}
}

func TestMultipleNonCombinationParentheses(t *testing.T) {
	input := "(Left side (information) source [AND] (middle information and) [AND] right-side)"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, modified, err := parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Inline parentheses without embedded combinations (e.g., (dgjslkgjsø)) should not produce error")
	}

	if !(node.LogicalOperator == "AND" &&
		node.Left.LogicalOperator == "AND" &&
		node.Left.Left.Entry == "Left side (information) source" &&
		node.Left.Right.Entry == "(middle information and)" &&
		node.Right.Entry == "right-side") {
		t.Fatal("Returned node tree does not correspond to input")
	}

	// Test return information from parsing
	if modified != "((Left side (information) source [AND] (middle information and)) [AND] right-side)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	// Test reconstruction from tree
	if node.Stringify() != "((Left side (information) source [AND] (middle information and)) [AND] right-side)" {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}
}

/*
Tests whether inline annotations using the [text] notation are ignored
 */
func TestInlineAnnotations(t *testing.T) {
	input := "(Left side information [source] [AND] middle information and [AND] right-side)"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, modified, err := parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Inline annotations (e.g., [dgjslkgjsø]) should not produce error " + err.Error())
	}

	// Test return information from parsing
	if modified != "((Left side information [source] [AND] middle information and) [AND] right-side)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	// Test reconstruction from tree
	if node.Stringify() != "((Left side information [source] [AND] middle information and) [AND] right-side)" {
		t.Fatal("Stringified output does not correspond to input (Output: " + node.Stringify() + "')")
	}
}

/*
Tests that the parser captures non-AND operators on the same level
 */
func TestCombinedOperators(t *testing.T) {

	input := "(Left side information [AND] middle information and [XOR] right-side)"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, text, err := parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS {
		t.Fatal("Did not pick up on invalid logical operator combinations on given level")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

func TestAdjacentAndOperators(t *testing.T) {

	input := "(Left side information [AND] [AND] right-side)"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, text, err := parseDepth(input, &node)

	fmt.Println(err.Error())

	if err.ErrorCode != tree.PARSING_INVALID_COMBINATION {
		t.Fatal("Did not pick up on invalid combination expression")
	}

	fmt.Println(text)

	if text != "((Left side information [AND])[AND] right-side)" {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

func TestAdjacentNonAndOperators(t *testing.T) {

	input := "(Left side information [OR] [XOR] right-side)"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, text, err := parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS {
		t.Fatal("Did not pick up on invalid logical operator combinations on given level")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

func TestAdjacentAndAndNonAndOperators(t *testing.T) {

	input := "(Left side information [OR] [AND] right-side)"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, text, err := parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS {
		t.Fatal("Did not pick up on invalid logical operator combinations on given level")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

	input = "(Left side information [AND] [OR] right-side [XOR] something)"

	// Create root node
	node = tree.Node{}
	// Parse provided expression
	_, text, err = parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS {
		t.Fatal("Did not pick up on invalid logical operator combinations on given level")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

func TestIncompleteExpression(t *testing.T) {

	// Test empty right side (including whitespace)
	input := "(Left side information [OR] )"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, text, err := parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_ERROR_EMPTY_LEAF {
		t.Fatal("Did not pick up on empty leaf value")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

	// Test empty right side (without whitespace)
	input = "(Left side information [OR])"

	// Create root node
	node = tree.Node{}
	// Parse provided expression
	_, text, err = parseDepth(input, &node)

	fmt.Println(err.Error())

	if err.ErrorCode != tree.PARSING_INVALID_COMBINATION {
		t.Fatal("Did not pick up on invalid combinations")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

	// Test left side (including whitespace)
	input = "( [OR] right side info )"

	// Create root node
	node = tree.Node{}
	// Parse provided expression
	_, text, err = parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_ERROR_EMPTY_LEAF {
		t.Fatal("Did not pick up on empty leaf value")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

	// Test left side (without whitespace)
	input = "([OR] right side info )"

	// Create root node
	node = tree.Node{}
	// Parse provided expression
	_, text, err = parseDepth(input, &node)

	fmt.Println(err.Error())

	if err.ErrorCode != tree.PARSING_INVALID_COMBINATION {
		t.Fatal("Did not pick up on empty leaf value")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

func TestUnbalancedParentheses(t *testing.T) {

	input := "(Left side information [AND] middle information ))"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	_, text, err := parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_ERROR_IMBALANCED_PARENTHESES {
		t.Fatal("Did not pick up on imbalanced parentheses")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

	input = "((( Left side information [AND] middle information ))"

	// Create root node
	node = tree.Node{}
	// Parse provided expression
	_, text, err = parseDepth(input, &node)

	if err.ErrorCode != tree.PARSING_ERROR_IMBALANCED_PARENTHESES {
		t.Fatal("Did not pick up on imbalanced parentheses")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

