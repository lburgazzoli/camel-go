package camel

import "reflect"

// Processor --
type Processor func(*Exchange)

// IsProcessor --
func IsProcessor(ifc interface{}) bool {
	fv := reflect.ValueOf(ifc)
	ft := fv.Type()

	if fv.Kind() != reflect.Func {
		return false
	}

	if ft.NumIn() != 1 {
		return false
	}

	return ft.In(0) == reflect.TypeOf(&Exchange{})
}

// Trasformer --
type Trasformer func(*Exchange) *Exchange

// Predicate --
type Predicate func(*Exchange) bool

// IsPredicate --
func IsPredicate(ifc interface{}) bool {
	fv := reflect.ValueOf(ifc)
	ft := fv.Type()

	if fv.Kind() != reflect.Func {
		return false
	}

	if ft.NumIn() != 1 {
		return false
	}

	if ft.NumOut() != 1 {
		return false
	}

	return ft.In(0) == reflect.TypeOf(&Exchange{}) && ft.Out(0).Kind() == reflect.Bool
}
