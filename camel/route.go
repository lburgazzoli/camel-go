package camel

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
	route.services = append(route.services, service)
}

// Start --
func (route *Route) Start() {
	StartServices(route.services)
}

// Stop --
func (route *Route) Stop() {
	StopServices(route.services)
}
