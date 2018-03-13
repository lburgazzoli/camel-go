package camel

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/lburgazzoli/camel-go/types"
)

// ==========================
//
// Global Converters
//
// ==========================

// TypeConverters
var gTypeConverters = make([]types.TypeConverter, 0)
var gTypeConvertersLock = sync.RWMutex{}

// AddTypeConverter --
func AddTypeConverter(converter types.TypeConverter) {
	gTypeConvertersLock.Lock()
	gTypeConverters = append(gTypeConverters, converter)
	gTypeConvertersLock.Unlock()
}

// ==========================
//
// Global Converters
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

	name            string
	registryLoaders []RegistryLoader
	routes          []*Route
	components      map[string]Component
	converters      []types.TypeConverter
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
	return &Context{
		name:            name,
		registryLoaders: make([]RegistryLoader, 0),
		routes:          make([]*Route, 0),
		components:      make(map[string]Component),
		converters: []types.TypeConverter{
			types.ToIntConverter,
			types.ToDuratioinConverter,
			types.ToLogLevelConverter,
			types.ToBoolConverter,
		},
	}
}

// ==========================
//
//
//
// ==========================

// AddRegistryLoader --
func (context *Context) AddRegistryLoader(loader RegistryLoader) {
	context.registryLoaders = append(context.registryLoaders, loader)
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

		// Use global type converters
		gTypeConvertersLock.RLock()
		defer gTypeConvertersLock.RUnlock()
		for _, converter := range gTypeConverters {
			r, err := converter(source, targetType)
			if err == nil {
				return r, nil
			}
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
	context.components[name] = component
	context.components[name].SetContext(context)
}

// Component --
func (context *Context) Component(name string) (Component, error) {
	component, found := context.components[name]

	// check if the component has already been created or added to the context
	// component list
	if !found {
		for _, loader := range context.registryLoaders {
			component, err := loader.Load(name)

			if err != nil {
				return nil, err
			}

			if component == nil {
				continue
			}

			if _, ok := component.(Component); !ok {
				// not a component
				continue
			}

			if component != nil {
				break
			}
		}

		if component != nil {
			context.AddComponent(name, component)
		}
	}

	return component, nil
}

// AddRoute --
func (context *Context) AddRoute(route *Route) {
	context.routes = append(context.routes, route)
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
