package types

import (
	"reflect"
)

// ==========================
//
// Global Converters
//
// ==========================

// TypeConverter --
type TypeConverter func(source interface{}, targetType reflect.Type) (interface{}, error)
