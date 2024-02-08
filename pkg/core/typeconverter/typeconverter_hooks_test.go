package typeconverter

import (
	"testing"

	"github.com/lburgazzoli/camel-go/pkg/api"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawJsonConversion(t *testing.T) {

	t.Run("string_to_rawjson", func(t *testing.T) {

		var in string
		var out api.RawJSON

		in = `{ "foo": "bar" }`

		tc, err := NewDefaultTypeConverter()
		require.NoError(t, err)

		ok, err := tc.Convert(in, &out)
		require.NoError(t, err)
		assert.True(t, ok)

		assert.Equal(t, "bar", out["foo"])
	})

	t.Run("bytes_to_rawjson", func(t *testing.T) {

		var in []byte
		var out api.RawJSON

		in = []byte(`{ "foo": "bar" }`)

		tc, err := NewDefaultTypeConverter()
		require.NoError(t, err)

		ok, err := tc.Convert(in, &out)
		require.NoError(t, err)
		assert.True(t, ok)

		assert.Equal(t, "bar", out["foo"])
	})
}
