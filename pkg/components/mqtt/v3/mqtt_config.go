////go:build components_mqtt_v3 || components_all

package v3

type Config struct {
	Remaining string `mapstructure:"remaining"`
	Brokers   string `mapstructure:"brokers"`
	ClientID  string `mapstructure:"clientId"`
	Retained  bool   `mapstructure:"retained"`
	QoS       byte   `mapstructure:"qus"`
}
