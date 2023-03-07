////go:build components_mqtt || components_all

package mqtt

import (
	"github.com/lburgazzoli/camel-go/pkg/components"
)

func init() {
	components.Factories[Scheme] = NewComponent
}
