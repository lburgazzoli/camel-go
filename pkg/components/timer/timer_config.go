////go:build components_timer || components_all

package timer

import "time"

type Config struct {
	Remaining   string        `mapstructure:"remaining"`
	Period      time.Duration `mapstructure:"period"`
	RepeatCount uint64        `mapstructure:"repeatCount"`
}
