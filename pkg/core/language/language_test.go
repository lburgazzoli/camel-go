package language

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

const jqFullExpression = `
jq: 
  expression: '.foo'
`

const jqInlineExpression = `
jq: .bar
`

func TestJq(t *testing.T) {

	t.Run("full", func(t *testing.T) {
		l := Language{}

		err := yaml.Unmarshal([]byte(jqFullExpression), &l)
		require.NoError(t, err)
		require.Equal(t, ".foo", l.Jq.Expression)
	})

	t.Run("inline", func(t *testing.T) {
		l := Language{}

		err := yaml.Unmarshal([]byte(jqInlineExpression), &l)
		require.NoError(t, err)
		require.Equal(t, ".bar", l.Jq.Expression)
	})
}
