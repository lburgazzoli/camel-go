package api

// ==========================
//
// Route
//
// ==========================

// Route --
type Route struct {
	Service
	id       string
	services []Service
}

// NewRoute --
func NewRoute(id string) *Route {
	return &Route{
		id:       id,
		services: make([]Service, 0),
	}
}

// ID --
func (route *Route) ID() string {
	return route.id
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
