package tree

import "strconv"

const (
	ATTRIBUTES = "A"
	ATTRIBUTES_PROPERTY = "A,p"
	DEONTIC = "D"
	AIM = "I"
	DIRECT_OBJECT = "Bdir"
	DIRECT_OBJECT_PROPERTY = "Bdir,p"
	INDIRECT_OBJECT = "Bind"
	INDIRECT_OBJECT_PROPERTY = "Bind,p"
	ACTIVATION_CONDITION = "Cac"
	EXECUTION_CONSTRAINT = "Cex"
	CONSTITUTED_ENTITY = "E"
	CONSTITUTED_ENTITY_PROPERTY = "E,p"
	MODAL = "M"
	CONSTITUTIVE_FUNCTION = "F"
	CONSTITUTING_PROPERTIES = "P"
	CONSTITUTING_PROPERTIES_PROPERTY = "P,p"
	AND = "AND"
	OR = "OR"
	XOR = "XOR"
	NOT = "NOT"
	AND_BRACKETS = "[" + AND + "]"
	OR_BRACKETS = "[" + OR + "]"
	XOR_BRACKETS = "[" + XOR + "]"
	NOT_BRACKETS = "[" + NOT + "]"
	PARSING_MODE_LEFT = "PARSING_LEFT"
	PARSING_MODE_RIGHT = "PARSING_RIGHT"
	PARSING_MODE_OUTSIDE_EXPRESSION = "PARSING_OUTSIDE"
)

/*
type igComponent struct {
	ComponentName string
}

func IGComponent(componentName string) *igComponent {
	i := igComponent{}
	i.ComponentName = componentName
	return &i
}

func (i igComponent) String() string {
	return i.ComponentName
}
/*
Checks whether component value is valid (i.e., a valid IG Component symbol).
 */
/*func (c *igComponent) valid() bool {
	return StringInSlice(c.ComponentName, igComponents)
}*/

/*
IG 2.0 Component Symbols
 */
var IGComponents = []string{
	ATTRIBUTES,
	ATTRIBUTES_PROPERTY,
	DEONTIC,
	AIM,
	DIRECT_OBJECT,
	DIRECT_OBJECT_PROPERTY,
	INDIRECT_OBJECT,
	INDIRECT_OBJECT_PROPERTY,
	ACTIVATION_CONDITION,
	EXECUTION_CONSTRAINT,
	CONSTITUTED_ENTITY,
	CONSTITUTED_ENTITY_PROPERTY,
	MODAL,
	CONSTITUTIVE_FUNCTION,
	CONSTITUTING_PROPERTIES,
	CONSTITUTING_PROPERTIES_PROPERTY}

type igLogicalOperator struct {
	LogicalOperatorName string
}

/*
Checks whether operator value is valid (i.e., a valid logical operator symbol).
*/
func (o *igLogicalOperator) valid() bool {
	return StringInSlice(o.LogicalOperatorName, IGLogicalOperators)
}

func (o igLogicalOperator) String() string {
	return o.LogicalOperatorName
}

/*
Valid logical operators in IG 2.0
 */
var IGLogicalOperators = []string{
	AND,
	OR,
	XOR,
	NOT,
}

/**
Signals boundary values for detected combinations
 */
type Boundaries struct {
	Left int
	Operator int
	OperatorVal string
	Right int
	Complete bool
}

func (b *Boundaries) String() string {
	return "Boundaries{\n"+
		"  Left: " + strconv.Itoa(b.Left) + "\n" +
		"  Op Pos: " + strconv.Itoa(b.Operator) + "\n" +
		"  Operator: " + b.OperatorVal + "\n" +
		"  Right: " + strconv.Itoa(b.Right) + "\n}"
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
