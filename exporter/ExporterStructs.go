package exporter

import (
	"IG-Parser/tree"
	"fmt"
	"log"
)

/*
Returns operators and statement IDs for target node IDs associated with target node in given input
 */
func GetLogicalOperatorAndStatementRefs(sourceNode *tree.Node, targetNode *tree.Node, componentRefs map[*tree.Node][]int) ([]string, []int) {
	// Find operator links from source to target
	res, ops, err := tree.FindLogicalLinkage(sourceNode, targetNode)
	if err.ErrorCode != tree.TREE_NO_ERROR {
		errorMsg := fmt.Sprint("Error when parsing retrieving operator linkages: ", err.ErrorMessage)
		log.Println(errorMsg)
		fmt.Errorf("&v", errorMsg)
		return nil, nil
	}
	if res {
		// Extract IDs for target component
		ids := componentRefs[targetNode]
		// Return both operators and target IDs
		return ops, ids
	}
	return nil, nil
}