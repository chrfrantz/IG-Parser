package tree

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Statement struct {

	// Regulative Statement
	Attributes                            *Node
	AttributesPropertySimple              *Node
	AttributesPropertyComplex             *Node
	Deontic                               *Node
	Aim                                   *Node
	DirectObject                          *Node
	DirectObjectComplex                   *Node
	DirectObjectPropertySimple            *Node
	DirectObjectPropertyComplex           *Node
	IndirectObject                        *Node
	IndirectObjectComplex                 *Node
	IndirectObjectPropertySimple          *Node
	IndirectObjectPropertyComplex         *Node

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
	ActivationConditionSimple             *Node
	ActivationConditionComplex            *Node
	ExecutionConstraintSimple             *Node
	ExecutionConstraintComplex            *Node
	OrElse                                *Node
}

func (s *Statement) String() string {
	out := ""

	out = s.printComponent(out, s.Attributes, ATTRIBUTES, false)
	out = s.printComponent(out, s.AttributesPropertySimple, ATTRIBUTES_PROPERTY, false)
	out = s.printComponent(out, s.AttributesPropertyComplex, ATTRIBUTES_PROPERTY, true)
	out = s.printComponent(out, s.Deontic, DEONTIC, false)
	out = s.printComponent(out, s.Aim, AIM, false)
	out = s.printComponent(out, s.DirectObject, DIRECT_OBJECT, false)
	out = s.printComponent(out, s.DirectObjectComplex, DIRECT_OBJECT, true)
	out = s.printComponent(out, s.DirectObjectPropertySimple, DIRECT_OBJECT_PROPERTY, false)
	out = s.printComponent(out, s.DirectObjectPropertyComplex, DIRECT_OBJECT_PROPERTY, true)
	out = s.printComponent(out, s.IndirectObject, INDIRECT_OBJECT, false)
	out = s.printComponent(out, s.IndirectObjectComplex, INDIRECT_OBJECT, true)
	out = s.printComponent(out, s.IndirectObjectPropertySimple, INDIRECT_OBJECT_PROPERTY, false)
	out = s.printComponent(out, s.IndirectObjectPropertyComplex, INDIRECT_OBJECT_PROPERTY, true)

	out = s.printComponent(out, s.ActivationConditionSimple, ACTIVATION_CONDITION, false)
	out = s.printComponent(out, s.ActivationConditionComplex, ACTIVATION_CONDITION, true)
	out = s.printComponent(out, s.ExecutionConstraintSimple, EXECUTION_CONSTRAINT, false)
	out = s.printComponent(out, s.ExecutionConstraintComplex, EXECUTION_CONSTRAINT, true)

	out = s.printComponent(out, s.ConstitutedEntity, CONSTITUTED_ENTITY, false)
	out = s.printComponent(out, s.ConstitutedEntityPropertySimple, CONSTITUTED_ENTITY_PROPERTY, false)
	out = s.printComponent(out, s.ConstitutedEntityPropertyComplex, CONSTITUTED_ENTITY_PROPERTY, true)
	out = s.printComponent(out, s.Modal, MODAL, false)
	out = s.printComponent(out, s.ConstitutiveFunction, CONSTITUTIVE_FUNCTION, false)
	out = s.printComponent(out, s.ConstitutingProperties, CONSTITUTING_PROPERTIES, false)
	out = s.printComponent(out, s.ConstitutingPropertiesComplex, CONSTITUTING_PROPERTIES, true)
	out = s.printComponent(out, s.ConstitutingPropertiesPropertySimple, CONSTITUTING_PROPERTIES_PROPERTY, false)
	out = s.printComponent(out, s.ConstitutingPropertiesPropertyComplex, CONSTITUTING_PROPERTIES_PROPERTY, true)
	
    out = s.printComponent(out, s.OrElse, OR_ELSE, true)

	return out
}

/*
Appends component information for output string
Input:
- input string to append output to
- Node whose content is to be appended
- Symbol associated with component
- Indicator whether component is complex

Returns string for printing
 */
func (s *Statement) printComponent(inputString string, node *Node, nodeSymbol string, complex bool) string {

	sep := ": "
	suffix := "\n"
	complexPrefix := "{\n"
	complexSuffix := "\n}"

	// If node is not nil
	if node != nil {
		// Print symbol
		inputString += nodeSymbol + sep
		// Add core content
		if complex {
			// Complex (i.e., nested) node output
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
Internally, this would be represented as ((Sellers [AND] Buyers] [sAND] (Northern [OR] Southern))', and returned as separate
trees with index 0 (Sellers [AND] Buyers) and 1 (Northern [OR] Southern).
The second return indicates the fact that the first two entries in the first return type instance are of type ATTRIBUTES by holding
an entry '"ATTRIBUTES": 2', etc.

 */
func (s *Statement) GenerateLeafArrays() ([][]*Node, map[string]int) {
	return s.generateLeafArrays(0)
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

*/
func (s *Statement) GenerateLeafArraysSuffixOnly() ([][]*Node, map[string]int) {
	return s.generateLeafArrays(1)
}

/*
Generates map of arrays containing pointers to leaf nodes in each component.
Key is an incrementing index, and value is an array of the corresponding nodes.
It further returns an array containing the component keys alongside the number of leaf nodes per component,
in order to reconstruct the linkage between the index in the first return value and the components they relate to.

Input: level indicates selection of nodes considered in aggregation (0 --> all nodes, 1 --> nodes with non-nil suffix only)

Example: The first return may include two ATTRIBUTES component trees separated by synthetic AND connections (sAND)
based on different logical combination within the attributes component that are not genuine logical relationships (i.e.,
not signaled using [AND], [OR], or [XOR], but inferred during parsing based on the occurrence of multiple such combinations
within an Attributes component expression (e.g., A((Sellers [AND] Buyers) from (Northern [OR] Southern) states)).
Internally, this would be represented as ((Sellers [AND] Buyers] [sAND] (Northern [OR] Southern))', and returned as separate
trees with index 0 (Sellers [AND] Buyers) and 1 (Northern [OR] Southern).
The second return indicates the fact that the first two entries in the first return type instance are of type ATTRIBUTES by holding
an entry '"ATTRIBUTES": 2', etc.

*/
func (s *Statement) generateLeafArrays(level int) ([][]*Node, map[string]int) {

	// Map holding reference from component type (e.g., ATTRIBUTES) to number of entries (relevant for reconstruction)
	referenceMap := map[string]int{}

	// Counter for overall number of entries
	nodesMap := make([][]*Node, 0)

	// Regulative components
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Attributes, ATTRIBUTES, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.AttributesPropertySimple, ATTRIBUTES_PROPERTY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.AttributesPropertyComplex, ATTRIBUTES_PROPERTY_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Deontic, DEONTIC, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Aim, AIM, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObject, DIRECT_OBJECT, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObjectComplex, DIRECT_OBJECT_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObjectPropertySimple, DIRECT_OBJECT_PROPERTY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.DirectObjectPropertyComplex, DIRECT_OBJECT_PROPERTY_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObject, INDIRECT_OBJECT, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObjectComplex, INDIRECT_OBJECT_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObjectPropertySimple, INDIRECT_OBJECT_PROPERTY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.IndirectObjectPropertyComplex, INDIRECT_OBJECT_PROPERTY_REFERENCE, true, level)

	// Context
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ActivationConditionSimple, ACTIVATION_CONDITION, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ActivationConditionComplex, ACTIVATION_CONDITION_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ExecutionConstraintSimple, EXECUTION_CONSTRAINT, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ExecutionConstraintComplex, EXECUTION_CONSTRAINT_REFERENCE, true, level)

	// Constitutive components
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutedEntity, CONSTITUTED_ENTITY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutedEntityPropertySimple, CONSTITUTED_ENTITY_PROPERTY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutedEntityPropertyComplex, CONSTITUTED_ENTITY_PROPERTY_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.Modal, MODAL, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutiveFunction, CONSTITUTIVE_FUNCTION, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingProperties, CONSTITUTING_PROPERTIES, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingPropertiesComplex, CONSTITUTING_PROPERTIES_REFERENCE, true, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingPropertiesPropertySimple, CONSTITUTING_PROPERTIES_PROPERTY, false, level)
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.ConstitutingPropertiesPropertyComplex, CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE, true, level)

	// Shared components
	nodesMap, referenceMap = getComponentLeafArray(nodesMap, referenceMap, s.OrElse, OR_ELSE, true, level)

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
- Indicator whether all leaf nodes should be returned, or only one satisfying particular conditions
  (0 --> all nodes, 1 --> only ones with non-empty suffix).

Returns:
- Node map of nodes associated with components
- Reference map counting number of components
*/
func GetSingleComponentLeafArray(componentNode *Node, componentSymbol string, complex bool, level int) ([][]*Node, map[string]int) {

	// Map holding reference from component type (e.g., ATTRIBUTES) to number of entries (relevant for reconstruction)
	referenceMap := map[string]int{}

	// Counter for overall number of entries
	nodesMap := make([][]*Node, 0)

	return getComponentLeafArray(nodesMap, referenceMap, componentNode, componentSymbol, complex, level)
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
func getComponentLeafArray(nodesMap [][]*Node, referenceMap map[string]int, componentNode *Node, componentSymbol string, complex bool, level int) ([][]*Node, map[string]int) {

	if componentNode == nil {
		fmt.Println("No component node found - returning unmodified node and reference map ...")
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
		for _, v := range componentNode.GetLeafNodes() {
			nodesMap = append(nodesMap, v)
			i++
		}
		// Add number of nodes referring to a particular component
		referenceMap[componentSymbol] = i
	}

	// Return modified or generated structures
	return nodesMap, referenceMap
}

// NODE

type Node struct {
	// Linkage to parent
	Parent *Node
	// Linkage to left child
	Left *Node
	// Linkage to right child
	Right *Node
	// Indicates component type (i.e., name of component)
	ComponentType string
	// Substantive content of a leaf node
	Entry interface{}
	// Text shared across children to the left of a combination (e.g., (shared info (left val [AND] right val)))
	SharedLeft []string
	// Text shared across children to the right of a combination (e.g., ((left val [AND] right val) shared info))
	SharedRight []string
	// Logical operator that links left and right values/nodes
	LogicalOperator string
	// Implicitly holds element order by keeping non-shared elements and references to nodes in order of addition
	ElementOrder []interface{}
	// Suffix for distinctive references to related component instances (e.g., A,p1 pointing to A1)
	Suffix interface{}
	// Annotations for element - to be stored without surrounding brackets
	Annotations interface{}
	// Private links to given node (e.g., private properties)
	PrivateNodeLinks []*Node
}

/*
Returns parents' left shared elements in order of hierarchical (top first).
 */
func (n *Node) getParentsLeftSharedElements() []string {
	if n.Parent != nil && n.Parent.SharedLeft != nil && len(n.Parent.SharedLeft) != 0 {
		// Recursively return parents' shared elements, followed by respective children ones
		return append(n.Parent.getParentsLeftSharedElements(), n.Parent.SharedLeft...)
	} else if n.Parent != nil {
		// Return only parents' shared elements
		return n.Parent.getParentsLeftSharedElements()
	}
	// Return empty structure
	return []string{}
}

/*
Returns parents' right shared elements in order of hierarchical (top first).
*/
func (n *Node) getParentsRightSharedElements() []string {
	if n.Parent != nil && n.Parent.SharedRight != nil && len(n.Parent.SharedRight) != 0 {
		// Recursively return parents' shared elements, followed by respective children ones
		return append(n.Parent.getParentsRightSharedElements(), n.Parent.SharedRight...)
	} else if n.Parent != nil {
		// Return only parents' shared elements
		return n.Parent.getParentsRightSharedElements()
	}
	// Return empty structure
	return []string{}
}

/*
Returns left shared elements under consideration of SHARED_ELEMENT_INHERITANCE_MODE
 */
func (n *Node) GetSharedLeft() []string {
	switch SHARED_ELEMENT_INHERITANCE_MODE {
		case SHARED_ELEMENT_INHERIT_OVERRIDE:
			// Overwrite child with parent shared element values
			shared := n.getParentsLeftSharedElements()
			// If no shared components from parents ...
			if len(shared) == 0 {
				// ... return own shared components
				return n.SharedLeft
			}
			// else return parents' shared components
			return n.getParentsLeftSharedElements()
		case SHARED_ELEMENT_INHERIT_APPEND:
			parentsSharedLeft := n.getParentsLeftSharedElements()
			if len(n.SharedLeft) != 0 && len(parentsSharedLeft) != 0 {
				// Append child's to parents' elements
				return append(parentsSharedLeft, n.SharedLeft...)
			} else if len(n.SharedLeft) != 0 {
				// Return own node information
				return n.SharedLeft
			} else {
				// Return parent node information
				return n.getParentsLeftSharedElements()
			}
		case SHARED_ELEMENT_INHERIT_FROM_COMBINATION:
			if !n.IsCombination() && n.Parent != nil {
				if n.SharedLeft != nil {
					// Return parent and own shared left information
					return append(n.Parent.GetSharedLeft(), n.SharedLeft...)
				}
				// Return parent left shared
				return n.Parent.GetSharedLeft()
			} else {
				// Return shared left (may be of combination of leaf node)
				return n.SharedLeft
			}
		case SHARED_ELEMENT_INHERIT_NOTHING:
			// Simply return own elements
			return n.SharedLeft
	}
	return []string{}
}

/*
Returns right shared elements under consideration of SHARED_ELEMENT_INHERITANCE_MODE
*/
func (n *Node) GetSharedRight() []string {
	switch SHARED_ELEMENT_INHERITANCE_MODE {
	case SHARED_ELEMENT_INHERIT_OVERRIDE:
		// Overwrite child with parent shared element values
		shared := n.getParentsRightSharedElements()
		// If no shared components from parents ...
		if len(shared) == 0 {
			// ... return own shared components
			return n.SharedRight
		}
		// else return parents' shared components
		return n.getParentsRightSharedElements()
	case SHARED_ELEMENT_INHERIT_APPEND:
		parentsSharedRight := n.getParentsRightSharedElements()
		if len(n.SharedRight) != 0 && len(parentsSharedRight) != 0 {
			// Append child's to parents' elements
			return append(parentsSharedRight, n.SharedRight...)
		} else if len(n.SharedRight) != 0 {
			// Return own node information
			return n.SharedRight
		} else {
			// Return parent node information
			return n.getParentsRightSharedElements()
		}
	case SHARED_ELEMENT_INHERIT_FROM_COMBINATION:
		if !n.IsCombination() && n.Parent != nil {
			if n.SharedRight != nil {
				// Return parent and own shared right information
				return append(n.Parent.GetSharedRight(), n.SharedRight...)
			}
			// Return parent right shared
			return n.Parent.GetSharedRight()
		} else {
			// Return shared right (may be of combination of leaf node)
			return n.SharedRight
		}
	case SHARED_ELEMENT_INHERIT_NOTHING:
		// Simply return own elements
		return n.SharedRight
	}
	return []string{}
}

/*
Returns component name stored in component type field. Recursively
iterates through node hierarchy.
 */
func (n *Node) GetComponentName() string {
	// If value is filled
	if n.ComponentType != "" {
		// return content
		return n.ComponentType
	// else test parent node
	} else if n.Parent != nil {
		// retrieve parent information
		return n.Parent.GetComponentName()
	} else {
		// else simply return empty component name
		return n.ComponentType
	}
}

/*
Indicates if node has a primitive consisting of string value, or conversely,
a complex entry consisting of an institutional statement in its own right.
 */
func (n *Node) HasPrimitiveEntry() bool {
	// Check whether the entry is a string
	if _, ok := n.Entry.(string); ok {
		return true
	}
	return false
}

// Sort interface implementation for alphabetic ordering (not order of appearance in tree) of nodes
type ByEntry []*Node

func (e ByEntry) Len() int {
	return len(e)
}

func (e ByEntry) Less(i, j int) bool {
	if e[i].HasPrimitiveEntry() && e[j].HasPrimitiveEntry() {
		return e[i].Entry.(string) < e[j].Entry.(string)
	}
	return false
}

func (e ByEntry) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

/*
Counts the number of parents in the tree hierarchy for a given node.
 */
func (n *Node) CountParents() int {
	if n.Parent == nil {
		return 0
	} else {
		return 1 + n.Parent.CountParents()
	}
}

/*
Returns string representation of node tree structure. The output is
compatible to the input parser to reconstruct the tree.
 */
func (n *Node) Stringify() string {
	// Empty node
	if n.Left == nil && n.Right == nil && n.Entry == "" {
		return ""
	}
	// Leaf node
	if n.IsLeafNode() {
		if n.HasPrimitiveEntry() {
			return n.Entry.(string)
		} else {
			// Return string of statement
			// TODO: REVIEW TO ENSURE IT PRINTS CORRECTLY OR ADJUST TO STRINGIFY()
			val := n.Entry.(Statement)
			return val.String()
		}
	}
	// Walk the tree
	out := ""
	// Potential left shared elements
	if n.GetSharedLeft() != nil && len(n.GetSharedLeft()) != 0 {
		out += "(" + strings.Trim(fmt.Sprint(n.GetSharedLeft()), "[]") + " "
	}
	// Inner combination
	out += "("
	if n.Left != nil {
		out += n.Left.Stringify()
	}
	if n.LogicalOperator != "" {
		out += " [" + n.LogicalOperator + "] " // no extra spacing on left side; due to parsing
	}
	if n.Right != nil {
		out += n.Right.Stringify()
	}
	// Close inner combinations
	out += ")"
	// Potential right shared elements
	if n.GetSharedRight() != nil && len(n.GetSharedRight()) != 0 {
		out += " " + strings.Trim(fmt.Sprint(n.GetSharedRight()), "[]") + ")"
	}
	return out
}


var PrintValueOrder = false

/*
Prints node content in human-readable form (for printing on console).
For parseable version, look at Stringify().
 */
func (n *Node) String() string {

	if n == nil {
		return "Node is not initialized."
	} else if n.IsLeafNode() {
		retVal := "Leaf entry: "
		// TODO: Check for correct printing
		if n.Entry == nil {
			retVal = retVal + "nil (detected in String())"
		} else if n.HasPrimitiveEntry() {
			retVal = retVal + n.Entry.(string)
			// Assumes that suffix and annotations are in string form
			if n.Suffix != nil {
				retVal = retVal + " (Suffix: " + n.Suffix.(string) + ")"
			}
			if n.Annotations != nil {
				retVal = retVal + " (Annotation: " + n.Annotations.(string) + ")"
			}
			if n.PrivateNodeLinks != nil {
				retVal = retVal + " (Private links: " + fmt.Sprint(n.PrivateNodeLinks) + ")"
			}
		} else {
			// if not a string
			val := n.Entry.(Statement)
			retVal = retVal + val.String()
			// Assumes that suffix and annotations are in string format for nodes that have nested statements
			// TODO: see whether that needs to be adjusted
			if n.Suffix != nil {
				retVal = retVal + " (Suffix: " + n.Suffix.(string) + ")"
			}
			if n.Annotations != nil {
				retVal = retVal + " (Annotation: " + n.Annotations.(string) + ")"
			}
		}
		return retVal
		//return /*n.ComponentType + */"Leaf entry: " + n.Entry //+ "\n"
	} else {
		out := ""

		i := 0
		prefix := ""
		for i < n.CountParents() {
			prefix += "===="
			i++
		}

		if len(n.ElementOrder) > 0 && PrintValueOrder {
			i := 0
			for i < len(n.ElementOrder) {
				out += prefix + "Non-Shared: " + fmt.Sprintf("%v", n.ElementOrder[i]) + "\n"
				i++
			}
		}

		// Assumes that suffix and annotations are in string format for nodes that have nested statements
		// TODO: see whether that needs to be adjusted
		if n.Suffix != nil {
			out = " (Suffix: " + n.Suffix.(string) + ")"
		}
		if n.Annotations != nil {
			out = " (Annotation: " + n.Annotations.(string) + ")"
		}

		if n.GetSharedLeft() != nil && len(n.GetSharedLeft()) != 0 {
			fmt.Println("Own LEFT SHARED value (raw content): " + fmt.Sprint(n.SharedLeft) + ", Count: " + strconv.Itoa(len(n.SharedLeft)))
			out += prefix + "Shared (left): " + strings.Trim(fmt.Sprint(n.GetSharedLeft()), "[]") + "\n"
		}
		if n.GetSharedRight() != nil && len(n.GetSharedRight()) != 0 {
			fmt.Println("Own RIGHT SHARED value (raw content): " + fmt.Sprint(n.SharedRight) + ", Count: " + strconv.Itoa(len(n.SharedRight)))
			out += prefix + "Shared (right): " + strings.Trim(fmt.Sprint(n.GetSharedRight()), "[]") + "\n"
		}

		return "(\n" + out +
			prefix + "Left: " + n.Left.String() + "\n" +
			prefix + "Operator: " + n.LogicalOperator + "\n" +
			prefix + "Right: " + n.Right.String() + "\n" +
			prefix + ")"
	}
}

/*
Makes the given node parent of the current (calling node).
Should only be internally used, since it does not deal with child assignment on parent node.
Use InsertLeftNode() or InsertRightNode() instead.
 */
func (n *Node) assignParent(node *Node) (bool, NodeError) {
	if n == node {
		errorMsg := "Attempting to make node parent of itself"
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_SELF_LINKAGE, ErrorMessage: errorMsg}
	}
	n.Parent = node
	return true, NodeError{ErrorCode: TREE_NO_ERROR}
}

/*
Insert left subnode to node
*/
func (n *Node) InsertLeftNode(entry *Node) (bool, NodeError) {
	if n.Left != nil {
		errorMsg := "Attempting to add left node to node already containing left leaf. Node: " + n.String()
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_ADDITION, ErrorMessage: errorMsg}
	}
	if n == entry {
		errorMsg := "Attempting to assign node to itself"
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_SELF_LINKAGE, ErrorMessage: errorMsg}
	}
	if n.Entry != nil {
		errorMsg := "Attempting to add left leaf node to populated node (i.e., it has an entry itself). Node: " + n.String()
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_ADDITION, ErrorMessage: errorMsg}
	}
	entry.assignParent(n)
	n.Left = entry
	return true, NodeError{ErrorCode: TREE_NO_ERROR}
}

/*
Insert right subnode to node
*/
func (n *Node) InsertRightNode(entry *Node) (bool, NodeError) {
	if n.Right != nil {
		errorMsg := "Attempting to add right node to node already containing right leaf. Node: " + n.String()
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_ADDITION, ErrorMessage: errorMsg}
	}
	if n == entry {
		errorMsg := "Attempting to assign node to itself"
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_SELF_LINKAGE, ErrorMessage: errorMsg}
	}
	if n.Entry != nil {
		errorMsg := "Attempting to add right leaf node to populated node (i.e., it has an entry itself). Node: " + n.String()
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_ADDITION, ErrorMessage: errorMsg}
	}
	entry.assignParent(n)
	n.Right = entry
	return true, NodeError{ErrorCode: TREE_NO_ERROR}
}

/*
Insert left leaf to node based on string entry
 */
func (n *Node) InsertLeftLeaf(entry string) (bool, NodeError) {
	if n.Left != nil {
		errorMsg := "Attempting to add left leaf node to node already containing left leaf. Node: " + n.String()
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_ADDITION, ErrorMessage: errorMsg}
	}
	if n.Entry != nil {
		errorMsg := "Attempting to add left leaf node to populated node (i.e., it has an entry itself). Node: " + n.String()
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_ADDITION, ErrorMessage: errorMsg}
	}
	newNode := Node{}
	newNode.Entry = entry
	newNode.assignParent(n)
	n.Left = &newNode
	return true, NodeError{ErrorCode: TREE_NO_ERROR}
}

/*
Insert right leaf to node based on string entry
*/
func (n *Node) InsertRightLeaf(entry string) (bool, NodeError) {
	if n.Right != nil {
		errorMsg := "Attempting to add right leaf node to node already containing right leaf. Node: " + n.String()
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_ADDITION, ErrorMessage: errorMsg}
	}
	if n.Entry != nil {
		errorMsg := "Attempting to add right leaf node to populated node (i.e., it has an entry itself). Node: " + n.String()
		log.Println(errorMsg)
		return false, NodeError{ErrorCode: TREE_INVALID_NODE_ADDITION, ErrorMessage: errorMsg}
	}
	newNode := Node{}
	newNode.Entry = entry
	newNode.assignParent(n)
	n.Right = &newNode
	return true, NodeError{ErrorCode: TREE_NO_ERROR}
}

/*
Removes the given node from the tree structure it is embedded in, i.e.,
it does not have a parent and the parent is no longer aware of this child.

Returns boolean indicating success and potential error (in success case TREE_NO_ERROR).
 */
func RemoveNodeFromTree(node *Node) (bool, NodeError) {

	if node.Parent != nil {
		// Remove parent's reference to child, and collapse tree structure of necessary
		if node.Parent.Left == node {
			// If the parent is a combination and the node's parent has a parent
			if node.Parent.IsCombination() && node.Parent.Parent != nil {
				// If the sibling node on the right is not nil
				if node.Parent.Right != nil {
					// Check whether the combination sits on the left side of its parent
					if node.Parent.Parent.Left == node.Parent {
						// and if sitting on the left, assign the former right sibling in place of the combination
						node.Parent.Parent.Left = node.Parent.Right
						// and adjust former right node to link to new parent
						node.Parent.Right.assignParent(node.Parent.Parent)
					} else if node.Parent.Parent.Right == node.Parent {
						// and if sitting on the right, assign the former right sibling in place of the combination
						node.Parent.Parent.Right = node.Parent.Right
						// and adjust former right node to link to new parent
						node.Parent.Right.assignParent(node.Parent.Parent)
					}
				}
			} else if node.Parent.IsCombination() {
				// if the node's parent is a combination (but the parent does not have a parent on its own),
				// then simply assign former sibling as root (i.e., modify parent node of passed node)
				fmt.Println("Assigned right as root")
				*node.Parent = *node.Parent.Right
			}
		} else if node.Parent.Right == node {
			// If the parent is a combination and the node's parent has a parent
			if node.Parent.IsCombination() && node.Parent.Parent != nil {
				// if the sibling on the left is not nil
				if node.Parent.Left != nil {
					// Check whether the combination sits on the left side of its parent
					if node.Parent.Parent.Left == node.Parent {
						// and if sitting on the left, assign the former left sibling in place of the combination
						node.Parent.Parent.Left = node.Parent.Left
						// and adjust former left node to link to new parent
						node.Parent.Left.assignParent(node.Parent.Parent)
					} else if node.Parent.Parent.Right == node.Parent {
						// and if sitting on the right, assign the former left sibling in place of the combination
						node.Parent.Parent.Right = node.Parent.Left
						// and adjust former left node to link to new parent
						node.Parent.Left.assignParent(node.Parent.Parent)
					}
				}
			} else if node.Parent.IsCombination() {
					// if the node's parent is a combination (but the parent does not have a parent on its own),
					// then simply assign former sibling as root (i.e., modify parent node of passed node)
					fmt.Println("Assigned left as root")
					*node.Parent = *node.Parent.Left
			}
		} else {
			errorMsg := "Could not find linkage of parent node in tree structure to ensure proper rebalancing following removal of node."
			return false, NodeError{ErrorCode: TREE_INVALID_NODE_REMOVAL, ErrorMessage: errorMsg}
		}

		// Remove reference from child to parent
		node.Parent = nil

		// Now node should be disconnected and tree reorganized
		return true, NodeError{ErrorCode: TREE_NO_ERROR}
	}
	// else tag the removal as invalid
	errorMsg := "Attempted to remove already disconnected node from parent tree"
	return false, NodeError{ErrorCode: TREE_INVALID_NODE_REMOVAL, ErrorMessage: errorMsg}
}

/*
Finds logical linkages between a source and target node in the tree they are embedded in.
Returns true if a link is found, and provides the logical operators on that path.
It further returns an error in case of navigation challenges (with error case TREE_NO_ERROR
signaling successful navigation irrespective of outcome.
 */
func FindLogicalLinkage(sourceNode *Node, targetNode *Node) (bool, []string, NodeError) {

	// Test down first
	foundDownwards, ops, err := searchDownward(sourceNode, sourceNode, sourceNode, targetNode, []string{})
	if err.ErrorCode != TREE_NO_ERROR {
		return false, ops, err
	}
	// If found in downwards search, return
	if foundDownwards {
		return true, ops, err
	}

	// Test up
	foundUpwards, ops, err := searchUpward(sourceNode, sourceNode, targetNode, []string{})
	if err.ErrorCode != TREE_NO_ERROR {
		return false, ops, err
	}
	// If found in downwards search, return
	if foundUpwards {
		return true, ops, err
	}

	fmt.Println("Could not find target node ", targetNode, " from start node ", sourceNode)
	return false, nil, err
}

/*
Searches for a given node upward in the tree structure by recursively moving upwards, while excluding previously explored
unsuccessful branches (Parameter lastNode). The parameter originNode and targetNode are retained as reference to the search
origin and target. opsPath retains all logical operators along the path.
Returns true in case of successful outcome, with logical operators provided alongside. It further returns an error
(per default with error code TREE_NO_ERROR indicating that no navigation issue occurred throughout the tree - irrespective of the outcome).
 */
func searchUpward(originNode *Node, lastNode *Node, targetNode *Node, opsPath []string) (bool, []string, NodeError) {

	// If the node from which search is initiated does not have a parent, return false, since no upward exploration is possible
	if lastNode.Parent == nil {
		return false, opsPath, NodeError{ErrorCode: TREE_NO_ERROR}
	}

	fmt.Println("Searching downward from ", lastNode.Parent)
	// Search unexplored neighbouring path, starting with the input's node's parent - prevent repeated exploration of the input node path lastNode
	response, ops, err := searchDownward(originNode, lastNode, lastNode.Parent, targetNode, opsPath)
	if err.ErrorCode != TREE_NO_ERROR {
		return false, ops, err
	}

	// If not successful, recurse upwards, and attempt again, with reference to the explore parent as last node (to prevent repeated exploration)
	if !response {
		// Explicit include logical operator if moving upward
		opsPath = append(opsPath, lastNode.Parent.LogicalOperator)
		fmt.Println("Search one level higher above ", lastNode.Parent)
		response, ops, err = searchUpward(originNode, lastNode.Parent, targetNode, opsPath)
	}

	return response, ops, err
}

/*
Searches downwards from a given node (startNode), exploring any left branch first (recursively), followed by the right one (recursively).
The lastNode indicates the last node explored in preceding search (i.e., a nested node), preventing the exploration of the underlying path.
The originNode retains the reference to the search origin, and targetNode is the reference to the target.
opsPath retains all operators found along the path.
Returns true if successful, alongside the relevant logical operators along the path, as well as an error
(per default with error code TREE_NO_ERROR indicating that no navigation issue occurred throughout the tree - irrespective of the outcome).
 */
func searchDownward(originNode *Node, lastNode *Node, startNode *Node, targetNode *Node, opsPath []string) (bool, []string, NodeError) {

	// Check if input node is already target
	if startNode == targetNode {
		fmt.Println("Search node ", startNode, " is target node")
		return true, opsPath, NodeError{ErrorCode: TREE_NO_ERROR}
	}

	// if leaf and does not match, return false
	if startNode.IsLeafNode() {
		fmt.Println("Search node ", startNode, " is leaf, but does not corresponding to the target node. Node: ", targetNode)
		return false, opsPath, NodeError{ErrorCode: TREE_NO_ERROR}
	}

	// Carry over input operator
	ops := []string{}
	if opsPath != nil {
		fmt.Println("Inherits operators ", opsPath)
		ops = append(ops, opsPath...)
	}

	// Append start node operator
	ops = append(ops, startNode.LogicalOperator)
	fmt.Println("Added logical operator ", startNode.LogicalOperator)

	// Predefine response values
	err := NodeError{ErrorCode: TREE_NO_ERROR}

	// Test left first - it must not be nil, and not be the last explored node (i.e., a left child of the currently explored one)
	if startNode.Left != nil && startNode.Left != lastNode {
		fmt.Println("Test left branch ...")
		response, ops2, err := searchDownward(originNode, startNode.Left, startNode.Left, targetNode, ops)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops2, err
		}
		// If positive outcome
		if response {
			fmt.Println("Found target on left side")
			return true, ops2, err
		}
		// Delegate downwards
		fmt.Println("- Test left left")
		response, ops3, err := searchDownward(originNode, startNode.Left.Left, startNode.Left.Left, targetNode, ops2)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops3, err
		}
		// If positive outcome
		if response {
			fmt.Println("Found target on left left side")
			return true, ops3, err
		}
		fmt.Println("- Test left right")
		response, ops3, err = searchDownward(originNode, startNode.Left.Right, startNode.Left.Right, targetNode, ops2)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops3, err
		}
		// If positive outcome
		if response {
			fmt.Println("Found target on left right side")
			return true, ops3, err
		}
	}
	// Test right (will only be done if left was not successful)
	// Right node must not be nil, and not be the last explored node (i.e., a right child of the currently explored one)
	if startNode.Right != nil && startNode.Right != lastNode {
		fmt.Println("Testing right branch ...")
		response, ops2, err := searchDownward(originNode, startNode.Right, startNode.Right, targetNode, ops)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops2, err
		}
		// If positive outcome
		if response {
			fmt.Println("Found target on right side")
			return true, ops2, err
		}
		// Delegate downwards
		fmt.Println("- Test right left")
		response, ops3, err := searchDownward(originNode, startNode.Right.Left, startNode.Right.Left, targetNode, ops2)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops3, err
		}
		// If positive outcome
		if response {
			fmt.Println("Found target on right left side")
			return true, ops3, err
		}
		fmt.Println("- Test right right")
		response, ops3, err = searchDownward(originNode, startNode.Right.Right, startNode.Right.Right, targetNode, ops2)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops3, err
		}
		// If positive outcome
		if response {
			fmt.Println("Found target on right right side")
			return true, ops3, err
		}
	}
	fmt.Println("Final result: false")
	return false, ops, err
}

/*
Combines existing nodes into new node and returns newly generated node
 */
func Combine(leftNode *Node, rightNode *Node, logicalOperator string) *Node {

	if leftNode == nil && rightNode == nil {
		log.Fatal("Illegal call to Combine() with nil nodes")
	}
	if leftNode == nil || leftNode.IsEmptyNode() {
		fmt.Println("Combining nodes returns right node (other node is nil or empty)")
		return rightNode
	}
	if rightNode == nil || rightNode.IsEmptyNode() {
		fmt.Println("Combining nodes returns left node (other node is nil or empty)")
		return leftNode
	}
	// In all other cases, create new combination using provided logical operator
	newNode := Node{}
	newNode.Left = leftNode
	newNode.Left.Parent = &newNode
	newNode.Right = rightNode
	newNode.Right.Parent = &newNode
	newNode.LogicalOperator = logicalOperator
	return &newNode
}


/**
Adds non-shared values to the node, i.e., values that are not shared across subnodes, but attached to the
node itself.
 */
func (n *Node) InsertNonSharedValues(value string) {
	n.ElementOrder = append(n.ElementOrder, value)
}

/*
Creates a generic node, with various options
 */
func ComponentNode(entry string, leftValue string, rightValue string, componentType string, sharedValueLeft []string, sharedValueRight []string, logicalOperator string) *Node {

	// Validation (Entry cannot be mixed with the other fields)
	if entry != "" {
		if leftValue != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as left-hand node (" + leftValue + ").")
		}
		if rightValue != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as right-hand node (" + rightValue + ").")
		}
		if sharedValueLeft != nil || len(sharedValueLeft) > 0 {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as (left) shared content field (" + fmt.Sprint(sharedValueLeft) + ").")
		}
		if sharedValueRight != nil || len(sharedValueRight) > 0 {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as (right) shared content field (" + fmt.Sprint(sharedValueRight) + ").")
		}
		if logicalOperator != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry +
			"), as well as logical operator field (" + logicalOperator + ").")
		}
	}
	// Validation (Check whether left-, right-hand, and logical operator are filled)
	if entry == "" && (leftValue == "" || rightValue == "" || logicalOperator == "") {
		log.Fatal("Non-leaf node, but missing specification of left-hand, right-hand value, " +
		"or logical operator (Left hand: " + leftValue + "; Right hand: " + rightValue +
		"; Logical operator: " + logicalOperator + ")")
	}
	if logicalOperator != "" {
		res, _ := StringInSlice(logicalOperator, IGLogicalOperators)
		if !res {
			log.Fatal("Logical operator value invalid (Value: " + logicalOperator + ")")
		}
	}

	// Specification must be valid. Continue with creation ...

	node := Node{}

	// Assign parent as nil
	node.Parent = nil

	// Inherit (if not specified), or assign component name from parameters
	if componentType == "" {
		log.Println("Inheriting component type from parent ... (" + componentType + ")")
		if componentType != "" {
			node.ComponentType = componentType
		}
	} else {
		log.Println("Assigning component type " + componentType)
		node.ComponentType = componentType
	}

	// If leaf node, fill all relevant fields
	if entry != "" {
		node.Entry = entry
	} else {
		// if non-leaf, fill all relevant fields
		node.InsertLeftLeaf(leftValue)
		node.InsertRightLeaf(rightValue)
		node.LogicalOperator = logicalOperator
		node.SharedLeft = sharedValueLeft
		node.SharedRight = sharedValueRight
	}
	return &node
}

/*
Validates all nodes from this node downwards with respect to population as linking node or leaf node.
 */
func (n *Node) Validate() (bool, NodeError){
	if n.Entry == nil && (n.Left == nil || n.Right == nil) {
		errorMsg := "Non-leaf node, but missing specification of left and right child, " +
			"or both. Node: " + fmt.Sprint(n.String())
		return false, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: errorMsg}
	}
	if n.Entry != nil && (n.Left != nil  || n.Right != nil) {
		errorMsg := "Leaf node, but still has filled left or right node. Node: " + fmt.Sprint(n.String())
		return false, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: errorMsg}
	}
	if n.Left != nil && n.Right != nil && n.LogicalOperator == "" {
		errorMsg := "Did not specify logical operator in populated tree. Node: " + fmt.Sprint(n.String())
		return false, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: errorMsg}
	}

	if !n.IsLeafNode() {
		downwardResult := false
		err := NodeError{}
		// Move downwards
		if n.Left == nil {
			return false, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: "Empty left node"}
		} else {
			downwardResult, err = n.Left.Validate()
		}
		if !downwardResult {
			return false, err
		}
		if n.Right == nil {
			return false, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: "Empty right node"}
		} else {
			downwardResult, err = n.Right.Validate()
		}
		if !downwardResult {
			return false, err
		}
	}
	return true, NodeError{ErrorCode: TREE_NO_ERROR}
}

/*
Counts the number of leaves of node tree
 */
func (n *Node) CountLeaves() int {
	if n == nil {
		// Uninitialized node
		return 0
	}
	if n.Left == nil && n.Right == nil && n.Entry == "" {
		// Must be empty node
		return 0
	}
	if n.Left == nil && n.Right == nil && n.Entry != "" {
		// Must be single leaf node (entry)
		return 1
	}
	leftBreadth := 0
	rightBreadth := 0
	if n.Left != nil {
		leftBreadth = n.Left.CountLeaves()
	}
	if n.Right != nil {
		rightBreadth = n.Right.CountLeaves()
	}
	return leftBreadth + rightBreadth
}

/*
Returns root node of given tree the node is embedded in
up to the level at which nodes are linked by synthetic AND (sAND).
 */
func (n *Node) GetSyntheticRootNode() *Node {
	if n.Parent == nil || n.Parent.LogicalOperator == SAND {
		// Assume to be parent if no parent on its own,
		// or root in synthetic hierarchy if paired with sAND
		return n
	} else {
		// else delegate to parent
		return n.Parent.GetSyntheticRootNode()
	}
}

/*
Returns leaf nodes of a given node as arrays of arrays of nodes.
Note: currently, all leaf arrays are stored at array[0], thus not
exploiting the multi-dimensional nature
 */
func (n *Node) GetLeafNodes() [][]*Node {
	if n == nil {
		// Uninitialized node
		return nil
	}
	returnNode := make([][]*Node, 0)
	if n.Left == nil && n.Right == nil && n.Entry == "" {
		// Must be empty node
		return returnNode
	}
	if n.Left == nil && n.Right == nil && n.Entry != "" {
		// Must be single leaf node (entry)
		inner := []*Node{n}
		returnNode = append(returnNode, inner)
		return returnNode
	}
	leftNodes := [][]*Node{}
	rightNodes := [][]*Node{}

	// If both left and right children nodes exist, return those combined
	if n.Left != nil && n.Right != nil {
		leftNodes = n.Left.GetLeafNodes()
		rightNodes = n.Right.GetLeafNodes()
		if n.LogicalOperator == SAND {
			// Return as collection of node collection
			// Append individual nodes arrays
			for _, v := range leftNodes {
				returnNode = append(returnNode, v)
			}
			// Append individual nodes arrays
			for _, v := range rightNodes {
				returnNode = append(returnNode, v)
			}
			return returnNode
		} else {
			nodeArray := make([]*Node, 0)
			// return as individual nodes
			for _, v := range leftNodes {
				nodeArray = append(nodeArray, v...)
			}
			for _, v := range rightNodes {
				nodeArray = append(nodeArray, v...)
			}
			// Appends as first array (second remains empty)
			returnNode = append(returnNode, nodeArray)
			return returnNode
		}
	}
	// Process left nodes
	if n.Left != nil {
		leftNodes = n.Left.GetLeafNodes()
		for _, v := range leftNodes {
			returnNode[0] = append(returnNode[0], v...)
		}
	}
	// Process right nodes
	if n.Right != nil {
		rightNodes = n.Right.GetLeafNodes()
		for _, v := range rightNodes {
			returnNode[0] = append(returnNode[0], v...)
		}
	}
	return returnNode
}

/*
Calculate depth of node tree
 */
func (n *Node) CalculateDepth() int {
	if n == nil {
		return 0
	}
	if n.Left == nil && n.Right == nil {
		return 0
	}
	leftDepth := 0
	rightDepth := 0
	if n.Left != nil {
		leftDepth = 1 + n.Left.CalculateDepth()
	}
	if n.Right != nil {
		rightDepth = 1 + n.Right.CalculateDepth()
	}
	if leftDepth < rightDepth {
		return rightDepth
	} else {
		return leftDepth
	}
}

/*
Indicates whether node is leaf node
 */
func (n *Node) IsLeafNode() bool {
	return n == nil || (n.Left == nil && n.Right == nil)
}

/*
Indicates whether node contains valid combination (i.e., left and right and logical operator are populated).
*/
func (n *Node) IsCombination() bool {
	return n.Entry == nil && !n.Left.IsNil() &&
		!n.Right.IsNil() && n.LogicalOperator != ""
}

/*
Indicates whether node has populated logical operator, but does not check for proper assignment of left and right children.
 */
func (n *Node) hasLogicalOperator() bool {
	return n.LogicalOperator != ""
}

/*
Indicates whether node is empty
 */
func (n *Node) IsEmptyNode() bool {
	return n.IsLeafNode() && n.Entry == nil
}

/*
Indicates if node is nil
 */
func (n *Node) IsNil() bool {
	return n == nil
}

/*
Applies statement parsing function to all entries below a given node.
 */
func (n *Node) ParseAllEntries(function func(string) (Statement, ParsingError)) ParsingError {
	if n.IsNil() {
		return ParsingError{ErrorCode: PARSING_ERROR_NIL_ELEMENT, ErrorMessage: "Attempted to parse nil element."}
	}
	if !n.IsEmptyNode() && n.Entry != nil {

		// Execute actual function
		newEntry, err := function(n.Entry.(string))
		if err.ErrorCode != PARSING_NO_ERROR {
			return err
		}
		// and reassign parsed element
		n.Entry = newEntry
	}
	if !n.Left.IsNil() {
		n.Left.ParseAllEntries(function)
	}
	if !n.Right.IsNil() {
		n.Right.ParseAllEntries(function)
	}
	return ParsingError{ErrorCode: PARSING_NO_ERROR}
}