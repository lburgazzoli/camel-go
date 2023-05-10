package language

import (
	"github.com/lburgazzoli/camel-go/pkg/core/language/constant"
	"github.com/lburgazzoli/camel-go/pkg/core/language/jq"
	"github.com/lburgazzoli/camel-go/pkg/core/language/mustache"
	"github.com/lburgazzoli/camel-go/pkg/core/language/wasm"
)

type Language struct {
	Jq       *jq.Jq             `yaml:"jq,omitempty"`
	Mustache *mustache.Mustache `yaml:"mustache,omitempty"`
	Wasm     *wasm.Wasm         `yaml:"wasm,omitempty"`
	Constant *constant.Constant `yaml:"constant,omitempty"`
}
