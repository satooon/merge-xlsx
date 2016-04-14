package main

import "log"

type debugLogStruct struct {
	verbose bool
}

var debugLog debugLogStruct

func (l debugLogStruct) Println(args ...interface{}) {
	if l.verbose == false {
		return
	}
	log.Println(args...)
}
