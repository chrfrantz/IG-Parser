package tree

import (
	"fmt"
	"os"
	"testing"
)


func TestTreeCreation(t *testing.T) {

	root := Node{}

	leftNode := Node{}
	leftNode.LogicalOperator = "123"
	res, err := root.InsertLeftNode(&leftNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of left node should not have failed")
	}

	rightNode := Node{}
	rightNode.LogicalOperator = "456"
	res, err = root.InsertRightNode(&rightNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of right node should not have failed")
	}

	leftLeftChildNode := Node{Entry: "Left left data"}
	leftRightChildNode := Node{Entry: "Left right data"}

	res, err = leftNode.InsertLeftNode(&leftLeftChildNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of left left node should not have failed")
	}
	res, err = leftNode.InsertRightNode(&leftRightChildNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of left right node should not have failed")
	}

	rightLeftChildNode := Node{Entry: "Right left data"}
	rightRightChildNode := Node{Entry: "Right right data"}

	res, err = rightNode.InsertLeftNode(&rightLeftChildNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of right left node should not have failed")
	}
	res, err = rightNode.InsertRightNode(&rightRightChildNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of right right node should not have failed")
	}

	// Check for values and correct parent assignment
	if rightRightChildNode.Parent != &rightNode {
		t.Fatal("Parent is not correctly assigned on right right child")
	}

	if rightLeftChildNode.Parent != &rightNode {
		t.Fatal("Parent is not correctly assigned on right left child")
	}

	if rightRightChildNode.Parent.Parent != &root {
		t.Fatal("Grandparent is not correctly assigned on right right child")
	}

	if rightLeftChildNode.Parent.Parent != &root {
		t.Fatal("Grandparent is not correctly assigned on right right child")
	}

	if leftNode.LogicalOperator != "123" {
		t.Fatal("Logical operator information is incorrect")
	}

	if rightNode.LogicalOperator != "456" {
		t.Fatal("Logical operator information is incorred")
	}

	// Should have a single branch with four leaves, i.e., leaves[1][4]
	leaves := root.GetLeafNodes()

	if len(leaves) != 1 || len(leaves[0]) != 4 {
		t.Fatal("Wrong number of leaves returned")
	}

	if leftLeftChildNode.CountParents() != 2 {
		t.Fatal("Wrong parent count")
	}

	if leftRightChildNode.CountParents() != 2 {
		t.Fatal("Wrong parent count")
	}

	if rightLeftChildNode.CountParents() != 2 {
		t.Fatal("Wrong parent count")
	}

	if rightRightChildNode.CountParents() != 2 {
		t.Fatal("Wrong parent count")
	}

	if root.CalculateDepth() != 2 {
		t.Fatal("Depth of tree is wrong")
	}

	if root.CountLeaves() != 4 {
		t.Fatal("Calculation of leaf nodes is wrong")
	}

}

func TestRepeatedNodeAssignment(t *testing.T) {

	root := Node{}

	res, err := root.assignParent(&root)
	if res || err.ErrorCode != TREE_INVALID_NODE_SELF_LINKAGE {
		t.Fatal("Assignment to self as parent should fail")
	}
	res, err = root.InsertLeftNode(&root)
	if res || err.ErrorCode != TREE_INVALID_NODE_SELF_LINKAGE {
		t.Fatal("Assignment to self to left child should fail")
	}
	res, err = root.InsertRightNode(&root)
	if res || err.ErrorCode != TREE_INVALID_NODE_SELF_LINKAGE {
		t.Fatal("Assignment to self to right child should fail")
	}

	res, err = root.InsertLeftLeaf("Test")
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Adding left leaf failed")
	}
	res, err = root.InsertLeftLeaf("Test3")
	if res || err.ErrorCode != TREE_INVALID_NODE_ADDITION {
		t.Fatal("Adding left leaf should fail")
	}
	res, err = root.InsertRightLeaf("Test")
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Adding right leaf failed")
	}
	res, err = root.InsertRightLeaf("Test3")
	if res || err.ErrorCode != TREE_INVALID_NODE_ADDITION {
		t.Fatal("Adding right leaf should fail")
	}
}

/*
Tests the distance function and retrieval of logical operators in between
 */
func TestNodeDistanceSearch(t *testing.T) {
	root := Node{}
	root.LogicalOperator = "789"

	leftNode := Node{}
	leftNode.LogicalOperator = "123"
	res, err := root.InsertLeftNode(&leftNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of left node should not have failed")
	}

	rightNode := Node{}
	rightNode.LogicalOperator = "456"
	res, err = root.InsertRightNode(&rightNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of right node should not have failed")
	}

	leftLeftChildNode := Node{Entry: "Left left data"}
	leftRightChildNode := Node{Entry: "Left right data"}

	res, err = leftNode.InsertLeftNode(&leftLeftChildNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of left left node should not have failed")
	}
	res, err = leftNode.InsertRightNode(&leftRightChildNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of left right node should not have failed")
	}

	rightLeftChildNode := Node{Entry: "Right left data"}
	rightRightChildNode := Node{Entry: "Right right data"}

	res, err = rightNode.InsertLeftNode(&rightLeftChildNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of right left node should not have failed")
	}
	res, err = rightNode.InsertRightNode(&rightRightChildNode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Insertion of right right node should not have failed")
	}

	// Left branch
	// Simple left to right
	res, ops, err := FindLogicalLinkage(&leftLeftChildNode, &leftRightChildNode, nil)

	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}

	fmt.Println(ops)
	if len(ops) != 1 || ops[0] != "123" {
		t.Fatal("Logical operators are incorrectly determined.")
	}

	// Simple right to left
	res, ops, err = FindLogicalLinkage(&leftRightChildNode, &leftLeftChildNode, nil)

	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}
	if len(ops) != 1 || ops[0] != "123" {
		t.Fatal("Logical operators are incorrectly determined.")
	}

	// Right branch
	// Simple left to right
	res, ops, err = FindLogicalLinkage(&rightLeftChildNode, &rightRightChildNode, nil)

	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}
	if len(ops) != 1 || ops[0] != "456" {
		t.Fatal("Logical operators are incorrectly determined.")
	}

	// Simple right to left
	res, ops, err = FindLogicalLinkage(&rightRightChildNode, &rightLeftChildNode, nil)

	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}
	if len(ops) != 1 || ops[0] != "456" {
		t.Fatal("Logical operators are incorrectly determined.")
	}

	// Across branches
	res, ops, err = FindLogicalLinkage(&leftLeftChildNode, &rightRightChildNode, nil)

	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}

	fmt.Println(ops)

	if len(ops) != 3 || ops[0] != "123" || ops[1] != "789" || ops[2] != "456" {
		t.Fatal("Logical operators are incorrectly determined.")
	}

	// Test unbalanced branches

	subnode := Node{}

	leftSubnode := Node{Entry: "left sub"}
	rightSubnode := Node{Entry: "right sub"}

	res, err = subnode.InsertLeftNode(&leftSubnode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when adding new node should not happen")
	}

	res, err = subnode.InsertRightNode(&rightSubnode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when adding new node should not happen")
	}

	res, err = leftRightChildNode.InsertRightNode(&subnode)
	if res || err.ErrorCode != TREE_INVALID_NODE_ADDITION {
		t.Fatal("Should pick up on invalid addition of node into existing node")
	}

	// Manually clean node
	leftRightChildNode.Entry = ""
	res, err = leftRightChildNode.InsertRightNode(&subnode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Addition of node to empty node should work. Error: ", err)
	}

	// Assessment should fail due to missing logical operator and left leaf value
	res, err = leftRightChildNode.Validate()
	if res || err.ErrorCode != TREE_INVALID_TREE {
		t.Fatal("Reconfigured tree should throw problem. Error: ", err)
	}

	// Manually add operator (should still fail)
	leftRightChildNode.LogicalOperator = "AND"

	// Node-level validation should fail
	res, err = leftRightChildNode.Validate()
	if res || err.ErrorCode != TREE_INVALID_TREE {
		t.Fatal("Reconfigured tree should still throw problem. Error: ", err)
	}

	// Now validate entire tree - which should fail
	res, err = root.Validate()
	if res || err.ErrorCode != TREE_INVALID_TREE {
		t.Fatal("Validation of reconfigured tree from root should throw problem. Error: ", err)
	}

	// Detect logical operators
	res, ops, err = FindLogicalLinkage(&leftLeftChildNode, &rightSubnode, nil)
	if res || err.ErrorCode != TREE_INPUT_VALIDATION {
		t.Fatal("Link between nodes should not be found because of empty nodes.")
	}

	// Fix node by adding left leaf
	res, err = leftRightChildNode.InsertLeftLeaf("left sub")
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Reconfigured tree by adding left leaf should not throw problem. Error: ", err)
	}

	// Should still fail, because of missing operator on right nested node
	res, err = leftRightChildNode.Validate()
	if res || err.ErrorCode != TREE_INVALID_TREE {
		t.Fatal("Reconfigured tree should throw problem on validation due to missing logical operator. Error: ", err)
	}

	// Adding operator should resolve this
	leftRightChildNode.Right.LogicalOperator = "OR"
	res, err = leftRightChildNode.Validate()
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Reconfigured tree should not throw problem. Error: ", err)
	}

	// Now validate entire tree again - should not fail
	res, err = root.Validate()
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Reconfigured tree should not throw problem. Error: ", err)
	}

	res, ops, err = FindLogicalLinkage(&leftLeftChildNode, &rightSubnode, nil)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}
	if len(ops) != 3 || ops[0] != "123" || ops[1] != "AND" || ops[2] != "OR" {
		t.Fatal("Logical operators are incorrectly detected.")
	}

	os.Exit(0)
	res, ops, _ = FindLogicalLinkage(&leftLeftChildNode, &rightRightChildNode, nil)
	if !res {
		t.Fatal("Link between nodes could not be found.")
	}
	fmt.Println(ops)

}

// Test combination of nodes
