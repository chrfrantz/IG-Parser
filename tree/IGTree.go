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
Prints node content in human-readable form (for printing on console).
For parseable version, look at Stringify().
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
	if n.Entry != "" {
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
	if n.Entry != "" {
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
Returns the logical linkages between source and target node, where
both are leaf nodes
 */
func FindLogicalLinkage(sourceNode *Node, targetNode *Node, opsOnPath []string) (bool, []string, NodeError) {

	if sourceNode == nil || targetNode == nil {
		return false, opsOnPath, NodeError{ErrorCode: TREE_INPUT_VALIDATION, ErrorMessage: "No source or target node provided. " +
			"This is often the case when navigating tree structures with missing/empty leaves. Consider validating tree to " +
			"ensure the absence of structural gaps."}
	}
	// Override source node for search if provided
	n := sourceNode

	if n.Parent == nil || n.Parent.IsEmptyNode() {
		return false, opsOnPath, NodeError{ErrorCode: TREE_INPUT_VALIDATION,
			ErrorMessage: fmt.Sprint("Can't search for related node, since no parent node. Node: ", n)}
	}

	// Inherit operators if existing
	ops := []string{}
	if opsOnPath != nil {
		ops = opsOnPath
	}

	// If input is output, just return
	if sourceNode == targetNode {
		return true, ops, NodeError{ErrorCode: TREE_NO_ERROR}
	}
	fmt.Println("Searching for " + targetNode.String() + " on node " + n.String())


	// Search immediate nodes
	if n.Left == targetNode {
		//fmt.Println("Before operators: ", ops)
		//ops = append(ops, n.LogicalOperator)
		fmt.Println("Found target on the left side.")
		//fmt.Println("Added operator ", n.LogicalOperator)
		//fmt.Println("Current operators: ", ops)
		return true, ops, NodeError{ErrorCode: TREE_NO_ERROR}
	}
	if n.Right == targetNode {
		//fmt.Println("Before operators: ", ops)
		//ops = append(ops, n.LogicalOperator)
		fmt.Println("Found target on the right side.")
		//fmt.Println("Added operator ", n.LogicalOperator)
		//fmt.Println("Current operators: ", ops)
		return true, ops, NodeError{ErrorCode: TREE_NO_ERROR}
	}

	outError := NodeError{}
	hitParent := false

	// Switch to signal whether target has been found
	result := false
	// if on the left side
	if n.Parent.Left == n {
		//fmt.Println("Searching on the right side of node ", n.Parent)
		// Search on the right side
		if n.Parent.Right == targetNode {
			fmt.Println("Found node on right side")

			if len(ops) == 0 {
				fmt.Println("Added operator ", n.Parent.LogicalOperator)
				ops = append(ops, n.Parent.LogicalOperator)
				fmt.Println("Current operators: ", ops)
			}

			return true, ops, NodeError{ErrorCode: TREE_NO_ERROR}
		}
		// if right sibling is non-leaf, delegate
		if !n.Parent.Right.IsLeafNode() {
			// Add first path element
			if len(ops) == 0 {
				fmt.Println("Added operator ", n.Parent.LogicalOperator)
				ops = append(ops, n.Parent.LogicalOperator)
				fmt.Println("Current operators: ", ops)
			}

			fmt.Println("Search left side downwards on right node ", n.Parent.Right.Left)

			// Add operator of nested node to search downwards
			fmt.Println("Added operator ", n.Parent.Right.LogicalOperator)
			ops = append(ops, n.Parent.Right.LogicalOperator)
			fmt.Println("Current operators: ", ops)

			// Delegate search the left child of neighbouring right combination
			result, ops, outError = FindLogicalLinkage(n.Parent.Right.Left, targetNode, ops)
			if outError.ErrorCode != TREE_NO_ERROR {
				return result, ops, outError
			}
			// Intentionally not checking for errors here, since search on empty is internally permissible
			if !result {
				fmt.Println("Search right side downwards on right node ", n.Parent.Right.Right)
				// Delegate search the right child of neighbouring right combination
				result, ops, outError = FindLogicalLinkage(n.Parent.Right.Right, targetNode, ops)
				if outError.ErrorCode != TREE_NO_ERROR {
					return result, ops, outError
				}
			}
		}
	} else 	// if on the right
	if n.Parent.Right == n {
		//fmt.Println("Searching on the left side of node ", n.Parent)
		// Search on the left side
		if n.Parent.Left == targetNode {
			fmt.Println("Found node on left side")

			if len(ops) == 0 {
				fmt.Println("Added operator ", n.Parent.LogicalOperator)
				ops = append(ops, n.Parent.LogicalOperator)
				fmt.Println("Current operators: ", ops)
			}

			return true, ops, NodeError{ErrorCode: TREE_NO_ERROR}
		}
		// if left sibling is non-leaf, delegate
		if !n.Parent.Left.IsLeafNode() {
			// Add first path element
			if len(ops) == 0 {
				fmt.Println("Added operator ", n.Parent.LogicalOperator)
				ops = append(ops, n.Parent.LogicalOperator)
				fmt.Println("Current operators: ", ops)
			}

			fmt.Println("Search left side downwards on left node ", n.Parent.Left.Left)

			// Add operator of nested node to search downwards
			fmt.Println("Added operator ", n.Parent.Left.LogicalOperator)
			ops = append(ops, n.Parent.Left.LogicalOperator)
			fmt.Println("Current operators: ", ops)

			// Delegate search the left child of neighbouring left combination
			result, ops, outError = FindLogicalLinkage(n.Parent.Left.Left, targetNode, ops)
			if !result {
				fmt.Println("Search right side downwards on left node ", n.Parent.Left.Right)
				// Delegate search the right child of neighbouring left combination
				result, ops, outError = FindLogicalLinkage(n.Parent.Left.Right, targetNode, ops)
			}
		}
	}
	// If nothing has been found until here, go up in the hierarchy - unless we have been there
	if !result && !hitParent {
		fmt.Println("Delegating to parent: ", n.Parent.Parent)
		// Add first path element (needs to be inverted at the end)
		// Own operator (i.e., in own parent) is appended by default - any other leaf needs to be linked by logical operator
		ops = append(ops, n.Parent.LogicalOperator)
		fmt.Println("Added operator ", n.Parent.LogicalOperator)
		fmt.Println("Current operators: ", ops)
		// Parent's parent operator
		ops = append(ops, n.Parent.Parent.LogicalOperator)
		fmt.Println("Added operator ", n.Parent.Parent.LogicalOperator)
		fmt.Println("Current operators: ", ops)
		// Mark whether root node has been reached
		if n.Parent.Parent.Parent == nil {
			fmt.Println("Hit ceiling")
			hitParent = true
		}
		if n.Parent.Parent.Left == n.Parent {
			fmt.Println("Pushing down the right side ... on ", n.Parent.Parent.Right)
			// Push towards right side
			ops = append(ops, n.Parent.Parent.Right.LogicalOperator)
			result, ops, outError = FindLogicalLinkage(n.Parent.Parent.Right, targetNode, ops)
		} else {
			fmt.Println("Pushing down the left side ... on ", n.Parent.Parent.Left)
			// Push towards left side
			ops = append(ops, n.Parent.Parent.Left.LogicalOperator)
			result, ops, outError = FindLogicalLinkage(n.Parent.Parent.Left, targetNode, ops)
		}
	}
	if !result {
		fmt.Println("No links between ", sourceNode, " and ", targetNode, " found.")
		//ops = append(ops, n.LogicalOperator)
	}
	return result, ops, outError
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
/*
func (n *Node) InsertLeafNode(entry string) *Node {
	return n.InsertChildNode(entry, "", "", n.ComponentType, nil, nil, "")
}

func ComponentLeafNode(entry string, componentType string) *Node {
	return ComponentNode(entry, "", "", componentType, nil, nil, "")
}*/

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
		error := NodeError{}
		// Move downwards
		if n.Left == nil {
			return false, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: "Empty left node"}
		} else {
			downwardResult, error = n.Left.Validate()
		}
		if !downwardResult {
			return false, error
		}
		if n.Right == nil {
			return false, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: "Empty right node"}
		} else {
			downwardResult, error = n.Right.Validate()
		}
		if !downwardResult {
			return false, error
		}
	}
	return true, NodeError{ErrorCode: TREE_NO_ERROR}
}

/*
Creates a new node based on parameter specification. If non-root node, the node should be created
within parent node (to ensure proper association as either left- or right-hand node).
If entry value is specified, the node is presumed to be leaf node; in all other instances, the
nodes is interpreted as combination, and left and right values are moved into respective leaf nodes.
Component type name is saved in node.
 */
/*
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
}*/

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