package main

import (
	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"

	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// plugin entry-pooint
//
// ==========================

// trace --
func trace(e api.Exchange) {
	count, found := e.Headers().LookupAs("timer.fire.count", camel.TypeInt)

	zlog.Info().Msgf("count: %+v, found=%v", count, found)
}

// Create --
func Create() interface{} {
	return trace
}
