package camel

type Component interface {
	Service

	Process(message string)
}
