package containers

import (
	"fmt"
	"log/slog"
	"os"

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

func NewSlogLogConsumer(name string) testcontainers.LogConsumer {
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

func NewSlogLogger(name string) testcontainers.Logging {
	return &SlogLogger{
		Name: name,
		l:    Log.WithGroup(name),
	}
}

type SlogLogger struct {
	Name string
	l    *slog.Logger
}

func (g *SlogLogger) Printf(format string, v ...interface{}) {
	g.l.Info(
		fmt.Sprintf(format, v...),
	)
}
