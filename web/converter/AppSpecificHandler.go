package converter

import (
	"IG-Parser/core/endpoints"
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/core/tree"
	"IG-Parser/web/converter/shared"
	"fmt"
	"log"
	"net/http"
)

/*
Third-level handler generating tabular output in response to web request.
Should be invoked by #converterHandler().
*/
func handleTabularOutput(w http.ResponseWriter, codedStmt string, stmtId string, retStruct shared.ReturnStruct, dynamicOutput bool, produceIGExtendedOutput bool, includeAnnotations bool, outputType string, printHeaders bool) {
	// Run default configuration
	shared.SetDefaultConfig()
	// Now, adjust to user settings based on UI output
	// Define whether output is dynamic
	Println("Setting dynamic output:", dynamicOutput)
	tabular.SetDynamicOutput(dynamicOutput)
	// Define whether output is IG Extended (component-level nesting)
	Println("Setting IG Extended output:", produceIGExtendedOutput)
	tabular.SetProduceIGExtendedOutput(produceIGExtendedOutput)
	// Define whether annotations are included
	Println("Setting annotations:", includeAnnotations)
	tabular.SetIncludeAnnotations(includeAnnotations)
	// Define whether header row is included
	Println("Setting header row:", printHeaders)
	tabular.SetIncludeHeaders(printHeaders)
	// Output type
	Println("Output type:", outputType)
	// Convert input
	output, err2 := endpoints.ConvertIGScriptToTabularOutput(codedStmt, stmtId, outputType, "", true, tabular.IncludeHeader())
	if err2.ErrorCode != tree.PARSING_NO_ERROR {
		retStruct.Success = false
		retStruct.Error = true
		retStruct.CodedStmt = codedStmt
		// Deal with potential errors and prepopulate return message
		switch err2.ErrorCode {
		case tree.PARSING_ERROR_EMPTY_LEAF:
			retStruct.Message = shared.ERROR_INPUT_NO_STATEMENT
		default:
			retStruct.Message = "Parsing error (" + err2.ErrorCode + "): " + err2.ErrorMessage
		}
		// Execute template
		err3 := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_TABULAR, retStruct)
		if err3 != nil {
			log.Println("Error processing default template:", err3.Error())
			http.Error(w, "Could not process request.", http.StatusInternalServerError)
		}

		// Final comment in log
		Println("Error: " + fmt.Sprint(err2))
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
	tabularOutput := ""
	for _, v := range output {
		tabularOutput += v.Output
	}
	retStruct.Output = tabularOutput
	err := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_TABULAR, retStruct)
	if err != nil {
		log.Println("Error processing default template:", err.Error())
		http.Error(w, "Could not process request.", http.StatusInternalServerError)
	}

	// Final comment in log
	Println("Success")
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
func handleVisualOutput(w http.ResponseWriter, codedStmt string, stmtId string, retStruct shared.ReturnStruct, flatOutput bool, binaryOutput bool, moveActivationConditionsToTop bool, dynamicOutput bool, produceIGExtendedOutput bool, includeAnnotations bool, includeDoV bool) {
	// Run default configuration
	shared.SetDefaultConfig()
	// Now, adjust to user settings based on UI output
	// Define whether output is dynamic
	Println("Setting dynamic output:", dynamicOutput)
	tabular.SetDynamicOutput(dynamicOutput)
	// Define whether output is IG Extended (component-level nesting)
	Println("Setting IG Extended output:", produceIGExtendedOutput)
	tabular.SetProduceIGExtendedOutput(produceIGExtendedOutput)
	// Define whether annotations are included
	Println("Setting annotations:", includeAnnotations)
	tabular.SetIncludeAnnotations(includeAnnotations)
	// Define whether Degree of Variability is included
	Println("Setting Degree of Variability (DoV):", includeDoV)
	tabular.SetIncludeDegreeOfVariability(includeDoV)
	// Setting flat printing
	Println("Setting flat printing of properties:", flatOutput)
	tree.SetFlatPrinting(flatOutput)
	Println("Setting binary tree printing:", binaryOutput)
	tree.SetBinaryPrinting(binaryOutput)
	Println("Setting activation condition on top in visual output:", moveActivationConditionsToTop)
	tree.SetMoveActivationConditionsToFront(moveActivationConditionsToTop)
	// Convert input
	output, err2 := endpoints.ConvertIGScriptToVisualTree(codedStmt, stmtId, "")
	if err2.ErrorCode != tree.PARSING_NO_ERROR {
		retStruct.Success = false
		retStruct.Error = true
		retStruct.CodedStmt = codedStmt
		switch err2.ErrorCode {
		case tree.PARSING_ERROR_EMPTY_LEAF:
			retStruct.Message = shared.ERROR_INPUT_NO_STATEMENT
		default:
			retStruct.Message = "Parsing error (" + err2.ErrorCode + "): " + err2.ErrorMessage
		}
		err3 := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_PARSER_VISUAL, retStruct)
		if err3 != nil {
			log.Println("Error processing default template:", err3.Error())
			http.Error(w, "Could not process request.", http.StatusInternalServerError)
		}

		// Final comment in log
		Println("Error: " + fmt.Sprint(err2))
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
	Println("Success")
	// Ensure logging is terminated
	err3 := terminateOutput(SUCCESS_SUFFIX)
	if err3 != nil {
		log.Println("Error when finalizing log file: ", err3.Error())
	}
	return
}
