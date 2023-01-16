package converter

import "html/template"

/*
Struct for interacting with template via handler
*/
type ReturnStruct struct {
	// Indicates whether operation was successful
	Success bool
	// Indicates whether an error has occurred
	Error bool
	// Message shown to user
	Message string
	// Original unparsed statement
	RawStmt string
	// IG-Script annotated statement
	CodedStmt string
	// Statement ID
	StmtId string
	// Dynamic output indicator
	DynamicOutput string
	// IG Extended output indicator
	IGExtendedOutput string
	// Annotation inclusion indicator
	IncludeAnnotations string
	// Degree of Variability inclusion indicator
	IncludeDoV string
	// Include headers in output
	IncludeHeaders string
	// Output type indicator (e.g., Google Sheets, CSV)
	OutputType string
	// Output types (to populate UI)
	OutputTypes []string
	// Property tree printing indicator
	PrintPropertyTree string
	// Binary tree printing indicator (as opposed to tree aggregation based on logical operator by component)
	PrintBinaryTree string
	// Binary indicator whether activation conditions should be output on top of visual tree, or their regular position
	ActivationConditionsOnTop string
	// Generated output to be rendered (e.g., tabular, visual)
	Output string
	// Width of output canvas
	Width int
	// Height of output canvas
	Height int
	// Transaction ID
	TransactionId string
	// Help message for raw statement
	RawStmtHelp string
	// Help text indicating reference to help page
	CodedStmtHelpRef string
	// Help message for coded statement - needs to be provided as parseable HTML for templating.
	CodedStmtHelp template.HTML
	// Help message for statement ID
	StmtIdHelp string
	// Help message for parameters
	ParametersHelp string
	// Help message for output format
	OutputTypeHelp string
	// Help message for report tooltip
	ReportHelp string
}
