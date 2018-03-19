package camel

// ==========================
//
// Route
//
// ==========================

// Route --
type Route struct {
	Service

	services []Service
}

// NewRoute --
func NewRoute() *Route {
	return &Route{
		services: make([]Service, 0),
	}
}

// AddService --
func (route *Route) AddService(service Service) {
	if service != nil {
		route.services = append(route.services, service)
	}
}

// Start --
func (route *Route) Start() {
	StartServices(route.services)
}

// Stop --
func (route *Route) Stop() {
	StopServices(route.services)
}

// ==========================
//
// Route Loader
//
// ==========================

// RouteLoader --
type RouteLoader interface {
	Load() (*Route, error)
}
