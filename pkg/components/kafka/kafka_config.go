////go:build components_kafka || components_all

package kafka

type Config struct {
	Brokers       string `mapstructure:"brokers"`
	Topics        string `mapstructure:"topics"`
	ConsumerGroup string `mapstructure:"consumerGroup"`
}
