package shared

import (
	"IG-Parser/core/exporter/tabular"
)

/*
Contains the default configuration for tabular output.
*/
func SetDefaultConfig() {
	tabular.SetIncludeSharedElementsInTabularOutput(true)
	tabular.SetProduceIGExtendedOutput(false)
}
