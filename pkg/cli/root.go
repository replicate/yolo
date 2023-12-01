package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/logs"
	"github.com/replicate/yolo/pkg/auth"
	"github.com/replicate/yolo/pkg/version"
	"github.com/spf13/cobra"
)

func NewRootCommand() (*cobra.Command, error) {
	rootCmd := cobra.Command{
		Use:           "yolo",
		Short:         "remix the web",
		Version:       version.GetVersion(),
		SilenceErrors: true,
	}

	rootCmd.AddCommand(
		newCloneCommand(),
		newFetchCommand(),
		newPushCommand(),
	)
	logs.Warn = log.New(os.Stderr, "gcr WARN: ", log.LstdFlags)
	logs.Progress = log.New(os.Stderr, "gcr: ", log.LstdFlags)

	return &rootCmd, nil
}

func authenticate() authn.Authenticator {
	if sToken == "" {
		sToken = os.Getenv("REPLICATE_API_TOKEN")
	}

	if sToken == "" {
		sToken = os.Getenv("COG_TOKEN")
	}

	if sToken != "" {
		u, err := auth.VerifyCogToken(sRegistry, sToken)
		if err != nil {
			fmt.Fprintln(os.Stderr, "authentication error, invalid token or registry host error")
			return nil
		}
		return authn.FromConfig(authn.AuthConfig{Username: u, Password: sToken})
	}

	return authn.Anonymous
}

func ensureRegistry(baseRef string) string {
	if !strings.Contains(baseRef, sRegistry) {
		return sRegistry + "/" + baseRef
	}
	return baseRef
}
