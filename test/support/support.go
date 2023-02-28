package support

import (
	"testing"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/stretchr/testify/assert"
)

func Reify(t *testing.T, c camel.Context, r processors.Reifyable) (string, error) {
	t.Helper()

	id, err := r.Reify(nil, c)
	assert.Nil(t, err)
	assert.NotNil(t, id)

	return id, err
}
