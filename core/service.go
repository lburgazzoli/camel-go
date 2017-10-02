package core

import "github.com/lburgazzoli/camel-go/camel"

// NewService --
func NewService() camel.Service {
	return &defaultService{
		order:  0,
		status: camel.ServiceStatusSTOPPED,
	}
}

// NewServiceWithOrder --
func NewServiceWithOrder(order int) camel.Service {
	return &defaultService{
		order:  order,
		status: camel.ServiceStatusSTOPPED,
	}
}

// DefaultService --
type defaultService struct {
	order  int
	status camel.ServiceStatus
}

// Order --
func (service *defaultService) Order() int {
	return service.order
}

// Start --
func (service *defaultService) Start() {
}

// Stop --
func (service *defaultService) Stop() {
}

// Status --
func (service *defaultService) Status() camel.ServiceStatus {
	return service.status
}
