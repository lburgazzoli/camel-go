package registry

import (
	"context"
	"io"
	"os"
	"path"
	"path/filepath"

	gregName "github.com/google/go-containerregistry/pkg/name"
	gregv1 "github.com/google/go-containerregistry/pkg/v1"
	gregRemote "github.com/google/go-containerregistry/pkg/v1/remote"
	gregTypes "github.com/google/go-containerregistry/pkg/v1/types"
)

const OCIContentWasm gregTypes.MediaType = "application/vnd.module.wasm.content.layer.v1+wasm"
const OCIAnnotationImageTitle string = "org.opencontainers.image.title"

// Pull download an image from the given registry and copy the content to a local temporary folder.
func Pull(ctx context.Context, imageName string) (string, error) {
	ref, err := gregName.ParseReference(imageName)
	if err != nil {
		return "", err
	}

	image, err := gregRemote.Image(
		ref,
		gregRemote.WithContext(ctx))

	if err != nil {
		return "", err
	}

	manifest, err := image.Manifest()
	if err != nil {
		return "", err
	}

	root, err := os.MkdirTemp("", "camel-")
	if err != nil {
		return "", err
	}

	for _, descriptor := range manifest.Layers {
		if descriptor.MediaType != OCIContentWasm {
			continue
		}

		tile, ok := descriptor.Annotations[OCIAnnotationImageTitle]
		if !ok {
			continue
		}

		l, err := image.LayerByDigest(descriptor.Digest)
		if err != nil {
			return "", err
		}

		err = copyLayer(l, path.Join(root, tile))
		if err != nil {
			return "", err
		}
	}

	return root, nil
}

// Blob read a blob from the registry.
func Blob(ctx context.Context, imageName string, layerName string) (io.ReadCloser, error) {
	ref, err := gregName.ParseReference(imageName)
	if err != nil {
		return nil, err
	}

	image, err := gregRemote.Image(
		ref,
		gregRemote.WithContext(ctx))

	if err != nil {
		return nil, err
	}

	manifest, err := image.Manifest()
	if err != nil {
		return nil, err
	}

	for _, descriptor := range manifest.Layers {
		if descriptor.MediaType != OCIContentWasm {
			continue
		}

		tile, ok := descriptor.Annotations[OCIAnnotationImageTitle]
		if !ok {
			continue
		}

		if layerName != tile {
			continue
		}

		l, err := image.LayerByDigest(descriptor.Digest)
		if err != nil {
			return nil, err
		}

		return l.Compressed()
	}

	return nil, nil
}

func copyLayer(layer gregv1.Layer, target string) error {
	parent := filepath.Dir(target)
	if parent != "" {
		if err := os.MkdirAll(parent, os.ModePerm); err != nil {
			return err
		}
	}

	out, err := os.Create(target)
	if err != nil {
		return err
	}

	defer func() {
		_ = out.Close()
	}()

	content, err := layer.Compressed()
	if err != nil {
		return err
	}

	defer func() {
		_ = content.Close()
	}()

	_, err = io.Copy(out, content)
	if err != nil {
		return err
	}

	return nil
}
