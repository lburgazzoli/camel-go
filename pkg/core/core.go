package core

import (
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/context"
	"go.uber.org/zap"
)

var L *zap.Logger

func NewContext(logger *zap.Logger, opts ...context.Option) camel.Context {
	return context.NewDefaultContext(logger, opts...)
}
