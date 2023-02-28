////go:build components_timer || components_all

package timer

import "time"

type Config struct {
	Remaining string        `mapstructure:"remaining"`
	Interval  time.Duration `mapstructure:"interval"`
}
