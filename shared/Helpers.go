package shared

import "strings"

/*
Replaces/escapes selected symbols in as far as relevant for export (e.g., quotation marks).
*/
func EscapeSymbolsForExport(rawValue string) string {
	// Replace quotation marks with single quotes
	return strings.ReplaceAll(rawValue, "\"", "'")
}
