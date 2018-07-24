package timer

import "github.com/lburgazzoli/camel-go/camel"

// ==========================
//
// Init
//
// ==========================

func init() {
	camel.RootContext.Registry().Bind("timer", NewComponent())
}

// ==========================
//
// this is where constant
// shoud be se
//
// ==========================
