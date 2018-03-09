package camel

import (
	"fmt"
	"reflect"
	"strconv"
)

// DefaultTypeConverters --
func DefaultTypeConverters() []TypeConverter {
	return []TypeConverter{
		&FromStringConverter{},
	}
}

// ==========================
//
// String converter
//
// ==========================

// FromStringConverter --
type FromStringConverter struct {
	TypeConverter
}

// CanConvert --
func (converter *FromStringConverter) CanConvert(sourceType reflect.Type, targetType reflect.Type) bool {
	return sourceType.Kind() == reflect.String
}

// Convert --
func (converter *FromStringConverter) Convert(source interface{}, targetType reflect.Type) (interface{}, error) {
	sourceType := reflect.TypeOf(source)
	sourceKind := sourceType.Kind()
	targetKind := targetType.Kind()

	if sourceKind != reflect.String {
		return nil, fmt.Errorf("Incompatible source, expected:%v got:%v", reflect.String, sourceKind)
	}
	if !converter.CanConvert(sourceType, targetType) {
		return nil, fmt.Errorf("Unable to convert from:%v to:%v", sourceType, targetType)
	}

	str := source.(string)
	var answer interface{}
	var err error

	switch targetKind {
	case reflect.String:
		answer = source
	case reflect.Int8:
		i, err := strconv.ParseInt(str, 10, 8)
		if err == nil {
			answer = int8(i)
		}
	case reflect.Int16:
		i, err := strconv.ParseInt(str, 10, 16)
		if err == nil {
			answer = int16(i)
		}
	case reflect.Int32:
		i, err := strconv.ParseInt(str, 10, 32)
		if err == nil {
			answer = int32(i)
		}
	case reflect.Int64:
		i, err := strconv.ParseInt(str, 10, 64)
		if err == nil {
			answer = int32(i)
		}
	default:
		answer = nil
		err = fmt.Errorf("Unable to convert from:%v to:%v, error:%v", sourceType, targetType, err)
	}

	if err != nil {
		return nil, fmt.Errorf("Unable to convert from:%v to:%v, error:%v", sourceType, targetType, err)
	}

	return answer, nil
}
