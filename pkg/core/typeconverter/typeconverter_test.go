package typeconverter

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestSimpleConversion(t *testing.T) {
	in := "1"
	out := 0

	tc, err := NewDefaultTypeConverter()

	require.NoError(t, err)

	ok, err := tc.Convert(&in, &out)
	require.NoError(t, err)
	assert.True(t, ok)

	assert.Equal(t, 1, out)
}
