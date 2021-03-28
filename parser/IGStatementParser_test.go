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
