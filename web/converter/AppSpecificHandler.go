package converter

import (
	"IG-Parser/app"
	"IG-Parser/exporter"
	"IG-Parser/tree"
	"fmt"
	"log"
	"net/http"
)

/*
Third-level handler generating tabular output in response to web request.
Should be invoked by #converterHandler().
*/
func handleGoogleSheetsOutput(w http.ResponseWriter, codedStmt string, stmtId string, retStruct ReturnStruct, dynamicOutput bool, produceIGExtendedOutput bool, includeAnnotations bool) {
	// Run default configuration
	SetDefaultConfig()
	// Now, adjust to user settings based on UI output
	// Define whether output is dynamic
	fmt.Println("Setting dynamic output: ", dynamicOutput)
	exporter.SetDynamicOutput(dynamicOutput)
	// Define whether output is IG Extended (component-level nesting)
	fmt.Println("Setting IG Extended output: ", produceIGExtendedOutput)
	exporter.SetProduceIGExtendedOutput(produceIGExtendedOutput)
	// Define whether annotations are included
	fmt.Println("Setting annotations: ", includeAnnotations)
	exporter.SetIncludeAnnotations(includeAnnotations)
	// Convert input
	output, err2 := app.ConvertIGScriptToGoogleSheets(codedStmt, stmtId, "")
	if err2.ErrorCode != tree.PARSING_NO_ERROR {
		retStruct.Success = false
		retStruct.Error = true
		retStruct.CodedStmt = codedStmt
		switch err2.ErrorCode {
		case tree.PARSING_ERROR_EMPTY_LEAF:
			retStruct.Message = ERROR_INPUT_NO_STATEMENT
		default:
			retStruct.Message = "Parsing error (" + err2.ErrorCode + "): " + err2.ErrorMessage
		}
		err3 := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_SHEETS, retStruct)
		if err3 != nil {
			log.Println("Error processing default template:", err3.Error())
			http.Error(w, "Could not process request.", http.StatusInternalServerError)
		}

		// Final comment in log
		fmt.Println("Error: " + fmt.Sprint(err2))
		// Ensure logging is terminated
		err := terminateOutput(ERROR_SUFFIX)
		if err != nil {
			log.Println("Error when finalizing log file: ", err.Error())
		}

		return
	}
	// Return success if parsing was successful
	retStruct.Success = true
	retStruct.CodedStmt = codedStmt
	retStruct.Output = output
	err := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_SHEETS, retStruct)
	if err != nil {
		log.Println("Error processing default template:", err.Error())
		http.Error(w, "Could not process request.", http.StatusInternalServerError)
	}

	// Final comment in log
	fmt.Println("Success")
	// Ensure logging is terminated
	err3 := terminateOutput(SUCCESS_SUFFIX)
	if err3 != nil {
		log.Println("Error when finalizing log file: ", err3.Error())
	}

	return
}

/*
Third-level handler generating visual tree output in response to web request.
Should be invoked by #converterHandler().
*/
func handleVisualOutput(w http.ResponseWriter, codedStmt string, stmtId string, retStruct ReturnStruct, flatOutput bool, binaryOutput bool, dynamicOutput bool, produceIGExtendedOutput bool, includeAnnotations bool) {
	// Run default configuration
	SetDefaultConfig()
	// Now, adjust to user settings based on UI output
	// Define whether output is dynamic
	fmt.Println("Setting dynamic output: ", dynamicOutput)
	exporter.SetDynamicOutput(dynamicOutput)
	// Define whether output is IG Extended (component-level nesting)
	fmt.Println("Setting IG Extended output: ", produceIGExtendedOutput)
	exporter.SetProduceIGExtendedOutput(produceIGExtendedOutput)
	// Define whether annotations are included
	fmt.Println("Setting annotations: ", includeAnnotations)
	exporter.SetIncludeAnnotations(includeAnnotations)
	// Setting flat printing
	fmt.Println("Setting flat printing of properties: ", flatOutput)
	tree.SetFlatPrinting(flatOutput)
	fmt.Println("Setting binary tree printing: ", binaryOutput)
	tree.SetBinaryPrinting(binaryOutput)
	// Convert input
	output, err2 := app.ConvertIGScriptToVisualTree(codedStmt, stmtId, "")
	if err2.ErrorCode != tree.PARSING_NO_ERROR {
		retStruct.Success = false
		retStruct.Error = true
		retStruct.CodedStmt = codedStmt
		switch err2.ErrorCode {
		case tree.PARSING_ERROR_EMPTY_LEAF:
			retStruct.Message = ERROR_INPUT_NO_STATEMENT
		default:
			retStruct.Message = "Parsing error (" + err2.ErrorCode + "): " + err2.ErrorMessage
		}
		err3 := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_VISUAL, retStruct)
		if err3 != nil {
			log.Println("Error processing default template:", err3.Error())
			http.Error(w, "Could not process request.", http.StatusInternalServerError)
		}

		// Final comment in log
		fmt.Println("Error: " + fmt.Sprint(err2))
		// Ensure logging is terminated
		err := terminateOutput(ERROR_SUFFIX)
		if err != nil {
			log.Println("Error when finalizing log file: ", err.Error())
		}
		return
	}
	// Return success if parsing was successful
	retStruct.Success = true
	retStruct.CodedStmt = codedStmt
	retStruct.Output = output
	err := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_VISUAL, retStruct)
	if err != nil {
		log.Println("Error processing default template:", err.Error())
		http.Error(w, "Could not process request.", http.StatusInternalServerError)
	}

	// Final comment in log
	fmt.Println("Success")
	// Ensure logging is terminated
	err3 := terminateOutput(SUCCESS_SUFFIX)
	if err3 != nil {
		log.Println("Error when finalizing log file: ", err3.Error())
	}

	return
}
