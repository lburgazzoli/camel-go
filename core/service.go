package core

import "github.com/lburgazzoli/camel-go/camel"

// DefaultService --
type DefaultService struct {
	order  int
	status camel.ServiceStatus
}

// Order --
func (service *DefaultService) Order() int {
	return service.order
}

// Start --
func (service *DefaultService) Start() {
}

// Stop --
func (service *DefaultService) Stop() {
}

// Status --
func (service *DefaultService) Status() camel.ServiceStatus {
	return service.status
}
