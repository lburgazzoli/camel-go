package support

import (
	"testing"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/stretchr/testify/assert"
)

func Reify(t *testing.T, c api.Context, r processors.Reifyable) *actor.PID {
	t.Helper()

	pid, err := r.Reify(c)
	assert.Nil(t, err)
	assert.NotNil(t, pid)

	return pid
}
