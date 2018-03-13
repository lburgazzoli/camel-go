package camel

import (
	"fmt"
	"net/url"
	"reflect"

	"github.com/lburgazzoli/camel-go/types"
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

	name       string
	registry   *Registry
	routes     []*Route
	converters []types.TypeConverter
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
	}

	// Type conversion
	context.AddTypeConverter(types.ToIntConverter)
	context.AddTypeConverter(types.ToDuratioinConverter)
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

// AddComponent --
func (context *Context) AddComponent(name string, component Component) {
	component.SetContext(context)

	context.registry.Bind(name, component)
}

// Component --
func (context *Context) Component(name string) (Component, error) {
	value, err := context.registry.Lookup(name)

	if err != nil {
		return nil, err
	}

	if value != nil {
		if component, ok := value.(Component); ok {
			return component, nil
		}
	}

	return nil, fmt.Errorf("Unable toi find component %s", name)
}

// AddRoute --
func (context *Context) AddRoute(route *Route) {
	context.routes = append(context.routes, route)
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
	for _, service := range context.routes {
		service.Start()
	}
}

// Stop --
func (context *Context) Stop() {
	for _, service := range context.routes {
		service.Stop()
	}
}
