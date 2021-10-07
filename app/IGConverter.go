package app

import (
	"IG-Parser/exporter"
	"IG-Parser/parser"
	"IG-Parser/tree"
	"fmt"
	"log"
)

/*
Consumes statements as input and produces outfile
Arguments include the IGScript-annotated statement, statement ID based on which substatements are generated,
and a filename for the output. If the filename is empty, no output will be written.
Returns Google Sheets output as string, and error (defaults to tree.PARSING_NO_ERROR).
 */
func ConvertIGScriptToGoogleSheets(statement string, stmtId string, filename string) (string, tree.ParsingError) {

	log.Println(" Step: Parse input statement")
	// Explicitly activate printing of shared elements
	//exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Parse IGScript statement into tree
	s, err := parser.ParseStatement(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	// Run composite generation and return output and error. Will write file if filename != ""
	output, statementMap, statementHeader, statementHeaderNames, err := exporter.GenerateGoogleSheetsOutputFromParsedStatement(s, stmtId, "", tree.AGGREGATE_IMPLICIT_LINKAGES)

	fmt.Println(statementHeader)
	fmt.Println(statementHeaderNames)
	fmt.Println(statementMap)

	output, _ = exporter.GenerateGoogleSheetsOutput(statementMap, statementHeader, statementHeaderNames, filename)

	return output, err

}
