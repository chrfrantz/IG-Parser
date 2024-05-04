package tree

import (
	"fmt"
	"testing"
)

func TestNode_IsEmptyNode(t *testing.T) {
	node := Node{}

	if !node.IsEmptyOrNilNode() {
		t.Fatal("Node should be considered empty")
	}

	if node.IsNil() {
		t.Fatal("Node should not be nil")
	}

	node.Entry = "some value"
	if node.IsEmptyOrNilNode() {
		t.Fatal("Node should not be considered empty")
	}

	if !node.IsLeafNode() {
		t.Fatal("Node should be leaf node")
	}

}

func TestNode_IsEmptyNodeNonEntryValue(t *testing.T) {
	node := Node{}

	if !node.IsEmptyOrNilNode() {
		t.Fatal("Node should be considered empty")
	}

	if node.IsNil() {
		t.Fatal("Node should not be nil")
	}

	node.Annotations = "testAnnotations"
	if node.IsEmptyOrNilNode() {
		t.Fatal("Node should not be considered empty")
	}

	if !node.IsLeafNode() {
		t.Fatal("Node should be leaf node")
	}

}

func TestTreeCreation(t *testing.T) {

	root := Node{}
	if root.IsNil() {
		t.Fatal("Node should not be nil (but be empty).")
	}
	if !root.IsEmptyOrNilNode() {
		t.Fatal("Node has not been correctly detected as empty.")
	}

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

	if !leftNode.IsCombination() {
		t.Fatal("Combination has not been detected.")
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

	if !rightNode.IsCombination() {
		t.Fatal("Combination has not been detected.")
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
	leaves := root.GetLeafNodes(true)

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
Tests the distance function and retrieval of logical operators in between; applies upward, and downward search,
and across branches
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
	res, ops, err := FindLogicalLinkage(&leftLeftChildNode, &leftRightChildNode)

	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}

	fmt.Println(ops)
	if len(ops) != 1 || ops[0] != "123" {
		t.Fatal("Logical operators are incorrectly determined.")
	}

	// Simple right to left
	res, ops, err = FindLogicalLinkage(&leftRightChildNode, &leftLeftChildNode)

	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}
	if len(ops) != 1 || ops[0] != "123" {
		t.Fatal("Logical operators are incorrectly determined.")
	}

	// Right branch
	// Simple left to right
	res, ops, err = FindLogicalLinkage(&rightLeftChildNode, &rightRightChildNode)

	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}
	if len(ops) != 1 || ops[0] != "456" {
		t.Fatal("Logical operators are incorrectly determined.")
	}

	// Find adjacent leaves combined by single operator
	res, ops, err = FindLogicalLinkage(&rightRightChildNode, &rightLeftChildNode)

	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}
	if len(ops) != 1 || ops[0] != "456" {
		t.Fatal("Logical operators are incorrectly determined.")
	}

	// Across branches
	res, ops, err = FindLogicalLinkage(&leftLeftChildNode, &rightRightChildNode)

	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}

	if len(ops) != 3 || ops[0] != "123" || ops[1] != "789" || ops[2] != "456" {
		t.Fatal("Logical operators are incorrectly determined.")
	}

	// Unbalanced tree structures

	subnode := Node{}
	leftSubnode := Node{Entry: "left subsub"}
	rightSubnode := Node{Entry: "right subsub"}

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
	leftRightChildNode.Entry = nil
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

	// Detect logical operators across broken tree (should work - nil values are ignored)
	res, ops, err = FindLogicalLinkage(&leftLeftChildNode, &rightSubnode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Nil nodes may have caused error in search.")
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

	// Find links within branch
	res, ops, err = FindLogicalLinkage(&leftLeftChildNode, &rightSubnode)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Link between nodes could not be found.")
	}
	if len(ops) != 3 || ops[0] != "123" || ops[1] != "AND" || ops[2] != "OR" {
		t.Fatal("Logical operators are incorrectly detected: ", ops)
	}

	// Find links across root from left to right
	res, ops, _ = FindLogicalLinkage(&leftLeftChildNode, &rightRightChildNode)
	if !res {
		t.Fatal("Link between nodes could not be found.")
	}
	if len(ops) != 3 || ops[0] != "123" || ops[1] != "789" || ops[2] != "456" {
		t.Fatal("Logical operators are incorrectly detected: ", ops)
	}

	// Find links across root from right to left
	res, ops, _ = FindLogicalLinkage(&rightRightChildNode, &leftLeftChildNode)
	if !res {
		t.Fatal("Link between nodes could not be found.")
	}
	if len(ops) != 3 || ops[0] != "456" || ops[1] != "789" || ops[2] != "123" {
		t.Fatal("Logical operators are incorrectly detected: ", ops)
	}

}

/*
Tests the combination of nodes into new node
*/
func TestNodeCombination(t *testing.T) {

	node1 := Node{}
	node1.ComponentType = "Type2"
	leftSubnode1 := Node{Entry: "left subvalue1"}
	leftSubnode1.ComponentType = "Type1"
	rightSubnode1 := Node{Entry: "right subvalue1"}
	rightSubnode1.ComponentType = "Type1"

	res, err := node1.InsertLeftNode(&leftSubnode1)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when adding new node should not happen")
	}

	res, err = node1.InsertRightNode(&rightSubnode1)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when adding new node should not happen")
	}

	node2 := Node{}
	node2.ComponentType = "Type2"
	leftSubnode2 := Node{Entry: "left subvalue2"}
	rightSubnode2 := Node{Entry: "right subvalue2"}

	res, err = node2.InsertLeftNode(&leftSubnode2)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when adding new node should not happen")
	}

	res, err = node2.InsertRightNode(&rightSubnode2)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when adding new node should not happen")
	}

	combinedNode, err := Combine(&node1, &node2, "AND")
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error during combination of nodes.")
	}

	if combinedNode.LogicalOperator != "AND" {
		t.Fatal("Logical operator was not correctly assigned in combined node.")
	}

	if combinedNode.GetComponentName() != "Type2" {
		t.Fatal("Wrong component name introduced in combination.")
	}

	if combinedNode.Left != &node1 {
		t.Fatal("Left node was not correctly assigned in combined node.")
	}

	if combinedNode.Right != &node2 {
		t.Fatal("Right node was not correctly assigned in combined node.")
	}

	if combinedNode.Left.Left.Entry != "left subvalue1" {
		t.Fatal("Entry value of leaf was not correctly assigned in combined node.")
	}

	if combinedNode.Left.Right.Entry != "right subvalue1" {
		t.Fatal("Entry value of leaf was not correctly assigned in combined node.")
	}

	if combinedNode.Right.Left.Entry != "left subvalue2" {
		t.Fatal("Entry value of leaf was not correctly assigned in combined node.")
	}

	if combinedNode.Right.Right.Entry != "right subvalue2" {
		t.Fatal("Entry value of leaf was not correctly assigned in combined node.")
	}

}

/*
Tests the combination of nodes into new node under consideration of conflicting types.
*/
func TestNodeCombinationConflictingTypes(t *testing.T) {

	node1 := Node{}
	node1.ComponentType = "Type2"
	leftSubnode1 := Node{Entry: "left subvalue1"}
	leftSubnode1.ComponentType = "Type1"
	rightSubnode1 := Node{Entry: "right subvalue1"}
	rightSubnode1.ComponentType = "Type1"

	res, err := node1.InsertLeftNode(&leftSubnode1)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when adding new node should not happen")
	}

	res, err = node1.InsertRightNode(&rightSubnode1)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when adding new node should not happen")
	}

	node2 := Node{}
	node2.ComponentType = "Type1"
	leftSubnode2 := Node{Entry: "left subvalue2"}
	rightSubnode2 := Node{Entry: "right subvalue2"}

	res, err = node2.InsertLeftNode(&leftSubnode2)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when adding new node should not happen")
	}

	res, err = node2.InsertRightNode(&rightSubnode2)
	if !res || err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when adding new node should not happen")
	}

	_, err = Combine(&node1, &node2, "AND")
	if err.ErrorCode == TREE_NO_ERROR {
		t.Fatal("Did not pick up on combination of incompatible components")
	}

}

/*
Test for inheriting shared elements using the append inheritance mode.
*/
func TestExtractInheritAppend(t *testing.T) {

	SHARED_ELEMENT_INHERITANCE_MODE = SHARED_ELEMENT_INHERIT_APPEND

	root := Node{}
	root.SharedLeft = []string{"Shared top left"}
	root.SharedRight = []string{"Shared top right"}

	leftChild := Node{LogicalOperator: "AND"}
	leftChild.SharedLeft = []string{"Shared middle left"}

	leftLeftChild := Node{Entry: "left left"}
	_, err := leftChild.InsertLeftNode(&leftLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	leftRightChild := Node{Entry: "left right"}
	_, err = leftChild.InsertRightNode(&leftRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertLeftNode(&leftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightChild := Node{LogicalOperator: "XOR"}
	rightChild.SharedRight = []string{"Shared middle right"}

	rightLeftChild := Node{Entry: "right left"}
	_, err = rightChild.InsertLeftNode(&rightLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightRightChild := Node{Entry: "right right"}
	_, err = rightChild.InsertRightNode(&rightRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertRightNode(&rightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	fmt.Println(root.String())

	if len(leftLeftChild.GetSharedLeft()) != 2 {
		t.Fatal("Left left child does not have the correct number of shared left elements")
	}

	if leftLeftChild.GetSharedLeft()[0] != "Shared top left" {
		t.Fatal("Left left child's first element is incorrect.")
	}

	if leftLeftChild.GetSharedLeft()[1] != "Shared middle left" {
		t.Fatal("Left left child's second element is incorrect.")
	}

	if len(rightLeftChild.GetSharedRight()) != 2 {
		t.Fatal("Right left child does not have the correct number of shared right elements: ", rightLeftChild.GetSharedRight())
	}

	if rightLeftChild.GetSharedRight()[0] != "Shared top right" {
		t.Fatal("Right left child's first element is incorrect.")
	}

	if rightLeftChild.GetSharedRight()[1] != "Shared middle right" {
		t.Fatal("Left left child's second element is incorrect.")
	}

}

/*
Test for inheriting shared elements using the overrid inheritance mode.
*/
func TestExtractInheritOverride(t *testing.T) {

	SHARED_ELEMENT_INHERITANCE_MODE = SHARED_ELEMENT_INHERIT_OVERRIDE

	root := Node{}
	root.SharedLeft = []string{"Shared top left"}
	root.SharedRight = []string{"Shared top right"}

	leftChild := Node{LogicalOperator: "AND"}
	leftChild.SharedLeft = []string{"Shared middle left"}

	leftLeftChild := Node{Entry: "left left"}
	_, err := leftChild.InsertLeftNode(&leftLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	leftRightChild := Node{Entry: "left right"}
	_, err = leftChild.InsertRightNode(&leftRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertLeftNode(&leftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightChild := Node{LogicalOperator: "XOR"}
	rightChild.SharedRight = []string{"Shared middle right"}

	rightLeftChild := Node{Entry: "right left"}
	_, err = rightChild.InsertLeftNode(&rightLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightRightChild := Node{Entry: "right right"}
	_, err = rightChild.InsertRightNode(&rightRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertRightNode(&rightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	fmt.Println(root.String())

	if len(leftLeftChild.GetSharedLeft()) != 2 {
		t.Fatal("Left left child does not have the correct number of shared left elements")
	}

	if leftLeftChild.GetSharedLeft()[0] != "Shared top left" {
		t.Fatal("Left left child's first element is incorrect.")
	}

	if leftLeftChild.GetSharedLeft()[1] != "Shared middle left" {
		t.Fatal("Left left child's second element is incorrect.")
	}

	if len(rightLeftChild.GetSharedRight()) != 2 {
		t.Fatal("Right left child does not have the correct number of shared right elements: ", rightLeftChild.GetSharedRight())
	}

	if rightLeftChild.GetSharedRight()[0] != "Shared top right" {
		t.Fatal("Right left child's first element is incorrect.")
	}

	if rightLeftChild.GetSharedRight()[1] != "Shared middle right" {
		t.Fatal("Left left child's second element is incorrect.")
	}

}

/*
Test for inheriting shared elements using the override inheritance mode, even if child node is present
*/
func TestExtractInheritOverrideChildElement(t *testing.T) {

	SHARED_ELEMENT_INHERITANCE_MODE = SHARED_ELEMENT_INHERIT_OVERRIDE

	root := Node{}
	root.SharedLeft = []string{"Shared top left"}
	root.SharedRight = []string{"Shared top right"}

	leftChild := Node{LogicalOperator: "AND"}
	leftChild.SharedLeft = []string{"Shared middle left"}

	leftLeftChild := Node{Entry: "left left"}
	_, err := leftChild.InsertLeftNode(&leftLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	leftRightChild := Node{Entry: "left right"}
	_, err = leftChild.InsertRightNode(&leftRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertLeftNode(&leftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightChild := Node{LogicalOperator: "XOR"}
	rightChild.SharedRight = []string{"Shared middle right"}

	rightLeftChild := Node{Entry: "right left"}
	_, err = rightChild.InsertLeftNode(&rightLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightRightChild := Node{Entry: "right right"}
	_, err = rightChild.InsertRightNode(&rightRightChild)
	rightRightChild.SharedRight = []string{"lower right"}
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertRightNode(&rightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	fmt.Println(root.String())

	if len(leftLeftChild.GetSharedLeft()) != 2 {
		t.Fatal("Left left child does not have the correct number of shared left elements")
	}

	if leftLeftChild.GetSharedLeft()[0] != "Shared top left" {
		t.Fatal("Left left child's first element is incorrect.")
	}

	if leftLeftChild.GetSharedLeft()[1] != "Shared middle left" {
		t.Fatal("Left left child's second element is incorrect.")
	}

	if len(rightLeftChild.GetSharedRight()) != 2 {
		t.Fatal("Right left child does not have the correct number of shared right elements: ", rightLeftChild.GetSharedRight())
	}

	if rightLeftChild.GetSharedRight()[0] != "Shared top right" {
		t.Fatal("Right left child's first element is incorrect.")
	}

	if rightLeftChild.GetSharedRight()[1] != "Shared middle right" {
		t.Fatal("Left left child's second element is incorrect.")
	}

}

/*
Test for inheriting shared elements using the no inheritance mode, even if child node is present
*/
func TestExtractInheritNoOverride(t *testing.T) {

	// Nothing is inherited, only content of leaf node
	SHARED_ELEMENT_INHERITANCE_MODE = SHARED_ELEMENT_INHERIT_NOTHING

	root := Node{}
	root.SharedLeft = []string{"Shared top left"}
	root.SharedRight = []string{"Shared top right"}

	leftChild := Node{LogicalOperator: "AND"}
	leftChild.SharedLeft = []string{"Shared middle left"}

	leftLeftChild := Node{Entry: "left left"}
	_, err := leftChild.InsertLeftNode(&leftLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	leftRightChild := Node{Entry: "left right"}
	_, err = leftChild.InsertRightNode(&leftRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertLeftNode(&leftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightChild := Node{LogicalOperator: "XOR"}
	rightChild.SharedRight = []string{"Shared middle right"}

	rightLeftChild := Node{Entry: "right left"}
	// special case of leaf-level shared entry
	rightLeftChild.SharedRight = []string{"right left leaf shared entry"}
	_, err = rightChild.InsertLeftNode(&rightLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightRightChild := Node{Entry: "right right"}
	_, err = rightChild.InsertRightNode(&rightRightChild)
	rightRightChild.SharedRight = []string{"lower right"}
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertRightNode(&rightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	fmt.Println(root.String())

	if len(leftLeftChild.GetSharedLeft()) != 0 {
		t.Fatal("Left left child does not have the correct number of shared left elements: ", leftLeftChild.GetSharedLeft())
	}

	if len(rightLeftChild.GetSharedRight()) != 1 {
		t.Fatal("Right left child does not have the correct number of shared right elements: ", rightLeftChild.GetSharedRight())
	}

	if rightLeftChild.GetSharedRight()[0] != "right left leaf shared entry" {
		t.Fatal("Right left child's first element is incorrect.")
	}

}

/*
Test for inheriting shared elements using the inheritance from next higher combination, even if child node is present
*/
func TestExtractInheritFromNextHigherCombination(t *testing.T) {

	// Only elements from next-higher combination are inherited (including leaf-level elements)
	SHARED_ELEMENT_INHERITANCE_MODE = SHARED_ELEMENT_INHERIT_FROM_COMBINATION

	root := Node{}
	root.SharedLeft = []string{"Shared top left"}
	root.SharedRight = []string{"Shared top right"}

	leftChild := Node{LogicalOperator: "AND"}
	leftChild.SharedLeft = []string{"Shared middle left"}

	leftLeftChild := Node{Entry: "left left"}
	_, err := leftChild.InsertLeftNode(&leftLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	leftRightChild := Node{Entry: "left right"}
	_, err = leftChild.InsertRightNode(&leftRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertLeftNode(&leftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightChild := Node{LogicalOperator: "XOR"}
	rightChild.SharedRight = []string{"Shared middle right"}

	rightLeftChild := Node{Entry: "right left"}
	// special case of leaf-level shared entry
	rightLeftChild.SharedRight = []string{"right left leaf shared entry"}
	_, err = rightChild.InsertLeftNode(&rightLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightRightChild := Node{Entry: "right right"}
	_, err = rightChild.InsertRightNode(&rightRightChild)
	rightRightChild.SharedRight = []string{"lower right"}
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertRightNode(&rightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	fmt.Println(root.String())

	if len(leftLeftChild.GetSharedLeft()) != 1 {
		t.Fatal("Left left child does not have the correct number of shared left elements: ", leftLeftChild.GetSharedLeft())
	}

	if leftLeftChild.GetSharedLeft()[0] != "Shared middle left" {
		t.Fatal("Left left child's first element is incorrect.")
	}

	if len(rightLeftChild.GetSharedRight()) != 2 {
		t.Fatal("Right left child does not have the correct number of shared right elements: ", rightLeftChild.GetSharedRight())
	}

	if rightLeftChild.GetSharedRight()[0] != "Shared middle right" {
		t.Fatal("Right left child's first element is incorrect.")
	}

	if rightLeftChild.GetSharedRight()[1] != "right left leaf shared entry" {
		t.Fatal("Right left child's second element is incorrect.")
	}

}

/*
Tests the removal of nodes from tree.
*/
func TestRemoveNodeFromTree(t *testing.T) {

	root := Node{}
	root.SharedLeft = []string{"Shared top left"}
	root.SharedRight = []string{"Shared top right"}
	root.LogicalOperator = "AND"

	leftChild := Node{LogicalOperator: "AND"}
	leftChild.SharedLeft = []string{"Shared middle left"}

	leftLeftChild := Node{Entry: "left left"}
	_, err := leftChild.InsertLeftNode(&leftLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	leftRightChild := Node{Entry: "left right"}
	_, err = leftChild.InsertRightNode(&leftRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertLeftNode(&leftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightChild := Node{LogicalOperator: "XOR"}
	rightChild.SharedRight = []string{"Shared middle right"}

	rightLeftChild := Node{Entry: "right left"}
	_, err = rightChild.InsertLeftNode(&rightLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightRightChild := Node{LogicalOperator: "AND"}
	_, err = rightChild.InsertRightNode(&rightRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightRightLeftChild := Node{Entry: "right right left"}
	_, err = rightRightChild.InsertLeftNode(&rightRightLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightRightRightChild := Node{Entry: "right right right"}
	_, err = rightRightChild.InsertRightNode(&rightRightRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertRightNode(&rightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	fmt.Println("Before:", root)

	// Remove right left child

	// Attempt to remove node
	res, err := RemoveNodeFromTree(&rightLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Removal of node from tree failed (Error:", err.ErrorMessage, ")")
	}
	if !res {
		t.Fatal("Removal of node from tree returned false, but did not produce error.")
	}

	if rightLeftChild.Parent != nil {
		t.Fatal("Removed node should have empty parent.")
	}

	if rightRightChild.Parent != &root {
		t.Fatal("Sibling node of removed node has not been properly reassigned.")
	}

	if rightRightRightChild.Parent != &rightRightChild {
		t.Fatal("Reassignment of children to collapsed node failed.")
	}

	if rightRightLeftChild.Parent != &rightRightChild {
		t.Fatal("Reassignment of children to collapsed node failed.")
	}

	// Remove left left child

	res, err = RemoveNodeFromTree(&leftLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Removal of node from tree failed (Error:", err.ErrorMessage, ")")
	}
	if !res {
		t.Fatal("Removal of node from tree returned false, but did not produce error.")
	}

	if leftLeftChild.Parent != nil {
		t.Fatal("Removed node should have empty parent.")
	}

	if leftRightChild.Parent != &root {
		t.Fatal("Sibling node of removed node has not been properly reassigned.")
	}

}

/*
Tests root node retrieval functions in variants for synthetic root nodes and actual root nodes.
*/
func TestNode_GetSyntheticRootNode(t *testing.T) {
	root := Node{}
	root.SharedLeft = []string{"Shared top left"}
	root.SharedRight = []string{"Shared top right"}
	root.LogicalOperator = "bAND"

	leftChild := Node{LogicalOperator: "AND"}
	leftChild.SharedLeft = []string{"Shared middle left"}

	leftLeftChild := Node{Entry: "left left"}
	_, err := leftChild.InsertLeftNode(&leftLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	leftRightChild := Node{Entry: "left right"}
	_, err = leftChild.InsertRightNode(&leftRightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertLeftNode(&leftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightChild := Node{LogicalOperator: "XOR"}
	rightChild.SharedRight = []string{"Shared middle right"}

	rightLeftChild := Node{Entry: "right left"}
	_, err = rightChild.InsertLeftNode(&rightLeftChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	rightRightChild := Node{Entry: "right right"}
	_, err = rightChild.InsertRightNode(&rightRightChild)
	rightRightChild.SharedRight = []string{"lower right"}
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	_, err = root.InsertRightNode(&rightChild)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree.")
	}

	if rightRightChild.GetNodeBelowSyntheticRootNode() != &rightChild {
		t.Fatal("Wrong identification of synthetic root node:", rightRightChild.GetNodeBelowSyntheticRootNode())
	}

	if rightRightChild.GetRootNode() != &root {
		t.Fatal("Wrong identification of root node:", rightRightChild.GetRootNode())
	}
}

/*
Tests whether HasPrivateNodes function works correctly.
*/
func TestNode_HasPrivateNodes(t *testing.T) {
	n := Node{Entry: "random content"}

	if n.HasPrivateNodes() {
		t.Fatal("Node should not have private nodes.")
	}

	n2 := Node{Entry: "content of linked node"}
	n.PrivateNodeLinks = append(n.PrivateNodeLinks, &n2)

	if !n.HasPrivateNodes() {
		t.Fatal("Private node was not detected.")
	}

}

/*
Test retrieval of annotations across tree with and without inheritance.
Tests both GetAnnotations() and HasAnnotations().
*/
func TestNode_GetAnnotationsHasAnnotations(t *testing.T) {
	//oneNode := Node{Entry: "entry1", Annotations: "top annotation"}
	twoNode := Node{LogicalOperator: SAND_BETWEEN_COMPONENTS, Annotations: "upper annotation"}
	//topRightNode := Node{Entry: "right unused entry"}
	threeNode := Node{Entry: "three entry"}
	fourNode := Node{LogicalOperator: AND, Annotations: "lower annotation"}
	fiveNode := Node{Entry: "entry1-1"}
	sixNode := Node{Entry: "entry1-2"}

	res, err := fourNode.InsertLeftNode(&fiveNode)
	if !res {
		t.Fatal("Error when inserting node. Error:", err.Error())
	}

	res, err = fourNode.InsertRightNode(&sixNode)
	if !res {
		t.Fatal("Error when inserting node. Error:", err.Error())
	}

	res, err = twoNode.InsertLeftNode(&threeNode)
	if !res {
		t.Fatal("Error when inserting node. Error:", err.Error())
	}

	res, err = twoNode.InsertRightNode(&fourNode)
	if !res {
		t.Fatal("Error when inserting node. Error:", err.Error())
	}

	if !twoNode.HasAnnotations() {
		t.Fatal("Error No annotations found, but should have been.")
	}

	if twoNode.GetAnnotations() != "upper annotation" {
		t.Fatal("Error during annotation retrieval. Found:", twoNode.GetAnnotations())
	}

	if !fourNode.HasAnnotations() {
		t.Fatal("Error No annotations found, but should have been.")
	}

	if fourNode.GetAnnotations() != "lower annotation" {
		t.Fatal("Error during annotation retrieval. Found:", fourNode.GetAnnotations())
	}

	if !fiveNode.HasAnnotations() {
		t.Fatal("Error No annotations found, but should have been.")
	}

	if fiveNode.GetAnnotations() == "" || fiveNode.GetAnnotations() != "lower annotation" {
		t.Fatal("Error during annotation retrieval via inheritance. Found:", fiveNode.GetAnnotations())
	}

	if !sixNode.HasAnnotations() {
		t.Fatal("Error No annotations found, but should have been.")
	}

	if sixNode.GetAnnotations() == "" || sixNode.GetAnnotations() != "lower annotation" {
		t.Fatal("Error during annotation retrieval via inheritance. Found:", sixNode.GetAnnotations())
	}

	if threeNode.HasAnnotations() {
		t.Fatal("Error annotations found but should not have any.")
	}

	if threeNode.GetAnnotations() != nil {
		t.Fatal("Error during annotation retrieval. Found:", threeNode.GetAnnotations())
	}
}

/*
Tests retrieval of suffix from nodes across tree (or parent if existing).
*/
func TestNode_GetSuffix(t *testing.T) {

	left := Node{LogicalOperator: XOR, Suffix: "1"}
	// No suffix
	leftLeft := Node{Entry: "leftLeft"}
	leftRight := Node{Entry: "leftRight"}
	rightLeft := Node{Entry: "rightLeft"}
	rightRight := Node{Entry: "rightRight", Suffix: "24"}
	right := Node{LogicalOperator: AND}
	res, err := right.InsertLeftNode(&rightLeft)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	res, err = right.InsertRightNode(&rightRight)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	root := Node{LogicalOperator: OR, Suffix: "17"}
	res, err = root.InsertLeftNode(&left)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	// Add to left node after adding to root
	res, err = left.InsertLeftNode(&leftLeft)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	// Add to left node after adding to root
	res, err = left.InsertRightNode(&leftRight)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	res, err = root.InsertRightNode(&right)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}

	// Assigned suffix
	if left.GetSuffix() != "1" {
		t.Fatal("Did not extract correct suffix: " + left.GetSuffix())
	}

	// Inherit from top-level node, i.e., inherit across multiple levels (top-level node has logical operator and suffix)
	if rightLeft.GetSuffix() != "17" {
		t.Fatal("Did not extract correct suffix: " + rightLeft.GetSuffix())
	}

	// Explicitly assigned suffix
	if rightRight.GetSuffix() != "24" {
		t.Fatal("Did not extract correct suffix: " + rightRight.GetSuffix())
	}

	// only inherit until higher logical operator
	if right.GetSuffix() != "17" {
		t.Fatal("Did not extract correct suffix: " + right.GetSuffix())
	}

	fmt.Println(root.String())

	// Inherit from next higher node (which has logical operator and suffix)
	if leftLeft.GetSuffix() != "1" {
		t.Fatal("Did not extract correct suffix: " + leftLeft.GetSuffix())
	}

	// Inherit from next higher node (which has logical operator and suffix)
	if leftRight.GetSuffix() != "1" {
		t.Fatal("Did not extract correct suffix: " + leftRight.GetSuffix())
	}

	// Returns own suffix
	if root.GetSuffix() != "17" {
		t.Fatal("Did not extract correct suffix: " + root.GetSuffix())
	}
}

/*
Tests retrieval of component names across tree (or parent if existing).
*/
func TestNode_GetComponentName(t *testing.T) {

	left := Node{Entry: "Left", ComponentType: "leftComp"}
	// No suffix
	rightLeft := Node{Entry: "rightLeft"}
	rightRight := Node{Entry: "rightRight"}
	right := Node{LogicalOperator: AND, ComponentType: "rightComp"}
	res, err := right.InsertLeftNode(&rightLeft)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	res, err = right.InsertRightNode(&rightRight)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	root := Node{LogicalOperator: OR, ComponentType: "topComp"}
	res, err = root.InsertLeftNode(&left)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	res, err = root.InsertRightNode(&right)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}

	if left.GetComponentName() != "leftComp" {
		t.Fatal("Did not extract correct component: " + left.GetComponentName())
	}

	if rightLeft.GetComponentName() != "rightComp" {
		t.Fatal("Did not extract correct component: " + rightLeft.GetComponentName())
	}

	if rightRight.GetComponentName() != "rightComp" {
		t.Fatal("Did not extract correct component: " + rightRight.GetComponentName())
	}

	if right.GetComponentName() != "rightComp" {
		t.Fatal("Did not extract correct component: " + right.GetComponentName())
	}

	if root.GetComponentName() != "topComp" {
		t.Fatal("Did not extract correct component: " + root.GetComponentName())
	}
}

/*
Tests calculation of state complexity
*/
func TestNode_CalculateStateComplexity(t *testing.T) {

	left := Node{Entry: "Left", ComponentType: "leftComp"}
	rightLeft := Node{Entry: "rightLeft"}
	rightRight := Node{Entry: "rightRight"}
	right := Node{LogicalOperator: AND, ComponentType: "rightComp"}
	res, err := right.InsertLeftNode(&rightLeft)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	res, err = right.InsertRightNode(&rightRight)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	root := Node{LogicalOperator: OR, ComponentType: "topComp"}
	res, err = root.InsertLeftNode(&left)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	res, err = root.InsertRightNode(&right)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}

	r := &root

	complexity, err := r.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when calculating complexity:", err)
	}
	if complexity != 3 {
		t.Error("Test returning wrong state complexity. Value:", complexity)
	}

}

/*
Tests substitution of node in tree structure on intermediate level.
*/
func TestNode_SubstituteNodeReferenceIntermediateNode(t *testing.T) {

	left := Node{LogicalOperator: XOR, Suffix: "1"}
	leftLeft := Node{Entry: "leftLeft"}
	leftRight := Node{Entry: "leftRight"}
	rightLeft := Node{Entry: "rightLeft"}
	rightRight := Node{Entry: "rightRight", Suffix: "24"}
	right := Node{LogicalOperator: AND}
	res, err := right.InsertLeftNode(&rightLeft)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	res, err = right.InsertRightNode(&rightRight)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	root := Node{LogicalOperator: OR, Suffix: "17"}
	res, err = root.InsertLeftNode(&left)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	// Add to left node after adding to root
	res, err = left.InsertLeftNode(&leftLeft)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	// Add to left node after adding to root
	res, err = left.InsertRightNode(&leftRight)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	res, err = root.InsertRightNode(&right)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}

	// Create new node
	newNode := Node{Parent: nil, Entry: "newNode"}

	// Substitute intermediate node on left side
	err2 := substituteNodeReferenceInTree(&left, &newNode, &root, true, true)
	if err2.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when substituting node in tree. Error:", err2)
	}

	if root.Left != &newNode {
		t.Fatal("Failed to substitute new node on left side.")
	}

	if newNode.Left != &leftLeft {
		t.Fatal("Failed to connect left child following substitution.")
	}

	if newNode.Right != &leftRight {
		t.Fatal("Failed to connect right child following substitution.")
	}

	if newNode.Parent != &root {
		t.Fatal("Failed to connect root as new parent following substitution.")
	}

	// Create another new node
	anotherNewNode := Node{Parent: nil, Entry: "anotherNewNode"}

	// Substitute intermediate node on right side, but with different parameters (no parent and children linkage)
	err3 := substituteNodeReferenceInTree(&right, &anotherNewNode, &root, false, false)
	if err3.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when substituting node in tree. Error:", err3)
	}

	if root.Right != &anotherNewNode {
		t.Fatal("Failed to substitute new node on left side.")
	}

	if anotherNewNode.Left != nil {
		t.Fatal("Should not have left child following substitution, but have child ", anotherNewNode.Left)
	}

	if anotherNewNode.Right != nil {
		t.Fatal("Should not have right child following substitution, but have child ", anotherNewNode.Right)
	}

	if anotherNewNode.Parent != nil {
		t.Fatal("Should not have parent following substitution, but have parent ", anotherNewNode.Parent)
	}

}

/*
Tests substitution of node in tree structure on leaf level. This assesses recursive calls as well as nil values.
*/
func TestNode_SubstituteNodeReferenceLeafNode(t *testing.T) {

	left := Node{LogicalOperator: XOR, Suffix: "1"}
	leftLeft := Node{Entry: "leftLeft"}
	leftRight := Node{Entry: "leftRight"}
	rightLeft := Node{Entry: "rightLeft"}
	rightRight := Node{Entry: "rightRight", Suffix: "24"}
	right := Node{LogicalOperator: AND}
	res, err := right.InsertLeftNode(&rightLeft)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	res, err = right.InsertRightNode(&rightRight)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	root := Node{LogicalOperator: OR, Suffix: "17"}
	res, err = root.InsertLeftNode(&left)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	// Add to left node after adding to root
	res, err = left.InsertLeftNode(&leftLeft)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	// Add to left node after adding to root
	res, err = left.InsertRightNode(&leftRight)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}
	res, err = root.InsertRightNode(&right)
	if err.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when populating tree. Error:", err)
	}
	if !res {
		t.Fatal("Error when populating tree. Error:", res)
	}

	// Create new node
	newNode := Node{Parent: nil, Entry: "newNode"}

	// Substitute intermediate node on left side
	err2 := substituteNodeReferenceInTree(&leftLeft, &newNode, &root, true, true)
	if err2.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when substituting node in tree. Error:", err2)
	}

	if root.Left.Left != &newNode {
		t.Fatal("Failed to substitute new node on left side.")
	}

	if newNode.Left != nil {
		t.Fatal("Should not have left child following substitution, but have child ", newNode.Left)
	}

	if newNode.Right != nil {
		t.Fatal("Should not have right child following substitution, but have child ", newNode.Right)
	}

	if newNode.Parent != &left {
		t.Fatal("Failed to connect root as new parent following substitution.")
	}

	// Create another new node
	anotherNewNode := Node{Parent: nil, Entry: "anotherNewNode"}

	// Substitute intermediate node on right side, including parent and children linkage
	err3 := substituteNodeReferenceInTree(&rightRight, &anotherNewNode, &root, true, true)
	if err3.ErrorCode != TREE_NO_ERROR {
		t.Fatal("Error when substituting node in tree. Error:", err3)
	}

	if root.Right.Right != &anotherNewNode {
		t.Fatal("Failed to substitute new node on left side.")
	}

	if anotherNewNode.Left != nil {
		t.Fatal("Should not have left child following substitution, but have child ", anotherNewNode.Left)
	}

	if anotherNewNode.Right != nil {
		t.Fatal("Should not have right child following substitution, but have child ", anotherNewNode.Right)
	}

	if anotherNewNode.Parent != root.Right {
		t.Fatal("Should have different parent following substitution, but have parent ", anotherNewNode.Parent)
	}

}

/*
Tests substitution of node in tree structure with nil input.
*/
func TestNode_SubstituteNodeReferenceNilInput(t *testing.T) {

	// Create another new node
	anotherNewNode := Node{Parent: nil, Entry: "anotherNewNode"}

	// Substitute intermediate node on right side, including parent and children linkage
	err := substituteNodeReferenceInTree(&anotherNewNode, nil, nil, true, true)

	if err.ErrorCode != TREE_ERROR_NIL_NODE {
		t.Fatal("Substitution should have thrown nil node error, but returned ", err.ErrorCode)
	}

	// Substitute intermediate node on right side, including parent and children linkage
	err = substituteNodeReferenceInTree(nil, &anotherNewNode, nil, true, true)

	if err.ErrorCode != TREE_ERROR_NIL_NODE {
		t.Fatal("Substitution should have thrown nil node error, but returned ", err.ErrorCode)
	}

	// Substitute intermediate node on right side, including parent and children linkage
	err = substituteNodeReferenceInTree(nil, nil, &anotherNewNode, true, true)

	if err.ErrorCode != TREE_ERROR_NIL_NODE {
		t.Fatal("Substitution should have thrown nil node error, but returned ", err.ErrorCode)
	}

}

//Collapse adjacent entries in logical operators - CollapseAdjacentOperators()
