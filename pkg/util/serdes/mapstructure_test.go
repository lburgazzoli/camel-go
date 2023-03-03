package serdes

import (
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleEndpoint(t *testing.T) {
	var in string
	var out int

	dec, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		Result:           &out,

		// custom hooks
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
		),
	})

	in = "1"

	assert.Nil(t, err)
	assert.Nil(t, dec.Decode(&in))
	assert.Equal(t, 1, out)

}
