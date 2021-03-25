package exporter

import "IG-Parser/tree"

type LogicalOperatorLinkage struct {
	// Index in component structure of statement
	Index int
	// Source component pointer
	Component *tree.Node
	// References to all ids the source component itself occurs in
	OwnStatements []int
	// Structure holding target node as key, and array of referenced entries (row id)
	LinkedStatements map[*tree.Node][]int
	// Structure holding reference to target component and associated operator (to be used in conjunction with LinkedStatements)
	LinkedComponentOperator map[*tree.Node][]string
	// Keeping operators in order of precedence (first = higher precedence)
	OperatorPrecedence []string
}

type OperatorLinkages struct {
	Operator []string
	Refs []int
}
