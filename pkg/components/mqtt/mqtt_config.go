////go:build components_mqtt || components_all

package mqtt

type Config struct {
	Remaining string `mapstructure:"remaining"`
	Brokers   string `mapstructure:"brokers"`
	ClientID  string `mapstructure:"clientId"`
}
