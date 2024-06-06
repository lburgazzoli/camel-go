package log

import (
	"log/slog"

	"github.com/asynkron/protoactor-go/actor"
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/errors"
)

func NewLogHandler(ctx camel.Context) errors.Handler {
	return &Log{
		context: ctx,
	}
}

type Log struct {
	context camel.Context
}

func (a *Log) HandleFailure(
	_ *actor.ActorSystem,
	_ actor.Supervisor,
	_ *actor.PID,
	_ *actor.RestartStatistics,
	reason interface{},
	message interface{},
) {
	a.context.Logger().Info(
		"failure",
		slog.Any("reason", reason),
		slog.Any("message", message),
	)
}
