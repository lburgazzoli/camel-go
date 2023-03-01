package registry

import (
	"context"
	"os"
	"strings"

	"oras.land/oras-go/v2"
	"oras.land/oras-go/v2/content/file"
	"oras.land/oras-go/v2/registry/remote"
	"oras.land/oras-go/v2/registry/remote/auth"
	"oras.land/oras-go/v2/registry/remote/retry"
)

// Pull download an image from the given registry and copy the content to a local temporary folder.
func Pull(ctx context.Context, image string) (string, error) {
	repo := strings.SplitAfter(image, ":")[0]
	repo = strings.TrimSuffix(repo, ":")

	tag := strings.SplitAfter(image, ":")[1]

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
