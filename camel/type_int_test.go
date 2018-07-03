package camel

import (
	"reflect"
	"testing"

	"github.com/lburgazzoli/camel-go/api"
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
	r, e := ToIntConverter("1", api.TypeString)

	assert.Nil(t, r)
	assert.Error(t, e)
}
