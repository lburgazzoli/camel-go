// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"testing"

	"github.com/lburgazzoli/camel-go/camel"
	"github.com/rs/zerolog"

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
	assert.Equal(t, zerolog.DebugLevel, le.level)
	assert.Equal(t, "test-log", le.logger)
}

func TestLogLevelAsUint8(t *testing.T) {
	context := camel.NewContext()

	component := NewComponent()
	component.SetContext(context)

	options := make(map[string]interface{})
	options["level"] = uint8(zerolog.WarnLevel)

	endpoint, err := component.CreateEndpoint("test-log", options)

	assert.NoError(t, err)
	assert.NotNil(t, endpoint)

	le, ok := endpoint.(*logEndpoint)
	assert.True(t, ok)
	assert.Equal(t, zerolog.WarnLevel, le.level)
	assert.Equal(t, "test-log", le.logger)
}

func TestLogLevelAsLevel(t *testing.T) {
	context := camel.NewContext()

	component := NewComponent()
	component.SetContext(context)

	options := make(map[string]interface{})
	options["level"] = zerolog.FatalLevel

	endpoint, err := component.CreateEndpoint("test-log", options)

	assert.NoError(t, err)
	assert.NotNil(t, endpoint)

	le, ok := endpoint.(*logEndpoint)
	assert.True(t, ok)
	assert.Equal(t, zerolog.FatalLevel, le.level)
	assert.Equal(t, "test-log", le.logger)
}

func TestLogNameOverride(t *testing.T) {
	context := camel.NewContext()

	component := NewComponent()
	component.SetContext(context)

	options := make(map[string]interface{})
	options["level"] = zerolog.FatalLevel
	options["logger"] = "override"

	endpoint, err := component.CreateEndpoint("test-log", options)

	assert.NoError(t, err)
	assert.NotNil(t, endpoint)

	le, ok := endpoint.(*logEndpoint)
	assert.True(t, ok)
	assert.Equal(t, zerolog.FatalLevel, le.level)
	assert.Equal(t, "override", le.logger)
}
