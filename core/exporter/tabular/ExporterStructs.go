package tabular

import "IG-Parser/core/tree"

/*
The TabularOutputResult contains the results of the tabular parsing process,
while maintaining metadata (e.g., headers) relevant for post-processing or
output generation. It is primarily used from tabular.TabularOutputGenerator.go,
but is the primary structure on which any prospective output format should be based.
*/
type TabularOutputResult struct {
	Output        string
	StatementMap  []map[string]string
	HeaderSymbols []string
	HeaderNames   []string
	Error         tree.ParsingError
}
