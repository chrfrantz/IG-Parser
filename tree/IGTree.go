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
	ActivationConditionComplex *Statement
	ExecutionConstraintSimple *Node
	ExecutionConstraintComplex *Statement
	OrElse *Statement
}

func (s *Statement) String() string {
	out := ""
	sep := ": "
	suffix := "\n"
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
		out += s.ActivationConditionComplex.String()
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
		//out += " "
	}
	return out
}


func (s *Statement) Stringify() string {
	log.Fatal("Stringify() is not yet implemented.")
	return ""
}

/*
Generates map of arrays containing pointers to leaf nodes in each component. Key is component name constant (e.g., ATTRIBUTES).
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

	// Counter for number of elements in given component
	i = 0
	// Add array of leaf nodes attached to general array
	for _, v := range s.ExecutionConstraintSimple.GetLeafNodes() {
		nodesMap = append(nodesMap, v)
		i++
		ct++
	}
	referenceMap[EXECUTION_CONSTRAINT] = i

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
	Entry string
	// Text shared across children to the left of a combination (e.g., (shared info (left val [AND] right val)))
	SharedLeft []string
	// Text shared across children to the right of a combination (e.g., ((left val [AND] right val) shared info))
	SharedRight []string
	// Logical operator that links left and right values/nodes
	LogicalOperator string
	// Implicitly holds element order by keeping non-shared elements and references to nodes in order of addition
	ElementOrder []interface{}
}

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
		return n.Entry
	}
	// Walk the tree
	out := ""
	// Potential left shared elements
	if n.SharedLeft != nil && len(n.SharedLeft) != 0 {
		out += "(" + strings.Trim(fmt.Sprint(n.SharedLeft), "[]") + " "
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
	if n.SharedRight != nil && len(n.SharedRight) != 0 {
		out += " " + strings.Trim(fmt.Sprint(n.SharedRight), "[]") + ")"
	}
	return out
}


var PrintValueOrder = false

/*
Prints node content
 */
func (n *Node) String() string {

	if n == nil {
		return "Node is not initialized."
	} else if n.IsLeafNode() {
		return /*n.ComponentType + */"Leaf entry: " + n.Entry //+ "\n"
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

		if n.SharedLeft != nil && len(n.SharedLeft) != 0 {
			fmt.Println("LEFT CONTAINS: " + fmt.Sprint(n.SharedLeft) + strconv.Itoa(len(n.SharedLeft)))
			out += prefix + "Shared (left): " + strings.Trim(fmt.Sprint(n.SharedLeft), "[]") + "\n"
		}
		if n.SharedRight != nil && len(n.SharedRight) != 0 {
			fmt.Println("RIGHT CONTAINS: " + fmt.Sprint(n.SharedRight) + strconv.Itoa(len(n.SharedRight)))
			out += prefix + "Shared (right): " + strings.Trim(fmt.Sprint(n.SharedRight), "[]") + "\n"
		}

		return "(\n" + out +
			prefix + "Left: " + n.Left.String() + "\n" +
			prefix + "Operator: " + n.LogicalOperator + "\n" +
			prefix + "Right: " + n.Right.String() + "\n" +
			prefix + ")"
	}
}

/*
Insert left leaf to node
 */
func (n *Node) InsertLeftLeaf(entry string) {
	if n.Left != nil {
		log.Fatal("Attempting to add left leaf node to node already containing left leaf. Node: " + n.String())
	}
	if n.Entry != "" {
		log.Fatal("Attempting to add left leaf node to populated node (i.e., it has an entry itself). Node: " + n.String())
	}
	newNode := Node{}
	newNode.Entry = entry
	newNode.Parent = n
	n.Left = &newNode
}

/*
Insert right leaf to node
*/
func (n *Node) InsertRightLeaf(entry string) {
	if n.Right != nil {
		log.Fatal("Attempting to add right leaf node to node already containing right leaf. Node: " + n.String())
	}
	if n.Entry != "" {
		log.Fatal("Attempting to add right leaf node to populated node (i.e., it has an entry itself). Node: " + n.String())
	}
	newNode := Node{}
	newNode.Entry = entry
	newNode.Parent = n
	n.Right = &newNode
}

/*
Returns the logical linkages between source and target node, where
both are leaf nodes
 */
func FindLogicalLinkage(sourceNode *Node, targetNode *Node, opsOnPath []string) (bool, []string) {

	if sourceNode == nil || targetNode == nil {
		log.Fatal("No source or target node provided")
	}
	// Override source node for search if provided
	n := sourceNode

	if n.Parent == nil || n.Parent.IsEmptyNode() {
		log.Fatal("Can't search for related node, since no parent node. Node: ", n)
	}

	// Inherit operators if existing
	ops := []string{}
	if opsOnPath != nil {
		ops = opsOnPath
	}

	// Switch to signal whether target has been found
	result := false
	// if on the left side
	if n.Parent.Left == n {
		//fmt.Println("Searching on the right side of node ", n.Parent)
		// Search on the right side
		if n.Parent.Right == targetNode {
			//fmt.Println("Found node on right side")
			ops = append(ops, n.Parent.LogicalOperator)
			return true, ops
		}
		// if right sibling is non-leaf, delegate
		if !n.Parent.Right.IsLeafNode() {
			//fmt.Println("Search right side downwards on node ", n.Parent.Right.Left)
			// Add first path element (needs to be inverted at the end)
			ops = append(ops, n.Parent.LogicalOperator)
			// Delegate search the left child of neighbouring right combination
			result, ops = FindLogicalLinkage(n.Parent.Right.Left, targetNode, ops)
			if !result {
				result, ops = FindLogicalLinkage(n.Parent.Right.Right, targetNode, ops)
			}
		}
	} else 	// if on the right
	if n.Parent.Right == n {
		//fmt.Println("Searching on the left side of node ", n.Parent)
		// Search on the left side
		if n.Parent.Left == targetNode {
			fmt.Println("Found node on left side")
			ops = append(ops, n.Parent.LogicalOperator)
			return true, ops
		}
		// if left sibling is non-leaf, delegate
		if !n.Parent.Left.IsLeafNode() {
			//fmt.Println("Search left side downwards on node ", n.Parent.Left.Left)
			// Add first path element (needs to be inverted at the end)
			ops = append(ops, n.Parent.LogicalOperator)
			// Delegate search the left child of neighbouring left combination
			result, ops = FindLogicalLinkage(n.Parent.Left.Left, targetNode, ops)
			if !result {
				result, ops = FindLogicalLinkage(n.Parent.Left.Right, targetNode, ops)
			}
		}
	}
	// If nothing has been found until here, go up in the hierarchy
	if !result {
		//fmt.Println("Delegating to parent: ", n.Parent.Parent)
		// Add first path element (needs to be inverted at the end)
		ops = append(ops, n.Parent.LogicalOperator)
		result, ops = FindLogicalLinkage(n.Parent, targetNode, ops)
	}
	if !result {
		fmt.Println("No links between ", sourceNode, " and ", targetNode, " found.")
	}
	return result, ops
}

/*
Combines existing nodes into new node and returns newly generated node
 */
func Combine(leftnode *Node, rightNode *Node, logicalOperator string) *Node {

	if leftnode == nil && rightNode == nil {
		log.Fatal("Illegal call to Combine() with nil nodes")
	}
	if leftnode == nil || leftnode.IsEmptyNode() {
		fmt.Println("Combining nodes returns right node (other node is nil or empty)")
		return rightNode
	}
	if rightNode == nil || rightNode.IsEmptyNode() {
		fmt.Println("Combining nodes returns left node (other node is nil or empty)")
		return leftnode
	}
	// In all other cases, create new combination using provided logical operator
	newNode := Node{}
	newNode.Left = leftnode
	newNode.Left.Parent = &newNode
	newNode.Right = rightNode
	newNode.Right.Parent = &newNode
	newNode.LogicalOperator = logicalOperator
	return &newNode
}

/*
Inserts node under the current node (not replacing it), and associates any existing node with an AND combination,
pushing it to the bottom of the tree.
The added node itself can be of any kind, i.e., either have a nested structure or be a leaf node.
 */
/*func (n *Node) Insert(node *Node, logicalOperator string) *Node {

	fmt.Println("Insert into ... ")
	// If this node is empty, assign new one to it
	if n.IsEmptyNode() {
		node.Parent = nil
		n = node
		fmt.Println("Empty node (Overwrite)")
		return n
	} else if n.IsLeafNode() {
		// make combination
		fmt.Println("Leaf node --> make combination (old left, new right)")
		//newNode := Node{}
		newNode := n
		n.Left = newNode//&Node{Entry: n.Entry,ComponentType: n.ComponentType, Parent: n}
		n.Left.Parent = n
		n.LogicalOperator = logicalOperator
		n.Entry = ""
		n.Right = node
		n.Right.Parent = n
		n.ElementOrder = append(n.ElementOrder, node)
		return n
	}*/

	/*
	if n.Left == nil {
		// If left is empty, assign there, ...
		// Assign new parent
		node.Parent = n
		n.Left = node
		n.ElementOrder = append(n.ElementOrder, node)
		return n
	} else if n.Left != nil && n.Right == nil {
		// else try on the right
		// Assign new parent
		n.Right = node
		node.Parent = n
		n.LogicalOperator = AND
		n.ElementOrder = append(n.ElementOrder, node)
		return n
	} else if n.Left != nil && n.Right != nil {*/
		// Delegate to right child node to deal with it ...
		// TODO: Insert on right side to retain order (should be reviewed for balance, but ok for now)
		//fmt.Println("Delegate to right side ...")
		//return n.Right.Insert(node, logicalOperator)
	//}
	//return n
//}

/**
Adds non-shared values to the node, i.e., values that are not shared across subnodes, but attached to the
node itself.
 */
func (n *Node) InsertNonSharedValues(value string) {
	n.ElementOrder = append(n.ElementOrder, value)
}

/*
Inserts a leaf node under a given node and inherits its component type
 */
func (n *Node) InsertLeafNode(entry string) *Node {
	return n.InsertChildNode(entry, "", "", n.ComponentType, nil, nil, "")
}

func ComponentLeafNode(entry string, componentType string) *Node {
	return ComponentNode(entry, "", "", componentType, nil, nil, "")
}

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
		node.Left = node.InsertLeafNode(leftValue)
		node.Right = node.InsertLeafNode(rightValue)
		node.LogicalOperator = logicalOperator
		node.SharedLeft = sharedValueLeft
		node.SharedRight = sharedValueRight
	}
	return &node
}

/*
Creates a new node based on parameter specification. If non-root node, the node should be created
within parent node (to ensure proper association as either left- or right-hand node).
If entry value is specified, the node is presumed to be leaf node; in all other instances, the
nodes is interpreted as combination, and left and right values are moved into respective leaf nodes.
Component type name is saved in node.
 */
func (n *Node) InsertChildNode(entry string, leftValue string, rightValue string, componentType string, sharedValueLeft []string, sharedValueRight []string, logicalOperator string) *Node {
	// Validation (Entry cannot be mixed with the other fields)
	if entry != "" {
		if leftValue != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as left-hand node (" + leftValue + ").")
		}
		if rightValue != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as right-hand node (" + rightValue + ").")
		}
		if sharedValueLeft != nil && len(sharedValueLeft) != 0 {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as (left) shared content field (" + fmt.Sprint(sharedValueLeft) + ").")
		}
		if sharedValueRight != nil && len(sharedValueRight) != 0 {
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
	// Assign node on which this function is called as parent
	node.Parent = n

	// Inherit (if not specified), or assign component name from parameters
	if componentType == "" {
		log.Println("Inheriting component type from parent ... (" + n.ComponentType + ")")
		if n.ComponentType != "" {
			node.ComponentType = n.ComponentType
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
		node.Left = n.InsertLeafNode(leftValue)
		node.Right = n.InsertLeafNode(rightValue)
		node.LogicalOperator = logicalOperator
		node.SharedLeft = sharedValueLeft
		node.SharedRight = sharedValueRight
	}
	return &node
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
			returnNode = append(returnNode, nodeArray)
			return returnNode
		}
	}
	if n.Left != nil {
		leftNodes = n.Left.GetLeafNodes()
		for _, v := range leftNodes {
			returnNode[0] = append(returnNode[0], v...)
		}
	}
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