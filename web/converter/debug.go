package converter

import (
	"IG-Parser/core/config"
	"log"
)

/*
Prints output corresponding to debug settings.
*/
func Println(content ...interface{}) {
	if config.DEBUG_ALL || config.DEBUG_FRONTEND {
		log.Println(content...)
	}
}
