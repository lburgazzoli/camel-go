package main

import (
	"math/rand"

	"github.com/lburgazzoli/camel-go/api"
)

// ==========================
//
// plugin entry-pooint
//
// ==========================

// SetHeader --
func setHeader(e api.Exchange) {
	e.Headers().Bind("plugin.header", rand.Int())
}

// Create --
func Create() interface{} {
	return setHeader
}
