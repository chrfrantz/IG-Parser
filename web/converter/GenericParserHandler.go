package converter

import (
	"IG-Parser/exporter"
	"IG-Parser/web/helper"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
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
	formValueRawStmt := r.FormValue(PARAM_RAW_STATEMENT)
	formValueCodedStmt := r.FormValue(PARAM_CODED_STATEMENT)
	formValueStmtId := r.FormValue(PARAM_STATEMENT_ID)
	formValueDynamicOutput := r.FormValue(PARAM_DYNAMIC_SCHEMA)
	formValueIncludeAnnotations := r.FormValue(PARAM_LOGICO_OUTPUT)
	formValueIncludeDoV := r.FormValue(PARAM_DOV)
	formValueIgExtendedOutput := r.FormValue(PARAM_EXTENDED_OUTPUT)
	formValueIncludeHeaders := r.FormValue(PARAM_PRINT_HEADERS)
	formValueOutputType := r.FormValue(PARAM_OUTPUT_TYPE)
	formValuePropertyTree := r.FormValue(PARAM_PROPERTY_TREE)
	formValueBinaryTree := r.FormValue(PARAM_BINARY_TREE)
	formValueMoveActivationConditionsToTop := r.FormValue(PARAM_ACTIVATION_CONDITION_ON_TOP)
	formValueCanvasHeightValue := r.FormValue(PARAM_HEIGHT)
	formValueCanvasWidthValue := r.FormValue(PARAM_WIDTH)

	// EVALUATE INDIVIDUAL CHECKBOX INPUTS

	// Dynamic output
	dynamicOutput := false
	fmt.Println("Dynamic: ", formValueDynamicOutput)
	if formValueDynamicOutput == CHECKBOX_ON {
		formValueDynamicOutput = CHECKBOX_CHECKED
		dynamicOutput = true
	} else {
		formValueDynamicOutput = CHECKBOX_UNCHECKED
		dynamicOutput = false
	}

	// Annotations in output
	includeAnnotations := false
	fmt.Println("Annotations: ", formValueIncludeAnnotations)
	if formValueIncludeAnnotations == CHECKBOX_ON {
		formValueIncludeAnnotations = CHECKBOX_CHECKED
		includeAnnotations = true
	} else {
		formValueIncludeAnnotations = CHECKBOX_UNCHECKED
		includeAnnotations = false
	}

	// DoV in output
	includeDoV := false
	fmt.Println("DoV: ", formValueIncludeDoV)
	if formValueIncludeDoV == CHECKBOX_ON {
		formValueIncludeDoV = CHECKBOX_CHECKED
		includeDoV = true
	} else {
		formValueIncludeDoV = CHECKBOX_UNCHECKED
		includeDoV = false
	}

	// Component-level nesting in output
	produceIGExtendedOutput := false
	fmt.Println("IG Extended output: ", formValueIgExtendedOutput)
	if formValueIgExtendedOutput == CHECKBOX_ON {
		formValueIgExtendedOutput = CHECKBOX_CHECKED
		produceIGExtendedOutput = true
	} else {
		formValueIgExtendedOutput = CHECKBOX_UNCHECKED
		produceIGExtendedOutput = false
	}

	// Print headers in output
	printHeaders := false
	// If not received by POST, set headers as default setting
	if formValueIncludeHeaders == "" && r.Method != http.MethodPost {
		formValueIncludeHeaders = CHECKBOX_ON
	}
	fmt.Println("Include headers in output: ", formValueIncludeHeaders)
	if formValueIncludeHeaders == CHECKBOX_ON {
		formValueIncludeHeaders = CHECKBOX_CHECKED
		printHeaders = true
	} else {
		formValueIncludeHeaders = CHECKBOX_UNCHECKED
		printHeaders = false
	}

	// Private property printing in output
	printFlatProperties := false
	fmt.Println("Private property tree printing: ", formValuePropertyTree)
	if formValuePropertyTree == CHECKBOX_ON {
		formValuePropertyTree = CHECKBOX_CHECKED
		printFlatProperties = false
	} else {
		formValuePropertyTree = CHECKBOX_UNCHECKED
		printFlatProperties = true
	}

	// Binary tree printing in output
	printBinaryTree := false
	fmt.Println("Binary tree printing: ", formValueBinaryTree)
	if formValueBinaryTree == CHECKBOX_ON {
		formValueBinaryTree = CHECKBOX_CHECKED
		printBinaryTree = true
	} else {
		formValueBinaryTree = CHECKBOX_UNCHECKED
		printBinaryTree = false
	}

	// Activation condition on top in output
	printActivationConditionsOnTop := false
	fmt.Println("Activation conditions on top: ", formValueMoveActivationConditionsToTop)
	if formValueMoveActivationConditionsToTop == CHECKBOX_ON {
		formValueMoveActivationConditionsToTop = CHECKBOX_CHECKED
		printActivationConditionsOnTop = true
	} else {
		formValueMoveActivationConditionsToTop = CHECKBOX_UNCHECKED
		printActivationConditionsOnTop = false
	}

	// Checkbox interpretation finished

	// Prepare return structure with prepopulated information (to be refined prior to return)
	retStruct := ReturnStruct{
		Success:                   false,
		Error:                     false,
		Message:                   message,
		RawStmt:                   formValueRawStmt,
		CodedStmt:                 formValueCodedStmt,
		StmtId:                    formValueStmtId,
		DynamicOutput:             formValueDynamicOutput,
		IGExtendedOutput:          formValueIgExtendedOutput,
		IncludeAnnotations:        formValueIncludeAnnotations,
		IncludeDoV:                formValueIncludeDoV,
		IncludeHeaders:            formValueIncludeHeaders,
		OutputType:                formValueOutputType,
		OutputTypes:               exporter.OUTPUT_TYPES,
		PrintPropertyTree:         formValuePropertyTree,
		PrintBinaryTree:           formValueBinaryTree,
		ActivationConditionsOnTop: formValueMoveActivationConditionsToTop,
		Width:                     WIDTH,
		Height:                    HEIGHT,
		TransactionId:             transactionID,
		RawStmtHelp:               HELP_RAW_STMT,
		CodedStmtHelpRef:          HELP_REF,
		CodedStmtHelp:             template.HTML(strings.Replace(HELP_CODED_STMT, "\n", "<br>", -1)),
		StmtIdHelp:                HELP_STMT_ID,
		ParametersHelp:            HELP_PARAMETERS,
		OutputTypeHelp:            HELP_OUTPUT_TYPE,
		ReportHelp:                HELP_REPORT}

	// Assign width for UI rendering
	if formValueCanvasWidthValue != "" {
		widthVal, err := strconv.Atoi(formValueCanvasWidthValue)
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
	if formValueCanvasHeightValue != "" {
		heightVal, err := strconv.Atoi(formValueCanvasHeightValue)
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

		// TABULAR OUTPUT PARAMETERS

		// Parameter: Statement ID
		val, suc = extractUrlParameters(r, PARAM_STATEMENT_ID)
		if suc {
			retStruct.StmtId = val
		}

		// Parameter: Dynamic output
		val, suc = extractUrlParameters(r, PARAM_DYNAMIC_SCHEMA)
		check := evaluateBooleanUrlParameters(PARAM_DYNAMIC_SCHEMA, val, suc)
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

		// Parameter: IG Logico annotations
		val, suc = extractUrlParameters(r, PARAM_LOGICO_OUTPUT)
		check = evaluateBooleanUrlParameters(PARAM_LOGICO_OUTPUT, val, suc)
		// Assign values
		if check {
			retStruct.IncludeAnnotations = CHECKBOX_CHECKED
			includeAnnotations = true
		} else {
			retStruct.IncludeAnnotations = CHECKBOX_UNCHECKED
			includeAnnotations = false
		}

		// Parameter: Header row printing
		val, suc = extractUrlParameters(r, PARAM_PRINT_HEADERS)
		check = evaluateBooleanUrlParameters(PARAM_PRINT_HEADERS, val, suc)
		// Sets default if no information is passed along in form (e.g., deactivation of headers);
		// Note that formValueIncludeHeaders will be prepopulated based on earlier form processing
		if formValueIncludeHeaders != CHECKBOX_ON {

			// Fall back and assess whether the URL contained relevant parameter
			if !suc {
				check = true
			}
			// Assign values
			if check {
				retStruct.IncludeHeaders = CHECKBOX_CHECKED
				printHeaders = true
			} else {
				retStruct.IncludeHeaders = CHECKBOX_UNCHECKED
				printHeaders = false
			}
		}

		// Parameter: Output type
		val, suc = extractUrlParameters(r, PARAM_OUTPUT_TYPE)
		if val != "" {
			// Read from parameter
			retStruct.OutputType = val
		} else {
			// Use default parameter
			retStruct.OutputType = exporter.DEFAULT_OUTPUT_TYPES
			fmt.Println("Set default output type: " + exporter.DEFAULT_OUTPUT_TYPES)
		}

		// VISUAL PARAMETERS

		// Parameter: Property tree
		val, suc = extractUrlParameters(r, PARAM_PROPERTY_TREE)
		check = evaluateBooleanUrlParameters(PARAM_PROPERTY_TREE, val, suc)
		// Manually override if not set - effectively defines default setting
		if !suc {
			check = true
		}
		// Assign values
		if check {
			retStruct.PrintPropertyTree = CHECKBOX_CHECKED
			printFlatProperties = false
		} else {
			retStruct.PrintPropertyTree = CHECKBOX_UNCHECKED
			printFlatProperties = true
		}

		// Parameter: DoV
		val, suc = extractUrlParameters(r, PARAM_DOV)
		check = evaluateBooleanUrlParameters(PARAM_DOV, val, suc)
		// Manually override if not set - effectively defines default setting
		if !suc {
			check = false
		}
		// Assign values
		if check {
			retStruct.IncludeDoV = CHECKBOX_CHECKED
			includeDoV = true
		} else {
			retStruct.IncludeDoV = CHECKBOX_UNCHECKED
			includeDoV = false
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

		// Parameter: Activation condition on top
		val, suc = extractUrlParameters(r, PARAM_ACTIVATION_CONDITION_ON_TOP)
		check = evaluateBooleanUrlParameters(PARAM_ACTIVATION_CONDITION_ON_TOP, val, suc)
		// Assign values
		if check {
			retStruct.ActivationConditionsOnTop = CHECKBOX_CHECKED
			printActivationConditionsOnTop = true
		} else {
			retStruct.ActivationConditionsOnTop = CHECKBOX_UNCHECKED
			printActivationConditionsOnTop = false
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
			fmt.Println("Google Sheets output requested")
			handleTabularOutput(w, retStruct.CodedStmt, retStruct.StmtId, retStruct, dynamicOutput, produceIGExtendedOutput, includeAnnotations, retStruct.OutputType, printHeaders)
		} else if templateName == TEMPLATE_NAME_PARSER_VISUAL {
			fmt.Println("Visual output requested")
			handleVisualOutput(w, retStruct.CodedStmt, retStruct.StmtId, retStruct, printFlatProperties, printBinaryTree, printActivationConditionsOnTop, dynamicOutput, produceIGExtendedOutput, includeAnnotations, includeDoV)
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
