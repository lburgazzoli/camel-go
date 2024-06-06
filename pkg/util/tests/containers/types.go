package containers

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/imdario/mergo"
	"github.com/testcontainers/testcontainers-go"
)

const (
	FileModeRead   os.FileMode = 0o600 // For secret files.
	FileModeShared os.FileMode = 0o644 // For normal files.
	FileModeExec   os.FileMode = 0o755 // For directory or execute files.
)

var (
	Log = slog.Default().WithGroup("container")
)

type OverrideContainerRequestOption func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest

var NoopOverrideContainerRequest = func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
	return req
}

func OverrideContainerRequest(r testcontainers.ContainerRequest) func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
	return func(req testcontainers.ContainerRequest) testcontainers.ContainerRequest {
		if err := mergo.Merge(&req, r, mergo.WithOverride); err != nil {
			slog.Default().WithGroup("container").Info(
				fmt.Sprintf("error merging container request %v. Keeping the default one: %v", err, req),
			)

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
		l:    Log.WithGroup(name),
	}
}

type SlogLogConsumer struct {
	Name string
	l    *slog.Logger
}

func (g *SlogLogConsumer) Accept(l testcontainers.Log) {
	g.l.Info(
		string(l.Content),
		slog.String("stream", l.LogType),
	)
}
