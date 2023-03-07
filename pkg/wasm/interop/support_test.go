package interop

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupport(t *testing.T) {
	m1 := Message{}

	data := make([]byte, 0, 1024)

	err := EncodeMessage(m1, data)
	assert.Nil(t, err)

	m2 := DecodeMessage(data)
	assert.Equal(t, m1.ID, m2.ID)
}
