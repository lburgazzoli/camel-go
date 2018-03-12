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

// TypeConverter --
type TypeConverter func(source interface{}, targetType reflect.Type) (interface{}, error)

var gTypeConverters = make([]TypeConverter, 0)
var gTypeConvertersLock = sync.RWMutex{}

// AddTypeConverter --
func AddTypeConverter(converter TypeConverter) {
	gTypeConvertersLock.Lock()
	gTypeConverters = append(gTypeConverters, converter)
	gTypeConvertersLock.Unlock()
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
	for _, converter := range typeConverter.converters {
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
