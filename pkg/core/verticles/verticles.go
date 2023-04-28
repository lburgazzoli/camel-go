package verticles

import (
	"context"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
)

func Spawn(ac actor.Context, verticle camel.Verticle) (*actor.PID, error) {

	f := func() actor.Actor { return verticle }
	p := actor.PropsFromProducer(f)

	pid, err := ac.SpawnNamed(p, verticle.ID())
	if err != nil {
		return nil, err
	}

	return pid, nil
}

func NewContext(cc camel.Context, ac actor.Context) context.Context {
	c := context.Background()
	c = context.WithValue(c, camel.ContextKeyCamelContext, cc)
	c = context.WithValue(c, camel.ContextKeyActorContext, ac)

	return c
}

func Contains(pids []*actor.PID, pid *actor.PID) bool {
	for _, c := range pids {
		if c.Equal(pid) {
			return true
		}
	}

	return false
}
