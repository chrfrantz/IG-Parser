package exporter

import "IG-Parser/core/tree"

type TabularOutputResult struct {
	Output        string
	StatementMap  []map[string]string
	HeaderSymbols []string
	HeaderNames   []string
	Error         tree.ParsingError
}
