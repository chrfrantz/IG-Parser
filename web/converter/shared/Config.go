package shared

import (
	"IG-Parser/core/exporter/tabular"
)

/*
Contains the default configuration for tabular output.
*/
func SetDefaultConfig() {
	tabular.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
	tabular.SetProduceIGExtendedOutput(false)
}
