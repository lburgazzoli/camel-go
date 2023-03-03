package properties

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/pkg/errors"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/lburgazzoli/camel-go/pkg/api"
)

const delimiter = "."
const envPrefix = "CAMEL_"
const placeholderBegin = "{{"
const placeholderEnd = "}}"

func NewDefaultProperties() (api.Properties, error) {
	p := defaultProperties{
		konf: koanf.New(delimiter),
	}

	envProvider := env.Provider(envPrefix, delimiter, func(s string) string {
		s = strings.TrimPrefix(s, envPrefix)
		s = strings.ReplaceAll(s, "_", delimiter)
		s = strings.ToLower(s)

		return s
	})

	err := p.konf.Load(envProvider, nil)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

type defaultProperties struct {
	konf *koanf.Koanf
}

func (r *defaultProperties) AddSource(path string) error {
	fi, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	// TODO: check for symlinks
	if fi.IsDir() {
		return fmt.Errorf("file %s is a directory", path)
	}

	switch filepath.Ext(path) {
	case ".yaml":
		if err := r.konf.Load(file.Provider(path), yaml.Parser()); err != nil {
			return errors.Wrapf(err, "error loading config")
		}
	case ".json":
		if err := r.konf.Load(file.Provider(path), json.Parser()); err != nil {
			return errors.Wrapf(err, "error loading config")
		}
	case ".toml":
		if err := r.konf.Load(file.Provider(path), toml.Parser()); err != nil {
			return errors.Wrapf(err, "error loading config")
		}
	}

	return nil
}

func (r *defaultProperties) String(data string) string {
	// TODO: must handle multiple placeholders, recursion, etc.

	key := data

	if strings.HasPrefix(data, placeholderBegin) && strings.HasSuffix(data, placeholderEnd) {
		key = strings.TrimPrefix(key, placeholderBegin)
		key = strings.TrimSuffix(key, placeholderEnd)
	}

	result := r.konf.String(key)
	if result != "" {
		return result
	}

	return data
}