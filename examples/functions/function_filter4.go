package main

import (
	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"	
)

// ==========================
//
// plugin entry-pooint
//
// ==========================

// Filter4 --
func filter4(e api.Exchange) bool {
	if count, found := e.Headers().LookupAs("timer.fire.count", camel.TypeInt); found {
		return count.(int) != 4
	}

	return true
}

// Create --
func Create() interface{} {
	return filter4
}
