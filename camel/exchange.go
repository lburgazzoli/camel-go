package camel

import (
	"log"
	"reflect"
)

// Exchange --
type Exchange struct {
	context    *Context
	body       interface{}
	headers    map[string]interface{}
	properties map[string]interface{}
}

// ==========================
//
// Initialize an exchnage
//
// ==========================

// NewExchange --
func NewExchange(context *Context) *Exchange {
	return &Exchange{
		context:    context,
		body:       nil,
		headers:    make(map[string]interface{}),
		properties: make(map[string]interface{}),
	}
}

// ==========================
//
//
//
// ==========================

// Body --
func (exchange *Exchange) Body() interface{} {
	return exchange.body
}

// BodyAs --
func (exchange *Exchange) BodyAs(asType reflect.Type) interface{} {
	answer := exchange.Body()

	if answer != nil {
		converter := exchange.context.TypeConverter()
		result, err := converter.Convert(answer, asType)

		if err != nil {
			log.Fatalf("Unable to covert body to required type: %v", asType)
		}

		return result
	}

	return answer
}

// SetBody --
func (exchange *Exchange) SetBody(body interface{}) {
	exchange.body = body
}

// Header --
func (exchange *Exchange) Header(name string) interface{} {
	return exchange.headers[name]
}

// HeaderAs --
func (exchange *Exchange) HeaderAs(name string, asType reflect.Type) interface{} {
	answer := exchange.Header(name)

	if answer != nil {
		converter := exchange.context.TypeConverter()
		result, err := converter.Convert(answer, asType)

		if err != nil {
			log.Fatalf("Unable to covert header: %s to required type: %v", name, asType)
		}

		return result
	}

	return answer
}

// HeaderOrDefault --
func (exchange *Exchange) HeaderOrDefault(name string, defaultValue interface{}) interface{} {
	answer := exchange.Header(name)
	if answer == nil {
		answer = defaultValue
	}

	return answer
}

// SetHeader --
func (exchange *Exchange) SetHeader(name string, value interface{}) {
	exchange.headers[name] = value
}

// Property --
func (exchange *Exchange) Property(name string) interface{} {
	return exchange.properties[name]
}

// SetProperty --
func (exchange *Exchange) SetProperty(name string, value interface{}) {
	exchange.properties[name] = value
}
