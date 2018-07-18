package api

// ==========================
//
//
//
// ==========================

// Subscription --
type Subscription interface {
	Cancel()
}

// ==========================
//
//
//
// ==========================

// Publisher --
type Publisher interface {
	Publish(Exchange)
}

// Subscriber --
type Subscriber interface {
	Subscribe(func(Exchange)) Subscription
}

// Processor --
type Processor interface {
	Publisher
	Subscriber
}

// ==========================
//
//
//
// ==========================

// ProcessingService --
type ProcessingService interface {
	Service
	Processor
}
