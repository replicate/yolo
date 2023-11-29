package cli

import (
	"fmt"
	"os"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/replicate/yolo/pkg/auth"
	"github.com/replicate/yolo/pkg/images"
	"github.com/spf13/cobra"
)

func newCloneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "clone",
		Short:  "clone an existing image",
		Hidden: false,
		RunE:   cloneCommmand,
		Args:   cobra.ExactArgs(0),
	}

	cmd.Flags().StringVarP(&sToken, "token", "t", "", "replicate api token")
	cmd.Flags().StringVarP(&sRegistry, "registry", "r", "r8.im", "registry host")
	cmd.Flags().StringVarP(&baseRef, "base", "b", "", "base image reference.  examples: owner/model or r8.im/owner/model@sha256:hexdigest")
	cmd.Flags().StringVarP(&dest, "dest", "d", "", "destination image. examples: owner/model or r8.im/owner/model")
	cmd.MarkFlagRequired("base")
	cmd.MarkFlagRequired("dest")

	return cmd
}

func cloneCommmand(cmd *cobra.Command, args []string) error {
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
	dest = ensureRegistry(dest)
	image_id, err := images.Affix(baseRef, dest, nil, "", "", session)
	if err != nil {
		return err
	}
	fmt.Println(image_id)

	return nil
}
