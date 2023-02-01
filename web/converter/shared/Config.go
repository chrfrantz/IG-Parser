package shared

import (
	"IG-Parser/core/exporter"
)

/*
Contains the default configuration for tabular output.
*/
func SetDefaultConfig() {
	exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	exporter.SetProduceIGExtendedOutput(false)
}
