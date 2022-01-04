package tree

import (
	"fmt"
	"log"
)

type Statement struct {

	// Regulative Statement
	Attributes                    *Node
	AttributesPropertySimple      *Node
	AttributesPropertyComplex     *Node
	Deontic                       *Node
	Aim                           *Node
	DirectObject                  *Node
	DirectObjectComplex           *Node
	DirectObjectPropertySimple    *Node
	DirectObjectPropertyComplex   *Node
	IndirectObject                *Node
	IndirectObjectComplex         *Node
	IndirectObjectPropertySimple  *Node
	IndirectObjectPropertyComplex *Node

	//Constitutive Statement
	ConstitutedEntity                     *Node
	ConstitutedEntityPropertySimple       *Node
	ConstitutedEntityPropertyComplex      *Node
	Modal                                 *Node
	ConstitutiveFunction                  *Node
	ConstitutingProperties                *Node
	ConstitutingPropertiesComplex         *Node
	ConstitutingPropertiesPropertySimple  *Node
	ConstitutingPropertiesPropertyComplex *Node

	// Shared Components
	ActivationConditionSimple  *Node
	ActivationConditionComplex *Node
	ExecutionConstraintSimple  *Node
	ExecutionConstraintComplex *Node
	OrElse                     *Node
}

/*
Returns statement as formatted string that reflects tree structure (vertical orientation, with indentation of nested elements).
*/
func (s *Statement) String() string {
	out := ""

	out = s.printComponent(out, s.Attributes, ATTRIBUTES, false, false)
	out = s.printComponent(out, s.AttributesPropertySimple, ATTRIBUTES_PROPERTY, false, false)
	out = s.printComponent(out, s.AttributesPropertyComplex, ATTRIBUTES_PROPERTY, true, false)
	out = s.printComponent(out, s.Deontic, DEONTIC, false, false)
	out = s.printComponent(out, s.Aim, AIM, false, false)
	out = s.printComponent(out, s.DirectObject, DIRECT_OBJECT, false, false)
	out = s.printComponent(out, s.DirectObjectComplex, DIRECT_OBJECT, true, false)
	out = s.printComponent(out, s.DirectObjectPropertySimple, DIRECT_OBJECT_PROPERTY, false, false)
	out = s.printComponent(out, s.DirectObjectPropertyComplex, DIRECT_OBJECT_PROPERTY, true, false)
	out = s.printComponent(out, s.IndirectObject, INDIRECT_OBJECT, false, false)
	out = s.printComponent(out, s.IndirectObjectComplex, INDIRECT_OBJECT, true, false)
	out = s.printComponent(out, s.IndirectObjectPropertySimple, INDIRECT_OBJECT_PROPERTY, false, false)
	out = s.printComponent(out, s.IndirectObjectPropertyComplex, INDIRECT_OBJECT_PROPERTY, true, false)

	out = s.printComponent(out, s.ActivationConditionSimple, ACTIVATION_CONDITION, false, false)
	out = s.printComponent(out, s.ActivationConditionComplex, ACTIVATION_CONDITION, true, false)
	out = s.printComponent(out, s.ExecutionConstraintSimple, EXECUTION_CONSTRAINT, false, false)
	out = s.printComponent(out, s.ExecutionConstraintComplex, EXECUTION_CONSTRAINT, true, false)

	out = s.printComponent(out, s.ConstitutedEntity, CONSTITUTED_ENTITY, false, false)
	out = s.printComponent(out, s.ConstitutedEntityPropertySimple, CONSTITUTED_ENTITY_PROPERTY, false, false)
	out = s.printComponent(out, s.ConstitutedEntityPropertyComplex, CONSTITUTED_ENTITY_PROPERTY, true, false)
	out = s.printComponent(out, s.Modal, MODAL, false, false)
	out = s.printComponent(out, s.ConstitutiveFunction, CONSTITUTIVE_FUNCTION, false, false)
	out = s.printComponent(out, s.ConstitutingProperties, CONSTITUTING_PROPERTIES, false, false)
	out = s.printComponent(out, s.ConstitutingPropertiesComplex, CONSTITUTING_PROPERTIES, true, false)
	out = s.printComponent(out, s.ConstitutingPropertiesPropertySimple, CONSTITUTING_PROPERTIES_PROPERTY, false, false)
	out = s.printComponent(out, s.ConstitutingPropertiesPropertyComplex, CONSTITUTING_PROPERTIES_PROPERTY, true, false)

	out = s.printComponent(out, s.OrElse, OR_ELSE, true, false)

	return out
}

/*
Appends component information for output string
Input:
- input string to append output to
- Node whose content is to be appended
- Symbol associated with component
- Indicator whether component is complex
- Indicator whether output to be constructed should be inherently human-readable (no symbols, no linebreaks, just content)

Returns string for printing
*/
func (s *Statement) printComponent(inputString string, node *Node, nodeSymbol string, complex bool, flatOutput bool) string {

	sep := ": "
	suffix := "\n"
	complexPrefix := "{\n"
	complexSuffix := "\n}"

	// If node is not nil
	if node != nil {
		// Only flat output of component (human-readable, no symbols)
		if flatOutput {
			// Generate output for node
			content := node.StringFlat()
			if len(content) > 0 {
				// If output present, append to existing output and append whitespace (to be pruned prior to final print)
				return inputString + content + " "
			} else {
				// Else simply forward input information
				return inputString
			}
		}

		// Print symbol
		inputString += nodeSymbol + sep
		// Add core content
		if complex {
			// Complex (i.e., nested) node output

			// Append complex node-specific information to the end of nested statement
			// Assumes that suffix and annotations are in string format for nodes that have nested statements
			// TODO: see whether that needs to be adjusted
			if node.Suffix != nil {
				complexSuffix += " (Suffix: " + node.Suffix.(string) + ")"
			}
			if node.Annotations != nil {
				complexSuffix += " (Annotation: " + node.Annotations.(string) + ")"
			}
			if node.PrivateNodeLinks != nil {
				complexSuffix += " (Private links: " + fmt.Sprint(node.PrivateNodeLinks) + ")"
			}
			if node.GetComponentName() != "" {
				complexSuffix += " (Component name: " + fmt.Sprint(node.GetComponentName()) + ")"
			}

			inputString += complexPrefix + node.String() + complexSuffix
		} else {
			// Simple output
			inputString += node.String()
		}
		// Append suffix
		inputString += suffix
	}
	return inputString
}

/*
Return flat string of embedded statement (human-readable output (no symbols, no linebreaks); not IG Script)
*/
func (s *Statement) StringFlat() string {
	out := ""

	out = s.printComponent(out, s.Attributes, ATTRIBUTES, false, true)
	out = s.printComponent(out, s.AttributesPropertySimple, ATTRIBUTES_PROPERTY, false, true)
	out = s.printComponent(out, s.AttributesPropertyComplex, ATTRIBUTES_PROPERTY, true, true)
	out = s.printComponent(out, s.Deontic, DEONTIC, false, true)
	out = s.printComponent(out, s.Aim, AIM, false, true)
	out = s.printComponent(out, s.DirectObject, DIRECT_OBJECT, false, true)
	out = s.printComponent(out, s.DirectObjectComplex, DIRECT_OBJECT, true, true)
	out = s.printComponent(out, s.DirectObjectPropertySimple, DIRECT_OBJECT_PROPERTY, false, true)
	out = s.printComponent(out, s.DirectObjectPropertyComplex, DIRECT_OBJECT_PROPERTY, true, true)
	out = s.printComponent(out, s.IndirectObject, INDIRECT_OBJECT, false, true)
	out = s.printComponent(out, s.IndirectObjectComplex, INDIRECT_OBJECT, true, true)
	out = s.printComponent(out, s.IndirectObjectPropertySimple, INDIRECT_OBJECT_PROPERTY, false, true)
	out = s.printComponent(out, s.IndirectObjectPropertyComplex, INDIRECT_OBJECT_PROPERTY, true, true)

	out = s.printComponent(out, s.ActivationConditionSimple, ACTIVATION_CONDITION, false, true)
	out = s.printComponent(out, s.ActivationConditionComplex, ACTIVATION_CONDITION, true, true)
	out = s.printComponent(out, s.ExecutionConstraintSimple, EXECUTION_CONSTRAINT, false, true)
	out = s.printComponent(out, s.ExecutionConstraintComplex, EXECUTION_CONSTRAINT, true, true)

	out = s.printComponent(out, s.ConstitutedEntity, CONSTITUTED_ENTITY, false, true)
	out = s.printComponent(out, s.ConstitutedEntityPropertySimple, CONSTITUTED_ENTITY_PROPERTY, false, true)
	out = s.printComponent(out, s.ConstitutedEntityPropertyComplex, CONSTITUTED_ENTITY_PROPERTY, true, true)
	out = s.printComponent(out, s.Modal, MODAL, false, true)
	out = s.printComponent(out, s.ConstitutiveFunction, CONSTITUTIVE_FUNCTION, false, true)
	out = s.printComponent(out, s.ConstitutingProperties, CONSTITUTING_PROPERTIES, false, true)
	out = s.printComponent(out, s.ConstitutingPropertiesComplex, CONSTITUTING_PROPERTIES, true, true)
	out = s.printComponent(out, s.ConstitutingPropertiesPropertySimple, CONSTITUTING_PROPERTIES_PROPERTY, false, true)
	out = s.printComponent(out, s.ConstitutingPropertiesPropertyComplex, CONSTITUTING_PROPERTIES_PROPERTY, true, true)

	out = s.printComponent(out, s.OrElse, OR_ELSE, true, true)

	// Cut last element if any output string exists (since it will contain an appended whitespace)
	if out != "" {
		out = out[:len(out)-1]
	}
	return out
}

/*
Stringifies institutional statement
*/
func (s *Statement) Stringify() string {
	log.Fatal("Stringify() is not yet implemented.")
	return ""
}

/*
Generates map of arrays containing pointers to leaf nodes in each component.
Key is an incrementing index, and value is an array of the corresponding nodes.
It further returns an array containing the component keys alongside the number of leaf nodes per component,
in order to reconstruct the linkage between the index in the first return value and the components they relate to.

Example: The first return may include two ATTRIBUTES component trees separated by synthetic AND connections (sAND)
based on different logical combination within the attributes component that are not genuine logical relationships (i.e.,
not signaled using [AND], [OR], or [XOR], but inferred during parsing based on the occurrence of multiple such combinations
within an Attributes component expression (e.g., A((Sellers [AND] Buyers) from (Northern [OR] Southern) states)).
Internally, this would be represented as ((Sellers [AND] Buyers] [bAND] (Northern [OR] Southern))', and returned as separate
trees with index 0 (Sellers [AND] Buyers) and 1 (Northern [OR] Southern).
The second return indicates the fact that the first two entries in the first return type instance are of type ATTRIBUTES by holding
an entry '"ATTRIBUTES": 2', etc.

The parameter aggregateImplicitLinkages specifies whether implicit linkages (based on bAND) are actually aggregated, or
returned as separate node trees.

*/
func (s *Statement) GenerateLeafArrays(aggregateImplicitLinkages bool) ([][]*Node, map[string]int) {
	return s.generateLeafArrays(aggregateImplicitLinkages, 0)
}

/*
Generates map of arrays containing pointers to leaf nodes in each component.
Key is an incrementing index, and value is an array of the corresponding nodes.
It further returns an array containing the component keys alongside the number of leaf nodes per component,
in order to reconstruct the linkage between the index in the first return value and the components they relate to.

Note: This variant only returns nodes that have a non-nil suffix.

Example: The first return may include two ATTRIBUTES component trees separated by synthetic AND connections (sAND)
based on different logical combination within the attributes component that are not genuine logical relationships (i.e.,
not signaled using [AND], [OR], or [XOR], but inferred during parsing based on the occurrence of multiple such combinations
within an Attributes component expression (e.g., A((Sellers [AND] Buyers) from (Northern [OR] Southern) states)).
Internally, this would be represented as ((Sellers [AND] Buyers] [sAND] (Northern [OR] Southern))', and returned as separate
trees with index 0 (Sellers [AND] Buyers) and 1 (Northern [OR] Southern).
The second return indicates the fact that the first two entries in the first return type instance are of type ATTRIBUTES by holding
an entry '"ATTRIBUTES": 2', etc.

The parameter aggregateImplicitLinkages indicates whether implicitly linked trees of nodes should be returned as a single
tree, or as separate trees.
The parameter level indicates whether all nodes should be returned, or only ones that contain suffix information.

*/
func (s *Statement) GenerateLeafArraysSuffixOnly(aggregateImplicitLinkages bool) ([][]*Node, map[string]int) {
	return s.generateLeafArrays(aggregateImplicitLinkages, 1)
}

/*
Generates map of arrays containing pointers to leaf nodes in each component.
Key is an incrementing index, and value is an array of the corresponding nodes.
It further returns an array containing the component keys alongside the number of leaf nodes per component,
in order to reconstruct the linkage between the index in the first return value and the components they relate to.

Input: level indicates selection of nodes considered in aggregation (0 --> all nodes, 1 --> nodes with non-nil suffix only)

Example: The first return may include two ATTRIBUTES component trees separated by synthetic AND connections (bAND)
based on different logical combination within the attributes component that are not genuine logical relationships (i.e.,
not signaled using [AND], [OR], or [XOR], but inferred during parsing based on the occurrence of multiple such combinations
within an Attributes component expression (e.g., A((Sellers [AND] Buyers) from (Northern [OR] Southern) states)).
Internally, this would be represented as ((Sellers [AND] Buyers] [bAND] (Northern [OR] Southern))', and returned as separate
trees with index 0 (Sellers [AND] Buyers) and 1 (Northern [OR] Southern).
The second return indicates the fact that the first two entries in the first return type instance are of type ATTRIBUTES by holding
an entry '"ATTRIBUTES": 2', etc.

The parameter aggregateImplicitLinkages indicates whether implicitly linked trees of nodes should be returned as a single
tree, or as separate trees.
The parameter level indicates whether all nodes should be returned, or only ones that contain suffix information.

*/
func (s *Statement) generateLeafArrays(aggregateImplicitLinkages bool, level int) ([][]*Node, map[string]int) {

	// Map holding reference from component type (e.g., ATTRIBUTES) to number of entries (relevant for reconstruction)
	referenceMap := map[string]int{}

	// Counter for overall number of entries
	nodesMap := make([][]*Node, 0)

	// Regulative components
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Attributes, ATTRIBUTES, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.AttributesPropertySimple, ATTRIBUTES_PROPERTY, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.AttributesPropertyComplex, ATTRIBUTES_PROPERTY_REFERENCE, true, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Deontic, DEONTIC, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Aim, AIM, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObject, DIRECT_OBJECT, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObjectComplex, DIRECT_OBJECT_REFERENCE, true, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObjectPropertySimple, DIRECT_OBJECT_PROPERTY, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObjectPropertyComplex, DIRECT_OBJECT_PROPERTY_REFERENCE, true, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObject, INDIRECT_OBJECT, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObjectComplex, INDIRECT_OBJECT_REFERENCE, true, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObjectPropertySimple, INDIRECT_OBJECT_PROPERTY, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObjectPropertyComplex, INDIRECT_OBJECT_PROPERTY_REFERENCE, true, aggregateImplicitLinkages, level)

	// Context
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ActivationConditionSimple, ACTIVATION_CONDITION, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ActivationConditionComplex, ACTIVATION_CONDITION_REFERENCE, true, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ExecutionConstraintSimple, EXECUTION_CONSTRAINT, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ExecutionConstraintComplex, EXECUTION_CONSTRAINT_REFERENCE, true, aggregateImplicitLinkages, level)

	// Constitutive components
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutedEntity, CONSTITUTED_ENTITY, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutedEntityPropertySimple, CONSTITUTED_ENTITY_PROPERTY, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutedEntityPropertyComplex, CONSTITUTED_ENTITY_PROPERTY_REFERENCE, true, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Modal, MODAL, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutiveFunction, CONSTITUTIVE_FUNCTION, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingProperties, CONSTITUTING_PROPERTIES, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingPropertiesComplex, CONSTITUTING_PROPERTIES_REFERENCE, true, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingPropertiesPropertySimple, CONSTITUTING_PROPERTIES_PROPERTY, false, aggregateImplicitLinkages, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingPropertiesPropertyComplex, CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE, true, aggregateImplicitLinkages, level)

	// Shared components
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.OrElse, OR_ELSE, true, aggregateImplicitLinkages, level)

	return nodesMap, referenceMap
}

/*
Generates a leaf array for a given component under consideration of node as being of simple or complex nature.
Appends to existing structure if provided (i.e., not nil) to allow for iterative invocation.
For a version that allows for iterative invocation, consider #getComponentLeafArray.
For returning only leaves that contain suffix information, consider #getComponentLeafArrayWithSuffix.

Input:
- Reference to component node for which leaf elements are to be extracted
- Component symbol associated with component
- Indicator whether element embedded in node is complex (i.e., nested statement)
- The parameter aggregateImplicitLinkages indicates whether implicitly linked trees of nodes should be returned as a single tree, or as separate trees.
- The parameter level indicates whether all nodes should be returned, or only ones that contain suffix information.

- Indicator whether all leaf nodes should be returned, or only one satisfying particular conditions
  (0 --> all nodes, 1 --> only ones with non-empty suffix).

Returns:
- Node map of nodes associated with components
- Reference map counting number of components
*/
/*func GetSingleComponentLeafArray(componentNode *Node, componentSymbol string, complex bool, aggregateImplicitLinkages bool, level int) ([][]*Node, map[string]int) {

	// Map holding reference from component type (e.g., ATTRIBUTES) to number of entries (relevant for reconstruction)
	referenceMap := map[string]int{}

	// Counter for overall number of entries
	nodesMap := make([][]*Node, 0)

	return getComponentLeafArray(nodesMap, referenceMap, componentNode, componentSymbol, complex, aggregateImplicitLinkages, level)
}*/

/*
Generates a leaf array for a given component under consideration of node as being of simple or complex nature.
Appends to existing structure if provided (i.e., not nil) to allow for iterative invocation.
For returning only leaves that contain suffix information consider #getComponentLeafArrayWithSuffix.
Input:
- maps of nodes potentially including existing nodes for other components. Will be created internally if nil
  (to allow iterative invocation).
- reference map that indexes the number of nodes associated with a specific component (to retain association).
  Will be created internally if nil (to allow iterative invocation).
- Reference to component node for which leaf elements are to be extracted
- Component symbol associated with component
- Indicator whether element embedded in node is complex (i.e., nested statement)
- Indicator whether all leaf nodes should be returned, or only one satisfying particular conditions
  (0 --> all nodes, 1 --> only ones with non-empty suffix).

Returns:
- Node map of nodes associated with components
- Reference map counting number of components
*/
func getComponentLeafArray(nodesMap [][]*Node, referenceMap map[string]int, componentNode *Node, componentSymbol string, complex bool, aggregateImplicitLinkages bool, level int) ([][]*Node, map[string]int) {

	if componentNode == nil {
		Println("No component node found - returning unmodified node and reference map ...")
		return nodesMap, referenceMap
	}

	// Initialize data structures if nil
	if nodesMap == nil {
		nodesMap = make([][]*Node, 1)
	}

	if referenceMap == nil {
		referenceMap = make(map[string]int, 1)
	}

	// Check for complex content
	if complex {
		// Embed nested statement in node structure, before adding to node map
		nodesMap = append(nodesMap, []*Node{componentNode})

		// since statements can be combined, they are returned as a single element
		referenceMap[componentSymbol] = 1
	} else {
		// Counter for number of elements in given component
		i := 0
		// Add array of leaf nodes attached to general array
		for _, v := range componentNode.GetLeafNodes(aggregateImplicitLinkages) {
			nodesMap = append(nodesMap, v)
			i++
		}
		// Add number of nodes referring to a particular component
		referenceMap[componentSymbol] = i
	}

	// Return modified or generated structures
	return nodesMap, referenceMap
}

/*
Returns the property node corresponding to the current component. If the component does
not possess a corresponding property, or the node itself is nil, the function returns an empty array.
Otherwise, the properties node(s) is/are returned. Variably allows for return of primitive nodes only,
or also complex ones (i.e., nested statements). Where multiple primitive nodes exist, those are
returned as combinations.
*/
func (s *Statement) GetPropertyComponent(n *Node, complex bool) []*Node {
	out := make([]*Node, 0)

	// Check whether node is actually not nil
	if n == nil {
		return out
	}

	// Explore mapping of components and properties
	switch n.GetComponentName() {
	case ATTRIBUTES:
		if s.AttributesPropertySimple != nil {
			out = append(out, s.AttributesPropertySimple)
		}
		if complex && s.AttributesPropertyComplex != nil {
			out = append(out, s.AttributesPropertyComplex)
		}
	case DIRECT_OBJECT:
		if s.DirectObjectPropertySimple != nil {
			out = append(out, s.DirectObjectPropertySimple)
		}
		if complex && s.DirectObjectPropertyComplex != nil {
			out = append(out, s.DirectObjectPropertyComplex)
		}
	case INDIRECT_OBJECT:
		if s.IndirectObjectPropertySimple != nil {
			out = append(out, s.IndirectObjectPropertySimple)
		}
		if complex && s.IndirectObjectPropertyComplex != nil {
			out = append(out, s.IndirectObjectPropertyComplex)
		}
	case CONSTITUTED_ENTITY:
		if s.ConstitutedEntityPropertySimple != nil {
			out = append(out, s.ConstitutedEntityPropertySimple)
		}
		if complex && s.ConstitutedEntityPropertyComplex != nil {
			out = append(out, s.ConstitutedEntityPropertyComplex)
		}
	case CONSTITUTING_PROPERTIES:
		if s.ConstitutingPropertiesPropertySimple != nil {
			out = append(out, s.ConstitutingPropertiesPropertySimple)
		}
		if complex && s.ConstitutingPropertiesPropertyComplex != nil {
			out = append(out, s.ConstitutingPropertiesPropertyComplex)
		}
	}
	return out
}
