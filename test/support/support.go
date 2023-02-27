package support

import (
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/stretchr/testify/assert"
)

func Reify(t *testing.T, c api.Context, r processors.Reifyable) (string, error) {
	t.Helper()

	id, err := r.Reify(c)
	assert.Nil(t, err)
	assert.NotNil(t, id)

	return id, err
}
