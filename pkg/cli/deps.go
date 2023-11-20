package cli

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/replicate/yolo/pkg/auth"
	"github.com/replicate/yolo/pkg/images"
	"github.com/spf13/cobra"
)

func newDepsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "deps",
		Short:  "install dependencies for existing image",
		Hidden: false,
		RunE:   depsCommmand,
		Args:   cobra.ExactArgs(1),
	}

	cmd.Flags().StringVarP(&sToken, "token", "t", "", "replicate api token")
	cmd.Flags().StringVarP(&sRegistry, "registry", "r", "r8.im", "registry host")

	return cmd
}

func depsCommmand(cmd *cobra.Command, args []string) error {
	baseRef := args[0]
	var session authn.Authenticator

	if sToken == "" {
		sToken = os.Getenv("REPLICATE_API_TOKEN")
	}

	if sToken == "" {
		sToken = os.Getenv("COG_TOKEN")
	}

	if sToken == "" {
		session = authn.Anonymous
	} else {
		u, err := auth.VerifyCogToken(sRegistry, sToken)
		if err != nil {
			fmt.Fprintln(os.Stderr, "authentication error, invalid token or registry host error")
			return err
		}
		session = authn.FromConfig(authn.AuthConfig{Username: u, Password: sToken})
	}

	baseRef = ensureRegistry(baseRef)
	config, err := images.Config(baseRef, session)
	if err != nil {
		return err
	}

	if len(config.Build.SystemPackages) > 0 {
		apt := "apt-get install -y"
		for _, pkg := range config.Build.SystemPackages {
			apt += " \"" + pkg + "\""
		}
		fmt.Println(apt)
	}

	if len(config.Build.PythonPackages) > 0 {
		pip := "python3 -m pip install --no-cache-dir"
		for _, pkg := range config.Build.PythonPackages {
			pip += " \"" + pkg + "\""
		}
		fmt.Println(pip)
	}

	for _, run := range config.Build.Run {
		fmt.Println(run.Command)
	}

	return nil
}
