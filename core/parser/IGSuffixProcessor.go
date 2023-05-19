package parser

import (
	"IG-Parser/core/tree"
	"reflect"
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
		Println("Source node:", v)
		val := v
		if val.GetSuffix() != "" {
			rawSuffix := val.GetSuffix()
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
				Println("Target node to test:", val2)
				if val2.GetSuffix() != "" {
					rawTargetSuffix := val2.GetSuffix()
					Println("Found target suffix", rawTargetSuffix)
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
Processes reorganization of statements by iterating over statement to identify private linkages amongst
components and properties identified by suffix. Converts parsed elements into private nodes and attaches
those to associated components. It operates directly on provided statement.

Parameters include the statement to operate on, as well as the indication whether to identify linkages
for simple (for basic statements) or complex properties (properties that themselves consist of nested statements/structures).
*/
func ProcessPrivateComponentLinkages(s *tree.Statement, complex bool) {

	Println("Statement before reviewing linkages: ", s)

	// Find all leaves that have suffix
	leafArrays, _ := s.GenerateLeafArraysSuffixOnly(true)

	if len(leafArrays) == 0 {
		Println("No leaf entries found, hence no suffix linkages.")
		return
	}

	// Pairs of matched component-private property linkages
	type Pair struct {
		Src *tree.Node
		Tgt *tree.Node
	}

	// Array to keep track of identified linkages for post-processing (i.e., removal from tree)
	identifiedLinkages := []Pair{}

	// Identify links starting from top-level components
	for _, v := range leafArrays {

		// Initialize linked leaves structure
		linkedLeaves := map[*tree.Node][]*tree.Node{}

		for elem := range v {

			// Extract source component for which associated candidate private components are sought
			sourceComponentElement := v[elem]

			switch sourceComponentElement.GetComponentName() {
			case tree.ATTRIBUTES:
				if complex {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.AttributesPropertyComplex)
				} else {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.AttributesPropertySimple)
				}
			case tree.AIM:
				if complex {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ExecutionConstraintComplex)
				} else {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ExecutionConstraintSimple)
				}
			case tree.DIRECT_OBJECT:
				if complex {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.DirectObjectPropertyComplex)
				} else {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.DirectObjectPropertySimple)
				}
			case tree.INDIRECT_OBJECT:
				if complex {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.IndirectObjectPropertyComplex)
				} else {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.IndirectObjectPropertySimple)
				}
			case tree.CONSTITUTED_ENTITY:
				if complex {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ConstitutedEntityPropertyComplex)
				} else {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ConstitutedEntityPropertySimple)
				}
			case tree.CONSTITUTING_PROPERTIES:
				if complex {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ConstitutingPropertiesPropertyComplex)
				} else {
					linkedLeaves = FindNodesLinkedViaSuffix(sourceComponentElement, s.ConstitutingPropertiesPropertySimple)
				}
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
						// Keep track of linked nodes for later removal -- see below
						identifiedLinkages = append(identifiedLinkages, Pair{Src: sourceComponentElement, Tgt: tgtComp})

						Println("Component", srcComp.GetComponentName(), ": Added private link for node", tgtComp)
					}
				}
			}
		}
	}
	Println("Private node linkages to remove from main tree:", identifiedLinkages)

	// Post process linkages for removal of private nodes from property tree, and potential removal of source node if empty
	for _, pair := range identifiedLinkages {

		Println("-> Processing removal of identified private node from statement tree structure. Node: " + pair.Tgt.String())

		// Copy potentially inherited component name into target node directly (which will be lost after node removal)
		pair.Tgt.ComponentType = pair.Tgt.GetComponentName()

		// Remove private node from original tree structure
		rt, err := tree.RemoveNodeFromTree(pair.Tgt)
		if err.ErrorCode != tree.TREE_NO_ERROR {
			// Do not deal with error, since false will always refer to need for node removal from statement
			Println("Response when attempting to remove private node from tree (likely final element in component):", err.Error())
		}

		// If return value is false (e.g., node is disconnected), this may imply that the node is the last element, confirm (if statement is its own root) and reset node linkage from statement perspective
		if !rt {
			// Identify corresponding element from statement and remove from statement tree
			Println("Statement-side reset of linkages to removed node. Resetting of corresponding property if element", pair.Tgt, "is last (private) element")
			switch pair.Src.GetComponentName() {
			case tree.ATTRIBUTES:
				if reflect.DeepEqual(pair.Tgt, s.AttributesPropertySimple) {
					Println("Reset A,p ...")
					s.AttributesPropertySimple = nil
				}
				if reflect.DeepEqual(pair.Tgt, s.AttributesPropertyComplex) {
					Println("Reset A,p (complex) ...")
					s.AttributesPropertyComplex = nil
				}
			case tree.AIM:
				if reflect.DeepEqual(pair.Tgt, s.ExecutionConstraintSimple) {
					Println("Reset Cex ...")
					s.ExecutionConstraintSimple = nil
				}
				if reflect.DeepEqual(pair.Tgt, s.ExecutionConstraintComplex) {
					Println("Reset Cex (complex) ...")
					s.ExecutionConstraintComplex = nil
				}
			case tree.DIRECT_OBJECT:
				if reflect.DeepEqual(pair.Tgt, s.DirectObjectPropertySimple) {
					Println("Reset Bdir,p ...")
					s.DirectObjectPropertySimple = nil
				}
				if reflect.DeepEqual(pair.Tgt, s.DirectObjectPropertyComplex) {
					Println("Reset Bdir,p (complex) ...")
					s.DirectObjectPropertyComplex = nil
				}
			case tree.INDIRECT_OBJECT:
				if reflect.DeepEqual(pair.Tgt, s.IndirectObjectPropertySimple) {
					Println("Reset Bind,p ...")
					s.IndirectObjectPropertySimple = nil
				}
				if reflect.DeepEqual(pair.Tgt, s.IndirectObjectPropertyComplex) {
					Println("Reset Bind,p (complex) ...")
					s.IndirectObjectPropertyComplex = nil
				}
			case tree.CONSTITUTED_ENTITY:
				if reflect.DeepEqual(pair.Tgt, s.ConstitutedEntityPropertySimple) {
					Println("Reset E,p ...")
					s.ConstitutedEntityPropertySimple = nil
				}
				if reflect.DeepEqual(pair.Tgt, s.ConstitutedEntityPropertyComplex) {
					Println("Reset E,p (complex) ...")
					s.ConstitutedEntityPropertyComplex = nil
				}
			case tree.CONSTITUTING_PROPERTIES:
				if reflect.DeepEqual(pair.Tgt, s.ConstitutingPropertiesPropertySimple) {
					Println("Reset P,p ...")
					s.ConstitutingPropertiesPropertySimple = nil
				}
				if reflect.DeepEqual(pair.Tgt, s.ConstitutingPropertiesPropertyComplex) {
					Println("Reset P,p (complex) ...")
					s.ConstitutingPropertiesPropertyComplex = nil
				}
			default:
				Println("Node deletion from tree: Could not find match for component name", pair.Src.GetComponentName())
			}
		}
	}
	// Now linkages should have been inferred, and deleted from property tree (as well as deletion of empty source nodes)
}
