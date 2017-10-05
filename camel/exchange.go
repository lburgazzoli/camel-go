package camel

// Exchange --
type Exchange struct {
	body    interface{}
	headers map[string]interface{}
}

// ==========================
//
// Initialize an exchnage
//
// ==========================

// NewExchange --
func NewExchange() *Exchange {
	return &Exchange{
		body:    nil,
		headers: make(map[string]interface{}),
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

// SetBody --
func (exchange *Exchange) SetBody(body interface{}) {
	exchange.body = body
}

// Header --
func (exchange *Exchange) Header(name string) interface{} {
	return exchange.headers[name]
}

// SetHeader --
func (exchange *Exchange) SetHeader(name string, value interface{}) {
	exchange.headers[name] = value
}
