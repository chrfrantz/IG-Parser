package tree

import (
	"IG-Parser/core/config"
	"fmt"
	"log"
)

/*
Prints output corresponding to debug settings.
*/
func Print(content ...interface{}) {
	if config.DEBUG_ALL || config.DEBUG_TREE_OPERATIONS {
		log.Print(content...)
	}
}

/*
Prints output corresponding to debug settings.
*/
func Println(content ...interface{}) {
	if config.DEBUG_ALL || config.DEBUG_TREE_OPERATIONS {
		fmt.Println(content...)
	}
}
