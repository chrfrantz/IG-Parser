package exporter

import "IG-Parser/tree"

/*
Returns operators and statement IDs for target node IDs associated with target node in given input
 */
func GetLogicalOperatorAndStatementRefs(sourceNode *tree.Node, targetNode *tree.Node, componentRefs map[*tree.Node][]int) ([]string, []int) {
	// Find link from source to target
	res, ops := tree.FindLogicalLinkage(sourceNode, targetNode, nil)
	if res {
		// Extract IDs for target component
		ids := componentRefs[targetNode]
		// Return both operators and target IDs
		return ops, ids
	}
	return nil, nil
}