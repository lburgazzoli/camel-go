package log

import (
	"testing"

	"github.com/lburgazzoli/camel-go/camel"
	zl "github.com/rs/zerolog"

	"github.com/stretchr/testify/assert"
)

func TestLogLevelAsString(t *testing.T) {
	context := camel.NewContext()

	component := NewComponent()
	component.SetContext(context)

	options := make(map[string]interface{})
	options["level"] = "debug"

	endpoint, err := component.CreateEndpoint("test-log", options)

	assert.NoError(t, err)
	assert.NotNil(t, endpoint)

	le, ok := endpoint.(*logEndpoint)
	assert.True(t, ok)
	assert.Equal(t, zl.DebugLevel, le.level)
	assert.Equal(t, "test-log", le.name)
}

func TestLogLevelAsUint8(t *testing.T) {
	context := camel.NewContext()

	component := NewComponent()
	component.SetContext(context)

	options := make(map[string]interface{})
	options["level"] = uint8(zl.WarnLevel)

	endpoint, err := component.CreateEndpoint("test-log", options)

	assert.NoError(t, err)
	assert.NotNil(t, endpoint)

	le, ok := endpoint.(*logEndpoint)
	assert.True(t, ok)
	assert.Equal(t, zl.WarnLevel, le.level)
	assert.Equal(t, "test-log", le.name)
}

func TestLogLevelAsLevel(t *testing.T) {
	context := camel.NewContext()

	component := NewComponent()
	component.SetContext(context)

	options := make(map[string]interface{})
	options["level"] = zl.FatalLevel

	endpoint, err := component.CreateEndpoint("test-log", options)

	assert.NoError(t, err)
	assert.NotNil(t, endpoint)

	le, ok := endpoint.(*logEndpoint)
	assert.True(t, ok)
	assert.Equal(t, zl.FatalLevel, le.level)
	assert.Equal(t, "test-log", le.name)
}
