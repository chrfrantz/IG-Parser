package converter

import (
	"embed"
	"html/template"
	"log"
	"net/http"
)

/*
Template reference
*/
var tmpl *template.Template

// Frontend templates for user interaction
const TEMPLATE_NAME_PARSER_SHEETS = "ig-parser-sheets.html"
const TEMPLATE_NAME_PARSER_VISUAL = "ig-parser-visualizer.html"

// Embed templates in compiled binary
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
	/*dir, err := os.Getwd()
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
		RelativePathPrefix = "./web/converter/"
	}*/
	RelativePathPrefix = ""
	// Load all templates in folder, and address specific ones during writing by name (see TEMPLATE_NAME_ constants).
	tpl, err := template.ParseFS(files, RelativePathPrefix+"templates/*")
	if err != nil {
		log.Fatal("Failed to load website templates. Error:", err)
	}
	// Assign to global variable upon successful load
	tmpl = tpl
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
