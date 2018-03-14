package timer

import (
	"errors"
	"time"

	"github.com/lburgazzoli/camel-go/camel"
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

func (endpoint *timerEndpoint) Component() camel.Component {
	return endpoint.component
}

func (endpoint *timerEndpoint) CreateProducer(pipe *camel.Pipe) (camel.Producer, error) {
	return nil, errors.New("log is Consumer only")
}

func (endpoint *timerEndpoint) CreateConsumer(pipe *camel.Pipe) (camel.Consumer, error) {
	return &timerConsumer{endpoint: endpoint, pipe: pipe}, nil
}
