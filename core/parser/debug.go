package parser

import (
	"IG-Parser/core/config"
	"fmt"
	"log"
)

/*
Prints output corresponding to debug settings.
*/
func Print(content ...interface{}) {
	if config.DEBUG_ALL || config.DEBUG_STATEMENT_PARSING {
		log.Print(content...)
	}
}

/*
Prints output corresponding to debug settings.
*/
func Println(content ...interface{}) {
	if config.DEBUG_ALL || config.DEBUG_STATEMENT_PARSING {
		fmt.Println(content...)
	}
}
