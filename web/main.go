package main

import (
	"IG-Parser/web/converter"
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

func main() {

	// Initializes templating and determines correct relative path for templates and CSS
	converter.Init()

	http.HandleFunc("/", converter.ConverterHandler)
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir(converter.RelativePathPrefix + "css"))))

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

	// Check for logging path (if logging is enabled)
	if converter.Logging == true {
		logPath := os.Getenv(ENV_VAR_LOGGING_PATH)
		if logPath != "" {
			converter.LoggingPath = logPath
			log.Println("Found logging path: " + logPath)
		}
	}

	addr := ":" + port
	log.Printf("Listening on %s ...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))

}
