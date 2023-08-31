package core

import (
	"log/slog"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/context"
)

func NewContext(logger *slog.Logger, opts ...context.Option) camel.Context {
	return context.NewDefaultContext(logger, opts...)
}
