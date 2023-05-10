////go:build components_dapr_pubsub || components_all

package pubsub

import (
	"github.com/lburgazzoli/camel-go/pkg/components"
)

func init() {
	components.Factories[Scheme] = NewComponent
}
