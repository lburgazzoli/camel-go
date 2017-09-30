package camel

// Exchange --
type Exchange interface {

	Body() interface{}

	SetBody(body interface{})

	Header(name string) interface{}

	SetHeader(name string, value interface{})
}