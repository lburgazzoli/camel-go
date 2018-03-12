package camel

import (
	"fmt"
	"reflect"
)

// Context --
type Context struct {
	name            string
	registryLoaders []RegistryLoader
	components      map[string]Component
	converters      []TypeConverter
}

// ==========================
//
// Initialize a camel context
//
// ==========================

// NewContext --
func NewContext() *Context {
	return &Context{
		name:            "camel",
		registryLoaders: make([]RegistryLoader, 0),
		components:      make(map[string]Component),
		converters: []TypeConverter{
			ToIntConverter,
			ToDuratioinConverter,
			ToLogLevelConverter,
		},
	}
}

// NewContextWithName --
func NewContextWithName(name string) *Context {
	return &Context{
		name:            name,
		registryLoaders: make([]RegistryLoader, 0),
		components:      make(map[string]Component),
		converters: []TypeConverter{
			ToIntConverter,
			ToDuratioinConverter,
			ToLogLevelConverter,
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
func (context *Context) AddTypeConverter(converter TypeConverter) {
	context.converters = append(context.converters, converter)
}

// TypeConverter --
func (context *Context) TypeConverter() TypeConverter {
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

// ==========================
//
// Lyfecycle
//
// ==========================

// Start --
func (context *Context) Start() {
}

// Stop --
func (context *Context) Stop() {
}
