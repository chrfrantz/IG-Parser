package tree

/*
Indicates inheritance mode of shared elements from parent to children nodes (see constants below)
*/
var SHARED_ELEMENT_INHERITANCE_MODE = SHARED_ELEMENT_INHERIT_APPEND

// Indicates separator used when appending elements to components (e.g., inheriting from parent node)
const INHERITANCE_DELIMITER = ","

// Disables inheriting shared elements from parent nodes
const SHARED_ELEMENT_INHERIT_NOTHING = "SHARED_ELEMENT_INHERIT_NOTHING"
/*
 Indicates that child nodes inherit parent node's shared elements if both have AND operators by overwriting child
 elements in child node with parent values if parent values are non-empty
*/
const SHARED_ELEMENT_INHERIT_OVERRIDE = "SHARED_ELEMENT_INHERIT_OVERRIDE"
/*
 Indicates that child nodes inherit parent node's shared elements if both have AND operators by appending child
 elements to parent elements in child node
*/
const SHARED_ELEMENT_INHERIT_APPEND = "SHARED_ELEMENT_INHERIT_APPEND"
