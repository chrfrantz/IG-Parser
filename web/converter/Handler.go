package converter

import (
	"IG-Parser/app"
	"IG-Parser/tree"
	"IG-Parser/web/helper"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

/*
Template reference
 */
var tmpl *template.Template
/*
Dummy function in case logging is not activated
*/
var terminateOutput = func() {}
/*
Indicates whether logging occurs
 */
var Logging = true
/*
Indicates folder to log to
 */
var LoggingPath = ""

/*
Init needs to be called from main to instantiate templates.
 */
func Init() {
	tmpl = template.Must(template.ParseFiles("./web/templates/IG-Parser-Form.html"))
}

func ConverterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Just show empty form
		tmpl.Execute(w, ReturnStruct{Success: false, Message: "", CodedStmt: ANNOTATED_STATEMENT, StmtId: STATEMENT_ID})
		return
	}
	message := ""
	transactionID := ""
	rawStmt := r.FormValue("rawStatement")
	codedStmt := r.FormValue("annotatedStatement")
	stmtId := r.FormValue("stmtId")
	retStruct := ReturnStruct{Success: false, Error: false, Message: message, RawStmt: rawStmt, CodedStmt: codedStmt, StmtId: stmtId, TransactionId: transactionID}

	// Initiatialize request-specific logfile first
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
			fmt.Println("Error when redirecting output: " + err.Error())
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
		tmpl.Execute(w, retStruct)

		// Final comment in log
		fmt.Println("Error: No input to parse.")
		// Ensure logging is terminated
		terminateOutput()

		return
	} else {
		// Check for statement ID
		id, err := strconv.Atoi(stmtId)
		if err != nil {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = app.ERROR_INPUT_STATEMENT_ID
			tmpl.Execute(w, retStruct)

			// Final comment in log
			fmt.Println("Error: " + fmt.Sprint(err))
			// Ensure logging is terminated
			terminateOutput()

			return
		}
		// Convert input
		output, err2 := app.ConvertIGScriptToGoogleSheets(codedStmt, id, "")
		if err2.ErrorCode != tree.PARSING_NO_ERROR {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.CodedStmt = codedStmt
			switch err2.ErrorCode {
			case tree.PARSING_ERROR_EMPTY_LEAF:
				retStruct.Message = app.ERROR_INPUT_NO_STATEMENT
			default:
				retStruct.Message = "Parsing error (" + err2.ErrorCode + "): " + err2.ErrorMessage
			}
			tmpl.Execute(w, retStruct)

			// Final comment in log
			fmt.Println("Error: " + fmt.Sprint(err2))
			// Ensure logging is terminated
			terminateOutput()

			return
		}
		// Return success if parsing was successful
		retStruct.Success = true
		retStruct.CodedStmt = codedStmt
		retStruct.TabularOutput = output
		tmpl.Execute(w, retStruct)

		// Final comment in log
		fmt.Println("Success")
		// Ensure logging is terminated
		terminateOutput()

		return
	}
}
