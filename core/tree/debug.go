package tree

import (
	"IG-Parser/core/config"
	"fmt"
)

/*
Prints output corresponding to debug settings.
*/
func Print(content ...interface{}) {
	if config.DEBUG_ALL || config.DEBUG_TREE_OPERATIONS {
		fmt.Print(content...)
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
