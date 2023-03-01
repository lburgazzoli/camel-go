package serdes

import (
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/stretchr/testify/assert"
)

func TestSerde(t *testing.T) {
	m, err := message.New()
	assert.Nil(t, err)

	encoded, err := EncodeMessage(m)
	assert.Nil(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := DecodeMessage(encoded)
	assert.Nil(t, err)

	assert.Equal(t, m.GetID(), decoded.GetID())
	assert.Equal(t, m.GetTime().UnixMilli(), decoded.GetTime().UnixMilli())
}
