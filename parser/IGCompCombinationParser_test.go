package parser

import (
	"IG-Parser/tree"
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

	if !node.Left.IsLeafNode() {
		t.Error("Left leaf node not detected.")
	}
	if !node.Right.IsLeafNode() {
		t.Error("Right leaf node not detected.")
	}

	if node.Left.Entry != "inspect and " {
		t.Error("Left leaf node has wrong value.")
	}

	if node.Right.Entry != " party" {
		t.Error("Right leaf node has wrong value.")
	}

	if node.CalculateDepth() != 1 {
		t.Error("Tree depth is wrong: " + strconv.Itoa(node.CalculateDepth()))
	}

	if node.CountLeaves() != 2 {
		t.Error("Tree leaf count is wrong: " + strconv.Itoa(node.CountLeaves()))
	}

}

func TestTwoLevelTree(t *testing.T) {
	input := "((inspect and [OR] party) [AND] sing)"

	// Create root node
	node := tree.Node{}
	// Parse provided expression
	parseDepth(input, &node)

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

	if node.CalculateDepth() != 6 {
		t.Error("Tree depth is wrong: " + strconv.Itoa(node.CalculateDepth()))
	}

	if node.CountLeaves() != 9 {
		t.Error("Tree leaf count is wrong: " + strconv.Itoa(node.CountLeaves()))
	}
}
