package main

import (
	"IG-Parser/app"
	"IG-Parser/tree"
	"IG-Parser/web/helper"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var ANNOTATED_STATEMENT = "(National Organic Program's Program Manager), Cex(on behalf of the Secretary), D(may) I(inspect and), I(sustain (review [AND] (refresh [AND] drink))) Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) Cex(for compliance with the (Act or [XOR] regulations in this part))"
var STATEMENT_ID = "650"

// Indicates whether logging occurs
var logging = true

var tmpl *template.Template

type ReturnStruct struct{
	// Indicates whether operation was successful
	Success bool;
	// Indicates whether an error has occurred
	Error bool;
	// Message shown to user
	Message string;
	// Original unparsed statement
	RawStmt string;
	// IG-Script annotated statement
	CodedStmt string;
	// Statement ID
	StmtId string;
	// Generated tabular output
	TabularOutput string
	// Transaction ID
	TransactionId string;
}

/*
Dummy function in case logging is not activated
 */
var terminateOutput = func() {}

func init() {
	tmpl = template.Must(template.ParseFiles("./web/templates/IG-Parser-Form.html"))
}

func parserHandler(w http.ResponseWriter, r *http.Request) {
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

	fmt.Println(retStruct)

	// Check for input statement first
	if codedStmt == "" {
		retStruct.Success = false
		retStruct.Error = true
		retStruct.Message = app.ERROR_INPUT_NO_STATEMENT
		tmpl.Execute(w, retStruct)
		return
	} else {
		// Check for statement ID
		id, err := strconv.Atoi(stmtId)
		if err != nil {
			retStruct.Success = false
			retStruct.Error = true
			retStruct.Message = app.ERROR_INPUT_STATEMENT_ID
			tmpl.Execute(w, retStruct)
			return
		}
		// Only then parse input
		if logging {
			tID, filename := helper.GenerateUniqueIdAndFilename()
			// Assign transaction ID
			retStruct.TransactionId = tID
			terminateOutput = helper.SaveOutput(filename)
			fmt.Println("TRANSACTION ID: " + retStruct.TransactionId)
		}
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

func main() {
	http.HandleFunc("/", parserHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))

	// Make service Heroku-compatible
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Listening on %s ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
