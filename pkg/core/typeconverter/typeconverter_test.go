package typeconverter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleConversion(t *testing.T) {
	in := "1"
	out := 0

	tc, err := NewDefaultTypeConverter()
	assert.Nil(t, err)

	ok, err := tc.Convert(&in, &out)
	assert.Nil(t, err)
	assert.True(t, ok)

	assert.Equal(t, 1, out)
}
