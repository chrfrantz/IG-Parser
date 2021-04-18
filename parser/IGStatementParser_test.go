package parser

import (
	"IG-Parser/tree"
	"fmt"
	"testing"
)

func TestStatementParsingIncludingSyntheticANDs(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(as well as (review [AND] (audit [AND] challenge))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	s, err := ParseStatement(text)

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	if s.Attributes.Entry != "National Organic Program's Program Manager" {
		t.Fatal("Parsed element value is incorrect")
	}

	if s.Attributes.CountLeaves() != 1 || s.Attributes.CalculateDepth() != 0 {
		t.Fatal("Wrong leaf count or depth calculation.")
	}

	if s.Deontic.Entry != "may" {
		t.Fatal("Parsed element value is incorrect")
	}

	if s.Deontic.CountLeaves() != 1 || s.Deontic.CalculateDepth() != 0 {
		t.Fatal("Wrong leaf count or depth calculation")
	}

	if s.Aim.LogicalOperator != "sAND" {
		t.Fatal("Parsed element value is incorrect")
	}

	if s.Aim.Left.Entry != "inspect" || s.Aim.Right.Left.Entry != "review" ||
		s.Aim.Right.Right.LogicalOperator != "AND" || s.Aim.Right.Right.Left.Entry != "audit" ||
		s.Aim.Right.Right.Right.Entry != "challenge" {
		t.Fatal("Parsed element values are incorrect")
	}

	if s.Aim.CountLeaves() != 4 || s.Aim.CalculateDepth() != 3 {
		t.Fatal("Wrong leaf count or depth calculation")
	}

	if s.DirectObject.Left.Left.Entry != "certified production and" ||
		s.DirectObject.Left.Right.Entry != "handling operations and" ||
		s.DirectObject.LogicalOperator != "AND" ||
		s.DirectObject.Right.Entry != "accredited certifying agents" {
		t.Fatal("Parsed element values are incorrect")
	}

	if s.DirectObject.CountLeaves() != 3 || s.DirectObject.CalculateDepth() != 2 {
		t.Fatal("Wrong leaf count or depth calculation")
	}

	if s.ExecutionConstraintSimple.Left.Entry != "on behalf of the Secretary" ||
		s.ExecutionConstraintSimple.LogicalOperator != "sAND" ||
		s.ExecutionConstraintSimple.Right.Left.Entry != "Act or" ||
		s.ExecutionConstraintSimple.Right.LogicalOperator != "XOR" ||
		s.ExecutionConstraintSimple.Right.Right.Entry != "regulations in this part" {
		t.Fatal("Parsed element values are incorrect")
	}

	if s.ExecutionConstraintSimple.CountLeaves() != 3 || s.ExecutionConstraintSimple.CalculateDepth() != 2 {
		t.Fatal("Wrong leaf count or depth calculation")
	}

}

/*
Test the correct generation of leaf arrays from statements.
 */
func TestLeafArrayGeneration(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	s, err := ParseStatement(text)

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	nodeArray, componentIdx := s.GenerateLeafArrays()

	if nodeArray == nil {
		t.Fatal("Generated array should not be empty.")
	}

	if len(nodeArray) != 7 {
		t.Fatal("Wrong number of array elements in generated leaf component array: ", len(nodeArray))
	}

	element := nodeArray[0]

	if len(element) != 1 || element[0].Entry != "National Organic Program's Program Manager" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// Deontic
	element = nodeArray[1]

	if len(element) != 1 || element[0].Entry != "may" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// first aim entry
	element = nodeArray[2]

	if len(element) != 1 || element[0].Entry != "inspect and" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// second aim array
	element = nodeArray[3]

	if len(element) != 3 || element[0].Entry != "review" ||
		element[1].Entry != "refresh" || element[2].Entry != "drink" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// second aim array
	element = nodeArray[3]

	if len(element) != 3 || element[0].Entry != "review" ||
		element[1].Entry != "refresh" || element[2].Entry != "drink" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// object
	element = nodeArray[4]

	if len(element) != 3 || element[0].Entry != "certified production and" ||
		element[1].Entry != "handling operations and" || element[2].Entry != "accredited certifying agents" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// execution constraint
	element = nodeArray[5]

	if len(element) != 1 || element[0].Entry != "on behalf of the Secretary" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// execution constraint
	element = nodeArray[6]

	if len(element) != 2 || element[0].Entry != "Act or" || element[1].Entry != "regulations in this part" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	fmt.Println(componentIdx)

	if componentIdx[tree.ATTRIBUTES] != 1 || componentIdx[tree.DIRECT_OBJECT] != 1 ||
		componentIdx[tree.EXECUTION_CONSTRAINT] != 2 || componentIdx[tree.DEONTIC] != 1 ||
		componentIdx[tree.AIM] != 2 {
		t.Fatal("Component element count is incorrect.")
	}

	if componentIdx[tree.CONSTITUTED_ENTITY] != 0 || componentIdx[tree.CONSTITUTED_ENTITY_PROPERTY] != 0 ||
		componentIdx[tree.CONSTITUTIVE_FUNCTION] != 0 || componentIdx[tree.CONSTITUTING_PROPERTIES] != 0 ||
		componentIdx[tree.CONSTITUTING_PROPERTIES_PROPERTY] != 0 {
		t.Fatal("Component element count is not empty for some elements.")
	}
}

func TestSyntheticRootRetrieval(t *testing.T) {

	text := "I(inspect and), I(sustain (review [AND] (refresh [AND] drink)))"

	s, err := ParseStatement(text)

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	nodeArray, componentIdx := s.GenerateLeafArrays()

	fmt.Println(nodeArray)
	fmt.Println(componentIdx)

	if len(nodeArray) != 2 {
		t.Fatal("Wrong number of array entries returned.")
	}

	// Test basic root detection function
	if nodeArray[1][0].GetSyntheticRootNode().LogicalOperator != "AND" ||
		nodeArray[1][0].GetSyntheticRootNode().Left.Entry != "review" ||
		nodeArray[1][0].GetSyntheticRootNode().Right.LogicalOperator != "AND" {
		t.Fatal("Root node was wrongly detected.")
	}

	// Now link both leaves with synthetic AND (sAND)
	newRoot := tree.Node{}
	newRoot.LogicalOperator = "XOR"
	res, errAdd := newRoot.InsertLeftNode(nodeArray[0][0])
	if !res || errAdd.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Addition of left node failed. Error: ", errAdd)
	}
	res, errAdd = newRoot.InsertRightNode(nodeArray[1][0].GetSyntheticRootNode())
	if !res || errAdd.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Addition of right node failed. Error: ", errAdd)
	}

	if nodeArray[1][0].GetSyntheticRootNode().LogicalOperator != "XOR" ||
		nodeArray[1][0].GetSyntheticRootNode().Left.Entry != "inspect and" ||
		nodeArray[1][0].GetSyntheticRootNode().Right.LogicalOperator != "AND" ||
		nodeArray[1][0].GetSyntheticRootNode().Right.Left.Entry != "review" ||
		nodeArray[1][0].GetSyntheticRootNode().Right.Right.Left.Entry != "refresh" ||
		nodeArray[1][0].GetSyntheticRootNode().Right.Right.Right.Entry != "drink" {
		t.Fatal("Root node in new node combination was wrongly detected.")
	}

	newRoot.LogicalOperator = tree.SAND

	if nodeArray[1][0].GetSyntheticRootNode().LogicalOperator != "AND" ||
		nodeArray[1][0].GetSyntheticRootNode().Left.Entry != "review" ||
		nodeArray[1][0].GetSyntheticRootNode().Right.LogicalOperator != "AND" {
		t.Fatal("Root node in new node combination was wrongly detected.")
	}

}

func TestExcessiveParentheses(t *testing.T) {

	// Test excessive right parentheses
	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may)) " +
		"I(inspect), I(as well as (review [AND] (audit [AND] challenge))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	_, err := ParseStatement(text)

	if err.ErrorCode != tree.PARSING_ERROR_IMBALANCED_PARENTHESES {
		t.Fatal("Test did not pick up on unbalanced parentheses")
	}

	// Test excessive left parentheses
	text = "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(as well as (review [AND] (audit [AND] challenge))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	_, err = ParseStatement(text)

	if err.ErrorCode != tree.PARSING_ERROR_IMBALANCED_PARENTHESES {
		t.Fatal("Test did not pick up on unbalanced parentheses")
	}

}

func TestComponentLevelNestedStatement(t *testing.T) {
	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{A(Programme Manager) I(suspects) Bdir(violations) Cac{A(NOP Manager) I(orders) Bdir(review)}}"

	s, err := ParseStatement(text)

	fmt.Println(s.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	nodeArray, componentIdx := s.GenerateLeafArrays()

	if nodeArray == nil {
		t.Fatal("Generated array should not be empty.")
	}

	if len(nodeArray) != 8 {
		t.Fatal("Wrong number of array elements in generated leaf component array: ", len(nodeArray), "Contents: ", nodeArray)
	}

	element := nodeArray[0]

	if len(element) != 1 || element[0].Entry != "National Organic Program's Program Manager" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// Deontic
	element = nodeArray[1]

	if len(element) != 1 || element[0].Entry != "may" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// first aim entry
	element = nodeArray[2]

	if len(element) != 1 || element[0].Entry != "inspect and" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// second aim array
	element = nodeArray[3]

	if len(element) != 3 || element[0].Entry != "review" ||
		element[1].Entry != "refresh" || element[2].Entry != "drink" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// second aim array
	element = nodeArray[3]

	if len(element) != 3 || element[0].Entry != "review" ||
		element[1].Entry != "refresh" || element[2].Entry != "drink" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// object
	element = nodeArray[4]

	if len(element) != 3 || element[0].Entry != "certified production and" ||
		element[1].Entry != "handling operations and" || element[2].Entry != "accredited certifying agents" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	// activation condition
	element = nodeArray[5]

	if element[0].HasPrimitiveEntry() {
		t.Fatal("Did detect activation condition as primitive entry")
	}

	if !element[0].Entry.(tree.Statement).Attributes.HasPrimitiveEntry() {
		t.Fatal("Did not detect attribute as primitive entry")
	}

	if element[0].Entry.(tree.Statement).Attributes.Entry == nil {
		t.Fatal("Did detect attribute entry as nil")
	}

	if element[0].Entry.(tree.Statement).Attributes.Entry.(string) != "Programme Manager" {
		t.Fatal("Incorrectly detected attribute in nested activation condition")
	}

	if element[0].Entry.(tree.Statement).Aim.Entry.(string) != "suspects" {
		t.Fatal("Incorrectly detected attribute in nested activation condition")
	}

	// Test for nested elements
	nestedStmt := element[0].Entry.(tree.Statement)
	leaves, _ := nestedStmt.GenerateLeafArrays()
	if len(leaves) != 4 {
		t.Fatal("Did not leaf elements of first-order nested statement correctly")
	}

	if !element[0].Entry.(tree.Statement).ActivationConditionSimple.IsNil() {
		t.Fatal("Simple activation condition field of nested statement should be nil")
	}

	// Test for second-order nested statements
	if element[0].Entry.(tree.Statement).ActivationConditionComplex.IsNil() {
		t.Fatal("Complex activation condition field of nested statement should not be nil")
	}

	if element[0].Entry.(tree.Statement).ActivationConditionComplex.Entry.(tree.Statement).Attributes.Entry != "NOP Manager" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	if element[0].Entry.(tree.Statement).ActivationConditionComplex.Entry.(tree.Statement).Aim.Entry != "orders" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	if element[0].Entry.(tree.Statement).ActivationConditionComplex.Entry.(tree.Statement).DirectObject.Entry != "review" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	nestedStmt = element[0].Entry.(tree.Statement).ActivationConditionComplex.Entry.(tree.Statement)
	leaves, _ = nestedStmt.GenerateLeafArrays()
	if len(leaves) != 3 {
		t.Fatal("Did not leaf elements of second-order nested statement correctly")
	}

	// execution constraint
	element = nodeArray[6]

	if len(element) != 1 || element[0].Entry != "on behalf of the Secretary" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect. Element values: ", element)
	}

	// execution constraint
	element = nodeArray[7]

	if len(element) != 2 || element[0].Entry != "Act or" || element[1].Entry != "regulations in this part" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	fmt.Println(componentIdx)

	// All fields and activation condition reference should be filled
	if componentIdx[tree.ATTRIBUTES] != 1 || componentIdx[tree.DIRECT_OBJECT] != 1 ||
		componentIdx[tree.EXECUTION_CONSTRAINT] != 2 || componentIdx[tree.DEONTIC] != 1 ||
		componentIdx[tree.AIM] != 2 || componentIdx[tree.ACTIVATION_CONDITION_REFERENCE] != 1 {
		t.Fatal("Component element count is incorrect.")
	}

	// All unused fields should be empty (including simple activation condition)
	if componentIdx[tree.CONSTITUTED_ENTITY] != 0 || componentIdx[tree.CONSTITUTED_ENTITY_PROPERTY] != 0 ||
		componentIdx[tree.CONSTITUTIVE_FUNCTION] != 0 || componentIdx[tree.CONSTITUTING_PROPERTIES] != 0 ||
		componentIdx[tree.CONSTITUTING_PROPERTIES_PROPERTY] != 0 || componentIdx[tree.ACTIVATION_CONDITION] != 0 {
		t.Fatal("Component element count is not empty for some elements.")
	}

}
