package exporter

import "IG-Parser/tree"

/*
Indicates whether shared elements are included in output
 */
var INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

/*
Indicates whether tabular export produces dynamic output based on present components
(better for individual statements), or produces fixed predefined structure (better for datasets)
*/
var CREATE_DYNAMIC_TABULAR_OUTPUT = false

/*
Returns a fixed schema for tabular output
 */
func GetStaticTabularOutputSchema() map[string]int {

	// Generate static headers
	staticComponentFrequency := make(map[string]int)

	// Regulative side
	staticComponentFrequency[tree.ATTRIBUTES_PROPERTY] = 1
	staticComponentFrequency[tree.ATTRIBUTES_PROPERTY_REFERENCE] = 1
	staticComponentFrequency[tree.ATTRIBUTES] = 1
	staticComponentFrequency[tree.DEONTIC] = 1
	staticComponentFrequency[tree.AIM] = 1
	staticComponentFrequency[tree.DIRECT_OBJECT_PROPERTY] = 1
	staticComponentFrequency[tree.DIRECT_OBJECT_PROPERTY_REFERENCE] = 1
	staticComponentFrequency[tree.DIRECT_OBJECT] = 1
	staticComponentFrequency[tree.DIRECT_OBJECT_REFERENCE] = 1
	staticComponentFrequency[tree.INDIRECT_OBJECT_PROPERTY] = 1
	staticComponentFrequency[tree.INDIRECT_OBJECT_PROPERTY_REFERENCE] = 1
	staticComponentFrequency[tree.INDIRECT_OBJECT] = 1
	staticComponentFrequency[tree.INDIRECT_OBJECT_REFERENCE] = 1

	// Shared elements
	staticComponentFrequency[tree.ACTIVATION_CONDITION] = 1
	staticComponentFrequency[tree.ACTIVATION_CONDITION_REFERENCE] = 1
	staticComponentFrequency[tree.EXECUTION_CONSTRAINT] = 1
	staticComponentFrequency[tree.EXECUTION_CONSTRAINT_REFERENCE] = 1

	// Constitutive side
	staticComponentFrequency[tree.CONSTITUTED_ENTITY_PROPERTY] = 1
	staticComponentFrequency[tree.CONSTITUTED_ENTITY_PROPERTY_REFERENCE] = 1
	staticComponentFrequency[tree.CONSTITUTED_ENTITY] = 1
	staticComponentFrequency[tree.MODAL] = 1
	staticComponentFrequency[tree.CONSTITUTIVE_FUNCTION] = 1
	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_PROPERTY] = 1
	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE] = 1
	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES] = 1
	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_REFERENCE] = 1

	return staticComponentFrequency
}