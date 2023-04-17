package wasm

import (
	"context"
	"os"
	"path"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/util/registry"
	"github.com/lburgazzoli/camel-go/pkg/wasm"
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

	r, err := wasm.NewRuntime(ctx, wasm.Options{})
	if err != nil {
		return nil, err
	}

	f, err := r.Load(ctx, path.Join(rootPath, l.Path))
	if err != nil {
		return nil, err
	}

	p := func(ctx context.Context, m camel.Message) error {
		result, err := f.Invoke(ctx, m)
		if err != nil {
			return err
		}

		_ = m.SetID(result.GetID())
		_ = m.SetSource(result.GetSource())
		_ = m.SetType(result.GetType())
		_ = m.SetSubject(result.GetSubject())
		_ = m.SetDataContentType(result.GetDataContentType())
		_ = m.SetDataSchema(result.GetDataSchema())
		//_ = msg.SetTime(time.UnixMilli(result.Time))

		m.SetContent(result.Content())

		result.ForEachAnnotation(func(k string, v string) {
			m.SetAnnotation(k, v)
		})

		return nil
	}

	return p, nil
}
