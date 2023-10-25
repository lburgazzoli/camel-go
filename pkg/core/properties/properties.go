package properties

import (
	"fmt"
	"maps"
	"os"
	"path/filepath"

	"github.com/knadh/koanf/providers/confmap"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/pkg/errors"

	"github.com/knadh/koanf/v2"
	"github.com/lburgazzoli/camel-go/pkg/api"
)

const delimiter = "."
const envPrefix = "CAMEL_"

func NewDefaultProperties() (api.Properties, error) {
	k, err := newKoanf()
	if err != nil {
		return nil, err
	}

	p := defaultProperties{
		konf: k,
	}

	return &p, nil
}

type defaultProperties struct {
	konf *koanf.Koanf
}

func (r *defaultProperties) Add(source map[string]any) error {
	err := r.konf.Load(
		confmap.Provider(source, "."),
		nil,
	)

	if err != nil {
		return errors.Wrapf(err, "error loading config")
	}

	return nil
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

func (r *defaultProperties) Parameters() api.Parameters {
	return r.konf.All()
}

func (r *defaultProperties) View(path string) api.PropertiesResolver {
	return &defaultProperties{
		konf: r.konf.Cut(path),
	}
}

func (r *defaultProperties) Expand(in string) (string, bool) {
	answer := os.Expand(in, func(s string) string {
		v := r.konf.Get(s)
		if v == nil {
			return ""
		}

		if i, ok := v.(string); ok {
			return i
		}

		return ""
	})

	if answer == "" {
		answer = in
	}

	return answer, answer != in
}

func (r *defaultProperties) ExpandAll(in map[string]any) map[string]any {
	answer := maps.Clone(in)

	for k := range answer {
		if v, ok := answer[k].(string); ok {
			answer[k], _ = r.Expand(v)
		}
	}

	return answer
}

func (r *defaultProperties) Merge(in map[string]any) (api.PropertiesResolver, error) {
	konf := r.konf.Copy()
	for k, v := range in {
		if err := konf.Set(k, v); err != nil {
			return nil, err
		}
	}

	return &defaultProperties{konf: konf}, nil
}
