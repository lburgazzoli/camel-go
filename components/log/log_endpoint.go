package log

import (
	"errors"
	"os"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/rs/zerolog"
)

// ==========================
//
// Endpoint
//
// ==========================

type logEndpoint struct {
	component  *Component
	logger     string
	level      zerolog.Level
	logHeaders bool
}

func (endpoint *logEndpoint) Start() {
}

func (endpoint *logEndpoint) Stop() {
}

func (endpoint *logEndpoint) Stage() api.ServiceStage {
	return api.ServiceStageEndpoint
}

func (endpoint *logEndpoint) Component() api.Component {
	return endpoint.component
}

func (endpoint *logEndpoint) CreateProducer() (api.Producer, error) {
	// need to be replaced with better configuration from camel logging
	newlog := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger := newlog.With().Str("logger", endpoint.logger).Logger()

	return newLogProducer(endpoint, &logger), nil
}

func (endpoint *logEndpoint) CreateConsumer() (api.Consumer, error) {
	return nil, errors.New("log is Producer only")
}

func (endpoint *logEndpoint) SetLogger(logger string) {
	endpoint.logger = logger
}

func (endpoint *logEndpoint) SetLevel(level zerolog.Level) {
	endpoint.level = level
}

func (endpoint *logEndpoint) SetLogHeaders(logHeaders bool) {
	endpoint.logHeaders = logHeaders
}
