package cli

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/replicate/yolo/pkg/auth"
	"github.com/replicate/yolo/pkg/images"
	"github.com/spf13/cobra"
)

func newFetchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "fetch",
		Short:  "fetch an existing image",
		Hidden: false,
		RunE:   fetchCommmand,
		Args:   cobra.ExactArgs(1),
	}

	cmd.Flags().StringVarP(&sToken, "token", "t", "", "replicate cog token")
	cmd.Flags().StringVarP(&sRegistry, "registry", "r", "r8.im", "registry host")
	cmd.Flags().StringVarP(&baseRef, "base", "b", "", "base image reference - include tag: r8.im/username/modelname@sha256:hexdigest")
	cmd.MarkFlagRequired("base")

	return cmd
}

func fetchCommmand(cmd *cobra.Command, args []string) error {
	dest := args[0]
	if sToken == "" {
		sToken = os.Getenv("COG_TOKEN")
	}

	u, err := auth.VerifyCogToken(sRegistry, sToken)
	if err != nil {
		fmt.Fprintln(os.Stderr, "authentication error, invalid token or registry host error")
		return err
	}
	auth := authn.FromConfig(authn.AuthConfig{Username: u, Password: sToken})

	err = images.Extract(baseRef, dest, auth)
	if err != nil {
		return err
	}

	return nil
}
