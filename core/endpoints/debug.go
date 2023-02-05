package endpoints

import (
	"IG-Parser/core/config"
	"log"
)

/*
Prints output corresponding to debug settings.
*/
func Println(content ...interface{}) {
	if config.DEBUG_ALL || config.DEBUG_FINAL_OUTPUT {
		log.Println(content...)
	}
}
