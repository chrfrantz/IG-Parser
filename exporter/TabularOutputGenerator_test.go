package exporter

import (
	"IG-Parser/parser"
	"IG-Parser/tree"
	"fmt"
	"io/ioutil"
	"testing"
)
/*
Tests the header generation function for tabular output.
Note that this test is implicitly contained in IGStatementParser_test.go
 */
func TestHeaderRowGeneration(t *testing.T) {
	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	s, err := parser.ParseStatement(text)

	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during parsing: ", err.Error())
	}

	// This is tested in IGStatementParser_test.go as well
	nodeArray, componentIdx := s.GenerateLeafArrays()

	if nodeArray == nil || componentIdx == nil {
		t.Fatal("Generated array or component header array should not be empty.")
	}

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

	output := "Init, "
	output, outArr, _ := generateHeaderRow(output, componentIdx, ";")

	if output == "" || len(outArr) == 0 {
		t.Fatal("Generate header row did not return filled data structures")
	}

	if output != "Init, Attributes;Deontic;Aim_1;Aim_2;Direct Object;Execution Constraint_1;Execution Constraint_2" {
		t.Fatal("Generated header row is wrong. Output: " + output)
	}

	if fmt.Sprint(outArr) != "[A D I_1 I_2 Bdir Cex_1 Cex_2]" {
		t.Fatal("Generated component array is wrong: " + fmt.Sprint(outArr))
	}

}

func TestTabularOutput(t *testing.T) {
	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Error("Error during parsing of statement")
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays()

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, error := ioutil.ReadFile("GeneratedGoogleSheetsOutput.test")
	if error != nil {
		t.Fatal("Error attempting to read test text input. Error: ", error.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	_, statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "sanity")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		WriteToFile("errorOutput.error", output)
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

func TestTabularOutputWithNestedComponent(t *testing.T) {
	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{E(Program Manager) F(is) P((approved [AND] committed)) Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}}"

	s,err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Error("Error during parsing of statement")
	}

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays()

	res, err := GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		//t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, error := ioutil.ReadFile("GeneratedGoogleSheetsOutput.test")
	if error != nil {
		t.Fatal("Error attempting to read test text input. Error: ", error.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	_, statementMap, statementHeaders, statementHeadersNames, err := generateTabularStatementOutput(res, componentRefs, links, "650")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := GenerateGoogleSheetsOutput(statementMap, statementHeaders, statementHeadersNames, "sanity")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		WriteToFile("errorOutput.error", output)
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

// ensure ordering of column headers

// introduce for statement combinations

// introduce feature for other components

// test with invalid statement and empty input nodes, unbalanced parentheses, missing ID
