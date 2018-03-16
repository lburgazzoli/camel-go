package camel

// Processor --
type Processor func(*Exchange)

// Transformer --
type Trasformer func(*Exchange) *Exchange

// Predicate
type Predicate func(*Exchange) bool
