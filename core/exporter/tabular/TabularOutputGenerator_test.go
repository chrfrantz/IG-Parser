package tabular

import (
	"IG-Parser/core/parser"
	"IG-Parser/core/tree"
	"fmt"
	"os"
	"testing"
)

/*
Tests the header generation function for tabular output, here as variant that does not assume component aggregation.
Note that this test is implicitly contained in IGStatementParser_test.go
*/
func TestHeaderRowGenerationWithoutComponentAggregation(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// This is tested in IGStatementParser_test.go as well
	nodeArray, componentIdx := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	if nodeArray == nil || componentIdx == nil {
		t.Fatal("Generated array or component header array should not be empty.")
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

	output := "Init, "
	output, outArr, _, err := generateHeaderRow(output, componentIdx, ";")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during header generation:", err)
	}

	// Unnecessary check is intentional to allow for flexible row specification change
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

/*
Tests the header generation function for tabular output, here as variant that applies component aggregation.
Note that this test is implicitly contained in IGStatementParser_test.go
*/
func TestHeaderRowGenerationWithComponentAggregation(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// OVERRIDE dynamic output setting
	tree.AGGREGATE_IMPLICIT_LINKAGES = true

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	// This is tested in IGStatementParser_test.go as well
	nodeArray, componentIdx := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	if nodeArray == nil || componentIdx == nil {
		t.Fatal("Generated array or component header array should not be empty.")
	}

	fmt.Println(componentIdx)

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

	output := "Init, "
	output, outArr, _, err := generateHeaderRow(output, componentIdx, ";")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during header generation:", err)
	}

	// Unnecessary check is intentional to allow for flexible row specification change
	if output == "" || len(outArr) == 0 {
		t.Fatal("Generate header row did not return filled data structures")
	}

	if output != "Init, Attributes;Deontic;Aim;Direct Object;Execution Constraint" {
		t.Fatal("Generated header row is wrong. Output: " + output)
	}

	if fmt.Sprint(outArr) != "[A D I Bdir Cex]" {
		t.Fatal("Generated component array is wrong: " + fmt.Sprint(outArr))
	}

}

/*
Tests simple tabular output without any combinations or nesting.
*/
func TestSimpleTabularOutput(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect)" +
		"Bdir(certified production facilities) "

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)
	// OVERRIDE dynamic output setting
	tree.AGGREGATE_IMPLICIT_LINKAGES = true

	// Override cell separator symbol
	CellSeparator = ";"

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 5 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputSimpleTabularOutput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Take separator for Google Sheets output
	separator := ";"

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests simple tabular output without any combinations or nesting, without header row in output.
*/
func TestSimpleTabularOutputNoHeaders(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect)" +
		"Bdir(certified production facilities) "

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(false)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)
	// OVERRIDE dynamic output setting
	tree.AGGREGATE_IMPLICIT_LINKAGES = true

	// Override cell separator symbol
	CellSeparator = ";"

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 5 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputSimpleTabularOutputNoHeaders.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Take separator for Google Sheets output
	separator := ";"

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting, but component-level combinations; no implicit combinations
*/
func TestBasicTabularOutputCombinations(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) "

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 5 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputSimpleNoNestingWithCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting, but component-level combinations; no implicit combinations,
no headers
*/
func TestBasicTabularOutputCombinationsNoHeaders(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) "

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(false)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 5 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputSimpleNoNestingWithCombinationsNoHeaders.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting, but implicitly sAND-linked components
*/
func TestBasicTabularOutputImplicitAnd(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(review) I(sustain) " +
		"Bdir(certified production facilities) "

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 6 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputSimpleNoNestingImplicitAnd.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting, but component-level combinations and implicitly linked components
*/
func TestTabularOutputCombinationsImplicitAnd(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputSimpleNoNestingWithCombinationsImplicitAnd.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting - only component-level combinations,
but including shared left elements in output.
*/
func TestTabularOutputWithSharedLeftElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputSimpleNoNestingWithSharedLeftElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting - only component-level combinations,
but including shared right elements in output.
*/
func TestTabularOutputWithSharedRightElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I((review [AND] (refresh [AND] drink)) rightShared) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputSimpleNoNestingWithSharedRightElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic tabular output without statement-level nesting - only component-level combinations,
but including shared left and right elements in output.
*/
func TestTabularOutputWithSharedLeftAndRightElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(leftShared (review [AND] (refresh [AND] revise)) rightShared) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part))."

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value: ", len(links))
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputSimpleNoNestingWithSharedLeftAndRightElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests multi-level nesting on statements, i.e., activation with own activation condition
*/
func TestTabularOutputWithTwoLevelNestedComponent(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// This is the tricky lines, specifically the second Cac{}
		"Cac{E(Program Manager) F(is) P((approved [AND] committed)) Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputTwoLevelComplexNesting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests tabular output with combination of two-level statement-level nested component and
simple activation condition (to be linked by implicit AND).
*/
func TestTabularOutputWithCombinationOfSimpleAndTwoLevelNestedComponent(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Simple activation condition
		"Cac(Upon approval)" +
		// Complex activation condition, including two-level nesting (Cac{Cac{}})
		"Cac{E(Program Manager) F(is) P((approved [AND] committed)) Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 9 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputSimpleAndTwoLevelComplexNesting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests combination of two nested activation conditions (single level)
*/
func TestTabularOutputWithCombinationOfTwoNestedComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputTwoNestedComplexComponents.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests combination of three nested activation conditions (single level), including embedded component-level nesting
*/
func TestTabularOutputWithCombinationOfThreeNestedComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}" +
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputThreeNestedComplexComponents.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests partial AND-linked statement-level combinations, expects the inference of implicit AND to non-linked complex component
*/
func TestTabularOutputWithNestedStatementCombinationsImplicitAnd(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Proper combination
		"Cac{Cac{E(Program Manager) F(is) P(approved)} [OR] " +
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Implicitly linked nested statement
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputNestedComplexCombinationsImplicitAnd.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests flattened representation (IG Core) of partial AND-linked statement-level combinations, expects the inference of implicit AND to non-linked complex component
*/
func TestTabularOutputWithNestedStatementCombinationsImplicitAndIGCore(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Proper combination
		"Cac{Cac{E(Program Manager) F(is) P(approved)} [OR] " +
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Implicitly linked nested statement
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Core output
	SetProduceIGExtendedOutput(false)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputNestedComplexCombinationsImplicitAndIGCore.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests partial XOR-linked statement-level combinations, expects the inference of implicit AND to non-linked complex component
*/
func TestTabularOutputWithNestedStatementCombinationsXOR(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Complex expression with XOR linkage
		"Cac{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// should be automatically linked using AND
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputNestedComplexCombinationsXor.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement-level combinations, alongside embedded component-level combinations to ensure the
filtering of within-statement component-level combinations are filtered prior to statement assembly.
*/
func TestTabularOutputWithNestedStatementCombinationsAndComponentCombinations(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputNestedStatementCombinationsAndComponentCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement-level combinations, alongside embedded component-level combinations to ensure the
filtering of within-statement component-level combinations are filtered prior to statement assembly.
Includes generation of output with shared elements.
*/
func TestTabularOutputWithNestedStatementCombinationsAndComponentCombinationsWithSharedElements(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	fmt.Println(text)

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputNestedStatementCombinationsAndComponentCombinationsWithSharedElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement-level combinations on multiple components, alongside selected embedded component-level nesting,
alongside combination with non-nested components and implicit linkage of combinations and simple nesting.

Includes testing of reclassification of mistakenly identified combinations and reparsing as nested statements.

Dynamic output variant
*/

func TestTabularOutputWithMultipleNestedStatementsAndSimpleComponentsAcrossDifferentComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		// This is another nested component - should have implicit link to regular Bdir
		"Bdir{A(farmers) that I((apply [OR] plan to apply)) for Bdir(organic farming status)}" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Dynamic output
	SetDynamicOutput(true)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != false {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Generated Component References: ", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 9 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputMultipleNestedStatementsAndSimpleComponentsAcrossDifferentComponents.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement-level combinations on multiple components, alongside selected embedded component-level nesting,
alongside combination with non-nested components and implicit linkage of combinations and simple nesting.

Includes testing of reclassification of mistakenly identified combinations and reparsing as nested statements.

Static output variant
*/

func TestStaticTabularOutputWithMultipleNestedStatementsAndSimpleComponentsAcrossDifferentComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		// This is another nested component - should have implicit link to regular Bdir
		"Bdir{A(farmers) that I((apply [OR] plan to apply)) for Bdir(organic farming status)}" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This is the tricky line, the multiple aims
		"Cac{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Generated Component References: ", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaMultipleNestedStatementsAndSimpleComponentsAcrossDifferentComponents.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement with multiple levels of component-level nesting, alongside component combination.
Note: Generally initially identified as combination, but recategorized as regular nested statement.
Exemplifies exception handling of parser.
*/
func TestStaticTabularOutputWithMultiLevelNestingAndComponentLevelCombinations(t *testing.T) {

	text := "Cac{When A(Program Manager) I(reveals)	Bdir{A,p(accredited) A(certifying agent) I([is not in compliance]) " +
		// Parser miscategorizes as combination based on nested logical operators, but should recover
		"with the Bdir((Act [OR] regulations in this part))} " +
		"Cac{When A(Program Manager) I((([inspects] [OR] [reviews]) [OR] [investigates])) " +
		"Bind,p(accredited) Bind(certifying agent)}} " +
		"A([Program Manager]) D(shall) I([send]) Bdir(notification) Bdir,p(of non-compliance) to the " +
		"Bind,p(accredited) Bind(certifying agent)."

	// Dynamic output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := "|"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Generated Component References: ", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaMultiLevelNestingAndComponentLevelCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests or else statements with statement-level annotations.
*/
func TestStaticTabularOutputOrElseAnnotations(t *testing.T) {

	text := "A(farmer) D([must]) I(submit) Bdir,p(an organic systems) Bdir(plan) Cex(by the end of the " +
		"calendar year) O[consequence]{the A(certifier) D(may) I(suspend) the Bdir,p(farmer’s) Bdir(operating license)}"

	// Dynamic output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Activate annotations
	SetIncludeAnnotations(true)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// No shared elements
	SetIncludeSharedElementsInTabularOutput(false)

	// Take separator for Google Sheets output
	separator := "|"

	// Test for correct configuration for dynamic output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Generated Component References: ", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticOrElseAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests statement-level combinations with incompatible component symbols
*/
func TestTabularOutputWithStatementCombinationsOfIncompatibleComponents(t *testing.T) {

	text := "A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		"Cac{Cac{E(Program Manager) F(is) P(approved)} [XOR] " +
		// This combination mixes the previous Cac with Cex - and should fail
		"Cex{A(NOP Official) I((recognizes [AND] accepts)) Bdir(Program Manager)}} " +
		// non-linked additional activation condition (should be linked by implicit AND)
		"Cac{A(Another Official) I(complains) Bdir(Program Manager) Cex(daily)}"

	_, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_ERROR_INVALID_TYPES_IN_NESTED_STATEMENT_COMBINATION {
		t.Fatal("Parser should have picked up on invalid component combinations.")
	}

}

/*
Tests combination of two nested activation conditions (single level) for static output
*/
func TestStaticTabularOutputBasicStatement(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(approved) Bdir,p(certified) Bdir((production [operations] [AND] handling operations)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaBasicStatement.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests nested combinations with focus on proper resolution of references.
*/
func TestStaticTabularOutputNestedCombinations(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(approved) Bdir,p(certified) Bdir((production [operations] [AND] handling operations)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{Cac{E(Program Manager) F(is) P((approved [AND] committed))} [XOR] " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNestedCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests nested combinations and implicit AND-linked statement with focus on proper resolution of references.
*/
func TestStaticTabularOutputNestedCombinationsImplicitAnd(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(approved) Bdir,p(certified) Bdir((production [operations] [AND] handling operations)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{Cac{E(Program Manager) F(is) P((approved [AND] committed))} [XOR] " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Activation condition 3 (to be linked implicitly)
		"Cac{A(Another official) I(does) Bdir(something else)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNestedCombinationsImplicitAnd.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests flat printing (IG Core) of nested combinations and implicit AND-linked statement with focus on proper resolution of references.
*/
func TestStaticTabularOutputNestedCombinationsImplicitAndIGCore(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(approved) Bdir,p(certified) Bdir((production [operations] [AND] handling operations)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{Cac{E(Program Manager) F(is) P((approved [AND] committed))} [XOR] " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Activation condition 3 (to be linked implicitly)
		"Cac{A(Another official) I(does) Bdir(something else)}"

	// Static output
	SetDynamicOutput(false)
	// IG Core output (no component-level nesting)
	SetProduceIGExtendedOutput(false)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNestedCombinationsImplicitAndIGCore.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests printing of nested properties and proper resolution of references.
*/
func TestStaticTabularOutputNestedProperties(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(approved) Bdir,p(certified) Bdir,p{E(operations) that F(have experience) with P(farming)} " +
		"Bdir((production [operations] [AND] handling operations)) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{Cac{E(Program Manager) F(is) P((approved [AND] committed))} [XOR] " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}} " +
		// Activation condition 3 (to be linked implicitly)
		"Cac{A(Another official) I(does) Bdir(something else)}"

	// Static output
	SetDynamicOutput(false)
	// IG Core output (no component-level nesting)
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 9 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNestedProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Correct parsing of shared elements left and right (on Aim and Cex).
*/
func TestStaticTabularOutputBasicStatementSharedLeftAndRightElements(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink)) rightShared) " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part) and beyond) "

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 5 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaBasicStatementLeftAndRightElements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Correct parsing of private properties (mix of shared and private) for static output
*/
func TestStaticTabularOutputBasicStatementSharedAndPrivateProperties(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p1(accredited) Bdir1(certifying agents) Bdir(other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A(NOP Official) I(recognizes) Bdir(Program Manager)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaBasicStatementPrivateProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests private properties on indexed components.
*/
func TestStaticTabularOutputBasicStatementPrivatePropertiesOnly(t *testing.T) {

	text := "Bdir1,p(organic farming) Bdir1(provisions) and Bdir2,p(improper) Bdir2(rules)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 1 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaBasicStatementPrivatePropertiesOnly.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests private and shared properties combinations with combinations of indexed components.
*/
func TestStaticTabularOutputBasicStatementMixSharedPrivatePropertyComponents(t *testing.T) {

	text := "Bdir1((left [OR] right)) Bdir1,p((private [AND] public)) Bdir(general object) Bdir,p((shared [XOR] non-shared))"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 2 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaBasicStatementMixSharedPrivatePropertyComponents.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests combination of two nested activation conditions (single level) for static output and tests for a mix
of shared and private properties (on top level) and private properties only on nested level
*/
func TestStaticTabularOutputBasicStatementMixSharedPrivateAndNestedPrivatePropertiesOnly(t *testing.T) {

	text := "A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p1(accredited) Bdir1(certifying agents) Bdir(other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2 -- note that associated is wrongly annotated, leading to linkage to both Bdirs
		"Cac{A(NOP Official) I(recognizes) Bdir1,p1(responsible) Bdir1(Program Manager) and Bdir,p2(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// With shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaStatementSharedAndPrivateOnlyProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests presence of shared and private properties on indexed components (e.g., A1(content)), as opposed
to indexed properties only (e.g., A,p1(content)).
*/
func TestStaticTabularOutputBasicStatementComponentLevelIndexedProperties(t *testing.T) {

	text := "A1,p(National Organic Program's) A1(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p(accredited) Bdir1(certifying agents) Bdir(other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2 - shared properties implied in this activation condition due to wrong syntax
		"Cac{A(NOP Official) I(recognizes) Bdir,p1(responsible) Bdir1(Program Manager) and Bdir,p2(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaStatementComponentLevelIndexedProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests presence of shared and private properties on indexed components (e.g., A1(content)), as opposed
to indexed properties only (e.g., A,p1(content)).
*/
func TestStaticTabularOutputBasicStatementComponentLevelIndexedPropertiesAnnotations(t *testing.T) {

	text := "A1,p[prop=qualitative](National Organic Program's) A1[type=animate](Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p(accredited) Bdir1(certifying agents) Bdir(other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2 - shared properties implied in this activation condition
		"Cac{A(NOP Official) I(recognizes) Bdir,p1(responsible) Bdir1(Program Manager) and Bdir,p2(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 7 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaStatementComponentLevelIndexedProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests parsing of shared and private properties (via index) on component-level nested components.
*/
func TestStaticTabularOutputBasicStatementComponentLevelNestedPrivateAndSharedProperties(t *testing.T) {

	text := "Bdir{A1,p(first) A,p(shared) (A1(farmer) [OR] A2(citizen))}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 1 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComponentLevelNestedPrivateAndSharedProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests combination of two nested activation conditions (single level) , tests for a mix
of shared and private properties (on top level) and private properties only on nested level,
and includes *deactivated* annotations on various components.
*/
func TestStaticTabularOutputBasicStatementMixedPropertiesAnnotationsDeactivated(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p1[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		// Activation condition 2 - misspelt second property is intentional (to make it shared property)
		"Cac{A(NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1(Program Manager) and Bdir,p1(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaStatementSharedAndPrivateOnlyProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests combination of two nested activation conditions (single level) , tests for a mix
of shared and private properties (on top level) and private properties only on nested level,
and includes *activated* annotations on various components.
*/
func TestStaticTabularOutputBasicStatementMixedPropertiesAnnotationsActivated(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A[type=enforcer](NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1[type=main object](Program Manager) and Bdir2,p2[type=third party](associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaStatementMixedPropertiesAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests for proper replacement of different types of quotations marks for preprocessing Google Sheets output. Includes
complexity of previous tests.
*/
func TestStaticTabularOutputBasicStatementEmbeddedQuotationSymbolsGoogleSheets(t *testing.T) {

	text := "A,p(National Organic Program's) A(\"Program Manager\"), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p(accredited) Bdir1(\"certifying agents) Bdir(\"other agents\")" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A('NOP Official') I(recognizes) Bdir1,p(responsible) Bdir1(Program Manager) and Bdir2,p(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaStatementEscapedQuotationMarksGoogleSheets.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests for proper replacement of different types of quotations marks for preprocessing CSV output. Includes
complexity of previous tests.
*/
func TestStaticTabularOutputBasicStatementEmbeddedQuotationSymbolsCSV(t *testing.T) {

	text := "A,p(National Organic Program's) A(\"Program Manager\"), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I(inspect), " +
		"I(sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p(recognized) Bdir1,p(accredited) Bdir1(\"certifying agents) Bdir(\"other agents\")" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{E(Program Manager) F(is) P((approved [AND] committed))} " +
		// Activation condition 2
		"Cac{A('NOP Official') I(recognizes) Bdir1,p(responsible) Bdir1(Program Manager) and Bdir2,p(associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(false)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for CSV output
	separator := "|"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaStatementEscapedQuotationMarksCSV.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_CSV, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateCSVOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper export of annotations on nested components.
*/
func TestStaticTabularOutputNestedStatementsAnnotations(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac1[ctx=stAte]{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		// Activation condition 2
		"CacB[annotation2]{A[type=enforcer](NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1[type=main object](Program Manager) and Bdir2,p2[type=third party](associated) Bdir2(inspectors)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNestedStatementsAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper export of annotations on nested component combinations.
*/
func TestStaticTabularOutputNestedStatementCombinationAnnotations(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{Cac1[ctx=stAte]{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		"[XOR] " +
		// Activation condition 2
		"CacB[annotation2]{A[type=enforcer](NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1[type=main object](Program Manager) and Bdir2,p2[type=third party](associated) Bdir2(inspectors)}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNestedStatementCombinationsAnnotations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper export of annotations on nested component combinations and individual nested statements. Tested on Google Sheets output.
*/
func TestStaticTabularOutputNestedStatementsAndCombinationMixAnnotationsGoogleSheets(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		// Implicitly linked nested direct object with annotations and invalid structure (no aim and context)
		"Bdir{A[type=animate](another actor) A,p[prop=qualitative](who does not comply)} " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{Cac1[ctx=stAte]{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		"[XOR] " +
		// Activation condition 2
		"CacB[annotation2]{A[type=enforcer](NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1[type=main object](Program Manager) and Bdir2,p2[type=third party](associated) Bdir2(inspectors)}} " +
		// Implicitly linked activation condition with diverse annotations
		"CacC[ABdir]{A[type=animate](further entity) I[act=violate](violates) Bdir[entity=law](part of provisions)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 9 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNestedStatementsAndCombinationMixAnnotationsGoogleSheets.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper export of annotations on nested component combinations and individual nested statements. Tested on CSV output.
*/
func TestStaticTabularOutputNestedStatementsAndCombinationMixAnnotationsCSV(t *testing.T) {

	text := "A,p[type=animate](National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), " +
		"D(may) " +
		"I[act=main](inspect), " +
		"I[act=variable](sustain (review [AND] (refresh [AND] drink))) " +
		"Bdir,p[shared](recognized) Bdir1,p[private](accredited) Bdir1[type=main object](certifying agents) Bdir[type=third party](other agents)" +
		// Implicitly linked nested direct object with annotations and invalid structure (no aim and context)
		"Bdir{A[type=animate](another actor) A,p[prop=qualitative](who does not comply)} " +
		"Cex(for compliance with the (Act or [XOR] regulations in this part)) " +
		// Activation condition 1
		"Cac{Cac1[ctx=stAte]{E(Program Manager) F[cfunc=state](is) P((approved [AND] committed))} " +
		"[XOR] " +
		// Activation condition 2
		"CacB[annotation2]{A[type=enforcer](NOP Official) I[act=main](recognizes) Bdir1,p1(responsible) Bdir1[type=main object](Program Manager) and Bdir2,p2[type=third party](associated) Bdir2(inspectors)}} " +
		// Implicitly linked activation condition with diverse annotations
		"CacC[ABdir]{A[type=animate](further entity) I[act=violate](violates) Bdir[entity=law](part of provisions)}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Indicates whether header row is included in output.
	SetIncludeHeaders(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 9 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNestedStatementsAndCombinationMixAnnotationsCSV.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_CSV, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateCSVOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests proper parsing and export of within-component combination linkages (wAND).
*/
func TestStaticTabularOutputParsingOfWithinComponentLinkages(t *testing.T) {

	// Simple output with multiple combinations within Cex
	text := "Cex[exampleConstraint](for compliance with (left [AND] right) as well as (left1 [XOR] right1) shared) Cex(outlier)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := "|"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 1 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaBasicWithinComponentCombinationLinkages.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}
}

/*
Tests tabular output for the default example statement that showcases most IG features.
*/
func TestTabularOutputDefaultExample(t *testing.T) {

	// Default example
	text := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), " +
		"D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) " +
		"Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and " +
		"Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) " +
		"Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) " +
		"under the condition that Cac{Cac[state]{A[role=monitored,type=animate](Operations) " +
		"I[act=violate](were (non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] " +
		"Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)}}."

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 8 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputTabularDefaultExample.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests complex multi-level nesting in tabular output (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNesting(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNesting.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic component pair (which extrapolates into multiple statements).
*/
func TestTabularOutputBasicComponentPairs(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) {I(aim1) Bdir(object1) [XOR] I(aim2) Bdir(object2)} Cac(condition)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaBasicComponentPairs.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic component pair and within-component combination (which extrapolate into multiple statements).
*/
func TestTabularOutputBasicComponentPairAndWithinComponentCombinations(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) {I(aim1) Bdir(object1) [XOR] I(aim2) Bdir(object2)} Cac((condition1 [OR] condition2))"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaBasicComponentPairsAndWithinComponentCombinations.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests basic component pair and within-component combination and component-level nesting (extrapolating into multiple statements).
*/
func TestTabularOutputBasicComponentPairAndWithinComponentCombinationsAndNestedStatements(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) {I(aim1) Bdir(object1) Cac{A(actor2) I(aim3)} [XOR] Cac{A(actor4) I(aim4)} I(aim2) Bdir(object2)} Bind((indirectobject1 [OR] indirectobject2))"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaBasicComponentPairsAndWithinComponentCombinationsAndNestedStatements.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests standard statement with component pairs instead of component combination (which extrapolates into multiple statements).
*/
func TestTabularOutputStandardStatementComponentPairs(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), " +
		"D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) " +
		"Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and " +
		"Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) " +
		"Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) " +
		"under the condition that {Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate]((were non-compliant [OR] violated)) " +
		"Bdir[type=inanimate](organic farming provisions)} [AND] Cac[state]{A[role=enforcer,type=animate](Manager) " +
		"I[act=terminate](has concluded) Bdir[type=activity](investigation)}}."

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaStandardStatementComponentPairs.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests complex nested statements without component pairs.
*/
func TestTabularOutputComplexNestedCombinationsWithoutComponentPairs(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "Cac{Cac{Cac{A(actor1) I(aim1) Bdir(object1)} [AND] Cac{A(actor2) I(aim2) Bdir(object2)} [AND] Cac{A(actor4) I(aim4)}} [OR] Cac{Cac{A(actor3) I(aim3) Bdir(object3)} [XOR] Cac{Cac{A(actor6) I(aim6) Bdir(object6)} [AND] Cac{Cac{A(actor7) I(aim7) Bdir(object7)} [XOR] Cac{A(actor8) I(aim8)}}}}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexNestedCombinationsWithoutComponentPairs.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests complex nested statements *with* component pairs.
*/
func TestTabularOutputComplexNestedCombinationsWithComponentPairs(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "{Cac{Cac{A(actor1) I(aim1) Bdir(object1)} [AND] Cac{A(actor2) I(aim2) Bdir(object2)} [AND] Cac{A(actor4) I(aim4)}} [OR] Cac{Cac{A(actor3) I(aim3) Bdir(object3)} [XOR] Cac{Cac{A(actor6) I(aim6) Bdir(object6)} [AND] Cac{Cac{A(actor7) I(aim7) Bdir(object7)} [XOR] Cac{A(actor8) I(aim8)}}}}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexNestedCombinationsWithComponentPairs.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the combined use of component combinations, component pairs, as well as nested component statements.
*/
func TestTabularOutputNestedCombinationsComponentLevelNestingAndComponentPairs(t *testing.T) {

	// Statement with multi-level nesting, combinations as well as component pairs
	text := "A(actor1) I(aim1) Cac{Cac{A(actor2) I(aim2)} [XOR] Cac{A(actor3) I(aim3)}} {Bdir(directobject1) Bind(indirectobject1) [OR] Bdir{ A(actor4) I(aim4) Bdir(directobject2) Cac{A(actor5) I(aim5)}} Bind(indirectobject2)} "

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNestedCombinationsComponentLevelNestingAndComponentPairs.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests correct identification of logical operators between neighbouring nodes.
*/
func TestTabularOutputLogicalOperatorsNeighbouringStatementsComponentPairs(t *testing.T) {

	text := "{ A(actor1) I(aim1) [XOR] {A(actor2) I(aim2) [AND] A(actor3) I(aim3)}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaLogicalOperatorsNeighbouringStatementsComponentPairs.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests correct parsing of component pairs in nested components.
*/
func TestTabularOutputComponentPairsInNestedComponents(t *testing.T) {

	text := "A(Individuals) D(must) { I(monitor) Bdir(compliance) [AND] I(report) Bdir(violation) } Cac(in the case of (repeated offense [OR] other reasons)) O{ A(actor2) D(must) {I(enforce) Bdir(compliance) [OR] I(delegate) Bdir(enforcement)}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComponentPairsInNestedComponents.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests correct parsing of component pairs in nested components in IG Core.
*/
func TestTabularOutputComponentPairsInNestedComponentsIGCore(t *testing.T) {

	text := "A(Individuals) D(must) { I(monitor) Bdir(compliance) [AND] I(report) Bdir(violation) } Cac(in the case of (repeated offense [OR] other reasons)) O{ A(actor2) D(must) {I(enforce) Bdir(compliance) [OR] I(delegate) Bdir(enforcement)}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(false)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComponentPairsInNestedComponentsIGCore.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests the expansion of missing combination parentheses in input statements (for coder convenience).
Example: Bdir(left [AND] right) should be expanded to Bdir((left [AND] right)) and then correctly encoded.
*/
func TestTabularOutputExpandComponentCombinationsWithMissingParentheses(t *testing.T) {

	text := " A(actor) D(may) {I(leftAim) Bdir(leftObject) [OR] I(rightAim) Bdir(rightObject)} Cac{ {A(actor2) I(aim2 [XOR] aim4) [XOR] A(actor3) I(aim3)} }"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaExpandComponentCombinationsWithMissingParentheses.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests visual output for nested properties, including component pairs, simple nested statements and primitive properties.
*/
func TestTabularOutputNestedPropertiesIncludingComponentPairsAndNestedPropertyAndPrimitiveProperty(t *testing.T) {

	// Statement with multi-level nesting, combinations as well as component pairs
	text := "Such E(notification) M(shall) F(provide): (1) A P(description of each noncompliance); " +
		"(2) The P(facts upon which the notification of noncompliance is based); and " +
		"(3) The P1(date) " +
		// Component-pair property
		"P1,p{by which the A(certified operation) D(must) " +
		"{I(rebut [XOR] correct) Bdir,p(each) Bdir(noncompliance) [AND] I(submit) Bdir,p(supporting) " +
		"Bdir(documentation) of Bdir,p(each such correction) Cac(when correction is possible)}} " +
		// Primitive property
		"P1,p(private component) " +
		// nested property
		"P1,p{where E(date) F(is defined) in the P(Gregorian calendar)}."

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Define separator for output generation
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaPrivateNestedProperties.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests generation of tabular output with linebreaks in IG Script input.
*/
func TestTabularOutputWithLinebreaksInIGScriptInput(t *testing.T) {

	// Statement with line break
	text := "A(actor) I(aim) Bdir(object1) Bind(object2)\n Cac(condition1)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Define separator for output generation
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaLinebreakInIgScriptInput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests generation of tabular output with non-cell separator in IG Script input.
*/
func TestTabularOutputWithNonCellSeparatorInIGScriptInput(t *testing.T) {

	// Statement with cell separator
	text := "A(actor) I(aim) Bdir(object1) Bind(object2)| Cac(condition1)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Define separator for output generation
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNonCellSeparatorInIgScriptInput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests generation of tabular output with cell separator in IG Script input.
*/
func TestTabularOutputWithCellSeparatorInIGScriptInput(t *testing.T) {

	// Statement with cell separator
	text := "A(actor) I(aim) Bdir(object1) Bind(object2)| Cac(condition1)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Define separator for output generation
	separator := "|"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", "", text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaCellSeparatorInIgScriptInput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests generation of tabular output with linebreaks in Original Statement input.
*/
func TestTabularOutputWithLinebreaksInOriginalStatementInput(t *testing.T) {

	// Original statement (with linebreak)
	originalStatement := "actor aim object1\n object2 condition1"
	// Coded statement (no linebreak)
	text := "A(actor) I(aim) Bdir(object1) Bind(object2) Cac(condition1)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Define separator for output generation
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", originalStatement, text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaLinebreakInOriginalStatementInput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests generation of tabular output with non-cell separator in Original Statement input.
*/
func TestTabularOutputWithNonCellSeparatorInOriginalStatementInput(t *testing.T) {

	// Original statement (with non-cell separator)
	originalStatement := "actor aim object1| object2 condition1"
	// Coded statement (no cell separator)
	text := "A(actor) I(aim) Bdir(object1) Bind(object2) Cac(condition1)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Define separator for output generation
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", originalStatement, text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNonCellSeparatorInOriginalStatementInput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests generation of tabular output with cell separator in Original Statement input.
*/
func TestTabularOutputWithCellSeparatorInOriginalStatementInput(t *testing.T) {

	// Original statement (with cell separator)
	originalStatement := "actor aim object1| object2 |condition1"
	// Coded statement (no cell separator)
	text := "A(actor) I(aim) Bdir(object1) Bind(object2) Cac(condition1)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Define separator for output generation
	separator := "|"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", originalStatement, text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaCellSeparatorInOriginalStatementInput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests generation of tabular output with linebreaks in Original Statement and IG Script input.
*/
func TestTabularOutputWithLinebreaksInOriginalStatementAndIgScriptInput(t *testing.T) {

	// Original statement (with linebreak)
	originalStatement := "actor \naim \nobject1\n object2 condition1"
	// Coded statement (no linebreak)
	text := "A(actor) I(aim)\n Bdir(object1) Bind(object2)\n Cac(condition1)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Define separator for output generation
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", originalStatement, text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaLinebreakInOriginalStatementAndIgScriptInput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests generation of tabular output with non-cell separator in Original Statement and IG Script input.
*/
func TestTabularOutputWithNonCellSeparatorInOriginalStatementAndIgScriptInput(t *testing.T) {

	// Original statement
	originalStatement := "actor |aim |object1| object2 condition1"
	// Coded statement
	text := "A(actor) I(aim)| Bdir(object1) Bind(object2)| Cac(condition1)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Define separator for output generation
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", originalStatement, text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaNonCellSeparatorInOriginalStatementAndIgScriptInput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests generation of tabular output with cell separator in Original Statement and IG Script input.
*/
func TestTabularOutputWithCellSeparatorInOriginalStatementAndIgScriptInput(t *testing.T) {

	// Original statement
	originalStatement := "actor |aim |object1| object2 condition1"
	// Coded statement
	text := "A(actor) I(aim)| Bdir(object1) Bind(object2)| Cac(condition1)"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Define separator for output generation
	separator := "|"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) != 1 {
		t.Fatal("Wrong statement count: ", stmts)
	}

	results := GenerateTabularOutputFromParsedStatements(stmts, "", originalStatement, text, "123", "", true, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error when generating output:", err)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaCellSeparatorInOriginalStatementAndIgScriptInput.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	// Aggregate output if multiple results
	output := ""
	for _, v := range results {
		output += v.Output
	}

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests IG Script output (with single entry) as part of Google Sheets tabular output for complex multi-level nesting (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNestingWithIgScriptOutputSingleEntry(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNestingIgScriptOutputSingleRowEntry.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_FIRST_ENTRY)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests IG Script output (for all entries) as part of Google Sheets tabular output for complex multi-level nesting (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNestingWithIgScriptOutputAllEntries(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNestingIgScriptOutputAllRowEntries.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests Original Statement output (with single entry) as part of Google Sheets tabular output for complex multi-level nesting (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNestingWithOriginalStatementOutputSingleEntry(t *testing.T) {

	// Original statement to be included in output
	originalStatement := "actor1 aim1 that actor2 aim2 under the condition that actor3 aim3 something or actor4 aim4 something else"

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNestingOriginalStatementOutputSingleRowEntry.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateGoogleSheetsOutput(statementMap, originalStatement, text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_FIRST_ENTRY, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests Original Statement output (for all entries) as part of Google Sheets tabular output for complex multi-level nesting (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNestingWithOriginalStatementOutputAllEntries(t *testing.T) {

	// Original statement to be included in output
	originalStatement := "actor1 aim1 that actor2 aim2 under the condition that actor3 aim3 something or actor4 aim4 something else"

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNestingOriginalStatementOutputAllRowEntries.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateGoogleSheetsOutput(statementMap, originalStatement, text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_NONE)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests Original Statement and IG Script output (both with single entry) as part of Google Sheets tabular output for complex multi-level nesting (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNestingWithOriginalStatementAndIgScriptOutputSingleEntry(t *testing.T) {

	// Original statement to be included in output
	originalStatement := "actor1 aim1 that actor2 aim2 under the condition that actor3 aim3 something or actor4 aim4 something else"

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNestingOriginalStatementOutputSingleRowAndIgScriptSingleRowEntry.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateGoogleSheetsOutput(statementMap, originalStatement, text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_FIRST_ENTRY, IG_SCRIPT_OUTPUT_FIRST_ENTRY)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests Original Statement and IG Script output (both with single entry) as part of CSV tabular output for complex multi-level nesting (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNestingWithOriginalStatementAndIgScriptOutputSingleEntryCSV(t *testing.T) {

	// Original statement to be included in output
	originalStatement := "actor1 aim1 that actor2 aim2 under the condition that actor3 aim3 something or actor4 aim4 something else"

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNestingOriginalStatementOutputSingleRowAndIgScriptSingleRowEntryCSV.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateCSVOutput(statementMap, originalStatement, text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_FIRST_ENTRY, IG_SCRIPT_OUTPUT_FIRST_ENTRY)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests Original Statement and IG Script output (both for all entries) as part of Google Sheets tabular output for complex multi-level nesting (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNestingWithOriginalStatementAndIgScriptOutputAllEntries(t *testing.T) {

	// Original statement to be included in output
	originalStatement := "actor1 aim1 that actor2 aim2 under the condition that actor3 aim3 something or actor4 aim4 something else"

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNestingOriginalStatementAndIgScriptOutputAllRowEntries.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateGoogleSheetsOutput(statementMap, originalStatement, text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests Original Statement and IG Script output (both for all entries) as part of CSV tabular output for complex multi-level nesting (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNestingWithOriginalStatementAndIgScriptOutputAllEntriesCSV(t *testing.T) {

	// Original statement to be included in output
	originalStatement := "actor1 aim1 that actor2 aim2 under the condition that actor3 aim3 something or actor4 aim4 something else"

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNestingOriginalStatementAndIgScriptOutputAllRowEntriesCSV.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateCSVOutput(statementMap, originalStatement, text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests Original Statement (in first entry) and IG Script output (for all entries) as part of Google Sheets tabular output for complex multi-level nesting (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNestingWithOriginalStatementSingleAndIgScriptOutputAllEntry(t *testing.T) {

	// Original statement to be included in output
	originalStatement := "actor1 aim1 that actor2 aim2 under the condition that actor3 aim3 something or actor4 aim4 something else"

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNestingOriginalStatementSingleRowAndIgScriptOutputAllRowEntries.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateGoogleSheetsOutput(statementMap, originalStatement, text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_FIRST_ENTRY, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests Original Statement (for all entries) and IG Script output (for first entry) as part of Google Sheets tabular output for complex multi-level nesting (e.g., Cac{{Cac{} [OR] Cac{}}}.
*/
func TestTabularOutputComplexMultilevelNestingWithOriginalStatementAllAndIgScriptOutputSingleEntry(t *testing.T) {

	// Original statement to be included in output
	originalStatement := "actor1 aim1 that actor2 aim2 under the condition that actor3 aim3 something or actor4 aim4 something else"

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "A(actor1) I(aim1) Bdir{A(actor2) I(aim2) Cac{   Cac{A(actor3) I(aim3) Bdir(something)  }   [OR]   Cac{  A(actor4) I(aim4) Bdir(something else)  }}}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 3 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestOutputStaticSchemaComplexMultilevelNestingOriginalStatementAllAndIgScriptOutputSingleRowEntry.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateGoogleSheetsOutput(statementMap, originalStatement, text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_FIRST_ENTRY)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests complex combination of component-level nesting with embedded combinations with respect to correct logical linkages
*/
func TestTabularOutputCombinationOfComponentLevelNestingAndCombinationsLogicalOperatorInference(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ Cac{ A(actor2) I(aim2) Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } Bdir(object2) } [OR] Cac{ A(actor5) I(aim5) Cac{ A(actor6) I(aim6) } } }}"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 1 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestTabularOutputCombinationOfComponentLevelNestingAndCombinationsLogicalOperatorInference.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}

/*
Tests complex combination of component-level nesting with embedded combinations with respect to correct logical linkages,
including the inference of an implicit AND.
*/
func TestTabularOutputCombinationOfComponentLevelNestingAndCombinationsLogicalOperatorInferenceImplicitAnd(t *testing.T) {

	// Statement with multi-level nesting with embedded nested statement combinations (erratic spacing is intentional)
	text := "Cac{ Cac{ A(actor6) I(actor6) } [XOR] Cac{ A(actor7) I(actor7) }} Cac{ Cac{ A(actor1) I(aim1) } [OR] Cac{ A(actor0) I(aim0) Cac{ A(actor2) I(aim2) Cac{ Cac{ A(actor3) I(aim3) } [OR] Cac{ A(actor4) I(aim4) } } Bdir(object2) } [OR] Cac{ A(actor5) I(aim5)} }}\n"

	// Static output
	SetDynamicOutput(false)
	// IG Extended output
	SetProduceIGExtendedOutput(true)
	// Indicates whether annotations are included in output.
	SetIncludeAnnotations(true)
	// Deactivate DoV
	SetIncludeDegreeOfVariability(false)
	// Include shared elements
	SetIncludeSharedElementsInTabularOutput(true)

	// Take separator for Google Sheets output
	separator := ";"

	// Test for correct configuration for static output
	if tree.AGGREGATE_IMPLICIT_LINKAGES != true {
		t.Fatal("SetDynamicOutput() did not properly configure implicit link aggregation")
	}

	stmts, err := parser.ParseStatement(text)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during parsing of statement", err.Error())
	}

	if len(stmts) > 1 {
		t.Fatal("Too many statements identified: ", stmts)
	}

	s := stmts[0].Entry.(*tree.Statement)

	fmt.Println(s.String())

	// This is tested in IGStatementParser_test.go as well as in TestHeaderRowGeneration() (above)
	leafArrays, componentRefs := s.GenerateLeafArrays(tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println("Component refs:", componentRefs)

	res, err := tree.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Unexpected error during array generation.")
	}

	fmt.Println("Input arrays: ", res)

	links := tree.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	fmt.Println("Links: ", links)

	// Content of statement links is tested in ArrayCombinationGenerator_test.go
	if len(links) != 1 {
		t.Fatal("Number of statement reference links is incorrect. Value:", len(links), "Links:", links)
	}

	// Read reference file
	content, err2 := os.ReadFile("TestTabularOutputCombinationOfComponentLevelNestingAndCombinationsLogicalOperatorInferenceImplicitAnd.test")
	if err2 != nil {
		t.Fatal("Error attempting to read test text input. Error: ", err2.Error())
	}

	// Extract expected output
	expectedOutput := string(content)

	statementMap, statementHeaders, statementHeadersNames, err := generateStatementMatrix(res, nil, "", componentRefs, links, "650", separator, OUTPUT_TYPE_GOOGLE_SHEETS, IncludeHeader())
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Generating tabular output should not fail. Error: " + fmt.Sprint(err.Error()))
	}

	// Here is the relevant parameterization of the output
	output, err := generateGoogleSheetsOutput(statementMap, "", text, statementHeaders, statementHeadersNames, separator, "", true, IncludeHeader(), ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES, IG_SCRIPT_OUTPUT_ALL_ENTRIES)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		t.Fatal("Error during Google Sheets generation. Error: " + fmt.Sprint(err.Error()))
	}

	fmt.Println("Output:", output)

	// Compare to actual output
	if output != expectedOutput {
		fmt.Println("Statement headers:\n", statementHeaders)
		fmt.Println("Statement map:\n", statementMap)
		fmt.Println("Produced output:\n", output)
		fmt.Println("Expected output:\n", expectedOutput)
		err3 := WriteToFile("errorOutput.error", output, true)
		if err3 != nil {
			t.Fatal("Error attempting to read test text input. Error: ", err3.Error())
		}
		t.Fatal("Output generation is wrong for given input statement. Wrote output to 'errorOutput.error'")
	}

}
