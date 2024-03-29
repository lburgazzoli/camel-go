package route

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

const loaderIn = `
- route:
    id: "foo" 
`

func TestLoader(t *testing.T) {
	data := strings.NewReader(loaderIn)

	routes, err := Load(data)

	require.NoError(t, err)
	assert.Len(t, routes, 1)
	assert.NotEmpty(t, routes[0].ID)
}
