package log

import (
	"errors"

	zl "github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/lburgazzoli/camel-go/camel"
)

// ==========================
//
//
//
// ==========================

// NewComponent --
func NewComponent() camel.Component {
	return &Component{}
}

// ==========================
//
// Component
//
// ==========================

// Component --
type Component struct {
	context *camel.Context
}

// SetContext --
func (component *Component) SetContext(context *camel.Context) {
	component.context = context
}

// Context --
func (component *Component) Context() *camel.Context {
	return component.context
}

// CreateEndpoint --
func (component *Component) CreateEndpoint(remaining string, options map[string]interface{}) (camel.Endpoint, error) {
	// Create the endpoint and set default values
	endpoint := logEndpoint{component: component}
	endpoint.name = remaining
	endpoint.level = zl.InfoLevel

	// bind options to endpoint
	camel.SetFields(component.context, &endpoint, options)

	return &endpoint, nil
}

// ==========================
//
// Endpoint
//
// ==========================

type logEndpoint struct {
	component camel.Component
	name      string
	level     zl.Level
}

func (endpoint *logEndpoint) Component() camel.Component {
	return endpoint.component
}

func (endpoint *logEndpoint) CreateProducer() (camel.Producer, error) {
	sublogger := zlog.With().Str("name", endpoint.name)
	logger := sublogger.Logger()

	return &logProducer{
		endpoint: endpoint,
		logger:   &logger,
	}, nil
}

func (endpoint *logEndpoint) CreateConsumer() (camel.Consumer, error) {
	return nil, errors.New("log is Producer only")
}

func (endpoint *logEndpoint) SetName(name string) {
	endpoint.name = name
}

func (endpoint *logEndpoint) SetLevel(level zl.Level) {
	endpoint.level = level
}

// ==========================
//
// Producer
//
// ==========================

type logProducer struct {
	endpoint *logEndpoint
	logger   *zl.Logger
}

func (producer *logProducer) Endpoint() camel.Endpoint {
	return producer.endpoint
}

func (producer *logProducer) Process(exchange camel.Exchange) {
	producer.logger.WithLevel(producer.endpoint.level).Msgf("%+v", exchange.Body())
}
