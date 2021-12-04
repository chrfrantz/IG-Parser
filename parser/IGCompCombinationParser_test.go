package parser

import (
	"IG-Parser/tree"
	"fmt"
	"strconv"
	"testing"
)

func TestNonCombination(t *testing.T) {

	input := "(inspect and )"

	node, _, err := ParseIntoNodeTree(input, false, "(", ")")

	fmt.Println(node.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR && err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
		t.Fatal("Parsing throws error where there should be none. Error: ", err.Error())
	}

	if !node.IsLeafNode() {
		t.Fatal("Node should be leaf node")
	}

	// TODO: What should it return? Probably raw text
	/*if node.Entry != input {
		t.Fatal("Leaf node entry should be filled with non-combination text")
	}*/

}

func TestBasicExpression(t *testing.T) {

	input := "(inspect and [OR] party)"

	// Parse provided expression
	node, _, err := ParseIntoNodeTree(input, false, "(", ")")

	fmt.Println("Returned node: " + node.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing throws error where there should be none.")
	}

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

func TestMissingParenthesesAndLogicalOperator(t *testing.T) {

	input := "left [AND] right"

	// Parse provided expression
	_, _, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_ERROR_LOGICAL_OPERATOR_OUTSIDE_COMBINATION {
		t.Fatal("Parsing does not pick up on invalid operator use in " + input)
	}
}

func TestTwoLevelTree(t *testing.T) {

	input := "((inspect and [OR] party) [AND] sing)"

	// Parse provided expression
	node, _, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing throws error where there should be none.")
	}

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

	// Parse provided expression
	node, _, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing throws error where there should be none.")
	}

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

	// Parse provided expression
	node, _, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing throws error where there should be none.")
	}

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

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

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

/*
Tests the parser's ability to handle multiple OR operators on the same level
*/
func TestAutomatedOrExpansion(t *testing.T) {

	input := "(Left side information [OR] middle information and [OR] right-side)"

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing throws error where there should be none.")
	}

	// Test return information from parsing
	if modified != "((Left side information [OR] middle information and) [OR] right-side)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	// Test reconstruction from tree
	if node.Stringify() != "((Left side information [OR] middle information and) [OR] right-side)" {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}
}

/*
Tests the parser's ability to handle multiple XOR operators on the same level
*/
func TestAutomatedXorExpansion(t *testing.T) {

	input := "(Left side information [XOR] middle information and [XOR] right-side)"

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing throws error where there should be none.")
	}

	// Test return information from parsing
	if modified != "((Left side information [XOR] middle information and) [XOR] right-side)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	// Test reconstruction from tree
	if node.Stringify() != "((Left side information [XOR] middle information and) [XOR] right-side)" {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}
}

/*
Tests the parser's ability to prevent mixed operators on the same level
*/
func TestMixedOperatorsOnSameLevel(t *testing.T) {

	input := "(Left side information [AND] middle information and [OR] right-side)"

	// Parse provided expression
	_, _, err := ParseIntoNodeTree(input, false, "(", ")")

	// Should fail based on mixed operators
	if err.ErrorCode != tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS {
		t.Fatal("Parsing throws error where there should be none. Error: ", err.ErrorCode)
	}

}

func TestNonCombinationParentheses(t *testing.T) {

	input := "(Left side information (source) [AND] middle information and [AND] right-side)"

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

	fmt.Println(node.String())

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

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

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
Tests parsing of shared elements in surrounding non-combination parentheses and inclusion in resulting node
 */
func TestSharedElements(t *testing.T) {

	input := "( shared left (Left side information [XOR] middle information) shared right)"

	tree.SHARED_ELEMENT_INHERITANCE_MODE = tree.SHARED_ELEMENT_INHERIT_NOTHING

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Shared elements, e.g., '(left shared (left [AND] right) right shared)', should not produce error ", err.Error())
	}

	// Test return information from parsing
	if modified != "( shared left (Left side information [XOR] middle information) shared right)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	if node.GetSharedLeft()[0] != "shared left" {
		t.Fatal("Parsed left shared value is not correct. Output: '" + fmt.Sprint(node.GetSharedLeft()) + "'")
	}

	if node.GetSharedRight()[0] != "shared right" {
		t.Fatal("Parsed right shared value is not correct. Output: '" + fmt.Sprint(node.GetSharedRight()) + "'")
	}

	// Test reconstruction from tree
	if node.Stringify() != "(shared left (Left side information [XOR] middle information) shared right)" {
		t.Fatal("Stringified output does not correspond to input (Output: " + node.Stringify() + "')")
	}
}

func TestInferenceOfLowerLevelSharedElements(t *testing.T) {


	text := "Cex(shared left1 (left1 [XOR] left2) mid (shared left2 (lefter [XOR] (left [AND] right)) shared right2) shared right1)"

	res1, res2 := ParseStatement(text)

	//res1, res2, res3 := ParseIntoNodeTree(text, false, "(", ")")

	fmt.Println(res1)
	fmt.Println(res2)
	//fmt.Println(res3)
}

/*
Tests for correct parsing of shared elements as well as decomposition of multiple AND combinations
 */
func TestSharedElementsAndAndCombinationWithInheritance(t *testing.T) {

	input := "( shared left (Left side information [AND] middle information [AND] right information) shared right)"

	tree.SHARED_ELEMENT_INHERITANCE_MODE = tree.SHARED_ELEMENT_INHERIT_OVERRIDE

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Shared elements, e.g., '(left shared (left [AND] right) right shared)', should not produce error ", err.Error())
	}

	// Test return information from parsing (strips shared elements)
	if modified != "( shared left ((Left side information [AND] middle information) [AND] right information) shared right)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	if len(node.GetSharedLeft()) == 0 || node.GetSharedLeft()[0] != "shared left" {
		t.Fatal("Parsed left shared value is not correct. Output: " + fmt.Sprint(node.GetSharedLeft()))
	}

	if len(node.GetSharedRight()) == 0 || node.GetSharedRight()[0] != "shared right" {
		t.Fatal("Parsed right shared value is not correct. Output: " + fmt.Sprint(node.GetSharedRight()))
	}

	if len(node.Left.GetSharedLeft()) == 0 || node.Left.GetSharedLeft()[0] != "shared left" {
		t.Fatal("Left-nested left node did not inherit shared value. Node value: " + fmt.Sprint(node.Left.GetSharedLeft()) + ". Expected output: " + fmt.Sprint(node.GetSharedLeft()))
	}

	if len(node.Left.GetSharedRight()) == 0 || node.Left.GetSharedRight()[0] != "shared right" {
		t.Fatal("Left-nested right node did not inherit shared value. Node value: " + fmt.Sprint(node.Left.GetSharedRight()) + ". Expected output: " + fmt.Sprint(node.GetSharedRight()))
	}

	// Test reconstruction from tree
	if node.Stringify() != "(shared left ((shared left (Left side information [AND] middle information) shared right) [AND] right information) shared right)" {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}
}

/*
Tests for correct parsing of shared elements as well as decomposition of multiple AND combinations in APPEND mode
*/
func TestSharedElementsAndAndCombinationWithInheritanceAppendMode(t *testing.T) {

	input := "( shared left ( inner left (innermost left ((left-most information [AND] Left side information [AND] middle information) [AND] right information) shared right)))"

	tree.SHARED_ELEMENT_INHERITANCE_MODE = tree.SHARED_ELEMENT_INHERIT_APPEND

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

	fmt.Println(node.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Shared elements, e.g., '(left shared (left [AND] right) right shared)', should not produce error ", err.Error())
	}

	// Test return information from parsing
	if modified != "( shared left ( inner left (innermost left (((left-most information [AND] Left side information) [AND] middle information) [AND] right information) shared right)))" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	if fmt.Sprint(node.GetSharedLeft()) != "[innermost left]" {
		t.Fatal("Parsed left shared value is not correct. Node value: " + fmt.Sprint(node.GetSharedLeft()) + ". Expected output: [innermost left]")
	}

	if fmt.Sprint(node.GetSharedRight()) != "[shared right]" {
		t.Fatal("Parsed right shared value is not correct. Node value: " + fmt.Sprint(node.GetSharedRight()) + ". Expected output: [shared right]")
	}

	if fmt.Sprint(node.Left.GetSharedLeft()) != "[innermost left]" {
		t.Fatal("Left-nested left node did not inherit shared value. Node value: " + fmt.Sprint(node.Left.GetSharedLeft()) + ". Expected output: [innermost left]")
	}

	if fmt.Sprint(node.Left.GetSharedRight()) != "[shared right]" {
		t.Fatal("Left-nested right node did not inherit shared value. Node value: " + fmt.Sprint(node.Left.GetSharedRight()) + ". Expected output: shared right")
	}

	fmt.Println(node)

	// Test reconstruction from tree
	if node.Stringify() != "(innermost left ((innermost left ((innermost left (left-most information [AND] Left side information) shared right) [AND] middle information) shared right) [AND] right information) shared right)" {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}
}

/*
Tests for correct parsing of shared elements as well as decomposition of multiple AND combinations in OVERRIDE mode
*/
func TestSharedElementsAndAndCombinationWithInheritanceOverrideMode(t *testing.T) {

	input := "( shared left ( ( innermost left (Far left side [AND] Left side information [AND] inner right)) [AND] right information) shared right)"

	tree.SHARED_ELEMENT_INHERITANCE_MODE = tree.SHARED_ELEMENT_INHERIT_OVERRIDE

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

	fmt.Println(node.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Shared elements, e.g., '(left shared (left [AND] right) right shared)', should not produce error ", err.Error())
	}

	// Test return information from parsing (strips shared elements)
	if modified != "( shared left ( ( innermost left ((Far left side [AND] Left side information) [AND] inner right)) [AND] right information) shared right)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	if len(node.GetSharedLeft()) == 0 || node.GetSharedLeft()[0] != "shared left" {
		t.Fatal("Parsed left shared value is not correct. Node value: " + fmt.Sprint(node.GetSharedLeft()) + ". Expected output: shared left,inner left")
	}

	if len(node.GetSharedRight()) == 0 || node.GetSharedRight()[0] != "shared right" {
		t.Fatal("Parsed right shared value is not correct. Node value: " + fmt.Sprint(node.GetSharedRight()) + ". Expected output: shared right")
	}

	if fmt.Sprint(node.Left.GetSharedLeft()) != "[shared left]" {
		t.Fatal("Left-nested left node did not inherit shared value. Node value: " + fmt.Sprint(node.Left.GetSharedLeft()) + ". Expected output: [shared left]")
	}

	if len(node.Left.GetSharedRight()) == 0 || node.Left.GetSharedRight()[0] != "shared right" {
		t.Fatal("Right-nested right node did not inherit shared value. Node value: " + fmt.Sprint(node.Left.GetSharedRight()) + ". Expected output: shared right")
	}

	// Test reconstruction from tree
	if node.Stringify() != "(shared left ((shared left ((shared left innermost left (Far left side [AND] Left side information) shared right) [AND] inner right) shared right) [AND] right information) shared right)" {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}
}

/*
Tests whether missing specification of logical operator between simple string and combination is picked up.
 */
func TestSharedElementsAndAndCombinationWithMissingCombination(t *testing.T) {

	input := "( shared left ( inner left (Far left side [AND] Left side information)) [AND] inner right [AND] right information)"

	tree.SHARED_ELEMENT_INHERITANCE_MODE = tree.SHARED_ELEMENT_INHERIT_APPEND

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")


	fmt.Println(node.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parser has produced error. Error: ", err.Error())
	}

	// Test return information from parsing (strips shared elements)
	if modified != "(( shared left ( inner left (Far left side [AND] Left side information)) [AND] inner right) [AND] right information)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	if len(node.GetSharedLeft()) != 0 || fmt.Sprint(node.GetSharedLeft()) != "[]" {
		t.Fatal("Parsed left shared value is not correct. Node value: " + fmt.Sprint(node.GetSharedLeft()) + ". Expected output: ")
	}

	if len(node.GetSharedRight()) != 0 {
		t.Fatal("Parsed right shared value is not correct. Node value: " + fmt.Sprint(node.GetSharedRight()) + ". Expected empty output")
	}

	if fmt.Sprint(node.Left.Left.GetSharedLeft()) != "[inner left]" {
		t.Fatal("Left-nested left node did not inherit shared value. Node value: " + fmt.Sprint(node.Left.Left.GetSharedLeft()) + ". Expected output: [inner left]")
	}

	// Test reconstruction from tree
	if node.Stringify() != "(((inner left (Far left side [AND] Left side information) [AND] inner right) [AND] right information)" {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}
}

/*
Tests for correct parsing of shared elements as well as decomposition of multiple AND combinations
*/
func TestSharedElementsAndAndCombinationWithoutInheritance(t *testing.T) {

	tree.SHARED_ELEMENT_INHERITANCE_MODE = tree.SHARED_ELEMENT_INHERIT_NOTHING

	input := "( shared left (Left side information [AND] middle information [AND] right information) shared right)"

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Shared elements, e.g., '(left shared (left [AND] right) right shared)', should not produce error ", err.Error())
	}

	// Test return information from parsing (strips shared elements)
	if modified != "( shared left ((Left side information [AND] middle information) [AND] right information) shared right)" {
		t.Fatal("Modified output does not correspond to input (Output: '" + modified + "')")
	}

	if node.GetSharedLeft()[0] != "shared left" {
		t.Fatal("Parsed left shared value is not correct. Output: " + fmt.Sprint(node.GetSharedLeft()))
	}

	if node.GetSharedRight()[0] != "shared right" {
		t.Fatal("Parsed right shared value is not correct. Output: " + fmt.Sprint(node.GetSharedRight()))
	}

	if len(node.Left.GetSharedLeft()) > 1 || node.Left.GetSharedLeft()[0] != "" {
		t.Fatal("Left-nested left node should not inherit shared value. Node value: " + fmt.Sprint(node.Left.GetSharedLeft()) + ". Expected output: ")
	}

	if len(node.Left.GetSharedRight()) > 1 || node.Left.GetSharedRight()[0] != "" {
		t.Fatal("Left-nested right node should not inherit shared value. Node value: " + fmt.Sprint(node.Left.GetSharedRight()) + ". Expected output: ")
	}

	// Test reconstruction from tree
	if node.Stringify() != "(shared left ((Left side information [AND] middle information) [AND] right information) shared right)" {
		t.Fatal("Stringified output does not correspond to input (Output: '" + node.Stringify() + "')")
	}
}

/*
Tests whether inline annotations using the [text] notation are ignored
 */
func TestInlineAnnotations(t *testing.T) {
	input := "(Left side information [source] [AND] middle information and [AND] right-side)"

	// Parse provided expression
	node, modified, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Inline annotations (e.g., [dgjslkgjsø]) should not produce error ", err.Error())
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

	// Parse provided expression
	_, text, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS {
		t.Fatal("Did not pick up on invalid logical operator combinations on given level")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

func TestAdjacentAndOperators(t *testing.T) {

	input := "(Left side information [AND] [AND] right-side)"

	// Parse provided expression
	_, text, err := ParseIntoNodeTree(input, false, "(", ")")

	fmt.Println(err.Error())

	if err.ErrorCode != tree.PARSING_INVALID_COMBINATION {
		t.Fatal("Did not pick up on invalid combination expression")
	}

	fmt.Println(text)

	if text != "((Left side information [AND]) [AND] right-side)" {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

func TestAdjacentNonAndOperators(t *testing.T) {

	input := "(Left side information [OR] [XOR] right-side)"

	// Parse provided expression
	_, text, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS {
		t.Fatal("Did not pick up on invalid logical operator combinations on given level")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

func TestAdjacentAndAndNonAndOperators(t *testing.T) {

	input := "(Left side information [OR] [AND] right-side)"

	// Parse provided expression
	_, text, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_ERROR_INVALID_OPERATOR_COMBINATIONS {
		t.Fatal("Did not pick up on invalid logical operator combinations on given level")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

	input = "(Left side information [AND] [OR] right-side [XOR] something)"

	// Parse provided expression
	_, text, err = ParseIntoNodeTree(input, false, "(", ")")

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

	// Parse provided expression
	_, text, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_ERROR_EMPTY_LEAF {
		t.Fatal("Did not pick up on empty leaf value")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

	// Test empty right side (without whitespace)
	input = "(Left side information [OR])"

	// Parse provided expression
	_, text, err = ParseIntoNodeTree(input, false, "(", ")")

	fmt.Println(err.Error())

	if err.ErrorCode != tree.PARSING_INVALID_COMBINATION {
		t.Fatal("Did not pick up on invalid combinations")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

	// Test left side (including whitespace)
	input = "( [OR] right side info )"

	// Parse provided expression
	_, text, err = ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_ERROR_EMPTY_LEAF {
		t.Fatal("Did not pick up on empty leaf value")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

	// Test left side (without whitespace)
	input = "([OR] right side info )"

	// Parse provided expression
	_, text, err = ParseIntoNodeTree(input, false, "(", ")")

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

	// Parse provided expression
	_, text, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_ERROR_IMBALANCED_PARENTHESES {
		t.Fatal("Did not pick up on imbalanced parentheses")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

	input = "((( Left side information [AND] middle information ))"

	// Parse provided expression
	_, text, err = ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_ERROR_IMBALANCED_PARENTHESES {
		t.Fatal("Did not pick up on imbalanced parentheses")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

/*
Tests for statements that contain parentheses, but no combinations.
 */
func TestSimpleStatementWithEmbeddedParentheses(t *testing.T) {

	input := "Left side informat(ion middle information ) sd glkjdslk ) dksljfs("

	// Parse provided expression
	_, text, err := ParseIntoNodeTree(input, false, "(", ")")

	if err.ErrorCode != tree.PARSING_NO_COMBINATIONS {
		t.Fatal("Did not pick up on simple entry")
	}

	if text != input {
		t.Fatal("Returned output does not correspond to input (Output: '" + text + "')")
	}

}

