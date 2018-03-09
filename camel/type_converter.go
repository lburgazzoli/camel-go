package camel

import (
	"errors"
	"reflect"
)

// TypeConverter --
type TypeConverter interface {
	CanConvert(sourceType reflect.Type, targetType reflect.Type) bool
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

// CanConvert --
func (typeConverter *DelegatingTypeConverter) CanConvert(sourceType reflect.Type, targetType reflect.Type) bool {
	if sourceType == targetType {
		return true
	}

	for _, converter := range typeConverter.converters {
		if converter.CanConvert(sourceType, targetType) {
			return true
		}
	}

	return false
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

	return nil, errors.New("Unsupported type conversion")
}
