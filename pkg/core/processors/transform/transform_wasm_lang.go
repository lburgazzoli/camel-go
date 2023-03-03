////go:build steps_transform || steps_all

package transform

import (
	"context"
	"os"
	"path"

	"github.com/lburgazzoli/camel-go/pkg/util/registry"
	"github.com/lburgazzoli/camel-go/pkg/wasm/serdes"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/wasm"
)

type LanguageWasm struct {
	Path  string `yaml:"path"`
	Image string `yaml:"image,omitempty"`
}

func newWasmProcessor(ctx context.Context, definition *LanguageWasm) (languageProcessor, error) {
	if definition.Path == "" {
		return nil, camelerrors.MissingParameterf("wasm.path", "failure processing %s", TAG)
	}

	rootPath := ""

	if definition.Image != "" {
		fp, err := registry.Pull(ctx, definition.Image)
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

	fd, err := os.Open(path.Join(rootPath, definition.Path))
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

	return &wasmProcessor{r: r, f: f}, nil
}

type wasmProcessor struct {
	r *wasm.Runtime
	f *wasm.Function
}

func (p *wasmProcessor) Process(ctx context.Context, m camel.Message) (camel.Message, error) {
	encoded, err := serdes.EncodeMessage(m)
	if err != nil {
		return nil, err
	}

	data, err := p.f.Invoke(ctx, encoded)
	if err != nil {
		return nil, err
	}

	return serdes.DecodeMessage(data)
}
