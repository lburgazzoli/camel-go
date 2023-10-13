package containers

import (
	"fmt"
	"log/slog"

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

func NewSlogLogConsumer(r *testcontainers.ContainerRequest) testcontainers.LogConsumer {
	name := r.Name
	if name == "" {
		name = r.Image
	}

	return &SlogLogConsumer{
		Name: name,
	}
}

type SlogLogConsumer struct {
	Name string
}

func (g *SlogLogConsumer) Accept(l testcontainers.Log) {
	slog.Default().WithGroup("container").Info(
		string(l.Content),
		slog.String("type", l.LogType),
		slog.String("name", g.Name),
	)
}
