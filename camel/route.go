package camel

import "github.com/lburgazzoli/camel-go/api"

// ==========================
//
// Route
//
// ==========================

// Route --
type Route struct {
	api.Service

	services []api.Service
}

// NewRoute --
func NewRoute() *Route {
	return &Route{
		services: make([]api.Service, 0),
	}
}

// AddService --
func (route *Route) AddService(service api.Service) {
	if service != nil {
		route.services = append(route.services, service)
	}
}

// Start --
func (route *Route) Start() {
	api.StartServices(route.services)
}

// Stop --
func (route *Route) Stop() {
	api.StopServices(route.services)
}
