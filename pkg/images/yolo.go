package images

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/stream"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

func Yolo(baseRef string, dest string, files []LayerFile, predictorToParse string, commit string, session authn.Authenticator) (string, error) {
	fmt.Fprintln(os.Stderr, "fetching metadata for", baseRef)
	base, err := crane.Pull(baseRef, crane.WithAuth(session))
	if err != nil {
		return "", fmt.Errorf("pulling %w", err)
	}

	yoloLess, err := removeYolo(base)
	if err != nil {
		return "", fmt.Errorf("removing existing yolo layers: %w", err)
	}

	// try to parse the predictor if it's provided
	if predictorToParse != "" {
		yoloLess, err = updatePredictor(yoloLess, predictorToParse)
		if err != nil {
			return "", fmt.Errorf("updating predictor: %w", err)
		}
	}

	if commit != "" {
		yoloLess, err = updateCommit(yoloLess, commit)
		if err != nil {
			return "", fmt.Errorf("updating commit: %w", err)
		}
	}

	fmt.Fprintln(os.Stderr, "appending as new layer")
	var img v1.Image

	yoloLayers, err := GetSourceLayers(base, false, true)
	if err != nil {
		return "", fmt.Errorf("getting source layers: %w", err)
	}

	newLayer, err := MakeTar(files, yoloLayers)
	if err != nil {
		return "", fmt.Errorf("making tar: %w", err)
	}

	img, err = appendLayer(yoloLess, newLayer)
	if err != nil {
		return "", fmt.Errorf("appending %v: %w", newLayer, err)
	}

	// --- pushing image
	start := time.Now()
	err = crane.Push(img, dest, crane.WithAuth(session))
	if err != nil {
		return "", fmt.Errorf("pushing %s: %w", dest, err)
	}
	fmt.Fprintln(os.Stderr, "pushing took", time.Since(start))

	return ImageId(dest, img)
}

// All of this code is from pkg/v1/mutate - so we can add history and use a tarball
func appendLayer(base v1.Image, tarball *bytes.Buffer) (v1.Image, error) {
	baseMediaType, err := base.MediaType()
	if err != nil {
		return nil, fmt.Errorf("getting base image media type: %w", err)
	}

	layerType := types.DockerLayer
	if baseMediaType == types.OCIManifestSchema1 {
		layerType = types.OCILayer
	}

	layer := stream.NewLayer(io.NopCloser(tarball), stream.WithMediaType(layerType))

	history := v1.History{
		CreatedBy: "cp . /src # yolo",
		Created:   v1.Time{Time: time.Now()},
		Author:    "yolo",
		Comment:   "",
	}

	return mutate.Append(base, mutate.Addendum{Layer: layer, History: history})
}

func updateCommit(img v1.Image, commit string) (v1.Image, error) {
	cfg, err := img.ConfigFile()
	if err != nil {
		return nil, err
	}

	cfg.Config.Labels["org.opencontainers.image.revision"] = commit

	return mutate.Config(img, cfg.Config)
}

func updatePredictor(img v1.Image, predictorToParse string) (v1.Image, error) {
	schema, err := getSchema(predictorToParse)
	if err != nil {
		return nil, fmt.Errorf("getting schema: %w", err)
	}

	cfg, err := img.ConfigFile()
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(os.Stderr, "updating predictor to schema with length", len(schema))

	cfg.Config.Labels["org.cogmodel.openapi_schema"] = schema
	cfg.Config.Labels["run.cog.openapi_schema"] = schema

	return mutate.Config(img, cfg.Config)
}

// we need to remove any existing yolo layers before adding more... otherwise
// we'll end up with a bunch of yolo layers
func removeYolo(orig v1.Image) (v1.Image, error) {
	layers, err := orig.Layers()
	if err != nil {
		return nil, fmt.Errorf("failed to get layers for original: %w", err)
	}

	config, err := orig.ConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get config for original: %w", err)
	}

	rebasedImage, err := mutate.Config(empty.Image, *config.Config.DeepCopy())
	if err != nil {
		return nil, fmt.Errorf("failed to create empty image with original config: %w", err)
	}

	idx := 0
	for _, h := range config.History {

		if h.CreatedBy != "cp . /src # yolo" {
			add := mutate.Addendum{
				Layer:   nil,
				History: h,
			}
			if !h.EmptyLayer {
				add.Layer = layers[idx]
			}

			fmt.Println("adding layer", add.Layer, "with history", h)
			rebasedImage, err = mutate.Append(rebasedImage, add)
			if err != nil {
				return nil, fmt.Errorf("failed to add layer: %w", err)
			}
		}

		if !h.EmptyLayer {
			idx++
		}
	}

	return rebasedImage, nil
}
