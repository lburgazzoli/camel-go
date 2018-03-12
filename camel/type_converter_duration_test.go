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

	r, e := ToDuratioinConverter(value, expectedType)

	assert.NoError(t, e)
	assert.NotNil(t, r)
	assert.Equal(t, expectedResult, r)
	assert.IsType(t, expectedResult, r)
}

func TestStringToDurationConverter(t *testing.T) {
	testStringToDuration(t, "1s", time.Duration(1*time.Second))
}
