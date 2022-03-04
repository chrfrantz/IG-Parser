package converter

import (
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

// Frontend templates for user interaction
const TEMPLATE_NAME_PARSER_SHEETS = "ig-parser-sheets.html"
const TEMPLATE_NAME_PARSER_VISUAL = "ig-parser-visualizer.html"

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
	// Load all templates in folder, and address specific ones during writing by name (see TEMPLATE_NAME_ constants).
	tmpl = template.Must(template.ParseGlob(RelativePathPrefix + "templates/*"))
}

/*
Handler for Google Sheets.
*/
func ConverterHandlerSheets(w http.ResponseWriter, r *http.Request) {
	converterHandler(w, r, TEMPLATE_NAME_PARSER_SHEETS)
}

/*
Handler for visualization.
*/
func ConverterHandlerVisual(w http.ResponseWriter, r *http.Request) {
	converterHandler(w, r, TEMPLATE_NAME_PARSER_VISUAL)
}

/*
Serves favicon.
*/
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received favicon request")
	http.ServeFile(w, r, "web/css/favicon.ico")
}

/*
Serves D3 library.
*/
func D3Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received D3 library request")
	http.ServeFile(w, r, "web/libraries/d3.v7.min.js")
}
