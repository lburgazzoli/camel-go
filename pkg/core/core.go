package core

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/api"
	context2 "github.com/lburgazzoli/camel-go/pkg/core/context"
)

// NewContext ---
func NewContext(ctx context.Context) api.Context {
	return context2.NewDefaultContext(ctx)
}
