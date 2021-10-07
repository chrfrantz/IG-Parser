package tree

/*
Indicates inheritance mode of shared elements from parent to children nodes (see constants below)
*/
var SHARED_ELEMENT_INHERITANCE_MODE = SHARED_ELEMENT_INHERIT_APPEND

/*
Indicates whether implicitly linked components (e.g., multiple Cex annotations) should be aggregated
into single tree, or be returned as separate trees (based on the implicit AND operator)
*/
var AGGREGATE_IMPLICIT_LINKAGES = true

/*
 Disables inheriting shared elements from parent nodes, not even combination nodes embedding the entity.
 It would only return elements associated with the node itself.
 */
const SHARED_ELEMENT_INHERIT_NOTHING = "SHARED_ELEMENT_INHERIT_NOTHING"

/*
 Limits inheritance to next higher logical combination node (i.e., the shared elements of the node embedding the
 referenced one and its potential sibling node; as well as potential shared elements of the leaf node itself)
 Example: in (shared left (leafLeft [AND] leafRight)), leafLeft would inherit "shared left")
 */
const SHARED_ELEMENT_INHERIT_FROM_COMBINATION = "SHARED_ELEMENT_INHERIT_FROM_COMBINATION"

/*
 Indicates that child nodes inherit parent node's shared elements by overwriting child
 elements in child node with parent values if parent values are non-empty
*/
const SHARED_ELEMENT_INHERIT_OVERRIDE = "SHARED_ELEMENT_INHERIT_OVERRIDE"
/*
 Indicates that child nodes inherit parent node's shared elements by appending child
 elements to parent elements in child node
*/
const SHARED_ELEMENT_INHERIT_APPEND = "SHARED_ELEMENT_INHERIT_APPEND"
