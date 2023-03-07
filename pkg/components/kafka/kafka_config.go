////go:build components_kafka || components_all

package kafka

type Config struct {
	Remaining     string `mapstructure:"remaining"`
	Brokers       string `mapstructure:"brokers"`
	ConsumerGroup string `mapstructure:"consumerGroup"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
}
