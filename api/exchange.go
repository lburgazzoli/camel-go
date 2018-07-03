package api

import (
	"reflect"
)

// Headers --
type Headers struct {
	Registry
}

// Properties --
type Properties struct {
	Registry
}

// Exchange --
type Exchange interface {
	Body() interface{}
	BodyAs(asType reflect.Type) interface{}

	SetBody(body interface{})
	Headers() *Headers
	Properties() *Properties
}
