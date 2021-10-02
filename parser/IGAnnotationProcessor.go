package parser

import (
	"IG-Parser/tree"
	"fmt"
	"strings"
)

/*
Identifies links via elements in suffices (if multiple, comma-separated ones, it extracts only first)
between leaf nodes in source and target nodes.
Returns map of linked nodes, with key being the source node, and the value being an array of matched target nodes.
 */
func FindNodesLinkedViaSuffix(sourceTree *tree.Node, targetTree *tree.Node) map[*tree.Node][]*tree.Node {

	// Result structure with source node as key, and suffix-matched target nodes as value
	linkMap := make(map[*tree.Node][]*tree.Node)

	// Retrieve source arrays
	sourceArrays := tree.Flatten(sourceTree.GetLeafNodes())
	fmt.Println("Source arrays: ", sourceArrays)
	if len(sourceArrays) == 0 {
		fmt.Println("Could not find leaf nodes in source tree.")
		return linkMap
	}

	// Retrieve target arrays
	targetArrays := tree.Flatten(targetTree.GetLeafNodes())
	fmt.Println("Target arrays: ", targetArrays)
	if len(targetArrays) == 0 {
		fmt.Println("Could not find leaf nodes in target tree.")
		return linkMap
	}

	// Iterate through source components
	for _, v := range sourceArrays {
		//fmt.Println("Val:", v)
		val := v
		if val.Suffix != nil && len(val.Suffix.(string)) > 0 {
			rawSuffix := val.Suffix.(string)
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

			// Now check target side to see if there is matching suffix
			for _, v2 := range targetArrays {
				val2 := v2
				//fmt.Println("Target val:", val2)
				if val2.Suffix != nil && len(val2.Suffix.(string)) > 0 {
					rawTargetSuffix := val2.Suffix.(string)
					//fmt.Println("Found target suffix", rawTargetSuffix)
					// Assign full suffix by default
					targetElem := rawTargetSuffix
					// Extract first element from candidate target suffix
					idx2 := strings.Index(rawTargetSuffix, SUFFIX_DELIMITER)
					if idx2 != -1 {
						targetElem = rawTargetSuffix[:idx2]
					}
					// Test for match between source and target; if the same, store association
					if sourceElem == targetElem {
						fmt.Println("Found suffix match on", sourceElem, "for components (Source:", val.Entry,
							", Name:", val.GetComponentName(), " and Target:", val2.Entry, ", Name:", val2.GetComponentName(), ")")
						valArr := []*tree.Node{}
						// Check if existing entry exists
						if linkMap[val] != nil {
							// and extract
							valArr = linkMap[val]
						}
						// Append to existing array if entry exists
						valArr = append(valArr, val2)
						linkMap[val] = valArr
					}
				}
			}
		}
	}
	return linkMap
}

/*
Processes reorganization of statements to convert parsed elements as private elements based on suffix-based linkages.
Operates directly on provided statement.
 */
func ProcessPrivateComponentLinkages(s *tree.Statement) {

	fmt.Println("Statement before reviewing linkages: ", s)

	// Find all leaves that have suffix
	leafArrays, _ := s.GenerateLeafArraysSuffixOnly()

	if len(leafArrays) == 0 {
		fmt.Println("No leaf entries found, hence no suffix linkages.")
		return
	}

	// Identify links starting from top-level components
	for _, v := range leafArrays {

		// Initialize linked leaves structure
		linkedLeaves := map[*tree.Node][]*tree.Node{}

		for elem := range v {

			sourceComponentElement := v[elem]
			//fmt.Println("Source:", sourceComponentElement)
			//fmt.Println("Source component:", sourceComponentElement.GetComponentName())

			switch sourceComponentElement.GetComponentName() {
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
				fmt.Println("Could not find match for component name", sourceComponentElement.GetComponentName())
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
	}
	//fmt.Println("Statement after reviewing linkages: ", s)
}