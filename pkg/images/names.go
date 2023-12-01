package images

import (
	"fmt"
	"strings"

	v1 "github.com/google/go-containerregistry/pkg/v1"
)

func EnsureRegistry(baseRef string) string {
	if !strings.Contains(baseRef, "r8.im") {
		return "r8.im/" + baseRef
	}
	return baseRef
}

func ImageId(baseRef string, img v1.Image) (string, error) {
	d, err := img.Digest()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s@%s", EnsureRegistry(baseRef), d), nil
}
