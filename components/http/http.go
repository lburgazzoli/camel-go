package http

import "github.com/lburgazzoli/camel-go/camel"

// ==========================
//
// Init
//
// ==========================

func init() {
	camel.RootContext.Registry().Bind("http", NewComponent())
}

// ==========================
//
// this is where constant
// shoud be se
//
// ==========================

// HTTPHeaderPrefix --
const HTTPHeaderPrefix = "http."

// HTTPHeaderPrefixLen --
const HTTPHeaderPrefixLen = len(HTTPHeaderPrefix)
