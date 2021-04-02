package app

import (
	"IG-Parser/exporter"
	"IG-Parser/parser"
	"IG-Parser/tree"
)

/*
Consumes statements as input and produces outfile
 */
func ConvertIGScriptToGoogleSheets(statement string, filename string) (string, tree.ParsingError) {

	// Parse IGScript statement into tree
	s,err := parser.ParseStatement(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	// Retrieve leaf arrays from generated tree (alongside frequency indications for components)
	leafArrays, componentRefs := s.GenerateLeafArrays()

	// Generate all permutations of logically-linked components to produce statements
	res := exporter.GenerateNodeArrayPermutations(leafArrays...)

	// Extract logical operator links
	links := exporter.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Export in Google Sheets format
	output, err := exporter.GenerateGoogleSheetsOutput(res, componentRefs, links, "650")
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	// Write to file
	errWrite := exporter.WriteToFile(filename, output)
	if errWrite != nil {
		// Wrap into own error, alongside generated (but not written) output
		return output, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_WRITE, ErrorMessage: errWrite.Error()}
	}

	// Return Google Sheets output and default error
	return output, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}
