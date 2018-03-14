package camel

// Route --
type Route struct {
	Service

	services []Service
}

// ToRoute --
type ToRoute interface {
	// ToRoute --
	ToRoute(context *Context) (*Route, error)
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
