package parser

import (
	"IG-Parser/tree"
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
	sourceArrays := tree.Flatten(sourceTree.GetLeafNodes(true))
	Println("Source arrays: ", sourceArrays)
	if len(sourceArrays) == 0 {
		Println("Could not find leaf nodes in source tree.")
		return linkMap
	}

	// Retrieve target arrays
	targetArrays := tree.Flatten(targetTree.GetLeafNodes(true))
	Println("Target arrays: ", targetArrays)
	if len(targetArrays) == 0 {
		Println("Could not find leaf nodes in target tree.")
		return linkMap
	}

	// Iterate through source components
	for _, v := range sourceArrays {
		//Println("Val:", v)
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
				Println("Complete processing of suffix with more than one element not yet supported. Remaining elements:", rawSuffix[idx:])
			}

			// Now check target side to see if there is matching suffix
			for _, v2 := range targetArrays {
				val2 := v2
				//Println("Target val:", val2)
				if val2.Suffix != nil && len(val2.Suffix.(string)) > 0 {
					rawTargetSuffix := val2.Suffix.(string)
					//Println("Found target suffix", rawTargetSuffix)
					// Assign full suffix by default
					targetElem := rawTargetSuffix
					// Extract first element from candidate target suffix
					idx2 := strings.Index(rawTargetSuffix, SUFFIX_DELIMITER)
					if idx2 != -1 {
						targetElem = rawTargetSuffix[:idx2]
					}
					// Test for match between source and target; if the same, store association
					if sourceElem == targetElem {
						Println("Found suffix match on", sourceElem, "for components (Source:", val.Entry,
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

	Println("Statement before reviewing linkages: ", s)

	// Find all leaves that have suffix
	leafArrays, _ := s.GenerateLeafArraysSuffixOnly(true)

	if len(leafArrays) == 0 {
		Println("No leaf entries found, hence no suffix linkages.")
		return
	}

	// Identify links starting from top-level components
	for _, v := range leafArrays {

		// Initialize linked leaves structure
		linkedLeaves := map[*tree.Node][]*tree.Node{}

		for elem := range v {

			sourceComponentElement := v[elem]
			//Println("Source:", sourceComponentElement)
			//Println("Source component:", sourceComponentElement.GetComponentName())

			switch sourceComponentElement.GetComponentName() {
			case tree.ATTRIBUTES:
				linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.AttributesPropertySimple)
			case tree.AIM:
				linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ExecutionConstraintSimple)
			case tree.DIRECT_OBJECT:
				linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.DirectObjectPropertySimple)
			case tree.INDIRECT_OBJECT:
				linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.IndirectObjectPropertySimple)
			case tree.CONSTITUTED_ENTITY:
				linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ConstitutedEntityPropertySimple)
			case tree.CONSTITUTING_PROPERTIES:
				linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ConstitutingPropertiesPropertySimple)
			default:
				Println("Could not find match for component name", sourceComponentElement.GetComponentName())
			}
			if len(linkedLeaves) > 0 {
				Println("Found following links for", sourceComponentElement.GetComponentName(), ":", linkedLeaves)

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
						rt, err := tree.RemoveNodeFromTree(tgtComp)
						if err.ErrorCode != tree.TREE_NO_ERROR {
							// Do not deal with error, since false will always refer to need for node removal from statement
						}
						// If return value is false, this implies that the remaining node is the last element - and to be removed from statement
						if !rt {
							// Identify corresponding element from statement and remove from statement tree
							Println("Element", srcComp, "will be removed from parent tree (last element)")
							switch sourceComponentElement.GetComponentName() {
							case tree.ATTRIBUTES:
								s.AttributesPropertySimple = nil
							case tree.AIM:
								s.ExecutionConstraintSimple = nil
							case tree.DIRECT_OBJECT:
								s.DirectObjectPropertySimple = nil
							case tree.INDIRECT_OBJECT:
								s.IndirectObjectPropertySimple = nil
							case tree.CONSTITUTED_ENTITY:
								s.ConstitutedEntityPropertySimple = nil
							case tree.CONSTITUTING_PROPERTIES:
								s.ConstitutingPropertiesPropertySimple = nil
							default:
								Println("Node deletion from tree: Could not find match for component name", sourceComponentElement.GetComponentName())
							}
						}
						Println("Component", srcComp.GetComponentName(), ": Added private link for node", tgtComp)
					}
				}
			}
		}
	}
}