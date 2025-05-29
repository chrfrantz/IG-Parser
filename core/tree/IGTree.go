package tree

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

/*
This includes basic tree parsing functionality underlying the statement tree construction.
*/

type Node struct {
	// Linkage to parent
	Parent *Node
	// Linkage to left child
	Left *Node
	// Linkage to right child
	Right *Node
	// Indicates component type (i.e., name of component) - Note: Access via GetComponentName(); not directly
	ComponentType string
	// Substantive content of a leaf node (can be nested node, string or nested statement)
	Entry interface{}
	// Text shared across children to the left of a combination (e.g., (shared info (left val [AND] right val)))
	SharedLeft []string
	// Text shared across children to the right of a combination (e.g., ((left val [AND] right val) shared info))
	SharedRight []string
	// Logical operator that links left and right values/nodes
	LogicalOperator string
	// Implicitly holds element order by keeping non-shared elements and references to nodes in order of addition
	ElementOrder []interface{}
	// Suffix for distinctive references to related component instances (e.g., A1,p pointing to A1)
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

	if n != nil && n.Parent != nil {
		// Only allows for shared elements that are truly not-empty (including preventing "" as first element)
		if n.Parent.SharedLeft != nil && len(n.Parent.SharedLeft) != 0 && n.Parent.SharedLeft[0] != "" {
			// Recursively return parents' shared elements, followed by respective children ones
			return append(n.Parent.getParentsLeftSharedElements(), n.Parent.SharedLeft...)
		} else {
			// Return only parents' shared elements
			return n.Parent.getParentsLeftSharedElements()
		}
	}
	// Return empty structure
	return nil
}

/*
Returns parents' right shared elements in order of hierarchical (top first).
*/
func (n *Node) getParentsRightSharedElements() []string {

	if n != nil && n.Parent != nil {
		// Only allows for shared elements that are truly not-empty (including preventing "" as first element)
		if n.Parent.SharedRight != nil && len(n.Parent.SharedRight) != 0 && n.Parent.SharedRight[0] != "" {
			// Recursively return parents' shared elements, followed by respective children ones
			return append(n.Parent.getParentsRightSharedElements(), n.Parent.SharedRight...)
		} else {
			// Return only parents' shared elements
			return n.Parent.getParentsRightSharedElements()
		}
	}
	// Return empty structure
	return nil
}

/*
Returns left shared elements under consideration of SHARED_ELEMENT_INHERITANCE_MODE
*/
func (n *Node) GetSharedLeft() []string {
	if n == nil {
		return nil
	}
	switch SHARED_ELEMENT_INHERITANCE_MODE {
	case SHARED_ELEMENT_INHERIT_OVERRIDE:
		// Overwrite child with parent shared element values
		shared := n.getParentsLeftSharedElements()
		// If no shared components from parents ...
		if shared == nil || len(shared) == 0 || shared[0] == "" {
			// ... return own shared components
			return n.SharedLeft
		}
		// else return parents' shared components
		return n.getParentsLeftSharedElements()
	case SHARED_ELEMENT_INHERIT_APPEND:
		parentsSharedLeft := n.getParentsLeftSharedElements()
		if n != nil && n.SharedLeft != nil && len(n.SharedLeft) != 0 && n.SharedLeft[0] != "" && len(parentsSharedLeft) != 0 {
			// Append child's to parents' elements
			return append(parentsSharedLeft, n.SharedLeft...)
		} else if n != nil && len(n.SharedLeft) != 0 && n.SharedLeft[0] != "" {
			// Return own node information
			return n.SharedLeft
		} else if n != nil && n.Parent != nil {
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
	return nil
}

/*
Returns right shared elements under consideration of SHARED_ELEMENT_INHERITANCE_MODE
*/
func (n *Node) GetSharedRight() []string {
	if n == nil {
		return nil
	}
	switch SHARED_ELEMENT_INHERITANCE_MODE {
	case SHARED_ELEMENT_INHERIT_OVERRIDE:
		// Overwrite child with parent shared element values
		shared := n.getParentsRightSharedElements()
		// If no shared components from parents ...
		if shared == nil || len(shared) == 0 || shared[0] == "" {
			// ... return own shared components
			return n.SharedRight
		}
		// else return parents' shared components
		return n.getParentsRightSharedElements()
	case SHARED_ELEMENT_INHERIT_APPEND:
		parentsSharedRight := n.getParentsRightSharedElements()
		if n != nil && n.SharedRight != nil && len(n.SharedRight) != 0 && n.SharedRight[0] != "" && len(parentsSharedRight) != 0 {
			// Append child's to parents' elements
			return append(parentsSharedRight, n.SharedRight...)
		} else if n != nil && len(n.SharedRight) != 0 && n.SharedRight[0] != "" {
			// Return own node information
			return n.SharedRight
		} else if n != nil && n.Parent != nil {
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
	return nil
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
	} else if n.Parent != nil && n.Parent != n {
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

	printSharedParentheses := false
	// Check whether left or right shared elements exist (to ensure correct syntax (parentheses))
	if n.GetSharedLeft() != nil || n.GetSharedRight() != nil {
		printSharedParentheses = true
	}

	// Print potentially needed opening parentheses (to balance right parentheses)
	if printSharedParentheses {
		// print left parenthesis
		out += "("
	}
	// Potential left shared elements; note: explicitly checks for empty elements (e.g., "")
	if n.GetSharedLeft() != nil && len(n.GetSharedLeft()) != 0 && n.GetSharedLeft()[0] != "" {
		out += strings.Trim(fmt.Sprint(n.GetSharedLeft()), "[]") + " "
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
	// Potential right shared elements; note: explicitly checks for empty elements (e.g., "")
	if n.GetSharedRight() != nil && len(n.GetSharedRight()) != 0 && n.GetSharedRight()[0] != "" {
		out += " " + strings.Trim(fmt.Sprint(n.GetSharedRight()), "[]")
	}
	// Print potentially needed closing parentheses (to balance left parentheses)
	if printSharedParentheses {
		// print right parenthesis
		out += ")"
	}
	return out
}

/*
Provide the flattest possible representation of contents without structure or linebreaks.
Cautious: Can produce log.Fatal in case of novel unknown type. Developer consideration.
*/
func (n *Node) StringFlat() string {

	if n.HasPrimitiveEntry() {
		// Simple entry
		return n.Entry.(string)
	} else if n.IsCombination() {
		out := ""
		// Prepend left shared elements
		if n.SharedLeft != nil && len(n.SharedLeft) != 0 && n.SharedLeft[0] != "" {
			for _, v := range n.SharedLeft {
				out += v + " "
			}
		}
		out += n.Left.StringFlat() + " " + n.LogicalOperator + " " + n.Right.StringFlat()
		// Append right shared elements
		if n.SharedRight != nil && len(n.SharedRight) != 0 && n.SharedRight[0] != "" {
			out += " "
			for _, v := range n.SharedRight {
				out += v + " "
			}
		}
		return out
	} else if n.Entry == nil {
		// else simply return empty string
		return ""
	} else {
		// Complex entry by reference
		if reflect.TypeOf(n.Entry) == reflect.TypeOf(&Statement{}) {
			stmt := n.Entry.(*Statement)
			return stmt.StringFlat(false)
		}
		// Component-level nested pair combination
		if reflect.TypeOf(n.Entry) == reflect.TypeOf([]*Node{}) {
			stmt := n.Entry.([]*Node)[0].Entry.(*Statement)
			return stmt.StringFlat(false)
		}
		// Complex entry by value (rarest occurrence, hence last)
		if reflect.TypeOf(n.Entry) == reflect.TypeOf(Statement{}) {
			stmt := n.Entry.(Statement)
			return stmt.StringFlat(false)
		}
		log.Fatal("Error when generating string output. Unknown type (developer concern):", reflect.TypeOf(n.Entry))
		return ""
	}
}

var PrintValueOrder = false

/*
Returns node state.
*/
func (n *Node) GetNodeState() string {
	if n == nil {
		return NODE_STATE_NIL
	} else {
		return NODE_STATE_NON_NIL
	}
}

/*
Prints node content in human-readable form (for printing on console).
For parseable version, look at Stringify().
*/
func (n *Node) String() string {
	return n.string(0)
}

/*
Prints node content in human-readable form (for printing on console).
For parseable version, look at Stringify().
*/
func (n *Node) StringLevel(level int) string {
	return n.string(level)
}

// Indentation unit for statement tree printing
const MinimumIndentPrefix = "===="

/*
Variant of String() to produce human-readable output, but allows
for parameterization of level of indentation for recursive calls
*/
func (n *Node) string(level int) string {

	// Indent any output depending on level of nesting
	prefix := ""

	if n == nil {
		return "Node is not initialized."
	} else if n.IsLeafNode() {
		retVal := ""

		if n.Entry == nil {
			retVal = retVal + "nil (detected in String())"
			if n.LogicalOperator != "" {
				retVal = retVal + "\nLogical Operator: " + n.LogicalOperator
			}
		} else if n.HasPrimitiveEntry() {
			// Primitive component

			// Indent only for entries that are part of an elementary combination (have logical operator as parent).
			// If parent is statement (i.e., primitive component in statement), the parent is nil
			if n.Parent != nil {
				i := 0
				for i < level {
					retVal += MinimumIndentPrefix
					i++
				}
			}

			retVal = retVal + n.Entry.(string)
			// Assumes that suffix and annotations are in string form
			if n.GetSuffix() != "" {
				retVal = retVal + " (Suffix: " + n.GetSuffix() + ")"
			}
			if n.Annotations != nil {
				retVal = retVal + " (Annotation: " + n.Annotations.(string) + ")"
			}
			if n.PrivateNodeLinks != nil {
				retVal = retVal + " (Private links: " + fmt.Sprint(n.PrivateNodeLinks) + ")"
			}
			if n.GetComponentName() != "" {
				retVal = retVal + " (Component name: " + fmt.Sprint(n.GetComponentName()) + ")"
			}
		} else {

			// Check whether it returns embedded node structure - points to incorrect implementation of parsing based on node-embedded return structures
			if reflect.TypeOf([]*Node{}) == reflect.TypeOf(n.Entry) {

				// Iterate through nodes and parse depending on whether statement or simple node
				entries := n.Entry.([]*Node)
				entriesPrint := ""
				for _, v := range entries {
					if reflect.TypeOf(Statement{}) == reflect.TypeOf(v.Entry) {
						// Check whether statement is embedded ...
						stmt := v.Entry.(Statement)
						entriesPrint += stmt.string(level)
					} else {
						// ... else print node entry as string
						entriesPrint += v.string(level)
					}
				}
				retVal = retVal + entriesPrint //"empty Node - CHECK for unmanaged output case"
			} else {
				if reflect.TypeOf(n.Entry) == reflect.TypeOf(Statement{}) {
					// Assume entry is statement
					val := n.Entry.(Statement)
					retVal = retVal + val.string(level+1)
				} else {
					// Assume entry is statement reference
					val := n.Entry.(*Statement)
					retVal = retVal + val.string(level+1)
				}
			}
		}
		return retVal
	} else {
		// Nodes: Combinations (e.g., AND-combined components)
		out := ""

		// Increase level, since nested element
		level += 1

		i := 0
		for i < level {
			prefix += MinimumIndentPrefix
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
			Println("Own LEFT SHARED value (raw content): " + fmt.Sprint(n.SharedLeft) + ", Count: " + strconv.Itoa(len(n.SharedLeft)))
			out += prefix + "Shared (left): " + strings.Trim(fmt.Sprint(n.GetSharedLeft()), "[]") + "\n"
		}
		if n.GetSharedRight() != nil && len(n.GetSharedRight()) != 0 {
			Println("Own RIGHT SHARED value (raw content): " + fmt.Sprint(n.SharedRight) + ", Count: " + strconv.Itoa(len(n.SharedRight)))
			out += prefix + "Shared (right): " + strings.Trim(fmt.Sprint(n.GetSharedRight()), "[]") + "\n"
		}

		// Higher-level nesting of combinations - indentation from current level
		retPrep := "\n" + prefix + "(\n" + //out +
			prefix + "Left: \n" + n.Left.string(level+1) + "\n"
		if n.Left.GetSharedLeft() != nil {
			retPrep += prefix + " - Left shared (left): " + fmt.Sprint(n.Left.GetSharedLeft()) + "\n"
		}
		if n.Left.GetSharedRight() != nil {
			retPrep += prefix + " - Left shared (right): " + fmt.Sprint(n.Left.GetSharedRight()) + "\n"
		}
		retPrep += prefix + "Operator: " + n.LogicalOperator + "\n" +
			prefix + "Right: \n" + n.Right.string(level+1) + "\n"
		if n.Right.GetSharedLeft() != nil {
			retPrep += prefix + " - Right shared (left): " + fmt.Sprint(n.Right.GetSharedLeft()) + "\n"
		}
		if n.Right.GetSharedRight() != nil {
			retPrep += prefix + " - Right shared (right): " + fmt.Sprint(n.Right.GetSharedRight()) + "\n"
		}
		retPrep += prefix + ")"

		// Assumes that suffix and annotations are in string format for nodes that have nested statements
		// TODO: see whether that needs to be adjusted
		if n.GetSuffix() != "" {
			retPrep += " (Suffix: " + n.GetSuffix() + ")"
		}
		if n.Annotations != nil {
			retPrep += " (Annotation: " + n.Annotations.(string) + ")"
		}
		if n.PrivateNodeLinks != nil {
			retPrep += " (Private links: " + fmt.Sprint(n.PrivateNodeLinks) + ")"
		}
		if n.GetComponentName() != "" {
			retPrep += " (Component name: " + fmt.Sprint(n.GetComponentName()) + ")"
		}

		return retPrep
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
	if !n.Left.IsEmptyOrNilNode() {
		errorMsg := "Attempting to add left node to node already containing left leaf. Node:" + n.Left.String()
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
	if !n.Right.IsEmptyOrNilNode() {
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
		// Remove parent's reference to child, and collapse tree structure if necessary
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
				Println("Assigned right as root")
				// Remove pointer to parent
				node.Parent.Right.Parent = nil
				// Reassign
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
				Println("Assigned left as root")
				// Remove pointer to parent
				node.Parent.Left.Parent = nil
				// Reassign
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
	errorMsg := "Attempted to remove already disconnected node (" + node.String() + ") from parent tree"
	return false, NodeError{ErrorCode: TREE_INVALID_NODE_REMOVAL, ErrorMessage: errorMsg}
}

/*
Finds logical linkages between a source and target node in the tree they are embedded in.
Returns true if a link is found, and provides the logical operators on that path.
It further returns an error in case of navigation challenges (with error case TREE_NO_ERROR
signaling successful navigation irrespective of outcome).
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
	// If found in upwards search, return
	if foundUpwards {
		return true, ops, err
	}

	Println("Could not find target node ", targetNode, " from start node ", sourceNode)
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

	//Println("Searching downward from ", lastNode.Parent)
	// Search unexplored neighbouring path, starting with the input's node's parent - prevent repeated exploration of the input node path lastNode
	response, ops, err := searchDownward(originNode, lastNode, lastNode.Parent, targetNode, opsPath)
	if err.ErrorCode != TREE_NO_ERROR {
		return false, ops, err
	}

	// If not successful, recurse upwards, and attempt again, with reference to the explore parent as last node (to prevent repeated exploration)
	if !response {

		// Explicitly include logical operator if moving upward (if populated)
		if lastNode.Parent.LogicalOperator != "" {
			opsPath = append(opsPath, lastNode.Parent.LogicalOperator)
		}
		//Println("Search one level higher above ", lastNode.Parent)
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
		//Println("Search node ", startNode, " is target node")
		return true, opsPath, NodeError{ErrorCode: TREE_NO_ERROR}
	}

	// if leaf and does not match && if entry does not contain Node[] collection (in which case there are embedded statements), return false
	if startNode.IsLeafNode() && (startNode == nil || startNode.Entry == nil || reflect.TypeOf(startNode.Entry) != reflect.TypeOf([]*Node{})) {
		Println("Search node ", startNode, " is leaf, but does not corresponding to the target node. Node: ", targetNode)
		return false, opsPath, NodeError{ErrorCode: TREE_NO_ERROR}
	}

	// Carry over input operator
	ops := []string{}
	if opsPath != nil {
		//Println("Inherits operators ", opsPath)
		ops = append(ops, opsPath...)
	}

	// Append start node operator
	if startNode.LogicalOperator != "" {
		ops = append(ops, startNode.LogicalOperator)
		//Println("Added logical operator ", startNode.LogicalOperator)
	}

	// Predefine response values
	errSuccess := NodeError{ErrorCode: TREE_NO_ERROR}

	// If left and right are empty, check for []Node in entry (extrapolated statement) - should only hold if no branch (i.e., left and right) is populated
	// NOTE: Avoid upward search following downward search into []Node, since that may lead to infinite loop
	if startNode.Left == nil && startNode.Right == nil && reflect.TypeOf(startNode.Entry) == reflect.TypeOf([]*Node{}) {
		for _, v := range startNode.Entry.([]*Node) {
			response, ops4, err := searchDownward(originNode, startNode, v, targetNode, ops)
			// return lacking success if appearing
			if err.ErrorCode != TREE_NO_ERROR {
				return false, ops4, err
			}
			// If positive outcome
			if response {
				Println("Found target in nested statement collection (extrapolated structure)")
				return true, ops4, errSuccess
			}
		}
	}

	// Test left first - it must not be nil, and not be the last explored node (i.e., a left child of the currently explored one)
	if startNode.Left != nil && startNode.Left != lastNode {
		//Println("Test left branch ...")
		response, ops2, err := searchDownward(originNode, startNode.Left, startNode.Left, targetNode, ops)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops2, err
		}
		// If positive outcome
		if response {
			//Println("Found target on left side")
			return true, ops2, errSuccess
		}
		// Delegate downwards
		//Println("- Test left left")
		response, ops3, err := searchDownward(originNode, startNode.Left.Left, startNode.Left.Left, targetNode, ops2)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops3, err
		}
		// If positive outcome
		if response {
			//Println("Found target on left left side")
			return true, ops3, errSuccess
		}
		//Println("- Test left right")
		response, ops3, err = searchDownward(originNode, startNode.Left.Right, startNode.Left.Right, targetNode, ops2)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops3, err
		}
		// If positive outcome
		if response {
			//Println("Found target on left right side")
			return true, ops3, errSuccess
		}
	}
	// Test right (will only be done if left was not successful)
	// Right node must not be nil, and not be the last explored node (i.e., a right child of the currently explored one)
	if startNode.Right != nil && startNode.Right != lastNode {
		//Println("Testing right branch ...")
		response, ops2, err := searchDownward(originNode, startNode.Right, startNode.Right, targetNode, ops)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops2, err
		}
		// If positive outcome
		if response {
			//Println("Found target on right side")
			return true, ops2, errSuccess
		}
		// Delegate downwards
		//Println("- Test right left")
		response, ops3, err := searchDownward(originNode, startNode.Right.Left, startNode.Right.Left, targetNode, ops2)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops3, err
		}
		// If positive outcome
		if response {
			//Println("Found target on right left side")
			return true, ops3, errSuccess
		}
		//Println("- Test right right")
		response, ops3, err = searchDownward(originNode, startNode.Right.Right, startNode.Right.Right, targetNode, ops2)
		// return lacking success if appearing
		if err.ErrorCode != TREE_NO_ERROR {
			return false, ops3, err
		}
		// If positive outcome
		if response {
			Println("Found target on right right side")
			return true, ops3, errSuccess
		}
	}

	//Println("Final result: false")
	return false, ops, errSuccess
}

/*
Combines existing nodes into new node and returns newly generated node.
Returns an error if component types of input nodes differ (should not be combined).
Returns tree.TREE_NO_ERROR in case of success.
*/
func Combine(leftNode *Node, rightNode *Node, logicalOperator string) (*Node, NodeError) {

	if leftNode == nil && rightNode == nil {
		log.Fatal("Illegal call to Combine() with nil nodes")
	}
	if leftNode == nil || leftNode.IsEmptyOrNilNode() {
		Println("Combining nodes returns right node (other node is nil or empty)")
		return rightNode, NodeError{ErrorCode: TREE_NO_ERROR}
	}
	if rightNode == nil || rightNode.IsEmptyOrNilNode() {
		Println("Combining nodes returns left node (other node is nil or empty)")
		return leftNode, NodeError{ErrorCode: TREE_NO_ERROR}
	}
	// In all other cases, create new combination using provided logical operator
	newNode := Node{}
	newNode.Left = leftNode
	newNode.Left.Parent = &newNode
	newNode.Right = rightNode
	newNode.Right.Parent = &newNode
	newNode.LogicalOperator = logicalOperator
	// Move left nodes component name to newly created parent node
	if leftNode.GetComponentName() != "" {
		newNode.ComponentType = leftNode.GetComponentName()
	}
	// Attach right node's type to parent if none is provided
	if leftNode.GetComponentName() == "" && rightNode.GetComponentName() != "" {
		newNode.ComponentType = rightNode.GetComponentName()
	}
	// Check whether both nodes have divergent component names - should not be allowed.
	if leftNode.GetComponentName() != "" && rightNode.GetComponentName() != "" && leftNode.GetComponentName() != rightNode.GetComponentName() {
		return nil, NodeError{TREE_INVALID_COMPONENT_COMBINATIONS, "Invalid component types for nodes to be combined (Left: " +
			leftNode.GetComponentName() + ", Right: " + rightNode.GetComponentName() + ")", nil}
	}
	return &newNode, NodeError{ErrorCode: TREE_NO_ERROR}
}

/*
*

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
func (n *Node) Validate() (bool, NodeError) {
	if n.Entry == nil && (n.Left == nil || n.Right == nil) {
		errorMsg := "Non-leaf node, but missing specification of left and right child, " +
			"or both. Node: " + fmt.Sprint(n.String())
		return false, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: errorMsg}
	}
	if n.Entry != nil && (n.Left != nil || n.Right != nil) {
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
Calculates state complexity for given node and returns the result.
Number of options on which the calculation is based can be retrieved using CountLeaves().
*/
func (n *Node) CalculateStateComplexity() (int, NodeError) {
	if n == nil {
		// Uninitialized node
		return 0, NodeError{ErrorCode: TREE_NO_ERROR}
	}
	if n.Left == nil && n.Right == nil && n.Entry == "" {
		// Must be empty node
		return 0, NodeError{ErrorCode: TREE_NO_ERROR}
	}
	if n.Left == nil && n.Right == nil && n.HasPrimitiveEntry() && n.Entry != "" {
		// Must be single leaf node (entry)
		return 1, NodeError{ErrorCode: TREE_NO_ERROR}
	}
	if n.Left == nil && n.Right == nil && !n.HasPrimitiveEntry() {
		// Must be nested statement, so delegate execution to statement, ...
		if reflect.TypeOf(n.Entry) == reflect.TypeOf(&Statement{}) {
			stmt := n.Entry.(*Statement)
			return stmt.CalculateComplexity().TotalStateComplexity, NodeError{ErrorCode: TREE_NO_ERROR}
		}
		// ... or component pair statement (i.e., extrapolated into multiple statements)
		if reflect.TypeOf(n.Entry) == reflect.TypeOf(&Node{}) {
			stmt := n.Entry.([]*Node)[0].Entry.(*Statement)
			return stmt.CalculateComplexity().TotalStateComplexity, NodeError{ErrorCode: TREE_NO_ERROR}
		}
		// Invalid type - not yet handled
		return -1, NodeError{ErrorCode: PARSING_ERROR_INVALID_TYPE_COMPLEXITY_CALCULATION,
			ErrorMessage: "No complexity calculation possible for type " + reflect.TypeOf(n.Entry).String()}
	}
	// Check if nested elements contain complexity
	if n.Left != nil && n.Right != nil {
		leftComplexity, err1 := n.Left.CalculateStateComplexity()
		if err1.ErrorCode != TREE_NO_ERROR {
			return -1, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: "Invalid tree structure in left nested tree: " + err1.ErrorMessage}
		}
		rightComplexity, err2 := n.Right.CalculateStateComplexity()
		if err2.ErrorCode != TREE_NO_ERROR {
			return -1, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: "Invalid tree structure in right nested tree: " + err2.ErrorMessage}
		}

		// Check for operator - note that AND can be explicit, or implicit (bAND, wAND)
		if n.LogicalOperator == AND || n.LogicalOperator == SAND_BETWEEN_COMPONENTS || n.LogicalOperator == SAND_WITHIN_COMPONENTS {
			// Discount for AND on this level, and return
			return leftComplexity + rightComplexity - 1, NodeError{ErrorCode: TREE_NO_ERROR}
		} else if n.LogicalOperator == XOR {
			// Return sum of left and right complexity (since either or applies)
			return leftComplexity + rightComplexity, NodeError{ErrorCode: TREE_NO_ERROR}
		} else if n.LogicalOperator == OR {
			// Return sum of both, alongside additional state (their combined applicability)
			return leftComplexity + rightComplexity + 1, NodeError{ErrorCode: TREE_NO_ERROR}
		}
	}
	return -1, NodeError{ErrorCode: TREE_INVALID_TREE, ErrorMessage: "Invalid tree structure absence of left or right leaf in combination."}
}

/*
   Returns root node of given tree the node is embedded in
   up to the level at which nodes are linked by synthetic AND (bAND and wAND).
   I.e., it returns the last node level below an sAND or bAND linkage.
*/
// TODO: Check for the need to refine considerations of SAND_WITHIN_COMPONENTS
func (n *Node) GetNodeBelowSyntheticRootNode() *Node {
	if n.Parent == nil || n.Parent.LogicalOperator == SAND_BETWEEN_COMPONENTS || n.Parent.LogicalOperator == SAND_WITHIN_COMPONENTS {
		// Assume to be parent if no parent on its own,
		// or root in synthetic hierarchy if paired with sAND
		return n
	} else {
		// else delegate to parent
		return n.Parent.GetNodeBelowSyntheticRootNode()
	}
}

/*
Indicates whether a given node has a linkage (wAND) within the same component
(e.g., Cex(shared (left [AND] right) middle (left2 [XOR] right2) shared)).
*/
func (n *Node) HasWithinComponentLinkage() bool {
	if n.Parent != nil && n.Parent.LogicalOperator == SAND_WITHIN_COMPONENTS {
		// Has linkage in parent
		return true
	} else if n.Parent != nil {
		// Delegate to parent
		return n.Parent.HasWithinComponentLinkage()
	} else {
		// if no parent, then no within-linkage
		return false
	}
}

/*
Returns all nodes that are in the same branch under a within-component linkage (i.e., wAND operators).
*/
func (n *Node) getAllNodesInWithinComponentLinkageBranch() []*Node {
	if n.Parent != nil && n.Parent.LogicalOperator == SAND_WITHIN_COMPONENTS {
		// Return own node if immediate parent is wAND linkage
		return []*Node{n}
	} else if n.Parent != nil {
		return n.Parent.getAllNodesInWithinComponentLinkageBranch()
	}
	// Return empty node
	return []*Node{}
}

/*
Returns root node of given tree independent of linking logical operators.
Internally it checks the chain of parent relationships. If no parent exists, the node itself is returned.
*/
func (n *Node) GetRootNode() *Node {
	if n.Parent == nil {
		// Assume to be parent if no parent on its own,
		// or root in synthetic hierarchy if paired with sAND
		return n
	} else {
		// else delegate to parent
		return n.Parent.GetRootNode()
	}
}

/*
Returns leaf nodes of a given node as arrays of arrays of nodes.
The two-dimensional array allows for separate storage of multiple arrays for a given component (e.g., multiple attributes, aims, etc.).
The parameter aggregateImplicitLinkages indicates whether the nodes for a given tree with implicitly linked branches
should be returned as a single tree, or multiple trees.
*/
func (n *Node) GetLeafNodes(aggregateImplicitLinkages bool) [][]*Node {
	return n.GetLeafNodesWithoutGivenNode(aggregateImplicitLinkages, nil)
}

/*
Returns leaf nodes of a given node as arrays of arrays of nodes, while ignoring a given node in the output (nil, if none to be ignored).
The two-dimensional array allows for separate storage of multiple arrays for a given component (e.g., multiple attributes, aims, etc.).
The parameter aggregateImplicitLinkages indicates whether the nodes for a given tree with implicitly linked branches
should be returned as a single tree, or multiple trees.
*/
func (n *Node) GetLeafNodesWithoutGivenNode(aggregateImplicitLinkages bool, nodeToBeIgnored *Node) [][]*Node {

	// Error checking first
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
		if n != nodeToBeIgnored {
			// Check whether to be ignored - only add if not be ignored
			returnNode = append(returnNode, inner)
		}
		return returnNode
	}

	// Output 2-dim arrays
	leftNodes := [][]*Node{}
	rightNodes := [][]*Node{}

	// If both left and right children nodes exist, return those combined
	if n.Left != nil && n.Right != nil {
		aggregate := 1
		if n.LogicalOperator == SAND_BETWEEN_COMPONENTS {
			Println("Found " + SAND_BETWEEN_COMPONENTS)
			if aggregateImplicitLinkages {
				aggregate = 0
			}
		} else if n.LogicalOperator == SAND_WITHIN_COMPONENTS {
			Println("Found " + SAND_WITHIN_COMPONENTS)
			// Generate combinations of internally linked component elements (wAND)
			aggregate = 1
		} else {
			Println("Found operator:", n.LogicalOperator)
			// Regular operators (i.e., AND, OR, XOR)
			aggregate = 0
		}

		res := aggregateNodes(aggregate, n.Left.GetLeafNodesWithoutGivenNode(aggregateImplicitLinkages, nodeToBeIgnored),
			n.Right.GetLeafNodesWithoutGivenNode(aggregateImplicitLinkages, nodeToBeIgnored), returnNode)
		Println("Returning following aggregated nodes for operator", n.LogicalOperator, ": ", res)
		return res
	}

	// Alternatively, process nodes individually

	// Process left nodes
	if n.Left != nil {
		leftNodes = n.Left.GetLeafNodes(aggregateImplicitLinkages)
		for _, v := range leftNodes {
			returnNode[0] = append(returnNode[0], v...)
		}
	}
	// Process right nodes
	if n.Right != nil {
		rightNodes = n.Right.GetLeafNodes(aggregateImplicitLinkages)
		for _, v := range rightNodes {
			returnNode[0] = append(returnNode[0], v...)
		}
	}
	return returnNode
}

/*
Returns top-level nodes containing statements. If none is found, returns nil.
*/
func (n *Node) GetTopLevelStatementNodes() []*Node {

	// Move to root first and work downwards
	root := n.GetRootNode()

	// Check that root is not empty
	if root.IsEmptyOrNilNode() {
		return nil
	}

	// Navigate downwards
	return findTopLevelStatementBelowNode(root, []*Node{})
}

/*
Retrieves top-level statements in nodes on or below given node.

Takes starting node as input as well as initialized collection for returned statements.
*/
func findTopLevelStatementBelowNode(node *Node, stmts []*Node) []*Node {

	// Test whether node itself is single and has statement embedded
	if node.Entry != nil && reflect.TypeOf(node.Entry) == reflect.TypeOf(&Statement{}) {
		// If root has statement in Entry, then return this as top-level statement
		return append(stmts, node)
	}

	// Iterate through nodes array (if present)
	if node.Entry != nil && reflect.TypeOf(node.Entry) == reflect.TypeOf([]*Node{}) {
		for _, v := range node.Entry.([]*Node) {
			stmts = findTopLevelStatementBelowNode(v, stmts)
		}
	}

	// test left side downwards
	if node.IsCombination() && !node.Left.IsEmptyOrNilNode() {
		if node.Left.IsCombination() || reflect.TypeOf(node.Left.Entry) != reflect.TypeOf([]*Statement{}) {
			stmts = findTopLevelStatementBelowNode(node.Left, stmts)
		} else if reflect.TypeOf(node.Left.Entry) == reflect.TypeOf(&Statement{}) {
			stmts = append(stmts, node.Left)
		}
	}

	// test right side downwards
	if node.IsCombination() && !node.Right.IsEmptyOrNilNode() {
		if node.Right.IsCombination() || reflect.TypeOf(node.Right.Entry) != reflect.TypeOf([]*Statement{}) {
			stmts = findTopLevelStatementBelowNode(node.Right, stmts)
		} else if reflect.TypeOf(node.Right.Entry) == reflect.TypeOf(&Statement{}) {
			stmts = append(stmts, node.Right)
		}
	}

	return stmts
}

/*
Enables different forms of node aggregation, where aggregationType 0 indicates flat combination of nodes in array ([ ..., one, two, ...]),
and aggregationType 1 indicates returning node arrays within node array ([ ..., [one, two], ...])
Takes populated leaf arrays as input and prepared return structure for output.
*/
func aggregateNodes(aggregationType int, leftNodes [][]*Node, rightNodes [][]*Node, returnNode [][]*Node) [][]*Node {
	switch aggregationType {
	case 0:
		// return as flat structure (i.e., individual nodes are returned in isolation)
		// (e.g., [ ..., one, two, ... ])
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
	case 1:
		// Append individual nested node arrays
		// (e.g., [ ... [one, two] ... ])
		for _, v := range leftNodes {
			returnNode = append(returnNode, v)
		}
		// Append individual node arrays
		for _, v := range rightNodes {
			returnNode = append(returnNode, v)
		}
		return returnNode
	}
	return returnNode
}

/*
Substitutes node references within a tree structure for given start node (tree) based on *downward* search
of branches (left, right).
Allows for indication whether new node should be linked to parent, and whether children from original node should be inherited.
*/
func substituteNodeReferenceInTree(originalNode *Node, newNode *Node, tree *Node, linkNewNodeToParent bool, inheritChildren bool) NodeError {
	if originalNode == nil {
		return NodeError{ErrorCode: TREE_ERROR_NIL_NODE, ErrorMessage: "Node to be substituted is nil."}
	}
	if newNode == nil {
		return NodeError{ErrorCode: TREE_ERROR_NIL_NODE, ErrorMessage: "Node to substitute existing node is nil."}
	}
	if tree == nil {
		return NodeError{ErrorCode: TREE_ERROR_NIL_NODE, ErrorMessage: "Node tree to perform substitution on is nil."}
	}

	// Search downwards
	if tree.Left != nil {
		if tree.Left == originalNode {
			if inheritChildren {
				// Inherit children on left side and relink parent
				if tree.Left.Left != nil {
					newNode.Left = tree.Left.Left
					tree.Left.Left.Parent = newNode
				}
				if tree.Left.Right != nil {
					newNode.Right = tree.Left.Right
					tree.Left.Right.Parent = newNode
				}
			}
			if linkNewNodeToParent {
				// Link to original node's parent
				newNode.Parent = tree
			}
			// Replace left child
			tree.Left = newNode
		} else {
			// if not nil, but left element, delegate down
			substituteNodeReferenceInTree(originalNode, newNode, tree.Left, linkNewNodeToParent, inheritChildren)
		}
	}
	if tree.Right != nil {
		if tree.Right == originalNode {
			if inheritChildren {
				// Inherit children on right side and relink parent
				if tree.Right.Left != nil {
					newNode.Left = tree.Right.Left
					tree.Right.Left.Parent = newNode
				}
				if tree.Right.Right != nil {
					newNode.Right = tree.Right.Right
					tree.Right.Right.Parent = newNode
				}
			}
			if linkNewNodeToParent {
				// Link to original node's parent
				newNode.Parent = tree
			}
			// Replace right child
			tree.Right = newNode
		} else {
			// if not nil, but right element, delegate down
			substituteNodeReferenceInTree(originalNode, newNode, tree.Right, linkNewNodeToParent, inheritChildren)
		}
	}
	return NodeError{ErrorCode: TREE_NO_ERROR}
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
Indicates whether node is empty (i.e., has "zero values" in fields) or is nil (implies empty node).
*/
func (n *Node) IsEmptyOrNilNode() bool {

	if n.GetNodeState() == NODE_STATE_NON_NIL {
		return !reflect.ValueOf(n.Entry).IsValid() &&
			n.LogicalOperator == "" &&
			n.ComponentType == "" &&
			n.Parent == nil &&
			n.Left == nil &&
			n.Right == nil &&
			n.Suffix == nil &&
			n.Annotations == nil &&
			n.SharedLeft == nil &&
			n.SharedRight == nil &&
			n.PrivateNodeLinks == nil &&
			n.ElementOrder == nil
	}
	// Node is nil, and hence implied empty
	return true
}

/*
Indicates if node is nil.
*/
func (n *Node) IsNil() bool {
	return n == nil
}

/*
Applies function to all entries below a given node and adds statements to their entries based on parsed input.
*/
func (n *Node) ParseAllEntries(function func(string) (*Statement, ParsingError)) ParsingError {
	if n.IsNil() {
		return ParsingError{ErrorCode: PARSING_ERROR_NIL_ELEMENT, ErrorMessage: "Attempted to parse nil element."}
	}
	if !n.IsEmptyOrNilNode() && n.Entry != nil {

		// Execute actual function
		newEntry, err := function(n.Entry.(string))
		if err.ErrorCode != PARSING_NO_ERROR {
			Println("Received error when parsing nested string to statement.")
			return err
		}
		// and reassign parsed element (i.e., substitute previous string element with node representation of statement)
		n.Entry = newEntry
	}
	if !n.Left.IsNil() {
		// Parse left child of combination
		err := n.Left.ParseAllEntries(function)
		if err.ErrorCode != PARSING_NO_ERROR {
			Println("Received error when parsing left-hand nested statement.")
			return err
		}
	}
	if !n.Right.IsNil() {
		// Parse right child of combination
		err := n.Right.ParseAllEntries(function)
		if err.ErrorCode != PARSING_NO_ERROR {
			Println("Received error when parsing right-hand nested statement.")
			return err
		}
	}
	return ParsingError{ErrorCode: PARSING_NO_ERROR}
}

/*
Returns suffix of given node, or, if nested, of parent node if it has logical operator.
Inheritance only works up to node with logical operator and populated suffix
(i.e., will walk upward until finding node with logical operator and suffix).
Returns empty string if no suffix present.
*/
func (n *Node) GetSuffix() string {
	if n.Suffix != nil {
		// Return suffix
		return n.Suffix.(string)
	} else {
		// Delegate to parent
		if n.Parent != nil && n.Parent.hasLogicalOperator() {
			return n.Parent.GetSuffix()
		} else if n.Parent != nil && !n.Parent.hasLogicalOperator() {
			fmt.Println("Parent node is missing logical operator (and hence child nodes do not inherit suffix)")
		}
	}
	// Default suffix value
	return ""
}

/*
Indicates whether node has own private nodes (referenced via PrivateNodes field)
*/
func (n *Node) HasPrivateNodes() bool {
	if n.PrivateNodeLinks != nil && len(n.PrivateNodeLinks) > 0 {
		return true
	} else {
		return false
	}
}

/*
Indicates whether node has annotations (referenced via Annotations field)
*/
func (n *Node) HasAnnotations() bool {
	if n.GetAnnotations() != nil && n.GetAnnotations() != nil && len(n.GetAnnotations().(string)) > 0 {
		return true
	} else {
		return false
	}
}

/*
Returns annotations for specific node. If non-synthetic parent nodes hold annotations,
those are inherited; if both parent and child annotations exist, both are combined (parent + child).
*/
func (n *Node) GetAnnotations() interface{} {

	// Return nil if node is nil
	if n == nil {
		return nil
	}

	// Otherwise, retrieve annotations

	// Retrieve parent annotations
	var parentAnnotations interface{}
	if n.Parent != nil && n.Parent.LogicalOperator != SAND_BETWEEN_COMPONENTS {
		// Delegate to parent component
		parentAnnotations = n.Parent.GetAnnotations()
	}
	// Retrieve child annotations
	var childAnnotations interface{}
	if n.Entry != nil && reflect.TypeOf(n.Entry) == reflect.TypeOf([]*Node{}) && len(n.Entry.([]*Node)) > 0 {
		childAnnotations = n.Entry.([]*Node)[0].GetAnnotations()
	}

	// If child annotations are populated, return only those
	if childAnnotations != nil {
		// Return child annotations (including recursively inherited annotations) only
		return childAnnotations
	}

	// else process all other variations ...

	// Return combined parent and child annotations if existing, ...
	if parentAnnotations != nil && n.Annotations != nil {
		return parentAnnotations.(string) + n.Annotations.(string)
	} else if n.Annotations != nil {
		// ... else only child annotations, ...
		return n.Annotations.(string)
	} else if parentAnnotations != nil {
		// ... or only parent annotations
		return parentAnnotations.(string)
	}
	return nil
}

/*
Helper function to pretty-print node arrays
*/
func PrintNodes(nodes []*Node) string {
	outString := "\n"
	for i, v := range nodes {
		outString += "Node " + strconv.Itoa(i) + ": \n"
		outString += "Annotations: " + fmt.Sprintf("%v", v.Annotations) + "\n"
		outString += v.String() + "\n"
	}
	return outString
}
