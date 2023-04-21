//go:build tinygo.wasm

// nolint
package main

import (
	"context"
	"strings"

	"github.com/lburgazzoli/camel-go/pkg/wasm/plugin/processor"
)

// main is required for TinyGo to compile to Wasm.
func main() {
	processor.RegisterProcessors(MyProcessor{})
}

type MyProcessor struct{}

func (m MyProcessor) Process(_ context.Context, request *processor.Message) (*processor.Message, error) {
	request.Data = []byte(strings.ToUpper(string(request.Data)))

	return request, nil
}
