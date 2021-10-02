package parser

import (
	"IG-Parser/tree"
	"fmt"
	"strings"
	"testing"
)

func TestStatementParsingIncludingSyntheticANDs(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(as well as (review [AND] (audit [AND] challenge))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	s, err := ParseStatement(text)

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

	nodeArray, componentIdx := s.GenerateLeafArrays()

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

	nodeArray, componentIdx := s.GenerateLeafArrays()

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

	// TODO: Check for the need to consider SAND_WITHIN_COMPONENTS
	newRoot.LogicalOperator = tree.SAND_BETWEEN_COMPONENTS

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

	s, err := ParseStatement(text)

	fmt.Println(s.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	nodeArray, componentIdx := s.GenerateLeafArrays()

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

	s, err := ParseStatement(text)

	fmt.Println(s.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	nodeArray, componentIdx := s.GenerateLeafArrays()

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

	if element[0].IsEmptyNode() {
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

	s, err := ParseStatement(text)

	fmt.Println(s.String())

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	nodeArray, componentIdx := s.GenerateLeafArrays()

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

	if !element[0].Left.Entry.(tree.Statement).Attributes.HasPrimitiveEntry() {
		t.Fatal("Did not detect attribute as primitive entry (Entry:", element[0].Entry, ")")
	}

	if element[0].Left.Entry.(tree.Statement).Attributes.Entry == nil {
		t.Fatal("Did detect attribute entry as nil")
	}

	if element[0].Left.Entry.(tree.Statement).Attributes.Entry.(string) != "Programme Manager" {
		t.Fatal("Incorrectly detected attribute in nested activation condition")
	}

	if element[0].Left.Entry.(tree.Statement).Aim.Entry.(string) != "suspects" {
		t.Fatal("Incorrectly detected attribute in nested activation condition")
	}

	// Test for nested elements
	nestedStmt := element[0].Left.Entry.(tree.Statement)
	leaves, _ := nestedStmt.GenerateLeafArrays()
	if len(leaves) != 4 {
		t.Fatal("Did not leaf elements of first-order nested statement correctly")
	}

	if !element[0].Left.Entry.(tree.Statement).ActivationConditionSimple.IsNil() {
		t.Fatal("Simple activation condition field of nested statement should be nil")
	}

	// Test for second-order nested statements
	if element[0].Left.Entry.(tree.Statement).ActivationConditionComplex.IsNil() {
		t.Fatal("Complex activation condition field of nested statement should not be nil")
	}

	if element[0].Left.Entry.(tree.Statement).ActivationConditionComplex.Entry.(tree.Statement).Attributes.Entry != "NOP Manager" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	if element[0].Left.Entry.(tree.Statement).ActivationConditionComplex.Entry.(tree.Statement).Aim.Entry != "orders" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	if element[0].Left.Entry.(tree.Statement).ActivationConditionComplex.Entry.(tree.Statement).DirectObject.Entry != "review" {
		t.Fatal("Did not correctly detect second-order nested activation condition element")
	}

	nestedStmt = element[0].Left.Entry.(tree.Statement).ActivationConditionComplex.Entry.(tree.Statement)
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
	flatCombo := tree.Flatten(combo.GetLeafNodes())
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

	if len(tree.Flatten(combo.GetLeafNodes())) != 3 {
		t.Fatal("Wrong number of parsed string nodes.")
	}

	// Parse all entries in tree from string to statement
	err := combo.ParseAllEntries(func(oldValue string) (tree.Statement, tree.ParsingError) {
		stmt, errStmt := ParseStatement(oldValue[strings.Index(oldValue, LEFT_BRACE)+1 : strings.LastIndex(oldValue, RIGHT_BRACE)])
		if errStmt.ErrorCode != tree.PARSING_NO_ERROR {
			return stmt, errStmt
		}
		return stmt, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
	})
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Conversion of string entries to parsed statements failed:", err.Error())
	}

	if len(tree.Flatten(combo.GetLeafNodes())) != 3 {
		t.Fatal("Wrong number of parsed string nodes.")
	}

	fmt.Println(combo.String())

	if combo.Left.Entry.(tree.Statement).ConstitutedEntity.Entry != "Program Manager" ||
		combo.Left.Entry.(tree.Statement).ConstitutingProperties.Entry != "qualified" ||
		combo.Left.Entry.(tree.Statement).ConstitutiveFunction.Entry != "is" ||
		combo.LogicalOperator != "AND" ||
		combo.Right.LogicalOperator != "XOR" ||
		combo.Right.Left.Entry.(tree.Statement).ConstitutedEntity.Entry != "Program Participant2" ||
		combo.Right.Left.Entry.(tree.Statement).ConstitutiveFunction.Entry != "is2" ||
		combo.Right.Left.Entry.(tree.Statement).ConstitutingProperties.Entry != "employed2" ||
		combo.Right.Right.Entry.(tree.Statement).ConstitutedEntity.Entry != "Program Participant" ||
		combo.Right.Right.Entry.(tree.Statement).ConstitutiveFunction.Entry != "is" ||
		combo.Right.Right.Entry.(tree.Statement).ConstitutingProperties.Entry != "employed" {

			t.Fatal("Parsing into statements failed.")
	}


}

/*
Tests parsing of special characters in regular and nested components
*/
func TestSpecialCharacters(t *testing.T) {

	input := "A(A&dsisgj=) I(=#) Bdir((l$ef% [AND] Ri@g¤#)) Bind((`?a€v [XOR] (dg/sg) !sdg£jd*s)) Cac{A(/sd-g$s%d) D(s%k£g=js) I(s§d€k+l/g#j!ds)}"

	s, err := ParseStatement(input)

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

	if s.DirectObject.Left.Entry != "l$ef%" {
		t.Fatal("Failed to detect Direct object left")
	}

	if s.DirectObject.Right.Entry != "Ri@g¤#" {
		t.Fatal("Failed to detect Direct object right")
	}

	if s.IndirectObject.Left.Entry != "`?a€v" {
		t.Fatal("Failed to detect Indirect object left")
	}

	if s.IndirectObject.Right.Entry != "(dg/sg) !sdg£jd*s" {
		t.Fatal("Failed to detect Indirect object right")
	}

	if s.ActivationConditionComplex.Entry.(tree.Statement).Attributes.Entry != "/sd-g$s%d" {
		t.Fatal("Failed to detect nested Attribute")
	}

	if s.ActivationConditionComplex.Entry.(tree.Statement).Deontic.Entry != "s%k£g=js" {
		t.Fatal("Failed to detect nested Deontic")
	}

	if s.ActivationConditionComplex.Entry.(tree.Statement).Aim.Entry != "s§d€k+l/g#j!ds" {
		t.Fatal("Failed to detect nested Aim")
	}

}

/*
Tests whether parser does not mistakenly pick up component properties (e.g., A,p) as components (e.g., A).
 */
func TestUnambiguousExtractionOfComponentAndRelatedProperties(t *testing.T) {

	input := "A,p(property) A,p1(another prop) A(value)"

	s, err := ParseStatement(input)

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

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", text, "(", ")")
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

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", text, "(", ")")
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

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", text, "(", ")")
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

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", text, "(", ")")
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

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", text, "(", ")")
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

	suffix, annotation, content, err := extractSuffixAndAnnotations("A", text, "(", ")")
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

	stmt, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	// Check Attributes

	if stmt.Attributes.GetLeafNodes()[0][0].Suffix != "1" {
		t.Fatal("Suffix should be 1, but is:", stmt.Attributes.GetLeafNodes()[0][0].Suffix)
	}

	if stmt.Attributes.GetLeafNodes()[1][0].Suffix != "2" {
		t.Fatal("Suffix should be 2, but is:", stmt.Attributes.GetLeafNodes()[1][0].Suffix)
	}

	if stmt.Attributes.GetLeafNodes()[2][0].Suffix != "3" {
		t.Fatal("Suffix should be 3, but is:", stmt.Attributes.GetLeafNodes()[2][0].Suffix)
	}

	if stmt.Attributes.GetLeafNodes()[0][0].Annotations != "[annotation1]" {
		t.Fatal("Suffix should be [annotation1], but is:", stmt.Attributes.GetLeafNodes()[0][0].Annotations)
	}

	if stmt.Attributes.GetLeafNodes()[1][0].Annotations != "[annotation2]" {
		t.Fatal("Suffix should be [annotation2], but is:", stmt.Attributes.GetLeafNodes()[1][0].Annotations)
	}

	if stmt.Attributes.GetLeafNodes()[2][0].Annotations != nil {
		t.Fatal("Suffix should be nil, but is:", stmt.Attributes.GetLeafNodes()[2][0].Annotations)
	}

	// Check Aim

	if stmt.Aim.GetLeafNodes()[0][0].Suffix != "4" {
		t.Fatal("Suffix should be 4, but is:", stmt.Aim.GetLeafNodes()[0][0].Suffix)
	}

	if stmt.Aim.GetLeafNodes()[0][0].Annotations != "[annotation=(left,right)]" {
		t.Fatal("Suffix should be [annotation=(left,right)], but is:", stmt.Aim.GetLeafNodes()[0][0].Annotations)
	}

}

/*
Tests whether complete statements are parsed and suffices and annotations stored accordingly in the underlying node structure.
This test specifically looks at nested statements
*/
func TestNodeParsingOfSuffixAndAnnotationsNestedStatement(t *testing.T) {

	text := "Cac1[leftAnno]{A1[annotation=(left,right)](content) A2[annot](content2) I[regfunc=initiate](action)} Cac2[rightAnno]{A5[|exampleAnnotation](actor)}"

	stmt, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Extraction should not have failed.")
	}

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
	if stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Attributes.Left.Suffix.(string) != "1" {
		t.Fatal("Suffix should be 1, but is:", stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Attributes.Left.Suffix.(string))
	}

	if stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Attributes.Left.Annotations.(string) != "[annotation=(left,right)]" {
		t.Fatal("Annotation should be [annotation=(left,right)], but is:", stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Attributes.Left.Annotations.(string))
	}

	// Second attribute in left activation condition
	if stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Attributes.Right.Suffix.(string) != "2" {
		t.Fatal("Suffix should be 2, but is:", stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Attributes.Right.Suffix.(string))
	}

	if stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Attributes.Right.Annotations.(string) != "[annot]" {
		t.Fatal("Annotation should be [annot], but is:", stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Attributes.Right.Annotations.(string))
	}

	// Aim
	if stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Aim.Suffix != nil {
		t.Fatal("Suffix should be nil, but is:", stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Aim.Suffix)
	}

	if stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Aim.Annotations.(string) != "[regfunc=initiate]" {
		t.Fatal("Annotation should be [regfunc=initiate], but is:", stmt.ActivationConditionComplex.Left.Entry.(tree.Statement).Aim.Annotations)
	}

	// Attributes in right activation condition
	if stmt.ActivationConditionComplex.Right.Entry.(tree.Statement).Attributes.Suffix.(string) != "5" {
		t.Fatal("Suffix should be 5, but is:", stmt.ActivationConditionComplex.Right.Entry.(tree.Statement).Attributes.Suffix.(string))
	}

	if stmt.ActivationConditionComplex.Right.Entry.(tree.Statement).Attributes.Annotations.(string) != "[|exampleAnnotation]" {
		t.Fatal("Annotation should be [|exampleAnnotation], but is:", stmt.ActivationConditionComplex.Right.Entry.(tree.Statement).Attributes.Annotations.(string))
	}

}

/*
Test proper resolution of component name for primitive element, combination header and elements, and nested statement
 */
func TestComponentNameIdentification(t *testing.T) {

	text := "A(Single Element) D( must) I((combLeft [AND] combRight)) Cac{A(Nested Element) I(perform) Bdir(something)}"

	// Indicates whether implicitly linked components (e.g., I(one) I(two)) are aggregated into a single component
	tree.AGGREGATE_IMPLICIT_LINKAGES = false

	stmt, err := ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Parsing error for statement", text)
	}

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

	if stmt.ActivationConditionComplex.Entry.(tree.Statement).Attributes.GetComponentName() != tree.ATTRIBUTES {
		t.Fatal("Incorrect identification of component name for nested statement's attribute")
	}

	if stmt.ActivationConditionComplex.Entry.(tree.Statement).Aim.GetComponentName() != tree.AIM {
		t.Fatal("Incorrect identification of component name for nested statement's aim")
	}

	if stmt.ActivationConditionComplex.Entry.(tree.Statement).DirectObject.GetComponentName() != tree.DIRECT_OBJECT {
		t.Fatal("Incorrect identification of component name for nested statement's direct object")
	}

}