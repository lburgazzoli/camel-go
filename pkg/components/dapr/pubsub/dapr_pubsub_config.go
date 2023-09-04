// //go:build components_dapr_pubsub || components_all

package pubsub

type Config struct {
	Remaining string `mapstructure:"remaining"`
	Raw       bool   `mapstructure:"raw"`
}
