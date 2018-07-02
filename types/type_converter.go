package types

import (
	"fmt"
	"reflect"
)

// ==========================
//
// Global Converters
//
// ==========================

// TypeConverter --
type TypeConverter func(source interface{}, targetType reflect.Type) (interface{}, error)

// NewConbinedTypeConverter --
func NewConbinedTypeConverter(converter TypeConverter, converters ...TypeConverter) TypeConverter {
	return func(source interface{}, targetType reflect.Type) (interface{}, error) {
		var answer interface{}
		var err error

		answer, err = converter(source, targetType)
		if answer != nil && err == nil {
			return answer, err
		}

		for _, c := range converters {
			answer, err = c(source, targetType)
			if answer != nil && err == nil {
				return answer, err
			}
		}

		return nil, fmt.Errorf("unsupported type conversion (source:%v, target:%v", source, targetType)
	}
}
