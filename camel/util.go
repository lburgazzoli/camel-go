package camel

import "reflect"

// IsInt --
func IsInt(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int:
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
