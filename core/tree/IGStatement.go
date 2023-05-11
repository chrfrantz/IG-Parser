package tree

import (
	"IG-Parser/core/shared"
	"fmt"
	"log"
	"strings"
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
	return s.string(0)
}

/*
Variant of String() with indentation level passed into output (for pretty printing of nested structures)
*/
func (s *Statement) string(indentationLevel int) string {
	out := ""

	out = s.printComponent(out, indentationLevel, s.Attributes, ATTRIBUTES, false, false, false)
	out = s.printComponent(out, indentationLevel, s.AttributesPropertySimple, ATTRIBUTES_PROPERTY, false, false, false)
	out = s.printComponent(out, indentationLevel, s.AttributesPropertyComplex, ATTRIBUTES_PROPERTY, true, false, false)
	out = s.printComponent(out, indentationLevel, s.Deontic, DEONTIC, false, false, false)
	out = s.printComponent(out, indentationLevel, s.Aim, AIM, false, false, false)
	out = s.printComponent(out, indentationLevel, s.DirectObject, DIRECT_OBJECT, false, false, false)
	out = s.printComponent(out, indentationLevel, s.DirectObjectComplex, DIRECT_OBJECT, true, false, false)
	out = s.printComponent(out, indentationLevel, s.DirectObjectPropertySimple, DIRECT_OBJECT_PROPERTY, false, false, false)
	out = s.printComponent(out, indentationLevel, s.DirectObjectPropertyComplex, DIRECT_OBJECT_PROPERTY, true, false, false)
	out = s.printComponent(out, indentationLevel, s.IndirectObject, INDIRECT_OBJECT, false, false, false)
	out = s.printComponent(out, indentationLevel, s.IndirectObjectComplex, INDIRECT_OBJECT, true, false, false)
	out = s.printComponent(out, indentationLevel, s.IndirectObjectPropertySimple, INDIRECT_OBJECT_PROPERTY, false, false, false)
	out = s.printComponent(out, indentationLevel, s.IndirectObjectPropertyComplex, INDIRECT_OBJECT_PROPERTY, true, false, false)

	out = s.printComponent(out, indentationLevel, s.ActivationConditionSimple, ACTIVATION_CONDITION, false, false, false)
	out = s.printComponent(out, indentationLevel, s.ActivationConditionComplex, ACTIVATION_CONDITION, true, false, false)
	out = s.printComponent(out, indentationLevel, s.ExecutionConstraintSimple, EXECUTION_CONSTRAINT, false, false, false)
	out = s.printComponent(out, indentationLevel, s.ExecutionConstraintComplex, EXECUTION_CONSTRAINT, true, false, false)

	out = s.printComponent(out, indentationLevel, s.ConstitutedEntity, CONSTITUTED_ENTITY, false, false, false)
	out = s.printComponent(out, indentationLevel, s.ConstitutedEntityPropertySimple, CONSTITUTED_ENTITY_PROPERTY, false, false, false)
	out = s.printComponent(out, indentationLevel, s.ConstitutedEntityPropertyComplex, CONSTITUTED_ENTITY_PROPERTY, true, false, false)
	out = s.printComponent(out, indentationLevel, s.Modal, MODAL, false, false, false)
	out = s.printComponent(out, indentationLevel, s.ConstitutiveFunction, CONSTITUTIVE_FUNCTION, false, false, false)
	out = s.printComponent(out, indentationLevel, s.ConstitutingProperties, CONSTITUTING_PROPERTIES, false, false, false)
	out = s.printComponent(out, indentationLevel, s.ConstitutingPropertiesComplex, CONSTITUTING_PROPERTIES, true, false, false)
	out = s.printComponent(out, indentationLevel, s.ConstitutingPropertiesPropertySimple, CONSTITUTING_PROPERTIES_PROPERTY, false, false, false)
	out = s.printComponent(out, indentationLevel, s.ConstitutingPropertiesPropertyComplex, CONSTITUTING_PROPERTIES_PROPERTY, true, false, false)

	out = s.printComponent(out, indentationLevel, s.OrElse, OR_ELSE, true, false, false)

	return out
}

/*
Appends component information for output string
Input:
- input string to append output to
- indentation level (indicating nesting level for pretty printing)
- Node whose content is to be appended
- Symbol associated with component
- Indicator whether component is complex
- Indicator whether output to be constructed should be inherently human-readable (no symbols, no linebreaks, just content)
- Indicator whether component symbol should be included in the output (as opposed to merely content of component) - intended to finetune output

Returns string for printing
*/
func (s *Statement) printComponent(inputString string, indentationLevel int, node *Node, nodeSymbol string, complex bool, flatOutput bool, includeComponentSymbol bool) string {

	b := strings.Builder{}
	b.WriteString(inputString)

	sep := ": "
	suffix := "\n"

	indentString := ""
	i := 0
	for i < indentationLevel {
		indentString += MinimumIndentPrefix
		i++
	}

	// If node is not nil
	if node != nil {
		// Prepend indentationLevel
		b.WriteString(indentString)
		// Only flat output of component (human-readable, no symbols)
		if flatOutput {
			// Generate output for node
			content := node.StringFlat()
			if len(content) > 0 {
				if includeComponentSymbol {
					// Includes component symbol in output
					b.WriteString(nodeSymbol)
					b.WriteString("(")
					b.WriteString(content)
					b.WriteString(")")
					b.WriteString(" ")
					return b.String()
				} else {
					// If output present, append to existing output and append whitespace (to be pruned prior to final print)
					b.WriteString(content)
					b.WriteString(" ")
					return b.String()
				}
			} else {
				// Else simply forward input information
				return inputString
			}
		}

		// Print symbol
		b.WriteString(nodeSymbol)
		b.WriteString(sep)
		// Add core content
		if complex {
			// Complex (i.e., nested) node output

			complexPrefix := "{\n"
			// Indent closing brace at deeper level due to indentation of content
			complexSuffix := "\n" + indentString + MinimumIndentPrefix + "}"
			complexSuffixBuilder := strings.Builder{}
			complexSuffixBuilder.WriteString(complexSuffix)

			// Append complex node-specific information to the end of nested statement
			// Assumes that suffix and annotations are in string format for nodes that have nested statements
			// TODO: see whether that needs to be adjusted
			if node.GetSuffix() != "" {
				complexSuffixBuilder.WriteString(" (Suffix: ")
				complexSuffixBuilder.WriteString(node.GetSuffix())
				complexSuffixBuilder.WriteString(")")
			}
			if node.Annotations != nil {
				complexSuffixBuilder.WriteString(" (Annotation: ")
				complexSuffixBuilder.WriteString(node.Annotations.(string))
				complexSuffixBuilder.WriteString(")")
			}
			if node.PrivateNodeLinks != nil {
				complexSuffixBuilder.WriteString(" (Private links: ")
				complexSuffixBuilder.WriteString(fmt.Sprint(node.PrivateNodeLinks))
				complexSuffixBuilder.WriteString(")")
			}
			if node.GetComponentName() != "" {
				complexSuffixBuilder.WriteString(" (Component name: ")
				complexSuffixBuilder.WriteString(fmt.Sprint(node.GetComponentName()))
				complexSuffixBuilder.WriteString(")")
			}

			// Write out pre-generated information ...
			// Write prefix
			b.WriteString(complexPrefix)
			// Write node component content (may include nested structure)
			b.WriteString(node.string(indentationLevel + 1))
			// Write complex node suffix
			b.WriteString(complexSuffixBuilder.String())
		} else {
			// Simple output
			b.WriteString(node.string(indentationLevel))
		}
		// Append suffix
		b.WriteString(suffix)
	}
	return b.String()
}

/*
Return flat string of embedded statement (human-readable output (no linebreaks); not full IG Script, but potentially component symbols for nested structures)
- Indicator whether component symbols should be included for nested statements (e.g., statements in properties).
*/
func (s Statement) StringFlatStatement(includeComponentSymbol bool) string {
	out := ""

	out = s.printComponent(out, 0, s.Attributes, ATTRIBUTES, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.AttributesPropertySimple, ATTRIBUTES_PROPERTY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.AttributesPropertyComplex, ATTRIBUTES_PROPERTY, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.Deontic, DEONTIC, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.Aim, AIM, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.DirectObject, DIRECT_OBJECT, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.DirectObjectComplex, DIRECT_OBJECT, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.DirectObjectPropertySimple, DIRECT_OBJECT_PROPERTY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.DirectObjectPropertyComplex, DIRECT_OBJECT_PROPERTY, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.IndirectObject, INDIRECT_OBJECT, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.IndirectObjectComplex, INDIRECT_OBJECT, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.IndirectObjectPropertySimple, INDIRECT_OBJECT_PROPERTY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.IndirectObjectPropertyComplex, INDIRECT_OBJECT_PROPERTY, true, true, includeComponentSymbol)

	out = s.printComponent(out, 0, s.ActivationConditionSimple, ACTIVATION_CONDITION, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ActivationConditionComplex, ACTIVATION_CONDITION, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ExecutionConstraintSimple, EXECUTION_CONSTRAINT, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ExecutionConstraintComplex, EXECUTION_CONSTRAINT, true, true, includeComponentSymbol)

	out = s.printComponent(out, 0, s.ConstitutedEntity, CONSTITUTED_ENTITY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutedEntityPropertySimple, CONSTITUTED_ENTITY_PROPERTY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutedEntityPropertyComplex, CONSTITUTED_ENTITY_PROPERTY, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.Modal, MODAL, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutiveFunction, CONSTITUTIVE_FUNCTION, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutingProperties, CONSTITUTING_PROPERTIES, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutingPropertiesComplex, CONSTITUTING_PROPERTIES, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutingPropertiesPropertySimple, CONSTITUTING_PROPERTIES_PROPERTY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutingPropertiesPropertyComplex, CONSTITUTING_PROPERTIES_PROPERTY, true, true, includeComponentSymbol)

	out = s.printComponent(out, 0, s.OrElse, OR_ELSE, true, true, includeComponentSymbol)

	// Cut last element if any output string exists (since it will contain an appended whitespace)
	if out != "" {
		out = out[:len(out)-1]
	}
	return out
}

/*
Return flat string of embedded statement (human-readable output (no linebreaks); no full IG Script, but potentially component symbols for nested structures)
- Indicator whether component symbols should be included for nested statements (e.g., statements in properties).
*/
func (s *Statement) StringFlat(includeComponentSymbol bool) string {
	out := ""

	out = s.printComponent(out, 0, s.Attributes, ATTRIBUTES, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.AttributesPropertySimple, ATTRIBUTES_PROPERTY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.AttributesPropertyComplex, ATTRIBUTES_PROPERTY, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.Deontic, DEONTIC, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.Aim, AIM, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.DirectObject, DIRECT_OBJECT, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.DirectObjectComplex, DIRECT_OBJECT, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.DirectObjectPropertySimple, DIRECT_OBJECT_PROPERTY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.DirectObjectPropertyComplex, DIRECT_OBJECT_PROPERTY, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.IndirectObject, INDIRECT_OBJECT, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.IndirectObjectComplex, INDIRECT_OBJECT, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.IndirectObjectPropertySimple, INDIRECT_OBJECT_PROPERTY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.IndirectObjectPropertyComplex, INDIRECT_OBJECT_PROPERTY, true, true, includeComponentSymbol)

	out = s.printComponent(out, 0, s.ActivationConditionSimple, ACTIVATION_CONDITION, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ActivationConditionComplex, ACTIVATION_CONDITION, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ExecutionConstraintSimple, EXECUTION_CONSTRAINT, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ExecutionConstraintComplex, EXECUTION_CONSTRAINT, true, true, includeComponentSymbol)

	out = s.printComponent(out, 0, s.ConstitutedEntity, CONSTITUTED_ENTITY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutedEntityPropertySimple, CONSTITUTED_ENTITY_PROPERTY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutedEntityPropertyComplex, CONSTITUTED_ENTITY_PROPERTY, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.Modal, MODAL, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutiveFunction, CONSTITUTIVE_FUNCTION, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutingProperties, CONSTITUTING_PROPERTIES, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutingPropertiesComplex, CONSTITUTING_PROPERTIES, true, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutingPropertiesPropertySimple, CONSTITUTING_PROPERTIES_PROPERTY, false, true, includeComponentSymbol)
	out = s.printComponent(out, 0, s.ConstitutingPropertiesPropertyComplex, CONSTITUTING_PROPERTIES_PROPERTY, true, true, includeComponentSymbol)

	out = s.printComponent(out, 0, s.OrElse, OR_ELSE, true, true, includeComponentSymbol)

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
			if level == 1 {
				Println("Leaf nodes to consider for suffix:", v)
				// Iterate through nodes to detect suffix
				for _, v2 := range v {
					Println("Node to check for suffix:", v2)
					// Check for presence of suffix before adding individually
					if v2.GetSuffix() != "" {
						Println("Found suffix in node:", v2)
						nodesMap = append(nodesMap, []*Node{v2})
						i++
					}
				}
			} else {
				// In all other cases, simple add all leaf nodes (no checking for suffix-only nodes)
				nodesMap = append(nodesMap, v)
				i++
			}
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

/*
Copies components from a given statement into statement on which the function is called. Checks whether target
component is empty before copying.
*/
func CopyComponentsFromStatement(stmtToCopyTo *Statement, stmtToCopyFrom *Statement) *Statement {

	// Regulative
	stmtToCopyTo.Attributes = copyComponentValue(stmtToCopyTo.Attributes, stmtToCopyFrom.Attributes)
	stmtToCopyTo.AttributesPropertySimple = copyComponentValue(stmtToCopyTo.AttributesPropertySimple, stmtToCopyFrom.AttributesPropertySimple)
	stmtToCopyTo.AttributesPropertyComplex = copyComponentValue(stmtToCopyTo.AttributesPropertyComplex, stmtToCopyFrom.AttributesPropertyComplex)
	stmtToCopyTo.Deontic = copyComponentValue(stmtToCopyTo.Deontic, stmtToCopyFrom.Deontic)
	stmtToCopyTo.Aim = copyComponentValue(stmtToCopyTo.Aim, stmtToCopyFrom.Aim)
	stmtToCopyTo.DirectObject = copyComponentValue(stmtToCopyTo.DirectObject, stmtToCopyFrom.DirectObject)
	stmtToCopyTo.DirectObjectComplex = copyComponentValue(stmtToCopyTo.DirectObjectComplex, stmtToCopyFrom.DirectObjectComplex)
	stmtToCopyTo.DirectObjectPropertySimple = copyComponentValue(stmtToCopyTo.DirectObjectPropertySimple, stmtToCopyFrom.DirectObjectPropertySimple)
	stmtToCopyTo.DirectObjectPropertyComplex = copyComponentValue(stmtToCopyTo.DirectObjectPropertyComplex, stmtToCopyFrom.DirectObjectPropertyComplex)
	stmtToCopyTo.IndirectObject = copyComponentValue(stmtToCopyTo.IndirectObject, stmtToCopyFrom.IndirectObject)
	stmtToCopyTo.IndirectObjectComplex = copyComponentValue(stmtToCopyTo.IndirectObjectComplex, stmtToCopyFrom.IndirectObjectComplex)
	stmtToCopyTo.IndirectObjectPropertySimple = copyComponentValue(stmtToCopyTo.IndirectObjectPropertySimple, stmtToCopyFrom.IndirectObjectPropertySimple)
	stmtToCopyTo.IndirectObjectPropertyComplex = copyComponentValue(stmtToCopyTo.IndirectObjectPropertyComplex, stmtToCopyFrom.IndirectObjectPropertyComplex)

	// Constitutive
	stmtToCopyTo.ConstitutedEntity = copyComponentValue(stmtToCopyTo.ConstitutedEntity, stmtToCopyFrom.ConstitutedEntity)
	stmtToCopyTo.ConstitutedEntityPropertySimple = copyComponentValue(stmtToCopyTo.ConstitutedEntityPropertySimple, stmtToCopyFrom.ConstitutedEntityPropertySimple)
	stmtToCopyTo.ConstitutedEntityPropertyComplex = copyComponentValue(stmtToCopyTo.ConstitutedEntityPropertyComplex, stmtToCopyFrom.ConstitutedEntityPropertyComplex)
	stmtToCopyTo.Modal = copyComponentValue(stmtToCopyTo.Modal, stmtToCopyFrom.Modal)
	stmtToCopyTo.ConstitutiveFunction = copyComponentValue(stmtToCopyTo.ConstitutiveFunction, stmtToCopyFrom.ConstitutiveFunction)
	stmtToCopyTo.ConstitutingProperties = copyComponentValue(stmtToCopyTo.ConstitutingProperties, stmtToCopyFrom.ConstitutingProperties)
	stmtToCopyTo.ConstitutingPropertiesComplex = copyComponentValue(stmtToCopyTo.ConstitutingPropertiesComplex, stmtToCopyFrom.ConstitutingPropertiesComplex)
	stmtToCopyTo.ConstitutingPropertiesPropertySimple = copyComponentValue(stmtToCopyTo.ConstitutingPropertiesPropertySimple, stmtToCopyFrom.ConstitutingPropertiesPropertySimple)
	stmtToCopyTo.ConstitutingPropertiesPropertyComplex = copyComponentValue(stmtToCopyTo.ConstitutingPropertiesPropertyComplex, stmtToCopyFrom.ConstitutingPropertiesPropertyComplex)

	// Shared components
	stmtToCopyTo.ActivationConditionSimple = copyComponentValue(stmtToCopyTo.ActivationConditionSimple, stmtToCopyFrom.ActivationConditionSimple)
	stmtToCopyTo.ActivationConditionComplex = copyComponentValue(stmtToCopyTo.ActivationConditionComplex, stmtToCopyFrom.ActivationConditionComplex)
	stmtToCopyTo.ExecutionConstraintSimple = copyComponentValue(stmtToCopyTo.ExecutionConstraintSimple, stmtToCopyFrom.ExecutionConstraintSimple)
	stmtToCopyTo.ExecutionConstraintComplex = copyComponentValue(stmtToCopyTo.ExecutionConstraintComplex, stmtToCopyFrom.ExecutionConstraintComplex)
	stmtToCopyTo.OrElse = copyComponentValue(stmtToCopyTo.OrElse, stmtToCopyFrom.OrElse)

	return stmtToCopyTo

}

/*
Overwrites the component value of the statement on which the function is called with the given
component value. Combines input with target component if it is populated.
*/
func copyComponentValue(targetComponent *Node, sourceComponent *Node) *Node {
	if sourceComponent != nil {
		// If target component is empty, simply substitute ...
		if targetComponent == nil {
			targetComponent = sourceComponent
		} else {
			// ... else combine with existing value (using synthetic AND (bAND))
			combinedComponent, err := Combine(targetComponent, sourceComponent, SAND_BETWEEN_COMPONENTS)
			if err.ErrorCode != TREE_NO_ERROR {
				Println("Component copying failed. Error: ", err)
			}
			targetComponent = combinedComponent
		}
	}
	return targetComponent
}

/*
Calculates a statement's complexity and returns a populated StateComplexity struct
that contains options per component, associated states per component, and total
state complexity. It takes into account nested complex statements, as well as leaf elements.
*/
func (s *Statement) CalculateComplexity() StateComplexity {

	// Prepare results structure
	results := StateComplexity{}

	// Regulative

	// Attributes

	results.AttributesOptions = s.Attributes.CountLeaves()
	attributesComplexity, err := s.Attributes.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.AttributesComplexity = attributesComplexity

	// Attributes Property

	results.AttributesPropertySimpleOptions = s.AttributesPropertySimple.CountLeaves()
	attributesPropertySimpleComplexity, err := s.AttributesPropertySimple.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.AttributesPropertySimpleComplexity = attributesPropertySimpleComplexity

	results.AttributesPropertyComplexOptions = s.AttributesPropertyComplex.CountLeaves()
	attributesPropertyComplexComplexity, err := s.AttributesPropertyComplex.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.AttributesPropertyComplexComplexity = attributesPropertyComplexComplexity

	// Deontic

	results.DeonticOptions = s.Deontic.CountLeaves()
	deonticComplexity, err := s.Deontic.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.DeonticComplexity = deonticComplexity

	// Aim

	results.AimOptions = s.Aim.CountLeaves()
	aimComplexity, err := s.Aim.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.AimComplexity = aimComplexity

	// Direct Object

	results.DirectObjectSimpleOptions = s.DirectObject.CountLeaves()
	directObjectSimpleComplexity, err := s.DirectObject.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.DirectObjectSimpleComplexity = directObjectSimpleComplexity

	results.DirectObjectComplexOptions = s.DirectObjectComplex.CountLeaves()
	directObjectComplexComplexity, err := s.DirectObjectComplex.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.DirectObjectComplexComplexity = directObjectComplexComplexity

	// Direct Object Property

	results.DirectObjectPropertySimpleOptions = s.DirectObjectPropertySimple.CountLeaves()
	directObjectPropertySimpleComplexity, err := s.DirectObjectPropertySimple.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.DirectObjectPropertySimpleComplexity = directObjectPropertySimpleComplexity

	results.DirectObjectPropertyComplexOptions = s.DirectObjectPropertyComplex.CountLeaves()
	directObjectPropertyComplexComplexity, err := s.DirectObjectPropertyComplex.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.DirectObjectComplexComplexity = directObjectPropertyComplexComplexity

	// Indirect Object

	results.IndirectObjectSimpleOptions = s.IndirectObject.CountLeaves()
	indirectObjectSimpleComplexity, err := s.IndirectObject.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.IndirectObjectSimpleComplexity = indirectObjectSimpleComplexity

	results.IndirectObjectComplexComplexity = s.IndirectObjectComplex.CountLeaves()
	indirectObjectComplexComplexity, err := s.IndirectObjectComplex.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.IndirectObjectComplexComplexity = s.IndirectObjectComplex.CountLeaves()

	// Indirect Object Property

	results.IndirectObjectPropertySimpleOptions = s.IndirectObjectPropertySimple.CountLeaves()
	indirectObjectPropertySimpleComplexity, err := s.IndirectObjectPropertySimple.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.IndirectObjectSimpleComplexity = indirectObjectPropertySimpleComplexity

	results.IndirectObjectPropertyComplexComplexity = s.IndirectObjectPropertyComplex.CountLeaves()
	indirectObjectPropertyComplexComplexity, err := s.IndirectObjectPropertyComplex.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.IndirectObjectComplexComplexity = indirectObjectPropertyComplexComplexity

	// Constitutive

	// Constituted Entity

	results.ConstitutedEntityOptions = s.ConstitutedEntity.CountLeaves()
	constitutedEntityComplexity, err := s.ConstitutedEntity.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ConstitutedEntityComplexity = constitutedEntityComplexity

	// Constituted Entity Property

	results.ConstitutedEntityPropertySimpleOptions = s.ConstitutedEntityPropertySimple.CountLeaves()
	constitutedEntityPropertySimpleComplexity, err := s.ConstitutedEntityPropertySimple.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ConstitutedEntityPropertySimpleComplexity = constitutedEntityPropertySimpleComplexity

	results.ConstitutedEntityPropertyComplexOptions = s.ConstitutedEntityPropertyComplex.CountLeaves()
	constitutedEntityPropertyComplexComplexity, err := s.ConstitutedEntityPropertyComplex.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ConstitutedEntityPropertyComplexComplexity = constitutedEntityPropertyComplexComplexity

	// Modal

	results.ModalOptions = s.Modal.CountLeaves()
	modalComplexity, err := s.Modal.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ModalComplexity = modalComplexity

	// Constitutive Function

	results.ConstitutiveFunctionOptions = s.ConstitutiveFunction.CountLeaves()
	constitutiveFunctionComplexity, err := s.ConstitutiveFunction.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ConstitutiveFunctionComplexity = constitutiveFunctionComplexity

	// Constituting Properties

	results.ConstitutingPropertiesSimpleOptions = s.ConstitutingProperties.CountLeaves()
	constitutingPropertiesSimpleComplexity, err := s.ConstitutingProperties.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ConstitutingPropertiesSimpleComplexity = constitutingPropertiesSimpleComplexity

	results.ConstitutingPropertiesComplexOptions = s.ConstitutingPropertiesComplex.CountLeaves()
	constitutingPropertiesComplexComplexity, err := s.ConstitutingPropertiesComplex.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ConstitutingPropertiesComplexComplexity = constitutingPropertiesComplexComplexity

	// Constituting Properties Property

	results.ConstitutingPropertiesPropertiesSimpleOptions = s.ConstitutingPropertiesPropertySimple.CountLeaves()
	constitutingPropertiesPropertiesSimpleComplexity, err := s.ConstitutingPropertiesPropertySimple.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ConstitutingPropertiesPropertiesSimpleComplexity = constitutingPropertiesPropertiesSimpleComplexity

	results.ConstitutingPropertiesPropertiesComplexOptions = s.ConstitutingPropertiesPropertyComplex.CountLeaves()
	constitutingPropertiesPropertiesComplexComplexity, err := s.ConstitutingPropertiesPropertyComplex.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ConstitutingPropertiesPropertiesComplexComplexity = constitutingPropertiesPropertiesComplexComplexity

	// Context

	// Activation conditions

	results.ActivationConditionSimpleOptions = s.ActivationConditionSimple.CountLeaves()
	activationConditionSimpleComplexity, err := s.ActivationConditionSimple.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ActivationConditionSimpleComplexity = activationConditionSimpleComplexity

	results.ActivationConditionComplexOptions = s.ActivationConditionComplex.CountLeaves()
	activationConditionComplexComplexity, err := s.ActivationConditionComplex.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ActivationConditionComplexComplexity = activationConditionComplexComplexity

	// Execution constraints

	results.ExecutionConstraintSimpleOptions = s.ExecutionConstraintSimple.CountLeaves()
	executionConstraintSimpleComplexity, err := s.ExecutionConstraintSimple.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ExecutionConstraintSimpleComplexity = executionConstraintSimpleComplexity

	results.ExecutionConstraintComplexOptions = s.ExecutionConstraintComplex.CountLeaves()
	executionConstraintComplexComplexity, err := s.ExecutionConstraintComplex.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.ExecutionConstraintComplexComplexity = executionConstraintComplexComplexity

	// Or else
	orElseComplexity, err := s.OrElse.CalculateStateComplexity()
	if err.ErrorCode != TREE_NO_ERROR {
		fmt.Println(err)
	}
	results.OrElseComplexity = orElseComplexity

	// Composing overall complexity

	// Find highest state complexity on given level
	leadingStmtStates :=
		[]int{ // regulative components
			attributesComplexity,
			attributesPropertySimpleComplexity,
			attributesPropertyComplexComplexity,
			deonticComplexity,
			aimComplexity,
			directObjectSimpleComplexity,
			directObjectComplexComplexity,
			directObjectPropertySimpleComplexity,
			directObjectPropertyComplexComplexity,
			indirectObjectSimpleComplexity,
			indirectObjectComplexComplexity,
			indirectObjectPropertySimpleComplexity,
			indirectObjectPropertyComplexComplexity,
			// constitutive components
			constitutedEntityComplexity,
			constitutedEntityPropertySimpleComplexity,
			constitutedEntityPropertyComplexComplexity,
			modalComplexity,
			constitutiveFunctionComplexity,
			constitutingPropertiesSimpleComplexity,
			constitutingPropertiesComplexComplexity,
			constitutingPropertiesPropertiesSimpleComplexity,
			constitutingPropertiesPropertiesComplexComplexity,
			// execution constraints
			executionConstraintSimpleComplexity,
			executionConstraintComplexComplexity}

	statesOnGivenLevel := shared.AggregateIfGreaterThan(leadingStmtStates, 1, 1)

	// Conditions are handled separately, since they are preconditions, and possibly
	// statements on their own. Default state is 1 (at all times).
	// TODO: Review for accuracy
	conditionsComplexity :=
		shared.FindMaxValue([]int{activationConditionSimpleComplexity +
			activationConditionComplexComplexity}, 1)

	// Consequences are handled separately, since they are separate activities

	// Multiplication of preconditions and leading statement complexity
	results.TotalStateComplexity =
		statesOnGivenLevel * conditionsComplexity

	return results

}
