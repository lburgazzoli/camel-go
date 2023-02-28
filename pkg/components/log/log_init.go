////go:build components_log || components_all

package log

import (
	"github.com/lburgazzoli/camel-go/pkg/components"
)

func init() {
	components.Factories[Scheme] = NewComponent
}
