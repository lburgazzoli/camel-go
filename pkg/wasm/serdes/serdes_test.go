package serdes

import (
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerde(t *testing.T) {
	m, err := message.New()
	assert.Nil(t, err)

	encoded, err := Encode(m)
	assert.Nil(t, err)
	assert.NotEmpty(t, encoded)

	decoded, err := Decode(encoded)
	assert.Nil(t, err)

	assert.Equal(t, m.GetID(), decoded.GetID())
	assert.Equal(t, m.GetTime().UnixMilli(), decoded.GetTime().UnixMilli())
}
