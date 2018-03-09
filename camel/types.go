package camel

import (
	"reflect"
)

// TypeInt --
var TypeInt = reflect.TypeOf(int(0))

// TypeInt8 --
var TypeInt8 = reflect.TypeOf(int8(0))

// TypeInt16 --
var TypeInt16 = reflect.TypeOf(int16(0))

// TypeInt32 --
var TypeInt32 = reflect.TypeOf(int32(0))

// TypeInt64 --
var TypeInt64 = reflect.TypeOf(int64(0))

// IsInt --
func IsInt(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int:
		return true
	case reflect.Int8:
		return true
	case reflect.Int16:
		return true
	case reflect.Int32:
		return true
	case reflect.Int64:
		return true
	default:
		return false
	}
}
