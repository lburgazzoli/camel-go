package typeconverter

import (
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/api"

	"github.com/stretchr/testify/assert"
)

func TestRawJsonConversion(t *testing.T) {

	t.Run("string_to_rawjson", func(t *testing.T) {

		var in string
		var out api.RawJSON

		in = `{ "foo": "bar" }`

		tc, err := NewDefaultTypeConverter()
		assert.Nil(t, err)

		ok, err := tc.Convert(in, &out)
		assert.Nil(t, err)
		assert.True(t, ok)

		assert.Equal(t, "bar", out["foo"])
	})

	t.Run("bytes_to_rawjson", func(t *testing.T) {

		var in []byte
		var out api.RawJSON

		in = []byte(`{ "foo": "bar" }`)

		tc, err := NewDefaultTypeConverter()
		assert.Nil(t, err)

		ok, err := tc.Convert(in, &out)
		assert.Nil(t, err)
		assert.True(t, ok)

		assert.Equal(t, "bar", out["foo"])
	})
}
