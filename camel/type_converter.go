package camel

import (
	"fmt"
	"reflect"
)

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

	for _, converter := range typeConverter.converters {
		r, err := converter.Convert(source, targetType)
		if err == nil {
			return r, nil
		}
	}

	return nil, fmt.Errorf("Unsupported type conversion (source:%v, target:%v", sourceType, targetType)
}
