////go:build components_mqtt_v5 || components_all

package v5

type Config struct {
	Remaining string  `mapstructure:"remaining"`
	Broker    string  `mapstructure:"broker"`
	ClientID  string  `mapstructure:"clientId"`
	Retained  bool    `mapstructure:"retained"`
	QoS       byte    `mapstructure:"qus"`
	Username  string  `mapstructure:"username"`
	Password  string  `mapstructure:"password"`
	Keepalive *uint16 `mapstructure:"keepalive"`
}
