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
	converter := ToIntConverter{}

	r, e := converter.Convert(value, expectedType)

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
	converter := ToIntConverter{}

	r, e := converter.Convert("1", TypeString)

	assert.Nil(t, r)
	assert.Error(t, e)
}
