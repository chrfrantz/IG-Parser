package app

import (
	"IG-Parser/exporter"
	"IG-Parser/parser"
	"IG-Parser/tree"
	"strconv"
)

/*
Consumes statements as input and produces outfile
Arguments include the IGScript-annotated statement, statement ID based on which substatements are generated,
and a filename for the output. If the filename is empty,
no output will be written.
Returns Google Sheets output as string, and error (defaults to tree.PARSING_NO_ERROR).
 */
func ConvertIGScriptToGoogleSheets(statement string, stmtId int, filename string) (string, tree.ParsingError) {

	// Parse IGScript statement into tree
	s,err := parser.ParseStatement(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	// Retrieve leaf arrays from generated tree (alongside frequency indications for components)
	leafArrays, componentRefs := s.GenerateLeafArrays()

	// Generate all permutations of logically-linked components to produce statements
	res, err := exporter.GenerateNodeArrayPermutations(leafArrays...)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	// Extract logical operator links
	links := exporter.GenerateLogicalOperatorLinkagePerCombination(res, true, true)

	// Export in Google Sheets format
	output, err := exporter.GenerateGoogleSheetsOutput(res, componentRefs, links, strconv.Itoa(stmtId))
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	// Outfile will only be written if filename is specified
	if filename != "" {
		// Write to file
		errWrite := exporter.WriteToFile(filename, output)
		if errWrite != nil {
			// Wrap into own error, alongside generated (but not written) output
			return output, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_WRITE, ErrorMessage: errWrite.Error()}
		}
	}

	// Return Google Sheets output and default error
	return output, tree.ParsingError{ErrorCode: tree.PARSING_NO_ERROR}
}
