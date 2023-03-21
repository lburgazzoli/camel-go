package wasm

import (
	"context"
	"os"
	"path"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/util/registry"
	"github.com/lburgazzoli/camel-go/pkg/wasm"
	"github.com/lburgazzoli/camel-go/pkg/wasm/serdes"
)

type Wasm struct {
	Path  string `yaml:"path"`
	Image string `yaml:"image,omitempty"`
}

func (l *Wasm) Processor(ctx context.Context, camelContext camel.Context) (camel.Processor, error) {
	if l.Path == "" {
		return nil, camelerrors.MissingParameterf("wasm.path", "failure configuring wasm processor")
	}

	rootPath := ""

	if l.Image != "" {
		fp, err := registry.Pull(ctx, l.Image)
		if err != nil {
			return nil, err
		}

		rootPath = fp
	}

	defer func() {
		if rootPath != "" {
			_ = os.RemoveAll(rootPath)
		}
	}()

	fd, err := os.Open(path.Join(rootPath, l.Path))
	if err != nil {
		return nil, err
	}

	defer func() { _ = fd.Close() }()

	r, err := wasm.NewRuntime(ctx, wasm.Options{})
	if err != nil {
		return nil, err
	}

	f, err := r.Load(ctx, "process", fd)
	if err != nil {
		return nil, err
	}

	p := func(ctx context.Context, m camel.Message) error {
		encoded, err := serdes.EncodeMessage(m)
		if err != nil {
			return err
		}

		data, err := f.Invoke(ctx, encoded)
		if err != nil {
			return err
		}

		return serdes.DecodeMessageTo(data, m)
	}

	return p, nil
}
