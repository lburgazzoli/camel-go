package log

import (
	"errors"
	"os"

	"github.com/lburgazzoli/camel-go/camel"
	"github.com/rs/zerolog"
)

// ==========================
//
// Endpoint
//
// ==========================

type logEndpoint struct {
	component  camel.Component
	name       string
	level      zerolog.Level
	logHeaders bool
}

func (endpoint *logEndpoint) Start() {
}

func (endpoint *logEndpoint) Stop() {
}

func (endpoint *logEndpoint) Component() camel.Component {
	return endpoint.component
}

func (endpoint *logEndpoint) CreateProducer() (camel.Producer, error) {
	// need to be replaced with better configuration from camel logging
	newlog := zerolog.New(os.Stderr).With().Timestamp().Logger()
	logger := newlog.With().Str("name", endpoint.name).Logger()

	return &logProducer{
		endpoint: endpoint,
		logger:   &logger,
	}, nil
}

func (endpoint *logEndpoint) CreateConsumer(producer camel.Producer) (camel.Consumer, error) {
	return nil, errors.New("log is Producer only")
}

func (endpoint *logEndpoint) SetName(name string) {
	endpoint.name = name
}

func (endpoint *logEndpoint) SetLevel(level zerolog.Level) {
	endpoint.level = level
}
