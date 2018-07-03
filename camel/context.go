package camel

import (
	"fmt"
	"net/url"
	"reflect"
	"sync"

	"github.com/lburgazzoli/camel-go/api"
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
	api.Service

	parent       *Context
	name         string
	registry     api.LoadingRegistry
	routes       []*Route
	converters   []api.TypeConverter
	services     []api.Service
	servicesLock sync.RWMutex
}

// RootContext --
var RootContext = NewContextWithParentAndName(nil, "root")

// ==========================
//
// Initialize a camel context
//
// ==========================

// NewContext --
func NewContext() *Context {
	return NewContextWithParentAndName(RootContext, "camel")
}

// NewContextWithParent --
func NewContextWithParent(parent *Context) *Context {
	return NewContextWithParentAndName(parent, "camel")
}

// NewContextWithName --
func NewContextWithName(name string) *Context {
	return NewContextWithParentAndName(RootContext, name)
}

// NewContextWithParentAndName --
func NewContextWithParentAndName(paretn *Context, name string) *Context {
	context := Context{
		parent:     paretn,
		name:       name,
		routes:     make([]*Route, 0),
		converters: make([]api.TypeConverter, 0),
		services:   make([]api.Service, 0),
	}

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
func (context *Context) Registry() api.LoadingRegistry {
	return context.registry
}

// AddTypeConverter --
func (context *Context) AddTypeConverter(converter api.TypeConverter) {
	context.converters = append(context.converters, converter)
}

// TypeConverter --
func (context *Context) TypeConverter() api.TypeConverter {
	converter := func(source interface{}, targetType reflect.Type) (interface{}, error) {
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

	if context.parent != nil {
		return api.NewConbinedTypeConverter(converter, context.parent.TypeConverter())
	}

	return converter
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

// Component --
func (context *Context) Component(name string) (Component, error) {
	value, found := context.lookup(name)

	if found && value != nil {
		if ca, ok := value.(ContextAware); ok {
			ca.SetContext(context)
		}

		if component, ok := value.(Component); ok {
			added := context.addService(component)
			if added {
				zlog.Debug().Msgf("Component with scheme %s registered as service", name)
			}

			return component, nil
		}
	}

	if context.parent != nil {
		return context.parent.Component(name)
	}

	return nil, fmt.Errorf("unable to find component with scheme: %s", name)
}

// Endpoint --
func (context *Context) Endpoint(uri string) (Endpoint, error) {
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

func (context *Context) lookup(name string) (interface{}, bool) {
	value, found := context.registry.Lookup(name)

	if found {
		if ca, ok := value.(ContextAware); ok {
			if ca.Context() == nil {
				ca.SetContext(context)
			}
		}
	}

	return value, found
}

func (context *Context) addService(service api.Service) bool {
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
	var s api.Service
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
			if s, ok := p.(api.Service); ok {
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
