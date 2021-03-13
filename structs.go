package igTree


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
	return stringInSlice(c.ComponentName, igComponents)
}*/

/*
IG 2.0 Component Symbols
 */
var igComponents = []string{
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
	return stringInSlice(o.LogicalOperatorName, igLogicalOperators)
}

func (o igLogicalOperator) String() string {
	return o.LogicalOperatorName
}

/*
Valid logical operators in IG 2.0
 */
var igLogicalOperators = []string{
	AND,
	OR,
	XOR,
	NOT,
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
