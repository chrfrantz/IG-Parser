package converter

import (
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

	// Reading form to prepopulate response
	message := ""
	transactionID := ""
	rawStmt := r.FormValue(PARAM_RAW_STATEMENT)
	codedStmt := r.FormValue(PARAM_CODED_STATEMENT)
	stmtId := r.FormValue(PARAM_STATEMENT_ID)
	dynChk := r.FormValue(PARAM_DYNAMIC_SCHEMA)
	inclAnnotations := r.FormValue(PARAM_LOGICO_OUTPUT)
	igExtended := r.FormValue(PARAM_EXTENDED_OUTPUT)
	propertyTree := r.FormValue(PARAM_PROPERTY_TREE)
	binaryTree := r.FormValue(PARAM_BINARY_TREE)
	heightValue := r.FormValue(PARAM_HEIGHT)
	widthValue := r.FormValue(PARAM_WIDTH)

	// EVALUATE INDIVIDUAL CHECKBOX INPUTS

	// Dynamic output
	dynamicOutput := false
	fmt.Println("Dynamic: ", dynChk)
	if dynChk == CHECKBOX_ON {
		dynChk = CHECKBOX_CHECKED
		dynamicOutput = true
	} else {
		dynChk = CHECKBOX_UNCHECKED
		dynamicOutput = false
	}

	// Annotations in output
	includeAnnotations := false
	fmt.Println("Annotations: ", inclAnnotations)
	if inclAnnotations == CHECKBOX_ON {
		inclAnnotations = CHECKBOX_CHECKED
		includeAnnotations = true
	} else {
		inclAnnotations = CHECKBOX_UNCHECKED
		includeAnnotations = false
	}

	// Component-level nesting in output
	produceIGExtendedOutput := false
	fmt.Println("IG Extended output: ", igExtended)
	if igExtended == CHECKBOX_ON {
		igExtended = CHECKBOX_CHECKED
		produceIGExtendedOutput = true
	} else {
		igExtended = CHECKBOX_UNCHECKED
		produceIGExtendedOutput = false
	}

	// Private property printing in output
	printFlatProperties := false
	fmt.Println("Private property tree printing: ", propertyTree)
	if propertyTree == CHECKBOX_ON {
		propertyTree = CHECKBOX_CHECKED
		printFlatProperties = false
	} else {
		propertyTree = CHECKBOX_UNCHECKED
		printFlatProperties = true
	}

	// Binary tree printing in output
	printBinaryTree := false
	fmt.Println("Binary tree printing: ", binaryTree)
	if binaryTree == CHECKBOX_ON {
		binaryTree = CHECKBOX_CHECKED
		printBinaryTree = true
	} else {
		binaryTree = CHECKBOX_UNCHECKED
		printBinaryTree = false
	}

	// Checkbox interpretation finished

	// Prepare return structure with prepopulated information (to be refined prior to return)
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
		if err != nil || widthVal < MIN_WIDTH {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = ERROR_INPUT_WIDTH
			err := tmpl.ExecuteTemplate(w, templateName, retStruct)
			if err != nil {
				log.Println("Error generating error response for template processing:", err.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}
			// Stop execution and return to UI
			return
		}
		// Assign input value only in case of success
		retStruct.Width = widthVal
	}

	// Assign height for UI rendering
	if heightValue != "" {
		heightVal, err := strconv.Atoi(heightValue)
		if err != nil || heightVal < MIN_HEIGHT {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = ERROR_INPUT_HEIGHT
			err := tmpl.ExecuteTemplate(w, templateName, retStruct)
			if err != nil {
				log.Println("Error generating error response for template processing:", err.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}
			// Stop execution and return to UI
			return
		}
		// Assign input value only in case of success
		retStruct.Height = heightVal
	}

	// All form parameters are processed ...

	if r.Method != http.MethodPost {
		// Just show empty form with prepopulated elements
		retStruct.RawStmt = RAW_STATEMENT
		retStruct.CodedStmt = ANNOTATED_STATEMENT
		retStruct.StmtId = STATEMENT_ID

		// CHECK FOR URL PARAMETERS TO CUSTOMIZE EXECUTION

		// Set switch to indicate potential need to align raw and coded statement field entries
		resetValues := false

		// Parameter: Raw Statement
		val, suc := extractUrlParameters(r, PARAM_RAW_STATEMENT)
		if suc {
			retStruct.RawStmt = val
			resetValues = true
		}

		// Parameter: IG Script-coded statement - consider interaction with raw statement
		val, suc = extractUrlParameters(r, PARAM_CODED_STATEMENT)
		if suc {
			retStruct.CodedStmt = val
			// Check for raw statement if it is still default, while the coded stmt changed; if so, then reset
			if retStruct.CodedStmt != ANNOTATED_STATEMENT && retStruct.RawStmt == RAW_STATEMENT {
				retStruct.RawStmt = ""
			}
		} else if resetValues {
			// Reset value, since the default coded statement will likely not correspond.
			retStruct.CodedStmt = ""
		}

		// Parameter: Statement ID
		val, suc = extractUrlParameters(r, PARAM_STATEMENT_ID)
		if suc {
			retStruct.StmtId = val
		}

		// Parameter: IG Logico annotations
		val, suc = extractUrlParameters(r, PARAM_LOGICO_OUTPUT)
		check := evaluateBooleanUrlParameters(PARAM_LOGICO_OUTPUT, val, suc)
		// Assign values
		if check {
			retStruct.IncludeAnnotations = CHECKBOX_CHECKED
			includeAnnotations = true
		} else {
			retStruct.IncludeAnnotations = CHECKBOX_UNCHECKED
			includeAnnotations = false
		}

		// Parameter: Dynamic output
		val, suc = extractUrlParameters(r, PARAM_DYNAMIC_SCHEMA)
		check = evaluateBooleanUrlParameters(PARAM_DYNAMIC_SCHEMA, val, suc)
		// Assign values
		if check {
			retStruct.DynamicOutput = CHECKBOX_CHECKED
			dynamicOutput = true
		} else {
			retStruct.DynamicOutput = CHECKBOX_UNCHECKED
			dynamicOutput = false
		}

		// Parameter: Component-level nesting (IG Extended output)
		val, suc = extractUrlParameters(r, PARAM_EXTENDED_OUTPUT)
		check = evaluateBooleanUrlParameters(PARAM_EXTENDED_OUTPUT, val, suc)
		// Assign values
		if check {
			retStruct.IGExtendedOutput = CHECKBOX_CHECKED
			produceIGExtendedOutput = true
		} else {
			retStruct.IGExtendedOutput = CHECKBOX_UNCHECKED
			produceIGExtendedOutput = false
		}

		// Parameter: Property tree
		val, suc = extractUrlParameters(r, PARAM_PROPERTY_TREE)
		check = evaluateBooleanUrlParameters(PARAM_PROPERTY_TREE, val, suc)
		// Assign values
		if check {
			retStruct.PrintPropertyTree = CHECKBOX_CHECKED
			printFlatProperties = false
		} else {
			retStruct.PrintPropertyTree = CHECKBOX_UNCHECKED
			printFlatProperties = true
		}

		// Parameter: Binary tree
		val, suc = extractUrlParameters(r, PARAM_BINARY_TREE)
		check = evaluateBooleanUrlParameters(PARAM_BINARY_TREE, val, suc)
		// Assign values
		if check {
			retStruct.PrintBinaryTree = CHECKBOX_CHECKED
			printBinaryTree = true
		} else {
			retStruct.PrintBinaryTree = CHECKBOX_UNCHECKED
			printBinaryTree = false
		}

		// Parameter: Canvas width
		val, suc = extractUrlParameters(r, PARAM_WIDTH)
		if suc {
			w, err := strconv.Atoi(val)
			if err != nil {
				log.Println("Error when interpreting URL parameter '"+PARAM_WIDTH+"':", err)
			} else {
				retStruct.Width = w
			}
		}

		// Parameter: Canvas height
		val, suc = extractUrlParameters(r, PARAM_HEIGHT)
		if suc {
			h, err := strconv.Atoi(val)
			if err != nil {
				log.Println("Error when interpreting URL parameter '"+PARAM_HEIGHT+"':", err)
			} else {
				retStruct.Height = h
			}
		}

		// Parameter: Execution - invokes immediate execution
		val, suc = extractUrlParameters(r, PARAM_EXECUTE_PARSER)
		check = evaluateBooleanUrlParameters(PARAM_EXECUTE_PARSER, val, suc)
		if !check {

			// All URL parameters processed, but this is only returned if no immediate execution is requested ...

			// BY DEFAULT, THIS IS EXECUTED

			err := tmpl.ExecuteTemplate(w, templateName, retStruct)
			if err != nil {
				log.Println("Error processing default template:", err.Error())
				http.Error(w, "Could not process request.", http.StatusInternalServerError)
			}
			return
		}
		log.Println("Automated execution of parser requested via URL parameter.")
		// if reaching here, automated execution is requested via URL parameters ...
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
	if retStruct.CodedStmt == "" {
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
			handleGoogleSheetsOutput(w, retStruct.CodedStmt, retStruct.StmtId, retStruct, dynamicOutput, produceIGExtendedOutput, includeAnnotations)
		} else if templateName == TEMPLATE_NAME_PARSER_VISUAL {
			fmt.Println("Visual output requested")
			handleVisualOutput(w, retStruct.CodedStmt, retStruct.StmtId, retStruct, printFlatProperties, printBinaryTree, dynamicOutput, produceIGExtendedOutput, includeAnnotations)
		} else {
			log.Fatal("Output variant " + templateName + " not found.")
		}
	}
}

/*
Extract URL parameters for further processing. It takes the parameter key (as specified in UrlParameters.go) as input
and returns the associated value, as well as a flag indicating success of extraction (i.e., whether a parameter had
been present in the first place).
*/
func extractUrlParameters(r *http.Request, parameterKey string) (string, bool) {
	keys, ok := r.URL.Query()[parameterKey]
	if ok && len(keys[0]) > 0 {

		// Assume single item
		key := keys[0]

		// Return entry as string and signal success (even if value is empty)
		return string(key), true
	}
	// Return empty string (and signal absence of URL parameter)
	return "", false
}

/*
Evaluates URL parameter value input for boolean variants:
- "true", "t", and "1" are interpreted as "on"
- "false", "f", and "0" are interpreted as "off"
- If success is false, the function returns false (i.e., no URL parameter of the given name found).
*/
func evaluateBooleanUrlParameters(parameter string, value string, success bool) bool {
	if success {
		switch value {
		case "t":
			return true
		case "true":
			return true
		case "1":
			return true
		case "f":
			return false
		case "false":
			return false
		case "0":
			return false
		default:
			log.Println("Invalid URL parameter value for parameter '" + parameter + "': " + value)
		}
	}
	return false
}
