package exporter

import (
	"IG-Parser/tree"
	"log"
)

/*
Indicates whether shared elements are included in output
 */
var INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

/*
Indicates whether tabular export produces dynamic output based on present components
(better for individual statements), or produces fixed predefined structure (better for datasets)
Should not be directly modified, but rather using SetDynamicOutput().
*/
var create_DYNAMIC_TABULAR_OUTPUT = false

/*
Sets whether produced output should be dynamic or static.
 */
func SetDynamicOutput(dynamic bool) {
	create_DYNAMIC_TABULAR_OUTPUT = dynamic
	if create_DYNAMIC_TABULAR_OUTPUT {
		log.Println("Activated dynamic output.")
		tree.AGGREGATE_IMPLICIT_LINKAGES = false
	} else {
		log.Println("Activated static output.")
		tree.AGGREGATE_IMPLICIT_LINKAGES = true
	}
}

/*
Queries whether dynamic (vs. static) output is activated.
 */
func ProduceDynamicOutput() bool {
	return create_DYNAMIC_TABULAR_OUTPUT
}

/*
Returns a fixed schema for tabular output
 */
func GetStaticTabularOutputSchema() map[string]int {

	// Generate static headers
	staticComponentFrequency := make(map[string]int)

	// Regulative side
	staticComponentFrequency[tree.ATTRIBUTES_PROPERTY] = 1
	staticComponentFrequency[tree.ATTRIBUTES_PROPERTY_REFERENCE] = 1
	staticComponentFrequency[tree.ATTRIBUTES_PROPERTY_ANNOTATION] = 1

	staticComponentFrequency[tree.ATTRIBUTES] = 1
	staticComponentFrequency[tree.ATTRIBUTES_ANNOTATION] = 1

	staticComponentFrequency[tree.DEONTIC] = 1
	staticComponentFrequency[tree.DEONTIC_ANNOTATION] = 1

	staticComponentFrequency[tree.AIM] = 1
	staticComponentFrequency[tree.AIM_ANNOTATION] = 1

	staticComponentFrequency[tree.DIRECT_OBJECT_PROPERTY] = 1
	staticComponentFrequency[tree.DIRECT_OBJECT_PROPERTY_REFERENCE] = 1
	staticComponentFrequency[tree.DIRECT_OBJECT_PROPERTY_ANNOTATION] = 1

	staticComponentFrequency[tree.DIRECT_OBJECT] = 1
	staticComponentFrequency[tree.DIRECT_OBJECT_REFERENCE] = 1
	staticComponentFrequency[tree.DIRECT_OBJECT_ANNOTATION] = 1

	staticComponentFrequency[tree.INDIRECT_OBJECT_PROPERTY] = 1
	staticComponentFrequency[tree.INDIRECT_OBJECT_PROPERTY_REFERENCE] = 1
	staticComponentFrequency[tree.INDIRECT_OBJECT_PROPERTY_ANNOTATION] = 1

	staticComponentFrequency[tree.INDIRECT_OBJECT] = 1
	staticComponentFrequency[tree.INDIRECT_OBJECT_REFERENCE] = 1
	staticComponentFrequency[tree.INDIRECT_OBJECT_ANNOTATION] = 1

	// Shared elements
	staticComponentFrequency[tree.ACTIVATION_CONDITION] = 1
	staticComponentFrequency[tree.ACTIVATION_CONDITION_REFERENCE] = 1
	staticComponentFrequency[tree.ACTIVATION_CONDITION_ANNOTATION] = 1

	staticComponentFrequency[tree.EXECUTION_CONSTRAINT] = 1
	staticComponentFrequency[tree.EXECUTION_CONSTRAINT_REFERENCE] = 1
	staticComponentFrequency[tree.EXECUTION_CONSTRAINT_ANNOTATION] = 1

	// Constitutive side
	staticComponentFrequency[tree.CONSTITUTED_ENTITY_PROPERTY] = 1
	staticComponentFrequency[tree.CONSTITUTED_ENTITY_PROPERTY_REFERENCE] = 1
	staticComponentFrequency[tree.CONSTITUTED_ENTITY_PROPERTY_ANNOTATION] = 1

	staticComponentFrequency[tree.CONSTITUTED_ENTITY] = 1
	staticComponentFrequency[tree.CONSTITUTED_ENTITY_ANNOTATION] = 1

	staticComponentFrequency[tree.MODAL] = 1
	staticComponentFrequency[tree.MODAL_ANNOTATION] = 1

	staticComponentFrequency[tree.CONSTITUTIVE_FUNCTION] = 1
	staticComponentFrequency[tree.CONSTITUTIVE_FUNCTION_ANNOTATION] = 1

	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_PROPERTY] = 1
	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE] = 1
	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_PROPERTY_ANNOTATION] = 1

	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES] = 1
	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_REFERENCE] = 1
	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_ANNOTATION] = 1

	return staticComponentFrequency
}