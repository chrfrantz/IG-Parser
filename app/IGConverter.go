package app

import (
	"IG-Parser/exporter"
	"IG-Parser/parser"
	"IG-Parser/tree"
	"log"
)

/*
Consumes statement as input and produces outfile
Arguments include the IGScript-annotated statement, statement ID based on which substatements are generated,
and a filename for the output. If the filename is empty, no output will be written.
Returns Google Sheets output as string, and error (defaults to tree.PARSING_NO_ERROR).
*/
func ConvertIGScriptToGoogleSheets(statement string, stmtId string, filename string) (string, tree.ParsingError) {

	// Use separator specified by default
	separator := exporter.CellSeparator

	log.Println(" Step: Parse input statement")
	// Explicitly activate printing of shared elements
	//exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Parse IGScript statement into tree
	s, err := parser.ParseStatement(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	Println("Parsed statement:", s.String())

	// Run composite generation and return output and error. Will write file if filename != ""
	output, statementMap, statementHeader, statementHeaderNames, err := exporter.GenerateGoogleSheetsOutputFromParsedStatement(s, "", stmtId, "", tree.AGGREGATE_IMPLICIT_LINKAGES, separator)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	Println("Results:")
	Println(statementHeader)
	Println(statementHeaderNames)
	Println(statementMap)

	log.Println("Output generation complete.")
	log.Println("Writing to file ...")

	output, _ = exporter.GenerateGoogleSheetsOutput(statementMap, statementHeader, statementHeaderNames, separator, filename)

	log.Println("Writing completed.")

	return output, err

}

/*
Consumes statement as input and produces outfile reflecting visual tree structure consumable by D3.
Arguments include the IGScript-annotated statement, statement ID based on which substatements are generated,
and a filename for the output. If the filename is empty, no output will be written.
Returns Visual tree structure as string, and error (defaults to tree.PARSING_NO_ERROR).
*/
func ConvertIGScriptToVisualTree(statement string, stmtId string, filename string) (string, tree.ParsingError) {

	log.Println(" Step: Parse input statement")
	// Explicitly activate printing of shared elements
	//exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	// Parse IGScript statement into tree
	s, err := parser.ParseStatement(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	output := s.PrintTree(nil, exporter.IncludeAnnotations())

	Println("Generated visual tree:", output)

	log.Println("Output generation complete.")

	if filename != "" {
		log.Println("Writing to file ...")

		err2 := exporter.WriteToFile(filename, output)
		if err2 != nil {
			log.Println("Problems when writing file "+filename+", Error:", err2)
		}

		log.Println("Writing completed.")
	}

	return output, err

}
