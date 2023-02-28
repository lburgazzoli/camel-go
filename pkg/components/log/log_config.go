////go:build components_log || components_all

package log

type Config struct {
	Remaining string `mapstructure:"remaining"`
}
