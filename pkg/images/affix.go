package images

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/stream"
	"github.com/google/go-containerregistry/pkg/v1/types"
)

//go:embed ast_openapi_schema.py
var script string

func Affix(baseRef string, dest string, newLayer *bytes.Buffer, predictorToParse string, commit string, auth authn.Authenticator) (string, error) {

	var base v1.Image
	var err error

	fmt.Fprintln(os.Stderr, "fetching metadata for", baseRef)

	start := time.Now()
	base, err = crane.Pull(baseRef, crane.WithAuth(auth))
	if err != nil {
		return "", fmt.Errorf("pulling %w", err)
	}
	fmt.Fprintln(os.Stderr, "pulling took", time.Since(start))

	// try to parse the predictor if it's provided
	if predictorToParse != "" {
		base, err = updatePredictor(base, predictorToParse)
		if err != nil {
			return "", fmt.Errorf("updating predictor: %w", err)
		}
	}

	if commit != "" {
		base, err = updateCommit(base, commit)
		if err != nil {
			return "", fmt.Errorf("updating predictor: %w", err)
		}
	}

	// FIXME(ja): find any YOLOs in the history and remove them?  We don't want to grow the history forever

	fmt.Fprintln(os.Stderr, "appending as new layer")

	start = time.Now()
	img, err := appendLayer(base, newLayer)
	if err != nil {
		return "", fmt.Errorf("appending %v: %w", newLayer, err)
	}
	fmt.Fprintln(os.Stderr, "appending took", time.Since(start))

	// --- pushing image
	start = time.Now()

	err = crane.Push(img, dest, crane.WithAuth(auth))
	if err != nil {
		return "", fmt.Errorf("pushing %s: %w", dest, err)
	}

	fmt.Fprintln(os.Stderr, "pushing took", time.Since(start))

	d, err := img.Digest()
	if err != nil {
		return "", err
	}
	image_id := fmt.Sprintf("%s@%s", dest, d)
	return image_id, nil
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

	fmt.Println("updating predictor to schema with length", len(schema))

	cfg.Config.Labels["org.cogmodel.openapi_schema"] = schema
	cfg.Config.Labels["run.cog.openapi_schema"] = schema

	return mutate.Config(img, cfg.Config)
}

func getSchema(predictorToParse string) (string, error) {
	cmd := exec.Command("python3", "-c", script, predictorToParse)

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("running ast_openapi_schema.py: %w", err)
	}

	schema := string(bytes.TrimSpace(out.Bytes()))
	return schema, nil
}
