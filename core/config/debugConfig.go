package config

/*
This file holds various switches that allows for selective suppression of focal output,
both for development and production use (see default settings).
These switches are primarily used in core packages; to date only the DEBUG_FRONTEND switch
is reserved for frontend application (e.g., web frontend, third-party applications).
Note: If logging is enabled in the web application (web/converter/Handler.go), the
transaction-related output (i.e., output related to input processing, not page invocation),
will be preserved in corresponding log files (named after transaction ID).
*/

// Override debug switch that, if set to true, overrides all settings below (default: false)
var DEBUG_ALL = false

// Debug related to the parsing of input statements (default: false)
var DEBUG_STATEMENT_PARSING = false

// Debug related to tree operations specifically (default: false)
var DEBUG_TREE_OPERATIONS = false

// Debug related to the generation process (preceding final output processing)
var DEBUG_OUTPUT_GENERATION = false

// Debug related to the final output preparation (e.g., CSV outfile)
var DEBUG_FINAL_OUTPUT = false

// Debug information related to the frontend (e.g., web frontend). Can be used by any third-party application.
var DEBUG_FRONTEND = true
