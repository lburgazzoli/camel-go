////go:build components_timer || components_all

package timer

import "time"

type Config struct {
	Interval time.Duration `mapstructure:"interval"`
}
