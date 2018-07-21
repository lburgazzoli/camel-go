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
package camel

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ==========================
//
// Duration converter
//
// ==========================

func testStringToDuration(t *testing.T, value string, expectedResult interface{}) {
	expectedType := reflect.TypeOf(expectedResult)

	r, e := ToDurationConverter(value, expectedType)

	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, expectedResult, r)
	assert.IsType(t, expectedResult, r)
}

func TestStringToDurationConverter(t *testing.T) {
	testStringToDuration(t, "1s", time.Duration(1*time.Second))
}
