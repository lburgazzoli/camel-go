////go:build components_mqtt_v5 || components_all

package v5

import (
	"github.com/lburgazzoli/camel-go/pkg/components"
)

func init() {
	components.Factories[Scheme] = NewComponent
	components.Factories[SchemeAlias] = NewComponent
}
