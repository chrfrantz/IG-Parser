package converter

import (
	"IG-Parser/core/endpoints"
	"IG-Parser/core/exporter/tabular"
	"IG-Parser/core/tree"
	"IG-Parser/web/converter/shared"
	"fmt"
	"log"
	"net/http"
	"strings"
)

/*
Contains output-specific handler to be invoked by GenericParserHandler.go.
*/

/*
Third-level handler generating tabular output in response to web request.
Should be invoked by #converterHandler().
*/
func handleTabularOutput(w http.ResponseWriter, originalStatement string, codedStmt string, stmtId string, retStruct shared.ReturnStruct, dynamicOutput bool, produceIGExtendedOutput bool, includeAnnotations bool, outputType string, printHeaders bool, printOriginalStatement string, printIgScriptInput string) {
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
	// Indicate whether Original Statement input is included in output
	Println("Include Original Statement input in generated output:", printOriginalStatement)
	// Indicate whether IG Script input is included in output
	Println("Include IG Script input in generated output:", printIgScriptInput)
	// Output type
	Println("Output type:", outputType)
	// Prepopulate coded statement in return structure
	retStruct.CodedStmt = codedStmt
	// Convert input
	output, err2 := endpoints.ConvertIGScriptToTabularOutput(originalStatement, codedStmt, stmtId, outputType, "", true, tabular.IncludeHeader(), printOriginalStatement, printIgScriptInput)
	// Stringified output delivered back to client in case of no error or warning
	finalOutput := ""
	if err2.ErrorCode == tree.PARSING_NO_ERROR || err2.ErrorCode == tree.PARSING_WARNING_POSSIBLY_NON_PARSED_CONTENT {
		tabularOutput := ""
		// Decompose array output into string
		for _, v := range output {
			tabularOutput += v.Output
		}
		finalOutput = tabularOutput
	}
	// Deliver parsed content back to client
	deliverParsedOutput(w, retStruct, TEMPLATE_NAME_PARSER_TABULAR, finalOutput, err2)
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
	// Prepopulate coded statement in return structure
	retStruct.CodedStmt = codedStmt
	// Convert input
	output, err2 := endpoints.ConvertIGScriptToVisualTree(codedStmt, stmtId, "")
	// Deliver parsed content back to client
	deliverParsedOutput(w, retStruct, TEMPLATE_NAME_PARSER_VISUAL, output, err2)
}

/*
Processes parsed output and delivers it back to client. Takes input from any kind of parser.
- w: Response writer to deliver output to client
- retStruct: Return structure to be populated with parsed output and error information
- template: Name of template to be used for output
- output: Stringified output to be delivered
- parsingError: Parsing error object to be processed in frontend
*/
func deliverParsedOutput(w http.ResponseWriter, retStruct shared.ReturnStruct, template string, output string, parsingError tree.ParsingError) {
	if parsingError.ErrorCode != tree.PARSING_NO_ERROR {
		retStruct.Success = false
		retStruct.Error = true
		// Deal with potential errors and prepopulate return message
		switch parsingError.ErrorCode {
		case tree.PARSING_ERROR_EMPTY_LEAF:
			retStruct.Message = shared.ERROR_INPUT_NO_STATEMENT
		case tree.PARSING_ERROR_EMPTY_STATEMENT:
			retStruct.Message = shared.ERROR_INPUT_NO_STATEMENT
		case tree.PARSING_ERROR_IGNORED_NESTED_ELEMENTS:
			// Parts that have been ignored in parsing due to errors
			retStruct.Message = shared.ERROR_INPUT_IGNORED_ELEMENTS + "\"" + strings.Join(parsingError.ErrorIgnoredElements, ", ") + "\""
		case tree.PARSING_WARNING_POSSIBLY_NON_PARSED_CONTENT:
			// Parts that have been ignored in parsing without errors, but could potentially needed to be parsed
			retStruct.Message = shared.WARNING_INPUT_NON_PARSED_ELEMENTS + "\"" + strings.Join(parsingError.ErrorIgnoredElements, ", ") + "\""
			// Still allow it to show
			retStruct.Success = true
			retStruct.Output = output
		default:
			retStruct.Message = "Parsing error (" + parsingError.ErrorCode + "): " + parsingError.ErrorMessage
		}
		// Execute template
		err3 := tmpl.ExecuteTemplate(w, template, retStruct)
		if err3 != nil {
			log.Println("Error processing template:", err3.Error())
			http.Error(w, "Could not process request.", http.StatusInternalServerError)
		}

		// Determine suffix for final output in log (warning or error)
		errSuffix := ""
		if retStruct.Success {
			// Warning
			Println("Warning: " + fmt.Sprint(parsingError))
			errSuffix = WARNING_SUFFIX
		} else {
			// Hard error
			Println("Error: " + fmt.Sprint(parsingError))
			errSuffix = ERROR_SUFFIX
		}
		// Ensure logging is terminated
		err := terminateOutput(errSuffix)
		if err != nil {
			log.Println("Error when finalizing log file: ", err.Error())
		}
		return
	}
	// Return success if parsing was successful
	retStruct.Success = true
	retStruct.Output = output
	err := tmpl.ExecuteTemplate(w, template, retStruct)
	if err != nil {
		log.Println("Error processing template:", err.Error())
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
