package converter

import (
	"IG-Parser/web/converter/shared"
	"html/template"
	"log"
	"net/http"
	"strings"
)

/*
This file contains the top-level handler to serve built-in
help functionality.
*/

/*
Handler serving help information.
*/
func HelpHandler(w http.ResponseWriter, r *http.Request) {

	// Converts help text to HTML and populates generic struct
	data := shared.ReturnStruct{CodedStmtHelp: template.HTML(strings.Replace(shared.HELP_CODED_STMT, "\n", "<br>", -1))}

	// Populate help template with text
	err := tmpl.ExecuteTemplate(w, TEMPLATE_NAME_HELP, data)
	if err != nil {
		log.Println("Error generating error response for template processing:", err.Error())
		http.Error(w, "Could not process request.", http.StatusInternalServerError)
	}

}
