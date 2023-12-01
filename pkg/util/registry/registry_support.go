package registry

import (
	"os"
)

func WasmContainerImage() string {
	image := os.Getenv("WASM_CONTAINER_IMAGE")
	if image == "" {
		image = "quay.io/lburgazzoli/camel-go-wasm:latest"
	}

	return image
}
