package camel

type Service interface {
	Start()
	Stop()
}

type Context interface {
	Service

	GetComponent(name string) (Component, error)
}

type Component interface {
	Service

	Init(context Context) error

    Process(message string)
}
