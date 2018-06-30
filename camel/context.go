package camel

import (
	"fmt"
	"net/url"
	"reflect"
	"sync"

	"github.com/lburgazzoli/camel-go/types"
	zlog "github.com/rs/zerolog/log"
)

// ==========================
//
//
//
// ==========================

// HasContext --
type HasContext interface {
	Context() *Context
}

// ContextAware --
type ContextAware interface {
	HasContext

	SetContext(context *Context)
}

// Context --
type Context struct {
	Service

	name         string
	registry     *Registry
	routes       []*Route
	converters   []types.TypeConverter
	services     []Service
	servicesLock sync.RWMutex
}

// ==========================
//
// Initialize a camel context
//
// ==========================

// NewContext --
func NewContext() *Context {
	return NewContextWithName("camel")
}

// NewContextWithName --
func NewContextWithName(name string) *Context {
	context := Context{
		name:       name,
		routes:     make([]*Route, 0),
		converters: make([]types.TypeConverter, 0),
		services:   make([]Service, 0),
	}

	// Type conversion
	context.AddTypeConverter(types.ToIntConverter)
	context.AddTypeConverter(types.ToDurationConverter)
	context.AddTypeConverter(types.ToLogLevelConverter)
	context.AddTypeConverter(types.ToBoolConverter)

	// Set the registry
	context.registry = NewRegistry(context.TypeConverter())

	return &context
}

// ==========================
//
//
//
// ==========================

// Registry --
func (context *Context) Registry() *Registry {
	return context.registry
}

// AddTypeConverter --
func (context *Context) AddTypeConverter(converter types.TypeConverter) {
	context.converters = append(context.converters, converter)
}

// TypeConverter --
func (context *Context) TypeConverter() types.TypeConverter {
	return func(source interface{}, targetType reflect.Type) (interface{}, error) {
		sourceType := reflect.TypeOf(source)

		// Don't convert same type
		if sourceType == targetType {
			return source, nil
		}

		// Context type converters
		for _, converter := range context.converters {
			r, err := converter(source, targetType)
			if err == nil {
				return r, nil
			}
		}

		// Try implicit go conversion
		if sourceType.ConvertibleTo(targetType) {
			v := reflect.ValueOf(source).Convert(targetType)
			if v.IsValid() {
				return v.Interface(), nil
			}
		}

		return nil, fmt.Errorf("unsupported type conversion (source:%v, target:%v", sourceType, targetType)
	}
}

// Component --
func (context *Context) Component(name string) (Component, error) {
	value, err := context.lookup(name)

	if err != nil {
		return nil, err
	}

	if value != nil {
		if component, ok := value.(Component); ok {
			if added := context.addService(component); added {
				zlog.Debug().Msgf("Component with scheme %s registered as service", name)
			}

			return component, nil
		}
	}

	return nil, fmt.Errorf("unable to find component with scheme: %s", name)
}

// AddRouteDefinition --
func (context *Context) AddRouteDefinition(definition Definition) {
	route := &Route{}

	// Find the root
	for definition.Parent() != nil {
		definition = definition.Parent()
	}

	context.addDefinitionsToRoute(route, nil, definition)
	context.routes = append(context.routes, route)

	context.addService(route)
}

// CreateEndpointFromURI --
func (context *Context) CreateEndpointFromURI(uri string) (Endpoint, error) {
	var err error
	var endpointURL *url.URL
	var component Component
	var endpoint Endpoint

	if endpointURL, err = url.Parse(uri); err != nil {
		return nil, err
	}

	scheme := endpointURL.Scheme
	opts := make(map[string]interface{})
	vals := make(url.Values)

	if vals, err = url.ParseQuery(endpointURL.RawQuery); err != nil {
		return nil, err
	}

	for k, v := range vals {
		opts[k] = v[0]
	}

	if component, err = context.Component(scheme); err != nil {
		return nil, err
	}

	remaining := ""
	if endpointURL.Opaque != "" {
		if remaining, err = url.PathUnescape(endpointURL.Opaque); err != nil {
			return nil, err
		}
	} else {
		remaining = endpointURL.Host

		if endpointURL.RawPath != "" {
			path, err := url.PathUnescape(endpointURL.RawPath)
			if err != nil {
				return nil, err
			}

			remaining += path
		}
	}

	if endpoint, err = component.CreateEndpoint(remaining, opts); err != nil {
		return nil, err
	}

	return endpoint, nil
}

// ==========================
//
// Lyfecycle
//
// ==========================

// Start --
func (context *Context) Start() {
	for _, service := range context.services {
		service.Start()
	}
}

// Stop --
func (context *Context) Stop() {
	for _, service := range context.services {
		service.Stop()
	}
}

// ==========================
//
// Helpers
//
// ==========================

func (context *Context) lookup(name string) (interface{}, error) {
	value, err := context.registry.Lookup(name)

	if err != nil {
		return nil, err
	}

	if ca, ok := value.(ContextAware); ok {
		if ca.Context() == nil {
			ca.SetContext(context)
		}
	}

	return value, err
}

func (context *Context) addService(service Service) bool {
	context.servicesLock.Lock()
	defer context.servicesLock.Unlock()

	for _, s := range context.services {
		if s == service {
			return false
		}
	}

	context.services = append(context.services, service)

	return true
}

func (context *Context) addDefinitionsToRoute(route *Route, processor Processor, definition Definition) Processor {
	var s Service
	var e error

	p := processor

	if u, ok := definition.(Unwrappable); ok {
		p, s, e = u.Unwrap(context, p)

		if e != nil {
			zlog.Fatal().Msgf("unable to load processor %v (%s)", definition, e)
		}

		if e == nil && s != nil {
			route.AddService(s)
		}

		if p != nil {
			if s, ok := p.(Service); ok {
				route.AddService(s)
			}
		} else {
			p = processor
		}
	}

	for _, c := range definition.Children() {
		p = context.addDefinitionsToRoute(route, p, c)
	}

	return p
}
