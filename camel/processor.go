package camel

// Processor --
type Processor func(*Exchange)

// Trasformer --
type Trasformer func(*Exchange) *Exchange
