package converter

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

/*
Template reference
*/
var tmpl *template.Template

// Frontend templates for user interaction
const TEMPLATE_NAME_PARSER_TABULAR = "ig-parser-tabular.html"
const TEMPLATE_NAME_PARSER_VISUAL = "ig-parser-visualizer.html"

// Help template
const TEMPLATE_NAME_HELP = "ig-parser-user-guide.html"

// Embed templates in compiled binary
//
//go:embed templates/*
var files embed.FS

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
	// Load all templates in folder, and address specific ones during processing by name (see TEMPLATE_NAME_ constants).
	tpl, err := template.ParseFS(files, "templates/*")
	if err != nil {
		log.Fatal("Failed to load website templates. Error:", err)
	}
	// Assign to global variable upon successful load
	tmpl = tpl
}

/*
Handler for tabular output.
*/
func ConverterHandlerTabular(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Invoked TABULAR output handler")
	converterHandler(w, r, TEMPLATE_NAME_PARSER_TABULAR)
}

/*
Handler for visual output.
*/
func ConverterHandlerVisual(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Invoked VISUAL output handler")
	converterHandler(w, r, TEMPLATE_NAME_PARSER_VISUAL)
}
