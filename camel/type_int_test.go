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

	"github.com/stretchr/testify/assert"
)

// ==========================
//
// Int converter
//
// ==========================

func testStringToInt(t *testing.T, value string, expectedResult interface{}) {
	expectedType := reflect.TypeOf(expectedResult)

	r, e := ToIntConverter(value, expectedType)

	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, expectedResult, r)
	assert.IsType(t, expectedResult, r)
}

func TestStringToIntConverter(t *testing.T) {
	testStringToInt(t, "1", int(1))
	testStringToInt(t, "1", int8(1))
	testStringToInt(t, "1", int16(1))
	testStringToInt(t, "1", int32(1))
	testStringToInt(t, "1", int64(1))
}

func TestStringToIntConverterWithInvalidType(t *testing.T) {
	r, e := ToIntConverter("1", TypeString)

	assert.Nil(t, r)
	assert.Error(t, e)
}
