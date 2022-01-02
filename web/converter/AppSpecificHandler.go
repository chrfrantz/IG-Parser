package converter

import (
	"IG-Parser/app"
	"IG-Parser/exporter"
	"IG-Parser/tree"
	"IG-Parser/web/helper"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

/*
Second-level general handler that retrieves and preprocesses information from input.
Delegates to third-order handler for output-specific generation.
Should be invoked by #ConverterHandlerSheets() and #ConverterHandlerVisual().
*/
func converterHandler(w http.ResponseWriter, r *http.Request, templateName string) {

	// Prepopulate response
	message := ""
	transactionID := ""
	rawStmt := r.FormValue("rawStatement")
	codedStmt := r.FormValue("annotatedStatement")
	stmtId := r.FormValue("stmtId")
	dynChk := r.FormValue("dynamicOutput")
	inclAnnotations := r.FormValue("annotations")
	igExtended := r.FormValue("compLevelNesting")
	propertyTree := r.FormValue("propertyTree")
	binaryTree := r.FormValue("binaryTree")
	heightValue := r.FormValue("heightValue")
	widthValue := r.FormValue("widthValue")

	// Dynamic output
	dynamicOutput := false
	fmt.Println("Dynamic: ", dynChk)
	if dynChk == "on" {
		dynChk = "checked"
		dynamicOutput = true
	} else {
		dynChk = "unchecked"
		dynamicOutput = false
	}

	// Annotations in output
	includeAnnotations := false
	fmt.Println("Annotations: ", inclAnnotations)
	if inclAnnotations == "on" {
		inclAnnotations = "checked"
		includeAnnotations = true
	} else {
		inclAnnotations = "unchecked"
		includeAnnotations = false
	}

	// Component-level nesting in output
	produceIGExtendedOutput := false
	fmt.Println("IG Extended output: ", igExtended)
	if igExtended == "on" {
		igExtended = "checked"
		produceIGExtendedOutput = true
	} else {
		igExtended = "unchecked"
		produceIGExtendedOutput = false
	}

	// Private property printing in output
	printFlatProperties := false
	fmt.Println("Private property tree printing: ", propertyTree)
	if propertyTree == "on" {
		propertyTree = "checked"
		printFlatProperties = false
	} else {
		propertyTree = "unchecked"
		printFlatProperties = true
	}

	// Binary tree printing in output
	printBinaryTree := false
	fmt.Println("Binary tree printing: ", binaryTree)
	if binaryTree == "on" {
		binaryTree = "checked"
		printBinaryTree = true
	} else {
		binaryTree = "unchecked"
		printBinaryTree = false
	}

	retStruct := ReturnStruct{
		Success:            false,
		Error:              false,
		Message:            message,
		RawStmt:            rawStmt,
		CodedStmt:          codedStmt,
		StmtId:             stmtId,
		DynamicOutput:      dynChk,
		IGExtendedOutput:   igExtended,
		IncludeAnnotations: inclAnnotations,
		PrintPropertyTree:  propertyTree,
		PrintBinaryTree:    binaryTree,
		Width:              WIDTH,
		Height:             HEIGHT,
		TransactionId:      transactionID,
		RawStmtHelp:        HELP_RAW_STMT,
		CodedStmtHelp:      HELP_CODED_STMT,
		StmtIdHelp:         HELP_STMT_ID,
		ParametersHelp:     HELP_PARAMETERS,
		ReportHelp:         HELP_REPORT}

	// Assign width for UI rendering
	if widthValue != "" {
		widthVal, err := strconv.Atoi(widthValue)
		if err != nil || widthVal < 100 {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = ERROR_INPUT_WIDTH
			err := tmpl.ExecuteTemplate(w, templateName, retStruct)
			if err != nil {
				log.Println("Error generating error response for empty input:", err.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}
		}
		retStruct.Width = widthVal
	}

	// Assign height for UI rendering
	if heightValue != "" {
		heightVal, err := strconv.Atoi(heightValue)
		if err != nil || heightVal < 100 {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = ERROR_INPUT_HEIGHT
			err := tmpl.ExecuteTemplate(w, templateName, retStruct)
			if err != nil {
				log.Println("Error generating error response for empty input:", err.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}
		}
		retStruct.Height = heightVal
	}

	// All parameters are processed ...

	if r.Method != http.MethodPost {
		// Just show empty form with prepopulated elements
		retStruct.RawStmt = RAW_STATEMENT
		retStruct.CodedStmt = ANNOTATED_STATEMENT
		retStruct.StmtId = STATEMENT_ID

		// Check for parameters that customize input
		// Set switch to indicate potential need to align raw and coded statement field entries
		resetValues := false
		// Parameter: Raw Statement
		keys, ok := r.URL.Query()[PARAM_RAW_STATEMENT]
		if ok && len(keys[0]) > 0 {

			// Assume single item
			key := keys[0]

			//log.Println("Url Param 'rawStmt' is: " + string(key))
			// Assign value instead
			retStruct.RawStmt = string(key)
			// Set switch to indicate reset of raw statement if not specified as parameter
			resetValues = true
		}

		// Parameter: IG Script-coded statement - consider interaction with raw statement
		keys, ok = r.URL.Query()[PARAM_CODED_STATEMENT]
		if ok && len(keys[0]) > 0 {

			// Assume single item
			key := keys[0]

			//log.Println("Url Param 'codedStmt' is: " + string(key))
			// Assign value instead
			retStruct.CodedStmt = string(key)
			// Check for raw statement if it is still default; then reset
			if retStruct.RawStmt == RAW_STATEMENT {
				retStruct.RawStmt = ""
			}
		} else if resetValues {
			// Reset value, since the default coded statement will likely not correspond.
			retStruct.CodedStmt = ""
		}

		// Parameter: Statement ID
		keys, ok = r.URL.Query()[PARAM_STATEMENT_ID]
		if ok && len(keys[0]) > 0 {

			// Assume single item
			key := keys[0]

			//log.Println("Url Param 'stmtId' is: " + string(key))
			// Assign value instead
			retStruct.StmtId = string(key)
		}

		err := tmpl.ExecuteTemplate(w, templateName, retStruct)
		if err != nil {
			log.Println("Error processing default template:", err.Error())
			http.Error(w, "Could not process request.", http.StatusInternalServerError)
		}
		return
	}

	// Initialize request-specific logfile first
	if Logging {
		log.Println("Logging enabled")
		tID, filename := helper.GenerateUniqueIdAndFilename()
		// Assign transaction ID
		retStruct.TransactionId = tID
		// Check whether logging path has terminating slash
		if LoggingPath != "" {
			if LoggingPath[len(LoggingPath)-1:] != "/" {
				LoggingPath += "/"
			}
		}
		// Perform the file redirection
		var err error
		terminateOutput, err = helper.SaveOutput(LoggingPath + filename)

		fmt.Println("TRANSACTION ID: " + retStruct.TransactionId)
		if err != nil {
			fmt.Println("Error when initializing logging: " + err.Error())
		}
	}

	fmt.Println("Input values:\n" +
		"Raw statement: " + retStruct.RawStmt + "\n" +
		"Annotated statement: " + retStruct.CodedStmt + "\n" +
		"Full input value struct: " + fmt.Sprint(retStruct))

	// Check for input statement first
	if codedStmt == "" {
		retStruct.Success = false
		retStruct.Error = true
		retStruct.Message = ERROR_INPUT_NO_STATEMENT
		err := tmpl.ExecuteTemplate(w, templateName, retStruct)
		if err != nil {
			log.Println("Error generating error response for empty input:", err.Error())
			http.Error(w, "Could not process request.", http.StatusInternalServerError)
		}

		// Final comment in log
		fmt.Println("Error: No input to parse.")
		// Ensure logging is terminated
		err2 := terminateOutput(ERROR_SUFFIX)
		if err2 != nil {
			log.Println("Error when finalizing log file: ", err2.Error())
		}

		return
	} else {
		// Produce return actual output
		if templateName == TEMPLATE_NAME_PARSER_SHEETS {
			handleGoogleSheetsOutput(w, codedStmt, stmtId, retStruct, dynamicOutput, produceIGExtendedOutput, includeAnnotations)
		} else if templateName == TEMPLATE_NAME_PARSER_VISUAL {
			fmt.Println("Visual output requested")
			handleVisualOutput(w, codedStmt, stmtId, retStruct, printFlatProperties, printBinaryTree, dynamicOutput, produceIGExtendedOutput, includeAnnotations)
		} else {
			log.Fatal("Output variant " + templateName + " not found.")
		}

	}
}

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
