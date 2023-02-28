package core

import (
	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/context"
)

func NewContext() camel.Context {
	return context.NewDefaultContext()
}
