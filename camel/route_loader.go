package camel

import "github.com/lburgazzoli/camel-go/api"

// ==========================
//
// Route Loader
//
// ==========================

// RouteLoader --
type RouteLoader interface {
	Load() ([]*api.Route, error)
}
