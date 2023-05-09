////go:build components_dapr || components_all

package dapr

type Config struct {
	Remaining string `mapstructure:"remaining"`
}
