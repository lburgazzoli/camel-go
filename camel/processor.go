package camel

// Processor --
type Processor func(*Exchange)

// Trasformer --
type Trasformer func(*Exchange) *Exchange

// Predicate
type Predicate func(*Exchange) bool
