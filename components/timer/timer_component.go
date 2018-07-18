package timer

import (
	"fmt"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/introspection"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
// Init
//
// ==========================

func init() {
	camel.RootContext.Registry().Bind("timer", NewComponent())
}

// ==========================
//
//
//
// ==========================

// NewComponent --
func NewComponent() api.Component {
	component := &Component{
		logger:         zlog.With().Str("logger", "timer.Component").Logger(),
		serviceSupport: api.NewServiceSupport(),
	}

	component.serviceSupport.Transition(api.ServiceStatusSTOPPED, api.ServiceStatusSTARTED, component.doStart)
	component.serviceSupport.Transition(api.ServiceStatusSTARTED, api.ServiceStatusSTOPPED, component.doStop)

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
	serviceSupport api.ServiceSupport
	context        api.Context
}

// SetContext --
func (component *Component) SetContext(context api.Context) {
	component.context = context
}

// Context --
func (component *Component) Context() api.Context {
	return component.context
}

// Start --
func (component *Component) Start() {
	component.serviceSupport.To(api.ServiceStatusSTARTED)
}

// Stop --
func (component *Component) Stop() {
	component.serviceSupport.To(api.ServiceStatusSTOPPED)
}

// Stage --
func (component *Component) Stage() api.ServiceStage {
	return api.ServiceStageComponent
}

// CreateEndpoint --
func (component *Component) CreateEndpoint(remaining string, options map[string]interface{}) (api.Endpoint, error) {
	// Create the endpoint and set default values
	endpoint := timerEndpoint{}
	endpoint.component = component

	// endpoint option validation
	if _, ok := options["period"]; !ok {
		return nil, fmt.Errorf("missing mandatory option: period")
	}

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
	logger := zlog.With().Str("logger", "timer.Component").Logger()
	logger.Info().Msg("Started")
}

func (component *Component) doStop() {
	logger := zlog.With().Str("logger", "timer.Component").Logger()
	logger.Info().Msg("Stopped")
}
