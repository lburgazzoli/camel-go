package api

import "sync/atomic"

// ==========================
//
// ServiceStatus
//
// ==========================

// ServiceStatus --
type ServiceStatus uint32

// Is --
func (status ServiceStatus) Is(other ServiceStatus) bool {
	return status == other
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
	serviceStatusTRANSITION ServiceStatus = iota

	// ServiceStatusSTOPPED --
	ServiceStatusSTOPPED

	// ServiceStatusSUSPENDED --
	ServiceStatusSUSPENDED

	// ServiceStatusSTARTED --
	ServiceStatusSTARTED
)

// ==========================
//
// ServiceStage
//
// ==========================

// ServiceStage --
type ServiceStage uint32

const (
	// ServiceStageContext --
	ServiceStageContext ServiceStage = iota

	// ServiceStageComponent --
	ServiceStageComponent

	// ServiceStageEndpoint --
	ServiceStageEndpoint

	// ServiceStageMisc --
	ServiceStageMisc
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
// ServiceSupport
//
// ==========================

// NewServiceSupport --
func NewServiceSupport() ServiceSupport {
	return ServiceSupport{
		status:   uint32(ServiceStatusSTOPPED),
		handlers: make(map[ServiceStatus]map[ServiceStatus]func()),
	}
}

// ServiceSupport  --
type ServiceSupport struct {
	status uint32

	handlers map[ServiceStatus]map[ServiceStatus]func()
}

// Transition --
func (support *ServiceSupport) Transition(from ServiceStatus, to ServiceStatus, transition func()) {
	handlers, ok := support.handlers[to]
	if !ok {
		handlers = make(map[ServiceStatus]func())
		support.handlers[to] = handlers
	}

	handlers[from] = transition
}

// To --
func (support *ServiceSupport) To(to ServiceStatus) bool {
	if handlers, ok := support.handlers[to]; ok {
		for from, function := range handlers {
			// Set the state to the reserved 'in transition' state
			if atomic.CompareAndSwapUint32(&support.status, uint32(from), uint32(serviceStatusTRANSITION)) {
				function()

				// set the targer state
				atomic.StoreUint32(&support.status, uint32(to))

				return true
			}
		}
	}

	return false
}
