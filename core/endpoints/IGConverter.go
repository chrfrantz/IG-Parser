package endpoints

import (
	"IG-Parser/core/exporter"
	"IG-Parser/core/parser"
	"IG-Parser/core/tree"
)

/*
This file contains the application endpoints that integrate the core parsing features, as well as file/output
handling. Both can be invoked with IG Script-encoded institutional statements to produce either tabular or
visual output for downstream processing, serving as endpoints for the use by specific applications, such as
web applications, console tools, etc.
*/

/*
Consumes statement as input and produces outfile.
Arguments include the IGScript-annotated statement, statement ID based on which substatements are generated,
the nature of the output type (see TabularOutputGeneratorConfig #OUTPUT_TYPE_CSV, #OUTPUT_TYPE_GOOGLE_SHEETS)
and a filename for the output. If the filename is empty, no output will be written. The parameter overwrite
indicates whether the target file will be overwritten upon repeated write.
If printHeaders is set, the output includes the header row.
Returns tabular output as string, and error (defaults to tree.PARSING_NO_ERROR).
*/
func ConvertIGScriptToTabularOutput(statement string, stmtId string, outputType string, filename string, overwrite bool, printHeaders bool) ([]exporter.TabularOutputResult, tree.ParsingError) {

	// Use separator specified by default
	separator := exporter.CellSeparator

	Println(" Step: Parse input statement")
	// Explicitly activate printing of shared elements
	//exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

	// Parse IGScript statement into tree
	stmts, err := parser.ParseStatement(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return nil, err
	}

	Println("Parsed statement:", stmts)

	// Run composite generation and return output and error. Will write file if filename != ""
	results := exporter.GenerateTabularOutputFromParsedStatements(stmts, "", stmtId, filename, overwrite, tree.AGGREGATE_IMPLICIT_LINKAGES, separator, outputType, printHeaders)
	for _, res := range results {
		if res.Error.ErrorCode != tree.PARSING_NO_ERROR {
			return results, res.Error
		}
	}

	Println("  - Results: ", results)

	Println("  - Output generation complete.")

	return results, err

}

/*
Consumes statement as input and produces outfile reflecting visual tree structure consumable by D3.
Arguments include the IGScript-annotated statement, statement ID (currently not used in visualization),
and a filename for the output. If the filename is empty, no output will be written.
Returns Visual tree structure as string, and error (defaults to tree.PARSING_NO_ERROR).
*/
func ConvertIGScriptToVisualTree(statement string, stmtId string, filename string) (string, tree.ParsingError) {

	Println(" Step: Parse input statement")

	// Parse IGScript statement into tree
	stmts, err := parser.ParseStatement(statement)
	if err.ErrorCode != tree.PARSING_NO_ERROR {
		return "", err
	}

	output := ""

	err2 := tree.NodeError{}

	// Prepare visual output for nodes
	Println(" Step: Generate visual output structure (combined statements)")
	output, err2 = stmts[0].PrintNodeTree(nil, tree.FlatPrinting(), tree.BinaryPrinting(), exporter.IncludeAnnotations(),
		//output, err2 = tree.PrintHeader(stmts[0], tree.FlatPrinting(), tree.BinaryPrinting(), exporter.IncludeAnnotations(),
		exporter.IncludeDegreeOfVariability(), tree.MoveActivationConditionsToFront(), 0)
	if err2.ErrorCode != tree.TREE_NO_ERROR {
		return output, tree.ParsingError{ErrorCode: tree.PARSING_ERROR_EMBEDDED_NODE_ERROR, ErrorMessage: err2.ErrorMessage}
	}

	Println("  - Generated visual tree:", output)

	Println("  - Output generation complete.")

	if filename != "" {
		Println("  - Writing to file ...")

		err3 := exporter.WriteToFile(filename, output, true)
		if err3 != nil {
			Println("  - Problems when writing file "+filename+", Error:", err3)
		}

		Println("  - Writing completed.")
	}

	return output, err
}
