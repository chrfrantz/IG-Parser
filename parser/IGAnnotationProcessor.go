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
/*func ExtractLinkBetweenProperties(s tree.Statement) map[*tree.Node]*tree.Node {

	// Find linked nodes from Properties to Attributes
	linkMap := CombineMaps(FindNodesLinkedViaSuffix(s.AttributesPropertySimple, s.Attributes),
		FindNodesLinkedViaSuffix(s.DirectObjectPropertySimple, s.DirectObject), false)

	return linkMap
}*/

/*func ConvertSuffixedNodesToPrivateNodes(statement tree.Statement) {
	gen
}*/

/*
Identifies links via elements in suffices (if multiple, comma-separated ones, it extracts only first)
between leaf nodes in source and target nodes.
Returns map of linked nodes, with key being the source node, and the value being an array of matched target nodes.
 */
func FindNodesLinkedViaSuffix(sourceTree *tree.Node, targetTree *tree.Node) map[*tree.Node][]*tree.Node {

	linkMap := make(map[*tree.Node][]*tree.Node)

	// Store origin of linkages between properties and components
	sourceArrays := sourceTree.GetLeafNodes()

	fmt.Println("Source arrays: ", sourceArrays)
	if len(sourceArrays) == 0 {
		fmt.Println("Could not find leaf nodes in source tree.")
		return linkMap
	}

	targetArrays := targetTree.GetLeafNodes()
	fmt.Println("Target arrays: ", targetArrays)
	if len(targetArrays) == 0 {
		fmt.Println("Could not find leaf nodes in target tree.")
		return linkMap
	}


	for _, v := range sourceArrays {
		fmt.Println("Val:", v)
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
				fmt.Println("Complete processing of suffix with more than one element not yet supported. Remaining elements:", rawSuffix[idx:])
			}

			fmt.Println("Element count: ", targetArrays)

			// Now check target side to see if there is matching suffix
			for _, v2 := range targetArrays {
				fmt.Println("Target val:", v2)
				if v2[0].Suffix != nil && len(v2[0].Suffix.(string)) > 0 {
					rawTargetSuffix := v2[0].Suffix.(string)
					fmt.Println("Found target suffix", rawTargetSuffix)
					// Assign full suffix by default
					targetElem := rawTargetSuffix
					// Extract first element from candidate target suffix
					idx2 := strings.Index(rawTargetSuffix, SUFFIX_DELIMITER)
					if idx2 != -1 {
						targetElem = rawTargetSuffix[:idx2]
					}
					// Test for match between source and target; if the same, store association
					if sourceElem == targetElem {
						fmt.Println("Found suffix match on", sourceElem, "for components (Source:", sourceTree.GetComponentName(), ", Target:", targetTree.GetComponentName(), ")")
						valArr := []*tree.Node{}
						// Check if existing entry exists
						if linkMap[v[0]] != nil {
							// and extract
							valArr = linkMap[v[0]]
						}
						// Append to existing array if entry exists
						valArr = append(valArr, v2[0])
						linkMap[v[0]] = valArr
					}
				}
			}
		}
	}
	return linkMap
}

func ProcessPrivateComponentLinkages(s *tree.Statement) {

	fmt.Println("Statement before reviewing linkages: ", s)

	// Find all leaves that have suffix
	leafArrays, _ := s.GenerateLeafArraysSuffixOnly()

	fmt.Println(leafArrays)

	if len(leafArrays) == 0 {
		fmt.Println("No leaf entries found, hence no suffix linkages.")
		return
	}

	fmt.Println("Array: ", leafArrays)
	// Identify links starting from top-level components
	for _, v := range leafArrays {

		// Initialize linked leaves structure
		linkedLeaves := map[*tree.Node][]*tree.Node{}

		sourceComponentElement := v[0]
		fmt.Println("Source:", sourceComponentElement)
		fmt.Println("Source component:", sourceComponentElement.GetComponentName())

		switch (sourceComponentElement.GetComponentName()) {
		case tree.ATTRIBUTES:
			linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.AttributesPropertySimple)
		case tree.DIRECT_OBJECT:
			linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.DirectObjectPropertySimple)
		case tree.INDIRECT_OBJECT:
			linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.IndirectObjectPropertySimple)
		case tree.CONSTITUTED_ENTITY:
			linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ConstitutedEntityPropertySimple)
		case tree.CONSTITUTING_PROPERTIES:
			linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ConstitutingPropertiesPropertySimple)
		default:
			fmt.Println("Could not find match for component name.")
		}
		if len(linkedLeaves) > 0 {
			fmt.Println("Found following links for", sourceComponentElement.GetComponentName(), ":", linkedLeaves)

			// Draw direct linkage between source and target component
			for srcComp, tgtCompArr := range linkedLeaves {
				for _, tgtComp := range tgtCompArr {
					linkComps := []*tree.Node{}
					// Retrieve potentially existing node links
					if srcComp.PrivateNodeLinks != nil {
						linkComps = srcComp.PrivateNodeLinks
					}
					// Add target component
					linkComps = append(linkComps, tgtComp)
					// Attach to node
					srcComp.PrivateNodeLinks = linkComps
					// Remove private node from original tree structure
					tree.RemoveNodeFromTree(tgtComp)
					fmt.Println("Component", srcComp.GetComponentName(), ": Added private link for node", tgtComp)
				}
			}
		}
	}
	fmt.Println("Statement after reviewing linkages: ", s)
	// Return original statement
	//return s
}