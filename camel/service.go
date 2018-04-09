package camel

import "sync/atomic"

// ==========================
//
// ServiceStatus
//
// ==========================

// ServiceStatus --
type ServiceStatus int32

// Is --
func (status ServiceStatus) Is(other ServiceStatus) bool {
	return status == other
}

// IsInTransition --
func (status ServiceStatus) IsInTransition() bool {
	return status.Is(ServiceStatusTRANSITION)
}

// IsStopped --
func (status ServiceStatus) IsStopped() bool {
	return status.Is(ServiceStatusSTOPPED)
}

// IsSuspended --
func (status ServiceStatus) IsSuspended() bool {
	return status.Is(ServiceStatusSUSPENDED)
}

// IsStarted --
func (status ServiceStatus) IsStarted() bool {
	return status.Is(ServiceStatusSTARTED)
}

const (
	// ServiceStatusTRANSITION --
	ServiceStatusTRANSITION ServiceStatus = iota

	// ServiceStatusSTOPPED --
	ServiceStatusSTOPPED

	// ServiceStatusSUSPENDED --
	ServiceStatusSUSPENDED

	// ServiceStatusSTARTED --
	ServiceStatusSTARTED
)

// ==========================
//
// Service
//
// ==========================

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

// ==========================
//
// ServiceState
//
// ==========================

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

// ==========================
//
// ServiceX
//
// ==========================

type stateHandlers map[ServiceStatus]func()

// ServiceX  --
type ServiceX struct {
	status int32
}
