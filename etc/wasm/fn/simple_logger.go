//go:build tinygo.wasm

// nolint
package main

import (
	"context"
	"fmt"

	"github.com/lburgazzoli/camel-go/pkg/wasm/plugin/processor"
)

// main is required for TinyGo to compile to Wasm.
func main() {
	processor.RegisterProcessors(Process)
}

func Process(_ context.Context, request *processor.Message) (*processor.Message, error) {
	fmt.Println("Processing message ", request.ID)

	return request, nil
}
