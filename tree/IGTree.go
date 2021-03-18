package tree

import (
	"log"
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
	if s.Attributes != nil {
		out += s.Attributes.String()
		out += " "
	}
	if s.AttributesProperty != nil {
		out += s.AttributesProperty.String()
		out += " "
	}
	if s.Deontic != nil {
		out += s.Deontic.String()
		out += " "
	}
	if s.Aim != nil {
		out += s.Aim.String()
		out += " "
	}
	if s.DirectObject != nil {
		out += s.DirectObject.String()
		out += " "
	}
	if s.DirectObjectProperty != nil {
		out += s.DirectObjectProperty.String()
		out += " "
	}
	if s.IndirectObject != nil {
		out += s.IndirectObject.String()
		out += " "
	}
	if s.IndirectObjectProperty != nil {
		out += s.IndirectObjectProperty.String()
		out += " "
	}
	if s.ConstitutedEntity != nil {
		out += s.ConstitutedEntity.String()
		out += " "
	}
	if s.ConstitutedEntityProperty != nil {
		out += s.ConstitutedEntityProperty.String()
		out += " "
	}
	if s.Modal != nil {
		out += s.Modal.String()
		out += " "
	}
	if s.ConstitutiveFunction != nil {
		out += s.ConstitutiveFunction.String()
		out += " "
	}
	if s.ConstitutingProperties != nil {
		out += s.ConstitutingProperties.String()
		out += " "
	}
	if s.ConstitutingPropertiesProperty != nil {
		out += s.ConstitutingPropertiesProperty.String()
		out += " "
	}
	if s.ActivationConditionSimple != nil {
		out += s.ActivationConditionSimple.String()
		out += " "
	}
	if s.ActivationConditionComplex != nil {
		out += s.ActivationConditionComplex.String()
		out += " "
	}
	if s.ExecutionConstraintSimple != nil {
		out += s.ExecutionConstraintSimple.String()
		out += " "
	}
	if s.ExecutionConstraintComplex != nil {
		out += s.ExecutionConstraintComplex.String()
		out += " "
	}
	if s.OrElse != nil {
		out += s.OrElse.String()
		out += " "
	}
	return out
}

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
	SharedLeft string
	// Text shared across children to the right of a combination (e.g., ((left val [AND] right val) shared info))
	SharedRight string
	// Logical operator that links left and right values/nodes
	LogicalOperator string
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
	if n.SharedLeft != "" {
		out += "(" + n.SharedLeft + " "
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
	if n.SharedRight != "" {
		out += " " + n.SharedRight + ")"
	}
	return out
}

/*
Prints node content
 */
func (n *Node) String() string {
	/*if n.Parent != nil {
		for n.Parent
	}*/
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

		if n.SharedLeft != "" {
			out += prefix + "Shared (left): " + n.SharedLeft + "\n"
		}
		if n.SharedRight != "" {
			out += prefix + "Shared (right): " + n.SharedRight + "\n"
		}

		return "(\n" + out +
			prefix + "Left: " + n.Left.String() + "\n" +
			prefix + "Operator: " + n.LogicalOperator + "\n" +
			prefix + "Right: " + n.Right.String() + "\n" +
			prefix + ")"
	}
}

/*
Inserts a leaf node under a given node and inherits its component type
 */
func (n *Node) InsertLeafNode(entry string) *Node {
	return n.InsertChildNode(entry, "", "", n.ComponentType, "", "", "")
}

func ComponentLeafNode(entry string, componentType string) *Node {
	return ComponentNode(entry, "", "", componentType, "", "", "")
}

func ComponentNode(entry string, leftValue string, rightValue string, componentType string, sharedValueLeft string, sharedValueRight string, logicalOperator string) *Node {

	// Validation (Entry cannot be mixed with the other fields)
	if entry != "" {
		if leftValue != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as left-hand node (" + leftValue + ").")
		}
		if rightValue != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as right-hand node (" + rightValue + ").")
		}
		if sharedValueLeft != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as (left) shared content field (" + sharedValueLeft + ").")
		}
		if sharedValueRight != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as (right) shared content field (" + sharedValueRight + ").")
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

	/*if parent != nil {
		fmt.Println("Node is non-root node.")
		node.Parent = parent
	} else {
		fmt.Println("Node is root node.")
	}*/

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
func (n *Node) InsertChildNode(entry string, leftValue string, rightValue string, componentType string, sharedValueLeft string, sharedValueRight string, logicalOperator string) *Node {
	// Validation (Entry cannot be mixed with the other fields)
	if entry != "" {
		if leftValue != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as left-hand node (" + leftValue + ").")
		}
		if rightValue != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as right-hand node (" + rightValue + ").")
		}
		if sharedValueLeft != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as (left) shared content field (" + sharedValueLeft + ").")
		}
		if sharedValueRight != "" {
			log.Fatal("Invalid node specification. Entry field is filled (" + entry + "), as well as (right) shared content field (" + sharedValueRight + ").")
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

	/*if parent != nil {
		fmt.Println("Node is non-root node.")
		node.Parent = parent
	} else {
		fmt.Println("Node is root node.")
	}*/

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