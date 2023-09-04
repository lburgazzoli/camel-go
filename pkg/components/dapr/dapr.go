package dapr

import (
	"fmt"
	"os"
)

const (
	// TODO: this should be moved to generic dapr component.

	DefaultPort     = 8082
	DefaultProtocol = "http"
	DefaultPortName = "dapr"

	EnvVarAddress = "CAMEL_DAPR_ADDRESS"

	AnnotationAEnabled    = "dapr.io/enabled"
	AnnotationAppID       = "dapr.io/app-id"
	AnnotationAppPort     = "dapr.io/app-port"
	AnnotationAppProtocol = "dapr.io/app-protocol"

	MetaRawPayload = "rawPayload"
)

func Address() string {
	address := os.Getenv(EnvVarAddress)
	if address == "" {
		address = fmt.Sprintf(":%d", DefaultPort)
	}

	return address
}
