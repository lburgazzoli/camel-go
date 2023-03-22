package errors

import "github.com/asynkron/protoactor-go/actor"

// Handler ---
// TODO: this likely need to be moved to the API package but
//
//	requires some more refinement
type Handler interface {
	actor.SupervisorStrategy
}
