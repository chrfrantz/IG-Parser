package tabular

import (
	"IG-Parser/core/tree"
	"log"
)

/*
Output generated for Google Sheets processing
*/
const OUTPUT_TYPE_GOOGLE_SHEETS = "Google Sheets"

/*
Output generated for generic CSV processing (e.g., Excel Text-to-Columns)
*/
const OUTPUT_TYPE_CSV = "CSV format"

/*
No flat output to be generated
*/
const OUTPUT_TYPE_NONE = "NONE"

/*
Output types available for output generation
*/
var OUTPUT_TYPES = []string{OUTPUT_TYPE_GOOGLE_SHEETS, OUTPUT_TYPE_CSV}

/*
Default tabular output type.
*/
var DEFAULT_OUTPUT_TYPES = OUTPUT_TYPE_GOOGLE_SHEETS

/*
Values for Original Statement inclusion in tabular output
*/
const ORIGINAL_STATEMENT_OUTPUT_NONE = "No inclusion of Original Statement in output (i.e., no additional column)"                                            // no Original Statement output (no additional column in output)
const ORIGINAL_STATEMENT_OUTPUT_FIRST_ENTRY = "Include Original Statement for first atomic statement only (i.e., in first row following optional header row)" // Original Statement output in first atomic entry only (implies additional column, but only single entry on first atomic statement)
const ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES = "Include Original Statement for each atomic statement (i.e., in each row)"                                      // Original Statement output in all atomic entries (implies additional column, and entry for each atomic statement)

/*
Original Statement inclusion options
*/
var ORIGINAL_STATEMENT_INCLUSION_OPTIONS = []string{ORIGINAL_STATEMENT_OUTPUT_NONE, ORIGINAL_STATEMENT_OUTPUT_FIRST_ENTRY, ORIGINAL_STATEMENT_OUTPUT_ALL_ENTRIES}

/*
Default value for Original Statement inclusion in tabular output
*/
var DEFAULT_ORIGINAL_STATEMENT_OUTPUT = ORIGINAL_STATEMENT_OUTPUT_NONE

/*
Values for IG Script inclusion in tabular output
*/
const IG_SCRIPT_OUTPUT_NONE = "No inclusion of IG Script coding in output (i.e., no additional column)"                                                       // no IG Script output (no additional column in output)
const IG_SCRIPT_OUTPUT_FIRST_ENTRY = "Include IG Script-encoded statement for first atomic statement only (i.e., in first row following optional header row)" // IG Script output in first atomic entry only (implies additional column, but only single entry on first atomic statement)
const IG_SCRIPT_OUTPUT_ALL_ENTRIES = "Include IG Script-encoded statement for each atomic statement (i.e., in each row)"                                      // IG Script output in all atomic entries (implies additional column, and entry for each atomic statement)

/*
IG Script inclusion options
*/
var IG_SCRIPT_INCLUSION_OPTIONS = []string{IG_SCRIPT_OUTPUT_NONE, IG_SCRIPT_OUTPUT_FIRST_ENTRY, IG_SCRIPT_OUTPUT_ALL_ENTRIES}

/*
Default value for IG Script inclusion in tabular output
*/
var DEFAULT_IG_SCRIPT_OUTPUT = IG_SCRIPT_OUTPUT_NONE

/*
Indicates whether shared elements are included in output
*/
var include_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true

/*
Sets indicator whether shared elements are included in tabular output.
*/
func SetIncludeSharedElementsInTabularOutput(includeSharedElements bool) {
	include_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = includeSharedElements
}

/*
Indicates whether shared elements are included in tabular output.
*/
func IncludeSharedElementsInTabularOutput() bool {
	return include_SHARED_ELEMENTS_IN_TABULAR_OUTPUT
}

/*
Indicates whether the coding follows the IG Core decomposition level (or IG Extended)
*/
var create_IG_EXTENDED_OUTPUT = true

/*
Indicates whether tabular export produces dynamic output based on present components
(better for individual statements), or produces fixed predefined structure (better for datasets)
Should not be directly modified, but rather using SetDynamicOutput().
*/
var create_DYNAMIC_TABULAR_OUTPUT = false

/*
Indicates whether annotations should be included in output.
Should not be directly modified, but rather using SetIncludeAnnotations().
*/
var include_ANNOTATIONS = false

/*
Indicates whether Degree of Variability should be included in output.
Should not be directly modified, but rather using SetIncludeDegreeOfVariability().
*/
var include_DEGREE_OF_VARIABILITY = false

/*
Indicates whether header row is to be included in output.
Should not be directly modified, but rather using SetIncludeHeaders().
*/
var include_HEADERS = true

/*
Indicates whether adjacent operators should be collapsed (right now AND, sAND and bAND).
Should not be directly modified, but rather using SetCollapseOperators().
*/
var collapse_OPERATORS = true

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
Defines whether annotations should be included in output.
*/
func SetIncludeAnnotations(include bool) {
	include_ANNOTATIONS = include
}

/*
Indicates whether annotations should be included in output.
*/
func IncludeAnnotations() bool {
	return include_ANNOTATIONS
}

/*
Defines whether Degree of Variability should be included in output.
*/
func SetIncludeDegreeOfVariability(include bool) {
	include_DEGREE_OF_VARIABILITY = include
}

/*
Indicates whether Degree of Variability should be included in output.
*/
func IncludeDegreeOfVariability() bool {
	return include_DEGREE_OF_VARIABILITY
}

/*
Defines whether header is included in tabular output.
*/
func SetIncludeHeaders(include bool) {
	include_HEADERS = include
}

/*
Indicates whether header row is included in tabular output.
*/
func IncludeHeader() bool {
	return include_HEADERS
}

/*
Sets whether operators should be collapsed.
*/
func SetCollapseOperators(collapse bool) {
	collapse_OPERATORS = collapse
}

/*
Indicates whether operators should be collapsed.
*/
func CollapseOperators() bool {
	return collapse_OPERATORS
}

/*
Sets whether output should be IG Extended (component-level nesting) or IG Core.
*/
func SetProduceIGExtendedOutput(extendedOutput bool) {
	create_IG_EXTENDED_OUTPUT = extendedOutput
}

/*
Indicates whether output included IG Extended features (specifically component-level nesting).
*/
func ProduceIGExtendedOutput() bool {
	return create_IG_EXTENDED_OUTPUT
}

/*
Returns a fixed schema for tabular output
*/
func GetStaticTabularOutputSchema() map[string]int {

	// Generate static headers
	staticComponentFrequency := make(map[string]int)

	if include_ANNOTATIONS {
		// Statement annotation
		staticComponentFrequency[tree.STATEMENT_ANNOTATION] = 1
	}

	// Regulative side
	staticComponentFrequency[tree.ATTRIBUTES_PROPERTY] = 1
	staticComponentFrequency[tree.ATTRIBUTES_PROPERTY_REFERENCE] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.ATTRIBUTES_PROPERTY_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.ATTRIBUTES] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.ATTRIBUTES_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.DEONTIC] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.DEONTIC_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.AIM] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.AIM_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.DIRECT_OBJECT_PROPERTY] = 1
	staticComponentFrequency[tree.DIRECT_OBJECT_PROPERTY_REFERENCE] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.DIRECT_OBJECT_PROPERTY_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.DIRECT_OBJECT] = 1
	staticComponentFrequency[tree.DIRECT_OBJECT_REFERENCE] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.DIRECT_OBJECT_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.INDIRECT_OBJECT_PROPERTY] = 1
	staticComponentFrequency[tree.INDIRECT_OBJECT_PROPERTY_REFERENCE] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.INDIRECT_OBJECT_PROPERTY_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.INDIRECT_OBJECT] = 1
	staticComponentFrequency[tree.INDIRECT_OBJECT_REFERENCE] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.INDIRECT_OBJECT_ANNOTATION] = 1
	}

	// Shared elements
	staticComponentFrequency[tree.ACTIVATION_CONDITION] = 1
	staticComponentFrequency[tree.ACTIVATION_CONDITION_REFERENCE] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.ACTIVATION_CONDITION_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.EXECUTION_CONSTRAINT] = 1
	staticComponentFrequency[tree.EXECUTION_CONSTRAINT_REFERENCE] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.EXECUTION_CONSTRAINT_ANNOTATION] = 1
	}

	// Constitutive side
	staticComponentFrequency[tree.CONSTITUTED_ENTITY_PROPERTY] = 1
	staticComponentFrequency[tree.CONSTITUTED_ENTITY_PROPERTY_REFERENCE] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.CONSTITUTED_ENTITY_PROPERTY_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.CONSTITUTED_ENTITY] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.CONSTITUTED_ENTITY_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.MODAL] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.MODAL_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.CONSTITUTIVE_FUNCTION] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.CONSTITUTIVE_FUNCTION_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_PROPERTY] = 1
	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_PROPERTY_REFERENCE] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_PROPERTY_ANNOTATION] = 1
	}

	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES] = 1
	staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_REFERENCE] = 1
	if include_ANNOTATIONS {
		staticComponentFrequency[tree.CONSTITUTING_PROPERTIES_ANNOTATION] = 1
	}

	// Or else only exists as reference
	staticComponentFrequency[tree.OR_ELSE_REFERENCE] = 1

	return staticComponentFrequency
}
