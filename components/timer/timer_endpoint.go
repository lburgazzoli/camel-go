package timer

import (
	"errors"
	"time"

	"github.com/lburgazzoli/camel-go/api"
)

// ==========================
//
// Endpoint
//
// ==========================

type timerEndpoint struct {
	component *Component
	period    time.Duration
}

func (endpoint *timerEndpoint) SetPeriod(period time.Duration) {
	endpoint.period = period
}

func (endpoint *timerEndpoint) Start() {
}

func (endpoint *timerEndpoint) Stop() {
}

func (endpoint *timerEndpoint) Stage() api.ServiceStage {
	return api.ServiceStageEndpoint
}

func (endpoint *timerEndpoint) Component() api.Component {
	return endpoint.component
}

func (endpoint *timerEndpoint) CreateProducer() (api.Producer, error) {
	return nil, errors.New("log is Consumer only")
}

func (endpoint *timerEndpoint) CreateConsumer() (api.Consumer, error) {
	return newTimerConsumer(endpoint), nil
}
