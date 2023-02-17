package message

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	m := New()

	assert.NotNil(t, m.ID())
	assert.NotNil(t, m.Time())
	assert.NotZero(t, m.Time())
}
