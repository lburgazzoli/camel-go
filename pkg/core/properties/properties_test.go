package properties

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProperties(t *testing.T) {
	p1, err := NewDefaultProperties()
	require.NoError(t, err)
	require.NotNil(t, t, p1)

	require.NoError(
		t,
		p1.Add(map[string]any{
			"key1":             "1",
			"nested.key2":      "2",
			"deep.nested.key3": "3",
		}),
	)

	v, ok := p1.Expand("${key1}")
	require.True(t, ok)
	require.Equal(t, "1", v)

	v, ok = p1.Expand("${key1}-${nested.key2}")
	require.True(t, ok)
	require.Equal(t, "1-2", v)

	v, ok = p1.Expand("${nested.key2}")
	require.True(t, ok)
	require.Equal(t, "2", v)

	p2 := p1.View("nested")

	v, ok = p2.Expand("${key2}")
	require.True(t, ok)
	require.Equal(t, "2", v)

	v, ok = p2.Expand("${key1}")
	require.False(t, ok)
	require.Equal(t, "${key1}", v)

	p3 := p1.View("deep.nested")

	v, ok = p3.Expand("${key3}")
	require.True(t, ok)
	require.Equal(t, "3", v)

	v, ok = p3.Expand("${nested.key2}")
	require.False(t, ok)
	require.Equal(t, "${nested.key2}", v)
}
