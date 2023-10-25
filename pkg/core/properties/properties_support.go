package properties

import (
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

func newKoanf() (*koanf.Koanf, error) {
	k := koanf.New(delimiter)

	envProvider := env.Provider(envPrefix, delimiter, func(s string) string {
		s = strings.TrimPrefix(s, envPrefix)
		s = strings.ReplaceAll(s, "_", delimiter)
		s = strings.ToLower(s)

		return s
	})

	err := k.Load(envProvider, nil)
	if err != nil {
		return nil, err
	}

	return k, nil
}
