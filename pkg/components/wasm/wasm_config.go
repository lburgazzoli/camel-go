////go:build components_wasm || components_all

package wasm

type Config struct {
	Remaining string `mapstructure:"remaining"`
	Path      string `mapstructure:"path"`
	Image     string `mapstructure:"image,omitempty"`

	Other map[string]string `mapstructure:",remain"`
}
