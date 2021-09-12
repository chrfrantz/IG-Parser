package parser

import (
	"IG-Parser/tree"
	"fmt"
	"strings"
)

/*
Extracts links between nodes of given statement.
Returns map of links between source and target nodes.
 */
func ExtractLinkBetweenProperties(s tree.Statement) map[*tree.Node]*tree.Node {

	// Find linked nodes from Properties to Attributes
	linkMap := CombineMaps(FindNodesLinkedViaSuffix(s.AttributesPropertySimple, s.Attributes),
		FindNodesLinkedViaSuffix(s.DirectObjectPropertySimple, s.DirectObject), false)


	return linkMap
}


/*
Identifies links via elements in suffices (if multiple, comma-separated ones, it extracts only first)
between leaf nodes in source and target nodes.
Returns map of linked nodes.
 */
func FindNodesLinkedViaSuffix(sourceTree *tree.Node, targetTree *tree.Node) map[*tree.Node]*tree.Node {

	linkMap := make(map[*tree.Node]*tree.Node)

	// Store origin of linkages between properties and components
	sourceArrays := sourceTree.GetLeafNodes()
	for k, v := range sourceArrays {
		fmt.Println(k, "Val:", v)
		if v[0].Suffix != nil && len(v[0].Suffix.(string)) > 0 {
			rawSuffix := v[0].Suffix.(string)
			// Assign full suffix by default
			sourceElem := rawSuffix
			// Extract first element from suffix
			idx := strings.Index(rawSuffix, SUFFIX_DELIMITER)
			// If delimiter is found, extract element, else take whole content
			if idx != -1 {
				sourceElem = rawSuffix[:idx]
			}

			if len(rawSuffix) != len(sourceElem) {
				fmt.Println("Complete processing of suffix not yet supported. Remaining elements:", rawSuffix[idx:])
			}
			for _, v2 := range targetTree.GetLeafNodes() {
				if v2[0].Suffix != nil && len(v2[0].Suffix.(string)) > 0 {
					rawTargetSuffix := v2[0].Suffix.(string)
					// Assign full suffix by default
					targetElem := rawTargetSuffix
					// Extract first element from candidate target suffix
					idx2 := strings.Index(rawTargetSuffix, SUFFIX_DELIMITER)
					if idx2 != -1 {
						targetElem = rawTargetSuffix[:idx2]
					}
					// Test for match between source and target
					if sourceElem == targetElem {
						fmt.Println("Found match between Attributes Properties and Attributes components")
						linkMap[v[0]] = v2[0]
					}
				}
			}
		}
	}
	return linkMap
}