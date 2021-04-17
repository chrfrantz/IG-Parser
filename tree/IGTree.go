package tree

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Statement struct {
	// Regulative Statement
	Attributes *Node
	AttributesProperty *Node
	Deontic *Node
	Aim *Node
	DirectObject *Node
	DirectObjectProperty *Node
	IndirectObject *Node
	IndirectObjectProperty *Node
	//Constitutive Statement
	ConstitutedEntity *Node
	ConstitutedEntityProperty *Node
	Modal *Node
	ConstitutiveFunction *Node
	ConstitutingProperties *Node
	ConstitutingPropertiesProperty *Node
	// Shared Components
	ActivationConditionSimple *Node
	ActivationConditionComplex *Node
	ExecutionConstraintSimple *Node
	ExecutionConstraintComplex *Node
	OrElse *Node
}

func (s *Statement) String() string {
	out := ""
	sep := ": "
	suffix := "\n"
	complexPrefix := "{\n"
	complexSuffix := "\n}"
	if s.Attributes != nil {
		out += ATTRIBUTES + sep
		out += s.Attributes.String()
		out += suffix
	}
	if s.AttributesProperty != nil {
		out += ATTRIBUTES_PROPERTY + sep
		out += s.AttributesProperty.String()
		out += suffix
	}
	if s.Deontic != nil {
		out += DEONTIC + sep
		out += s.Deontic.String()
		out += suffix
	}
	if s.Aim != nil {
		out += AIM + sep
		out += s.Aim.String()
		out += suffix
	}
	if s.DirectObject != nil {
		out += DIRECT_OBJECT + sep
		out += s.DirectObject.String()
		out += suffix
	}
	if s.DirectObjectProperty != nil {
		out += DIRECT_OBJECT_PROPERTY + sep
		out += s.DirectObjectProperty.String()
		out += suffix
	}
	if s.IndirectObject != nil {
		out += INDIRECT_OBJECT + sep
		out += s.IndirectObject.String()
		out += suffix
	}
	if s.IndirectObjectProperty != nil {
		out += INDIRECT_OBJECT_PROPERTY + sep
		out += s.IndirectObjectProperty.String()
		out += suffix
	}
	if s.ConstitutedEntity != nil {
		out += CONSTITUTED_ENTITY + sep
		out += s.ConstitutedEntity.String()
		out += suffix
	}
	if s.ConstitutedEntityProperty != nil {
		out += CONSTITUTING_PROPERTIES_PROPERTY + sep
		out += s.ConstitutedEntityProperty.String()
		out += suffix
	}
	if s.Modal != nil {
		out += MODAL + sep
		out += s.Modal.String()
		out += suffix
	}
	if s.ConstitutiveFunction != nil {
		out += CONSTITUTIVE_FUNCTION + sep
		out += s.ConstitutiveFunction.String()
		out += suffix
	}
	if s.ConstitutingProperties != nil {
		out += CONSTITUTING_PROPERTIES + sep
		out += s.ConstitutingProperties.String()
		out += suffix
	}
	if s.ConstitutingPropertiesProperty != nil {
		out += CONSTITUTING_PROPERTIES_PROPERTY + sep
		out += s.ConstitutingPropertiesProperty.String()
		out += suffix
	}
	if s.ActivationConditionSimple != nil {
		out += ACTIVATION_CONDITION + sep
		out += s.ActivationConditionSimple.String()
		out += suffix
	}
	if s.ActivationConditionComplex != nil {
		out += ACTIVATION_CONDITION + sep
		out += complexPrefix + s.ActivationConditionComplex.String() + complexSuffix
		out += suffix
	}
	if s.ExecutionConstraintSimple != nil {
		out += EXECUTION_CONSTRAINT + sep
		out += s.ExecutionConstraintSimple.String()
		out += suffix
	}
	if s.ExecutionConstraintComplex != nil {
		out += EXECUTION_CONSTRAINT + sep
		out += s.ExecutionConstraintComplex.String()
		out += suffix
	}
	if s.OrElse != nil {
		out += s.OrElse.String()
	}
	return out
}


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

	// Map holding reference from component type (e.g., ATTRIBUTES) to number of entries (relevant for reconstruction)
	referenceMap := map[string]int{}

	// Counter for overall number of entries
	ct := 0
	nodesMap := make([][]*Node, 0)

	// Counter for number of elements in given component
	i := 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.Attributes.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[ATTRIBUTES] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.AttributesProperty.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[ATTRIBUTES_PROPERTY] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.Deontic.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[DEONTIC] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.Aim.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[AIM] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.DirectObject.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[DIRECT_OBJECT] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.DirectObjectProperty.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[DIRECT_OBJECT_PROPERTY] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.IndirectObject.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[INDIRECT_OBJECT] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.IndirectObjectProperty.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[INDIRECT_OBJECT_PROPERTY] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.ActivationConditionSimple.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[ACTIVATION_CONDITION] = i

	// Check for complex activation conditions
	if s.ActivationConditionComplex != nil {
		nodesMap = append(nodesMap, []*Node{s.ActivationConditionComplex})
		referenceMap[ACTIVATION_CONDITION_REFERENCE] = 1
	}
	// Add array of leaf nodes attached to general array
	/*for _, v := range s.ActivationConditionComplex.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[ACTIVATION_CONDITION] = i*/

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.ExecutionConstraintSimple.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[EXECUTION_CONSTRAINT] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.ConstitutedEntity.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[CONSTITUTED_ENTITY] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.ConstitutedEntityProperty.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[CONSTITUTED_ENTITY_PROPERTY] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.Modal.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[MODAL] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.ConstitutiveFunction.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[CONSTITUTIVE_FUNCTION] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.ConstitutingProperties.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[CONSTITUTING_PROPERTIES] = i

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.ConstitutingPropertiesProperty.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[CONSTITUTING_PROPERTIES_PROPERTY] = i

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
	// Indicates component type
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
}

/*
Returns parents' left shared elements in order of hierarchical (top first).
 */
func (n *Node) getParentsLeftSharedElements() []string {
	if n.Parent != nil && n.Parent.SharedLeft != nil && len(n.Parent.SharedLeft) != 0 {
		// Recursively return parents' shared elements, followed by respective children ones
		return append(n.Parent.getParentsLeftSharedElements(), n.Parent.SharedLeft...)
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
			if len(n.SharedLeft) != 0 && len(n.getParentsLeftSharedElements()) != 0 {
				// Append child's to parents' elements
				return append(n.getParentsLeftSharedElements(), n.SharedLeft...)
			} else if len(n.SharedLeft) != 0 {
				// Return own node information
				return n.SharedLeft
			} else {
				// Return parent node information
				return n.getParentsLeftSharedElements()
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
		if len(n.SharedRight) != 0 && len(n.getParentsRightSharedElements()) != 0 {
			// Append child's to parents' elements
			return append(n.getParentsRightSharedElements(), n.SharedRight...)
		} else if len(n.SharedRight) != 0 {
			// Return own node information
			return n.SharedRight
		} else {
			// Return parent node information
			return n.getParentsRightSharedElements()
		}
	case SHARED_ELEMENT_INHERIT_NOTHING:
		// Simply return own elements
		return n.SharedRight
	}
	return []string{}
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
			return retVal + "nil (detected in String())"
		} else if n.HasPrimitiveEntry() {
			return retVal + n.Entry.(string)
		} else {
			val := n.Entry.(Statement)
			return retVal + val.String()
		}
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
Insert right subnode to node
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
	if n.Entry != "" {
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
	if n.Entry != "" {
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
	response := false
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
	fmt.Println("Final result: ", response)
	return response, ops, err
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
		if !StringInSlice(logicalOperator, IGLogicalOperators) {
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
	if n.Entry == "" && (n.Left == nil || n.Right == nil) {
		errorMsg := "Non-leaf node, but missing specification of left and right child, " +
			"or both. Node: " + fmt.Sprint(n.String())
		return false, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: errorMsg}
	}
	if n.Entry != "" && (n.Left != nil  || n.Right != nil) {
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
Returns leaf nodes of a given node as arrays of arrays of nodes
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
Indicates whether node is empty
 */
func (n *Node) IsEmptyNode() bool {
	return n.IsLeafNode() && n.Entry == ""
}