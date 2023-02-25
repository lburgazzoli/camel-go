package core

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	context2 "github.com/lburgazzoli/camel-go/pkg/core/context"
)

func NewContext() api.Context {
	return context2.NewDefaultContext()
}
