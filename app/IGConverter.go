package app

import (
	"IG-Parser/exporter"
	"IG-Parser/parser"
	"IG-Parser/tree"
	"log"
)

/*
Consumes statements as input and produces outfile
Arguments include the IGScript-annotated statement, statement ID based on which substatements are generated,
and a filename for the output. If the filename is empty, no output will be written.
Returns Google Sheets output as string, and error (defaults to tree.PARSING_NO_ERROR).
 */
func ConvertIGScriptToGoogleSheets(statement string, stmtId string, filename string) (string, tree.ParsingError) {

	log.Println("Step: Parse input statement")
	// Parse IGScript statement into tree
	s, err := parser.ParseStatement(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	// Run composite generation and return output and error. Will write file if filename != ""
	return exporter.GenerateGoogleSheetsOutputFromParsedStatement(s, stmtId, filename)

}
