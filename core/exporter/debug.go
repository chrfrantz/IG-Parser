package exporter

import (
	"IG-Parser/core/config"
	"fmt"
)

/*
Prints output corresponding to debug settings.
*/
func Println(content ...interface{}) {
	if config.DEBUG_ALL || config.DEBUG_OUTPUT_GENERATION {
		fmt.Println(content...)
	}
}
