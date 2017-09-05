package camel

type Context interface {
	Service

	GetComponent(name string) (Component, error)
}
