package module

import (
	"fmt"
	"os"
	"plugin"

	zlog "github.com/rs/zerolog/log"
)

// LoadSymbol --
func LoadSymbol(path string, symbol string) (interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("Path is empty")
	}

	if symbol == "" {
		return nil, fmt.Errorf("Symbol is empty")
	}

	zlog.Debug().Msgf("load symbol \"%s\" from plugin %s", symbol, path)

	location := os.ExpandEnv(path)
	_, err := os.Stat(location)

	if os.IsNotExist(err) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	plug, err := plugin.Open(location)
	if err != nil {
		return nil, fmt.Errorf("failed to open plugin %s: %v", location, err)
	}

	answer, err := plug.Lookup(symbol)
	if err != nil {
		return nil, fmt.Errorf("plugin %s does not export symbol \"%s\"", location, symbol)
	}

	return answer, nil
}
