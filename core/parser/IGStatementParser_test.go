package parser

import (
	"IG-Parser/core/tree"
	"fmt"
	"strings"
	"testing"
)

/*
Tests complete statement parsing, including the consideration of synthetic ANDs.
*/
func TestStatementParsingIncludingSyntheticANDs(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(as well as (review [AND] (audit [AND] challenge))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	stmt, err := ParseStatement(text)

	s := stmt[0].Entry.(*tree.Statement)

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

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

	if s.Aim.LogicalOperator != "bAND" {
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
		s.ExecutionConstraintSimple.LogicalOperator != "bAND" ||
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
Test the correct generation of leaf arrays from statements without aggregation of implicitly linked components,
tolerating multiple components per type.
*/
func TestLeafArrayGenerationWithoutAggregationOfImplicitlyLinkedComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	s, err := ParseStatement(text)

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	nodeArray, componentIdx := s[0].Entry.(*tree.Statement).GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	if nodeArray == nil {
		t.Fatal("Generated array should not be empty.")
	}

	if len(nodeArray) != 7 {
		t.Fatal("Wrong number of array elements in generated leaf component array: ", len(nodeArray))
	}

	// Attributes
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

	if element[0].Parent.LogicalOperator != "bAND" {
		t.Fatal("Wrong logical operator linking aims:", element[0].Parent.LogicalOperator)
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

/*
Test the correct generation of leaf arrays from statements collation of implicitly linked components, returning
one top-level component per component.
*/
func TestLeafArrayGenerationWithAggregationOfImplicitlyLinkedComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	s, err := ParseStatement(text)

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = true

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	nodeArray, componentIdx := s[0].Entry.(*tree.Statement).GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	if nodeArray == nil {
		t.Fatal("Generated array should not be empty.")
	}

	// Identify collapsed in between component linkage
	if len(nodeArray) != 5 {
		t.Fatal("Wrong number of array elements in generated leaf component array: ", len(nodeArray))
	}

	// Attributes
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

	if len(element) != 4 || element[0].Entry != "inspect and" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect. Number of elements:", len(element))
	}

	if element[0].Parent.LogicalOperator != "bAND" {
		t.Fatal("Wrong logical operator linking aims:", element[0].Parent.LogicalOperator)
	}

	// object
	element = nodeArray[3]

	if len(element) != 3 || element[0].Entry != "certified production and" ||
		element[1].Entry != "handling operations and" || element[2].Entry != "accredited certifying agents" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect. Number of elements:", len(element))
	}

	// execution constraint
	element = nodeArray[4]

	if len(element) != 3 || element[0].Entry != "on behalf of the Secretary" ||
		element[1].Entry != "Act or" || element[2].Entry != "regulations in this part" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect. Number of elements:", len(element))
	}

	if element[0].Parent.LogicalOperator != "bAND" {
		t.Fatal("Wrong logical operator linking aims:", element[0].Parent.LogicalOperator)
	}

	if componentIdx[tree.ATTRIBUTES] != 1 || componentIdx[tree.DIRECT_OBJECT] != 1 ||
		componentIdx[tree.EXECUTION_CONSTRAINT] != 1 || componentIdx[tree.DEONTIC] != 1 ||
		componentIdx[tree.AIM] != 1 {
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

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	nodeArray, componentIdx := s[0].Entry.(*tree.Statement).GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println(nodeArray)
	fmt.Println(componentIdx)

	if len(nodeArray) != 2 {
		t.Fatal("Wrong number of array entries returned.")
	}

	// Test basic root detection function
	if nodeArray[1][0].GetNodeBelowSyntheticRootNode().LogicalOperator != "AND" ||
		nodeArray[1][0].GetNodeBelowSyntheticRootNode().Left.Entry != "review" ||
		nodeArray[1][0].GetNodeBelowSyntheticRootNode().Right.LogicalOperator != "AND" {
		t.Fatal("Root node was wrongly detected.")
	}

	// Now link both leaves with synthetic AND (sAND)
	newRoot := tree.Node{}
	newRoot.LogicalOperator = "XOR"
	res, errAdd := newRoot.InsertLeftNode(nodeArray[0][0])
	if !res || errAdd.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Addition of left node failed. Error: ", errAdd)
	}
	res, errAdd = newRoot.InsertRightNode(nodeArray[1][0].GetNodeBelowSyntheticRootNode())
	if !res || errAdd.ErrorCode != tree.TREE_NO_ERROR {
		t.Fatal("Addition of right node failed. Error: ", errAdd)
	}

	if nodeArray[1][0].GetNodeBelowSyntheticRootNode().LogicalOperator != "XOR" ||
		nodeArray[1][0].GetNodeBelowSyntheticRootNode().Left.Entry != "inspect and" ||
		nodeArray[1][0].GetNodeBelowSyntheticRootNode().Right.LogicalOperator != "AND" ||
		nodeArray[1][0].GetNodeBelowSyntheticRootNode().Right.Left.Entry != "review" ||
		nodeArray[1][0].GetNodeBelowSyntheticRootNode().Right.Right.Left.Entry != "refresh" ||
		nodeArray[1][0].GetNodeBelowSyntheticRootNode().Right.Right.Right.Entry != "drink" {
		t.Fatal("Root node in new node combination was wrongly detected.")
	}

	// TODO: Check for the need to consider SAND_WITHIN_COMPONENTS
	newRoot.LogicalOperator = tree.SAND_BETWEEN_COMPONENTS

	if nodeArray[1][0].GetNodeBelowSyntheticRootNode().LogicalOperator != "AND" ||
		nodeArray[1][0].GetNodeBelowSyntheticRootNode().Left.Entry != "review" ||
		nodeArray[1][0].GetNodeBelowSyntheticRootNode().Right.LogicalOperator != "AND" {
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

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

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

func TestComponentTwoLevelNestedStatement(t *testing.T) {
	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{A(Programme Manager) I(suspects) Bdir(violations) Cac{A(NOP Manager) I(orders) Bdir(review)}}"

	stmt, err := ParseStatement(text)

	s := stmt[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	nodeArray, componentIdx := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	if nodeArray == nil {
		t.Fatal("Generated array should not be empty.")
	}

	if len(nodeArray) != 8 {
		t.Fatal("Wrong number of array elements in generated leaf component array: ", len(nodeArray), "Contents: ", nodeArray)
	}

	// Attributes
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

	if !element[0].Entry.(*tree.Statement).Attributes.HasPrimitiveEntry() {
		t.Fatal("Did not detect attribute as primitive entry")
	}

	if element[0].Entry.(*tree.Statement).Attributes.Entry == nil {
		t.Fatal("Did detect attribute entry as nil")
	}

	if element[0].Entry.(*tree.Statement).Attributes.Entry.(string) != "Programme Manager" {
		t.Fatal("Incorrectly detected attribute in nested activation condition")
	}

	if element[0].Entry.(*tree.Statement).Aim.Entry.(string) != "suspects" {
		t.Fatal("Incorrectly detected attribute in nested activation condition")
	}

	// Test for nested elements
	nestedStmt := element[0].Entry.(*tree.Statement)
	leaves, _ := nestedStmt.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)
	if len(leaves) != 4 {
		t.Fatal("Did not leaf elements of first-order nested statement correctly")
	}

	if !element[0].Entry.(*tree.Statement).ActivationConditionSimple.IsNil() {
		t.Fatal("Simple activation condition field of nested statement should be nil")
	}

	// Test for second-order nested statements
	if element[0].Entry.(*tree.Statement).ActivationConditionComplex.IsNil() {
		t.Fatal("Complex activation condition field of nested statement should not be nil")
	}

	if element[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).Attributes.Entry != "NOP Manager" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	if element[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).Aim.Entry != "orders" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	if element[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).DirectObject.Entry != "review" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	nestedStmt = element[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement)
	leaves, _ = nestedStmt.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)
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

/*
Tests for component-level nesting in activation condition and separate simple combination (coexistence).
*/
func TestComponentTwoLevelNestedStatementAndSimpleCombination(t *testing.T) {
	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{A(Programme Manager) I(suspects) Bdir(violations) Cac{A(NOP Manager) I(orders) Bdir(review)}}" +
		"Cac((regular precondition [AND] another precondition))"

	stmt, err := ParseStatement(text)

	s := stmt[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	nodeArray, componentIdx := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	if nodeArray == nil {
		t.Fatal("Generated array should not be empty.")
	}

	if len(nodeArray) != 9 {
		t.Fatal("Wrong number of array elements in generated leaf component array: ", len(nodeArray), "Contents: ", nodeArray)
	}

	// Attribute
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

	// simple activation condition
	element = nodeArray[5]
	if !element[0].HasPrimitiveEntry() {
		t.Fatal("Did not detect activation condition as primitive entry")
	}

	if element[0].IsNil() {
		t.Fatal("Node should not be nil")
	}

	if element[0].IsEmptyOrNilNode() {
		t.Fatal("Node should not be empty")
	}

	if !element[0].IsLeafNode() {
		t.Fatal("Node should not be leaf node")
	}

	if element[0].Entry.(string) != "regular precondition" {
		t.Fatal("Left nested element in activation condition not correctly identified")
	}

	if !element[1].IsLeafNode() {
		t.Fatal("Node should not be leaf node")
	}

	if element[1].Entry.(string) != "another precondition" {
		t.Fatal("Left nested element in activation condition not correctly identified")
	}

	if element[0].Parent.LogicalOperator != tree.AND {
		t.Fatal("Logical operator linking both simple activation conditions is incorrect:", element[0].Parent.LogicalOperator)
	}

	// complex activation condition
	element = nodeArray[6]

	if !element[0].Entry.(*tree.Statement).Attributes.HasPrimitiveEntry() {
		t.Fatal("Did not detect attribute as primitive entry")
	}

	if element[0].Entry.(*tree.Statement).Attributes.Entry == nil {
		t.Fatal("Did detect attribute entry as nil")
	}

	if element[0].Entry.(*tree.Statement).Attributes.Entry.(string) != "Programme Manager" {
		t.Fatal("Incorrectly detected attribute in nested activation condition")
	}

	if element[0].Entry.(*tree.Statement).Aim.Entry.(string) != "suspects" {
		t.Fatal("Incorrectly detected attribute in nested activation condition")
	}

	// Test for nested elements
	nestedStmt := element[0].Entry.(*tree.Statement)
	leaves, _ := nestedStmt.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)
	if len(leaves) != 4 {
		t.Fatal("Did not leaf elements of first-order nested statement correctly")
	}

	if !element[0].Entry.(*tree.Statement).ActivationConditionSimple.IsNil() {
		t.Fatal("Simple activation condition field of nested statement should be nil")
	}

	// Test for second-order nested statements
	if element[0].Entry.(*tree.Statement).ActivationConditionComplex.IsNil() {
		t.Fatal("Complex activation condition field of nested statement should not be nil")
	}

	if element[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).Attributes.Entry != "NOP Manager" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	if element[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).Aim.Entry != "orders" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	if element[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).DirectObject.Entry != "review" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	nestedStmt = element[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement)
	leaves, _ = nestedStmt.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)
	if len(leaves) != 3 {
		t.Fatal("Did not leaf elements of second-order nested statement correctly")
	}

	// execution constraint
	element = nodeArray[7]

	if len(element) != 1 || element[0].Entry != "on behalf of the Secretary" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect. Element values: ", element)
	}

	// execution constraint
	element = nodeArray[8]

	if len(element) != 2 || element[0].Entry != "Act or" || element[1].Entry != "regulations in this part" {
		t.Fatal("Number of elements or element values in generated array is/are incorrect.")
	}

	fmt.Println(componentIdx)

	// All fields and activation condition reference should be filled
	if componentIdx[tree.ATTRIBUTES] != 1 || componentIdx[tree.DIRECT_OBJECT] != 1 ||
		componentIdx[tree.EXECUTION_CONSTRAINT] != 2 || componentIdx[tree.DEONTIC] != 1 ||
		componentIdx[tree.AIM] != 2 || componentIdx[tree.ACTIVATION_CONDITION] != 1 ||
		componentIdx[tree.ACTIVATION_CONDITION_REFERENCE] != 1 {
		t.Fatal("Component element count is incorrect.")
	}

	// All unused fields should be empty (including simple activation condition)
	if componentIdx[tree.CONSTITUTED_ENTITY] != 0 || componentIdx[tree.CONSTITUTED_ENTITY_PROPERTY] != 0 ||
		componentIdx[tree.CONSTITUTIVE_FUNCTION] != 0 || componentIdx[tree.CONSTITUTING_PROPERTIES] != 0 ||
		componentIdx[tree.CONSTITUTING_PROPERTIES_PROPERTY] != 0 {
		t.Fatal("Component element count is not empty for some elements.")
	}

}

/*
Tests multiple complex activation conditions in a statement
*/
func TestComponentMultipleHorizontallyNestedStatement(t *testing.T) {
	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{A(Programme Manager) I(suspects) Bdir(violations) Cac{A(NOP Manager) I(orders) Bdir(review)}}" +
		"Cac{E(Program Manager) F(is) P(qualified)}"

	stmt, err := ParseStatement(text)

	s := stmt[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	nodeArray, componentIdx := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	if nodeArray == nil {
		t.Fatal("Generated array should not be empty.")
	}

	if len(nodeArray) != 8 {
		t.Fatal("Wrong number of array elements in generated leaf component array: ", len(nodeArray), "Contents: ", nodeArray)
	}

	// Attributes
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

	fmt.Println(element)

	if element[0].HasPrimitiveEntry() {
		t.Fatal("Did detect activation condition as primitive entry")
	}

	if !element[0].Left.Entry.(*tree.Statement).Attributes.HasPrimitiveEntry() {
		t.Fatal("Did not detect attribute as primitive entry (Entry:", element[0].Entry, ")")
	}

	if element[0].Left.Entry.(*tree.Statement).Attributes.Entry == nil {
		t.Fatal("Did detect attribute entry as nil")
	}

	if element[0].Left.Entry.(*tree.Statement).Attributes.Entry.(string) != "Programme Manager" {
		t.Fatal("Incorrectly detected attribute in nested activation condition")
	}

	if element[0].Left.Entry.(*tree.Statement).Aim.Entry.(string) != "suspects" {
		t.Fatal("Incorrectly detected attribute in nested activation condition")
	}

	// Test for nested elements
	nestedStmt := element[0].Left.Entry.(*tree.Statement)
	leaves, _ := nestedStmt.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)
	if len(leaves) != 4 {
		t.Fatal("Did not leaf elements of first-order nested statement correctly")
	}

	if !element[0].Left.Entry.(*tree.Statement).ActivationConditionSimple.IsNil() {
		t.Fatal("Simple activation condition field of nested statement should be nil")
	}

	// Test for second-order nested statements
	if element[0].Left.Entry.(*tree.Statement).ActivationConditionComplex.IsNil() {
		t.Fatal("Complex activation condition field of nested statement should not be nil")
	}

	if element[0].Left.Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).Attributes.Entry != "NOP Manager" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	if element[0].Left.Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).Aim.Entry != "orders" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	if element[0].Left.Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).DirectObject.Entry != "review" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	nestedStmt = element[0].Left.Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement)
	leaves, _ = nestedStmt.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)
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

/*
Tests parsing of nested statement combinations
*/
func TestFlatteningAndParsingOfStatementCombinations(t *testing.T) {

	input := "{Cac{E(Program Manager) F(is) P(qualified)} [AND] " +
		"{Cac{E(Program Participant2) F(is2) P(employed2)} [XOR] " +
		"Cac{E(Program Participant) F(is) P(employed)}}}"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	combo, _, errStmt := ParseIntoNodeTree(input, false, LEFT_BRACE, RIGHT_BRACE)
	if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
		t.Error("Error when parsing nested statements: " + errStmt.ErrorCode)
	}

	// Check whether all leaves have the same prefix
	flatCombo := tree.Flatten(combo.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES))
	sharedPrefix := ""
	for _, node := range flatCombo {
		entry := node.Entry.(string)
		// Extract prefix for node
		prefix := entry[:strings.Index(entry, LEFT_BRACE)]
		if sharedPrefix == "" {
			// Cache it if not already done
			sharedPrefix = prefix
			continue
		}
		// Check if it deviates from previously cached element
		if prefix != sharedPrefix {
			t.Error("Invalid combination of component-level nested statements. Expected component: " +
				sharedPrefix + ", but found: " + prefix)
		}
	}

	if len(tree.Flatten(combo.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES))) != 3 {
		t.Fatal("Wrong number of parsed string nodes.")
	}

	// Parse all entries in tree from string to statement
	err := combo.ParseAllEntries("", func(oldValue string) (*tree.Statement, tree.ParsingError) {
		stmt, errStmt := ParseStatement(oldValue[strings.Index(oldValue, LEFT_BRACE)+1 : strings.LastIndex(oldValue, RIGHT_BRACE)])
		if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
			return stmt[0].Entry.(*tree.Statement), errStmt
		}
		return stmt[0].Entry.(*tree.Statement), tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
	})
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Conversion of string entries to parsed statements failed:", err.Error())
	}

	if len(tree.Flatten(combo.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES))) != 3 {
		t.Fatal("Wrong number of parsed string nodes.")
	}

	fmt.Println(combo.String())

	if combo.Left.Entry.(*tree.Statement).ConstitutedEntity.Entry != "Program Manager" ||
		combo.Left.Entry.(*tree.Statement).ConstitutingProperties.Entry != "qualified" ||
		combo.Left.Entry.(*tree.Statement).ConstitutiveFunction.Entry != "is" ||
		combo.LogicalOperator != "AND" ||
		combo.Right.LogicalOperator != "XOR" ||
		combo.Right.Left.Entry.(*tree.Statement).ConstitutedEntity.Entry != "Program Participant2" ||
		combo.Right.Left.Entry.(*tree.Statement).ConstitutiveFunction.Entry != "is2" ||
		combo.Right.Left.Entry.(*tree.Statement).ConstitutingProperties.Entry != "employed2" ||
		combo.Right.Right.Entry.(*tree.Statement).ConstitutedEntity.Entry != "Program Participant" ||
		combo.Right.Right.Entry.(*tree.Statement).ConstitutiveFunction.Entry != "is" ||
		combo.Right.Right.Entry.(*tree.Statement).ConstitutingProperties.Entry != "employed" {

		t.Fatal("Parsing into statements failed.")
	}

}

/*
Tests parsing of special characters in regular and nested components
*/
func TestSpecialCharacters(t *testing.T) {

	input := "A(A&dsisgj=) I(=#) Bdir((l$.ef% [AND] Ri@,g¤#)) Bind((`?a€v [XOR] (dg/sg) !sdg~£jd*s)) Cac{A(/sd<-g$s%d) D(s%k£g=>js) I(s§d€k+l/g#j!ds)}"

	stmt, err := ParseStatement(input)

	s := stmt[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	if s.Attributes.Entry != "A&dsisgj=" {
		t.Fatal("Failed to detect Attributes")
	}

	if s.Aim.Entry != "=#" {
		t.Fatal("Failed to detect Aim")
	}

	if s.DirectObject.Left.Entry != "l$.ef%" {
		t.Fatal("Failed to detect Direct object left")
	}

	if s.DirectObject.Right.Entry != "Ri@,g¤#" {
		t.Fatal("Failed to detect Direct object right")
	}

	if s.IndirectObject.Left.Entry != "`?a€v" {
		t.Fatal("Failed to detect Indirect object left")
	}

	if s.IndirectObject.Right.Entry != "(dg/sg) !sdg~£jd*s" {
		t.Fatal("Failed to detect Indirect object right")
	}

	if s.ActivationConditionComplex.Entry.(*tree.Statement).Attributes.Entry != "/sd<-g$s%d" {
		t.Fatal("Failed to detect nested Attribute")
	}

	if s.ActivationConditionComplex.Entry.(*tree.Statement).Deontic.Entry != "s%k£g=>js" {
		t.Fatal("Failed to detect nested Deontic")
	}

	if s.ActivationConditionComplex.Entry.(*tree.Statement).Aim.Entry != "s§d€k+l/g#j!ds" {
		t.Fatal("Failed to detect nested Aim")
	}

}

/*
Tests whether parser does not mistakenly pick up component properties (e.g., A,p) as components (e.g., A).
*/
func TestUnambiguousExtractionOfComponentAndRelatedProperties(t *testing.T) {

	input := "A,p(property) A,p1(another prop) A(value)"

	stmt, err := ParseStatement(input)

	s := stmt[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	if s.Attributes.CountLeaves() != 1 {
		t.Fatal("Attributes count should be 1, but is:", s.Attributes.CountLeaves())
	}

	if s.Attributes.Entry.(string) != "value" {
		t.Fatal("Attributes should be 'value', but is:", s.Attributes.Entry.(string))
	}

	if s.AttributesPropertySimple.CountLeaves() != 2 {
		t.Fatal("Attributes Properties count should be 2, but is:", s.AttributesPropertySimple.CountLeaves())
	}

	if s.AttributesPropertySimple.Left.Entry.(string) != "property" {
		t.Fatal("Left Attributes Property should be 'property', but is:", s.AttributesPropertySimple.Left.Entry.(string))
	}

	if s.AttributesPropertySimple.Right.Entry.(string) != "another prop" {
		t.Fatal("Right Attributes Property should be 'another prop', but is:", s.AttributesPropertySimple.Right.Entry.(string))
	}

}

/*
Tests extraction of suffix and annotation in statements in which only a single component of a given type is present.
*/
func TestExtractSuffixAndAnnotationsSingleComponentValue(t *testing.T) {

	// Single component entry
	text := "A1[annotation=(left,right)](content)"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", false, text, "(", ")")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	if suffix != "1" {
		t.Fatal("Suffix should be 1 (from first element), but is:", suffix)
	}

	if annotation != "[annotation=(left,right)]" {
		t.Fatal("Annotation should be [annotation=(left,right)] (from first element), but is:", annotation)
	}

	if content != "A(content)" {
		t.Fatal("Content should have been raw component entry without suffix or annotation of first element, but is:", content)
	}

	fmt.Println("Suffix:", suffix, "; Annotation:", annotation, "; Content:", content)

}

/*
Tests extraction of suffix only in statements in which single component of a given type is present.
*/
func TestExtractSuffixOnlySingleComponentValue(t *testing.T) {

	// Single component entry
	text := "A1(content)"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", false, text, "(", ")")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	if suffix != "1" {
		t.Fatal("Suffix should be 1 (from first element), but is:", suffix)
	}

	if annotation != "" {
		t.Fatal("Annotation should be empty, but is:", annotation)
	}

	if content != "A(content)" {
		t.Fatal("Content should have been raw component entry without suffix or annotation of first element, but is:", content)
	}

	fmt.Println("Suffix:", suffix, "; Annotation:", annotation, "; Content:", content)

}

/*
Tests extraction of annotation only in statements in which single component of a given type is present.
*/
func TestExtractAnnotationOnlySingleComponentValue(t *testing.T) {

	// Single component entry
	text := "A[abc=(left;right)](content)"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", false, text, "(", ")")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	if suffix != "" {
		t.Fatal("Suffix should be empty, but is:", suffix)
	}

	if annotation != "[abc=(left;right)]" {
		t.Fatal("Annotation should be [abc=(left;right)], but is:", annotation)
	}

	if content != "A(content)" {
		t.Fatal("Content should have been raw component entry without suffix or annotation of first element, but is:", content)
	}

	fmt.Println("Suffix:", suffix, "; Annotation:", annotation, "; Content:", content)

}

/*
Tests extraction of annotation only in statements with special characters.
*/
func TestExtractAnnotationOnlyWithSpecialCharacters(t *testing.T) {

	// Single component entry
	text := "A[abc=(left|right)](content)"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", false, text, "(", ")")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	if suffix != "" {
		t.Fatal("Suffix should be empty, but is:", suffix)
	}

	if annotation != "[abc=(left|right)]" {
		t.Fatal("Annotation should be [abc=(left|right)], but is:", annotation)
	}

	if content != "A(content)" {
		t.Fatal("Content should have been raw component entry without suffix or annotation of first element, but is:", content)
	}

	fmt.Println("Suffix:", suffix, "; Annotation:", annotation, "; Content:", content)

}

/*
Tests extraction of suffix and annotations in statements with special characters.
*/
func TestExtractSuffixOnlyWithSpecialCharacters(t *testing.T) {

	// Single component entry
	text := "A2#|(cont$ent)"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", false, text, "(", ")")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	if suffix != "2#|" {
		t.Fatal("Suffix should be 2#|, but is:", suffix)
	}

	if annotation != "" {
		t.Fatal("Annotation should be empty, but is:", annotation)
	}

	if content != "A(cont$ent)" {
		t.Fatal("Content should have been raw component entry without suffix or annotation of first element, but is:", content)
	}

	fmt.Println("Suffix:", suffix, "; Annotation:", annotation, "; Content:", content)

}

/*
Tests extraction of suffix and annotations in statements with special characters.
*/
func TestExtractSuffixAndAnnotationWithSpecialCharacters(t *testing.T) {

	// Single component entry
	text := "A2#|[abc=(le#ft|righ$t)](cont$ent)"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", false, text, "(", ")")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	if suffix != "2#|" {
		t.Fatal("Suffix should be empty, but is:", suffix)
	}

	if annotation != "[abc=(le#ft|righ$t)]" {
		t.Fatal("Annotation should be [abc=(left;right)], but is:", annotation)
	}

	if content != "A(cont$ent)" {
		t.Fatal("Content should have been raw component entry without suffix or annotation of first element, but is:", content)
	}

	fmt.Println("Suffix:", suffix, "; Annotation:", annotation, "; Content:", content)

}

/*
Tests whether complete statements are parsed and suffices and annotations stored accordingly in the underlying node structure.
*/
func TestNodeParsingOfSuffixAndAnnotationsAtomicStatement(t *testing.T) {

	text := "A1[annotation1](content1) A2[annotation2](content2) A3(content3) I4[annotation=(left,right)](aim1)"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	stmt := s[0].Entry.(*tree.Statement)

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	// Check Attributes

	if stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Suffix != "1" {
		t.Fatal("Suffix should be 1, but is:", stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Suffix)
	}

	if stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[1][0].Suffix != "2" {
		t.Fatal("Suffix should be 2, but is:", stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[1][0].Suffix)
	}

	if stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[2][0].Suffix != "3" {
		t.Fatal("Suffix should be 3, but is:", stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[2][0].Suffix)
	}

	if stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Annotations != "[annotation1]" {
		t.Fatal("Suffix should be [annotation1], but is:", stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Annotations)
	}

	if stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[1][0].Annotations != "[annotation2]" {
		t.Fatal("Suffix should be [annotation2], but is:", stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[1][0].Annotations)
	}

	if stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[2][0].Annotations != nil {
		t.Fatal("Suffix should be nil, but is:", stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[2][0].Annotations)
	}

	// Check Aim

	if stmt.Aim.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Suffix != "4" {
		t.Fatal("Suffix should be 4, but is:", stmt.Aim.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Suffix)
	}

	if stmt.Aim.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Annotations != "[annotation=(left,right)]" {
		t.Fatal("Suffix should be [annotation=(left,right)], but is:", stmt.Aim.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Annotations)
	}

}

/*
Tests whether complete statements are parsed and whether indexed suffices (e.g., A1,p(content)) and annotations are stored accordingly in the underlying node structure.
*/
func TestNodeParsingOfIndexedSuffixAndAnnotationsAtomicStatement(t *testing.T) {

	text := "A1[annotation1](content1) A1,p[annotation2](content2) A1,p2(content3) I4[annotation=(left,right)](aim1)"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	stmt := s[0].Entry.(*tree.Statement)

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	fmt.Println(stmt.String())

	// Check Attributes

	if stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Suffix != "1" {
		t.Fatal("Suffix should be 1, but is:", stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Suffix)
	}

	if stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Annotations != "[annotation1]" {
		t.Fatal("Suffix should be [annotation1], but is:", stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Annotations)
	}

	if stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].PrivateNodeLinks[0].Annotations != "[annotation2]" {
		t.Fatal("Suffix of private node should be [annotation2], but is:", stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[1][0].Annotations)
	}

	if stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].PrivateNodeLinks[1].Annotations != nil {
		t.Fatal("Suffix should be nil, but is:", stmt.Attributes.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[2][0].Annotations)
	}

	// Check Aim

	if stmt.Aim.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Suffix != "4" {
		t.Fatal("Suffix should be 4, but is:", stmt.Aim.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Suffix)
	}

	if stmt.Aim.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Annotations != "[annotation=(left,right)]" {
		t.Fatal("Suffix should be [annotation=(left,right)], but is:", stmt.Aim.GetLeafNodes(tree.AGGREGATE_IMPLICIT_LINKAGES)[0][0].Annotations)
	}

}

/*
Tests whether complete statements are parsed and suffices and annotations stored accordingly in the underlying node structure.
This test specifically looks at nested statements
*/
func TestNodeParsingOfSuffixAndAnnotationsNestedStatement(t *testing.T) {

	text := "Cac1[leftAnno]{A1[annotation=(left,right)](content) A2[annot](content2) I[regfunc=initiate](action)} Cac2[rightAnno]{A5[|exampleAnnotation](actor)}"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	stmt := s[0].Entry.(*tree.Statement)

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	// Left activation condition
	if stmt.ActivationConditionComplex.Left.Suffix.(string) != "1" {
		t.Fatal("Suffix should be 1, but is:", stmt.ActivationConditionComplex.Left.Suffix)
	}

	if stmt.ActivationConditionComplex.Left.Annotations.(string) != "[leftAnno]" {
		t.Fatal("Annotation should be [leftAnno], but is:", stmt.ActivationConditionComplex.Left.Annotations)
	}

	// Right activation condition
	if stmt.ActivationConditionComplex.Right.Suffix.(string) != "2" {
		t.Fatal("Suffix should be 2, but is:", stmt.ActivationConditionComplex.Right.Suffix)
	}

	if stmt.ActivationConditionComplex.Right.Annotations.(string) != "[rightAnno]" {
		t.Fatal("Annotation should be [rightAnno], but is:", stmt.ActivationConditionComplex.Right.Annotations)
	}

	// First attribute in left activation condition
	if stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Attributes.Left.Suffix.(string) != "1" {
		t.Fatal("Suffix should be 1, but is:", stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Attributes.Left.Suffix.(string))
	}

	if stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Attributes.Left.Annotations.(string) != "[annotation=(left,right)]" {
		t.Fatal("Annotation should be [annotation=(left,right)], but is:", stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Attributes.Left.Annotations.(string))
	}

	// Second attribute in left activation condition
	if stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Attributes.Right.Suffix.(string) != "2" {
		t.Fatal("Suffix should be 2, but is:", stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Attributes.Right.Suffix.(string))
	}

	if stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Attributes.Right.Annotations.(string) != "[annot]" {
		t.Fatal("Annotation should be [annot], but is:", stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Attributes.Right.Annotations.(string))
	}

	// Aim
	if stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Aim.Suffix != nil {
		t.Fatal("Suffix should be nil, but is:", stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Aim.Suffix)
	}

	if stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Aim.Annotations.(string) != "[regfunc=initiate]" {
		t.Fatal("Annotation should be [regfunc=initiate], but is:", stmt.ActivationConditionComplex.Left.Entry.(*tree.Statement).Aim.Annotations)
	}

	// Attributes in right activation condition
	if stmt.ActivationConditionComplex.Right.Entry.(*tree.Statement).Attributes.Suffix.(string) != "5" {
		t.Fatal("Suffix should be 5, but is:", stmt.ActivationConditionComplex.Right.Entry.(*tree.Statement).Attributes.Suffix.(string))
	}

	if stmt.ActivationConditionComplex.Right.Entry.(*tree.Statement).Attributes.Annotations.(string) != "[|exampleAnnotation]" {
		t.Fatal("Annotation should be [|exampleAnnotation], but is:", stmt.ActivationConditionComplex.Right.Entry.(*tree.Statement).Attributes.Annotations.(string))
	}

}

/*
Test proper resolution of component name for primitive element, combination header and elements, and nested statement
*/
func TestComponentNameIdentification(t *testing.T) {

	text := "A(Single Element) D( must) I((combLeft [AND] combRight)) Cac{A(Nested Element) I(perform) Bdir(something)}"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing error for statement", text)
	}

	stmt := s[0].Entry.(*tree.Statement)

	if stmt.Attributes.GetComponentName() != tree.ATTRIBUTES {
		t.Fatal("Incorrect identification of component name for single element")
	}

	if stmt.Deontic.GetComponentName() != tree.DEONTIC {
		t.Fatal("Incorrect identification of component name for single element")
	}

	if stmt.Aim.GetComponentName() != tree.AIM {
		t.Fatal("Incorrect identification of component name for combination node")
	}

	if stmt.Aim.Left.GetComponentName() != tree.AIM {
		t.Fatal("Incorrect identification of component name for combination's left element")
	}

	if stmt.Aim.Right.GetComponentName() != tree.AIM {
		t.Fatal("Incorrect identification of component name for combination's right element")
	}

	if stmt.ActivationConditionComplex.GetComponentName() != tree.ACTIVATION_CONDITION {
		t.Fatal("Incorrect identification of component name for nested component node")
	}

	if stmt.ActivationConditionComplex.Entry.(*tree.Statement).Attributes.GetComponentName() != tree.ATTRIBUTES {
		t.Fatal("Incorrect identification of component name for nested statement's attribute")
	}

	if stmt.ActivationConditionComplex.Entry.(*tree.Statement).Aim.GetComponentName() != tree.AIM {
		t.Fatal("Incorrect identification of component name for nested statement's aim")
	}

	if stmt.ActivationConditionComplex.Entry.(*tree.Statement).DirectObject.GetComponentName() != tree.DIRECT_OBJECT {
		t.Fatal("Incorrect identification of component name for nested statement's direct object")
	}

}

/*
Tests whether annotated statements (with suffix) are properly decomposed. Case mixes private and shared properties.
*/
func TestCorrectNodeRemovalWithSharedAndPrivateProperties(t *testing.T) {
	text := "A(Operations) I(were (non-compliant [OR] violated)) Bdir,p(proper) Bdir1,p1(organic farming) Bdir1(provisions) and Bdir2,p2(improper) Bdir2(rulesS)"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	stmt := s[0].Entry.(*tree.Statement)

	if stmt.DirectObjectPropertySimple.Entry != "proper" {
		t.Fatal("Shared property was not correctly detected.")
	}

	if stmt.DirectObject.Left.PrivateNodeLinks[0].Entry != "organic farming" {
		t.Fatal("Private node was incorrectly identified")
	}

	if stmt.DirectObject.Right.PrivateNodeLinks[0].Entry != "improper" {
		t.Fatal("Private node was incorrectly identified")
	}

}

/*
Tests whether annotated statements (with suffix) are properly decomposed. Case contains only private properties.
*/
func TestCorrectNodeRemovalWithPrivatePropertiesOnly(t *testing.T) {
	text := "A(Operations) I(were (non-compliant [OR] violated)) Bdir1,p(organic farming) Bdir1(provisions) and Bdir2,p(improper) Bdir2(rules)"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	stmt := s[0].Entry.(*tree.Statement)

	if stmt.DirectObjectPropertySimple != nil {
		t.Fatal("Shared property present, but should be absent. Value:", stmt.DirectObjectPropertySimple)
	}

	if stmt.DirectObject.Left.PrivateNodeLinks[0].Entry != "organic farming" {
		t.Fatal("Private node was incorrectly identified")
	}

	if stmt.DirectObject.Right.PrivateNodeLinks[0].Entry != "improper" {
		t.Fatal("Private node was incorrectly identified")
	}

}

/*
Tests automated extraction of component type as well as property and suffix.
*/
func TestExtractComponentTypeIncludingPropertyAndSuffix(t *testing.T) {

	// Primitive type

	input := "Cac("

	compType, prop, err := extractComponentType(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during component type extraction:", err)
	}
	if compType != "Cac" {
		t.Fatal("Wrong component type detected.")
	}
	if prop != false {
		t.Fatal("Wrong property characteristics detected.")
	}

	// complex type

	input = "Cac{"

	compType, prop, err = extractComponentType(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during component type extraction:", err)
	}
	if compType != "Cac" {
		t.Fatal("Wrong component type detected.")
	}
	if prop != false {
		t.Fatal("Wrong property characteristics detected.")
	}

	// Primitive type with property

	input = "Cac,p("

	compType, prop, err = extractComponentType(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during component type extraction:", err)
	}
	if compType != "Cac,p" {
		t.Fatal("Wrong component type detected. Type:", compType)
	}
	if prop != true {
		t.Fatal("Wrong property characteristics detected.")
	}

	// simple type with suffix

	input = "Cac1("

	compType, prop, err = extractComponentType(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during component type extraction:", err)
	}
	if compType != "Cac" {
		t.Fatal("Wrong component type detected.")
	}
	if prop != false {
		t.Fatal("Wrong property characteristics detected.")
	}

	// simple type with suffix and property

	input = "Cac1,p("

	compType, prop, err = extractComponentType(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during component type extraction:", err)
	}
	if compType != "Cac,p" {
		t.Fatal("Wrong component type detected.")
	}
	if prop != true {
		t.Fatal("Wrong property characteristics detected.")
	}

	// simple type with annotation

	input = "Cac[annotation]("

	compType, prop, err = extractComponentType(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during component type extraction:", err)
	}
	if compType != "Cac" {
		t.Fatal("Wrong component type detected.")
	}
	if prop != false {
		t.Fatal("Wrong property characteristics detected.")
	}

	// simple type with suffix and annotation

	input = "Cac1[annot=test]("

	compType, prop, err = extractComponentType(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during component type extraction:", err)
	}
	if compType != "Cac" {
		t.Fatal("Wrong component type detected.")
	}
	if prop != false {
		t.Fatal("Wrong property characteristics detected.")
	}

	// simple type with suffix and annotation and property

	input = "Cac1,p[annot=test]("

	compType, prop, err = extractComponentType(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during component type extraction:", err)
	}
	if compType != "Cac,p" {
		t.Fatal("Wrong component type detected.")
	}
	if prop != true {
		t.Fatal("Wrong property characteristics detected.")
	}

	// complex type with suffix and annotation

	input = "Cac1[testAnnotation]{"

	compType, prop, err = extractComponentType(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during component type extraction:", err)
	}
	if compType != "Cac" {
		t.Fatal("Wrong component type detected. Component type:", compType)
	}
	if prop != false {
		t.Fatal("Wrong property characteristics detected.")
	}

}

/*
Tests proper suffix and annotation parsing on regular nested statement.
*/
func TestSuffixAndAnnotationParsingOnSimpleNestedStatement(t *testing.T) {

	input := "A(Actor) D(must) I(review) Bdir(subjects) CacA[ctx=state]{A(Supervisor) I(appoints) Bdir(actor)}"

	stmt, err := ParseStatement(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing:", err)
	}

	s := stmt[0].Entry.(*tree.Statement)

	if s.ActivationConditionComplex == nil {
		t.Fatal("Complex nested activation condition was not correctly parsed.")
	}

	if s.ActivationConditionComplex.Suffix != "A" {
		t.Fatal("Suffix on complex nested statement is incorrect:", s.ActivationConditionComplex.Suffix)
	}

	if s.ActivationConditionComplex.Annotations.(string) != "[ctx=state]" {
		t.Fatal("Annotations on complex nested statement are incorrect:", s.ActivationConditionComplex.Annotations.(string))
	}

}

/*
Tests proper parsing of component specifications when parsing nested combinations.
*/
func TestSuffixAndAnnotationParsingOnNestedStatementCombination(t *testing.T) {

	input := "A(Actor) D(must) I(review) Bdir(subjects) " +
		"Cac{CacB[ctx=state]{A(Supervisor) I(appoints) Bdir(actor)} [XOR] " +
		"Cac1[annotA]{A(other actor) I(does) Bdir(something)}}"

	stmt, err := ParseStatement(input)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing:", err)
	}

	s := stmt[0].Entry.(*tree.Statement)

	if s.ActivationConditionComplex == nil {
		t.Fatal("Complex nested activation condition was not correctly parsed.")
	}

	if s.ActivationConditionComplex.LogicalOperator != "XOR" {
		t.Fatal("Logical operator incorrectly identified as", s.ActivationConditionComplex.LogicalOperator)
	}

	if s.ActivationConditionComplex.Left.Suffix != "B" {
		t.Fatal("Suffix on complex nested statement is incorrect:", s.ActivationConditionComplex.Left.Suffix)
	}

	if s.ActivationConditionComplex.Right.Suffix != "1" {
		t.Fatal("Suffix on complex nested statement is incorrect:", s.ActivationConditionComplex.Right.Suffix)
	}

	if s.ActivationConditionComplex.Left.Annotations.(string) != "[ctx=state]" {
		t.Fatal("Annotations on complex nested statement are incorrect:", s.ActivationConditionComplex.Left.Annotations.(string))
	}

	if s.ActivationConditionComplex.Right.Annotations.(string) != "[annotA]" {
		t.Fatal("Annotations on complex nested statement are incorrect:", s.ActivationConditionComplex.Right.Annotations.(string))
	}

}

/*
Tests detection of conflicting component specifications when parsing nested combinations.
*/
func TestMultipleComponentSpecificationWhenParsingOnNestedStatementCombination(t *testing.T) {

	// Parsing of this statement should fail, because CacA specification contains multiple component symbols
	input := "A(Actor) D(must) I(review) Bdir(subjects) " +
		"Cac{CacA[ctx=state]{A(Supervisor) I(appoints) Bdir(actor)} [XOR] " +
		"Cac1[annotA]{A(other actor) I(does) Bdir(something)}}"

	_, err := ParseStatement(input)
	if err.ErrorCode != tree.PARSING_ERROR_MULTIPLE_COMPONENTS_FOUND {
		t.Fatal("Error during parsing:", err)
	}
}

/*
Tests correct handling of invalid combination of nested statements.
*/
func TestInvalidNestedCombination(t *testing.T) {

	text := "Cac{Cac{A(actor) I(action)} [XOR] Cac(simple content)}"

	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_INVALID_COMBINATION {
		t.Fatal("Parsing did not pick up on erroneous combination.")
	}

	text = "Cac{Cac(simple content) [OR] Cac{A(actor) I(action)} [OR] Cac{A(actor2) I(action2)}}"

	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_INVALID_COMBINATION {
		t.Fatal("Parsing did not pick up on erroneous combination.")
	}

}

/*
Tests whether presence of within-linkage (wAND) is correctly detected.
*/
func TestHasWithinLinkage(t *testing.T) {

	text := "Cex(for compliance with (left [AND] right) as well as (left1 [XOR] right1) shared) Cex(outlier)"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing:", err)
	}

	res := s[0].Entry.(*tree.Statement)

	nodes := res.ExecutionConstraintSimple.GetLeafNodes(true)

	fmt.Println("Nodes:", nodes)

	if !nodes[0][0].HasWithinComponentLinkage() {
		t.Fatal("Node", nodes[0][0], "should have within-component linkage.")
	}

	if !nodes[0][1].HasWithinComponentLinkage() {
		t.Fatal("Node", nodes[0][1], "should have within-component linkage.")
	}

	if !nodes[0][2].HasWithinComponentLinkage() {
		t.Fatal("Node", nodes[0][2], "should have within-component linkage.")
	}

	if !nodes[0][3].HasWithinComponentLinkage() {
		t.Fatal("Node", nodes[0][3], "should have within-component linkage.")
	}

	if nodes[0][4].HasWithinComponentLinkage() {
		t.Fatal("Node", nodes[0][4], "should NOT have within-component linkage.")
	}

}

/*
Tests for the correct resolution of between- (bAND) and within-linked (wAND) nodes.
*/
func TestGetSyntheticRootNode(t *testing.T) {

	text := "Cex(for compliance with (left [AND] right) as well as (left1 [XOR] right1) shared) Cex(outlier)"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing:", err)
	}

	res := s[0].Entry.(*tree.Statement)

	nodes := res.ExecutionConstraintSimple.GetLeafNodes(true)

	fmt.Println("Nodes:", nodes)

	if nodes[0][0].GetRootNode() != res.ExecutionConstraintSimple {
		t.Fatal("Node", nodes[0][0], "should have have top-level root node.")
	}

	if nodes[0][0].GetNodeBelowSyntheticRootNode() != nodes[0][0].Parent {
		t.Fatal("Node", nodes[0][0], "has an incorrect implicit within-linked synthetic root node.")
	}

	if nodes[0][1].GetNodeBelowSyntheticRootNode() != nodes[0][1].Parent {
		t.Fatal("Node", nodes[0][1], "has an incorrect implicit within-linked synthetic root node.")
	}

	if nodes[0][2].GetNodeBelowSyntheticRootNode() != nodes[0][2].Parent {
		t.Fatal("Node", nodes[0][2], "has an incorrect implicit within-linked synthetic root node.")
	}

	if nodes[0][3].GetNodeBelowSyntheticRootNode() != nodes[0][3].Parent {
		t.Fatal("Node", nodes[0][3], "has an incorrect implicit within-linked synthetic root node.")
	}

	if nodes[0][4].GetNodeBelowSyntheticRootNode() != res.ExecutionConstraintSimple.Right {
		t.Fatal("Node", nodes[0][4], "should have have implicit between-linked synthetic root node.")
	}

}

/*
Tests the resolution of components and the associated properties based on the IGStatement's GetProperties() function.
*/
func TestComponentPropertyRelationshipResolution(t *testing.T) {

	text := "A,p(Qualified) A(actor) I(does activity) on Bdir,p(qualified) " +
		"Bdir,p{A(actor) I(does) Bdir(something)} Bdir(something) Bind(to someone) " +
		"E,p(some) E(definiendum) F(is defined as) " +
		// Provide multiple properties - mix of primitive and complex ones
		"P,p(this) and P,p(that) P,p{E(something) F(is established) Cex(before)} P(definiens)"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing:", err)
	}

	res := s[0].Entry.(*tree.Statement)

	if len(res.GetPropertyComponent(res.Attributes, false)) != 1 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.Attributes, false)))
	}

	if len(res.GetPropertyComponent(res.Attributes, true)) != 1 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.Attributes, true)))
	}

	if res.GetPropertyComponent(res.Attributes, false)[0] != res.AttributesPropertySimple {
		t.Fatal("Failed to identify component-associated properties.")
	}

	if len(res.GetPropertyComponent(res.DirectObject, false)) != 1 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.DirectObject, false)))
	}

	if len(res.GetPropertyComponent(res.DirectObject, true)) != 2 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.DirectObject, false)))
	}

	if res.GetPropertyComponent(res.DirectObject, false)[0] != res.DirectObjectPropertySimple {
		t.Fatal("Failed to identify component-associated properties.")
	}

	if len(res.GetPropertyComponent(res.IndirectObject, false)) != 0 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.IndirectObject, false)))
	}

	if len(res.GetPropertyComponent(res.IndirectObject, true)) != 0 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.IndirectObject, false)))
	}

	// Test component that does not have property ever
	if len(res.GetPropertyComponent(res.Deontic, false)) != 0 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.Deontic, false)))
	}

	if len(res.GetPropertyComponent(res.Deontic, true)) != 0 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.Deontic, false)))
	}

	if len(res.GetPropertyComponent(res.ConstitutedEntity, false)) != 1 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.ConstitutedEntity, false)))
	}

	if len(res.GetPropertyComponent(res.ConstitutedEntity, true)) != 1 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.ConstitutedEntity, true)))
	}

	// Test components that have combined properties - will only return one implicitly linked entry
	if len(res.GetPropertyComponent(res.ConstitutingProperties, false)) != 1 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.ConstitutingProperties, false)))
	}

	if len(res.GetPropertyComponent(res.ConstitutingProperties, true)) != 2 {
		t.Fatal("Wrong properties count for given component: ", len(res.GetPropertyComponent(res.ConstitutingProperties, true)))
	}

}

/*
Tests the proper error response when parsing nested statements (leading to the nested elements being ignored).
Here, this is provoked by excessive braces in input.
*/
func TestNestedStatementParsingError(t *testing.T) {

	// Statement with excess braces around first activation condition element (i.e, Cac{{ <-- should only be one)
	text := "{Cac{{when the A(Program Manager) I(believes) that Bdir{a A,p(certified) A(operation) I((has violated [OR] is not in compliance))" +
		" Bdir(with (the Act [OR] regulations in this part))}}}, [OR] Cac{when a A((certifying agent [OR] State organic program’s governing State official)) I(fails to enforce) Bdir((the Act [OR] regulations in this part)).}}"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS {
		t.Fatal("Parsing should have caused error "+tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS+", but returned error:", err)
	}

	res := s[0].Entry.(*tree.Statement)

	if err.ErrorIgnoredElements == nil || len(err.ErrorIgnoredElements) == 0 ||
		strings.Join(err.ErrorIgnoredElements[:], ",") != "{when the A(Program Manager) I(believes) that Bdir{a A,p(certified) A(operation) I((has violated [OR] is not in compliance)) Bdir(with (the Act [OR] regulations in this part))}}" {
		t.Fatal("The error should contain ignored elements, but contained: " + strings.Join(err.ErrorIgnoredElements[:], ","))
	}

	if res.String() != "" {
		t.Fatal("Returned statements should be empty, but is ", res.String())
	}
}

/*
Tests complexity calculation of parsed statements.
*/
func TestStatement_CalculateComplexity(t *testing.T) {

	testStmt := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that Cac{Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate]((were non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)}}."

	s, err := ParseStatement(testStmt)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing of statement should not have failed. Error:", err)
	}

	stmt := s[0].Entry.(*tree.Statement)

	// Calculate actual complexity
	complexity := stmt.CalculateComplexity()

	// Evaluate

	if complexity.AttributesOptions != 1 {
		t.Fatal("Attributes options are incorrect. Value: ", complexity.AttributesOptions)
	}

	if complexity.AttributesComplexity != 1 {
		t.Fatal("Attributes complexity is incorrect. Value: ", complexity.AttributesComplexity)
	}

	if complexity.AttributesPropertySimpleOptions != 1 {
		t.Fatal("Attributes Properties options are incorrect. Value: ", complexity.AttributesPropertySimpleOptions)
	}

	if complexity.AttributesPropertySimpleComplexity != 1 {
		t.Fatal("Attributes Properties complexity is incorrect. Value: ", complexity.AttributesPropertySimpleComplexity)
	}

	if complexity.DeonticOptions != 1 {
		t.Fatal("Deontic options are incorrect. Value: ", complexity.DeonticOptions)
	}

	if complexity.DeonticComplexity != 1 {
		t.Fatal("Deontic complexity is incorrect. Value: ", complexity.DeonticComplexity)
	}

	if complexity.AimOptions != 3 {
		t.Fatal("Aim options are incorrect. Value: ", complexity.AimOptions)
	}

	if complexity.AimComplexity != 2 {
		t.Fatal("Aim complexity is incorrect. Value: ", complexity.AimComplexity)
	}

	if complexity.DirectObjectSimpleOptions != 3 {
		t.Fatal("Direct Object options are incorrect. Value: ", complexity.DirectObjectSimpleOptions)
	}

	if complexity.DirectObjectSimpleComplexity != 1 {
		t.Fatal("Direct Object complexity is incorrect. Value: ", complexity.DirectObjectSimpleComplexity)
	}

	if complexity.DirectObjectPropertySimpleOptions != 1 {
		t.Fatal("Direct Object Property options are incorrect. Value: ", complexity.DirectObjectPropertySimpleOptions)
	}

	if complexity.DirectObjectPropertySimpleComplexity != 1 {
		t.Fatal("Direct Object Property complexity is incorrect. Value: ", complexity.DirectObjectPropertySimpleComplexity)
	}

	if complexity.ExecutionConstraintSimpleOptions != 3 {
		t.Fatal("Execution constraints options are incorrect. Value: ", complexity.ExecutionConstraintSimpleOptions)
	}

	if complexity.ExecutionConstraintSimpleComplexity != 2 {
		t.Fatal("Execution constraints complexity is incorrect. Value: ", complexity.ExecutionConstraintSimpleComplexity)
	}

	if complexity.ActivationConditionComplexComplexity != 3 {
		t.Fatal("Activation conditions complexity is incorrect. Value: ", complexity.ActivationConditionComplexComplexity)
	}

	if complexity.ActivationConditionComplexOptions != 2 {
		t.Fatal("Activation conditions options are incorrect. Value: ", complexity.ActivationConditionComplexOptions)
	}

	if complexity.ActivationConditionComplexComplexity != 3 {
		t.Fatal("Activation conditions complexity is incorrect. Value: ", complexity.ActivationConditionComplexComplexity)
	}

	if complexity.TotalStateComplexity != 12 {
		t.Fatal("Total State Complexity is incorrect. Value: ", complexity.TotalStateComplexity)
	}

}

/*
Tests the parsing of component pair combinations into tree structure.
*/
func TestComponentPairCombinationTreeParsing(t *testing.T) {
	// Moderately complex statement with nested component on left
	text := "{ Cac{A(precond)} Bdir(leftbdir) I(leftact) [XOR] Bdir(rightbdir) I(rightact)}"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing returned error:", err)
	}

	if len(s) > 1 {
		t.Fatal("Found more than one node:", len(s))
	}

	res := s[0]

	// Logical operator

	val1 := res.LogicalOperator

	if val1 != "XOR" {
		t.Fatal("Component has not been correctly parsed. Value: ", val1)
	}

	// Shared elements

	if res.GetSharedLeft() != nil {
		t.Fatal("Shared elements should be nil, but are", res.GetSharedLeft())
	}

	if res.GetSharedRight() != nil {
		t.Fatal("Shared elements should be nil, but are", res.GetSharedRight())
	}

	// Test left side

	// Shared elements

	if res.Left.GetSharedLeft() != nil {
		t.Fatal("Shared elements should be nil, but are", res.Left.GetSharedLeft())
	}

	if res.Left.GetSharedRight() != nil {
		t.Fatal("Shared elements should be nil, but are", res.Left.GetSharedRight())
	}

	// Components

	val := res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Entry

	if val != "leftact" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry

	if val != "leftbdir" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).Attributes.Entry

	if val != "precond" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Right side

	// Shared elements

	if res.Right.GetSharedLeft() != nil {
		t.Fatal("Shared elements should be nil, but are", res.Right.GetSharedLeft())
	}

	if res.Right.GetSharedRight() != nil {
		t.Fatal("Shared elements should be nil, but are", res.Right.GetSharedRight())
	}

	// Components

	val = res.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Entry

	if val != "rightact" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry

	if val != "rightbdir" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

}

/*
Tests component pair combinations with private properties
*/
func TestComponentPairCombinationTreeParsingPrivateComponents(t *testing.T) {
	// Statement with left private property
	text := "{ Bdir,p(privateleft) Bdir(leftbdir) [AND] Bdir(rightbdir)}"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing returned error:", err)
	}

	if len(s) > 1 {
		t.Fatal("Found more than one node:", len(s))
	}

	res := s[0]

	// Logical operator

	val1 := res.LogicalOperator

	if val1 != "AND" {
		t.Fatal("Component has not been correctly parsed. Value: ", val1)
	}

	// Shared elements

	if res.GetSharedLeft() != nil {
		t.Fatal("Shared elements should be nil, but are", res.GetSharedLeft())
	}

	if res.GetSharedRight() != nil {
		t.Fatal("Shared elements should be nil, but are", res.GetSharedRight())
	}

	// Test left side

	val := res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObjectPropertySimple.Entry

	if val != "privateleft" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry

	if val != "leftbdir" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Right side

	val = res.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry

	if val != "rightbdir" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

}

/*
Tests complex component pair example with atomic components.
Implicitly tests CopyComponent function which needs to reconcile populated component fields.
*/
func TestComponentPairCombinationTreeParsingWithAtomicComponents(t *testing.T) {

	text := "D(deontic) Cac(atomicCondition) (lkjsdkljs) Bind(indirectobject) Cac{A(atomicnestedcondition)} " +
		"{I(maintain) Bdir((order [AND] control))  Cac{A(sharednestedcondition)} [XOR] {I(sustain) Bdir(peace) [OR] I(prevent) Bdir(war)}} "

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing returned error:", err)
	}

	if len(s) > 1 {
		t.Fatal("Found more than one node:", len(s))
	}

	res := s[0]

	// Logical operator

	val1 := res.LogicalOperator

	if val1 != "XOR" {
		t.Fatal("Component has not been correctly parsed. Value: ", val1)
	}

	// Shared elements

	if res.GetSharedLeft() != nil {
		t.Fatal("Shared elements should be nil, but are", res.GetSharedLeft())
	}

	if res.GetSharedRight() != nil {
		t.Fatal("Shared elements should be nil, but are", res.GetSharedRight())
	}

	// Test left side

	val := res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Deontic.Entry

	if val != "deontic" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Entry

	if val != "maintain" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Nested combination

	val2 := res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.LogicalOperator

	if val2 != "AND" {
		t.Fatal("Component has not been correctly parsed. Value: ", val2)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry

	if val != nil {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Left.Entry

	if val != "order" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Right.Entry

	if val != "control" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Indirect object

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).IndirectObject.Entry

	if val != "indirectobject" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Activation condition

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionSimple.Entry

	if val != "atomicCondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Nested shared condition (combined with individual nested condition) - tests aggregation into node

	val5 := res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.LogicalOperator

	if val5 != "bAND" {
		t.Fatal("Component has not been correctly parsed. Value: ", val5)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Entry.(*tree.Statement).Attributes.Entry

	if val != "sharednestedcondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Right.Entry.(*tree.Statement).Attributes.Entry

	if val != "atomicnestedcondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Right side

	// nested logical operator

	val4 := res.Right.LogicalOperator

	if val4 != "OR" {
		t.Fatal("Component has not been correctly parsed. Value: ", val4)
	}

	// Right left nested

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Deontic.Entry

	if val != "deontic" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Entry

	if val != "sustain" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry

	if val != "peace" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).IndirectObject.Entry

	if val != "indirectobject" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// injected atomic condition

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionSimple.Entry

	if val != "atomicCondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// injected atomic nested condition

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).Attributes.Entry

	if val != "atomicnestedcondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Right right nested

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Deontic.Entry

	if val != "deontic" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Entry

	if val != "prevent" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry

	if val != "war" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).IndirectObject.Entry

	if val != "indirectobject" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// injected atomic condition

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionSimple.Entry

	if val != "atomicCondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// injected atomic nested condition

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Entry.(*tree.Statement).Attributes.Entry

	if val != "atomicnestedcondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

}

/*
Tests complex component pair example with atomic components and component combinations (combines all statement classes).
Implicitly tests CopyComponent function which needs to reconcile populated component fields.
*/
func TestComponentPairCombinationTreeParsingWithAtomicComponentsAndComponentCombinations(t *testing.T) {

	text := "D(deontic) Cac(atomicCondition) (lkjsdkljs) Bind(indirectobject) Cac{A(atomicnestedcondition)} " +
		"{I(maintain) Bdir((order [AND] control))  Cac{A(sharednestedcondition)} [XOR] {I(sustain) Bdir(peace) [OR] I(prevent) Bdir(war)}} " +
		" Cac{Cac{ A(leftcombo) I(leftaim) } [XOR] Cac{ Cac{ A(rightleftcombo) I(rightleftaim) } [AND] Cac{ A(rightrightcombo) I(rightrightaim) }}}"

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing returned error:", err)
	}

	if len(s) > 1 {
		t.Fatal("Found more than one node:", len(s))
	}

	res := s[0]

	// Logical operator

	val1 := res.LogicalOperator

	if val1 != "XOR" {
		t.Fatal("Component has not been correctly parsed. Value: ", val1)
	}

	// Shared elements

	if res.GetSharedLeft() != nil {
		t.Fatal("Shared elements should be nil, but are", res.GetSharedLeft())
	}

	if res.GetSharedRight() != nil {
		t.Fatal("Shared elements should be nil, but are", res.GetSharedRight())
	}

	// Test left side

	val := res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Deontic.Entry

	if val != "deontic" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Entry

	if val != "maintain" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Nested combination

	val2 := res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.LogicalOperator

	if val2 != "AND" {
		t.Fatal("Component has not been correctly parsed. Value: ", val2)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry

	if val != nil {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Left.Entry

	if val != "order" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Right.Entry

	if val != "control" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Indirect object

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).IndirectObject.Entry

	if val != "indirectobject" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Activation condition

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionSimple.Entry

	if val != "atomicCondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Nested shared condition (combined with individual nested condition) - tests aggregation into node

	val5 := res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.LogicalOperator

	if val5 != "bAND" {
		t.Fatal("Component has not been correctly parsed. Value: ", val5)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Entry.(*tree.Statement).Attributes.Entry

	if val != "sharednestedcondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Right.Left.Left.Entry.(*tree.Statement).Attributes.Entry

	if val != "leftcombo" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Right.Left.Left.Entry.(*tree.Statement).Aim.Entry

	if val != "leftaim" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Nested activation condition right, left, right branch

	val6 := res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Right.Left.Right.LogicalOperator

	if val6 != "AND" {
		t.Fatal("Component has not been correctly parsed. Value: ", val6)
	}

	// left subbranch

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Right.Left.Right.Left.Entry.(*tree.Statement).Attributes.Entry

	if val != "rightleftcombo" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Right.Left.Right.Left.Entry.(*tree.Statement).Aim.Entry

	if val != "rightleftaim" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// right subbranch

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Right.Left.Right.Right.Entry.(*tree.Statement).Attributes.Entry

	if val != "rightrightcombo" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Right.Left.Right.Right.Entry.(*tree.Statement).Aim.Entry

	if val != "rightrightaim" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Right side

	// nested logical operator

	val4 := res.Right.LogicalOperator

	if val4 != "OR" {
		t.Fatal("Component has not been correctly parsed. Value: ", val4)
	}

	// Right left nested

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Deontic.Entry

	if val != "deontic" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Entry

	if val != "sustain" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry

	if val != "peace" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).IndirectObject.Entry

	if val != "indirectobject" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// injected atomic condition

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionSimple.Entry

	if val != "atomicCondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// combined activation condition on right left left side

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Left.Entry.(*tree.Statement).Attributes.Entry

	if val != "leftcombo" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Left.Entry.(*tree.Statement).Aim.Entry

	if val != "leftaim" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// ... right left right side

	val7 := res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Right.LogicalOperator

	if val7 != "AND" {
		t.Fatal("Component has not been correctly parsed. Value: ", val7)
	}

	// left subbranch

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Right.Left.Entry.(*tree.Statement).Attributes.Entry

	if val != "rightleftcombo" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Right.Left.Entry.(*tree.Statement).Aim.Entry

	if val != "rightleftaim" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// right subbranch

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Right.Right.Entry.(*tree.Statement).Attributes.Entry

	if val != "rightrightcombo" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Left.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Right.Right.Entry.(*tree.Statement).Aim.Entry

	if val != "rightrightaim" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// Right right nested

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Deontic.Entry

	if val != "deontic" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).Aim.Entry

	if val != "prevent" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).DirectObject.Entry

	if val != "war" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).IndirectObject.Entry

	if val != "indirectobject" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// injected atomic condition

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionSimple.Entry

	if val != "atomicCondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// inject component combination

	val8 := res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.LogicalOperator

	if val8 != "AND" {
		t.Fatal("Component has not been correctly parsed. Value: ", val8)
	}

	val9 := res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.LogicalOperator

	if val9 != "XOR" {
		t.Fatal("Component has not been correctly parsed. Value: ", val9)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Left.Entry.(*tree.Statement).Attributes.Entry

	if val != "leftcombo" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Left.Entry.(*tree.Statement).Aim.Entry

	if val != "leftaim" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val10 := res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Right.LogicalOperator

	if val10 != "AND" {
		t.Fatal("Component has not been correctly parsed. Value: ", val10)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Right.Left.Entry.(*tree.Statement).Attributes.Entry

	if val != "rightleftcombo" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Right.Left.Entry.(*tree.Statement).Aim.Entry

	if val != "rightleftaim" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Right.Right.Entry.(*tree.Statement).Attributes.Entry

	if val != "rightrightcombo" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Left.Right.Right.Entry.(*tree.Statement).Aim.Entry

	if val != "rightrightaim" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

	// individual condition on right side

	val = res.Right.Right.Entry.([]*tree.Node)[0].Entry.(*tree.Statement).ActivationConditionComplex.Right.Entry.(*tree.Statement).Attributes.Entry

	if val != "atomicnestedcondition" {
		t.Fatal("Component has not been correctly parsed. Value: ", val)
	}

}

/*
Tests multiple component pairs on a given nesting level. Should lead to error message.
*/
func TestComponentPairCombinationTreeParsingMultipleComponentPairs(t *testing.T) {

	text := " A(actor) D(may) {I(leftAim) Bdir(leftObject) [OR] I(rightAim) Bdir(rightObject)} {Cac(leftCondition) Cex(leftConstraint) [XOR] Cac(rightCondition) Cex(rightConstraint)} "

	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_MULTIPLE_COMPONENT_PAIRS_ON_SAME_LEVEL {
		t.Fatal("Parsing should have returned error identifying multiple component pairs, but returned error ", err)
	}
}

/*
Tests single component pair on different nesting levels. Should parse.
*/
func TestComponentPairCombinationTreeParsingComponentPairOnDifferentLevels(t *testing.T) {

	text := " A(actor) D(may) {I(leftAim) Bdir(leftObject) [OR] I(rightAim) Bdir(rightObject)} Cac{ {A(actor2) I(aim2) [XOR] A(actor3) I(aim3)} } "

	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should not have returned error, but returned error ", err)
	}
}

/*
Tests the automated expansion of component combinations (e.g., 'Bdir(left [XOR] right)' to 'Bdir((left [XOR] right))') to support the coding.
*/
func TestAutomatedExpansionOfParenthesesForComponentCombinations(t *testing.T) {

	// Combinations in statement explicitly miss inner component parentheses
	text := " A(actor) D(may) I(leftAim [XOR] rightAim) Bdir((leftObject [XOR] middleObject) [AND] rightObject) Cex(constraint) "

	s, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should not have returned error, but returned error ", err)
	}

	fmt.Println(s)

	stmt := s[0].Entry.(*tree.Statement)

	fmt.Println(stmt)

	// Check Aim

	if !stmt.Aim.IsCombination() {
		t.Fatal("Aim should be combination")
	}

	if stmt.Aim.LogicalOperator != "XOR" {
		t.Fatal("Aim should have logical operator XOR, but has ", stmt.Aim.LogicalOperator)
	}

	if stmt.Aim.Left.Entry != "leftAim" {
		t.Fatal("Aim should have left value leftAim, but has ", stmt.Aim.Left.Entry)
	}

	if stmt.Aim.Right.Entry != "rightAim" {
		t.Fatal("Aim should have right value rightAim, but has ", stmt.Aim.Right.Entry)
	}

	// Check Object

	if !stmt.DirectObject.IsCombination() {
		t.Fatal("Direct object should be combination")
	}

	if stmt.DirectObject.LogicalOperator != "AND" {
		t.Fatal("Direct object component combination should have logical operator AND, but has ", stmt.DirectObject.LogicalOperator)
	}

	if stmt.DirectObject.Left.LogicalOperator != "XOR" {
		t.Fatal("Direct object left nested component combination should have logical operator XOR, but has ", stmt.DirectObject.Left.LogicalOperator)
	}

	if stmt.DirectObject.Left.Left.Entry != "leftObject" {
		t.Fatal("Direct object left left component should have value leftObject, but has ", stmt.DirectObject.Left.Left.Entry)
	}

	if stmt.DirectObject.Left.Right.Entry != "middleObject" {
		t.Fatal("Direct object left right component should have value middleObject, but has ", stmt.DirectObject.Left.Right.Entry)
	}

	if stmt.DirectObject.Right.Entry != "rightObject" {
		t.Fatal("Direct object right component should have value rightObject, but has ", stmt.DirectObject.Right.Entry)
	}
}

/*
Tests the separation of input strings into statement patterns (basic components including combinations, nested component, nested statement combinations and component pair combinations).

Types:
- Components (and combinations): A(actor); A((actor1 [XOR] actor2))
- Component nesting syntax: Cac{ A(actor) I(action) }
- Component combination syntax: Cac{ Cac{A(leftNestedA) I(leftNestedI)} [XOR] Cac{A(rightNestedA) I(rightNestedI)} }
- Component pair combination syntax: { Cac{A(leftNestedA) I(leftNestedI)} [XOR] Cac{A(rightNestedA) I(rightNestedI)} }
*/
func TestSeparateComponentsNestedStatementsCombinationsAndComponentPairs(t *testing.T) {

	text := "D(deontic) Cac(atomicCondition) (meaningless content) Bind(indirectobject) Cac{A(atomicnestedcondition)} " +
		"{I(maintain) Bdir((order [AND] control))  Cac{A(sharednestedcondition)} [XOR] {I(sustain) Bdir(peace) [OR] I(prevent) Bdir(war)}} " +
		" Cac{Cac{ A(leftcombo) I(leftaim) } [XOR] Cac{ Cac{ A(rightleftcombo) I(rightleftaim) } [AND] Cac{ A(rightrightcombo) I(rightrightaim) }}} " +
		" Bdir(left [XOR] right) A((actor1 [AND] actor2)) Bdir{ somecontent } " +
		" Cex{ Cex{ A(anotherLeft) I(anotherAim) } [XOR] Cex{ I(anotherRightAim) Cex(embeddedCex) }} " +
		" {A(actor1) I(aim1) [XOR] {A(actor2) I(aim2) [AND] A(actor3) I(aim3)}}"

	types, err := separateComponentsNestedStatementsCombinationsAndComponentPairs(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing. Error: ", err)
	}

	// Basic statements (note: includes all whitespaces not fitting elsewhere - filtered in downstream processing)
	if types[0][0] != "D(deontic) Cac(atomicCondition) (meaningless content) Bind(indirectobject)      Bdir(left [XOR] right) A((actor1 [AND] actor2))     " {
		t.Fatal("Wrong entry in basic components")
	}
	// nested components
	if types[1][0] != "Cac{A(atomicnestedcondition)}" {
		t.Fatal("Wrong entry in nested components")
	}
	if types[1][1] != "Bdir{ somecontent }" {
		t.Fatal("Wrong entry in nested components")
	}
	// nested statement combinations
	if types[2][0] != "Cac{Cac{ A(leftcombo) I(leftaim) } [XOR] Cac{ Cac{ A(rightleftcombo) I(rightleftaim) } [AND] Cac{ A(rightrightcombo) I(rightrightaim) }}}" {
		t.Fatal("Wrong entry in nested statement combinations")
	}
	if types[2][1] != "Cex{ Cex{ A(anotherLeft) I(anotherAim) } [XOR] Cex{ I(anotherRightAim) Cex(embeddedCex) }}" {
		t.Fatal("Wrong entry in nested statement combinations")
	}
	// component pair combinations
	if types[3][0] != "{I(maintain) Bdir((order [AND] control))  Cac{A(sharednestedcondition)} [XOR] {I(sustain) Bdir(peace) [OR] I(prevent) Bdir(war)}}" {
		t.Fatal("Wrong entry in component pair combinations")
	}
	if types[3][1] != "{A(actor1) I(aim1) [XOR] {A(actor2) I(aim2) [AND] A(actor3) I(aim3)}}" {
		t.Fatal("Wrong entry in component pair combinations")
	}

}

/*
Tests input with correct parentheses/braces counts, but wrong order (i.e., not matching).
Tests the search exhaustion in #GetComponentContent().
*/
func TestWrongParenthesisOrderButCorrectCount(t *testing.T) {

	// Problem area is Bdir,p and Bdir (parentheses)
	text := "A(actor) Bdir,p(left [AND] right)) Bdir((left [AND] right) Cac(condition) something else"

	// Test for error during parsing. Should pick up on ordering problem.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_UNABLE_TO_EXTRACT_COMPONENT_CONTENT {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_UNABLE_TO_EXTRACT_COMPONENT_CONTENT+", but returned error ", err)
	}

}

/*
Tests duplicate components in input (detection of duplicate statements).
*/
func TestDuplicateComponentsInStatement(t *testing.T) {

	// Input text with repeated components (copy-and-paste problem)
	text := "A,p(relevant) A(regulators) D(must) I(monitor [AND] enforce) Bdir(compliance). " +
		"A,p(relevant) A(regulators) D(must) I(monitor [AND] enforce) Bdir(compliance)."

	// Test for error during parsing. Should pick up on duplicates.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_DUPLICATE_COMPONENT_ENTRIES {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_DUPLICATE_COMPONENT_ENTRIES+", but returned error ", err)
	}

}

/*
Tests for missing or incorrect logical operator on component pair level but operator contained in
embedded component property combination
*/
func TestInvalidComponentPairCombinationWithMissingLogicalOperatorButNestedComponentPropertyCombinations(t *testing.T) {

	// Input text with missing logical operator on component pair level but component combination in one pair
	text := "{I(prioritise) Bdir,p(swift [AND] predictable) Bdir(emission reductions) AND I(enhance) Bdir(removals) Cex(by natural sinks)}"

	// Test for error during parsing. Should pick up on missing component pair.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_INVALID_COMPONENT_PAIR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_INVALID_COMPONENT_PAIR+", but returned error ", err)
	}
}

/*
Tests for missing or incorrect logical operator on component pair level but containing component combination
(will be initially mistaken as nested statement, but caught during deep parsing).
*/
func TestInvalidComponentPairCombinationWithMissingLogicalOperatorButComponentCombination(t *testing.T) {

	// Input text with missing logical operator on component pair level
	text := "A(actor) {I(action1 [XOR] action2) Bdir(object1) and I(action3) Bdir(object2)}"

	// Test for error during parsing. Should pick up on missing component pair.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS+", but returned error ", err)
	}
}

/*
Tests for missing or incorrect logical operator on component pair level (and no further logical operator in pair)
(will be initially mistaken as nested statement, but caught during deep parsing).
*/
func TestInvalidComponentPairCombinationWithMissingLogicalOperator(t *testing.T) {

	// Input text with missing logical operator on component pair level
	text := "A(actor) {I(action1) Bdir(object1) and I(action3) Bdir(object2)}"

	// Test for error during parsing. Should pick up on missing component pair.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS+", but returned error ", err)
	}
}

/*
Tests for empty statement
*/
func TestEmptyStatement(t *testing.T) {

	// Empty input text
	text := ""

	// Test for error during parsing. Should pick up on empty input.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_EMPTY_STATEMENT {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_EMPTY_STATEMENT+", but returned error ", err)
	}
}

/*
Tests for empty statement with whitespace
*/
func TestEmptyStatementWhitespace(t *testing.T) {

	// Empty input text
	text := " "

	// Test for error during parsing. Should pick up on empty input.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_EMPTY_STATEMENT {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_EMPTY_STATEMENT+", but returned error ", err)
	}
}

/*
Tests for uncoded statement
*/
func TestEmptyUncodedText(t *testing.T) {

	// Uncoded input text
	text := "Once policy comes into force, relevant regulators must monitor and enforce compliance."

	// Test for error during parsing. Should pick up on empty input.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_EMPTY_STATEMENT {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_EMPTY_STATEMENT+", but returned error ", err)
	}
}

/*
Tests for empty statement in for all single components.
*/
func TestEmptyStatementAgainstSingleComponents(t *testing.T) {

	// Input text with single component
	text := "A(value)"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "A,p(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "D(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "I(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "Bdir(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "Bdir,p(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "Bind(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "Bind,p(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "Cac(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "Cex(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "E(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "E,p(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "M(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "F(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "P(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

	// Input text with single component
	text = "P,p(value)"

	// Test for error during parsing. Should not throw error.
	_, err = ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}

}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentCac(t *testing.T) {

	// Input text with complex component
	text := "Cac{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentCex(t *testing.T) {

	// Input text with complex component
	text := "Cex{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentBdir(t *testing.T) {

	// Input text with complex component
	text := "Bdir{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentBdirp(t *testing.T) {

	// Input text with complex component
	text := "Bdir,p{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstInvalidComplexComponentA(t *testing.T) {

	// Input text with complex component
	text := "A{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentAp(t *testing.T) {

	// Input text with complex component
	text := "A,p{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentBind(t *testing.T) {

	// Input text with complex component
	text := "Bind{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentBindp(t *testing.T) {

	// Input text with complex component
	text := "Bind,p{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstInvalidComplexComponentE(t *testing.T) {

	// Input text with complex component
	text := "E{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentEp(t *testing.T) {

	// Input text with complex component
	text := "E,p{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentP(t *testing.T) {

	// Input text with complex component
	text := "P{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentPp(t *testing.T) {

	// Input text with complex component
	text := "P,p{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in complex component.
*/
func TestEmptyStatementAgainstComplexComponentO(t *testing.T) {

	// Input text with complex component
	text := "O{ A(actor) D(must) I(sanction) }"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in basic regulative statement.
*/
func TestEmptyStatementAgainstInvalidComponent(t *testing.T) {

	// Input text with single component
	text := "O(value)"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_EMPTY_STATEMENT {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_ERROR_EMPTY_STATEMENT+", but returned error ", err)
	}
}

/*
Tests for empty statement in basic regulative statement.
*/
func TestEmptyStatementAgainstBasicRegulativeStatement(t *testing.T) {

	// Input text with basic statement
	text := "A,p(relevant) A(regulators) D(must) I(monitor [AND] enforce) Bdir(compliance). "

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in basic constitutive statement.
*/
func TestEmptyStatementAgainstBasicConstitutiveStatement(t *testing.T) {

	// Input text with basic statement
	text := "E(regulators) F(have the right) to P(monitor [AND] enforce). "

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in regulative statement with complex component.
*/
func TestEmptyStatementAgainstStatementWithNestedComponent(t *testing.T) {

	// Input text with complex component
	text := "A(actor) D(must) I(act) Cac{Once E(policy) F(comes into force)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}

/*
Tests for empty statement in component pair.
*/
func TestEmptyStatementAgainstComponentPair(t *testing.T) {

	// Input text with missing logical operator on component pair level
	text := "A(actor) {I(action1) Bdir(object1) [AND] I(action3) Bdir(object2)}"

	// Test for error during parsing. Should not throw error.
	_, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing should have returned error "+
			tree.PARSING_NO_ERROR+", but returned error ", err)
	}
}
