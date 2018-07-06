package camel

import (
	"log"
	"reflect"

	"github.com/lburgazzoli/camel-go/api"
)

// ==========================
//
// Initialize an exchnage
//
// ==========================

// NewExchange --
func NewExchange(context api.Context) api.Exchange {
	converter := context.TypeConverter()

	return &DefaultExchange{
		converter: converter,
		headers: api.Headers{
			Registry: api.NewInMemoryRegistry(converter),
		},
		properties: api.Properties{
			Registry: api.NewInMemoryRegistry(converter),
		},
	}
}

// ==========================
//
//
//
// ==========================

// DefaultExchange --
type DefaultExchange struct {
	body       interface{}
	converter  api.TypeConverter
	headers    api.Headers
	properties api.Properties
}

// Body --
func (exchange *DefaultExchange) Body() interface{} {
	return exchange.body
}

// BodyAs --
func (exchange *DefaultExchange) BodyAs(asType reflect.Type) interface{} {
	answer := exchange.Body()

	if answer != nil {
		result, err := exchange.converter(answer, asType)

		if err != nil {
			log.Fatalf("unable to covert body to required type: %v", asType)
		}

		return result
	}

	return answer
}

// SetBody --
func (exchange *DefaultExchange) SetBody(body interface{}) {
	exchange.body = body
}

// Headers --
func (exchange *DefaultExchange) Headers() *api.Headers {
	return &exchange.headers
}

// Properties --
func (exchange *DefaultExchange) Properties() *api.Properties {
	return &exchange.properties
}
