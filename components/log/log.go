package log

import "github.com/lburgazzoli/camel-go/camel"

// ==========================
//
// Init
//
// ==========================

func init() {
	camel.RootContext.Registry().Bind("log", NewComponent())
}

// ==========================
//
// this is where constant
// shoud be se
//
// ==========================
