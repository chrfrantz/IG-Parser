package converter

import (
	"IG-Parser/app"
	"IG-Parser/exporter"
	"IG-Parser/tree"
	"IG-Parser/web/helper"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

/*
Template reference
 */
var tmpl *template.Template

/*
Dummy function in case logging is not activated
*/
var terminateOutput = func(string) error {
	return nil
}

/*
Indicates whether logging occurs
 */
var Logging = true

/*
Indicates folder to log to
 */
var LoggingPath = ""

/*
Relative path prefix for all web resources (templates, CSS files)
 */
var RelativePathPrefix = ""

/*
Success suffix
 */
const SUCCESS_SUFFIX = ".success"

/*
Error suffix
 */
const ERROR_SUFFIX = ".error"

/*
Init needs to be called from main to instantiate templates.
 */
func Init() {
	dir, err := os.Getwd()
	if err != nil {
		// Sensible to terminate in this case
		log.Fatal(err)
	}

	fmt.Println("Working directory: " + dir)
	// If in docker container
	if dir == "/" {
		// relative to web folder
		RelativePathPrefix = "../"
	} else {
		// else started from repository root
		RelativePathPrefix = "./web/"
	}
	tmpl = template.Must(template.ParseFiles(RelativePathPrefix + "templates/IG-Parser-Form.html"))
}

func ConverterHandler(w http.ResponseWriter, r *http.Request) {

	// Prepopulate response
	message := ""
	transactionID := ""
	rawStmt := r.FormValue("rawStatement")
	codedStmt := r.FormValue("annotatedStatement")
	stmtId := r.FormValue("stmtId")
	dynChk := r.FormValue("dynamicOutput")
	inclAnnotations := r.FormValue("annotations")

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

	retStruct := ReturnStruct{
		Success: false,
		Error: false,
		Message: message,
		RawStmt: rawStmt,
		CodedStmt: codedStmt,
		StmtId: stmtId,
		DynamicOutput: dynChk,
		IncludeAnnotations: inclAnnotations,
		TransactionId: transactionID,
		RawStmtHelp: HELP_RAW_STMT,
		CodedStmtHelp: HELP_CODED_STMT,
		StmtIdHelp: HELP_STMT_ID,
		ReportHelp: HELP_REPORT}

	if r.Method != http.MethodPost {
		// Just show empty form with prepopulated elements
		retStruct.CodedStmt = ANNOTATED_STATEMENT
		retStruct.StmtId = STATEMENT_ID

		err := tmpl.Execute(w, retStruct)
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
		retStruct.Message = app.ERROR_INPUT_NO_STATEMENT
		retStruct.DynamicOutput = dynChk
		retStruct.IncludeAnnotations = inclAnnotations
		err := tmpl.Execute(w, retStruct)
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
		// Define whether output is dynamic
		fmt.Println("Setting dynamic output: ", dynamicOutput)
		exporter.SetDynamicOutput(dynamicOutput)
		// Define whether annotations are included
		fmt.Println("Setting annotations: ", includeAnnotations)
		exporter.SetIncludeAnnotations(includeAnnotations)
		// Convert input
		output, err2 := app.ConvertIGScriptToGoogleSheets(codedStmt, stmtId, "")
		if err2.ErrorCode != tree.PARSING_NO_ERROR {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.CodedStmt = codedStmt
			retStruct.DynamicOutput = dynChk
			retStruct.IncludeAnnotations = inclAnnotations
			switch err2.ErrorCode {
			case tree.PARSING_ERROR_EMPTY_LEAF:
				retStruct.Message = app.ERROR_INPUT_NO_STATEMENT
			default:
				retStruct.Message = "Parsing error (" + err2.ErrorCode + "): " + err2.ErrorMessage
			}
			err3 := tmpl.Execute(w, retStruct)
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
		retStruct.TabularOutput = output
		retStruct.DynamicOutput = dynChk
		retStruct.IncludeAnnotations = inclAnnotations
		err := tmpl.Execute(w, retStruct)
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
}
