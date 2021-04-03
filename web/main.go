package main

import (
	"IG-Parser/app"
	"IG-Parser/tree"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

var tmpl *template.Template

type ReturnStruct struct{
	// Indicates whether operation was successful
	Success bool;
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
}

func init() {
	tmpl = template.Must(template.ParseFiles("./web/templates/IG-Parser-Form.html"))
}

func parserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		// Just show empty form
		tmpl.Execute(w, ReturnStruct{Success: true, Message: ""})
		return
	}
	success := false
	message := ""
	rawStmt := r.FormValue("rawStatement")
	codedStmt := r.FormValue("annotatedStatement")
	stmtId := r.FormValue("stmtId")
	retStruct := ReturnStruct{Success: success, Message: message, RawStmt: rawStmt, CodedStmt: codedStmt, StmtId: stmtId}

	fmt.Println(retStruct)

	id, err := strconv.Atoi(stmtId)
	if err != nil {
		retStruct.Success = false
		retStruct.Message = err.Error()
		tmpl.Execute(w, retStruct)
	}
	if codedStmt == "" {
		retStruct.Success = false
		retStruct.Message = "No parseable input data"
		tmpl.Execute(w, retStruct)
	} else {
		output, err := app.ConvertIGScriptToGoogleSheets(codedStmt, id, "")
		if err.ErrorCode != tree.PARSING_NO_ERROR {
			retStruct.Success = false
			retStruct.Message = err.ErrorMessage
			tmpl.Execute(w, retStruct)
			return
		}
		retStruct.Success = true
		retStruct.TabularOutput = output
		tmpl.Execute(w, retStruct)
	}
}

func main() {
	http.HandleFunc("/", parserHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./web/css"))))

	// Make it Heroku-compatible
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("Listening on %s ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
