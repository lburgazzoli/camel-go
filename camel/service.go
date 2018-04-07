package camel

import "sync/atomic"

// ServiceStatus --
type ServiceStatus int32

const (
	_ = iota

	// ServiceStatusSTOPPED --
	ServiceStatusSTOPPED

	// ServiceStatusSUSPENDED --
	ServiceStatusSUSPENDED

	// ServiceStatusSTARTED --
	ServiceStatusSTARTED
)

// Service --
type Service interface {
	Start()
	Stop()
}

// StartServices --
func StartServices(services []Service) {
	for _, service := range services {
		service.Start()
	}
}

// StopServices --
func StopServices(services []Service) {
	for _, service := range services {
		service.Stop()
	}
}

// TODO: state machine alike

// NewServiceState --
func NewServiceState(status ServiceStatus) ServiceState {
	return ServiceState{status: int32(status)}
}

// ServiceState --
type ServiceState struct {
	status int32
}

// Set --
func (state *ServiceState) Set(status ServiceStatus) {
	atomic.StoreInt32(&state.status, int32(status))
}

// Transition --
func (state *ServiceState) Transition(from ServiceStatus, to ServiceStatus, transition func()) bool {

	if atomic.CompareAndSwapInt32(&state.status, int32(from), int32(to)) {
		transition()

		return true
	}

	return false
}
