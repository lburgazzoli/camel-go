////go:build components_http|| components_all

package http

type Config struct {
	URI       string `mapstructure:"uri"`
	Remaining string `mapstructure:"remaining"`
	Method    string `mapstructure:"method"`
}
