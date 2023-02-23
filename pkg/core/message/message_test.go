package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	m, err := New()

	assert.Nil(t, err)
	assert.NotNil(t, m.GetID())
	assert.NotNil(t, m.GetTime())
	assert.NotZero(t, m.GetTime())
}
