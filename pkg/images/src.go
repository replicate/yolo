package images

import (
	"regexp"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

// returns /src layers created by yolo and/or cog, oldest layers first
func GetSourceLayers(base v1.Image, cog bool, yolo bool) ([]v1.Layer, error) {
	var srcLayers []v1.Layer

	cfg, err := base.ConfigFile()
	if err != nil {
		return nil, err
	}

	layers, err := base.Layers()
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(` COPY .*\/src ?(# .*)?$`)
	idx := 0
	for _, h := range cfg.History {
		if h.EmptyLayer {
			continue
		}

		if yolo && h.CreatedBy == "cp . /src # yolo" {
			srcLayers = append(srcLayers, layers[idx])
		}
		if cog && (h.CreatedBy == "COPY . /src # buildkit" || re.MatchString(h.CreatedBy)) {
			srcLayers = append(srcLayers, layers[idx])
		}
		idx++
	}

	return srcLayers, nil
}
