package images

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

func Clone(baseRef string, dest string, session authn.Authenticator) (string, error) {

	fmt.Fprintln(os.Stderr, "fetching metadata for", baseRef)
	base, err := crane.Pull(baseRef, crane.WithAuth(session))
	if err != nil {
		return "", fmt.Errorf("pulling %w", err)
	}

	// as r8.im fails if you push the same image to a second location,
	// we can work around this by adding a label "cloned" with the current time
	cfg, err := base.ConfigFile()
	if err != nil {
		return "", fmt.Errorf("getting config file: %w", err)
	}
	cfg.Config.Labels["cloned"] = time.Now().String()

	img, err := mutate.ConfigFile(base, cfg)
	if err != nil {
		return "", fmt.Errorf("mutating config file: %w", err)
	}

	start := time.Now()
	err = crane.Push(img, dest, crane.WithAuth(session))
	if err != nil {
		return "", fmt.Errorf("pushing %s: %w", dest, err)
	}
	fmt.Fprintln(os.Stderr, "pushing took", time.Since(start))

	return ImageId(dest, img)
}
