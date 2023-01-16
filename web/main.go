package main

import (
	"IG-Parser/web/converter"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

/*
Environment variables (port, logging activation)
*/
const ENV_VAR_PORT = "IG_PARSER_PORT"
const ENV_VAR_LOGGING = "IG_PARSER_LOGGING"
const ENV_VAR_LOGGING_PATH = "IG_PARSER_LOGGING_PATH"

const DEFAULT_PORT = "8080"

const VISUAL_PATH = "visual"

// Embed files in compiled binary

//go:embed css/default.css css/favicon.ico
var cssFiles embed.FS

//go:embed libraries/d3.v7.min.js libraries/ace/ace.js
var libraryFiles embed.FS

func main() {

	// Initializes templating and determines correct relative path for templates and CSS
	converter.Init()

	// Conventional tabular output handler
	http.HandleFunc("/", converter.ConverterHandlerSheets)
	// Visual tree output handler
	http.HandleFunc("/"+VISUAL_PATH+"/", converter.ConverterHandlerVisual)
	// D3 (served for visual output)
	http.Handle("/libraries/", http.FileServer(http.FS(libraryFiles)))
	// CSS folder mapping (for CSS and favicon)
	http.Handle("/css/", http.FileServer(http.FS(cssFiles)))

	// Check for custom port
	port := os.Getenv(ENV_VAR_PORT)
	if port == "" {
		port = DEFAULT_PORT
	}

	// Check for logging specification (default activated)
	logEnv := os.Getenv(ENV_VAR_LOGGING)
	if logEnv == "" || strings.ToLower(logEnv) == "true" {
		converter.Logging = true
	} else {
		converter.Logging = false
	}

	//converter.Logging = false

	// Check for logging path (if logging is enabled)
	if converter.Logging == true {
		logPath := os.Getenv(ENV_VAR_LOGGING_PATH)
		if logPath != "" {
			converter.LoggingPath = logPath
			log.Println("Found logging path: " + logPath)
		}
	}

	// Manual override for local testing
	converter.LoggingPath = "./logs"

	// Redirect stdout
	//temp := os.Stdout
	//os.Stdout = nil

	addr := ":" + port
	log.Printf("Listening on %s ...\n", addr)
	log.Println("Logging: " + fmt.Sprint(converter.Logging))
	log.Println("Logging path: " + fmt.Sprint(converter.LoggingPath))
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("Web service stopped. Error:", err)
	}

}
