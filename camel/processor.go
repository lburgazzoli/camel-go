package camel

// Processor --
type Processor interface {
	Process(exchange *Exchange)
}
