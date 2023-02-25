////go:build components_wasm || components_all

package wasm

import (
	"github.com/lburgazzoli/camel-go/pkg/components"
)

func init() {
	components.Factories[Scheme] = NewComponent
}
