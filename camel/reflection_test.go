package camel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==========================
//
// Duration converter
//
// ==========================

type MyTarget struct {
	Field1 string
	field2 string
	Field3 string
}

func (target *MyTarget) SetF2(value string) {
	target.field2 = value
}

func TestSetProperty(t *testing.T) {
	my := MyTarget{
		Field1: "f1",
		field2: "f2",
		Field3: "f3",
	}

	context := NewContext()
	r1 := SetProperty(context, &my, "Field1", "new-value-1")
	r2 := SetProperty(context, &my, "f2", "new-value-2")
	r3 := SetProperty(context, &my, "f3", "new-value-3")

	assert.True(t, r1)
	assert.Equal(t, "new-value-1", my.Field1)

	assert.True(t, r2)
	assert.Equal(t, "new-value-2", my.field2)

	assert.False(t, r3)
	assert.Equal(t, "f3", my.Field3)
}
