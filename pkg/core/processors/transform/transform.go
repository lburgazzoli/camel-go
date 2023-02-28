////go:build steps_transform || steps_all

package transform

import (
	"context"
	"os"
	"path"
	"strings"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"

	"github.com/asynkron/protoactor-go/actor"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
	"github.com/lburgazzoli/camel-go/pkg/wasm"

	"oras.land/oras-go/v2/registry/remote"
)

const TAG = "transform"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Transform{
			Identity: uuid.New(),
		}
	}
}

type Transform struct {
	camel.Identifiable
	camel.WithOutputs

	Identity string `yaml:"id"`
	Language `yaml:",inline"`

	context   camel.Context
	processor *wasm.Processor
	runtime   *wasm.Runtime
}

type Language struct {
	Wasm *LanguageWasm `yaml:"wasm,omitempty"`
}

type LanguageWasm struct {
	Path  string `yaml:"path"`
	Image string `yaml:"image,omitempty"`
}

func (t *Transform) ID() string {
	return t.Identity
}

func (t *Transform) Reify(ctx context.Context, camelContext camel.Context) (string, error) {

	if t.Wasm == nil {
		return "", camelerrors.MissingParameterf("wasm", "failure processing %s", TAG)
	}
	if t.Wasm.Path == "" {
		return "", camelerrors.MissingParameterf("wasm.path", "failure processing %s", TAG)
	}

	rootPath := ""

	if t.Wasm.Image != "" {
		fp, err := t.pull(ctx)
		if err != nil {
			return "", err
		}

		rootPath = fp
	}

	defer func() {
		if rootPath != "" {
			_ = os.RemoveAll(rootPath)
		}
	}()

	f, err := os.Open(path.Join(rootPath, t.Wasm.Path))
	if err != nil {
		return "", err
	}

	defer func() { _ = f.Close() }()

	r, err := wasm.NewRuntime(ctx, wasm.Options{})
	if err != nil {
		return "", err
	}

	m, err := r.Load(ctx, f)
	if err != nil {
		return "", err
	}

	t.runtime = r
	t.context = camelContext
	t.processor = m

	return t.Identity, camelContext.Spawn(t)
}

func (t *Transform) Receive(c actor.Context) {
	msg, ok := c.Message().(camel.Message)
	if ok {
		out, err := t.processor.Process(context.Background(), msg)
		if err != nil {
			panic(err)
		}

		for _, pid := range t.Outputs() {
			if err := t.context.Send(pid, out); err != nil {
				panic(err)
			}
		}
	}
}

func (t *Transform) pull(ctx context.Context) (string, error) {
	repo := strings.SplitAfter(t.Language.Wasm.Image, ":")[0]
	repo = strings.TrimSuffix(repo, ":")

	tag := strings.SplitAfter(t.Language.Wasm.Image, ":")[1]

	r, err := remote.NewRepository(repo)
	if err != nil {
		return "", err
	}

	r.Client = &auth.Client{
		Client: retry.DefaultClient,
		Cache:  auth.DefaultCache,
	}

	f, err := os.MkdirTemp("", "camel-")
	if err != nil {
		return "", err
	}

	store, err := file.New(f)
	if err != nil {
		return "", err
	}

	if _, err = oras.Copy(ctx, r, tag, store, tag, oras.DefaultCopyOptions); err != nil {
		return "", err
	}

	return f, nil
}
