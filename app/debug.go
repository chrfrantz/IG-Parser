package app

import (
	"IG-Parser/config"
	"fmt"
)

/*
Prints output corresponding to debug settings.
*/
func Println(content ...interface{}) {
	if config.DEBUG_ALL || config.DEBUG_FINAL_OUTPUT {
		fmt.Println(content...)
	}
}
