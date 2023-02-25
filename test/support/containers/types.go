package containers

import (
	"fmt"

	"github.com/imdario/mergo"
	"github.com/testcontainers/testcontainers-go"
)

type OverrideContainerRequestOption func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest

var NoopOverrideContainerRequest = func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
	return req
}

func OverrideContainerRequest(r testcontainers.ContainerRequest) func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
	return func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
		if err := mergo.Merge(&req, r, mergo.WithOverride); err != nil {
			fmt.Printf("error merging container request %v. Keeping the default one: %v", err, req)
			return req
		}

		return req
	}
}

type SysOutLogConsumer struct {
}

func (g *SysOutLogConsumer) Accept(l testcontainers.Log) {
	// _, _ = fmt.Fprintf(os.Stdout, string(l.Content))
}
