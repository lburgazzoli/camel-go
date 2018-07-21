// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"sort"
	"sync/atomic"

	zlog "github.com/rs/zerolog/log"
)

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
	ServiceStageContext ServiceStage = 0

	// ServiceStageComponent --
	ServiceStageComponent ServiceStage = 10

	// ServiceStageEndpoint --
	ServiceStageEndpoint ServiceStage = 20

	// ServiceStageProducer --
	ServiceStageProducer ServiceStage = 30

	// ServiceStageConsumer --
	ServiceStageConsumer ServiceStage = 40

	// ServiceStageOther --
	ServiceStageOther ServiceStage = 50
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

// StagedService --
type StagedService interface {
	Service

	Stage() ServiceStage
}

// StartServices --
func StartServices(services []Service) {
	svcs := make([]Service, len(services))
	copy(svcs, services)

	// Sort the
	sort.SliceStable(svcs, func(i int, j int) bool {
		istage := ServiceStageOther
		jstage := ServiceStageOther

		if stage, ok := svcs[i].(StagedService); ok {
			istage = stage.Stage()
		}
		if stage, ok := svcs[j].(StagedService); ok {
			jstage = stage.Stage()
		}

		return istage < jstage
	})

	for _, service := range svcs {
		if stage, ok := service.(StagedService); ok {
			zlog.Debug().Msgf("Starting service at stage: %d", stage.Stage())
		}

		service.Start()
	}
}

// StopServices --
func StopServices(services []Service) {
	svcs := make([]Service, len(services))
	copy(svcs, services)

	sort.SliceStable(svcs, func(i int, j int) bool {
		istage := ServiceStageOther
		jstage := ServiceStageOther

		if stage, ok := svcs[i].(StagedService); ok {
			istage = stage.Stage()
		}
		if stage, ok := svcs[j].(StagedService); ok {
			jstage = stage.Stage()
		}

		return istage > jstage
	})

	for _, service := range svcs {
		if stage, ok := service.(StagedService); ok {
			zlog.Debug().Msgf("Stopping service at stage: %d", stage.Stage())
		}

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
