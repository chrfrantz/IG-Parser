package main

import (
	"IG-Parser/web/converter"
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

func main() {

	// Initializes templating and determines correct relative path for templates and CSS
	converter.Init()

	// Conventional tabular output handler
	http.HandleFunc("/", converter.ConverterHandlerSheets)
	// Visual tree output handler
	http.HandleFunc("/"+VISUAL_PATH+"/", converter.ConverterHandlerVisual)
	// Favicon (served for regular and visual output)
	http.HandleFunc("/favicon.ico", converter.FaviconHandler)
	http.HandleFunc("/"+VISUAL_PATH+"/favicon.ico", converter.FaviconHandler)
	// D3 (served for visual output)
	http.HandleFunc("/"+VISUAL_PATH+"/d3.v7.min.js", converter.D3Handler)
	// CSS folder mapping
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(converter.RelativePathPrefix+"css"))))

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
	log.Fatal(http.ListenAndServe(addr, nil))

}
