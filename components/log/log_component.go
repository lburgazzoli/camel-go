package log

import (
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/introspection"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
//
//
// ==========================

// NewComponent --
func NewComponent() camel.Component {
	component := &Component{
		logger:         zlog.With().Str("logger", "log.Component").Logger(),
		serviceSupport: camel.NewServiceSupport(),
	}

	component.serviceSupport.Transition(camel.ServiceStatusSTOPPED, camel.ServiceStatusSTARTED, component.doStart)
	component.serviceSupport.Transition(camel.ServiceStatusSTARTED, camel.ServiceStatusSTOPPED, component.doStop)

	return component
}

// ==========================
//
// Component
//
// ==========================

// Component --
type Component struct {
	logger         zerolog.Logger
	serviceSupport camel.ServiceSupport
	context        *camel.Context
}

// SetContext --
func (component *Component) SetContext(context *camel.Context) {
	component.context = context
}

// Context --
func (component *Component) Context() *camel.Context {
	return component.context
}

// Start --
func (component *Component) Start() {
	component.serviceSupport.To(camel.ServiceStatusSTARTED)
}

// Stop --
func (component *Component) Stop() {
	component.serviceSupport.To(camel.ServiceStatusSTOPPED)
}

// CreateEndpoint --
func (component *Component) CreateEndpoint(remaining string, options map[string]interface{}) (camel.Endpoint, error) {
	// Create the endpoint and set default values
	endpoint := logEndpoint{}
	endpoint.component = component
	endpoint.logger = remaining
	endpoint.level = zerolog.InfoLevel

	// bind options to endpoint
	introspection.SetProperties(component.context, &endpoint, options)

	return &endpoint, nil
}

// ==========================
//
// Helpers
//
// ==========================

func (component *Component) doStart() {
	component.logger.Debug().Msg("Started")
}

func (component *Component) doStop() {
	component.logger.Debug().Msg("Stopped")
}
