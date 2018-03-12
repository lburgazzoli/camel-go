package camel

import (
	"fmt"
	"reflect"
	"sync"
)

// ==========================
//
// Global Converters
//
// ==========================

var gTypeConverters = make([]TypeConverter, 0)
var gTypeConvertersLock = sync.RWMutex{}

// AddTypeConverter --
func AddTypeConverter(converter TypeConverter) {
	gTypeConvertersLock.Lock()
	gTypeConverters = append(gTypeConverters, converter)
	gTypeConvertersLock.Unlock()
}

// TypeConverter --
type TypeConverter interface {
	Convert(source interface{}, targetType reflect.Type) (interface{}, error)
}

// ==========================
//
// DelegatingTypeConverter
//
// ==========================

// DelegatingTypeConverter --
type DelegatingTypeConverter struct {
	converters []TypeConverter
}

// AddConverter --
func (typeConverter *DelegatingTypeConverter) AddConverter(converter TypeConverter) {
	typeConverter.converters = append(typeConverter.converters, converter)
}

// Convert --
func (typeConverter *DelegatingTypeConverter) Convert(source interface{}, targetType reflect.Type) (interface{}, error) {
	sourceType := reflect.TypeOf(source)

	if sourceType == targetType {
		return source, nil
	}
	if sourceType.Kind() == targetType.Kind() {
		return source, nil
	}

	// Use global type converters
	gTypeConvertersLock.RLock()
	defer gTypeConvertersLock.RUnlock()
	for _, converter := range gTypeConverters {
		r, err := converter.Convert(source, targetType)
		if err == nil {
			return r, nil
		}
	}

	// Context type converters
	for _, converter := range typeConverter.converters {
		r, err := converter.Convert(source, targetType)
		if err == nil {
			return r, nil
		}
	}

	return nil, fmt.Errorf("unsupported type conversion (source:%v, target:%v", sourceType, targetType)
}
