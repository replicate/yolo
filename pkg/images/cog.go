package images

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
)

// {"build":{"gpu":true,"python_version":"3.9","python_packages":["diffusers==0.19.3","torch==2.0.1","transformers==4.31.0","invisible-watermark==0.2.0","accelerate==0.21.0","pandas==2.0.3","torchvision==0.15.2","numpy==1.25.1","pandas==2.0.3","fire==0.5.0","opencv-python\u003e=4.1.0.25","mediapipe==0.10.2"],"run":[{"command":"curl -o /usr/local/bin/pget -L \"https://github.com/replicate/pget/releases/download/v0.0.3/pget\" \u0026\u0026 chmod +x /usr/local/bin/pget"},{"command":"wget http://thegiflibrary.tumblr.com/post/11565547760 -O face_landmarker_v2_with_blendshapes.task -q https://storage.googleapis.com/mediapipe-models/face_landmarker/face_landmarker/float16/1/face_landmarker.task"}],"system_packages":["libgl1-mesa-glx","ffmpeg","libsm6","libxext6","wget"],"cuda":"11.8","cudnn":"8"},"predict":"predict.py:Predictor","train":"train.py:train"}

type RunItem struct {
	Command string `json:"command,omitempty" yaml:"command"`
	Mounts  []struct {
		Type   string `json:"type,omitempty" yaml:"type"`
		ID     string `json:"id,omitempty" yaml:"id"`
		Target string `json:"target,omitempty" yaml:"target"`
	} `json:"mounts,omitempty" yaml:"mounts"`
}

type Build struct {
	GPU            bool      `json:"gpu"`
	PythonVersion  string    `json:"python_version"`
	PythonPackages []string  `json:"python_packages"`
	SystemPackages []string  `json:"system_packages"`
	Run            []RunItem `json:"run"`
	Cuda           string    `json:"cuda"`
	Cudnn          string    `json:"cudnn"`
}

type Schema struct {
	Build   Build  `json:"build"`
	Predict string `json:"predict"`
}

func Config(imageName string, auth authn.Authenticator) (Schema, error) {
	var schema Schema

	image, err := crane.Pull(imageName, crane.WithAuth(auth))
	if err != nil {
		return schema, err
	}

	cfg, err := image.ConfigFile()
	if err != nil {
		return schema, err
	}

	schemaString := cfg.Config.Labels["org.cogmodel.config"]
	if schemaString == "" {
		schemaString = cfg.Config.Labels["run.cog.config"]
	}
	if schemaString == "" {
		return schema, fmt.Errorf("no cog config found")
	}

	// fmt.Println(schemaString)

	err = json.Unmarshal([]byte(schemaString), &schema)
	if err != nil {
		return schema, fmt.Errorf("error unmarshalling schema: %v", err)
	}

	return schema, nil
}
