// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package camel

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/lburgazzoli/camel-go/api"
)

// ==========================
//
//
//
// ==========================

// Context --
type defaultContext struct {
	api.Context

	parent       api.Context
	name         string
	registry     api.LoadingRegistry
	routes       []*api.Route
	converters   []api.TypeConverter
	converter    api.TypeConverter
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
func NewContext() api.Context {
	return NewContextWithParentAndName(RootContext, "camel")
}

// NewContextWithParent --
func NewContextWithParent(parent api.Context) api.Context {
	return NewContextWithParentAndName(parent, "camel")
}

// NewContextWithName --
func NewContextWithName(name string) api.Context {
	return NewContextWithParentAndName(RootContext, name)
}

// NewContextWithParentAndName --
func NewContextWithParentAndName(parent api.Context, name string) api.Context {
	context := defaultContext{
		parent:     parent,
		name:       name,
		routes:     make([]*api.Route, 0),
		converters: make([]api.TypeConverter, 0),
		services:   make([]api.Service, 0),
	}

	context.converter = func(source interface{}, targetType reflect.Type) (interface{}, error) {
		if source == nil {
			return nil, fmt.Errorf("unsupported type conversion (source:nil, target:%v", targetType)
		}

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

	if parent != nil {
		// Set the registry
		context.registry = api.NewCombinedRegistry(
			NewRegistry(context.TypeConverter()),
			parent.Registry(),
		)

		context.converter = api.NewConbinedTypeConverter(context.converter, parent.TypeConverter())
	} else {
		context.registry = NewRegistry(context.TypeConverter())
	}

	return &context
}

// ==========================
//
//
//
// ==========================

// Registry --
func (context *defaultContext) Registry() api.LoadingRegistry {
	return context.registry
}

// AddTypeConverter --
func (context *defaultContext) AddTypeConverter(converter api.TypeConverter) {
	context.converters = append(context.converters, converter)
}

// TypeConverter --
func (context *defaultContext) TypeConverter() api.TypeConverter {
	return context.converter
}

// AddRouteDefinition --
func (context *defaultContext) AddRoute(route *api.Route) {
	context.routes = append(context.routes, route)

	context.AddService(route)
}

func (context *defaultContext) AddService(service api.Service) bool {
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

// ==========================
//
// Lyfecycle
//
// ==========================

// Start --
func (context *defaultContext) Start() {
	for _, service := range context.services {
		service.Start()
	}
}

// Stop --
func (context *defaultContext) Stop() {
	for _, service := range context.services {
		service.Stop()
	}
}

// ==========================
//
// Helpers
//
// ==========================

func (context *defaultContext) lookup(name string) (interface{}, bool) {
	value, found := context.registry.Lookup(name)

	if found {
		if ca, ok := value.(api.ContextAware); ok {
			if ca.Context() == nil {
				ca.SetContext(context)
			}
		}
	}

	return value, found
}
